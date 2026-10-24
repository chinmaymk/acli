package jira

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEscapeNonASCII_PlainASCIIPassthrough(t *testing.T) {
	in := []byte(`{"a":"hello world","b":42}`)
	got := string(escapeNonASCII(in))
	if got != string(in) {
		t.Errorf("plain ASCII should pass through unchanged\n got: %s\nwant: %s", got, in)
	}
}

func TestEscapeNonASCII_BMPCharsEscaped(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"latin accent é", `{"x":"café"}`, `{"x":"caf\u00e9"}`},
		{"cyrillic", `{"x":"Привет"}`, `{"x":"\u041f\u0440\u0438\u0432\u0435\u0442"}`},
		{"japanese", `{"x":"こんにちは"}`, `{"x":"\u3053\u3093\u306b\u3061\u306f"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(escapeNonASCII([]byte(tc.in)))
			if got != tc.want {
				t.Errorf("escapeNonASCII(%q)\n got: %s\nwant: %s", tc.in, got, tc.want)
			}
		})
	}
}

func TestEscapeNonASCII_SurrogatePairsForSupplementaryPlane(t *testing.T) {
	// Emoji 😀 is U+1F600 → must be encoded as UTF-16 surrogate pair \uD83D\uDE00
	// because that's the only way to represent supplementary-plane chars in
	// JSON's \uXXXX escape syntax.
	in := []byte(`{"x":"😀"}`)
	want := `{"x":"\ud83d\ude00"}`
	got := string(escapeNonASCII(in))
	if got != want {
		t.Errorf("emoji surrogate pair encoding\n got: %s\nwant: %s", got, want)
	}
}

func TestEscapeNonASCII_DoesNotTouchKeysOutsideStrings(t *testing.T) {
	// The escaper is string-context-aware: bytes outside double-quoted regions
	// must not be touched. Numbers, structural chars, whitespace stay intact.
	in := []byte(`{"key":  42 ,  "list":[1,2,3]  }`)
	got := string(escapeNonASCII(in))
	if got != string(in) {
		t.Errorf("non-string bytes should be untouched\n got: %s\nwant: %s", got, in)
	}
}

func TestEscapeNonASCII_RespectsEscapedQuotes(t *testing.T) {
	// `\"` inside a string must NOT close the string context.
	in := []byte(`{"x":"he said \"héllo\""}`)
	want := `{"x":"he said \"h\u00e9llo\""}`
	got := string(escapeNonASCII(in))
	if got != want {
		t.Errorf("escaped quotes should not exit string context\n got: %s\nwant: %s", got, want)
	}
}

func TestEscapeNonASCII_RespectsEscapedBackslash(t *testing.T) {
	// A literal `\\` is two backslashes — the second one must NOT be treated
	// as an escape character that would consume the following byte.
	in := []byte(`{"path":"C:\\Users\\café"}`)
	want := `{"path":"C:\\Users\\caf\u00e9"}`
	got := string(escapeNonASCII(in))
	if got != want {
		t.Errorf("escaped backslash handling\n got: %s\nwant: %s", got, want)
	}
}

func TestEscapeNonASCII_OutputRemainsValidJSON(t *testing.T) {
	original := map[string]interface{}{
		"latin":   "résumé",
		"emoji":   "🚀 ship it",
		"cyril":   "тест",
		"escaped": `quote: " backslash: \ tab: ` + "\t",
		"nested":  map[string]interface{}{"deep": "déjà vu"},
		"arr":     []string{"αβγ", "abc"},
	}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	escaped := escapeNonASCII(encoded)

	var roundTripped map[string]interface{}
	if err := json.Unmarshal(escaped, &roundTripped); err != nil {
		t.Fatalf("escapeNonASCII output is not valid JSON: %v\noutput: %s", err, escaped)
	}
	if roundTripped["latin"] != "résumé" {
		t.Errorf("round-trip preserves accents: got %v", roundTripped["latin"])
	}
	if roundTripped["cyril"] != "тест" {
		t.Errorf("round-trip preserves cyrillic: got %v", roundTripped["cyril"])
	}
	if roundTripped["emoji"] != "🚀 ship it" {
		t.Errorf("round-trip preserves emoji: got %v", roundTripped["emoji"])
	}
}

func TestNewRequest_EncodesNonASCIIAndDoesNotEscapeHTML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s := string(body)
		// non-ASCII chars must be \uXXXX escaped on the wire (Atlassian quirk
		// — some endpoints reject raw UTF-8 in JSON bodies).
		if strings.Contains(s, "café") {
			t.Errorf("body should escape non-ASCII; got raw UTF-8: %s", s)
		}
		if !strings.Contains(s, `\u00e9`) {
			t.Errorf("body should contain \\u00e9 for é; got: %s", s)
		}
		// HTML chars (<, >, &) must NOT be Go-default escaped (\u003c etc.) —
		// json.NewEncoder.SetEscapeHTML(false) prevents that.
		if strings.Contains(s, `\u003c`) {
			t.Errorf("HTML escape should be disabled; got: %s", s)
		}
		if !strings.Contains(s, "<b>") {
			t.Errorf("HTML chars should be passed through literally; got: %s", s)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	body := map[string]string{
		"name":  "café",
		"html":  "<b>bold</b> & <i>italic</i>",
		"emoji": "🚀",
	}
	if err := c.Post("/test", body, nil); err != nil {
		t.Fatalf("Post failed: %v", err)
	}
}

func TestNewRequest_ContentTypeIncludesCharset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json; charset=utf-8" {
			t.Errorf("expected charset=utf-8 content type, got: %s", ct)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	if err := c.Post("/test", map[string]string{"x": "y"}, nil); err != nil {
		t.Fatalf("Post failed: %v", err)
	}
}

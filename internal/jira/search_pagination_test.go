package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchJQLPagination(t *testing.T) {
	const first = `{"issues":[{"key":"TEST-1"}],"isLast":false,"nextPageToken":"a+/="}`
	const empty = `{"issues":[],"nextPageToken":"b"}`
	const last = `{"issues":[{"key":"TEST-2"}],"isLast":true}`
	tests := []struct {
		name      string
		all       bool
		pages     []string
		tokens    []string
		wantCount int
		wantErr   string
	}{
		{"one page by default", false, []string{first}, []string{""}, 1, ""},
		{"follow cursors through empty page", true, []string{first, empty, last}, []string{"", "a+/=", "b"}, 2, ""},
		{"terminal without isLast", true, []string{first, `{"issues":[]}`}, []string{"", "a+/="}, 1, ""},
		{"repeated token", true, []string{first, empty, first}, []string{"", "a+/=", "b"}, 0, "repeated"},
		{"missing token", true, []string{`{"issues":[],"isLast":false}`}, []string{""}, 0, "without a nextPageToken"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if n >= len(tc.pages) {
					http.Error(w, "unexpected extra request", http.StatusBadRequest)
					return
				}
				q := r.URL.Query()
				if r.URL.Path != "/rest/api/3/search/jql" || q.Get("jql") != "project = TEST" ||
					q.Get("fields") != "summary,status" || q.Get("maxResults") != "1" ||
					q.Get("nextPageToken") != tc.tokens[n] || q.Has("startAt") {
					t.Errorf("unexpected request: %s", r.URL)
				}
				_, _ = fmt.Fprint(w, tc.pages[n])
				n++
			}))
			got, err := newTestClient(srv).SearchJQL("project = TEST", 0, 1, []string{"summary", "status"}, nil, tc.all)
			srv.Close()
			if n != len(tc.pages) {
				t.Fatalf("got %d requests, want %d", n, len(tc.pages))
			}
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) || got != nil {
					t.Fatalf("result = %v, error = %v", got, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Issues) != tc.wantCount || got.Issues[0].Key != "TEST-1" {
				t.Fatalf("issues = %+v", got.Issues)
			}
			if (!tc.all && got.NextPageToken != "a+/=") || (tc.all && got.NextPageToken != "") {
				t.Fatalf("unexpected continuation token %q", got.NextPageToken)
			}
		})
	}
}

func TestSearchJQLWithLegacyOffsets(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("startAt") {
		case "5":
			_, _ = fmt.Fprint(w, `{"total":7,"issues":[{"key":"TEST-6"}]}`)
		case "6":
			_, _ = fmt.Fprint(w, `{"total":7,"issues":[{"key":"TEST-7"}]}`)
		default:
			http.Error(w, "unexpected offset", http.StatusBadRequest)
		}
	}))
	defer srv.Close()
	got, err := newTestClient(srv).SearchJQL("project = TEST", 5, 1, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Issues) != 2 || got.Issues[1].Key != "TEST-7" {
		t.Fatalf("unexpected issues: %+v", got.Issues)
	}
}

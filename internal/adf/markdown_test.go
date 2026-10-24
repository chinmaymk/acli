package adf

import (
	"reflect"
	"testing"
)

func TestMarkdownToADF_HeadingLevels1Through6(t *testing.T) {
	cases := []struct {
		md            string
		expectedLevel int
	}{
		{"# h1", 1},
		{"## h2", 2},
		{"### h3", 3},
		{"#### h4", 4},
		{"##### h5", 5},
		{"###### h6", 6},
	}
	for _, tc := range cases {
		t.Run(tc.md, func(t *testing.T) {
			doc := markdownToADF(tc.md)
			content, _ := doc["content"].([]interface{})
			if len(content) != 1 {
				t.Fatalf("expected 1 block, got %d: %#v", len(content), content)
			}
			node, _ := content[0].(map[string]interface{})
			if node["type"] != "heading" {
				t.Fatalf("expected heading, got %v", node["type"])
			}
			attrs, _ := node["attrs"].(map[string]interface{})
			if attrs["level"] != tc.expectedLevel {
				t.Errorf("level mismatch: got %v, want %d", attrs["level"], tc.expectedLevel)
			}
		})
	}
}

func TestMarkdownToADF_SevenHashesIsParagraphNotHeading(t *testing.T) {
	// `####### x` has 7 `#`s — beyond Markdown's h6 cap, should fall through
	// to paragraph rendering rather than being treated as some "h7".
	doc := markdownToADF("####### too deep")
	content, _ := doc["content"].([]interface{})
	if len(content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(content))
	}
	node, _ := content[0].(map[string]interface{})
	if node["type"] == "heading" {
		t.Errorf("7-hash should NOT be a heading; got: %#v", node)
	}
}

func TestMarkdownToADF_HashWithoutSpaceIsParagraph(t *testing.T) {
	// `#tag` (no space) must be a paragraph, not a heading — common in
	// hashtag-style writing inside ticket bodies.
	doc := markdownToADF("#tag")
	content, _ := doc["content"].([]interface{})
	if len(content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(content))
	}
	node, _ := content[0].(map[string]interface{})
	if node["type"] == "heading" {
		t.Errorf("no-space `#tag` should NOT be a heading; got: %#v", node)
	}
}

func TestMarkdownToADF_HeadingContentIsParsedInline(t *testing.T) {
	doc := markdownToADF("### **bold** in heading")
	content, _ := doc["content"].([]interface{})
	node, _ := content[0].(map[string]interface{})
	inline, _ := node["content"].([]interface{})
	if len(inline) == 0 {
		t.Fatalf("heading should have parsed inline content; got: %#v", node)
	}
	first, _ := inline[0].(map[string]interface{})
	marks, _ := first["marks"].([]interface{})
	if len(marks) == 0 {
		t.Errorf("expected strong mark on heading inline; got: %#v", first)
	}
	if !reflect.DeepEqual(marks[0], map[string]interface{}{"type": "strong"}) {
		t.Errorf("expected strong mark; got %#v", marks[0])
	}
}

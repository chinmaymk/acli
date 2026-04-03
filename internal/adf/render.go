// Package adf provides rendering of Atlassian Document Format (ADF) to plain text.
package adf

import (
	"fmt"
	"strings"
)

// Render converts an ADF document to plain text.
// It accepts interface{} to match the Description/Body/Comment field types
// used in the Jira API types (which are interface{} since they can be
// either a string or an ADF document map).
func Render(doc interface{}) string {
	if doc == nil {
		return ""
	}
	switch d := doc.(type) {
	case string:
		return d
	case map[string]interface{}:
		result := renderNode(d)
		return strings.TrimRight(result, "\n")
	default:
		return fmt.Sprintf("%v", doc)
	}
}

func renderNode(node map[string]interface{}) string {
	nodeType, _ := node["type"].(string)

	switch nodeType {
	case "doc", "blockquote", "panel", "expand":
		return renderChildren(node)
	case "paragraph", "heading", "codeBlock":
		return renderChildren(node) + "\n"
	case "bulletList":
		return renderBulletList(node)
	case "orderedList":
		return renderOrderedList(node)
	case "listItem":
		return renderListItem(node)
	case "text":
		text, _ := node["text"].(string)
		return text
	case "hardBreak":
		return "\n"
	case "rule":
		return "---\n"
	case "table":
		return renderChildren(node)
	case "tableRow":
		return renderTableRow(node)
	case "tableHeader", "tableCell":
		return strings.TrimRight(renderChildren(node), "\n")
	case "mention":
		return nodeAttr(node, "text", "@unknown")
	case "inlineCard":
		return nodeAttr(node, "url", "")
	case "emoji":
		if text := nodeAttr(node, "text", ""); text != "" {
			return text
		}
		return nodeAttr(node, "shortName", "")
	case "date":
		return nodeAttr(node, "timestamp", "")
	case "status":
		return nodeAttr(node, "text", "")
	case "mediaSingle", "mediaInline", "media", "mediaGroup":
		return ""
	default:
		return renderChildren(node)
	}
}

// renderChildren renders all child nodes in the "content" array.
func renderChildren(node map[string]interface{}) string {
	content, ok := node["content"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, child := range content {
		if m, ok := child.(map[string]interface{}); ok {
			sb.WriteString(renderNode(m))
		}
	}
	return sb.String()
}

func renderBulletList(node map[string]interface{}) string {
	content, ok := node["content"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, child := range content {
		if m, ok := child.(map[string]interface{}); ok {
			text := strings.TrimRight(renderNode(m), "\n")
			sb.WriteString(text)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func renderOrderedList(node map[string]interface{}) string {
	content, ok := node["content"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for i, child := range content {
		if m, ok := child.(map[string]interface{}); ok {
			text := strings.TrimRight(renderNode(m), "\n")
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, text))
		}
	}
	return sb.String()
}

func renderListItem(node map[string]interface{}) string {
	content, ok := node["content"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	first := true
	for _, child := range content {
		if m, ok := child.(map[string]interface{}); ok {
			text := strings.TrimRight(renderNode(m), "\n")
			if text != "" {
				if !first {
					sb.WriteString("\n")
				}
				sb.WriteString(text)
				first = false
			}
		}
	}
	return sb.String()
}

func renderTableRow(node map[string]interface{}) string {
	content, ok := node["content"].([]interface{})
	if !ok {
		return ""
	}
	var cells []string
	for _, child := range content {
		if m, ok := child.(map[string]interface{}); ok {
			cells = append(cells, renderNode(m))
		}
	}
	return strings.Join(cells, "\t") + "\n"
}

// nodeAttr extracts a string attribute from node["attrs"][key],
// returning fallback if not found.
func nodeAttr(node map[string]interface{}, key, fallback string) string {
	attrs, ok := node["attrs"].(map[string]interface{})
	if !ok {
		return fallback
	}
	if val, ok := attrs[key].(string); ok {
		return val
	}
	return fallback
}

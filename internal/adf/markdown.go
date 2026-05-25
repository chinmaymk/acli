package adf

import (
	"strings"
	"unicode/utf8"
)

func markdownToADF(markdown string) map[string]interface{} {
	lines := strings.Split(markdown, "\n")
	content := []interface{}{}
	var codeBlock []string
	codeLang := ""
	inCode := false

	flushParagraph := func(text string) {
		if text == "" {
			return
		}
		inlineContent := parseInlineMarks(text)
		if len(inlineContent) > 0 {
			content = append(content, map[string]interface{}{
				"type":    "paragraph",
				"content": inlineContent,
			})
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "```") {
			if inCode {
				content = append(content, map[string]interface{}{
					"type": "codeBlock",
					"attrs": map[string]interface{}{
						"language": codeLang,
					},
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": strings.Join(codeBlock, "\n"),
						},
					},
				})
				codeBlock = nil
				inCode = false
			} else {
				inCode = true
				codeLang = strings.TrimPrefix(line, "```")
				codeBlock = []string{}
			}
			continue
		}

		if inCode {
			codeBlock = append(codeBlock, line)
			continue
		}

		if line == "---" || line == "----" {
			content = append(content, map[string]interface{}{"type": "rule"})
			continue
		}

		if len(line) > 0 && line[0] == '#' {
			level := 0
			for level < len(line) && line[level] == '#' {
				level++
			}
			if level >= 1 && level <= 6 && level < len(line) && line[level] == ' ' {
				content = append(content, heading(level, line[level+1:]))
				continue
			}
		}

		if strings.HasPrefix(line, "| ") && strings.HasSuffix(line, "|") {
			tableLines := []string{line}
			for i+1 < len(lines) && strings.HasPrefix(lines[i+1], "|") {
				i++
				tableLines = append(tableLines, lines[i])
			}
			table := parseTable(tableLines)
			if table != nil {
				content = append(content, table)
			}
			continue
		}

		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "1. ") {
			listLines := []string{line}
			for i+1 < len(lines) && (strings.HasPrefix(lines[i+1], "- ") || strings.HasPrefix(lines[i+1], "* ") || strings.HasPrefix(lines[i+1], "1. ") || (len(lines[i+1]) > 2 && lines[i+1][0] >= '0' && lines[i+1][0] <= '9' && strings.Contains(lines[i+1][:3], "."))) {
				i++
				listLines = append(listLines, lines[i])
			}
			ordered := strings.HasPrefix(listLines[0], "1.") || (len(listLines[0]) > 2 && listLines[0][0] >= '0' && listLines[0][0] <= '9')
			list := parseList(listLines, ordered)
			if list != nil {
				content = append(content, list)
			}
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		flushParagraph(line)
	}

	if inCode && len(codeBlock) > 0 {
		content = append(content, map[string]interface{}{
			"type": "codeBlock",
			"attrs": map[string]interface{}{
				"language": codeLang,
			},
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": strings.Join(codeBlock, "\n"),
				},
			},
		})
	}

	return map[string]interface{}{
		"type":    "doc",
		"version": 1,
		"content": content,
	}
}

func heading(level int, text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "heading",
		"attrs": map[string]interface{}{
			"level": level,
		},
		"content": parseInlineMarks(text),
	}
}

func parseInlineMarks(text string) []interface{} {
	result := []interface{}{}
	if text == "" {
		return result
	}

	i := 0
	current := ""
	for i < len(text) {
		if i < len(text)-1 && text[i] == '*' && text[i+1] == '*' {
			if current != "" {
				result = append(result, textNode(current, nil))
				current = ""
			}
			end := strings.Index(text[i+2:], "**")
			if end >= 0 {
				result = append(result, textNode(text[i+2:i+2+end], []map[string]interface{}{{"type": "strong"}}))
				i = i + 2 + end + 2
				continue
			}
		}
		if text[i] == '`' {
			if current != "" {
				result = append(result, textNode(current, nil))
				current = ""
			}
			end := strings.Index(text[i+1:], "`")
			if end >= 0 {
				result = append(result, textNode(text[i+1:i+1+end], []map[string]interface{}{{"type": "code"}}))
				i = i + 1 + end + 1
				continue
			}
		}
		if text[i] == '[' {
			linkEnd := strings.Index(text[i:], "](")
			if linkEnd >= 0 {
				urlEnd := strings.Index(text[i+linkEnd+2:], ")")
				if urlEnd >= 0 {
					if current != "" {
						result = append(result, textNode(current, nil))
						current = ""
					}
					linkText := text[i+1 : i+linkEnd]
					url := text[i+linkEnd+2 : i+linkEnd+2+urlEnd]
					result = append(result, map[string]interface{}{
						"type": "text",
						"text": linkText,
						"marks": []interface{}{
							map[string]interface{}{
								"type":  "link",
								"attrs": map[string]interface{}{"href": url},
							},
						},
					})
					i = i + linkEnd + 2 + urlEnd + 1
					continue
				}
			}
		}
		_, size := utf8.DecodeRuneInString(text[i:])
		current += text[i : i+size]
		i += size
	}

	if current != "" {
		result = append(result, textNode(current, nil))
	}
	return result
}

func textNode(text string, marks []map[string]interface{}) map[string]interface{} {
	node := map[string]interface{}{
		"type": "text",
		"text": text,
	}
	if len(marks) > 0 {
		m := []interface{}{}
		for _, mark := range marks {
			m = append(m, mark)
		}
		node["marks"] = m
	}
	return node
}

func parseTable(lines []string) map[string]interface{} {
	if len(lines) < 2 {
		return nil
	}
	rows := []interface{}{}
	for i, line := range lines {
		line = strings.TrimPrefix(line, "|")
		line = strings.TrimSuffix(line, "|")
		cells := strings.Split(line, "|")

		if i == 1 {
			allDashes := true
			for _, c := range cells {
				trimmed := strings.TrimSpace(c)
				if trimmed != "" && !strings.HasPrefix(trimmed, "---") && !strings.HasPrefix(trimmed, ":--") && !strings.HasPrefix(trimmed, "--:") {
					allDashes = false
				}
			}
			if allDashes {
				continue
			}
		}

		cellType := "tableCell"
		if i == 0 {
			cellType = "tableHeader"
		}
		rowCells := []interface{}{}
		for _, c := range cells {
			rowCells = append(rowCells, map[string]interface{}{
				"type": cellType,
				"content": []interface{}{
					map[string]interface{}{
						"type":    "paragraph",
						"content": parseInlineMarks(strings.TrimSpace(c)),
					},
				},
			})
		}
		rows = append(rows, map[string]interface{}{
			"type":    "tableRow",
			"content": rowCells,
		})
	}
	return map[string]interface{}{
		"type":    "table",
		"content": rows,
	}
}

func parseList(lines []string, ordered bool) map[string]interface{} {
	listType := "bulletList"
	if ordered {
		listType = "orderedList"
	}
	items := []interface{}{}
	for _, line := range lines {
		text := line
		if strings.HasPrefix(text, "- ") {
			text = strings.TrimPrefix(text, "- ")
		} else if strings.HasPrefix(text, "* ") {
			text = strings.TrimPrefix(text, "* ")
		} else {
			dotIdx := strings.Index(text, ". ")
			if dotIdx > 0 && dotIdx < 4 {
				text = text[dotIdx+2:]
			}
		}
		items = append(items, map[string]interface{}{
			"type": "listItem",
			"content": []interface{}{
				map[string]interface{}{
					"type":    "paragraph",
					"content": parseInlineMarks(text),
				},
			},
		})
	}
	return map[string]interface{}{
		"type":    listType,
		"content": items,
	}
}

// MarkdownToADF converts markdown text to Jira ADF document structure
func MarkdownToADF(markdown string) map[string]interface{} {
	return markdownToADF(markdown)
}

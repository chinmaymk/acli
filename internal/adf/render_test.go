package adf

import "testing"

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "string input",
			input:    "plain text description",
			expected: "plain text description",
		},
		{
			name:     "non-map non-string input",
			input:    42,
			expected: "42",
		},
		{
			name: "simple paragraph",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "A simple widget description.",
							},
						},
					},
				},
			},
			expected: "A simple widget description.",
		},
		{
			name: "headings and paragraphs",
			// Modeled on a feature spec: heading/paragraph pairs
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "heading",
						"attrs": map[string]interface{}{
							"level": 2,
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Overview",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Users need a centralized dashboard for managing all widgets.",
							},
						},
					},
					map[string]interface{}{
						"type": "heading",
						"attrs": map[string]interface{}{
							"level": 2,
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Motivation",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Currently each team manages widgets independently with no shared visibility.",
							},
						},
					},
					map[string]interface{}{
						"type": "heading",
						"attrs": map[string]interface{}{
							"level": 2,
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Proposed Solution",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "A single API endpoint to query and update widget status.",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "A web dashboard backed by the new API.",
							},
						},
					},
				},
			},
			expected: "Overview\n" +
				"Users need a centralized dashboard for managing all widgets.\n" +
				"Motivation\n" +
				"Currently each team manages widgets independently with no shared visibility.\n" +
				"Proposed Solution\n" +
				"A single API endpoint to query and update widget status.\n" +
				"A web dashboard backed by the new API.",
		},
		{
			name: "bullet list with strong and code marks",
			// Modeled on a task with bold labels and inline code references
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "bulletList",
						"content": []interface{}{
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Scope",
												"marks": []interface{}{
													map[string]interface{}{"type": "strong"},
												},
											},
											map[string]interface{}{
												"type": "text",
												"text": ": Process events from the ",
											},
											map[string]interface{}{
												"type": "text",
												"text": "event_queue",
												"marks": []interface{}{
													map[string]interface{}{"type": "code"},
												},
											},
											map[string]interface{}{
												"type": "text",
												"text": " table and dispatch notifications.",
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Result",
												"marks": []interface{}{
													map[string]interface{}{"type": "strong"},
												},
											},
											map[string]interface{}{
												"type": "text",
												"text": ": Reliable delivery with retry logic and audit trail.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: "Scope: Process events from the event_queue table and dispatch notifications.\n" +
				"Result: Reliable delivery with retry logic and audit trail.",
		},
		{
			name: "heading with bullet lists and code marks",
			// Modeled on a design doc with sections and feature lists
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Build a plugin system that loads extensions at startup.",
							},
						},
					},
					map[string]interface{}{
						"type": "heading",
						"attrs": map[string]interface{}{"level": 2},
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Architecture",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Core Library",
								"marks": []interface{}{
									map[string]interface{}{"type": "strong"},
								},
							},
						},
					},
					map[string]interface{}{
						"type": "bulletList",
						"content": []interface{}{
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Single ",
											},
											map[string]interface{}{
												"type": "text",
												"text": "register_plugin()",
												"marks": []interface{}{
													map[string]interface{}{"type": "code"},
												},
											},
											map[string]interface{}{
												"type": "text",
												"text": " API with JSON schema validation",
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Hot-reload support for development",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: "Build a plugin system that loads extensions at startup.\n" +
				"Architecture\n" +
				"Core Library\n" +
				"Single register_plugin() API with JSON schema validation\n" +
				"Hot-reload support for development",
		},
		{
			name: "ordered list with media skipped",
			// Modeled on a bug report with numbered steps and embedded screenshots
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Steps to reproduce the rendering issue:",
							},
						},
					},
					map[string]interface{}{
						"type": "orderedList",
						"attrs": map[string]interface{}{"order": 1},
						"content": []interface{}{
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Open the sample project and navigate to settings.",
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Click the export button.",
											},
										},
									},
									map[string]interface{}{
										"type":   "mediaSingle",
										"attrs":  map[string]interface{}{"layout": "align-start"},
										"content": []interface{}{
											map[string]interface{}{
												"type": "media",
												"attrs": map[string]interface{}{
													"type":       "file",
													"id":         "abc-123",
													"collection": "",
												},
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Observe the blank output where the chart should appear.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: "Steps to reproduce the rendering issue:\n" +
				"1. Open the sample project and navigate to settings.\n" +
				"2. Click the export button.\n" +
				"3. Observe the blank output where the chart should appear.",
		},
		{
			name: "hard breaks",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Line one of the note.",
							},
							map[string]interface{}{
								"type": "hardBreak",
							},
							map[string]interface{}{
								"type": "hardBreak",
							},
							map[string]interface{}{
								"type": "text",
								"text": "Line two after a gap.",
							},
						},
					},
				},
			},
			expected: "Line one of the note.\n\nLine two after a gap.",
		},
		{
			name: "code block",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Example configuration:",
							},
						},
					},
					map[string]interface{}{
						"type": "codeBlock",
						"attrs": map[string]interface{}{
							"language": "json",
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "{\"key\": \"value\"}",
							},
						},
					},
				},
			},
			expected: "Example configuration:\n{\"key\": \"value\"}",
		},
		{
			name: "horizontal rule",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Section one.",
							},
						},
					},
					map[string]interface{}{
						"type": "rule",
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Section two.",
							},
						},
					},
				},
			},
			expected: "Section one.\n---\nSection two.",
		},
		{
			name: "table",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "table",
						"content": []interface{}{
							map[string]interface{}{
								"type": "tableRow",
								"content": []interface{}{
									map[string]interface{}{
										"type": "tableHeader",
										"content": []interface{}{
											map[string]interface{}{
												"type": "paragraph",
												"content": []interface{}{
													map[string]interface{}{
														"type": "text",
														"text": "Name",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"type": "tableHeader",
										"content": []interface{}{
											map[string]interface{}{
												"type": "paragraph",
												"content": []interface{}{
													map[string]interface{}{
														"type": "text",
														"text": "Status",
													},
												},
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "tableRow",
								"content": []interface{}{
									map[string]interface{}{
										"type": "tableCell",
										"content": []interface{}{
											map[string]interface{}{
												"type": "paragraph",
												"content": []interface{}{
													map[string]interface{}{
														"type": "text",
														"text": "Widget A",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"type": "tableCell",
										"content": []interface{}{
											map[string]interface{}{
												"type": "paragraph",
												"content": []interface{}{
													map[string]interface{}{
														"type": "text",
														"text": "Active",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: "Name\tStatus\nWidget A\tActive",
		},
		{
			name: "mention",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Assigned to ",
							},
							map[string]interface{}{
								"type": "mention",
								"attrs": map[string]interface{}{
									"id":   "abc123",
									"text": "@Jane Smith",
								},
							},
							map[string]interface{}{
								"type": "text",
								"text": " for review.",
							},
						},
					},
				},
			},
			expected: "Assigned to @Jane Smith for review.",
		},
		{
			name: "mention without text attr",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "mention",
								"attrs": map[string]interface{}{
									"id": "abc123",
								},
							},
						},
					},
				},
			},
			expected: "@unknown",
		},
		{
			name: "inline card with URL",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "See ",
							},
							map[string]interface{}{
								"type": "inlineCard",
								"attrs": map[string]interface{}{
									"url": "https://example.com/doc",
								},
							},
						},
					},
				},
			},
			expected: "See https://example.com/doc",
		},
		{
			name: "emoji with text fallback",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Great work ",
							},
							map[string]interface{}{
								"type": "emoji",
								"attrs": map[string]interface{}{
									"shortName": ":thumbsup:",
									"text":      "\U0001F44D",
								},
							},
						},
					},
				},
			},
			expected: "Great work \U0001F44D",
		},
		{
			name: "date node",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Due: ",
							},
							map[string]interface{}{
								"type": "date",
								"attrs": map[string]interface{}{
									"timestamp": "2025-12-31",
								},
							},
						},
					},
				},
			},
			expected: "Due: 2025-12-31",
		},
		{
			name: "status node",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Current status: ",
							},
							map[string]interface{}{
								"type": "status",
								"attrs": map[string]interface{}{
									"text":  "IN PROGRESS",
									"color": "blue",
								},
							},
						},
					},
				},
			},
			expected: "Current status: IN PROGRESS",
		},
		{
			name: "empty doc",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{},
			},
			expected: "",
		},
		{
			name: "unknown node type with content recurses",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "futureNodeType",
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": "Content inside unknown node.",
									},
								},
							},
						},
					},
				},
			},
			expected: "Content inside unknown node.",
		},
		{
			name: "unknown node type without content is skipped",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Before.",
							},
						},
					},
					map[string]interface{}{
						"type": "someFutureWidget",
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "After.",
							},
						},
					},
				},
			},
			expected: "Before.\nAfter.",
		},
		{
			name: "blockquote",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "blockquote",
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": "A quoted remark.",
									},
								},
							},
						},
					},
				},
			},
			expected: "A quoted remark.",
		},
		{
			name: "panel with content",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "panel",
						"attrs": map[string]interface{}{
							"panelType": "info",
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": "This is an info panel.",
									},
								},
							},
						},
					},
				},
			},
			expected: "This is an info panel.",
		},
		{
			name: "nested bullet lists",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "bulletList",
						"content": []interface{}{
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Parent item",
											},
										},
									},
									map[string]interface{}{
										"type": "bulletList",
										"content": []interface{}{
											map[string]interface{}{
												"type": "listItem",
												"content": []interface{}{
													map[string]interface{}{
														"type": "paragraph",
														"content": []interface{}{
															map[string]interface{}{
																"type": "text",
																"text": "Child item one",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"type": "listItem",
												"content": []interface{}{
													map[string]interface{}{
														"type": "paragraph",
														"content": []interface{}{
															map[string]interface{}{
																"type": "text",
																"text": "Child item two",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							map[string]interface{}{
								"type": "listItem",
								"content": []interface{}{
									map[string]interface{}{
										"type": "paragraph",
										"content": []interface{}{
											map[string]interface{}{
												"type": "text",
												"text": "Second parent",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: "Parent item\nChild item one\nChild item two\nSecond parent",
		},
		{
			name: "expand node",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "expand",
						"attrs": map[string]interface{}{
							"title": "Click to expand",
						},
						"content": []interface{}{
							map[string]interface{}{
								"type": "paragraph",
								"content": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": "Hidden details here.",
									},
								},
							},
						},
					},
				},
			},
			expected: "Hidden details here.",
		},
		{
			name: "media nodes produce no output",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "See screenshot below.",
							},
						},
					},
					map[string]interface{}{
						"type": "mediaSingle",
						"attrs": map[string]interface{}{"layout": "center"},
						"content": []interface{}{
							map[string]interface{}{
								"type": "media",
								"attrs": map[string]interface{}{
									"type": "file",
									"id":   "img-001",
								},
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "End of report.",
							},
						},
					},
				},
			},
			expected: "See screenshot below.\nEnd of report.",
		},
		{
			name: "paragraph with inline media skipped",
			input: map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Refer to the attached file ",
							},
							map[string]interface{}{
								"type": "mediaInline",
								"attrs": map[string]interface{}{
									"type": "file",
									"id":   "file-001",
								},
							},
							map[string]interface{}{
								"type": "text",
								"text": " for details.",
							},
						},
					},
				},
			},
			expected: "Refer to the attached file  for details.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(tt.input)
			if got != tt.expected {
				t.Errorf("Render() =\n%q\nwant:\n%q", got, tt.expected)
			}
		})
	}
}

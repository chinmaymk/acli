package acli

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/chinmaymk/acli/internal/adf"
	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

var (
	searchWord = regexp.MustCompile(`[\pL\pN][\pL\pM\pN]*`)
	searchKey  = regexp.MustCompile(`(?i)^[a-z][a-z0-9_]*-[0-9]+$`)
	searchAge  = regexp.MustCompile(`^[1-9][0-9]*[dw]$`)
)

type issueSearchFilters struct {
	project, assignee, updatedSince string
	mine                            bool
}

func issueSearchQuery(cmd *cobra.Command, args []string, match string) (string, string, error) {
	var text string
	if len(args) > 0 {
		text = strings.TrimSpace(args[0])
		if text == "" {
			return "", "", fmt.Errorf("search text cannot be empty")
		}
	}
	jql, _ := cmd.Flags().GetString("jql")
	if cmd.Flags().Changed("jql") {
		if text != "" || strings.TrimSpace(jql) == "" {
			return "", "", fmt.Errorf("--jql requires a nonempty query and cannot be combined with search text")
		}
		return jql, "", nil
	}
	for _, flag := range []string{"project", "assignee", "updated-since"} {
		value, _ := cmd.Flags().GetString(flag)
		if cmd.Flags().Changed(flag) && strings.TrimSpace(value) == "" {
			return "", "", fmt.Errorf("--%s requires a value", flag)
		}
	}
	project, _ := cmd.Flags().GetString("project")
	assignee, _ := cmd.Flags().GetString("assignee")
	age, _ := cmd.Flags().GetString("updated-since")
	mine, _ := cmd.Flags().GetBool("mine")
	if match != "any" && match != "all" {
		return "", "", fmt.Errorf("--match must be all or any")
	}
	jql, err := textSearchJQL(text, match, issueSearchFilters{project: project, assignee: assignee, mine: mine, updatedSince: age})
	return jql, text, err
}

// Build JQL from literal words. Only --jql accepts query syntax.
func textSearchJQL(text, match string, f issueSearchFilters) (string, error) {
	var clauses []string
	if text != "" {
		if searchKey.MatchString(text) {
			clauses = append(clauses, "key = "+strconv.Quote(strings.ToUpper(text)))
		} else {
			words := searchWord.FindAllString(text, -1)
			if len(words) == 0 {
				return "", fmt.Errorf("search text must contain a letter or number")
			}
			terms := make([]string, 0, len(words))
			for _, word := range words {
				// Quote the word for both Jira's text parser and JQL.
				terms = append(terms, "text ~ "+strconv.Quote(strconv.Quote(word)))
			}
			if match == "any" && len(terms) > 1 {
				clauses = append(clauses, "("+strings.Join(terms, " OR ")+")")
			} else {
				clauses = append(clauses, terms...)
			}
		}
	}
	if f.project != "" {
		clauses = append(clauses, "project = "+strconv.Quote(f.project))
	}
	if f.mine && f.assignee != "" {
		return "", fmt.Errorf("choose either --mine or --assignee")
	}
	if f.mine {
		clauses = append(clauses, "assignee = currentUser()")
	}
	if f.assignee != "" {
		clauses = append(clauses, "assignee = "+strconv.Quote(f.assignee))
	}
	if f.updatedSince != "" {
		if !searchAge.MatchString(f.updatedSince) {
			return "", fmt.Errorf("--updated-since requires a positive number of days or weeks (e.g. 30d or 2w)")
		}
		clauses = append(clauses, "updated >= "+strconv.Quote("-"+f.updatedSince))
	}
	if len(clauses) == 0 {
		return "", fmt.Errorf("provide search text, a filter, or --jql")
	}
	return strings.Join(clauses, " AND "), nil
}

type issuePreview struct {
	matchedWords int    // Used only for sorting match-any previews.
	Key          string `json:"key"`
	Summary      string `json:"summary"`
	Status       string `json:"status"`
	URL          string `json:"url"`
	Excerpt      string `json:"excerpt,omitempty"`
}

func issuePreviews(issues []jira.IssueDetailed, text, baseURL, match string) []issuePreview {
	previews := make([]issuePreview, 0, len(issues))
	words := searchWord.FindAllString(strings.ToLower(text), -1)
	for _, issue := range issues {
		description := adf.Render(issue.Fields.Description)
		body := strings.Fields(description)
		var excerpt string
		for i, word := range body {
			if containsSearchWord(word, words) {
				start, end := max(0, i-5), min(len(body), i+14)
				runes := []rune(strings.Join(body[start:end], " "))
				excerpt = string(runes[:min(len(runes), 160)])
				break
			}
		}
		var matchedWords int
		if match == "any" {
			matchedWords = countSearchWords(issue.Fields.Summary+" "+description, words)
		}
		previews = append(previews, issuePreview{
			matchedWords: matchedWords,
			Key:          issue.Key, Summary: issue.Fields.Summary, Status: renderIssueField("status", issue),
			URL: baseURL + "/browse/" + url.PathEscape(issue.Key), Excerpt: excerpt,
		})
	}
	// Promote title matches within the retrieved page, keeping Jira's order for ties.
	sort.SliceStable(previews, func(i, j int) bool {
		if previews[i].matchedWords != previews[j].matchedWords {
			return previews[i].matchedWords > previews[j].matchedWords
		}
		return containsSearchWord(previews[i].Summary, words) && !containsSearchWord(previews[j].Summary, words)
	})
	return previews
}

// Count distinct query words visible in the summary and rendered description.
func countSearchWords(value string, words []string) int {
	remaining := make(map[string]struct{}, len(words))
	for _, word := range words {
		remaining[word] = struct{}{}
	}
	total := len(remaining)
	for _, word := range searchWord.FindAllString(strings.ToLower(value), -1) {
		delete(remaining, word)
	}
	return total - len(remaining)
}

func containsSearchWord(value string, words []string) bool {
	value = strings.ToLower(value)
	for _, word := range words {
		if strings.Contains(value, word) {
			return true
		}
	}
	return false
}

package jira

import (
	"encoding/json"
	"testing"
)

func TestIssueDetailedSerialization(t *testing.T) {
	issue := IssueDetailed{
		ID:   "10001",
		Key:  "TEST-1",
		Self: "https://test.atlassian.net/rest/api/3/issue/10001",
		Fields: IssueFields{
			Summary: "Test issue",
			IssueType: &IssueType{
				ID:   "10",
				Name: "Bug",
			},
			Status: &StatusDetails{
				Name: "Open",
				StatusCategory: StatusCategory{
					Name:      "To Do",
					ColorName: "blue-gray",
				},
			},
			Priority: &Priority{
				Name: "High",
			},
			Assignee: &UserDetails{
				DisplayName: "John Doe",
				AccountID:   "abc123",
			},
		},
	}

	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded IssueDetailed
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Key != "TEST-1" {
		t.Errorf("expected key TEST-1, got %s", decoded.Key)
	}
	if decoded.Fields.Summary != "Test issue" {
		t.Errorf("expected summary 'Test issue', got %s", decoded.Fields.Summary)
	}
	if decoded.Fields.IssueType.Name != "Bug" {
		t.Errorf("expected type Bug, got %s", decoded.Fields.IssueType.Name)
	}
	if decoded.Fields.Status.Name != "Open" {
		t.Errorf("expected status Open, got %s", decoded.Fields.Status.Name)
	}
	if decoded.Fields.Priority.Name != "High" {
		t.Errorf("expected priority High, got %s", decoded.Fields.Priority.Name)
	}
	if decoded.Fields.Assignee.DisplayName != "John Doe" {
		t.Errorf("expected assignee John Doe, got %s", decoded.Fields.Assignee.DisplayName)
	}
}

func TestIssueDetailedDeserialization(t *testing.T) {
	raw := `{
		"id": "10001",
		"key": "PROJ-42",
		"self": "https://test.atlassian.net/rest/api/3/issue/10001",
		"fields": {
			"summary": "Fix login bug",
			"issuetype": {"id": "1", "name": "Bug", "subtask": false},
			"status": {"name": "In Progress", "statusCategory": {"name": "In Progress", "colorName": "blue"}},
			"priority": {"name": "Medium"},
			"assignee": {"displayName": "Jane Smith", "accountId": "xyz789"},
			"reporter": {"displayName": "Bob", "accountId": "bob1"},
			"created": "2024-01-15T10:00:00.000+0000",
			"updated": "2024-01-16T14:30:00.000+0000"
		}
	}`

	var issue IssueDetailed
	if err := json.Unmarshal([]byte(raw), &issue); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if issue.Key != "PROJ-42" {
		t.Errorf("unexpected key: %s", issue.Key)
	}
	if issue.Fields.Summary != "Fix login bug" {
		t.Errorf("unexpected summary: %s", issue.Fields.Summary)
	}
	if issue.Fields.IssueType == nil || issue.Fields.IssueType.Name != "Bug" {
		t.Error("unexpected issue type")
	}
	if issue.Fields.Status == nil || issue.Fields.Status.Name != "In Progress" {
		t.Error("unexpected status")
	}
	if issue.Fields.Reporter == nil || issue.Fields.Reporter.DisplayName != "Bob" {
		t.Error("unexpected reporter")
	}
}

func TestIssueFieldsNilPointers(t *testing.T) {
	raw := `{
		"id": "10002",
		"key": "PROJ-43",
		"fields": {
			"summary": "Minimal issue"
		}
	}`

	var issue IssueDetailed
	if err := json.Unmarshal([]byte(raw), &issue); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// All optional fields should be nil
	if issue.Fields.IssueType != nil {
		t.Error("expected nil issue type")
	}
	if issue.Fields.Status != nil {
		t.Error("expected nil status")
	}
	if issue.Fields.Priority != nil {
		t.Error("expected nil priority")
	}
	if issue.Fields.Assignee != nil {
		t.Error("expected nil assignee")
	}
	if issue.Fields.Reporter != nil {
		t.Error("expected nil reporter")
	}
}

func TestPageBeanSerialization(t *testing.T) {
	page := PageBean[Project]{
		Pagination: Pagination{
			StartAt:    0,
			MaxResults: 50,
			Total:      2,
		},
		IsLast: true,
		Values: []Project{
			{Key: "P1", Name: "Project 1"},
			{Key: "P2", Name: "Project 2"},
		},
	}

	data, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded PageBean[Project]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Total != 2 {
		t.Errorf("expected total 2, got %d", decoded.Total)
	}
	if !decoded.IsLast {
		t.Error("expected isLast true")
	}
	if len(decoded.Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(decoded.Values))
	}
}

func TestUserDetailsSerialization(t *testing.T) {
	user := UserDetails{
		AccountID:    "abc123",
		DisplayName:  "Test User",
		EmailAddress: "test@example.com",
		Active:       true,
		TimeZone:     "US/Pacific",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded UserDetails
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.AccountID != "abc123" {
		t.Errorf("unexpected accountId: %s", decoded.AccountID)
	}
	if decoded.DisplayName != "Test User" {
		t.Errorf("unexpected displayName: %s", decoded.DisplayName)
	}
	if !decoded.Active {
		t.Error("expected active true")
	}
}

func TestStatusDetailsSerialization(t *testing.T) {
	status := StatusDetails{
		ID:   "1",
		Name: "Done",
		StatusCategory: StatusCategory{
			ID:        3,
			Key:       "done",
			Name:      "Done",
			ColorName: "green",
		},
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded StatusDetails
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Name != "Done" {
		t.Errorf("unexpected name: %s", decoded.Name)
	}
	if decoded.StatusCategory.ColorName != "green" {
		t.Errorf("unexpected color: %s", decoded.StatusCategory.ColorName)
	}
}

func TestSearchResultsDeserialization(t *testing.T) {
	raw := `{
		"startAt": 0,
		"maxResults": 50,
		"total": 1,
		"issues": [
			{"key": "TEST-1", "fields": {"summary": "Found issue"}}
		]
	}`

	var result SearchResults
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Key != "TEST-1" {
		t.Errorf("unexpected key: %s", result.Issues[0].Key)
	}
}

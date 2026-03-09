package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// CreateIssue creates a new Jira issue.
func (c *Client) CreateIssue(details *IssueUpdateDetails) (*CreatedIssue, error) {
	var result CreatedIssue
	err := c.Post("/rest/api/3/issue", details, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BulkCreateIssues creates multiple issues in a single request.
func (c *Client) BulkCreateIssues(req *BulkIssueCreateRequest) (*BulkIssueCreateResponse, error) {
	var result BulkIssueCreateResponse
	err := c.Post("/rest/api/3/issue/bulk", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssue retrieves a single issue by ID or key.
func (c *Client) GetIssue(issueIdOrKey string, fields []string, expand []string) (*IssueDetailed, error) {
	query := url.Values{}
	if len(fields) > 0 {
		query.Set("fields", strings.Join(fields, ","))
	}
	if len(expand) > 0 {
		query.Set("expand", strings.Join(expand, ","))
	}
	var result IssueDetailed
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey), query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EditIssue edits an existing issue.
func (c *Client) EditIssue(issueIdOrKey string, details *IssueUpdateDetails, notifyUsers bool) error {
	query := url.Values{}
	if !notifyUsers {
		query.Set("notifyUsers", "false")
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey)
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.Put(path, details, nil)
}

// DeleteIssue deletes an issue.
func (c *Client) DeleteIssue(issueIdOrKey string, deleteSubtasks bool) error {
	query := url.Values{}
	if deleteSubtasks {
		query.Set("deleteSubtasks", "true")
	}
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey), query)
}

// AssignIssue assigns an issue to a user.
func (c *Client) AssignIssue(issueIdOrKey string, accountID string) error {
	body := map[string]interface{}{
		"accountId": accountID,
	}
	return c.Put(fmt.Sprintf("/rest/api/3/issue/%s/assignee", issueIdOrKey), body, nil)
}

// GetIssueChangelog returns the changelog for an issue.
func (c *Client) GetIssueChangelog(issueIdOrKey string, startAt, maxResults int) (*ChangelogPage, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result ChangelogPage
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/changelog", issueIdOrKey), query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueComments returns the comments for an issue.
func (c *Client) GetIssueComments(issueIdOrKey string, startAt, maxResults int) (*CommentPage, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result CommentPage
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIdOrKey), query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddIssueComment adds a comment to an issue.
func (c *Client) AddIssueComment(issueIdOrKey string, body interface{}, visibility *Visibility) (*Comment, error) {
	reqBody := map[string]interface{}{
		"body": body,
	}
	if visibility != nil {
		reqBody["visibility"] = visibility
	}
	var result Comment
	err := c.Post(fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIdOrKey), reqBody, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueComment retrieves a single comment on an issue.
func (c *Client) GetIssueComment(issueIdOrKey, commentId string) (*Comment, error) {
	var result Comment
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIdOrKey, commentId), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueComment updates an existing comment on an issue.
func (c *Client) UpdateIssueComment(issueIdOrKey, commentId string, body interface{}, visibility *Visibility) (*Comment, error) {
	reqBody := map[string]interface{}{
		"body": body,
	}
	if visibility != nil {
		reqBody["visibility"] = visibility
	}
	var result Comment
	err := c.Put(fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIdOrKey, commentId), reqBody, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssueComment deletes a comment from an issue.
func (c *Client) DeleteIssueComment(issueIdOrKey, commentId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIdOrKey, commentId), nil)
}

// GetIssueTransitions returns the available transitions for an issue.
func (c *Client) GetIssueTransitions(issueIdOrKey string) (*TransitionsResponse, error) {
	var result TransitionsResponse
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DoIssueTransition performs a transition on an issue.
func (c *Client) DoIssueTransition(issueIdOrKey string, details *IssueUpdateDetails) error {
	return c.Post(fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIdOrKey), details, nil)
}

// GetIssueVotes returns the votes for an issue.
func (c *Client) GetIssueVotes(issueIdOrKey string) (*Votes, error) {
	var result Votes
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddIssueVote adds a vote to an issue for the authenticated user.
func (c *Client) AddIssueVote(issueIdOrKey string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIdOrKey), nil, nil)
}

// RemoveIssueVote removes a vote from an issue for the authenticated user.
func (c *Client) RemoveIssueVote(issueIdOrKey string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIdOrKey), nil)
}

// GetIssueWatchers returns the watchers for an issue.
func (c *Client) GetIssueWatchers(issueIdOrKey string) (*Watches, error) {
	var result Watches
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddIssueWatcher adds a watcher to an issue.
func (c *Client) AddIssueWatcher(issueIdOrKey string, accountID string) error {
	// The Jira API expects the account ID as a quoted JSON string in the body.
	quotedAccountID := fmt.Sprintf("%q", accountID)
	return c.Post(fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey), json.RawMessage(quotedAccountID), nil)
}

// RemoveIssueWatcher removes a watcher from an issue.
func (c *Client) RemoveIssueWatcher(issueIdOrKey string, accountID string) error {
	query := url.Values{}
	query.Set("accountId", accountID)
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey), query)
}

// GetIssueWorklogs returns the worklogs for an issue.
func (c *Client) GetIssueWorklogs(issueIdOrKey string, startAt, maxResults int) (*WorklogPage, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result WorklogPage
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIdOrKey), query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddIssueWorklog adds a worklog to an issue.
func (c *Client) AddIssueWorklog(issueIdOrKey string, worklog *Worklog) (*Worklog, error) {
	var result Worklog
	err := c.Post(fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIdOrKey), worklog, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueWorklog retrieves a single worklog entry.
func (c *Client) GetIssueWorklog(issueIdOrKey, worklogId string) (*Worklog, error) {
	var result Worklog
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIdOrKey, worklogId), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueWorklog updates a worklog entry.
func (c *Client) UpdateIssueWorklog(issueIdOrKey, worklogId string, worklog *Worklog) (*Worklog, error) {
	var result Worklog
	err := c.Put(fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIdOrKey, worklogId), worklog, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssueWorklog deletes a worklog entry.
func (c *Client) DeleteIssueWorklog(issueIdOrKey, worklogId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIdOrKey, worklogId), nil)
}

// AddIssueAttachment uploads an attachment to an issue.
func (c *Client) AddIssueAttachment(issueIdOrKey string, filePath string) ([]Attachment, error) {
	var result []Attachment
	err := c.UploadFile(fmt.Sprintf("/rest/api/3/issue/%s/attachments", issueIdOrKey), "file", filePath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetIssueRemoteLinks returns the remote links for an issue.
func (c *Client) GetIssueRemoteLinks(issueIdOrKey string) ([]RemoteIssueLink, error) {
	var result []RemoteIssueLink
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/remotelink", issueIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateIssueRemoteLink creates a remote link on an issue.
func (c *Client) CreateIssueRemoteLink(issueIdOrKey string, link *RemoteIssueLink) (*RemoteIssueLink, error) {
	var result RemoteIssueLink
	err := c.Post(fmt.Sprintf("/rest/api/3/issue/%s/remotelink", issueIdOrKey), link, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueRemoteLink retrieves a single remote link on an issue.
func (c *Client) GetIssueRemoteLink(issueIdOrKey string, linkId string) (*RemoteIssueLink, error) {
	var result RemoteIssueLink
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/remotelink/%s", issueIdOrKey, linkId), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueRemoteLink updates a remote link on an issue.
func (c *Client) UpdateIssueRemoteLink(issueIdOrKey string, linkId string, link *RemoteIssueLink) error {
	return c.Put(fmt.Sprintf("/rest/api/3/issue/%s/remotelink/%s", issueIdOrKey, linkId), link, nil)
}

// DeleteIssueRemoteLink deletes a remote link on an issue.
func (c *Client) DeleteIssueRemoteLink(issueIdOrKey string, linkId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/remotelink/%s", issueIdOrKey, linkId), nil)
}

// NotifyIssue sends a notification for an issue.
func (c *Client) NotifyIssue(issueIdOrKey string, notify *IssueNotifyRequest) error {
	return c.Post(fmt.Sprintf("/rest/api/3/issue/%s/notify", issueIdOrKey), notify, nil)
}

// GetIssueEditMeta returns the edit metadata for an issue.
func (c *Client) GetIssueEditMeta(issueIdOrKey string) (json.RawMessage, error) {
	var result json.RawMessage
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/editmeta", issueIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateMeta returns the create metadata for issues.
func (c *Client) GetCreateMeta(projectKeys []string, expand []string) (*CreateMeta, error) {
	query := url.Values{}
	if len(projectKeys) > 0 {
		query.Set("projectKeys", strings.Join(projectKeys, ","))
	}
	if len(expand) > 0 {
		query.Set("expand", strings.Join(expand, ","))
	}
	var result CreateMeta
	err := c.Get("/rest/api/3/issue/createmeta", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueProperties returns all property keys for an issue.
func (c *Client) GetIssueProperties(issueIdOrKey string) ([]EntityProperty, error) {
	var wrapper struct {
		Keys []EntityProperty `json:"keys"`
	}
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/properties", issueIdOrKey), nil, &wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Keys, nil
}

// GetIssueProperty retrieves a single property of an issue.
func (c *Client) GetIssueProperty(issueIdOrKey, propertyKey string) (*EntityProperty, error) {
	var result EntityProperty
	err := c.Get(fmt.Sprintf("/rest/api/3/issue/%s/properties/%s", issueIdOrKey, propertyKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetIssueProperty sets a property on an issue.
func (c *Client) SetIssueProperty(issueIdOrKey, propertyKey string, value interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/issue/%s/properties/%s", issueIdOrKey, propertyKey), value, nil)
}

// DeleteIssueProperty deletes a property from an issue.
func (c *Client) DeleteIssueProperty(issueIdOrKey, propertyKey string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issue/%s/properties/%s", issueIdOrKey, propertyKey), nil)
}

// BulkFetchIssues fetches multiple issues by their IDs.
func (c *Client) BulkFetchIssues(issueIDs []int, fields []string) (*SearchResults, error) {
	body := map[string]interface{}{
		"issueIds": issueIDs,
	}
	if len(fields) > 0 {
		body["fields"] = fields
	}
	var result SearchResults
	err := c.Post("/rest/api/3/issue/bulkfetch", body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

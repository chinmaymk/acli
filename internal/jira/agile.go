package jira

import (
	"fmt"
	"net/url"
	"strconv"
)

// --- Boards ---

// GetBoards returns a paginated list of boards.
func (c *Client) GetBoards(startAt, maxResults int, projectKeyOrID, boardType, name string) (*BoardList, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if projectKeyOrID != "" {
		query.Set("projectKeyOrId", projectKeyOrID)
	}
	if boardType != "" {
		query.Set("type", boardType)
	}
	if name != "" {
		query.Set("name", name)
	}
	var result BoardList
	if err := c.Get("/rest/agile/1.0/board", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoard returns a board by ID.
func (c *Client) GetBoard(boardID int) (*Board, error) {
	var result Board
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d", boardID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoardConfiguration returns a board's configuration.
func (c *Client) GetBoardConfiguration(boardID int) (*BoardConfiguration, error) {
	var result BoardConfiguration
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d/configuration", boardID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoardIssues returns issues for a board.
func (c *Client) GetBoardIssues(boardID, startAt, maxResults int, jql string) (*SprintIssuesResponse, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if jql != "" {
		query.Set("jql", jql)
	}
	var result SprintIssuesResponse
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d/issue", boardID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoardBacklog returns backlog issues for a board.
func (c *Client) GetBoardBacklog(boardID, startAt, maxResults int, jql string) (*SprintIssuesResponse, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if jql != "" {
		query.Set("jql", jql)
	}
	var result SprintIssuesResponse
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d/backlog", boardID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoardSprints returns sprints for a board.
func (c *Client) GetBoardSprints(boardID, startAt, maxResults int, state string) (*SprintList, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if state != "" {
		query.Set("state", state)
	}
	var result SprintList
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d/sprint", boardID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBoardEpics returns epics for a board.
func (c *Client) GetBoardEpics(boardID, startAt, maxResults int) (*EpicList, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result EpicList
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/board/%d/epic", boardID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Sprints ---

// GetSprint returns a sprint by ID.
func (c *Client) GetSprint(sprintID int) (*Sprint, error) {
	var result Sprint
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSprint creates a new sprint.
func (c *Client) CreateSprint(sprint map[string]interface{}) (*Sprint, error) {
	var result Sprint
	if err := c.Post("/rest/agile/1.0/sprint", sprint, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSprint updates a sprint.
func (c *Client) UpdateSprint(sprintID int, sprint map[string]interface{}) (*Sprint, error) {
	var result Sprint
	if err := c.Put(fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), sprint, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PartialUpdateSprint partially updates a sprint.
func (c *Client) PartialUpdateSprint(sprintID int, sprint map[string]interface{}) (*Sprint, error) {
	var result Sprint
	if err := c.Post(fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), sprint, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteSprint deletes a sprint.
func (c *Client) DeleteSprint(sprintID int) error {
	return c.Delete(fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), nil)
}

// GetSprintIssues returns issues in a sprint.
func (c *Client) GetSprintIssues(sprintID, startAt, maxResults int, jql string) (*SprintIssuesResponse, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if jql != "" {
		query.Set("jql", jql)
	}
	var result SprintIssuesResponse
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MoveIssuesToSprint moves issues to a sprint.
func (c *Client) MoveIssuesToSprint(sprintID int, issueKeys []string) error {
	body := map[string]interface{}{
		"issues": issueKeys,
	}
	return c.Post(fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID), body, nil)
}

// MoveIssuesToBacklog moves issues to the backlog.
func (c *Client) MoveIssuesToBacklog(issueKeys []string) error {
	body := map[string]interface{}{
		"issues": issueKeys,
	}
	return c.Post("/rest/agile/1.0/backlog/issue", body, nil)
}

// --- Epics ---

// GetEpic returns an epic by ID or key.
func (c *Client) GetEpic(epicIdOrKey string) (*Epic, error) {
	var result Epic
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/epic/%s", epicIdOrKey), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MoveIssuesToEpic moves issues to an epic.
func (c *Client) MoveIssuesToEpic(epicIdOrKey string, issueKeys []string) error {
	body := map[string]interface{}{
		"issues": issueKeys,
	}
	return c.Post(fmt.Sprintf("/rest/agile/1.0/epic/%s/issue", epicIdOrKey), body, nil)
}

// GetEpicIssues returns issues belonging to an epic.
func (c *Client) GetEpicIssues(epicIdOrKey string, startAt, maxResults int, jql string) (*SprintIssuesResponse, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if jql != "" {
		query.Set("jql", jql)
	}
	var result SprintIssuesResponse
	if err := c.Get(fmt.Sprintf("/rest/agile/1.0/epic/%s/issue", epicIdOrKey), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

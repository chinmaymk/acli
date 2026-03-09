package jira

import (
	"fmt"
	"net/url"
	"strconv"
)

// ============================================================================
// Components
// ============================================================================

// CreateComponent creates a new project component.
func (c *Client) CreateComponent(component map[string]interface{}) (*ProjectComponent, error) {
	var result ProjectComponent
	if err := c.Post("/rest/api/3/component", component, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetComponent returns a project component by ID.
func (c *Client) GetComponent(id string) (*ProjectComponent, error) {
	var result ProjectComponent
	if err := c.Get(fmt.Sprintf("/rest/api/3/component/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateComponent updates a project component.
func (c *Client) UpdateComponent(id string, component map[string]interface{}) (*ProjectComponent, error) {
	var result ProjectComponent
	if err := c.Put(fmt.Sprintf("/rest/api/3/component/%s", id), component, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteComponent deletes a project component.
func (c *Client) DeleteComponent(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/component/%s", id), nil)
}

// GetComponentIssueCount returns the issue count for a component.
func (c *Client) GetComponentIssueCount(id string) (*ComponentIssueCount, error) {
	var result ComponentIssueCount
	if err := c.Get(fmt.Sprintf("/rest/api/3/component/%s/relatedIssueCounts", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Versions
// ============================================================================

// CreateVersion creates a new project version.
func (c *Client) CreateVersion(version *Version) (*Version, error) {
	var result Version
	if err := c.Post("/rest/api/3/version", version, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVersion returns a project version by ID.
func (c *Client) GetVersion(id string) (*Version, error) {
	var result Version
	if err := c.Get(fmt.Sprintf("/rest/api/3/version/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateVersion updates a project version.
func (c *Client) UpdateVersion(id string, version *Version) (*Version, error) {
	var result Version
	if err := c.Put(fmt.Sprintf("/rest/api/3/version/%s", id), version, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteVersion deletes a project version.
func (c *Client) DeleteVersion(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/version/%s", id), nil)
}

// MergeVersions merges a version into another version.
func (c *Client) MergeVersions(id, moveIssuesTo string) error {
	return c.Put(fmt.Sprintf("/rest/api/3/version/%s/mergeto/%s", id, moveIssuesTo), nil, nil)
}

// MoveVersion moves a version to a new position.
func (c *Client) MoveVersion(id string, position map[string]interface{}) (*Version, error) {
	var result Version
	if err := c.Post(fmt.Sprintf("/rest/api/3/version/%s/move", id), position, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVersionRelatedIssueCounts returns issue counts related to a version.
func (c *Client) GetVersionRelatedIssueCounts(id string) (*VersionIssueCounts, error) {
	var result VersionIssueCounts
	if err := c.Get(fmt.Sprintf("/rest/api/3/version/%s/relatedIssueCounts", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVersionUnresolvedIssueCount returns the unresolved issue count for a version.
func (c *Client) GetVersionUnresolvedIssueCount(id string) (*VersionUnresolvedIssueCount, error) {
	var result VersionUnresolvedIssueCount
	if err := c.Get(fmt.Sprintf("/rest/api/3/version/%s/unresolvedIssueCount", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Users
// ============================================================================

// GetUser returns a user by account ID.
func (c *Client) GetUser(accountID string) (*UserDetails, error) {
	var result UserDetails
	query := url.Values{}
	query.Set("accountId", accountID)
	if err := c.Get("/rest/api/3/user", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateUser creates a new user.
func (c *Client) CreateUser(user map[string]interface{}) (*UserDetails, error) {
	var result UserDetails
	if err := c.Post("/rest/api/3/user", user, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteUser deletes a user by account ID.
func (c *Client) DeleteUser(accountID string) error {
	query := url.Values{}
	query.Set("accountId", accountID)
	return c.Delete("/rest/api/3/user", query)
}

// GetUsersBulk returns multiple users by account IDs.
func (c *Client) GetUsersBulk(accountIDs []string, startAt, maxResults int) (*PageBean[UserDetails], error) {
	var result PageBean[UserDetails]
	query := url.Values{}
	for _, id := range accountIDs {
		query.Add("accountId", id)
	}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/user/bulk", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FindUsers searches for users by query string.
func (c *Client) FindUsers(query string, startAt, maxResults int) ([]UserDetails, error) {
	var result []UserDetails
	params := url.Values{}
	params.Set("query", query)
	params.Set("startAt", strconv.Itoa(startAt))
	params.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/user/search", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FindUsersAssignable searches for users assignable to issues.
func (c *Client) FindUsersAssignable(query, project, issueKey string, startAt, maxResults int) ([]UserDetails, error) {
	var result []UserDetails
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if project != "" {
		params.Set("project", project)
	}
	if issueKey != "" {
		params.Set("issueKey", issueKey)
	}
	params.Set("startAt", strconv.Itoa(startAt))
	params.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/user/assignable/search", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetCurrentUser returns the currently authenticated user.
func (c *Client) GetCurrentUser() (*UserDetails, error) {
	var result UserDetails
	if err := c.Get("/rest/api/3/myself", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllUsers returns all users with pagination.
func (c *Client) GetAllUsers(startAt, maxResults int) ([]UserDetails, error) {
	var result []UserDetails
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/users/search", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ============================================================================
// Groups
// ============================================================================

// GetGroup returns a group by name.
func (c *Client) GetGroup(groupName string) (*Group, error) {
	var result Group
	query := url.Values{}
	query.Set("groupname", groupName)
	if err := c.Get("/rest/api/3/group", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateGroup creates a new group.
func (c *Client) CreateGroup(name string) (*Group, error) {
	var result Group
	body := map[string]string{"name": name}
	if err := c.Post("/rest/api/3/group", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteGroup deletes a group by name.
func (c *Client) DeleteGroup(groupName string) error {
	query := url.Values{}
	query.Set("groupname", groupName)
	return c.Delete("/rest/api/3/group", query)
}

// GetGroupMembers returns the members of a group.
func (c *Client) GetGroupMembers(groupName string, startAt, maxResults int) (*GroupMembers, error) {
	var result GroupMembers
	query := url.Values{}
	query.Set("groupname", groupName)
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/group/member", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddUserToGroup adds a user to a group.
func (c *Client) AddUserToGroup(groupName, accountID string) (*Group, error) {
	var result Group
	query := url.Values{}
	query.Set("groupname", groupName)
	body := map[string]string{"accountId": accountID}
	path := "/rest/api/3/group/user?" + query.Encode()
	if err := c.Post(path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveUserFromGroup removes a user from a group.
func (c *Client) RemoveUserFromGroup(groupName, accountID string) error {
	query := url.Values{}
	query.Set("groupname", groupName)
	query.Set("accountId", accountID)
	return c.Delete("/rest/api/3/group/user", query)
}

// GetBulkGroups returns groups with pagination.
func (c *Client) GetBulkGroups(startAt, maxResults int) (*PageBean[Group], error) {
	var result PageBean[Group]
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/group/bulk", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FindGroups searches for groups using the picker endpoint.
func (c *Client) FindGroups(query string, maxResults int) (*FoundGroups, error) {
	var result FoundGroups
	params := url.Values{}
	params.Set("query", query)
	params.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/groups/picker", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Issue Links
// ============================================================================

// CreateIssueLink creates a link between two issues.
func (c *Client) CreateIssueLink(link *IssueLink) error {
	return c.Post("/rest/api/3/issueLink", link, nil)
}

// GetIssueLink returns an issue link by ID.
func (c *Client) GetIssueLink(linkId string) (*IssueLink, error) {
	var result IssueLink
	if err := c.Get(fmt.Sprintf("/rest/api/3/issueLink/%s", linkId), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssueLink deletes an issue link.
func (c *Client) DeleteIssueLink(linkId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issueLink/%s", linkId), nil)
}

// GetIssueLinkTypes returns all issue link types.
func (c *Client) GetIssueLinkTypes() ([]IssueLinkType, error) {
	var result struct {
		IssueLinkTypes []IssueLinkType `json:"issueLinkTypes"`
	}
	if err := c.Get("/rest/api/3/issueLinkType", nil, &result); err != nil {
		return nil, err
	}
	return result.IssueLinkTypes, nil
}

// CreateIssueLinkType creates a new issue link type.
func (c *Client) CreateIssueLinkType(linkType *IssueLinkType) (*IssueLinkType, error) {
	var result IssueLinkType
	if err := c.Post("/rest/api/3/issueLinkType", linkType, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueLinkType returns an issue link type by ID.
func (c *Client) GetIssueLinkType(id string) (*IssueLinkType, error) {
	var result IssueLinkType
	if err := c.Get(fmt.Sprintf("/rest/api/3/issueLinkType/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueLinkType updates an issue link type.
func (c *Client) UpdateIssueLinkType(id string, linkType *IssueLinkType) (*IssueLinkType, error) {
	var result IssueLinkType
	if err := c.Put(fmt.Sprintf("/rest/api/3/issueLinkType/%s", id), linkType, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssueLinkType deletes an issue link type.
func (c *Client) DeleteIssueLinkType(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issueLinkType/%s", id), nil)
}

// ============================================================================
// Attachments
// ============================================================================

// GetAttachment returns an attachment by ID.
func (c *Client) GetAttachment(id string) (*Attachment, error) {
	var result Attachment
	if err := c.Get(fmt.Sprintf("/rest/api/3/attachment/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAttachment deletes an attachment.
func (c *Client) DeleteAttachment(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/attachment/%s", id), nil)
}

// GetAttachmentMeta returns attachment settings.
func (c *Client) GetAttachmentMeta() (*AttachmentMeta, error) {
	var result AttachmentMeta
	if err := c.Get("/rest/api/3/attachment/meta", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Issue Types
// ============================================================================

// GetAllIssueTypes returns all issue types.
func (c *Client) GetAllIssueTypes() ([]IssueType, error) {
	var result []IssueType
	if err := c.Get("/rest/api/3/issuetype", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateIssueType creates a new issue type.
func (c *Client) CreateIssueType(issueType map[string]interface{}) (*IssueType, error) {
	var result IssueType
	if err := c.Post("/rest/api/3/issuetype", issueType, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueType returns an issue type by ID.
func (c *Client) GetIssueType(id string) (*IssueType, error) {
	var result IssueType
	if err := c.Get(fmt.Sprintf("/rest/api/3/issuetype/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueType updates an issue type.
func (c *Client) UpdateIssueType(id string, issueType map[string]interface{}) (*IssueType, error) {
	var result IssueType
	if err := c.Put(fmt.Sprintf("/rest/api/3/issuetype/%s", id), issueType, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssueType deletes an issue type.
func (c *Client) DeleteIssueType(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issuetype/%s", id), nil)
}

// GetIssueTypeAlternatives returns alternative issue types for the given issue type.
func (c *Client) GetIssueTypeAlternatives(id string) ([]IssueType, error) {
	var result []IssueType
	if err := c.Get(fmt.Sprintf("/rest/api/3/issuetype/%s/alternatives", id), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetProjectIssueTypes returns issue types for a project.
func (c *Client) GetProjectIssueTypes(projectIdOrKey string) ([]IssueType, error) {
	var result []IssueType
	query := url.Values{}
	query.Set("projectId", projectIdOrKey)
	if err := c.Get("/rest/api/3/issuetype/project", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ============================================================================
// Priorities
// ============================================================================

// GetAllPriorities returns all priorities.
func (c *Client) GetAllPriorities() ([]Priority, error) {
	var result []Priority
	if err := c.Get("/rest/api/3/priority", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreatePriority creates a new priority.
func (c *Client) CreatePriority(priority map[string]interface{}) (*Priority, error) {
	var result Priority
	if err := c.Post("/rest/api/3/priority", priority, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPriority returns a priority by ID.
func (c *Client) GetPriority(id string) (*Priority, error) {
	var result Priority
	if err := c.Get(fmt.Sprintf("/rest/api/3/priority/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePriority updates a priority.
func (c *Client) UpdatePriority(id string, priority map[string]interface{}) (*Priority, error) {
	var result Priority
	if err := c.Put(fmt.Sprintf("/rest/api/3/priority/%s", id), priority, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePriority deletes a priority.
func (c *Client) DeletePriority(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/priority/%s", id), nil)
}

// SearchPriorities searches for priorities with pagination.
func (c *Client) SearchPriorities(startAt, maxResults int) (*PageBean[Priority], error) {
	var result PageBean[Priority]
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/priority/search", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Resolutions
// ============================================================================

// GetAllResolutions returns all resolutions.
func (c *Client) GetAllResolutions() ([]Resolution, error) {
	var result []Resolution
	if err := c.Get("/rest/api/3/resolution", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateResolution creates a new resolution.
func (c *Client) CreateResolution(resolution map[string]interface{}) (*Resolution, error) {
	var result Resolution
	if err := c.Post("/rest/api/3/resolution", resolution, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetResolution returns a resolution by ID.
func (c *Client) GetResolution(id string) (*Resolution, error) {
	var result Resolution
	if err := c.Get(fmt.Sprintf("/rest/api/3/resolution/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateResolution updates a resolution.
func (c *Client) UpdateResolution(id string, resolution map[string]interface{}) (*Resolution, error) {
	var result Resolution
	if err := c.Put(fmt.Sprintf("/rest/api/3/resolution/%s", id), resolution, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteResolution deletes a resolution.
func (c *Client) DeleteResolution(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/resolution/%s", id), nil)
}

// SearchResolutions searches for resolutions with pagination.
func (c *Client) SearchResolutions(startAt, maxResults int) (*PageBean[Resolution], error) {
	var result PageBean[Resolution]
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/resolution/search", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Statuses
// ============================================================================

// GetAllStatuses returns all statuses.
func (c *Client) GetAllStatuses() ([]StatusDetails, error) {
	var result []StatusDetails
	if err := c.Get("/rest/api/3/status", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetStatus returns a status by ID or name.
func (c *Client) GetStatus(idOrName string) (*StatusDetails, error) {
	var result StatusDetails
	if err := c.Get(fmt.Sprintf("/rest/api/3/status/%s", idOrName), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchStatuses searches for statuses with pagination.
func (c *Client) SearchStatuses(searchString string, startAt, maxResults int) (*PageBean[StatusDetails], error) {
	var result PageBean[StatusDetails]
	query := url.Values{}
	if searchString != "" {
		query.Set("searchString", searchString)
	}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/statuses/search", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStatusCategories returns all status categories.
func (c *Client) GetStatusCategories() ([]StatusCategory, error) {
	var result []StatusCategory
	if err := c.Get("/rest/api/3/statuscategory", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetStatusCategory returns a status category by ID or key.
func (c *Client) GetStatusCategory(idOrKey string) (*StatusCategory, error) {
	var result StatusCategory
	if err := c.Get(fmt.Sprintf("/rest/api/3/statuscategory/%s", idOrKey), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Labels
// ============================================================================

// GetLabels returns all labels with pagination.
func (c *Client) GetLabels(startAt, maxResults int) (*PageBean[string], error) {
	var result PageBean[string]
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/label", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Server & Config
// ============================================================================

// GetServerInfo returns Jira server information.
func (c *Client) GetServerInfo() (*ServerInfo, error) {
	var result ServerInfo
	if err := c.Get("/rest/api/3/serverInfo", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetConfiguration returns the Jira configuration.
func (c *Client) GetConfiguration() (*Configuration, error) {
	var result Configuration
	if err := c.Get("/rest/api/3/configuration", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAnnouncementBanner returns the announcement banner settings.
func (c *Client) GetAnnouncementBanner() (*AnnouncementBanner, error) {
	var result AnnouncementBanner
	if err := c.Get("/rest/api/3/announcementBanner", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAnnouncementBanner updates the announcement banner settings.
func (c *Client) SetAnnouncementBanner(banner *AnnouncementBanner) error {
	return c.Put("/rest/api/3/announcementBanner", banner, nil)
}

// GetAuditRecords returns audit records with pagination.
func (c *Client) GetAuditRecords(startAt, maxResults int) (*AuditRecords, error) {
	var result AuditRecords
	query := url.Values{}
	query.Set("offset", strconv.Itoa(startAt))
	query.Set("limit", strconv.Itoa(maxResults))
	if err := c.Get("/rest/api/3/auditing/record", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetApplicationRoles returns all application roles.
func (c *Client) GetApplicationRoles() ([]ApplicationRole, error) {
	var result []ApplicationRole
	if err := c.Get("/rest/api/3/applicationrole", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetApplicationRole returns an application role by key.
func (c *Client) GetApplicationRole(key string) (*ApplicationRole, error) {
	var result ApplicationRole
	if err := c.Get(fmt.Sprintf("/rest/api/3/applicationrole/%s", key), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Permissions
// ============================================================================

// GetMyPermissions returns the permissions for the current user.
func (c *Client) GetMyPermissions(projectKey, issueKey string) (map[string]UserPermission, error) {
	var result struct {
		Permissions map[string]UserPermission `json:"permissions"`
	}
	query := url.Values{}
	if projectKey != "" {
		query.Set("projectKey", projectKey)
	}
	if issueKey != "" {
		query.Set("issueKey", issueKey)
	}
	if err := c.Get("/rest/api/3/mypermissions", query, &result); err != nil {
		return nil, err
	}
	return result.Permissions, nil
}

// GetAllPermissions returns all permissions in the system.
func (c *Client) GetAllPermissions() (map[string]UserPermission, error) {
	var result struct {
		Permissions map[string]UserPermission `json:"permissions"`
	}
	if err := c.Get("/rest/api/3/permissions", nil, &result); err != nil {
		return nil, err
	}
	return result.Permissions, nil
}

// ============================================================================
// Tasks
// ============================================================================

// GetTask returns an async task result by ID.
func (c *Client) GetTask(taskId string) (*TaskResult, error) {
	var result TaskResult
	if err := c.Get(fmt.Sprintf("/rest/api/3/task/%s", taskId), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelTask cancels an async task.
func (c *Client) CancelTask(taskId string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/task/%s/cancel", taskId), nil, nil)
}


package jira

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// IssueTypeWithStatuses represents an issue type with its associated statuses.
type IssueTypeWithStatuses struct {
	Self     string          `json:"self,omitempty"`
	ID       string          `json:"id,omitempty"`
	Name     string          `json:"name,omitempty"`
	Subtask  bool            `json:"subtask,omitempty"`
	Statuses []StatusDetails `json:"statuses,omitempty"`
}

// ProjectType represents a Jira project type.
type ProjectType struct {
	Key                string `json:"key,omitempty"`
	FormattedKey       string `json:"formattedKey,omitempty"`
	DescriptionI18nKey string `json:"descriptionI18nKey,omitempty"`
	Icon               string `json:"icon,omitempty"`
	Color              string `json:"color,omitempty"`
}

// --- Projects ---

// GetAllProjects returns all projects visible to the user.
func (c *Client) GetAllProjects(expand string, recent int) ([]Project, error) {
	q := url.Values{}
	if expand != "" {
		q.Set("expand", expand)
	}
	if recent > 0 {
		q.Set("recent", strconv.Itoa(recent))
	}
	var projects []Project
	err := c.Get("/rest/api/3/project", q, &projects)
	return projects, err
}

// CreateProject creates a new project.
func (c *Client) CreateProject(project map[string]interface{}) (*Project, error) {
	var result Project
	err := c.Post("/rest/api/3/project", project, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProject returns a project by ID or key.
func (c *Client) GetProject(projectIdOrKey string, expand string) (*Project, error) {
	q := url.Values{}
	if expand != "" {
		q.Set("expand", expand)
	}
	var project Project
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey), q, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates a project.
func (c *Client) UpdateProject(projectIdOrKey string, project map[string]interface{}) (*Project, error) {
	var result Project
	err := c.Put(fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey), project, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(projectIdOrKey string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey), nil)
}

// SearchProjects searches for projects using a query string.
func (c *Client) SearchProjects(query string, startAt, maxResults int, expand string) (*PageBean[Project], error) {
	q := url.Values{}
	if query != "" {
		q.Set("query", query)
	}
	if startAt > 0 {
		q.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		q.Set("maxResults", strconv.Itoa(maxResults))
	}
	if expand != "" {
		q.Set("expand", expand)
	}
	var result PageBean[Project]
	err := c.Get("/rest/api/3/project/search", q, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRecentProjects returns recently accessed projects.
func (c *Client) GetRecentProjects(count int) ([]Project, error) {
	q := url.Values{}
	if count > 0 {
		q.Set("count", strconv.Itoa(count))
	}
	var projects []Project
	err := c.Get("/rest/api/3/project/recent", q, &projects)
	return projects, err
}

// GetProjectComponents returns all components for a project.
func (c *Client) GetProjectComponents(projectIdOrKey string) ([]ProjectComponent, error) {
	var components []ProjectComponent
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/components", projectIdOrKey), nil, &components)
	return components, err
}

// GetProjectVersions returns all versions for a project.
func (c *Client) GetProjectVersions(projectIdOrKey string) ([]Version, error) {
	var versions []Version
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/versions", projectIdOrKey), nil, &versions)
	return versions, err
}

// GetProjectVersionsPaginated returns paginated versions for a project.
func (c *Client) GetProjectVersionsPaginated(projectIdOrKey string, startAt, maxResults int) (*PageBean[Version], error) {
	q := url.Values{}
	if startAt > 0 {
		q.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		q.Set("maxResults", strconv.Itoa(maxResults))
	}
	var result PageBean[Version]
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/version", projectIdOrKey), q, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProjectStatuses returns the valid statuses for a project.
func (c *Client) GetProjectStatuses(projectIdOrKey string) ([]IssueTypeWithStatuses, error) {
	var statuses []IssueTypeWithStatuses
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/statuses", projectIdOrKey), nil, &statuses)
	return statuses, err
}

// GetProjectRoles returns all project roles for a project as a map of role name to URL.
func (c *Client) GetProjectRoles(projectIdOrKey string) (map[string]string, error) {
	var roles map[string]string
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/role", projectIdOrKey), nil, &roles)
	return roles, err
}

// GetProjectRole returns a specific project role.
func (c *Client) GetProjectRole(projectIdOrKey string, roleId int) (*ProjectRole, error) {
	var role ProjectRole
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/role/%d", projectIdOrKey, roleId), nil, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ArchiveProject archives a project.
func (c *Client) ArchiveProject(projectIdOrKey string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/project/%s/archive", projectIdOrKey), nil, nil)
}

// RestoreProject restores an archived project.
func (c *Client) RestoreProject(projectIdOrKey string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/project/%s/restore", projectIdOrKey), nil, nil)
}

// GetProjectProperties returns all properties for a project.
func (c *Client) GetProjectProperties(projectIdOrKey string) ([]EntityProperty, error) {
	var result struct {
		Keys []EntityProperty `json:"keys"`
	}
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/properties", projectIdOrKey), nil, &result)
	return result.Keys, err
}

// SetProjectProperty sets a property on a project.
func (c *Client) SetProjectProperty(projectIdOrKey, key string, value interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/project/%s/properties/%s", projectIdOrKey, key), value, nil)
}

// DeleteProjectProperty deletes a property from a project.
func (c *Client) DeleteProjectProperty(projectIdOrKey, key string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/project/%s/properties/%s", projectIdOrKey, key), nil)
}

// GetProjectFeatures returns the features for a project.
func (c *Client) GetProjectFeatures(projectIdOrKey string) (*ProjectFeaturesResponse, error) {
	var result ProjectFeaturesResponse
	err := c.Get(fmt.Sprintf("/rest/api/3/project/%s/features", projectIdOrKey), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetProjectFeature sets the state of a project feature.
func (c *Client) SetProjectFeature(projectIdOrKey, featureKey, state string) (*ProjectFeature, error) {
	body := map[string]string{"state": state}
	var result ProjectFeature
	err := c.Put(fmt.Sprintf("/rest/api/3/project/%s/features/%s", projectIdOrKey, featureKey), body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Project Categories ---

// GetProjectCategories returns all project categories.
func (c *Client) GetProjectCategories() ([]ProjectCategory, error) {
	var categories []ProjectCategory
	err := c.Get("/rest/api/3/projectCategory", nil, &categories)
	return categories, err
}

// CreateProjectCategory creates a new project category.
func (c *Client) CreateProjectCategory(cat *ProjectCategory) (*ProjectCategory, error) {
	var result ProjectCategory
	err := c.Post("/rest/api/3/projectCategory", cat, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProjectCategory returns a project category by ID.
func (c *Client) GetProjectCategory(id string) (*ProjectCategory, error) {
	var result ProjectCategory
	err := c.Get(fmt.Sprintf("/rest/api/3/projectCategory/%s", id), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProjectCategory updates a project category.
func (c *Client) UpdateProjectCategory(id string, cat *ProjectCategory) (*ProjectCategory, error) {
	var result ProjectCategory
	err := c.Put(fmt.Sprintf("/rest/api/3/projectCategory/%s", id), cat, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProjectCategory deletes a project category.
func (c *Client) DeleteProjectCategory(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/projectCategory/%s", id), nil)
}

// --- Project Validation ---

// ValidateProjectKey validates a project key.
func (c *Client) ValidateProjectKey(key string) (map[string]interface{}, error) {
	q := url.Values{}
	q.Set("key", key)
	var result map[string]interface{}
	err := c.Get("/rest/api/3/projectvalidate/key", q, &result)
	return result, err
}

// GetValidProjectKey returns a valid project key based on the provided key.
func (c *Client) GetValidProjectKey(key string) (string, error) {
	q := url.Values{}
	q.Set("key", key)
	data, err := c.GetRaw("/rest/api/3/projectvalidate/validProjectKey", q)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(data), "\" \n"), nil
}

// GetValidProjectName returns a valid project name based on the provided name.
func (c *Client) GetValidProjectName(name string) (string, error) {
	q := url.Values{}
	q.Set("name", name)
	data, err := c.GetRaw("/rest/api/3/projectvalidate/validProjectName", q)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(data), "\" \n"), nil
}

// --- Project Types ---

// GetAllProjectTypes returns all project types.
func (c *Client) GetAllProjectTypes() ([]ProjectType, error) {
	var types []ProjectType
	err := c.Get("/rest/api/3/project/type", nil, &types)
	return types, err
}

// GetAccessibleProjectTypes returns all accessible project types.
func (c *Client) GetAccessibleProjectTypes() ([]ProjectType, error) {
	var types []ProjectType
	err := c.Get("/rest/api/3/project/type/accessible", nil, &types)
	return types, err
}

// --- Search ---

// SearchJQL searches for issues using JQL via GET.
func (c *Client) SearchJQL(jql string, startAt, maxResults int, fields []string, expand []string) (*SearchResults, error) {
	q := url.Values{}
	if jql != "" {
		q.Set("jql", jql)
	}
	if startAt > 0 {
		q.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		q.Set("maxResults", strconv.Itoa(maxResults))
	}
	if len(fields) > 0 {
		q.Set("fields", strings.Join(fields, ","))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var result SearchResults
	err := c.Get("/rest/api/3/search", q, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchJQLPost searches for issues using JQL via POST.
func (c *Client) SearchJQLPost(req *SearchRequest) (*SearchResults, error) {
	var result SearchResults
	err := c.Post("/rest/api/3/search", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Filters ---

// CreateFilter creates a new filter.
func (c *Client) CreateFilter(filter *Filter) (*Filter, error) {
	var result Filter
	err := c.Post("/rest/api/3/filter", filter, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFilter returns a filter by ID.
func (c *Client) GetFilter(id string) (*Filter, error) {
	var result Filter
	err := c.Get(fmt.Sprintf("/rest/api/3/filter/%s", id), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFilter updates a filter.
func (c *Client) UpdateFilter(id string, filter *Filter) (*Filter, error) {
	var result Filter
	err := c.Put(fmt.Sprintf("/rest/api/3/filter/%s", id), filter, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFilter deletes a filter.
func (c *Client) DeleteFilter(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/filter/%s", id), nil)
}

// GetFavouriteFilters returns the user's favourite filters.
func (c *Client) GetFavouriteFilters() ([]Filter, error) {
	var filters []Filter
	err := c.Get("/rest/api/3/filter/favourite", nil, &filters)
	return filters, err
}

// GetMyFilters returns the user's own filters.
func (c *Client) GetMyFilters() ([]Filter, error) {
	var filters []Filter
	err := c.Get("/rest/api/3/filter/my", nil, &filters)
	return filters, err
}

// SearchFilters searches for filters by name.
func (c *Client) SearchFilters(name string, startAt, maxResults int) (*PageBean[Filter], error) {
	q := url.Values{}
	if name != "" {
		q.Set("filterName", name)
	}
	if startAt > 0 {
		q.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		q.Set("maxResults", strconv.Itoa(maxResults))
	}
	var result PageBean[Filter]
	err := c.Get("/rest/api/3/filter/search", q, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

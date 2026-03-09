package jira

import (
	"fmt"
	"net/url"
	"strconv"
)

// FieldContext represents a field context.
type FieldContext struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	IsGlobalContext bool   `json:"isGlobalContext,omitempty"`
	IsAnyIssueType  bool   `json:"isAnyIssueType,omitempty"`
}

// --- Dashboards ---

// GetDashboards returns a paginated list of dashboards.
func (c *Client) GetDashboards(startAt, maxResults int) (*DashboardList, error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result DashboardList
	if err := c.Get("/rest/api/3/dashboard", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDashboard creates a new dashboard.
func (c *Client) CreateDashboard(dashboard map[string]interface{}) (*Dashboard, error) {
	var result Dashboard
	if err := c.Post("/rest/api/3/dashboard", dashboard, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDashboard returns a dashboard by ID.
func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	var result Dashboard
	if err := c.Get(fmt.Sprintf("/rest/api/3/dashboard/%s", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDashboard updates a dashboard.
func (c *Client) UpdateDashboard(id string, dashboard map[string]interface{}) (*Dashboard, error) {
	var result Dashboard
	if err := c.Put(fmt.Sprintf("/rest/api/3/dashboard/%s", id), dashboard, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDashboard deletes a dashboard by ID.
func (c *Client) DeleteDashboard(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/dashboard/%s", id), nil)
}

// CopyDashboard copies a dashboard.
func (c *Client) CopyDashboard(id string, dashboard map[string]interface{}) (*Dashboard, error) {
	var result Dashboard
	if err := c.Post(fmt.Sprintf("/rest/api/3/dashboard/%s/copy", id), dashboard, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchDashboards searches for dashboards by name.
func (c *Client) SearchDashboards(name string, startAt, maxResults int) (*PageBean[Dashboard], error) {
	query := url.Values{}
	if name != "" {
		query.Set("dashboardName", name)
	}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Dashboard]
	if err := c.Get("/rest/api/3/dashboard/search", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDashboardGadgets returns the gadgets on a dashboard.
func (c *Client) GetDashboardGadgets(dashboardId string) (*DashboardGadgetList, error) {
	var result DashboardGadgetList
	if err := c.Get(fmt.Sprintf("/rest/api/3/dashboard/%s/gadget", dashboardId), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddDashboardGadget adds a gadget to a dashboard.
func (c *Client) AddDashboardGadget(dashboardId string, gadget map[string]interface{}) (*DashboardGadget, error) {
	var result DashboardGadget
	if err := c.Post(fmt.Sprintf("/rest/api/3/dashboard/%s/gadget", dashboardId), gadget, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDashboardGadget updates a gadget on a dashboard.
func (c *Client) UpdateDashboardGadget(dashboardId string, gadgetId string, gadget map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/dashboard/%s/gadget/%s", dashboardId, gadgetId), gadget, nil)
}

// RemoveDashboardGadget removes a gadget from a dashboard.
func (c *Client) RemoveDashboardGadget(dashboardId string, gadgetId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/dashboard/%s/gadget/%s", dashboardId, gadgetId), nil)
}

// --- Fields ---

// GetFields returns all fields.
func (c *Client) GetFields() ([]Field, error) {
	var result []Field
	if err := c.Get("/rest/api/3/field", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateCustomField creates a custom field.
func (c *Client) CreateCustomField(field map[string]interface{}) (*Field, error) {
	var result Field
	if err := c.Post("/rest/api/3/field", field, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCustomField updates a custom field.
func (c *Client) UpdateCustomField(fieldId string, field map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/field/%s", fieldId), field, nil)
}

// DeleteCustomField deletes a custom field.
func (c *Client) DeleteCustomField(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/field/%s", id), nil)
}

// SearchFields searches for fields.
func (c *Client) SearchFields(query string, startAt, maxResults int) (*PageBean[Field], error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	params.Set("startAt", strconv.Itoa(startAt))
	params.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Field]
	if err := c.Get("/rest/api/3/field/search", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFieldContexts returns contexts for a field.
func (c *Client) GetFieldContexts(fieldId string, startAt, maxResults int) (*PageBean[FieldContext], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[FieldContext]
	if err := c.Get(fmt.Sprintf("/rest/api/3/field/%s/context", fieldId), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TrashCustomField moves a custom field to trash.
func (c *Client) TrashCustomField(id string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/field/%s/trash", id), nil, nil)
}

// RestoreCustomField restores a custom field from trash.
func (c *Client) RestoreCustomField(id string) error {
	return c.Post(fmt.Sprintf("/rest/api/3/field/%s/restore", id), nil, nil)
}

// GetTrashedFields returns trashed fields.
func (c *Client) GetTrashedFields(startAt, maxResults int) (*PageBean[Field], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Field]
	if err := c.Get("/rest/api/3/field/search/trashed", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Screens ---

// GetScreens returns a paginated list of screens.
func (c *Client) GetScreens(startAt, maxResults int) (*PageBean[Screen], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Screen]
	if err := c.Get("/rest/api/3/screens", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateScreen creates a screen.
func (c *Client) CreateScreen(screen map[string]interface{}) (*Screen, error) {
	var result Screen
	if err := c.Post("/rest/api/3/screens", screen, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateScreen updates a screen.
func (c *Client) UpdateScreen(screenId int, screen map[string]interface{}) (*Screen, error) {
	var result Screen
	if err := c.Put(fmt.Sprintf("/rest/api/3/screens/%d", screenId), screen, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteScreen deletes a screen.
func (c *Client) DeleteScreen(screenId int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/screens/%d", screenId), nil)
}

// GetScreenTabs returns the tabs for a screen.
func (c *Client) GetScreenTabs(screenId int) ([]ScreenTab, error) {
	var result []ScreenTab
	if err := c.Get(fmt.Sprintf("/rest/api/3/screens/%d/tabs", screenId), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateScreenTab creates a tab on a screen.
func (c *Client) CreateScreenTab(screenId int, tab map[string]interface{}) (*ScreenTab, error) {
	var result ScreenTab
	if err := c.Post(fmt.Sprintf("/rest/api/3/screens/%d/tabs", screenId), tab, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateScreenTab updates a tab on a screen.
func (c *Client) UpdateScreenTab(screenId, tabId int, tab map[string]interface{}) (*ScreenTab, error) {
	var result ScreenTab
	if err := c.Put(fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d", screenId, tabId), tab, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteScreenTab deletes a tab from a screen.
func (c *Client) DeleteScreenTab(screenId, tabId int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d", screenId, tabId), nil)
}

// GetScreenTabFields returns the fields on a screen tab.
func (c *Client) GetScreenTabFields(screenId, tabId int) ([]ScreenField, error) {
	var result []ScreenField
	if err := c.Get(fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields", screenId, tabId), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddScreenTabField adds a field to a screen tab.
func (c *Client) AddScreenTabField(screenId, tabId int, fieldId string) (*ScreenField, error) {
	body := map[string]interface{}{"fieldId": fieldId}
	var result ScreenField
	if err := c.Post(fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields", screenId, tabId), body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveScreenTabField removes a field from a screen tab.
func (c *Client) RemoveScreenTabField(screenId, tabId int, fieldId string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields/%s", screenId, tabId, fieldId), nil)
}

// --- Screen Schemes ---

// GetScreenSchemes returns a paginated list of screen schemes.
func (c *Client) GetScreenSchemes(startAt, maxResults int) (*PageBean[ScreenScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[ScreenScheme]
	if err := c.Get("/rest/api/3/screenscheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateScreenScheme creates a screen scheme.
func (c *Client) CreateScreenScheme(scheme map[string]interface{}) (*ScreenScheme, error) {
	var result ScreenScheme
	if err := c.Post("/rest/api/3/screenscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateScreenScheme updates a screen scheme.
func (c *Client) UpdateScreenScheme(id int, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/screenscheme/%d", id), scheme, nil)
}

// DeleteScreenScheme deletes a screen scheme.
func (c *Client) DeleteScreenScheme(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/screenscheme/%d", id), nil)
}

// --- Workflows ---

// GetWorkflows returns all workflows.
func (c *Client) GetWorkflows() ([]Workflow, error) {
	var result []Workflow
	if err := c.Get("/rest/api/3/workflow", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SearchWorkflows searches for workflows.
func (c *Client) SearchWorkflows(query string, startAt, maxResults int) (*PageBean[Workflow], error) {
	params := url.Values{}
	if query != "" {
		params.Set("queryString", query)
	}
	params.Set("startAt", strconv.Itoa(startAt))
	params.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Workflow]
	if err := c.Get("/rest/api/3/workflow/search", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Workflow Schemes ---

// GetWorkflowSchemes returns a paginated list of workflow schemes.
func (c *Client) GetWorkflowSchemes(startAt, maxResults int) (*PageBean[WorkflowScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[WorkflowScheme]
	if err := c.Get("/rest/api/3/workflowscheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateWorkflowScheme creates a workflow scheme.
func (c *Client) CreateWorkflowScheme(scheme map[string]interface{}) (*WorkflowScheme, error) {
	var result WorkflowScheme
	if err := c.Post("/rest/api/3/workflowscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWorkflowScheme returns a workflow scheme by ID.
func (c *Client) GetWorkflowScheme(id int) (*WorkflowScheme, error) {
	var result WorkflowScheme
	if err := c.Get(fmt.Sprintf("/rest/api/3/workflowscheme/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateWorkflowScheme updates a workflow scheme.
func (c *Client) UpdateWorkflowScheme(id int, scheme map[string]interface{}) (*WorkflowScheme, error) {
	var result WorkflowScheme
	if err := c.Put(fmt.Sprintf("/rest/api/3/workflowscheme/%d", id), scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteWorkflowScheme deletes a workflow scheme.
func (c *Client) DeleteWorkflowScheme(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/workflowscheme/%d", id), nil)
}

// --- Issue Type Schemes ---

// GetIssueTypeSchemes returns a paginated list of issue type schemes.
func (c *Client) GetIssueTypeSchemes(startAt, maxResults int) (*PageBean[IssueTypeScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[IssueTypeScheme]
	if err := c.Get("/rest/api/3/issuetypescheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateIssueTypeScheme creates an issue type scheme.
func (c *Client) CreateIssueTypeScheme(scheme map[string]interface{}) (*IssueTypeScheme, error) {
	var result IssueTypeScheme
	if err := c.Post("/rest/api/3/issuetypescheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueTypeScheme updates an issue type scheme.
func (c *Client) UpdateIssueTypeScheme(id string, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/issuetypescheme/%s", id), scheme, nil)
}

// DeleteIssueTypeScheme deletes an issue type scheme.
func (c *Client) DeleteIssueTypeScheme(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issuetypescheme/%s", id), nil)
}

// --- Issue Type Screen Schemes ---

// GetIssueTypeScreenSchemes returns a paginated list of issue type screen schemes.
func (c *Client) GetIssueTypeScreenSchemes(startAt, maxResults int) (*PageBean[IssueTypeScreenScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[IssueTypeScreenScheme]
	if err := c.Get("/rest/api/3/issuetypescreenscheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateIssueTypeScreenScheme creates an issue type screen scheme.
func (c *Client) CreateIssueTypeScreenScheme(scheme map[string]interface{}) (*IssueTypeScreenScheme, error) {
	var result IssueTypeScreenScheme
	if err := c.Post("/rest/api/3/issuetypescreenscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueTypeScreenScheme updates an issue type screen scheme.
func (c *Client) UpdateIssueTypeScreenScheme(id string, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/issuetypescreenscheme/%s", id), scheme, nil)
}

// DeleteIssueTypeScreenScheme deletes an issue type screen scheme.
func (c *Client) DeleteIssueTypeScreenScheme(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issuetypescreenscheme/%s", id), nil)
}

// --- Field Configurations ---

// GetFieldConfigurations returns a paginated list of field configurations.
func (c *Client) GetFieldConfigurations(startAt, maxResults int) (*PageBean[FieldConfiguration], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[FieldConfiguration]
	if err := c.Get("/rest/api/3/fieldconfiguration", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFieldConfiguration creates a field configuration.
func (c *Client) CreateFieldConfiguration(config map[string]interface{}) (*FieldConfiguration, error) {
	var result FieldConfiguration
	if err := c.Post("/rest/api/3/fieldconfiguration", config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFieldConfiguration updates a field configuration.
func (c *Client) UpdateFieldConfiguration(id int, config map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/fieldconfiguration/%d", id), config, nil)
}

// DeleteFieldConfiguration deletes a field configuration.
func (c *Client) DeleteFieldConfiguration(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/fieldconfiguration/%d", id), nil)
}

// --- Field Configuration Schemes ---

// GetFieldConfigurationSchemes returns a paginated list of field configuration schemes.
func (c *Client) GetFieldConfigurationSchemes(startAt, maxResults int) (*PageBean[FieldConfigurationScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[FieldConfigurationScheme]
	if err := c.Get("/rest/api/3/fieldconfigurationscheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFieldConfigurationScheme creates a field configuration scheme.
func (c *Client) CreateFieldConfigurationScheme(scheme map[string]interface{}) (*FieldConfigurationScheme, error) {
	var result FieldConfigurationScheme
	if err := c.Post("/rest/api/3/fieldconfigurationscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFieldConfigurationScheme updates a field configuration scheme.
func (c *Client) UpdateFieldConfigurationScheme(id string, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/fieldconfigurationscheme/%s", id), scheme, nil)
}

// DeleteFieldConfigurationScheme deletes a field configuration scheme.
func (c *Client) DeleteFieldConfigurationScheme(id string) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/fieldconfigurationscheme/%s", id), nil)
}

// --- Permission Schemes ---

// permissionSchemesResponse is the wrapper for the permission schemes list endpoint.
type permissionSchemesResponse struct {
	PermissionSchemes []PermissionScheme `json:"permissionSchemes"`
}

// GetPermissionSchemes returns all permission schemes.
func (c *Client) GetPermissionSchemes() ([]PermissionScheme, error) {
	var result permissionSchemesResponse
	if err := c.Get("/rest/api/3/permissionscheme", nil, &result); err != nil {
		return nil, err
	}
	return result.PermissionSchemes, nil
}

// CreatePermissionScheme creates a permission scheme.
func (c *Client) CreatePermissionScheme(scheme map[string]interface{}) (*PermissionScheme, error) {
	var result PermissionScheme
	if err := c.Post("/rest/api/3/permissionscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPermissionScheme returns a permission scheme by ID.
func (c *Client) GetPermissionScheme(id int) (*PermissionScheme, error) {
	var result PermissionScheme
	if err := c.Get(fmt.Sprintf("/rest/api/3/permissionscheme/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePermissionScheme updates a permission scheme.
func (c *Client) UpdatePermissionScheme(id int, scheme map[string]interface{}) (*PermissionScheme, error) {
	var result PermissionScheme
	if err := c.Put(fmt.Sprintf("/rest/api/3/permissionscheme/%d", id), scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePermissionScheme deletes a permission scheme.
func (c *Client) DeletePermissionScheme(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/permissionscheme/%d", id), nil)
}

// --- Notification Schemes ---

// GetNotificationSchemes returns a paginated list of notification schemes.
func (c *Client) GetNotificationSchemes(startAt, maxResults int) (*PageBean[NotificationScheme], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[NotificationScheme]
	if err := c.Get("/rest/api/3/notificationscheme", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateNotificationScheme creates a notification scheme.
func (c *Client) CreateNotificationScheme(scheme map[string]interface{}) (*NotificationScheme, error) {
	var result NotificationScheme
	if err := c.Post("/rest/api/3/notificationscheme", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetNotificationScheme returns a notification scheme by ID.
func (c *Client) GetNotificationScheme(id int) (*NotificationScheme, error) {
	var result NotificationScheme
	if err := c.Get(fmt.Sprintf("/rest/api/3/notificationscheme/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateNotificationScheme updates a notification scheme.
func (c *Client) UpdateNotificationScheme(id int, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/notificationscheme/%d", id), scheme, nil)
}

// DeleteNotificationScheme deletes a notification scheme.
func (c *Client) DeleteNotificationScheme(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/notificationscheme/%d", id), nil)
}

// --- Issue Security Schemes ---

// issueSecuritySchemesResponse is the wrapper for the issue security schemes list endpoint.
type issueSecuritySchemesResponse struct {
	IssueSecuritySchemes []IssueSecurityScheme `json:"issueSecuritySchemes"`
}

// GetIssueSecuritySchemes returns all issue security schemes.
func (c *Client) GetIssueSecuritySchemes() ([]IssueSecurityScheme, error) {
	var result issueSecuritySchemesResponse
	if err := c.Get("/rest/api/3/issuesecurityschemes", nil, &result); err != nil {
		return nil, err
	}
	return result.IssueSecuritySchemes, nil
}

// CreateIssueSecurityScheme creates an issue security scheme.
func (c *Client) CreateIssueSecurityScheme(scheme map[string]interface{}) (*IssueSecurityScheme, error) {
	var result IssueSecurityScheme
	if err := c.Post("/rest/api/3/issuesecurityschemes", scheme, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueSecurityScheme returns an issue security scheme by ID.
func (c *Client) GetIssueSecurityScheme(id int) (*IssueSecurityScheme, error) {
	var result IssueSecurityScheme
	if err := c.Get(fmt.Sprintf("/rest/api/3/issuesecurityschemes/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateIssueSecurityScheme updates an issue security scheme.
func (c *Client) UpdateIssueSecurityScheme(id int, scheme map[string]interface{}) error {
	return c.Put(fmt.Sprintf("/rest/api/3/issuesecurityschemes/%d", id), scheme, nil)
}

// DeleteIssueSecurityScheme deletes an issue security scheme.
func (c *Client) DeleteIssueSecurityScheme(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/issuesecurityschemes/%d", id), nil)
}

// --- Roles ---

// GetAllRoles returns all project roles.
func (c *Client) GetAllRoles() ([]ProjectRole, error) {
	var result []ProjectRole
	if err := c.Get("/rest/api/3/role", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateRole creates a project role.
func (c *Client) CreateRole(role map[string]interface{}) (*ProjectRole, error) {
	var result ProjectRole
	if err := c.Post("/rest/api/3/role", role, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRole returns a project role by ID.
func (c *Client) GetRole(id int) (*ProjectRole, error) {
	var result ProjectRole
	if err := c.Get(fmt.Sprintf("/rest/api/3/role/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRole updates a project role.
func (c *Client) UpdateRole(id int, role map[string]interface{}) (*ProjectRole, error) {
	var result ProjectRole
	if err := c.Put(fmt.Sprintf("/rest/api/3/role/%d", id), role, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteRole deletes a project role.
func (c *Client) DeleteRole(id int) error {
	return c.Delete(fmt.Sprintf("/rest/api/3/role/%d", id), nil)
}

// --- Webhooks ---

// GetWebhooks returns a paginated list of webhooks.
func (c *Client) GetWebhooks(startAt, maxResults int) (*PageBean[Webhook], error) {
	query := url.Values{}
	query.Set("startAt", strconv.Itoa(startAt))
	query.Set("maxResults", strconv.Itoa(maxResults))
	var result PageBean[Webhook]
	if err := c.Get("/rest/api/3/webhook", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterWebhooks registers webhooks.
func (c *Client) RegisterWebhooks(webhooks map[string]interface{}) ([]Webhook, error) {
	var result struct {
		WebhookRegistrationResult []Webhook `json:"webhookRegistrationResult"`
	}
	if err := c.Post("/rest/api/3/webhook", webhooks, &result); err != nil {
		return nil, err
	}
	return result.WebhookRegistrationResult, nil
}

// DeleteWebhooks deletes webhooks by IDs.
func (c *Client) DeleteWebhooks(webhookIds []int) error {
	body := map[string]interface{}{"webhookIds": webhookIds}
	return c.DeleteWithBody("/rest/api/3/webhook", body, nil)
}

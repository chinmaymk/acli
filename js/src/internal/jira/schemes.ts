import { get, post, put, del, deleteWithBody } from './client.js';
import type { JiraClient } from './client.js';

// --- Dashboards ---

/**
 * Returns a paginated list of dashboards.
 */
export async function getDashboards(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/dashboard', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a new dashboard.
 */
export async function createDashboard(client: JiraClient, dashboard: unknown): Promise<any> {
  return post(client, '/rest/api/3/dashboard', dashboard);
}

/**
 * Returns a dashboard by ID.
 */
export async function getDashboard(client: JiraClient, id: string): Promise<any> {
  return get(client, `/rest/api/3/dashboard/${id}`);
}

/**
 * Updates a dashboard.
 */
export async function updateDashboard(
  client: JiraClient,
  id: string,
  dashboard: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/dashboard/${id}`, dashboard);
}

/**
 * Deletes a dashboard by ID.
 */
export async function deleteDashboard(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/dashboard/${id}`);
}

/**
 * Copies a dashboard.
 */
export async function copyDashboard(
  client: JiraClient,
  id: string,
  dashboard: unknown,
): Promise<any> {
  return post(client, `/rest/api/3/dashboard/${id}/copy`, dashboard);
}

/**
 * Searches for dashboards by name.
 */
export async function searchDashboards(
  client: JiraClient,
  name: string,
  startAt: number,
  maxResults: number,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (name) query.dashboardName = name;
  return get(client, '/rest/api/3/dashboard/search', query as Record<string, string>);
}

/**
 * Returns the gadgets on a dashboard.
 */
export async function getDashboardGadgets(client: JiraClient, dashboardId: string): Promise<any> {
  return get(client, `/rest/api/3/dashboard/${dashboardId}/gadget`);
}

/**
 * Adds a gadget to a dashboard.
 */
export async function addDashboardGadget(
  client: JiraClient,
  dashboardId: string,
  gadget: unknown,
): Promise<any> {
  return post(client, `/rest/api/3/dashboard/${dashboardId}/gadget`, gadget);
}

/**
 * Updates a gadget on a dashboard.
 */
export async function updateDashboardGadget(
  client: JiraClient,
  dashboardId: string,
  gadgetId: string,
  gadget: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/dashboard/${dashboardId}/gadget/${gadgetId}`, gadget);
}

/**
 * Removes a gadget from a dashboard.
 */
export async function removeDashboardGadget(
  client: JiraClient,
  dashboardId: string,
  gadgetId: string,
): Promise<void> {
  return del(client, `/rest/api/3/dashboard/${dashboardId}/gadget/${gadgetId}`);
}

// --- Fields ---

/**
 * Returns all fields.
 */
export async function getFields(client: JiraClient): Promise<any> {
  return get(client, '/rest/api/3/field');
}

/**
 * Creates a custom field.
 */
export async function createCustomField(client: JiraClient, field: unknown): Promise<any> {
  return post(client, '/rest/api/3/field', field);
}

/**
 * Updates a custom field.
 */
export async function updateCustomField(
  client: JiraClient,
  fieldId: string,
  field: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/field/${fieldId}`, field);
}

/**
 * Deletes a custom field.
 */
export async function deleteCustomField(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/field/${id}`);
}

/**
 * Searches for fields.
 */
export async function searchFields(
  client: JiraClient,
  query: string,
  startAt: number,
  maxResults: number,
): Promise<any> {
  const params: Record<string, string | number> = { startAt, maxResults };
  if (query) params.query = query;
  return get(client, '/rest/api/3/field/search', params as Record<string, string>);
}

/**
 * Returns contexts for a field.
 */
export async function getFieldContexts(
  client: JiraClient,
  fieldId: string,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, `/rest/api/3/field/${fieldId}/context`, {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Moves a custom field to trash.
 */
export async function trashCustomField(client: JiraClient, id: string): Promise<any> {
  return post(client, `/rest/api/3/field/${id}/trash`, null);
}

/**
 * Restores a custom field from trash.
 */
export async function restoreCustomField(client: JiraClient, id: string): Promise<any> {
  return post(client, `/rest/api/3/field/${id}/restore`, null);
}

/**
 * Returns trashed fields.
 */
export async function getTrashedFields(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/field/search/trashed', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

// --- Screens ---

/**
 * Returns a paginated list of screens.
 */
export async function getScreens(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/screens', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a screen.
 */
export async function createScreen(client: JiraClient, screen: unknown): Promise<any> {
  return post(client, '/rest/api/3/screens', screen);
}

/**
 * Updates a screen.
 */
export async function updateScreen(
  client: JiraClient,
  screenId: number,
  screen: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/screens/${screenId}`, screen);
}

/**
 * Deletes a screen.
 */
export async function deleteScreen(client: JiraClient, screenId: number): Promise<void> {
  return del(client, `/rest/api/3/screens/${screenId}`);
}

/**
 * Returns the tabs for a screen.
 */
export async function getScreenTabs(client: JiraClient, screenId: number): Promise<any> {
  return get(client, `/rest/api/3/screens/${screenId}/tabs`);
}

/**
 * Creates a tab on a screen.
 */
export async function createScreenTab(
  client: JiraClient,
  screenId: number,
  tab: unknown,
): Promise<any> {
  return post(client, `/rest/api/3/screens/${screenId}/tabs`, tab);
}

/**
 * Updates a tab on a screen.
 */
export async function updateScreenTab(
  client: JiraClient,
  screenId: number,
  tabId: number,
  tab: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}`, tab);
}

/**
 * Deletes a tab from a screen.
 */
export async function deleteScreenTab(
  client: JiraClient,
  screenId: number,
  tabId: number,
): Promise<void> {
  return del(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}`);
}

/**
 * Returns the fields on a screen tab.
 */
export async function getScreenTabFields(
  client: JiraClient,
  screenId: number,
  tabId: number,
): Promise<any> {
  return get(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields`);
}

/**
 * Adds a field to a screen tab.
 */
export async function addScreenTabField(
  client: JiraClient,
  screenId: number,
  tabId: number,
  fieldId: string,
): Promise<any> {
  return post(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields`, { fieldId });
}

/**
 * Removes a field from a screen tab.
 */
export async function removeScreenTabField(
  client: JiraClient,
  screenId: number,
  tabId: number,
  fieldId: string,
): Promise<void> {
  return del(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields/${fieldId}`);
}

// --- Screen Schemes ---

/**
 * Returns a paginated list of screen schemes.
 */
export async function getScreenSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/screenscheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a screen scheme.
 */
export async function createScreenScheme(client: JiraClient, scheme: unknown): Promise<any> {
  return post(client, '/rest/api/3/screenscheme', scheme);
}

/**
 * Updates a screen scheme.
 */
export async function updateScreenScheme(
  client: JiraClient,
  id: number,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/screenscheme/${id}`, scheme);
}

/**
 * Deletes a screen scheme.
 */
export async function deleteScreenScheme(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/screenscheme/${id}`);
}

// --- Workflows ---

/**
 * Returns all workflows.
 */
export async function getWorkflows(client: JiraClient): Promise<any> {
  return get(client, '/rest/api/3/workflow');
}

/**
 * Searches for workflows.
 */
export async function searchWorkflows(
  client: JiraClient,
  query: string,
  startAt: number,
  maxResults: number,
): Promise<any> {
  const params: Record<string, string | number> = { startAt, maxResults };
  if (query) params.queryString = query;
  return get(client, '/rest/api/3/workflow/search', params as Record<string, string>);
}

// --- Workflow Schemes ---

/**
 * Returns a paginated list of workflow schemes.
 */
export async function getWorkflowSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/workflowscheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a workflow scheme.
 */
export async function createWorkflowScheme(client: JiraClient, scheme: unknown): Promise<any> {
  return post(client, '/rest/api/3/workflowscheme', scheme);
}

/**
 * Returns a workflow scheme by ID.
 */
export async function getWorkflowScheme(client: JiraClient, id: number): Promise<any> {
  return get(client, `/rest/api/3/workflowscheme/${id}`);
}

/**
 * Updates a workflow scheme.
 */
export async function updateWorkflowScheme(
  client: JiraClient,
  id: number,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/workflowscheme/${id}`, scheme);
}

/**
 * Deletes a workflow scheme.
 */
export async function deleteWorkflowScheme(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/workflowscheme/${id}`);
}

// --- Issue Type Schemes ---

/**
 * Returns a paginated list of issue type schemes.
 */
export async function getIssueTypeSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/issuetypescheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates an issue type scheme.
 */
export async function createIssueTypeScheme(client: JiraClient, scheme: unknown): Promise<any> {
  return post(client, '/rest/api/3/issuetypescheme', scheme);
}

/**
 * Updates an issue type scheme.
 */
export async function updateIssueTypeScheme(
  client: JiraClient,
  id: string,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/issuetypescheme/${id}`, scheme);
}

/**
 * Deletes an issue type scheme.
 */
export async function deleteIssueTypeScheme(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/issuetypescheme/${id}`);
}

// --- Issue Type Screen Schemes ---

/**
 * Returns a paginated list of issue type screen schemes.
 */
export async function getIssueTypeScreenSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/issuetypescreenscheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates an issue type screen scheme.
 */
export async function createIssueTypeScreenScheme(
  client: JiraClient,
  scheme: unknown,
): Promise<any> {
  return post(client, '/rest/api/3/issuetypescreenscheme', scheme);
}

/**
 * Updates an issue type screen scheme.
 */
export async function updateIssueTypeScreenScheme(
  client: JiraClient,
  id: string,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/issuetypescreenscheme/${id}`, scheme);
}

/**
 * Deletes an issue type screen scheme.
 */
export async function deleteIssueTypeScreenScheme(
  client: JiraClient,
  id: string,
): Promise<void> {
  return del(client, `/rest/api/3/issuetypescreenscheme/${id}`);
}

// --- Field Configurations ---

/**
 * Returns a paginated list of field configurations.
 */
export async function getFieldConfigurations(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/fieldconfiguration', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a field configuration.
 */
export async function createFieldConfiguration(
  client: JiraClient,
  config: unknown,
): Promise<any> {
  return post(client, '/rest/api/3/fieldconfiguration', config);
}

/**
 * Updates a field configuration.
 */
export async function updateFieldConfiguration(
  client: JiraClient,
  id: number,
  config: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/fieldconfiguration/${id}`, config);
}

/**
 * Deletes a field configuration.
 */
export async function deleteFieldConfiguration(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/fieldconfiguration/${id}`);
}

// --- Field Configuration Schemes ---

/**
 * Returns a paginated list of field configuration schemes.
 */
export async function getFieldConfigurationSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/fieldconfigurationscheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a field configuration scheme.
 */
export async function createFieldConfigurationScheme(
  client: JiraClient,
  scheme: unknown,
): Promise<any> {
  return post(client, '/rest/api/3/fieldconfigurationscheme', scheme);
}

/**
 * Updates a field configuration scheme.
 */
export async function updateFieldConfigurationScheme(
  client: JiraClient,
  id: string,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/fieldconfigurationscheme/${id}`, scheme);
}

/**
 * Deletes a field configuration scheme.
 */
export async function deleteFieldConfigurationScheme(
  client: JiraClient,
  id: string,
): Promise<void> {
  return del(client, `/rest/api/3/fieldconfigurationscheme/${id}`);
}

// --- Permission Schemes ---

/**
 * Returns all permission schemes.
 */
export async function getPermissionSchemes(client: JiraClient): Promise<any[]> {
  const result = await get(client, '/rest/api/3/permissionscheme');
  return result.permissionSchemes;
}

/**
 * Creates a permission scheme.
 */
export async function createPermissionScheme(client: JiraClient, scheme: unknown): Promise<any> {
  return post(client, '/rest/api/3/permissionscheme', scheme);
}

/**
 * Returns a permission scheme by ID.
 */
export async function getPermissionScheme(client: JiraClient, id: number): Promise<any> {
  return get(client, `/rest/api/3/permissionscheme/${id}`);
}

/**
 * Updates a permission scheme.
 */
export async function updatePermissionScheme(
  client: JiraClient,
  id: number,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/permissionscheme/${id}`, scheme);
}

/**
 * Deletes a permission scheme.
 */
export async function deletePermissionScheme(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/permissionscheme/${id}`);
}

// --- Notification Schemes ---

/**
 * Returns a paginated list of notification schemes.
 */
export async function getNotificationSchemes(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/notificationscheme', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Creates a notification scheme.
 */
export async function createNotificationScheme(
  client: JiraClient,
  scheme: unknown,
): Promise<any> {
  return post(client, '/rest/api/3/notificationscheme', scheme);
}

/**
 * Returns a notification scheme by ID.
 */
export async function getNotificationScheme(client: JiraClient, id: number): Promise<any> {
  return get(client, `/rest/api/3/notificationscheme/${id}`);
}

/**
 * Updates a notification scheme.
 */
export async function updateNotificationScheme(
  client: JiraClient,
  id: number,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/notificationscheme/${id}`, scheme);
}

/**
 * Deletes a notification scheme.
 */
export async function deleteNotificationScheme(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/notificationscheme/${id}`);
}

// --- Issue Security Schemes ---

/**
 * Returns all issue security schemes.
 */
export async function getIssueSecuritySchemes(client: JiraClient): Promise<any[]> {
  const result = await get(client, '/rest/api/3/issuesecurityschemes');
  return result.issueSecuritySchemes;
}

/**
 * Creates an issue security scheme.
 */
export async function createIssueSecurityScheme(
  client: JiraClient,
  scheme: unknown,
): Promise<any> {
  return post(client, '/rest/api/3/issuesecurityschemes', scheme);
}

/**
 * Returns an issue security scheme by ID.
 */
export async function getIssueSecurityScheme(client: JiraClient, id: number): Promise<any> {
  return get(client, `/rest/api/3/issuesecurityschemes/${id}`);
}

/**
 * Updates an issue security scheme.
 */
export async function updateIssueSecurityScheme(
  client: JiraClient,
  id: number,
  scheme: unknown,
): Promise<any> {
  return put(client, `/rest/api/3/issuesecurityschemes/${id}`, scheme);
}

/**
 * Deletes an issue security scheme.
 */
export async function deleteIssueSecurityScheme(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/issuesecurityschemes/${id}`);
}

// --- Roles ---

/**
 * Returns all project roles.
 */
export async function getAllRoles(client: JiraClient): Promise<any> {
  return get(client, '/rest/api/3/role');
}

/**
 * Creates a project role.
 */
export async function createRole(client: JiraClient, role: unknown): Promise<any> {
  return post(client, '/rest/api/3/role', role);
}

/**
 * Returns a project role by ID.
 */
export async function getRole(client: JiraClient, id: number): Promise<any> {
  return get(client, `/rest/api/3/role/${id}`);
}

/**
 * Updates a project role.
 */
export async function updateRole(client: JiraClient, id: number, role: unknown): Promise<any> {
  return put(client, `/rest/api/3/role/${id}`, role);
}

/**
 * Deletes a project role.
 */
export async function deleteRole(client: JiraClient, id: number): Promise<void> {
  return del(client, `/rest/api/3/role/${id}`);
}

// --- Webhooks ---

/**
 * Returns a paginated list of webhooks.
 */
export async function getWebhooks(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<any> {
  return get(client, '/rest/api/3/webhook', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Registers webhooks.
 */
export async function registerWebhooks(client: JiraClient, webhooks: unknown): Promise<any[]> {
  const result = await post(client, '/rest/api/3/webhook', webhooks);
  return result.webhookRegistrationResult;
}

/**
 * Deletes webhooks by IDs.
 */
export async function deleteWebhooks(client: JiraClient, webhookIds: number[]): Promise<any> {
  return deleteWithBody(client, '/rest/api/3/webhook', { webhookIds });
}

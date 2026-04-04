import { get, post, put, del, deleteWithBody } from './client.js';

// --- Dashboards ---

/**
 * Returns a paginated list of dashboards.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getDashboards(client, startAt, maxResults) {
  return get(client, '/rest/api/3/dashboard', { startAt, maxResults });
}

/**
 * Creates a new dashboard.
 * @param {object} client
 * @param {object} dashboard
 * @returns {Promise<any>}
 */
export async function createDashboard(client, dashboard) {
  return post(client, '/rest/api/3/dashboard', dashboard);
}

/**
 * Returns a dashboard by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<any>}
 */
export async function getDashboard(client, id) {
  return get(client, `/rest/api/3/dashboard/${id}`);
}

/**
 * Updates a dashboard.
 * @param {object} client
 * @param {string} id
 * @param {object} dashboard
 * @returns {Promise<any>}
 */
export async function updateDashboard(client, id, dashboard) {
  return put(client, `/rest/api/3/dashboard/${id}`, dashboard);
}

/**
 * Deletes a dashboard by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteDashboard(client, id) {
  return del(client, `/rest/api/3/dashboard/${id}`);
}

/**
 * Copies a dashboard.
 * @param {object} client
 * @param {string} id
 * @param {object} dashboard
 * @returns {Promise<any>}
 */
export async function copyDashboard(client, id, dashboard) {
  return post(client, `/rest/api/3/dashboard/${id}/copy`, dashboard);
}

/**
 * Searches for dashboards by name.
 * @param {object} client
 * @param {string} name
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function searchDashboards(client, name, startAt, maxResults) {
  const query = { startAt, maxResults };
  if (name) query.dashboardName = name;
  return get(client, '/rest/api/3/dashboard/search', query);
}

/**
 * Returns the gadgets on a dashboard.
 * @param {object} client
 * @param {string} dashboardId
 * @returns {Promise<any>}
 */
export async function getDashboardGadgets(client, dashboardId) {
  return get(client, `/rest/api/3/dashboard/${dashboardId}/gadget`);
}

/**
 * Adds a gadget to a dashboard.
 * @param {object} client
 * @param {string} dashboardId
 * @param {object} gadget
 * @returns {Promise<any>}
 */
export async function addDashboardGadget(client, dashboardId, gadget) {
  return post(client, `/rest/api/3/dashboard/${dashboardId}/gadget`, gadget);
}

/**
 * Updates a gadget on a dashboard.
 * @param {object} client
 * @param {string} dashboardId
 * @param {string} gadgetId
 * @param {object} gadget
 * @returns {Promise<null>}
 */
export async function updateDashboardGadget(client, dashboardId, gadgetId, gadget) {
  return put(client, `/rest/api/3/dashboard/${dashboardId}/gadget/${gadgetId}`, gadget);
}

/**
 * Removes a gadget from a dashboard.
 * @param {object} client
 * @param {string} dashboardId
 * @param {string} gadgetId
 * @returns {Promise<null>}
 */
export async function removeDashboardGadget(client, dashboardId, gadgetId) {
  return del(client, `/rest/api/3/dashboard/${dashboardId}/gadget/${gadgetId}`);
}

// --- Fields ---

/**
 * Returns all fields.
 * @param {object} client
 * @returns {Promise<any>}
 */
export async function getFields(client) {
  return get(client, '/rest/api/3/field');
}

/**
 * Creates a custom field.
 * @param {object} client
 * @param {object} field
 * @returns {Promise<any>}
 */
export async function createCustomField(client, field) {
  return post(client, '/rest/api/3/field', field);
}

/**
 * Updates a custom field.
 * @param {object} client
 * @param {string} fieldId
 * @param {object} field
 * @returns {Promise<null>}
 */
export async function updateCustomField(client, fieldId, field) {
  return put(client, `/rest/api/3/field/${fieldId}`, field);
}

/**
 * Deletes a custom field.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteCustomField(client, id) {
  return del(client, `/rest/api/3/field/${id}`);
}

/**
 * Searches for fields.
 * @param {object} client
 * @param {string} query
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function searchFields(client, query, startAt, maxResults) {
  const params = { startAt, maxResults };
  if (query) params.query = query;
  return get(client, '/rest/api/3/field/search', params);
}

/**
 * Returns contexts for a field.
 * @param {object} client
 * @param {string} fieldId
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getFieldContexts(client, fieldId, startAt, maxResults) {
  return get(client, `/rest/api/3/field/${fieldId}/context`, { startAt, maxResults });
}

/**
 * Moves a custom field to trash.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function trashCustomField(client, id) {
  return post(client, `/rest/api/3/field/${id}/trash`, null);
}

/**
 * Restores a custom field from trash.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function restoreCustomField(client, id) {
  return post(client, `/rest/api/3/field/${id}/restore`, null);
}

/**
 * Returns trashed fields.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getTrashedFields(client, startAt, maxResults) {
  return get(client, '/rest/api/3/field/search/trashed', { startAt, maxResults });
}

// --- Screens ---

/**
 * Returns a paginated list of screens.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getScreens(client, startAt, maxResults) {
  return get(client, '/rest/api/3/screens', { startAt, maxResults });
}

/**
 * Creates a screen.
 * @param {object} client
 * @param {object} screen
 * @returns {Promise<any>}
 */
export async function createScreen(client, screen) {
  return post(client, '/rest/api/3/screens', screen);
}

/**
 * Updates a screen.
 * @param {object} client
 * @param {number} screenId
 * @param {object} screen
 * @returns {Promise<any>}
 */
export async function updateScreen(client, screenId, screen) {
  return put(client, `/rest/api/3/screens/${screenId}`, screen);
}

/**
 * Deletes a screen.
 * @param {object} client
 * @param {number} screenId
 * @returns {Promise<null>}
 */
export async function deleteScreen(client, screenId) {
  return del(client, `/rest/api/3/screens/${screenId}`);
}

/**
 * Returns the tabs for a screen.
 * @param {object} client
 * @param {number} screenId
 * @returns {Promise<any>}
 */
export async function getScreenTabs(client, screenId) {
  return get(client, `/rest/api/3/screens/${screenId}/tabs`);
}

/**
 * Creates a tab on a screen.
 * @param {object} client
 * @param {number} screenId
 * @param {object} tab
 * @returns {Promise<any>}
 */
export async function createScreenTab(client, screenId, tab) {
  return post(client, `/rest/api/3/screens/${screenId}/tabs`, tab);
}

/**
 * Updates a tab on a screen.
 * @param {object} client
 * @param {number} screenId
 * @param {number} tabId
 * @param {object} tab
 * @returns {Promise<any>}
 */
export async function updateScreenTab(client, screenId, tabId, tab) {
  return put(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}`, tab);
}

/**
 * Deletes a tab from a screen.
 * @param {object} client
 * @param {number} screenId
 * @param {number} tabId
 * @returns {Promise<null>}
 */
export async function deleteScreenTab(client, screenId, tabId) {
  return del(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}`);
}

/**
 * Returns the fields on a screen tab.
 * @param {object} client
 * @param {number} screenId
 * @param {number} tabId
 * @returns {Promise<any>}
 */
export async function getScreenTabFields(client, screenId, tabId) {
  return get(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields`);
}

/**
 * Adds a field to a screen tab.
 * @param {object} client
 * @param {number} screenId
 * @param {number} tabId
 * @param {string} fieldId
 * @returns {Promise<any>}
 */
export async function addScreenTabField(client, screenId, tabId, fieldId) {
  return post(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields`, { fieldId });
}

/**
 * Removes a field from a screen tab.
 * @param {object} client
 * @param {number} screenId
 * @param {number} tabId
 * @param {string} fieldId
 * @returns {Promise<null>}
 */
export async function removeScreenTabField(client, screenId, tabId, fieldId) {
  return del(client, `/rest/api/3/screens/${screenId}/tabs/${tabId}/fields/${fieldId}`);
}

// --- Screen Schemes ---

/**
 * Returns a paginated list of screen schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getScreenSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/screenscheme', { startAt, maxResults });
}

/**
 * Creates a screen scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createScreenScheme(client, scheme) {
  return post(client, '/rest/api/3/screenscheme', scheme);
}

/**
 * Updates a screen scheme.
 * @param {object} client
 * @param {number} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateScreenScheme(client, id, scheme) {
  return put(client, `/rest/api/3/screenscheme/${id}`, scheme);
}

/**
 * Deletes a screen scheme.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteScreenScheme(client, id) {
  return del(client, `/rest/api/3/screenscheme/${id}`);
}

// --- Workflows ---

/**
 * Returns all workflows.
 * @param {object} client
 * @returns {Promise<any>}
 */
export async function getWorkflows(client) {
  return get(client, '/rest/api/3/workflow');
}

/**
 * Searches for workflows.
 * @param {object} client
 * @param {string} query
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function searchWorkflows(client, query, startAt, maxResults) {
  const params = { startAt, maxResults };
  if (query) params.queryString = query;
  return get(client, '/rest/api/3/workflow/search', params);
}

// --- Workflow Schemes ---

/**
 * Returns a paginated list of workflow schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getWorkflowSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/workflowscheme', { startAt, maxResults });
}

/**
 * Creates a workflow scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createWorkflowScheme(client, scheme) {
  return post(client, '/rest/api/3/workflowscheme', scheme);
}

/**
 * Returns a workflow scheme by ID.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<any>}
 */
export async function getWorkflowScheme(client, id) {
  return get(client, `/rest/api/3/workflowscheme/${id}`);
}

/**
 * Updates a workflow scheme.
 * @param {object} client
 * @param {number} id
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function updateWorkflowScheme(client, id, scheme) {
  return put(client, `/rest/api/3/workflowscheme/${id}`, scheme);
}

/**
 * Deletes a workflow scheme.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteWorkflowScheme(client, id) {
  return del(client, `/rest/api/3/workflowscheme/${id}`);
}

// --- Issue Type Schemes ---

/**
 * Returns a paginated list of issue type schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getIssueTypeSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/issuetypescheme', { startAt, maxResults });
}

/**
 * Creates an issue type scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createIssueTypeScheme(client, scheme) {
  return post(client, '/rest/api/3/issuetypescheme', scheme);
}

/**
 * Updates an issue type scheme.
 * @param {object} client
 * @param {string} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateIssueTypeScheme(client, id, scheme) {
  return put(client, `/rest/api/3/issuetypescheme/${id}`, scheme);
}

/**
 * Deletes an issue type scheme.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteIssueTypeScheme(client, id) {
  return del(client, `/rest/api/3/issuetypescheme/${id}`);
}

// --- Issue Type Screen Schemes ---

/**
 * Returns a paginated list of issue type screen schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getIssueTypeScreenSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/issuetypescreenscheme', { startAt, maxResults });
}

/**
 * Creates an issue type screen scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createIssueTypeScreenScheme(client, scheme) {
  return post(client, '/rest/api/3/issuetypescreenscheme', scheme);
}

/**
 * Updates an issue type screen scheme.
 * @param {object} client
 * @param {string} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateIssueTypeScreenScheme(client, id, scheme) {
  return put(client, `/rest/api/3/issuetypescreenscheme/${id}`, scheme);
}

/**
 * Deletes an issue type screen scheme.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteIssueTypeScreenScheme(client, id) {
  return del(client, `/rest/api/3/issuetypescreenscheme/${id}`);
}

// --- Field Configurations ---

/**
 * Returns a paginated list of field configurations.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getFieldConfigurations(client, startAt, maxResults) {
  return get(client, '/rest/api/3/fieldconfiguration', { startAt, maxResults });
}

/**
 * Creates a field configuration.
 * @param {object} client
 * @param {object} config
 * @returns {Promise<any>}
 */
export async function createFieldConfiguration(client, config) {
  return post(client, '/rest/api/3/fieldconfiguration', config);
}

/**
 * Updates a field configuration.
 * @param {object} client
 * @param {number} id
 * @param {object} config
 * @returns {Promise<null>}
 */
export async function updateFieldConfiguration(client, id, config) {
  return put(client, `/rest/api/3/fieldconfiguration/${id}`, config);
}

/**
 * Deletes a field configuration.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteFieldConfiguration(client, id) {
  return del(client, `/rest/api/3/fieldconfiguration/${id}`);
}

// --- Field Configuration Schemes ---

/**
 * Returns a paginated list of field configuration schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getFieldConfigurationSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/fieldconfigurationscheme', { startAt, maxResults });
}

/**
 * Creates a field configuration scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createFieldConfigurationScheme(client, scheme) {
  return post(client, '/rest/api/3/fieldconfigurationscheme', scheme);
}

/**
 * Updates a field configuration scheme.
 * @param {object} client
 * @param {string} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateFieldConfigurationScheme(client, id, scheme) {
  return put(client, `/rest/api/3/fieldconfigurationscheme/${id}`, scheme);
}

/**
 * Deletes a field configuration scheme.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteFieldConfigurationScheme(client, id) {
  return del(client, `/rest/api/3/fieldconfigurationscheme/${id}`);
}

// --- Permission Schemes ---

/**
 * Returns all permission schemes.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getPermissionSchemes(client) {
  const result = await get(client, '/rest/api/3/permissionscheme');
  return result.permissionSchemes;
}

/**
 * Creates a permission scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createPermissionScheme(client, scheme) {
  return post(client, '/rest/api/3/permissionscheme', scheme);
}

/**
 * Returns a permission scheme by ID.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<any>}
 */
export async function getPermissionScheme(client, id) {
  return get(client, `/rest/api/3/permissionscheme/${id}`);
}

/**
 * Updates a permission scheme.
 * @param {object} client
 * @param {number} id
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function updatePermissionScheme(client, id, scheme) {
  return put(client, `/rest/api/3/permissionscheme/${id}`, scheme);
}

/**
 * Deletes a permission scheme.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deletePermissionScheme(client, id) {
  return del(client, `/rest/api/3/permissionscheme/${id}`);
}

// --- Notification Schemes ---

/**
 * Returns a paginated list of notification schemes.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getNotificationSchemes(client, startAt, maxResults) {
  return get(client, '/rest/api/3/notificationscheme', { startAt, maxResults });
}

/**
 * Creates a notification scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createNotificationScheme(client, scheme) {
  return post(client, '/rest/api/3/notificationscheme', scheme);
}

/**
 * Returns a notification scheme by ID.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<any>}
 */
export async function getNotificationScheme(client, id) {
  return get(client, `/rest/api/3/notificationscheme/${id}`);
}

/**
 * Updates a notification scheme.
 * @param {object} client
 * @param {number} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateNotificationScheme(client, id, scheme) {
  return put(client, `/rest/api/3/notificationscheme/${id}`, scheme);
}

/**
 * Deletes a notification scheme.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteNotificationScheme(client, id) {
  return del(client, `/rest/api/3/notificationscheme/${id}`);
}

// --- Issue Security Schemes ---

/**
 * Returns all issue security schemes.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getIssueSecuritySchemes(client) {
  const result = await get(client, '/rest/api/3/issuesecurityschemes');
  return result.issueSecuritySchemes;
}

/**
 * Creates an issue security scheme.
 * @param {object} client
 * @param {object} scheme
 * @returns {Promise<any>}
 */
export async function createIssueSecurityScheme(client, scheme) {
  return post(client, '/rest/api/3/issuesecurityschemes', scheme);
}

/**
 * Returns an issue security scheme by ID.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<any>}
 */
export async function getIssueSecurityScheme(client, id) {
  return get(client, `/rest/api/3/issuesecurityschemes/${id}`);
}

/**
 * Updates an issue security scheme.
 * @param {object} client
 * @param {number} id
 * @param {object} scheme
 * @returns {Promise<null>}
 */
export async function updateIssueSecurityScheme(client, id, scheme) {
  return put(client, `/rest/api/3/issuesecurityschemes/${id}`, scheme);
}

/**
 * Deletes an issue security scheme.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteIssueSecurityScheme(client, id) {
  return del(client, `/rest/api/3/issuesecurityschemes/${id}`);
}

// --- Roles ---

/**
 * Returns all project roles.
 * @param {object} client
 * @returns {Promise<any>}
 */
export async function getAllRoles(client) {
  return get(client, '/rest/api/3/role');
}

/**
 * Creates a project role.
 * @param {object} client
 * @param {object} role
 * @returns {Promise<any>}
 */
export async function createRole(client, role) {
  return post(client, '/rest/api/3/role', role);
}

/**
 * Returns a project role by ID.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<any>}
 */
export async function getRole(client, id) {
  return get(client, `/rest/api/3/role/${id}`);
}

/**
 * Updates a project role.
 * @param {object} client
 * @param {number} id
 * @param {object} role
 * @returns {Promise<any>}
 */
export async function updateRole(client, id, role) {
  return put(client, `/rest/api/3/role/${id}`, role);
}

/**
 * Deletes a project role.
 * @param {object} client
 * @param {number} id
 * @returns {Promise<null>}
 */
export async function deleteRole(client, id) {
  return del(client, `/rest/api/3/role/${id}`);
}

// --- Webhooks ---

/**
 * Returns a paginated list of webhooks.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<any>}
 */
export async function getWebhooks(client, startAt, maxResults) {
  return get(client, '/rest/api/3/webhook', { startAt, maxResults });
}

/**
 * Registers webhooks.
 * @param {object} client
 * @param {object} webhooks
 * @returns {Promise<any[]>}
 */
export async function registerWebhooks(client, webhooks) {
  const result = await post(client, '/rest/api/3/webhook', webhooks);
  return result.webhookRegistrationResult;
}

/**
 * Deletes webhooks by IDs.
 * @param {object} client
 * @param {number[]} webhookIds
 * @returns {Promise<null>}
 */
export async function deleteWebhooks(client, webhookIds) {
  return deleteWithBody(client, '/rest/api/3/webhook', { webhookIds });
}

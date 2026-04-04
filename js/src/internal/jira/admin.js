import { get, post, put, del } from './client.js';

// ============================================================================
// Components
// ============================================================================

/**
 * Creates a new project component.
 * @param {object} client
 * @param {Object} component
 * @returns {Promise<Object>}
 */
export async function createComponent(client, component) {
  return post(client, '/rest/api/3/component', component);
}

/**
 * Returns a project component by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getComponent(client, id) {
  return get(client, `/rest/api/3/component/${id}`);
}

/**
 * Updates a project component.
 * @param {object} client
 * @param {string} id
 * @param {Object} component
 * @returns {Promise<Object>}
 */
export async function updateComponent(client, id, component) {
  return put(client, `/rest/api/3/component/${id}`, component);
}

/**
 * Deletes a project component.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteComponent(client, id) {
  return del(client, `/rest/api/3/component/${id}`);
}

/**
 * Returns the issue count for a component.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getComponentIssueCount(client, id) {
  return get(client, `/rest/api/3/component/${id}/relatedIssueCounts`);
}

// ============================================================================
// Versions
// ============================================================================

/**
 * Creates a new project version.
 * @param {object} client
 * @param {Object} version
 * @returns {Promise<Object>}
 */
export async function createVersion(client, version) {
  return post(client, '/rest/api/3/version', version);
}

/**
 * Returns a project version by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getVersion(client, id) {
  return get(client, `/rest/api/3/version/${id}`);
}

/**
 * Updates a project version.
 * @param {object} client
 * @param {string} id
 * @param {Object} version
 * @returns {Promise<Object>}
 */
export async function updateVersion(client, id, version) {
  return put(client, `/rest/api/3/version/${id}`, version);
}

/**
 * Deletes a project version.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteVersion(client, id) {
  return del(client, `/rest/api/3/version/${id}`);
}

/**
 * Merges a version into another version.
 * @param {object} client
 * @param {string} id
 * @param {string} moveIssuesTo
 * @returns {Promise<any>}
 */
export async function mergeVersions(client, id, moveIssuesTo) {
  return put(client, `/rest/api/3/version/${id}/mergeto/${moveIssuesTo}`, null);
}

/**
 * Moves a version to a new position.
 * @param {object} client
 * @param {string} id
 * @param {Object} position
 * @returns {Promise<Object>}
 */
export async function moveVersion(client, id, position) {
  return post(client, `/rest/api/3/version/${id}/move`, position);
}

/**
 * Returns issue counts related to a version.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getVersionRelatedIssueCounts(client, id) {
  return get(client, `/rest/api/3/version/${id}/relatedIssueCounts`);
}

/**
 * Returns the unresolved issue count for a version.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getVersionUnresolvedIssueCount(client, id) {
  return get(client, `/rest/api/3/version/${id}/unresolvedIssueCount`);
}

// ============================================================================
// Users
// ============================================================================

/**
 * Returns a user by account ID.
 * @param {object} client
 * @param {string} accountId
 * @returns {Promise<Object>}
 */
export async function getUser(client, accountId) {
  return get(client, '/rest/api/3/user', { accountId });
}

/**
 * Creates a new user.
 * @param {object} client
 * @param {Object} user
 * @returns {Promise<Object>}
 */
export async function createUser(client, user) {
  return post(client, '/rest/api/3/user', user);
}

/**
 * Deletes a user by account ID.
 * @param {object} client
 * @param {string} accountId
 * @returns {Promise<null>}
 */
export async function deleteUser(client, accountId) {
  return del(client, '/rest/api/3/user', { accountId });
}

/**
 * Returns multiple users by account IDs.
 * @param {object} client
 * @param {string[]} accountIds
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function getUsersBulk(client, accountIds, startAt, maxResults) {
  const params = new URLSearchParams();
  for (const id of accountIds) {
    params.append('accountId', id);
  }
  params.set('startAt', String(startAt));
  params.set('maxResults', String(maxResults));
  return get(client, `/rest/api/3/user/bulk?${params.toString()}`);
}

/**
 * Searches for users by query string.
 * @param {object} client
 * @param {string} query
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object[]>}
 */
export async function findUsers(client, query, startAt, maxResults) {
  return get(client, '/rest/api/3/user/search', { query, startAt: String(startAt), maxResults: String(maxResults) });
}

/**
 * Searches for users assignable to issues.
 * @param {object} client
 * @param {string} query
 * @param {string} project
 * @param {string} issueKey
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object[]>}
 */
export async function findUsersAssignable(client, query, project, issueKey, startAt, maxResults) {
  const params = {};
  if (query) params.query = query;
  if (project) params.project = project;
  if (issueKey) params.issueKey = issueKey;
  params.startAt = String(startAt);
  params.maxResults = String(maxResults);
  return get(client, '/rest/api/3/user/assignable/search', params);
}

/**
 * Returns the currently authenticated user.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getCurrentUser(client) {
  return get(client, '/rest/api/3/myself');
}

/**
 * Returns all users with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object[]>}
 */
export async function getAllUsers(client, startAt, maxResults) {
  return get(client, '/rest/api/3/users/search', { startAt: String(startAt), maxResults: String(maxResults) });
}

// ============================================================================
// Groups
// ============================================================================

/**
 * Returns a group by name.
 * @param {object} client
 * @param {string} groupName
 * @returns {Promise<Object>}
 */
export async function getGroup(client, groupName) {
  return get(client, '/rest/api/3/group', { groupname: groupName });
}

/**
 * Creates a new group.
 * @param {object} client
 * @param {string} name
 * @returns {Promise<Object>}
 */
export async function createGroup(client, name) {
  return post(client, '/rest/api/3/group', { name });
}

/**
 * Deletes a group by name.
 * @param {object} client
 * @param {string} groupName
 * @returns {Promise<null>}
 */
export async function deleteGroup(client, groupName) {
  return del(client, '/rest/api/3/group', { groupname: groupName });
}

/**
 * Returns the members of a group.
 * @param {object} client
 * @param {string} groupName
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function getGroupMembers(client, groupName, startAt, maxResults) {
  return get(client, '/rest/api/3/group/member', {
    groupname: groupName,
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Adds a user to a group.
 * @param {object} client
 * @param {string} groupName
 * @param {string} accountId
 * @returns {Promise<Object>}
 */
export async function addUserToGroup(client, groupName, accountId) {
  const params = new URLSearchParams({ groupname: groupName });
  return post(client, `/rest/api/3/group/user?${params.toString()}`, { accountId });
}

/**
 * Removes a user from a group.
 * @param {object} client
 * @param {string} groupName
 * @param {string} accountId
 * @returns {Promise<null>}
 */
export async function removeUserFromGroup(client, groupName, accountId) {
  return del(client, '/rest/api/3/group/user', { groupname: groupName, accountId });
}

/**
 * Returns groups with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function getBulkGroups(client, startAt, maxResults) {
  return get(client, '/rest/api/3/group/bulk', { startAt: String(startAt), maxResults: String(maxResults) });
}

/**
 * Searches for groups using the picker endpoint.
 * @param {object} client
 * @param {string} query
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function findGroups(client, query, maxResults) {
  return get(client, '/rest/api/3/groups/picker', { query, maxResults: String(maxResults) });
}

// ============================================================================
// Issue Links
// ============================================================================

/**
 * Creates a link between two issues.
 * @param {object} client
 * @param {Object} link
 * @returns {Promise<null>}
 */
export async function createIssueLink(client, link) {
  return post(client, '/rest/api/3/issueLink', link);
}

/**
 * Returns an issue link by ID.
 * @param {object} client
 * @param {string} linkId
 * @returns {Promise<Object>}
 */
export async function getIssueLink(client, linkId) {
  return get(client, `/rest/api/3/issueLink/${linkId}`);
}

/**
 * Deletes an issue link.
 * @param {object} client
 * @param {string} linkId
 * @returns {Promise<null>}
 */
export async function deleteIssueLink(client, linkId) {
  return del(client, `/rest/api/3/issueLink/${linkId}`);
}

/**
 * Returns all issue link types.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getIssueLinkTypes(client) {
  const result = await get(client, '/rest/api/3/issueLinkType');
  return result.issueLinkTypes;
}

/**
 * Creates a new issue link type.
 * @param {object} client
 * @param {Object} linkType
 * @returns {Promise<Object>}
 */
export async function createIssueLinkType(client, linkType) {
  return post(client, '/rest/api/3/issueLinkType', linkType);
}

/**
 * Returns an issue link type by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getIssueLinkType(client, id) {
  return get(client, `/rest/api/3/issueLinkType/${id}`);
}

/**
 * Updates an issue link type.
 * @param {object} client
 * @param {string} id
 * @param {Object} linkType
 * @returns {Promise<Object>}
 */
export async function updateIssueLinkType(client, id, linkType) {
  return put(client, `/rest/api/3/issueLinkType/${id}`, linkType);
}

/**
 * Deletes an issue link type.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteIssueLinkType(client, id) {
  return del(client, `/rest/api/3/issueLinkType/${id}`);
}

// ============================================================================
// Attachments
// ============================================================================

/**
 * Returns an attachment by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getAttachment(client, id) {
  return get(client, `/rest/api/3/attachment/${id}`);
}

/**
 * Deletes an attachment.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteAttachment(client, id) {
  return del(client, `/rest/api/3/attachment/${id}`);
}

/**
 * Returns attachment settings.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getAttachmentMeta(client) {
  return get(client, '/rest/api/3/attachment/meta');
}

// ============================================================================
// Issue Types
// ============================================================================

/**
 * Returns all issue types.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getAllIssueTypes(client) {
  return get(client, '/rest/api/3/issuetype');
}

/**
 * Creates a new issue type.
 * @param {object} client
 * @param {Object} issueType
 * @returns {Promise<Object>}
 */
export async function createIssueType(client, issueType) {
  return post(client, '/rest/api/3/issuetype', issueType);
}

/**
 * Returns an issue type by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getIssueType(client, id) {
  return get(client, `/rest/api/3/issuetype/${id}`);
}

/**
 * Updates an issue type.
 * @param {object} client
 * @param {string} id
 * @param {Object} issueType
 * @returns {Promise<Object>}
 */
export async function updateIssueType(client, id, issueType) {
  return put(client, `/rest/api/3/issuetype/${id}`, issueType);
}

/**
 * Deletes an issue type.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteIssueType(client, id) {
  return del(client, `/rest/api/3/issuetype/${id}`);
}

/**
 * Returns alternative issue types for the given issue type.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object[]>}
 */
export async function getIssueTypeAlternatives(client, id) {
  return get(client, `/rest/api/3/issuetype/${id}/alternatives`);
}

/**
 * Returns issue types for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<Object[]>}
 */
export async function getProjectIssueTypes(client, projectIdOrKey) {
  return get(client, '/rest/api/3/issuetype/project', { projectId: projectIdOrKey });
}

// ============================================================================
// Priorities
// ============================================================================

/**
 * Returns all priorities.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getAllPriorities(client) {
  return get(client, '/rest/api/3/priority');
}

/**
 * Creates a new priority.
 * @param {object} client
 * @param {Object} priority
 * @returns {Promise<Object>}
 */
export async function createPriority(client, priority) {
  return post(client, '/rest/api/3/priority', priority);
}

/**
 * Returns a priority by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getPriority(client, id) {
  return get(client, `/rest/api/3/priority/${id}`);
}

/**
 * Updates a priority.
 * @param {object} client
 * @param {string} id
 * @param {Object} priority
 * @returns {Promise<Object>}
 */
export async function updatePriority(client, id, priority) {
  return put(client, `/rest/api/3/priority/${id}`, priority);
}

/**
 * Deletes a priority.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deletePriority(client, id) {
  return del(client, `/rest/api/3/priority/${id}`);
}

/**
 * Searches for priorities with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function searchPriorities(client, startAt, maxResults) {
  return get(client, '/rest/api/3/priority/search', { startAt: String(startAt), maxResults: String(maxResults) });
}

// ============================================================================
// Resolutions
// ============================================================================

/**
 * Returns all resolutions.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getAllResolutions(client) {
  return get(client, '/rest/api/3/resolution');
}

/**
 * Creates a new resolution.
 * @param {object} client
 * @param {Object} resolution
 * @returns {Promise<Object>}
 */
export async function createResolution(client, resolution) {
  return post(client, '/rest/api/3/resolution', resolution);
}

/**
 * Returns a resolution by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<Object>}
 */
export async function getResolution(client, id) {
  return get(client, `/rest/api/3/resolution/${id}`);
}

/**
 * Updates a resolution.
 * @param {object} client
 * @param {string} id
 * @param {Object} resolution
 * @returns {Promise<Object>}
 */
export async function updateResolution(client, id, resolution) {
  return put(client, `/rest/api/3/resolution/${id}`, resolution);
}

/**
 * Deletes a resolution.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteResolution(client, id) {
  return del(client, `/rest/api/3/resolution/${id}`);
}

/**
 * Searches for resolutions with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function searchResolutions(client, startAt, maxResults) {
  return get(client, '/rest/api/3/resolution/search', { startAt: String(startAt), maxResults: String(maxResults) });
}

// ============================================================================
// Statuses
// ============================================================================

/**
 * Returns all statuses.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getAllStatuses(client) {
  return get(client, '/rest/api/3/status');
}

/**
 * Returns a status by ID or name.
 * @param {object} client
 * @param {string} idOrName
 * @returns {Promise<Object>}
 */
export async function getStatus(client, idOrName) {
  return get(client, `/rest/api/3/status/${idOrName}`);
}

/**
 * Searches for statuses with pagination.
 * @param {object} client
 * @param {string} searchString
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function searchStatuses(client, searchString, startAt, maxResults) {
  const params = {};
  if (searchString) params.searchString = searchString;
  params.startAt = String(startAt);
  params.maxResults = String(maxResults);
  return get(client, '/rest/api/3/statuses/search', params);
}

/**
 * Returns all status categories.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getStatusCategories(client) {
  return get(client, '/rest/api/3/statuscategory');
}

/**
 * Returns a status category by ID or key.
 * @param {object} client
 * @param {string} idOrKey
 * @returns {Promise<Object>}
 */
export async function getStatusCategory(client, idOrKey) {
  return get(client, `/rest/api/3/statuscategory/${idOrKey}`);
}

// ============================================================================
// Labels
// ============================================================================

/**
 * Returns all labels with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function getLabels(client, startAt, maxResults) {
  return get(client, '/rest/api/3/label', { startAt: String(startAt), maxResults: String(maxResults) });
}

// ============================================================================
// Server & Config
// ============================================================================

/**
 * Returns Jira server information.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getServerInfo(client) {
  return get(client, '/rest/api/3/serverInfo');
}

/**
 * Returns the Jira configuration.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getConfiguration(client) {
  return get(client, '/rest/api/3/configuration');
}

/**
 * Returns the announcement banner settings.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getAnnouncementBanner(client) {
  return get(client, '/rest/api/3/announcementBanner');
}

/**
 * Updates the announcement banner settings.
 * @param {object} client
 * @param {Object} banner
 * @returns {Promise<any>}
 */
export async function setAnnouncementBanner(client, banner) {
  return put(client, '/rest/api/3/announcementBanner', banner);
}

/**
 * Returns audit records with pagination.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<Object>}
 */
export async function getAuditRecords(client, startAt, maxResults) {
  return get(client, '/rest/api/3/auditing/record', { offset: String(startAt), limit: String(maxResults) });
}

/**
 * Returns all application roles.
 * @param {object} client
 * @returns {Promise<Object[]>}
 */
export async function getApplicationRoles(client) {
  return get(client, '/rest/api/3/applicationrole');
}

/**
 * Returns an application role by key.
 * @param {object} client
 * @param {string} key
 * @returns {Promise<Object>}
 */
export async function getApplicationRole(client, key) {
  return get(client, `/rest/api/3/applicationrole/${key}`);
}

// ============================================================================
// Permissions
// ============================================================================

/**
 * Returns the permissions for the current user.
 * @param {object} client
 * @param {string} projectKey
 * @param {string} issueKey
 * @returns {Promise<Object>}
 */
export async function getMyPermissions(client, projectKey, issueKey) {
  const params = {};
  if (projectKey) params.projectKey = projectKey;
  if (issueKey) params.issueKey = issueKey;
  const result = await get(client, '/rest/api/3/mypermissions', params);
  return result.permissions;
}

/**
 * Returns all permissions in the system.
 * @param {object} client
 * @returns {Promise<Object>}
 */
export async function getAllPermissions(client) {
  const result = await get(client, '/rest/api/3/permissions');
  return result.permissions;
}

// ============================================================================
// Tasks
// ============================================================================

/**
 * Returns an async task result by ID.
 * @param {object} client
 * @param {string} taskId
 * @returns {Promise<Object>}
 */
export async function getTask(client, taskId) {
  return get(client, `/rest/api/3/task/${taskId}`);
}

/**
 * Cancels an async task.
 * @param {object} client
 * @param {string} taskId
 * @returns {Promise<null>}
 */
export async function cancelTask(client, taskId) {
  return post(client, `/rest/api/3/task/${taskId}/cancel`, null);
}

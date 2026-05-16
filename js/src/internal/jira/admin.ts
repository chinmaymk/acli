import { get, post, put, del, getBuffer } from './client.js';
import type { JiraClient } from './client.js';
import type {
  AnnouncementBanner,
  ApplicationRole,
  Attachment,
  AttachmentMeta,
  AuditRecords,
  ComponentIssueCount,
  Configuration,
  FoundGroups,
  Group,
  GroupMembers,
  IssueLink,
  IssueLinkType,
  IssueType,
  PageBean,
  Priority,
  ProjectComponent,
  Resolution,
  ServerInfo,
  StatusCategory,
  StatusDetails,
  TaskResult,
  UserDetails,
  UserPermission,
  Version,
  VersionIssueCounts,
  VersionUnresolvedIssueCount,
} from './types.js';

// ============================================================================
// Components
// ============================================================================

/**
 * Creates a new project component.
 */
export async function createComponent(client: JiraClient, component: Partial<ProjectComponent>): Promise<ProjectComponent> {
  return post(client, '/rest/api/3/component', component);
}

/**
 * Returns a project component by ID.
 */
export async function getComponent(client: JiraClient, id: string): Promise<ProjectComponent> {
  return get(client, `/rest/api/3/component/${id}`);
}

/**
 * Updates a project component.
 */
export async function updateComponent(
  client: JiraClient,
  id: string,
  component: Partial<ProjectComponent>,
): Promise<ProjectComponent> {
  return put(client, `/rest/api/3/component/${id}`, component);
}

/**
 * Deletes a project component.
 */
export async function deleteComponent(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/component/${id}`);
}

/**
 * Returns the issue count for a component.
 */
export async function getComponentIssueCount(client: JiraClient, id: string): Promise<ComponentIssueCount> {
  return get(client, `/rest/api/3/component/${id}/relatedIssueCounts`);
}

// ============================================================================
// Versions
// ============================================================================

/**
 * Creates a new project version.
 */
export async function createVersion(client: JiraClient, version: Partial<Version>): Promise<Version> {
  return post(client, '/rest/api/3/version', version);
}

/**
 * Returns a project version by ID.
 */
export async function getVersion(client: JiraClient, id: string): Promise<Version> {
  return get(client, `/rest/api/3/version/${id}`);
}

/**
 * Updates a project version.
 */
export async function updateVersion(
  client: JiraClient,
  id: string,
  version: Partial<Version>,
): Promise<Version> {
  return put(client, `/rest/api/3/version/${id}`, version);
}

/**
 * Deletes a project version.
 */
export async function deleteVersion(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/version/${id}`);
}

/**
 * Merges a version into another version.
 */
export async function mergeVersions(
  client: JiraClient,
  id: string,
  moveIssuesTo: string,
): Promise<void> {
  return put(client, `/rest/api/3/version/${id}/mergeto/${moveIssuesTo}`, null);
}

/**
 * Moves a version to a new position.
 */
export async function moveVersion(
  client: JiraClient,
  id: string,
  position: Record<string, string>,
): Promise<Version> {
  return post(client, `/rest/api/3/version/${id}/move`, position);
}

/**
 * Returns issue counts related to a version.
 */
export async function getVersionRelatedIssueCounts(client: JiraClient, id: string): Promise<VersionIssueCounts> {
  return get(client, `/rest/api/3/version/${id}/relatedIssueCounts`);
}

/**
 * Returns the unresolved issue count for a version.
 */
export async function getVersionUnresolvedIssueCount(
  client: JiraClient,
  id: string,
): Promise<VersionUnresolvedIssueCount> {
  return get(client, `/rest/api/3/version/${id}/unresolvedIssueCount`);
}

// ============================================================================
// Users
// ============================================================================

/**
 * Returns a user by account ID.
 */
export async function getUser(client: JiraClient, accountId: string): Promise<UserDetails> {
  return get(client, '/rest/api/3/user', { accountId });
}

/**
 * Creates a new user.
 */
export async function createUser(client: JiraClient, user: Partial<UserDetails>): Promise<UserDetails> {
  return post(client, '/rest/api/3/user', user);
}

/**
 * Deletes a user by account ID.
 */
export async function deleteUser(client: JiraClient, accountId: string): Promise<void> {
  return del(client, '/rest/api/3/user', { accountId });
}

/**
 * Returns multiple users by account IDs.
 */
export async function getUsersBulk(
  client: JiraClient,
  accountIds: string[],
  startAt: number,
  maxResults: number,
): Promise<PageBean<UserDetails>> {
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
 */
export async function findUsers(
  client: JiraClient,
  query: string,
  startAt: number,
  maxResults: number,
): Promise<UserDetails[]> {
  return get(client, '/rest/api/3/user/search', {
    query,
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Searches for users assignable to issues.
 */
export async function findUsersAssignable(
  client: JiraClient,
  query: string,
  project: string,
  issueKey: string,
  startAt: number,
  maxResults: number,
): Promise<UserDetails[]> {
  const params: Record<string, string> = {};
  if (query) params.query = query;
  if (project) params.project = project;
  if (issueKey) params.issueKey = issueKey;
  params.startAt = String(startAt);
  params.maxResults = String(maxResults);
  return get(client, '/rest/api/3/user/assignable/search', params);
}

/**
 * Returns the currently authenticated user.
 */
export async function getCurrentUser(client: JiraClient): Promise<UserDetails> {
  return get(client, '/rest/api/3/myself');
}

/**
 * Returns all users with pagination.
 */
export async function getAllUsers(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<UserDetails[]> {
  return get(client, '/rest/api/3/users/search', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

// ============================================================================
// Groups
// ============================================================================

/**
 * Returns a group by name.
 */
export async function getGroup(client: JiraClient, groupName: string): Promise<Group> {
  return get(client, '/rest/api/3/group', { groupname: groupName });
}

/**
 * Creates a new group.
 */
export async function createGroup(client: JiraClient, name: string): Promise<Group> {
  return post(client, '/rest/api/3/group', { name });
}

/**
 * Deletes a group by name.
 */
export async function deleteGroup(client: JiraClient, groupName: string): Promise<void> {
  return del(client, '/rest/api/3/group', { groupname: groupName });
}

/**
 * Returns the members of a group.
 */
export async function getGroupMembers(
  client: JiraClient,
  groupName: string,
  startAt: number,
  maxResults: number,
): Promise<GroupMembers> {
  return get(client, '/rest/api/3/group/member', {
    groupname: groupName,
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Adds a user to a group.
 */
export async function addUserToGroup(
  client: JiraClient,
  groupName: string,
  accountId: string,
): Promise<Group> {
  const params = new URLSearchParams({ groupname: groupName });
  return post(client, `/rest/api/3/group/user?${params.toString()}`, { accountId });
}

/**
 * Removes a user from a group.
 */
export async function removeUserFromGroup(
  client: JiraClient,
  groupName: string,
  accountId: string,
): Promise<void> {
  return del(client, '/rest/api/3/group/user', { groupname: groupName, accountId });
}

/**
 * Returns groups with pagination.
 */
export async function getBulkGroups(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<PageBean<Group>> {
  return get(client, '/rest/api/3/group/bulk', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

/**
 * Searches for groups using the picker endpoint.
 */
export async function findGroups(
  client: JiraClient,
  query: string,
  maxResults: number,
): Promise<FoundGroups> {
  return get(client, '/rest/api/3/groups/picker', { query, maxResults: String(maxResults) });
}

// ============================================================================
// Issue Links
// ============================================================================

/**
 * Creates a link between two issues.
 */
export async function createIssueLink(client: JiraClient, link: IssueLink): Promise<void> {
  return post(client, '/rest/api/3/issueLink', link);
}

/**
 * Returns an issue link by ID.
 */
export async function getIssueLink(client: JiraClient, linkId: string): Promise<IssueLink> {
  return get(client, `/rest/api/3/issueLink/${linkId}`);
}

/**
 * Deletes an issue link.
 */
export async function deleteIssueLink(client: JiraClient, linkId: string): Promise<void> {
  return del(client, `/rest/api/3/issueLink/${linkId}`);
}

/**
 * Returns all issue link types.
 */
export async function getIssueLinkTypes(client: JiraClient): Promise<IssueLinkType[]> {
  const result = await get<{ issueLinkTypes: IssueLinkType[] }>(client, '/rest/api/3/issueLinkType');
  return result.issueLinkTypes;
}

/**
 * Creates a new issue link type.
 */
export async function createIssueLinkType(client: JiraClient, linkType: Partial<IssueLinkType>): Promise<IssueLinkType> {
  return post(client, '/rest/api/3/issueLinkType', linkType);
}

/**
 * Returns an issue link type by ID.
 */
export async function getIssueLinkType(client: JiraClient, id: string): Promise<IssueLinkType> {
  return get(client, `/rest/api/3/issueLinkType/${id}`);
}

/**
 * Updates an issue link type.
 */
export async function updateIssueLinkType(
  client: JiraClient,
  id: string,
  linkType: Partial<IssueLinkType>,
): Promise<IssueLinkType> {
  return put(client, `/rest/api/3/issueLinkType/${id}`, linkType);
}

/**
 * Deletes an issue link type.
 */
export async function deleteIssueLinkType(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/issueLinkType/${id}`);
}

// ============================================================================
// Attachments
// ============================================================================

/**
 * Returns an attachment by ID.
 */
export async function getAttachment(client: JiraClient, id: string): Promise<Attachment> {
  return get(client, `/rest/api/3/attachment/${id}`);
}

/**
 * Deletes an attachment.
 */
export async function deleteAttachment(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/attachment/${id}`);
}

/**
 * Returns attachment settings.
 */
export async function getAttachmentMeta(client: JiraClient): Promise<AttachmentMeta> {
  return get(client, '/rest/api/3/attachment/meta');
}

/**
 * Downloads the binary content of an attachment by ID.
 */
export async function downloadAttachmentContent(client: JiraClient, id: string): Promise<Buffer> {
  return getBuffer(client, `/rest/api/3/attachment/content/${id}`);
}

// ============================================================================
// Issue Types
// ============================================================================

/**
 * Returns all issue types.
 */
export async function getAllIssueTypes(client: JiraClient): Promise<IssueType[]> {
  return get(client, '/rest/api/3/issuetype');
}

/**
 * Creates a new issue type.
 */
export async function createIssueType(client: JiraClient, issueType: Partial<IssueType>): Promise<IssueType> {
  return post(client, '/rest/api/3/issuetype', issueType);
}

/**
 * Returns an issue type by ID.
 */
export async function getIssueType(client: JiraClient, id: string): Promise<IssueType> {
  return get(client, `/rest/api/3/issuetype/${id}`);
}

/**
 * Updates an issue type.
 */
export async function updateIssueType(
  client: JiraClient,
  id: string,
  issueType: Partial<IssueType>,
): Promise<IssueType> {
  return put(client, `/rest/api/3/issuetype/${id}`, issueType);
}

/**
 * Deletes an issue type.
 */
export async function deleteIssueType(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/issuetype/${id}`);
}

/**
 * Returns alternative issue types for the given issue type.
 */
export async function getIssueTypeAlternatives(client: JiraClient, id: string): Promise<IssueType[]> {
  return get(client, `/rest/api/3/issuetype/${id}/alternatives`);
}

/**
 * Returns issue types for a project.
 */
export async function getProjectIssueTypes(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<IssueType[]> {
  return get(client, '/rest/api/3/issuetype/project', { projectId: projectIdOrKey });
}

// ============================================================================
// Priorities
// ============================================================================

/**
 * Returns all priorities.
 */
export async function getAllPriorities(client: JiraClient): Promise<Priority[]> {
  return get(client, '/rest/api/3/priority');
}

/**
 * Creates a new priority.
 */
export async function createPriority(client: JiraClient, priority: Partial<Priority>): Promise<Priority> {
  return post(client, '/rest/api/3/priority', priority);
}

/**
 * Returns a priority by ID.
 */
export async function getPriority(client: JiraClient, id: string): Promise<Priority> {
  return get(client, `/rest/api/3/priority/${id}`);
}

/**
 * Updates a priority.
 */
export async function updatePriority(
  client: JiraClient,
  id: string,
  priority: Partial<Priority>,
): Promise<Priority> {
  return put(client, `/rest/api/3/priority/${id}`, priority);
}

/**
 * Deletes a priority.
 */
export async function deletePriority(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/priority/${id}`);
}

/**
 * Searches for priorities with pagination.
 */
export async function searchPriorities(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<PageBean<Priority>> {
  return get(client, '/rest/api/3/priority/search', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

// ============================================================================
// Resolutions
// ============================================================================

/**
 * Returns all resolutions.
 */
export async function getAllResolutions(client: JiraClient): Promise<Resolution[]> {
  return get(client, '/rest/api/3/resolution');
}

/**
 * Creates a new resolution.
 */
export async function createResolution(client: JiraClient, resolution: Partial<Resolution>): Promise<Resolution> {
  return post(client, '/rest/api/3/resolution', resolution);
}

/**
 * Returns a resolution by ID.
 */
export async function getResolution(client: JiraClient, id: string): Promise<Resolution> {
  return get(client, `/rest/api/3/resolution/${id}`);
}

/**
 * Updates a resolution.
 */
export async function updateResolution(
  client: JiraClient,
  id: string,
  resolution: Partial<Resolution>,
): Promise<Resolution> {
  return put(client, `/rest/api/3/resolution/${id}`, resolution);
}

/**
 * Deletes a resolution.
 */
export async function deleteResolution(client: JiraClient, id: string): Promise<void> {
  return del(client, `/rest/api/3/resolution/${id}`);
}

/**
 * Searches for resolutions with pagination.
 */
export async function searchResolutions(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<PageBean<Resolution>> {
  return get(client, '/rest/api/3/resolution/search', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

// ============================================================================
// Statuses
// ============================================================================

/**
 * Returns all statuses.
 */
export async function getAllStatuses(client: JiraClient): Promise<StatusDetails[]> {
  return get(client, '/rest/api/3/status');
}

/**
 * Returns a status by ID or name.
 */
export async function getStatus(client: JiraClient, idOrName: string): Promise<StatusDetails> {
  return get(client, `/rest/api/3/status/${idOrName}`);
}

/**
 * Searches for statuses with pagination.
 */
export async function searchStatuses(
  client: JiraClient,
  searchString: string,
  startAt: number,
  maxResults: number,
): Promise<PageBean<StatusDetails>> {
  const params: Record<string, string> = {};
  if (searchString) params.searchString = searchString;
  params.startAt = String(startAt);
  params.maxResults = String(maxResults);
  return get(client, '/rest/api/3/statuses/search', params);
}

/**
 * Returns all status categories.
 */
export async function getStatusCategories(client: JiraClient): Promise<StatusCategory[]> {
  return get(client, '/rest/api/3/statuscategory');
}

/**
 * Returns a status category by ID or key.
 */
export async function getStatusCategory(client: JiraClient, idOrKey: string): Promise<StatusCategory> {
  return get(client, `/rest/api/3/statuscategory/${idOrKey}`);
}

// ============================================================================
// Labels
// ============================================================================

/**
 * Returns all labels with pagination.
 */
export async function getLabels(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<PageBean<string>> {
  return get(client, '/rest/api/3/label', {
    startAt: String(startAt),
    maxResults: String(maxResults),
  });
}

// ============================================================================
// Server & Config
// ============================================================================

/**
 * Returns Jira server information.
 */
export async function getServerInfo(client: JiraClient): Promise<ServerInfo> {
  return get(client, '/rest/api/3/serverInfo');
}

/**
 * Returns the Jira configuration.
 */
export async function getConfiguration(client: JiraClient): Promise<Configuration> {
  return get(client, '/rest/api/3/configuration');
}

/**
 * Returns the announcement banner settings.
 */
export async function getAnnouncementBanner(client: JiraClient): Promise<AnnouncementBanner> {
  return get(client, '/rest/api/3/announcementBanner');
}

/**
 * Updates the announcement banner settings.
 */
export async function setAnnouncementBanner(client: JiraClient, banner: Partial<AnnouncementBanner>): Promise<void> {
  return put(client, '/rest/api/3/announcementBanner', banner);
}

/**
 * Returns audit records with pagination.
 */
export async function getAuditRecords(
  client: JiraClient,
  startAt: number,
  maxResults: number,
): Promise<AuditRecords> {
  return get(client, '/rest/api/3/auditing/record', {
    offset: String(startAt),
    limit: String(maxResults),
  });
}

/**
 * Returns all application roles.
 */
export async function getApplicationRoles(client: JiraClient): Promise<ApplicationRole[]> {
  return get(client, '/rest/api/3/applicationrole');
}

/**
 * Returns an application role by key.
 */
export async function getApplicationRole(client: JiraClient, key: string): Promise<ApplicationRole> {
  return get(client, `/rest/api/3/applicationrole/${key}`);
}

// ============================================================================
// Permissions
// ============================================================================

/**
 * Returns the permissions for the current user.
 */
export async function getMyPermissions(
  client: JiraClient,
  projectKey: string,
  issueKey: string,
): Promise<Record<string, UserPermission>> {
  const params: Record<string, string> = {};
  if (projectKey) params.projectKey = projectKey;
  if (issueKey) params.issueKey = issueKey;
  const result = await get<{ permissions: Record<string, UserPermission> }>(client, '/rest/api/3/mypermissions', params);
  return result.permissions;
}

/**
 * Returns all permissions in the system.
 */
export async function getAllPermissions(client: JiraClient): Promise<Record<string, UserPermission>> {
  const result = await get<{ permissions: Record<string, UserPermission> }>(client, '/rest/api/3/permissions');
  return result.permissions;
}

// ============================================================================
// Tasks
// ============================================================================

/**
 * Returns an async task result by ID.
 */
export async function getTask(client: JiraClient, taskId: string): Promise<TaskResult> {
  return get(client, `/rest/api/3/task/${taskId}`);
}

/**
 * Cancels an async task.
 */
export async function cancelTask(client: JiraClient, taskId: string): Promise<void> {
  return post(client, `/rest/api/3/task/${taskId}/cancel`, null);
}

import { get, post, put, del, getRaw } from './client.js';
import type { JiraClient } from './client.js';
import type { JsonBody } from '../types.js';
import type {
  CreatedIssue,
  EntityProperty,
  Filter,
  IssueType,
  PageBean,
  Project,
  ProjectCategory,
  ProjectComponent,
  ProjectFeaturesResponse,
  ProjectKeyValidation,
  ProjectRole,
  ProjectType,
  SearchRequest,
  SearchResults,
  Version,
} from './types.js';

// --- Projects ---

/**
 * Returns all projects visible to the user.
 */
export async function getAllProjects(
  client: JiraClient,
  expand?: string,
  recent?: number,
): Promise<Project[]> {
  const q: Record<string, string> = {};
  if (expand) q.expand = expand;
  if (recent != null && recent > 0) q.recent = String(recent);
  return await get(client, '/rest/api/3/project', q);
}

/**
 * Creates a new project.
 */
export async function createProject(client: JiraClient, project: Partial<Project>): Promise<CreatedIssue> {
  return await post(client, '/rest/api/3/project', project);
}

/**
 * Returns a project by ID or key.
 */
export async function getProject(
  client: JiraClient,
  projectIdOrKey: string,
  expand?: string,
): Promise<Project> {
  const q: Record<string, string> = {};
  if (expand) q.expand = expand;
  return await get(client, `/rest/api/3/project/${projectIdOrKey}`, q);
}

/**
 * Updates a project.
 */
export async function updateProject(
  client: JiraClient,
  projectIdOrKey: string,
  project: Partial<Project>,
): Promise<Project> {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}`, project);
}

/**
 * Deletes a project.
 */
export async function deleteProject(client: JiraClient, projectIdOrKey: string): Promise<void> {
  return await del(client, `/rest/api/3/project/${projectIdOrKey}`);
}

/**
 * Searches for projects using a query string.
 */
export async function searchProjects(
  client: JiraClient,
  query?: string,
  startAt?: number,
  maxResults?: number,
  expand?: string,
): Promise<PageBean<Project>> {
  const q: Record<string, string> = {};
  if (query) q.query = query;
  if (startAt != null && startAt > 0) q.startAt = String(startAt);
  if (maxResults != null && maxResults > 0) q.maxResults = String(maxResults);
  if (expand) q.expand = expand;
  return await get(client, '/rest/api/3/project/search', q);
}

/**
 * Returns recently accessed projects.
 */
export async function getRecentProjects(client: JiraClient, count?: number): Promise<Project[]> {
  const q: Record<string, string> = {};
  if (count != null && count > 0) q.count = String(count);
  return await get(client, '/rest/api/3/project/recent', q);
}

/**
 * Returns all components for a project.
 */
export async function getProjectComponents(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<ProjectComponent[]> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/components`);
}

/**
 * Returns all versions for a project.
 */
export async function getProjectVersions(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<Version[]> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/versions`);
}

/**
 * Returns paginated versions for a project.
 */
export async function getProjectVersionsPaginated(
  client: JiraClient,
  projectIdOrKey: string,
  startAt?: number,
  maxResults?: number,
): Promise<PageBean<Version>> {
  const q: Record<string, string> = {};
  if (startAt != null && startAt > 0) q.startAt = String(startAt);
  if (maxResults != null && maxResults > 0) q.maxResults = String(maxResults);
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/version`, q);
}

/**
 * Returns the valid statuses for a project.
 */
export async function getProjectStatuses(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<IssueType[]> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/statuses`);
}

/**
 * Returns all project roles for a project as a map of role name to URL.
 */
export async function getProjectRoles(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<Record<string, string>> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/role`);
}

/**
 * Returns a specific project role.
 */
export async function getProjectRole(
  client: JiraClient,
  projectIdOrKey: string,
  roleId: number,
): Promise<ProjectRole> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/role/${roleId}`);
}

/**
 * Archives a project.
 */
export async function archiveProject(client: JiraClient, projectIdOrKey: string): Promise<void> {
  return await post(client, `/rest/api/3/project/${projectIdOrKey}/archive`, null);
}

/**
 * Restores an archived project.
 */
export async function restoreProject(client: JiraClient, projectIdOrKey: string): Promise<void> {
  return await post(client, `/rest/api/3/project/${projectIdOrKey}/restore`, null);
}

/**
 * Returns all properties for a project.
 */
export async function getProjectProperties(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<EntityProperty[]> {
  const result = await get<{ keys: EntityProperty[] }>(client, `/rest/api/3/project/${projectIdOrKey}/properties`);
  return result.keys;
}

/**
 * Sets a property on a project.
 */
export async function setProjectProperty(
  client: JiraClient,
  projectIdOrKey: string,
  key: string,
  value: JsonBody,
): Promise<void> {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}/properties/${key}`, value);
}

/**
 * Deletes a property from a project.
 */
export async function deleteProjectProperty(
  client: JiraClient,
  projectIdOrKey: string,
  key: string,
): Promise<void> {
  return await del(client, `/rest/api/3/project/${projectIdOrKey}/properties/${key}`);
}

/**
 * Returns the features for a project.
 */
export async function getProjectFeatures(
  client: JiraClient,
  projectIdOrKey: string,
): Promise<ProjectFeaturesResponse> {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/features`);
}

/**
 * Sets the state of a project feature.
 */
export async function setProjectFeature(
  client: JiraClient,
  projectIdOrKey: string,
  featureKey: string,
  state: string,
): Promise<ProjectFeaturesResponse> {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}/features/${featureKey}`, {
    state,
  });
}

// --- Project Categories ---

/**
 * Returns all project categories.
 */
export async function getProjectCategories(client: JiraClient): Promise<ProjectCategory[]> {
  return await get(client, '/rest/api/3/projectCategory');
}

/**
 * Creates a new project category.
 */
export async function createProjectCategory(client: JiraClient, cat: Partial<ProjectCategory>): Promise<ProjectCategory> {
  return await post(client, '/rest/api/3/projectCategory', cat);
}

/**
 * Returns a project category by ID.
 */
export async function getProjectCategory(client: JiraClient, id: string): Promise<ProjectCategory> {
  return await get(client, `/rest/api/3/projectCategory/${id}`);
}

/**
 * Updates a project category.
 */
export async function updateProjectCategory(
  client: JiraClient,
  id: string,
  cat: Partial<ProjectCategory>,
): Promise<ProjectCategory> {
  return await put(client, `/rest/api/3/projectCategory/${id}`, cat);
}

/**
 * Deletes a project category.
 */
export async function deleteProjectCategory(client: JiraClient, id: string): Promise<void> {
  return await del(client, `/rest/api/3/projectCategory/${id}`);
}

// --- Project Validation ---

/**
 * Validates a project key.
 */
export async function validateProjectKey(
  client: JiraClient,
  key: string,
): Promise<ProjectKeyValidation> {
  return await get(client, '/rest/api/3/projectvalidate/key', { key });
}

/**
 * Returns a valid project key based on the provided key.
 */
export async function getValidProjectKey(client: JiraClient, key: string): Promise<string> {
  const data = await getRaw(client, '/rest/api/3/projectvalidate/validProjectKey', { key });
  return data.replace(/^["'\s]+|["'\s]+$/g, '');
}

/**
 * Returns a valid project name based on the provided name.
 */
export async function getValidProjectName(client: JiraClient, name: string): Promise<string> {
  const data = await getRaw(client, '/rest/api/3/projectvalidate/validProjectName', { name });
  return data.replace(/^["'\s]+|["'\s]+$/g, '');
}

// --- Project Types ---

/**
 * Returns all project types.
 */
export async function getAllProjectTypes(client: JiraClient): Promise<ProjectType[]> {
  return await get(client, '/rest/api/3/project/type');
}

/**
 * Returns all accessible project types.
 */
export async function getAccessibleProjectTypes(client: JiraClient): Promise<ProjectType[]> {
  return await get(client, '/rest/api/3/project/type/accessible');
}

// --- Search ---

/**
 * Searches for issues using JQL via GET.
 */
export async function searchJQL(
  client: JiraClient,
  jql?: string,
  startAt?: number,
  maxResults?: number,
  fields?: string[],
  expand?: string[],
): Promise<SearchResults> {
  const q: Record<string, string> = {};
  if (jql) q.jql = jql;
  if (startAt != null && startAt > 0) q.startAt = String(startAt);
  if (maxResults != null && maxResults > 0) q.maxResults = String(maxResults);
  if (fields && fields.length > 0) q.fields = fields.join(',');
  if (expand && expand.length > 0) q.expand = expand.join(',');
  return await get(client, '/rest/api/3/search/jql', q);
}

/**
 * Searches for issues using JQL via POST.
 */
export async function searchJQLPost(client: JiraClient, req: SearchRequest): Promise<SearchResults> {
  return await post(client, '/rest/api/3/search/jql', req);
}

// --- Filters ---

/**
 * Creates a new filter.
 */
export async function createFilter(client: JiraClient, filter: Partial<Filter>): Promise<Filter> {
  return await post(client, '/rest/api/3/filter', filter);
}

/**
 * Returns a filter by ID.
 */
export async function getFilter(client: JiraClient, id: string): Promise<Filter> {
  return await get(client, `/rest/api/3/filter/${id}`);
}

/**
 * Updates a filter.
 */
export async function updateFilter(
  client: JiraClient,
  id: string,
  filter: Partial<Filter>,
): Promise<Filter> {
  return await put(client, `/rest/api/3/filter/${id}`, filter);
}

/**
 * Deletes a filter.
 */
export async function deleteFilter(client: JiraClient, id: string): Promise<void> {
  return await del(client, `/rest/api/3/filter/${id}`);
}

/**
 * Returns the user's favourite filters.
 */
export async function getFavouriteFilters(client: JiraClient): Promise<Filter[]> {
  return await get(client, '/rest/api/3/filter/favourite');
}

/**
 * Returns the user's own filters.
 */
export async function getMyFilters(client: JiraClient): Promise<Filter[]> {
  return await get(client, '/rest/api/3/filter/my');
}

/**
 * Searches for filters by name.
 */
export async function searchFilters(
  client: JiraClient,
  name?: string,
  startAt?: number,
  maxResults?: number,
): Promise<PageBean<Filter>> {
  const q: Record<string, string> = {};
  if (name) q.filterName = name;
  if (startAt != null && startAt > 0) q.startAt = String(startAt);
  if (maxResults != null && maxResults > 0) q.maxResults = String(maxResults);
  return await get(client, '/rest/api/3/filter/search', q);
}

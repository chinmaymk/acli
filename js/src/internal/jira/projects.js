import { get, post, put, del, getRaw } from './client.js';

// --- Projects ---

/**
 * Returns all projects visible to the user.
 * @param {object} client
 * @param {string} [expand]
 * @param {number} [recent]
 * @returns {Promise<any[]>}
 */
export async function getAllProjects(client, expand, recent) {
  const q = {};
  if (expand) q.expand = expand;
  if (recent > 0) q.recent = String(recent);
  return await get(client, '/rest/api/3/project', q);
}

/**
 * Creates a new project.
 * @param {object} client
 * @param {object} project
 * @returns {Promise<any>}
 */
export async function createProject(client, project) {
  return await post(client, '/rest/api/3/project', project);
}

/**
 * Returns a project by ID or key.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {string} [expand]
 * @returns {Promise<any>}
 */
export async function getProject(client, projectIdOrKey, expand) {
  const q = {};
  if (expand) q.expand = expand;
  return await get(client, `/rest/api/3/project/${projectIdOrKey}`, q);
}

/**
 * Updates a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {object} project
 * @returns {Promise<any>}
 */
export async function updateProject(client, projectIdOrKey, project) {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}`, project);
}

/**
 * Deletes a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<null>}
 */
export async function deleteProject(client, projectIdOrKey) {
  return await del(client, `/rest/api/3/project/${projectIdOrKey}`);
}

/**
 * Searches for projects using a query string.
 * @param {object} client
 * @param {string} [query]
 * @param {number} [startAt]
 * @param {number} [maxResults]
 * @param {string} [expand]
 * @returns {Promise<any>}
 */
export async function searchProjects(client, query, startAt, maxResults, expand) {
  const q = {};
  if (query) q.query = query;
  if (startAt > 0) q.startAt = String(startAt);
  if (maxResults > 0) q.maxResults = String(maxResults);
  if (expand) q.expand = expand;
  return await get(client, '/rest/api/3/project/search', q);
}

/**
 * Returns recently accessed projects.
 * @param {object} client
 * @param {number} [count]
 * @returns {Promise<any[]>}
 */
export async function getRecentProjects(client, count) {
  const q = {};
  if (count > 0) q.count = String(count);
  return await get(client, '/rest/api/3/project/recent', q);
}

/**
 * Returns all components for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any[]>}
 */
export async function getProjectComponents(client, projectIdOrKey) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/components`);
}

/**
 * Returns all versions for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any[]>}
 */
export async function getProjectVersions(client, projectIdOrKey) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/versions`);
}

/**
 * Returns paginated versions for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {number} [startAt]
 * @param {number} [maxResults]
 * @returns {Promise<any>}
 */
export async function getProjectVersionsPaginated(client, projectIdOrKey, startAt, maxResults) {
  const q = {};
  if (startAt > 0) q.startAt = String(startAt);
  if (maxResults > 0) q.maxResults = String(maxResults);
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/version`, q);
}

/**
 * Returns the valid statuses for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any[]>}
 */
export async function getProjectStatuses(client, projectIdOrKey) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/statuses`);
}

/**
 * Returns all project roles for a project as a map of role name to URL.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<Record<string, string>>}
 */
export async function getProjectRoles(client, projectIdOrKey) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/role`);
}

/**
 * Returns a specific project role.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {number} roleId
 * @returns {Promise<any>}
 */
export async function getProjectRole(client, projectIdOrKey, roleId) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/role/${roleId}`);
}

/**
 * Archives a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any>}
 */
export async function archiveProject(client, projectIdOrKey) {
  return await post(client, `/rest/api/3/project/${projectIdOrKey}/archive`, null);
}

/**
 * Restores an archived project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any>}
 */
export async function restoreProject(client, projectIdOrKey) {
  return await post(client, `/rest/api/3/project/${projectIdOrKey}/restore`, null);
}

/**
 * Returns all properties for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any[]>}
 */
export async function getProjectProperties(client, projectIdOrKey) {
  const result = await get(client, `/rest/api/3/project/${projectIdOrKey}/properties`);
  return result.keys;
}

/**
 * Sets a property on a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {string} key
 * @param {*} value
 * @returns {Promise<any>}
 */
export async function setProjectProperty(client, projectIdOrKey, key, value) {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}/properties/${key}`, value);
}

/**
 * Deletes a property from a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {string} key
 * @returns {Promise<null>}
 */
export async function deleteProjectProperty(client, projectIdOrKey, key) {
  return await del(client, `/rest/api/3/project/${projectIdOrKey}/properties/${key}`);
}

/**
 * Returns the features for a project.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @returns {Promise<any>}
 */
export async function getProjectFeatures(client, projectIdOrKey) {
  return await get(client, `/rest/api/3/project/${projectIdOrKey}/features`);
}

/**
 * Sets the state of a project feature.
 * @param {object} client
 * @param {string} projectIdOrKey
 * @param {string} featureKey
 * @param {string} state
 * @returns {Promise<any>}
 */
export async function setProjectFeature(client, projectIdOrKey, featureKey, state) {
  return await put(client, `/rest/api/3/project/${projectIdOrKey}/features/${featureKey}`, { state });
}

// --- Project Categories ---

/**
 * Returns all project categories.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getProjectCategories(client) {
  return await get(client, '/rest/api/3/projectCategory');
}

/**
 * Creates a new project category.
 * @param {object} client
 * @param {object} cat
 * @returns {Promise<any>}
 */
export async function createProjectCategory(client, cat) {
  return await post(client, '/rest/api/3/projectCategory', cat);
}

/**
 * Returns a project category by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<any>}
 */
export async function getProjectCategory(client, id) {
  return await get(client, `/rest/api/3/projectCategory/${id}`);
}

/**
 * Updates a project category.
 * @param {object} client
 * @param {string} id
 * @param {object} cat
 * @returns {Promise<any>}
 */
export async function updateProjectCategory(client, id, cat) {
  return await put(client, `/rest/api/3/projectCategory/${id}`, cat);
}

/**
 * Deletes a project category.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteProjectCategory(client, id) {
  return await del(client, `/rest/api/3/projectCategory/${id}`);
}

// --- Project Validation ---

/**
 * Validates a project key.
 * @param {object} client
 * @param {string} key
 * @returns {Promise<Record<string, any>>}
 */
export async function validateProjectKey(client, key) {
  return await get(client, '/rest/api/3/projectvalidate/key', { key });
}

/**
 * Returns a valid project key based on the provided key.
 * @param {object} client
 * @param {string} key
 * @returns {Promise<string>}
 */
export async function getValidProjectKey(client, key) {
  const data = await getRaw(client, '/rest/api/3/projectvalidate/validProjectKey', { key });
  return data.replace(/^["'\s]+|["'\s]+$/g, '');
}

/**
 * Returns a valid project name based on the provided name.
 * @param {object} client
 * @param {string} name
 * @returns {Promise<string>}
 */
export async function getValidProjectName(client, name) {
  const data = await getRaw(client, '/rest/api/3/projectvalidate/validProjectName', { name });
  return data.replace(/^["'\s]+|["'\s]+$/g, '');
}

// --- Project Types ---

/**
 * Returns all project types.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getAllProjectTypes(client) {
  return await get(client, '/rest/api/3/project/type');
}

/**
 * Returns all accessible project types.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getAccessibleProjectTypes(client) {
  return await get(client, '/rest/api/3/project/type/accessible');
}

// --- Search ---

/**
 * Searches for issues using JQL via GET.
 * @param {object} client
 * @param {string} [jql]
 * @param {number} [startAt]
 * @param {number} [maxResults]
 * @param {string[]} [fields]
 * @param {string[]} [expand]
 * @returns {Promise<any>}
 */
export async function searchJQL(client, jql, startAt, maxResults, fields, expand) {
  const q = {};
  if (jql) q.jql = jql;
  if (startAt > 0) q.startAt = String(startAt);
  if (maxResults > 0) q.maxResults = String(maxResults);
  if (fields && fields.length > 0) q.fields = fields.join(',');
  if (expand && expand.length > 0) q.expand = expand.join(',');
  return await get(client, '/rest/api/3/search/jql', q);
}

/**
 * Searches for issues using JQL via POST.
 * @param {object} client
 * @param {object} req
 * @returns {Promise<any>}
 */
export async function searchJQLPost(client, req) {
  return await post(client, '/rest/api/3/search/jql', req);
}

// --- Filters ---

/**
 * Creates a new filter.
 * @param {object} client
 * @param {object} filter
 * @returns {Promise<any>}
 */
export async function createFilter(client, filter) {
  return await post(client, '/rest/api/3/filter', filter);
}

/**
 * Returns a filter by ID.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<any>}
 */
export async function getFilter(client, id) {
  return await get(client, `/rest/api/3/filter/${id}`);
}

/**
 * Updates a filter.
 * @param {object} client
 * @param {string} id
 * @param {object} filter
 * @returns {Promise<any>}
 */
export async function updateFilter(client, id, filter) {
  return await put(client, `/rest/api/3/filter/${id}`, filter);
}

/**
 * Deletes a filter.
 * @param {object} client
 * @param {string} id
 * @returns {Promise<null>}
 */
export async function deleteFilter(client, id) {
  return await del(client, `/rest/api/3/filter/${id}`);
}

/**
 * Returns the user's favourite filters.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getFavouriteFilters(client) {
  return await get(client, '/rest/api/3/filter/favourite');
}

/**
 * Returns the user's own filters.
 * @param {object} client
 * @returns {Promise<any[]>}
 */
export async function getMyFilters(client) {
  return await get(client, '/rest/api/3/filter/my');
}

/**
 * Searches for filters by name.
 * @param {object} client
 * @param {string} [name]
 * @param {number} [startAt]
 * @param {number} [maxResults]
 * @returns {Promise<any>}
 */
export async function searchFilters(client, name, startAt, maxResults) {
  const q = {};
  if (name) q.filterName = name;
  if (startAt > 0) q.startAt = String(startAt);
  if (maxResults > 0) q.maxResults = String(maxResults);
  return await get(client, '/rest/api/3/filter/search', q);
}

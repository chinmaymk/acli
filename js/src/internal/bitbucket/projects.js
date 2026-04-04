import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listProjects(client, workspace, opts) {
  let path = `/workspaces/${encodeURIComponent(workspace)}/projects`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} projectKey
 * @returns {Promise<unknown>}
 */
export async function getProject(client, workspace, projectKey) {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects/${encodeURIComponent(projectKey)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ name: string, key: string, description?: string, is_private?: boolean }} req
 * @returns {Promise<unknown>}
 */
export async function createProject(client, workspace, req) {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} projectKey
 * @returns {Promise<void>}
 */
export async function deleteProject(client, workspace, projectKey) {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects/${encodeURIComponent(projectKey)}`;
  return deleteNoContent(client, path);
}

import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listRepoWebhooks(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
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
 * @param {string} repoSlug
 * @param {string} uid
 * @returns {Promise<unknown>}
 */
export async function getRepoWebhook(client, workspace, repoSlug, uid) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ description: string, url: string, active: boolean, events: string[] }} req
 * @returns {Promise<unknown>}
 */
export async function createRepoWebhook(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} uid
 * @param {{ description: string, url: string, active: boolean, events: string[] }} req
 * @returns {Promise<unknown>}
 */
export async function updateRepoWebhook(client, workspace, repoSlug, uid, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return put(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} uid
 * @returns {Promise<void>}
 */
export async function deleteRepoWebhook(client, workspace, repoSlug, uid) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listWorkspaceWebhooks(client, workspace, opts) {
  let path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
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
 * @param {{ description: string, url: string, active: boolean, events: string[] }} req
 * @returns {Promise<unknown>}
 */
export async function createWorkspaceWebhook(client, workspace, req) {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} uid
 * @returns {Promise<void>}
 */
export async function deleteWorkspaceWebhook(client, workspace, uid) {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

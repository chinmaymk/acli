import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listDeployKeys(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys`;
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
 * @param {number} keyID
 * @returns {Promise<unknown>}
 */
export async function getDeployKey(client, workspace, repoSlug, keyID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys/${keyID}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ key: string, label: string }} req
 * @returns {Promise<unknown>}
 */
export async function createDeployKey(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} keyID
 * @returns {Promise<void>}
 */
export async function deleteDeployKey(client, workspace, repoSlug, keyID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys/${keyID}`;
  return deleteNoContent(client, path);
}

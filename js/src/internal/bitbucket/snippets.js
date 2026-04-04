import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listSnippets(client, workspace, opts) {
  let path = `/snippets/${encodeURIComponent(workspace)}`;
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
 * @param {string} encodedID
 * @returns {Promise<unknown>}
 */
export async function getSnippet(client, workspace, encodedID) {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ title: string, is_private?: boolean, scm?: string, files?: Record<string, { content: string }> }} req
 * @returns {Promise<unknown>}
 */
export async function createSnippet(client, workspace, req) {
  const path = `/snippets/${encodeURIComponent(workspace)}`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} encodedID
 * @returns {Promise<void>}
 */
export async function deleteSnippet(client, workspace, encodedID) {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return deleteNoContent(client, path);
}

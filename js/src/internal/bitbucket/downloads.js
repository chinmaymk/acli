import { get, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listDownloads(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/downloads`;
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
 * @param {string} filename
 * @returns {Promise<void>}
 */
export async function deleteDownload(client, workspace, repoSlug, filename) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/downloads/${encodeURIComponent(filename)}`;
  return deleteNoContent(client, path);
}

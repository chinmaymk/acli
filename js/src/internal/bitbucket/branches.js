import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ q?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listBranches(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

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
 * @param {string} name
 * @returns {Promise<unknown>}
 */
export async function getBranch(client, workspace, repoSlug, name) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ name: string, target: { hash: string } }} req
 * @returns {Promise<unknown>}
 */
export async function createBranch(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} name
 * @returns {Promise<void>}
 */
export async function deleteBranch(client, workspace, repoSlug, name) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ q?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listTags(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

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
 * @param {string} name
 * @returns {Promise<unknown>}
 */
export async function getTag(client, workspace, repoSlug, name) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ name: string, target: { hash: string }, message?: string }} req
 * @returns {Promise<unknown>}
 */
export async function createTag(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} name
 * @returns {Promise<void>}
 */
export async function deleteTag(client, workspace, repoSlug, name) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ role?: string, q?: string, sort?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listRepositories(client, workspace, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}`;

  const params = new URLSearchParams();
  if (opts?.role) params.set('role', opts.role);
  if (opts?.q) params.set('q', opts.q);
  if (opts?.sort) params.set('sort', opts.sort);
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
 * @returns {Promise<unknown>}
 */
export async function getRepository(client, workspace, repoSlug) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {{ scm: string, name: string, slug?: string, is_private?: boolean, description?: string, language?: string, has_issues?: boolean, has_wiki?: boolean, fork_policy?: string, project?: { key: string } }} req
 * @returns {Promise<unknown>}
 */
export async function createRepository(client, workspace, req) {
  const slug = req.slug || req.name;
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(slug)}`;
  const body = { ...req };
  delete body.slug;
  return post(client, path, body);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @returns {Promise<void>}
 */
export async function deleteRepository(client, workspace, repoSlug) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ name?: string, workspace?: { slug: string }, is_private?: boolean, description?: string, language?: string }} req
 * @returns {Promise<unknown>}
 */
export async function forkRepository(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/forks`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listForks(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/forks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

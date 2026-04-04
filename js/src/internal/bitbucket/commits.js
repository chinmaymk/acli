import { get, getRaw, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ include?: string, exclude?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listCommits(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commits`;

  const params = new URLSearchParams();
  if (opts?.include) params.set('include', opts.include);
  if (opts?.exclude) params.set('exclude', opts.exclude);
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
 * @param {string} commitHash
 * @returns {Promise<unknown>}
 */
export async function getCommit(client, workspace, repoSlug, commitHash) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commit/${encodeURIComponent(commitHash)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} commitHash
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listCommitStatuses(client, workspace, repoSlug, commitHash, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commit/${encodeURIComponent(commitHash)}/statuses`;
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
 * @param {string} spec
 * @returns {Promise<string>}
 */
export async function getDiff(client, workspace, repoSlug, spec) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/diff/${encodeURIComponent(spec)}`;
  return getRaw(client, path);
}

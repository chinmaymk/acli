import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ q?: string, sort?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listIssues(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues`;

  const params = new URLSearchParams();
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
 * @param {number} issueID
 * @returns {Promise<unknown>}
 */
export async function getIssue(client, workspace, repoSlug, issueID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ title: string, content?: { raw: string }, kind?: string, priority?: string, state?: string, assignee?: { uuid: string }, component?: { name: string }, milestone?: { name: string }, version?: { name: string } }} req
 * @returns {Promise<unknown>}
 */
export async function createIssue(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} issueID
 * @param {{ title?: string, content?: { raw: string }, kind?: string, priority?: string, state?: string, assignee?: { uuid: string } }} req
 * @returns {Promise<unknown>}
 */
export async function updateIssue(client, workspace, repoSlug, issueID, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return put(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} issueID
 * @returns {Promise<void>}
 */
export async function deleteIssue(client, workspace, repoSlug, issueID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} issueID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listIssueComments(client, workspace, repoSlug, issueID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
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
 * @param {number} issueID
 * @param {string} content
 * @returns {Promise<unknown>}
 */
export async function createIssueComment(client, workspace, repoSlug, issueID, content) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
  return post(client, path, { content: { raw: content } });
}

import { get, getRaw, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ state?: string, author?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPullRequests(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests`;

  const params = new URLSearchParams();
  const qParts = [];

  if (opts?.state) {
    const state = opts.state.toUpperCase();
    qParts.push(`state="${state}"`);
  }
  if (opts?.author) {
    const author = opts.author;
    if (author.length > 2 && author[0] === '{' && author[author.length - 1] === '}') {
      qParts.push(`author.uuid="${author}"`);
    } else {
      qParts.push(`author.nickname="${author}"`);
    }
  }
  if (qParts.length > 0) params.set('q', qParts.join(' AND '));
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
 * @param {number} prID
 * @returns {Promise<unknown>}
 */
export async function getPullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ title: string, description?: string, sourceBranch: string, destinationBranch?: string, close_source_branch?: boolean }} req
 * @returns {Promise<unknown>}
 */
export async function createPullRequest(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests`;
  const body = {
    title: req.title,
    description: req.description,
    close_source_branch: req.close_source_branch,
    source: { branch: { name: req.sourceBranch } },
  };
  if (req.destinationBranch) {
    body.destination = { branch: { name: req.destinationBranch } };
  }
  return post(client, path, body);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ title?: string, description?: string, close_source_branch?: boolean }} req
 * @returns {Promise<unknown>}
 */
export async function updatePullRequest(client, workspace, repoSlug, prID, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}`;
  return put(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<unknown>}
 */
export async function approvePullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/approve`;
  return post(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<void>}
 */
export async function unapprovePullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/approve`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<unknown>}
 */
export async function declinePullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/decline`;
  return post(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ merge_strategy?: string, close_source_branch?: boolean, message?: string }} [req]
 * @returns {Promise<unknown>}
 */
export async function mergePullRequest(client, workspace, repoSlug, prID, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/merge`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<unknown>}
 */
export async function requestChangesPullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/request-changes`;
  return post(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<void>}
 */
export async function removeRequestChangesPullRequest(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/request-changes`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPRComments(client, workspace, repoSlug, prID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/comments`;
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
 * @param {number} prID
 * @param {string} content
 * @param {{ path: string, to: number }} [inline]
 * @returns {Promise<unknown>}
 */
export async function createPRComment(client, workspace, repoSlug, prID, content, inline) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/comments`;
  const body = { content: { raw: content } };
  if (inline) {
    body.inline = { path: inline.path, to: inline.to };
  }
  return post(client, path, body);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @returns {Promise<string>}
 */
export async function getPRDiff(client, workspace, repoSlug, prID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/diff`;
  return getRaw(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPRTasks(client, workspace, repoSlug, prID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks`;
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
 * @param {number} prID
 * @param {number} taskID
 * @returns {Promise<unknown>}
 */
export async function getPRTask(client, workspace, repoSlug, prID, taskID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ content: string, commentID?: number }} req
 * @returns {Promise<unknown>}
 */
export async function createPRTask(client, workspace, repoSlug, prID, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks`;
  const body = { content: { raw: req.content } };
  if (req.commentID != null) {
    body.comment = { id: req.commentID };
  }
  return post(client, path, body);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {number} taskID
 * @param {{ content?: string, state?: string }} req
 * @returns {Promise<unknown>}
 */
export async function updatePRTask(client, workspace, repoSlug, prID, taskID, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  const body = {};
  if (req.content != null) body.content = { raw: req.content };
  if (req.state) body.state = req.state;
  return put(client, path, body);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {number} taskID
 * @returns {Promise<void>}
 */
export async function deletePRTask(client, workspace, repoSlug, prID, taskID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  return deleteNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {number} prID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPRCommits(client, workspace, repoSlug, prID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/commits`;
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
 * @param {number} prID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function getPRDiffStat(client, workspace, repoSlug, prID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/diffstat`;
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
 * @param {number} prID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPRActivity(client, workspace, repoSlug, prID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/activity`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

const BASE_URL = 'https://api.bitbucket.org/2.0';

const DEFAULT_PAGE_LEN = 50;

/**
 * @typedef {{ email: string, api_token: string }} Profile
 * @typedef {{ page?: number, pageLen?: number, all?: boolean }} PaginationOptions
 */

/**
 * Creates a Bitbucket API client.
 * @param {Profile} profile
 */
function createClient(profile) {
  return { email: profile.email, token: profile.api_token };
}

/**
 * Core HTTP request handler.
 * @param {{ email: string, token: string }} client
 * @param {string} method
 * @param {string} path
 * @param {unknown} [body]
 * @returns {Promise<unknown>}
 */
async function doRequest(client, method, path, body) {
  const url = path.startsWith('http') ? path : BASE_URL + path;

  const headers = {};

  if (client.email) {
    const credentials = Buffer.from(`${client.email}:${client.token}`).toString('base64');
    headers['Authorization'] = `Basic ${credentials}`;
  } else {
    headers['Authorization'] = `Bearer ${client.token}`;
  }

  headers['Accept'] = 'application/json';

  let bodyStr;
  if (body !== undefined) {
    bodyStr = JSON.stringify(body);
    headers['Content-Type'] = 'application/json';
  }

  const resp = await fetch(url, { method, headers, body: bodyStr });

  const text = await resp.text();

  if (!resp.ok) {
    let message;
    try {
      const parsed = JSON.parse(text);
      message = parsed?.error?.message;
    } catch {
      // fall through
    }
    const err = new Error(message || text);
    err.status = resp.status;
    throw err;
  }

  if (!text) return undefined;
  return JSON.parse(text);
}

/**
 * GET request returning parsed JSON.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 */
async function get(client, path) {
  return doRequest(client, 'GET', path);
}

/**
 * GET request without Accept header, returning raw text (for logs, diffs, etc.).
 * @param {{ email: string, token: string }} client
 * @param {string} path
 * @returns {Promise<string>}
 */
async function getRaw(client, path) {
  const url = path.startsWith('http') ? path : BASE_URL + path;

  const headers = {};
  if (client.email) {
    const credentials = Buffer.from(`${client.email}:${client.token}`).toString('base64');
    headers['Authorization'] = `Basic ${credentials}`;
  } else {
    headers['Authorization'] = `Bearer ${client.token}`;
  }

  const resp = await fetch(url, { method: 'GET', headers });
  const text = await resp.text();

  if (!resp.ok) {
    const err = new Error(text);
    err.status = resp.status;
    throw err;
  }

  return text;
}

/**
 * POST request with JSON body returning parsed JSON.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 * @param {unknown} body
 */
async function post(client, path, body) {
  return doRequest(client, 'POST', path, body);
}

/**
 * PUT request with JSON body returning parsed JSON.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 * @param {unknown} body
 */
async function put(client, path, body) {
  return doRequest(client, 'PUT', path, body);
}

/**
 * DELETE request expecting 204 No Content.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 */
async function deleteNoContent(client, path) {
  await doRequest(client, 'DELETE', path);
}

/**
 * POST request expecting 204 No Content.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 * @param {unknown} body
 */
async function postNoContent(client, path, body) {
  await doRequest(client, 'POST', path, body);
}

/**
 * Paginated GET. Follows `next` links to collect all values across pages.
 * @param {{ email: string, token: string }} client
 * @param {string} path
 * @returns {Promise<unknown[]>}
 */
async function getAll(client, path) {
  const allValues = [];
  let next = path;

  while (next) {
    const page = await get(client, next);
    if (Array.isArray(page?.values)) {
      allValues.push(...page.values);
    }
    next = page?.next ?? null;
  }

  return allValues;
}

/**
 * Appends pagination query parameters to a path.
 * @param {string} path
 * @param {PaginationOptions} [opts]
 * @returns {string}
 */
function addPaginationParams(path, opts) {
  if (!opts) return path;

  const url = new URL(path.startsWith('http') ? path : BASE_URL + path);

  if (opts.page > 0) {
    url.searchParams.set('page', String(opts.page));
  }
  if (opts.pageLen > 0) {
    url.searchParams.set('pagelen', String(opts.pageLen));
  }
  if (!url.searchParams.has('pagelen')) {
    url.searchParams.set('pagelen', String(DEFAULT_PAGE_LEN));
  }

  // Return just the path+query if it was a relative path originally
  if (!path.startsWith('http')) {
    return url.pathname + (url.search || '');
  }
  return url.toString();
}

export {
  createClient,
  doRequest,
  get,
  getRaw,
  post,
  put,
  deleteNoContent,
  postNoContent,
  getAll,
  addPaginationParams,
};

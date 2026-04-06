import type { Profile } from '../config/config.js';

const BASE_URL = 'https://api.bitbucket.org/2.0';

const DEFAULT_PAGE_LEN = 50;

export interface BitbucketClient {
  token: string;
  email: string;
}

export interface PaginationOptions {
  page?: number;
  pageLen?: number;
  all?: boolean;
}

export interface PaginatedResponse<T = any> {
  size: number;
  page: number;
  pagelen: number;
  next?: string;
  previous?: string;
  values: T[];
}

interface HttpError extends Error {
  status: number;
}

/**
 * Creates a Bitbucket API client.
 */
function createClient(profile: Profile): BitbucketClient {
  return { email: profile.email, token: profile.api_token };
}

/**
 * Core HTTP request handler.
 */
async function doRequest<T = unknown>(client: BitbucketClient, method: string, path: string, body?: unknown): Promise<T> {
  const url = path.startsWith('http') ? path : BASE_URL + path;

  const headers: Record<string, string> = {};

  if (client.email) {
    const credentials = Buffer.from(`${client.email}:${client.token}`).toString('base64');
    headers['Authorization'] = `Basic ${credentials}`;
  } else {
    headers['Authorization'] = `Bearer ${client.token}`;
  }

  headers['Accept'] = 'application/json';

  let bodyStr: string | undefined;
  if (body !== undefined) {
    bodyStr = JSON.stringify(body);
    headers['Content-Type'] = 'application/json';
  }

  const resp = await fetch(url, { method, headers, body: bodyStr });

  const text = await resp.text();

  if (!resp.ok) {
    let message: string | undefined;
    try {
      const parsed = JSON.parse(text);
      message = parsed?.error?.message;
    } catch {
      // fall through
    }
    const err = new Error(message || text) as HttpError;
    err.status = resp.status;
    throw err;
  }

  if (!text) return undefined as T;
  return JSON.parse(text) as T;
}

/**
 * GET request returning parsed JSON.
 */
async function get<T = unknown>(client: BitbucketClient, path: string): Promise<T> {
  return doRequest<T>(client, 'GET', path);
}

/**
 * GET request without Accept header, returning raw text (for logs, diffs, etc.).
 */
async function getRaw(client: BitbucketClient, path: string): Promise<string> {
  const url = path.startsWith('http') ? path : BASE_URL + path;

  const headers: Record<string, string> = {};
  if (client.email) {
    const credentials = Buffer.from(`${client.email}:${client.token}`).toString('base64');
    headers['Authorization'] = `Basic ${credentials}`;
  } else {
    headers['Authorization'] = `Bearer ${client.token}`;
  }

  const resp = await fetch(url, { method: 'GET', headers });
  const text = await resp.text();

  if (!resp.ok) {
    const err = new Error(text) as HttpError;
    err.status = resp.status;
    throw err;
  }

  return text;
}

/**
 * POST request with JSON body returning parsed JSON.
 */
async function post<T = unknown>(client: BitbucketClient, path: string, body?: unknown): Promise<T> {
  return doRequest<T>(client, 'POST', path, body);
}

/**
 * PUT request with JSON body returning parsed JSON.
 */
async function put<T = unknown>(client: BitbucketClient, path: string, body: unknown): Promise<T> {
  return doRequest<T>(client, 'PUT', path, body);
}

/**
 * DELETE request expecting 204 No Content.
 */
async function deleteNoContent(client: BitbucketClient, path: string): Promise<void> {
  await doRequest(client, 'DELETE', path);
}

/**
 * POST request expecting 204 No Content.
 */
async function postNoContent(client: BitbucketClient, path: string, body?: unknown): Promise<void> {
  await doRequest(client, 'POST', path, body);
}

/**
 * Paginated GET. Follows `next` links to collect all values across pages.
 */
async function getAll<T = unknown>(client: BitbucketClient, path: string): Promise<T[]> {
  const allValues: T[] = [];
  let nextUrl: string | null = path;

  while (nextUrl !== null) {
    const page: PaginatedResponse<T> = await get<PaginatedResponse<T>>(client, nextUrl);
    if (Array.isArray(page?.values)) {
      allValues.push(...page.values);
    }
    nextUrl = page?.next ?? null;
  }

  return allValues;
}

/**
 * Appends pagination query parameters to a path.
 */
function addPaginationParams(path: string, opts?: PaginationOptions): string {
  if (!opts) return path;

  const url = new URL(path.startsWith('http') ? path : BASE_URL + path);

  if (opts.page !== undefined && opts.page > 0) {
    url.searchParams.set('page', String(opts.page));
  }
  if (opts.pageLen !== undefined && opts.pageLen > 0) {
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

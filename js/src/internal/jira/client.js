import fs from 'fs';
import path from 'path';

/**
 * Creates a new Jira API client from a profile object.
 * @param {{ name: string, atlassian_url: string, email: string, api_token: string }} profile
 * @returns {{ baseURL: string, email: string, apiToken: string }}
 */
export function createClient(profile) {
  return {
    baseURL: profile.atlassian_url.replace(/\/+$/, ''),
    email: profile.email,
    apiToken: profile.api_token,
  };
}

/**
 * Builds a full URL from base, path, and query object.
 * @param {string} base
 * @param {string} urlPath
 * @param {Record<string, any>} [query]
 * @returns {string}
 */
export function buildURL(base, urlPath, query) {
  const full = base + urlPath;
  if (!query || Object.keys(query).length === 0) {
    return full;
  }
  const params = new URLSearchParams();
  for (const [k, v] of Object.entries(query)) {
    if (v !== undefined && v !== null) {
      params.set(k, String(v));
    }
  }
  const qs = params.toString();
  return qs ? `${full}?${qs}` : full;
}

function authHeader(client) {
  if (client.email) {
    const encoded = Buffer.from(`${client.email}:${client.apiToken}`).toString('base64');
    return `Basic ${encoded}`;
  }
  return `Bearer ${client.apiToken}`;
}

async function handleError(response) {
  const text = await response.text();
  let errorMessages = [];
  let errors = {};
  try {
    const parsed = JSON.parse(text);
    errorMessages = parsed.errorMessages || [];
    errors = parsed.errors || {};
  } catch {
    errorMessages = [text];
  }
  const parts = [];
  if (errorMessages.length > 0) {
    parts.push(errorMessages.join(', '));
  }
  if (Object.keys(errors).length > 0) {
    parts.push(JSON.stringify(errors));
  }
  const detail = parts.join(' ');
  throw new Error(`API error ${response.status}: ${detail}`);
}

async function doRequest(client, method, url, { body, headers = {} } = {}) {
  const reqHeaders = {
    Accept: 'application/json',
    Authorization: authHeader(client),
    ...headers,
  };

  const options = { method, headers: reqHeaders };
  if (body !== undefined) {
    options.body = body;
  }

  const response = await fetch(url, options);
  if (!response.ok) {
    await handleError(response);
  }
  return response;
}

/**
 * GET request. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {Record<string, any>} [query]
 * @param {*} [_responseType] - unused in JS, JSON is always parsed
 * @returns {Promise<any>}
 */
export async function get(client, urlPath, query, _responseType) {
  const url = buildURL(client.baseURL, urlPath, query);
  const response = await doRequest(client, 'GET', url);
  if (response.status === 204) return null;
  return response.json();
}

/**
 * POST request with JSON body. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {*} body
 * @returns {Promise<any>}
 */
export async function post(client, urlPath, body) {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'POST', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null;
  return response.json();
}

/**
 * PUT request with JSON body. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {*} body
 * @returns {Promise<any>}
 */
export async function put(client, urlPath, body) {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'PUT', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null;
  return response.json();
}

/**
 * DELETE request with query params. Returns null.
 * @param {object} client
 * @param {string} urlPath
 * @param {Record<string, any>} [query]
 * @returns {Promise<null>}
 */
export async function del(client, urlPath, query) {
  const url = buildURL(client.baseURL, urlPath, query);
  await doRequest(client, 'DELETE', url);
  return null;
}

/**
 * DELETE request with JSON body. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {*} body
 * @returns {Promise<any>}
 */
export async function deleteWithBody(client, urlPath, body) {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'DELETE', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null;
  return response.json();
}

/**
 * PATCH request with JSON body. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {*} body
 * @returns {Promise<any>}
 */
export async function patch(client, urlPath, body) {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'PATCH', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null;
  return response.json();
}

/**
 * Multipart file upload. Parses JSON response.
 * @param {object} client
 * @param {string} urlPath
 * @param {string} fieldName
 * @param {string} filePath
 * @returns {Promise<any>}
 */
export async function uploadFile(client, urlPath, fieldName, filePath) {
  const url = buildURL(client.baseURL, urlPath);
  const fileBuffer = fs.readFileSync(filePath);
  const fileName = path.basename(filePath);

  const form = new FormData();
  form.append(fieldName, new Blob([fileBuffer]), fileName);

  const response = await doRequest(client, 'POST', url, {
    body: form,
    headers: { 'X-Atlassian-Token': 'no-check' },
  });
  if (response.status === 204) return null;
  return response.json();
}

/**
 * GET request that returns the raw response body as a string.
 * @param {object} client
 * @param {string} urlPath
 * @param {Record<string, any>} [query]
 * @returns {Promise<string>}
 */
export async function getRaw(client, urlPath, query) {
  const url = buildURL(client.baseURL, urlPath, query);
  const response = await doRequest(client, 'GET', url);
  return response.text();
}

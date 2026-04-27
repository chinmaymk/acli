import fs from 'fs';
import path from 'path';
import type { Profile } from '../config/config.js';
import type { JsonBody } from '../types.js';

export interface JiraClient {
  baseURL: string;
  email: string;
  apiToken: string;
}

export interface APIError {
  statusCode: number;
  errorMessages: string[];
  errors: Record<string, string>;
}

/**
 * Creates a new Jira API client from a profile object.
 */
export function createClient(profile: Profile): JiraClient {
  return {
    baseURL: profile.atlassian_url.replace(/\/+$/, ''),
    email: profile.email,
    apiToken: profile.api_token,
  };
}

/**
 * Builds a full URL from base, path, and query object.
 */
export function buildURL(base: string, urlPath: string, query?: Record<string, string>): string {
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

function authHeader(client: JiraClient): string {
  if (client.email) {
    const encoded = Buffer.from(`${client.email}:${client.apiToken}`).toString('base64');
    return `Basic ${encoded}`;
  }
  return `Bearer ${client.apiToken}`;
}

interface JiraErrorBody {
  errorMessages?: string[];
  errors?: Record<string, string>;
}

async function handleError(response: Response): Promise<never> {
  const text = await response.text();
  let errorMessages: string[] = [];
  let errors: Record<string, string> = {};
  try {
    const parsed = JSON.parse(text) as JiraErrorBody;
    errorMessages = parsed.errorMessages || [];
    errors = parsed.errors || {};
  } catch {
    errorMessages = [text];
  }
  const parts: string[] = [];
  if (errorMessages.length > 0) {
    parts.push(errorMessages.join(', '));
  }
  if (Object.keys(errors).length > 0) {
    parts.push(JSON.stringify(errors));
  }
  const detail = parts.join(' ');
  throw new Error(`API error ${response.status}: ${detail}`);
}

interface DoRequestOptions {
  body?: BodyInit;
  headers?: Record<string, string>;
}

async function doRequest(
  client: JiraClient,
  method: string,
  url: string,
  { body, headers = {} }: DoRequestOptions = {},
): Promise<Response> {
  const reqHeaders: Record<string, string> = {
    Accept: 'application/json',
    Authorization: authHeader(client),
    ...headers,
  };

  const options: RequestInit = { method, headers: reqHeaders };
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
 */
export async function get<T>(
  client: JiraClient,
  urlPath: string,
  query?: Record<string, string>,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath, query);
  const response = await doRequest(client, 'GET', url);
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * POST request with JSON body. Parses JSON response.
 */
export async function post<T>(
  client: JiraClient,
  urlPath: string,
  body?: JsonBody,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'POST', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * PUT request with JSON body. Parses JSON response.
 */
export async function put<T>(
  client: JiraClient,
  urlPath: string,
  body?: JsonBody,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'PUT', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * DELETE request with query params. Returns void.
 */
export async function del(
  client: JiraClient,
  urlPath: string,
  query?: Record<string, string>,
): Promise<void> {
  const url = buildURL(client.baseURL, urlPath, query);
  await doRequest(client, 'DELETE', url);
}

/**
 * DELETE request with JSON body. Parses JSON response.
 */
export async function deleteWithBody<T>(
  client: JiraClient,
  urlPath: string,
  body?: JsonBody,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'DELETE', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * PATCH request with JSON body. Parses JSON response.
 */
export async function patch<T>(
  client: JiraClient,
  urlPath: string,
  body?: JsonBody,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath);
  const response = await doRequest(client, 'PATCH', url, {
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * Multipart file upload. Parses JSON response.
 */
export async function uploadFile<T>(
  client: JiraClient,
  urlPath: string,
  fieldName: string,
  filePath: string,
): Promise<T> {
  const url = buildURL(client.baseURL, urlPath);
  const fileBuffer = fs.readFileSync(filePath);
  const fileName = path.basename(filePath);

  const form = new FormData();
  form.append(fieldName, new Blob([fileBuffer]), fileName);

  // Node's FormData (web spec) is an acceptable BodyInit, but TypeScript's
  // BodyInit type lacks the FormData overload here, so cast through the
  // structurally compatible interface explicitly.
  const response = await doRequest(client, 'POST', url, {
    body: form as FormData,
    headers: { 'X-Atlassian-Token': 'no-check' },
  });
  if (response.status === 204) return null as T;
  return response.json() as Promise<T>;
}

/**
 * GET request that returns the raw response body as a string.
 */
export async function getRaw(
  client: JiraClient,
  urlPath: string,
  query?: Record<string, string>,
): Promise<string> {
  const url = buildURL(client.baseURL, urlPath, query);
  const response = await doRequest(client, 'GET', url);
  return response.text();
}

import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Environment, CreateEnvironmentRequest } from './types.js';

export async function listEnvironments(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<Environment[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/environments`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<Environment>(client, path);
  }

  const page = await get<PaginatedResponse<Environment>>(client, path);
  return page.values ?? [];
}

export async function getEnvironment(client: BitbucketClient, workspace: string, repoSlug: string, envUUID: string): Promise<Environment> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/environments/${encodeURIComponent(envUUID)}`;
  return get(client, path);
}

export async function createEnvironment(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateEnvironmentRequest): Promise<Environment> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/environments`;
  return post(client, path, req);
}

export async function deleteEnvironment(client: BitbucketClient, workspace: string, repoSlug: string, envUUID: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/environments/${encodeURIComponent(envUUID)}`;
  return deleteNoContent(client, path);
}

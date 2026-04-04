import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface CreateProjectRequest {
  name: string;
  key: string;
  description?: string;
  is_private?: boolean;
}

export async function listProjects(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/projects`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getProject(client: BitbucketClient, workspace: string, projectKey: string): Promise<any> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects/${encodeURIComponent(projectKey)}`;
  return get(client, path);
}

export async function createProject(client: BitbucketClient, workspace: string, req: CreateProjectRequest): Promise<any> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects`;
  return post(client, path, req);
}

export async function deleteProject(client: BitbucketClient, workspace: string, projectKey: string): Promise<void> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/projects/${encodeURIComponent(projectKey)}`;
  return deleteNoContent(client, path);
}

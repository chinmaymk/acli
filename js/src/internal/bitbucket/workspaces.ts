import { get, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export async function listWorkspaces(client: BitbucketClient, opts?: PaginationOptions): Promise<unknown[]> {
  let path = '/workspaces';
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getWorkspace(client: BitbucketClient, workspace: string): Promise<any> {
  const path = `/workspaces/${encodeURIComponent(workspace)}`;
  return get(client, path);
}

export async function listWorkspaceMembers(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/members`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function listWorkspacePermissions(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/permissions`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

import { get, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { Workspace, WorkspaceMember, WorkspacePermission } from './types.js';

export async function listWorkspaces(client: BitbucketClient, opts?: PaginationOptions): Promise<Workspace[]> {
  let path = '/workspaces';
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<Workspace[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getWorkspace(client: BitbucketClient, workspace: string): Promise<Workspace> {
  const path = `/workspaces/${encodeURIComponent(workspace)}`;
  return get(client, path);
}

export async function listWorkspaceMembers(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<WorkspaceMember[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/members`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<WorkspaceMember[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function listWorkspacePermissions(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<WorkspacePermission[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/permissions`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<WorkspacePermission[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

import { get, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Workspace, WorkspaceMember, WorkspacePermission } from './types.js';

export async function listWorkspaces(client: BitbucketClient, opts?: PaginationOptions): Promise<Workspace[]> {
  let path = '/workspaces';
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<Workspace>(client, path);
  }

  const page = await get<PaginatedResponse<Workspace>>(client, path);
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
    return getAll<WorkspaceMember>(client, path);
  }

  const page = await get<PaginatedResponse<WorkspaceMember>>(client, path);
  return page.values ?? [];
}

export async function listWorkspacePermissions(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<WorkspacePermission[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/permissions`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<WorkspacePermission>(client, path);
  }

  const page = await get<PaginatedResponse<WorkspacePermission>>(client, path);
  return page.values ?? [];
}

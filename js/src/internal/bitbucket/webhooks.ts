import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { BBWebhook, CreateWebhookRequest } from './types.js';

export async function listRepoWebhooks(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<BBWebhook[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<BBWebhook[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string): Promise<BBWebhook> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return get(client, path);
}

export async function createRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateWebhookRequest): Promise<BBWebhook> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
  return post(client, path, req);
}

export async function updateRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string, req: CreateWebhookRequest): Promise<BBWebhook> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return put(client, path, req);
}

export async function deleteRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

export async function listWorkspaceWebhooks(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<BBWebhook[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<BBWebhook[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createWorkspaceWebhook(client: BitbucketClient, workspace: string, req: CreateWebhookRequest): Promise<BBWebhook> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
  return post(client, path, req);
}

export async function deleteWorkspaceWebhook(client: BitbucketClient, workspace: string, uid: string): Promise<void> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface WebhookRequest {
  description: string;
  url: string;
  active: boolean;
  events: string[];
}

export async function listRepoWebhooks(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return get(client, path);
}

export async function createRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, req: WebhookRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks`;
  return post(client, path, req);
}

export async function updateRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string, req: WebhookRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return put(client, path, req);
}

export async function deleteRepoWebhook(client: BitbucketClient, workspace: string, repoSlug: string, uid: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

export async function listWorkspaceWebhooks(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createWorkspaceWebhook(client: BitbucketClient, workspace: string, req: WebhookRequest): Promise<any> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks`;
  return post(client, path, req);
}

export async function deleteWorkspaceWebhook(client: BitbucketClient, workspace: string, uid: string): Promise<void> {
  const path = `/workspaces/${encodeURIComponent(workspace)}/hooks/${encodeURIComponent(uid)}`;
  return deleteNoContent(client, path);
}

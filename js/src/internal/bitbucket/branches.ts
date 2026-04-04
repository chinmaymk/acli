import { get, post, deleteNoContent, getAll } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface ListBranchesOptions extends PaginationOptions {
  q?: string;
}

export interface CreateBranchRequest {
  name: string;
  target: { hash: string };
}

export interface CreateTagRequest {
  name: string;
  target: { hash: string };
  message?: string;
}

export async function listBranches(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListBranchesOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getBranch(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return get(client, path);
}

export async function createBranch(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateBranchRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;
  return post(client, path, req);
}

export async function deleteBranch(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

export async function listTags(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListBranchesOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getTag(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return get(client, path);
}

export async function createTag(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateTagRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;
  return post(client, path, req);
}

export async function deleteTag(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

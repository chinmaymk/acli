import { get, post, deleteNoContent, getAll } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Branch, Tag, CreateBranchRequest, CreateTagRequest } from './types.js';

export interface ListBranchesOptions extends PaginationOptions {
  q?: string;
}

export async function listBranches(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListBranchesOptions): Promise<Branch[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll<Branch>(client, path);
  }

  const page = await get<PaginatedResponse<Branch>>(client, path);
  return page.values ?? [];
}

export async function getBranch(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<Branch> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return get<Branch>(client, path);
}

export async function createBranch(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateBranchRequest): Promise<Branch> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches`;
  return post<Branch>(client, path, req);
}

export async function deleteBranch(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/branches/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

export async function listTags(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListBranchesOptions): Promise<Tag[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll<Tag>(client, path);
  }

  const page = await get<PaginatedResponse<Tag>>(client, path);
  return page.values ?? [];
}

export async function getTag(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<Tag> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return get<Tag>(client, path);
}

export async function createTag(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateTagRequest): Promise<Tag> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags`;
  return post<Tag>(client, path, req);
}

export async function deleteTag(client: BitbucketClient, workspace: string, repoSlug: string, name: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/refs/tags/${encodeURIComponent(name)}`;
  return deleteNoContent(client, path);
}

import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Repository, CreateRepoRequest, ForkRepoRequest } from './types.js';

export interface ListRepositoriesOptions extends PaginationOptions {
  role?: string;
  q?: string;
  sort?: string;
}

export async function listRepositories(client: BitbucketClient, workspace: string, opts?: ListRepositoriesOptions): Promise<Repository[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}`;

  const params = new URLSearchParams();
  if (opts?.role) params.set('role', opts.role);
  if (opts?.q) params.set('q', opts.q);
  if (opts?.sort) params.set('sort', opts.sort);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll<Repository>(client, path);
  }

  const page = await get<PaginatedResponse<Repository>>(client, path);
  return page.values ?? [];
}

export async function getRepository(client: BitbucketClient, workspace: string, repoSlug: string): Promise<Repository> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}`;
  return get(client, path);
}

export async function createRepository(client: BitbucketClient, workspace: string, req: CreateRepoRequest): Promise<Repository> {
  const slug = req.name;
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(slug)}`;
  const body: Record<string, unknown> = { ...req };
  return post(client, path, body);
}

export async function deleteRepository(client: BitbucketClient, workspace: string, repoSlug: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}`;
  return deleteNoContent(client, path);
}

export async function forkRepository(client: BitbucketClient, workspace: string, repoSlug: string, req: ForkRepoRequest): Promise<Repository> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/forks`;
  return post(client, path, req);
}

export async function listForks(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<Repository[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/forks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<Repository>(client, path);
  }

  const page = await get<PaginatedResponse<Repository>>(client, path);
  return page.values ?? [];
}

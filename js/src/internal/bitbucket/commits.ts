import { get, getRaw, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Commit, CommitStatus } from './types.js';

export interface ListCommitsOptions extends PaginationOptions {
  include?: string;
  exclude?: string;
}

export async function listCommits(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListCommitsOptions): Promise<Commit[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commits`;

  const params = new URLSearchParams();
  if (opts?.include) params.set('include', opts.include);
  if (opts?.exclude) params.set('exclude', opts.exclude);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll<Commit>(client, path);
  }

  const page = await get<PaginatedResponse<Commit>>(client, path);
  return page.values ?? [];
}

export async function getCommit(client: BitbucketClient, workspace: string, repoSlug: string, commitHash: string): Promise<Commit> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commit/${encodeURIComponent(commitHash)}`;
  return get(client, path);
}

export async function listCommitStatuses(client: BitbucketClient, workspace: string, repoSlug: string, commitHash: string, opts?: PaginationOptions): Promise<CommitStatus[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/commit/${encodeURIComponent(commitHash)}/statuses`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<CommitStatus>(client, path);
  }

  const page = await get<PaginatedResponse<CommitStatus>>(client, path);
  return page.values ?? [];
}

export async function getDiff(client: BitbucketClient, workspace: string, repoSlug: string, spec: string): Promise<string> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/diff/${encodeURIComponent(spec)}`;
  return getRaw(client, path);
}

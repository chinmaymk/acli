import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { BranchRestriction, CreateBranchRestrictionRequest } from './types.js';

export async function listBranchRestrictions(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<BranchRestriction[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<BranchRestriction[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, id: number): Promise<BranchRestriction> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions/${id}`;
  return get(client, path);
}

export async function createBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateBranchRestrictionRequest): Promise<BranchRestriction> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions`;
  return post(client, path, req);
}

export async function deleteBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, id: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions/${id}`;
  return deleteNoContent(client, path);
}

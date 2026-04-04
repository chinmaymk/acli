import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface CreateBranchRestrictionRequest {
  kind: string;
  pattern: string;
  value?: number;
}

export async function listBranchRestrictions(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, id: number): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions/${id}`;
  return get(client, path);
}

export async function createBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateBranchRestrictionRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions`;
  return post(client, path, req);
}

export async function deleteBranchRestriction(client: BitbucketClient, workspace: string, repoSlug: string, id: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/branch-restrictions/${id}`;
  return deleteNoContent(client, path);
}

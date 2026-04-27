import { get, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Download } from './types.js';

export async function listDownloads(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<Download[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/downloads`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<Download>(client, path);
  }

  const page = await get<PaginatedResponse<Download>>(client, path);
  return page.values ?? [];
}

export async function deleteDownload(client: BitbucketClient, workspace: string, repoSlug: string, filename: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/downloads/${encodeURIComponent(filename)}`;
  return deleteNoContent(client, path);
}

import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface CreateSnippetRequest {
  title: string;
  is_private?: boolean;
  scm?: string;
  files?: Record<string, { content: string }>;
}

export async function listSnippets(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/snippets/${encodeURIComponent(workspace)}`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getSnippet(client: BitbucketClient, workspace: string, encodedID: string): Promise<any> {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return get(client, path);
}

export async function createSnippet(client: BitbucketClient, workspace: string, req: CreateSnippetRequest): Promise<any> {
  const path = `/snippets/${encodeURIComponent(workspace)}`;
  return post(client, path, req);
}

export async function deleteSnippet(client: BitbucketClient, workspace: string, encodedID: string): Promise<void> {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return deleteNoContent(client, path);
}

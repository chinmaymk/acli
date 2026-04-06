import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions, PaginatedResponse } from './client.js';
import type { Snippet, CreateSnippetRequest } from './types.js';

export async function listSnippets(client: BitbucketClient, workspace: string, opts?: PaginationOptions): Promise<Snippet[]> {
  let path = `/snippets/${encodeURIComponent(workspace)}`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll<Snippet>(client, path);
  }

  const page = await get<PaginatedResponse<Snippet>>(client, path);
  return page.values ?? [];
}

export async function getSnippet(client: BitbucketClient, workspace: string, encodedID: string): Promise<Snippet> {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return get(client, path);
}

export async function createSnippet(client: BitbucketClient, workspace: string, req: CreateSnippetRequest): Promise<Snippet> {
  const path = `/snippets/${encodeURIComponent(workspace)}`;
  return post(client, path, req);
}

export async function deleteSnippet(client: BitbucketClient, workspace: string, encodedID: string): Promise<void> {
  const path = `/snippets/${encodeURIComponent(workspace)}/${encodeURIComponent(encodedID)}`;
  return deleteNoContent(client, path);
}

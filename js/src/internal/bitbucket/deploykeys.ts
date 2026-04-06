import { get, post, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { DeployKey, CreateDeployKeyRequest } from './types.js';

export async function listDeployKeys(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<DeployKey[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<DeployKey[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getDeployKey(client: BitbucketClient, workspace: string, repoSlug: string, keyID: number): Promise<DeployKey> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys/${keyID}`;
  return get(client, path);
}

export async function createDeployKey(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateDeployKeyRequest): Promise<DeployKey> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys`;
  return post(client, path, req);
}

export async function deleteDeployKey(client: BitbucketClient, workspace: string, repoSlug: string, keyID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deploy-keys/${keyID}`;
  return deleteNoContent(client, path);
}

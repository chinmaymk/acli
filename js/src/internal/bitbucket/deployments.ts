import { get, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { Deployment } from './types.js';

export async function listDeployments(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<Deployment[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deployments`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<Deployment[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getDeployment(client: BitbucketClient, workspace: string, repoSlug: string, deploymentUUID: string): Promise<Deployment> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/deployments/${encodeURIComponent(deploymentUUID)}`;
  return get(client, path);
}

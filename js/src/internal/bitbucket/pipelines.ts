import { get, getRaw, post, deleteNoContent, postNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface ListPipelinesOptions extends PaginationOptions {
  status?: string;
  sort?: string;
}

export interface BranchPipelineTarget {
  type: string;
  ref_type: string;
  ref_name: string;
  selector?: { type: string; pattern: string };
}

export interface PipelineRequest {
  target: BranchPipelineTarget;
}

export async function listPipelines(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListPipelinesOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines`;

  const params = new URLSearchParams();
  if (opts?.status) params.set('status', opts.status);
  if (opts?.sort) params.set('sort', opts.sort);
  if (!params.has('sort')) params.set('sort', '-created_on');
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getPipeline(client: BitbucketClient, workspace: string, repoSlug: string, pipelineUUID: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}`;
  return get(client, path);
}

export function newBranchPipelineRequest(branch: string): PipelineRequest {
  return {
    target: {
      type: 'pipeline_ref_target',
      ref_type: 'branch',
      ref_name: branch,
    },
  };
}

export function newCustomPipelineRequest(branch: string, pattern: string): PipelineRequest {
  const req = newBranchPipelineRequest(branch);
  req.target.selector = { type: 'custom', pattern };
  return req;
}

export async function runPipeline(client: BitbucketClient, workspace: string, repoSlug: string, req: PipelineRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines`;
  return post(client, path, req);
}

export async function stopPipeline(client: BitbucketClient, workspace: string, repoSlug: string, pipelineUUID: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/stopPipeline`;
  return postNoContent(client, path);
}

export async function listPipelineSteps(client: BitbucketClient, workspace: string, repoSlug: string, pipelineUUID: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getPipelineStep(client: BitbucketClient, workspace: string, repoSlug: string, pipelineUUID: string, stepUUID: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps/${encodeURIComponent(stepUUID)}`;
  return get(client, path);
}

export async function getStepLog(client: BitbucketClient, workspace: string, repoSlug: string, pipelineUUID: string, stepUUID: string): Promise<string> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps/${encodeURIComponent(stepUUID)}/log`;
  return getRaw(client, path);
}

export async function listPipelineVariables(client: BitbucketClient, workspace: string, repoSlug: string, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createPipelineVariable(client: BitbucketClient, workspace: string, repoSlug: string, key: string, value: string, secured: boolean): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables`;
  return post(client, path, { key, value, secured });
}

export async function deletePipelineVariable(client: BitbucketClient, workspace: string, repoSlug: string, variableUUID: string): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables/${encodeURIComponent(variableUUID)}`;
  return deleteNoContent(client, path);
}

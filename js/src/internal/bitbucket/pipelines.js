import { get, getRaw, post, deleteNoContent, postNoContent, getAll, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ status?: string, sort?: string, page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPipelines(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines`;

  const params = new URLSearchParams();
  if (opts?.status) params.set('status', opts.status);
  if (opts?.sort) params.set('sort', opts.sort);
  if (!params.has('sort')) params.set('sort', '-created_on');
  if (opts?.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} pipelineUUID
 * @returns {Promise<unknown>}
 */
export async function getPipeline(client, workspace, repoSlug, pipelineUUID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}`;
  return get(client, path);
}

/**
 * Creates a pipeline target for a branch.
 * @param {string} branch
 * @returns {{ target: { type: string, ref_type: string, ref_name: string } }}
 */
export function newBranchPipelineRequest(branch) {
  return {
    target: {
      type: 'pipeline_ref_target',
      ref_type: 'branch',
      ref_name: branch,
    },
  };
}

/**
 * Creates a pipeline target for a custom pipeline on a branch.
 * @param {string} branch
 * @param {string} pattern
 * @returns {{ target: { type: string, ref_type: string, ref_name: string, selector: { type: string, pattern: string } } }}
 */
export function newCustomPipelineRequest(branch, pattern) {
  const req = newBranchPipelineRequest(branch);
  req.target.selector = { type: 'custom', pattern };
  return req;
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ target: object }} req
 * @returns {Promise<unknown>}
 */
export async function runPipeline(client, workspace, repoSlug, req) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines`;
  return post(client, path, req);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} pipelineUUID
 * @returns {Promise<void>}
 */
export async function stopPipeline(client, workspace, repoSlug, pipelineUUID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/stopPipeline`;
  return postNoContent(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} pipelineUUID
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPipelineSteps(client, workspace, repoSlug, pipelineUUID, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} pipelineUUID
 * @param {string} stepUUID
 * @returns {Promise<unknown>}
 */
export async function getPipelineStep(client, workspace, repoSlug, pipelineUUID, stepUUID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps/${encodeURIComponent(stepUUID)}`;
  return get(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} pipelineUUID
 * @param {string} stepUUID
 * @returns {Promise<string>}
 */
export async function getStepLog(client, workspace, repoSlug, pipelineUUID, stepUUID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines/${encodeURIComponent(pipelineUUID)}/steps/${encodeURIComponent(stepUUID)}/log`;
  return getRaw(client, path);
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<unknown[]>}
 */
export async function listPipelineVariables(client, workspace, repoSlug, opts) {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} key
 * @param {string} value
 * @param {boolean} secured
 * @returns {Promise<unknown>}
 */
export async function createPipelineVariable(client, workspace, repoSlug, key, value, secured) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables`;
  return post(client, path, { key, value, secured });
}

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} repoSlug
 * @param {string} variableUUID
 * @returns {Promise<void>}
 */
export async function deletePipelineVariable(client, workspace, repoSlug, variableUUID) {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pipelines_config/variables/${encodeURIComponent(variableUUID)}`;
  return deleteNoContent(client, path);
}

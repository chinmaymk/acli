import { get, addPaginationParams } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @param {string} workspace
 * @param {string} query
 * @param {{ page?: number, pageLen?: number, all?: boolean }} [opts]
 * @returns {Promise<{ size: number, page?: number, pagelen?: number, next?: string, values: unknown[] }>}
 */
export async function searchCode(client, workspace, query, opts) {
  const params = new URLSearchParams();
  params.set('search_query', query);
  if (opts?.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const basePath = `/workspaces/${encodeURIComponent(workspace)}/search/code?${params.toString()}`;

  if (opts?.all) {
    const allValues = [];
    let next = basePath;
    let totalSize = 0;

    while (next) {
      const result = await get(client, next);
      if (Array.isArray(result?.values)) {
        allValues.push(...result.values);
      }
      totalSize = result?.size ?? totalSize;
      next = result?.next ?? null;
    }

    return { size: totalSize, values: allValues };
  }

  return get(client, basePath);
}

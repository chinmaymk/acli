import { get } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { SearchResponse } from './types.js';

export async function searchCode(client: BitbucketClient, workspace: string, query: string, opts?: PaginationOptions): Promise<SearchResponse> {
  const params = new URLSearchParams();
  params.set('search_query', query);
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const basePath = `/workspaces/${encodeURIComponent(workspace)}/search/code?${params.toString()}`;

  if (opts?.all) {
    const allValues: SearchResponse['values'] = [];
    let nextUrl: string | null = basePath;
    let totalSize = 0;

    while (nextUrl !== null) {
      const result: SearchResponse = await get<SearchResponse>(client, nextUrl);
      if (Array.isArray(result?.values)) {
        allValues.push(...result.values);
      }
      totalSize = result?.size ?? totalSize;
      nextUrl = result?.next ?? null;
    }

    return { size: totalSize, page: 1, pagelen: allValues.length, next: '', values: allValues };
  }

  return get<SearchResponse>(client, basePath);
}

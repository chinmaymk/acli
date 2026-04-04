import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';

export interface ListIssuesOptions extends PaginationOptions {
  q?: string;
  sort?: string;
}

export interface CreateIssueRequest {
  title: string;
  content?: { raw: string };
  kind?: string;
  priority?: string;
  state?: string;
  assignee?: { uuid: string };
  component?: { name: string };
  milestone?: { name: string };
  version?: { name: string };
}

export interface UpdateIssueRequest {
  title?: string;
  content?: { raw: string };
  kind?: string;
  priority?: string;
  state?: string;
  assignee?: { uuid: string };
}

export async function listIssues(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListIssuesOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues`;

  const params = new URLSearchParams();
  if (opts?.q) params.set('q', opts.q);
  if (opts?.sort) params.set('sort', opts.sort);
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

export async function getIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return get(client, path);
}

export async function createIssue(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateIssueRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues`;
  return post(client, path, req);
}

export async function updateIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, req: UpdateIssueRequest): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return put(client, path, req);
}

export async function deleteIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return deleteNoContent(client, path);
}

export async function listIssueComments(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createIssueComment(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, content: string): Promise<any> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
  return post(client, path, { content: { raw: content } });
}

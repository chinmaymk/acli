import { get, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { BBIssue, CreateIssueRequest, UpdateIssueRequest, IssueComment } from './types.js';

export interface ListIssuesOptions extends PaginationOptions {
  q?: string;
  sort?: string;
}

export async function listIssues(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListIssuesOptions): Promise<BBIssue[]> {
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
    return getAll(client, path) as Promise<BBIssue[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number): Promise<BBIssue> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return get(client, path);
}

export async function createIssue(client: BitbucketClient, workspace: string, repoSlug: string, req: CreateIssueRequest): Promise<BBIssue> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues`;
  return post(client, path, req);
}

export async function updateIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, req: UpdateIssueRequest): Promise<BBIssue> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return put(client, path, req);
}

export async function deleteIssue(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}`;
  return deleteNoContent(client, path);
}

export async function listIssueComments(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, opts?: PaginationOptions): Promise<IssueComment[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<IssueComment[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createIssueComment(client: BitbucketClient, workspace: string, repoSlug: string, issueID: number, content: string): Promise<IssueComment> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/issues/${issueID}/comments`;
  return post(client, path, { content: { raw: content } });
}

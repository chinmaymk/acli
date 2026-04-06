import { get, getRaw, post, put, deleteNoContent, getAll, addPaginationParams } from './client.js';
import type { BitbucketClient, PaginationOptions } from './client.js';
import type { PullRequest, CreatePRRequest, UpdatePRRequest, MergePRRequest, Participant, PRComment, InlineCommentParams, PRTask, CreatePRTaskRequest, UpdatePRTaskRequest } from './types.js';

export interface ListPullRequestsOptions extends PaginationOptions {
  state?: string;
  author?: string;
}

export async function listPullRequests(client: BitbucketClient, workspace: string, repoSlug: string, opts?: ListPullRequestsOptions): Promise<PullRequest[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests`;

  const params = new URLSearchParams();
  const qParts: string[] = [];

  if (opts?.state) {
    const state = opts.state.toUpperCase();
    qParts.push(`state="${state}"`);
  }
  if (opts?.author) {
    const author = opts.author;
    if (author.length > 2 && author[0] === '{' && author[author.length - 1] === '}') {
      qParts.push(`author.uuid="${author}"`);
    } else {
      qParts.push(`author.nickname="${author}"`);
    }
  }
  if (qParts.length > 0) params.set('q', qParts.join(' AND '));
  if (opts?.page !== undefined && opts.page > 0) params.set('page', String(opts.page));
  if (opts?.pageLen !== undefined && opts.pageLen > 0) params.set('pagelen', String(opts.pageLen));
  if (!params.has('pagelen')) params.set('pagelen', '50');

  const qs = params.toString();
  if (qs) path += '?' + qs;

  if (opts?.all) {
    return getAll(client, path) as Promise<PullRequest[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getPullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<PullRequest> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}`;
  return get(client, path);
}

export async function createPullRequest(client: BitbucketClient, workspace: string, repoSlug: string, req: CreatePRRequest): Promise<PullRequest> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests`;
  const body: Record<string, unknown> = {
    title: req.title,
    description: req.description,
    close_source_branch: req.close_source_branch,
    source: { branch: { name: req.source_branch } },
  };
  if (req.destination_branch) {
    body['destination'] = { branch: { name: req.destination_branch } };
  }
  return post(client, path, body);
}

export async function updatePullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, req: UpdatePRRequest): Promise<PullRequest> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}`;
  return put(client, path, req);
}

export async function approvePullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<Participant> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/approve`;
  return post(client, path);
}

export async function unapprovePullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/approve`;
  return deleteNoContent(client, path);
}

export async function declinePullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<PullRequest> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/decline`;
  return post(client, path);
}

export async function mergePullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, req?: MergePRRequest): Promise<PullRequest> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/merge`;
  return post(client, path, req);
}

export async function requestChangesPullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<Participant> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/request-changes`;
  return post(client, path);
}

export async function removeRequestChangesPullRequest(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/request-changes`;
  return deleteNoContent(client, path);
}

export async function listPRComments(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, opts?: PaginationOptions): Promise<PRComment[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/comments`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<PRComment[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function createPRComment(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, content: string, inline?: InlineCommentParams): Promise<PRComment> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/comments`;
  const body: Record<string, unknown> = { content: { raw: content } };
  if (inline) {
    body['inline'] = { path: inline.path, to: inline.to };
  }
  return post(client, path, body);
}

export async function createPRCommentInline(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, content: string, inline: InlineCommentParams): Promise<PRComment> {
  return createPRComment(client, workspace, repoSlug, prID, content, inline);
}

export async function getPRDiff(client: BitbucketClient, workspace: string, repoSlug: string, prID: number): Promise<string> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/diff`;
  return getRaw(client, path);
}

export async function listPRTasks(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, opts?: PaginationOptions): Promise<PRTask[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path) as Promise<PRTask[]>;
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getPRTask(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, taskID: number): Promise<PRTask> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  return get(client, path);
}

export async function createPRTask(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, req: CreatePRTaskRequest): Promise<PRTask> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks`;
  const body: Record<string, unknown> = { content: { raw: req.content } };
  if (req.comment_id != null) {
    body['comment'] = { id: req.comment_id };
  }
  return post(client, path, body);
}

export async function updatePRTask(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, taskID: number, req: UpdatePRTaskRequest): Promise<PRTask> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  const body: Record<string, unknown> = {};
  if (req.content != null) body['content'] = { raw: req.content };
  if (req.state) body['state'] = req.state;
  return put(client, path, body);
}

export async function deletePRTask(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, taskID: number): Promise<void> {
  const path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/tasks/${taskID}`;
  return deleteNoContent(client, path);
}

export async function listPRCommits(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/commits`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function getPRDiffStat(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/diffstat`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

export async function listPRActivity(client: BitbucketClient, workspace: string, repoSlug: string, prID: number, opts?: PaginationOptions): Promise<unknown[]> {
  let path = `/repositories/${encodeURIComponent(workspace)}/${encodeURIComponent(repoSlug)}/pullrequests/${prID}/activity`;
  path = addPaginationParams(path, opts);

  if (opts?.all) {
    return getAll(client, path);
  }

  const page = await get(client, path);
  return page.values ?? [];
}

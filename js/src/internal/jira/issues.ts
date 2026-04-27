import { get, post, put, del, uploadFile } from './client.js';
import type { JiraClient } from './client.js';
import type { ADFNode, JsonBody, JsonValue } from '../types.js';
import type {
  Attachment,
  BulkIssueCreateRequest,
  BulkIssueCreateResponse,
  ChangelogPage,
  Comment,
  CommentPage,
  CreateMeta,
  CreatedIssue,
  EntityProperty,
  IssueDetailed,
  IssueNotifyRequest,
  IssueUpdateDetails,
  RemoteLink,
  TransitionsResponse,
  Visibility,
  Votes,
  Watches,
  Worklog,
  WorklogPage,
} from './types.js';

// createIssue creates a new Jira issue.
export async function createIssue(client: JiraClient, details: IssueUpdateDetails): Promise<CreatedIssue> {
  return await post(client, '/rest/api/3/issue', details);
}

// bulkCreateIssues creates multiple issues in a single request.
export async function bulkCreateIssues(client: JiraClient, req: BulkIssueCreateRequest): Promise<BulkIssueCreateResponse> {
  return await post(client, '/rest/api/3/issue/bulk', req);
}

// getIssue retrieves a single issue by ID or key.
export async function getIssue(
  client: JiraClient,
  issueIdOrKey: string,
  fields?: string[],
  expand?: string[],
): Promise<IssueDetailed> {
  const query: Record<string, string> = {};
  if (fields && fields.length > 0) {
    query.fields = fields.join(',');
  }
  if (expand && expand.length > 0) {
    query.expand = expand.join(',');
  }
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}`, query);
}

// editIssue edits an existing issue.
export async function editIssue(
  client: JiraClient,
  issueIdOrKey: string,
  details: IssueUpdateDetails,
  notifyUsers?: boolean,
): Promise<void> {
  const query: Record<string, string> = {};
  if (!notifyUsers) {
    query.notifyUsers = 'false';
  }
  let path = `/rest/api/3/issue/${issueIdOrKey}`;
  if (Object.keys(query).length > 0) {
    const params = new URLSearchParams(query);
    path = path + '?' + params.toString();
  }
  return await put(client, path, details);
}

// deleteIssue deletes an issue.
export async function deleteIssue(
  client: JiraClient,
  issueIdOrKey: string,
  deleteSubtasks?: boolean,
): Promise<void> {
  const query: Record<string, string> = {};
  if (deleteSubtasks) {
    query.deleteSubtasks = 'true';
  }
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}`, query);
}

// assignIssue assigns an issue to a user.
export async function assignIssue(
  client: JiraClient,
  issueIdOrKey: string,
  accountID: string,
): Promise<void> {
  const body = { accountId: accountID };
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/assignee`, body);
}

// getIssueChangelog returns the changelog for an issue.
export async function getIssueChangelog(
  client: JiraClient,
  issueIdOrKey: string,
  startAt: number,
  maxResults: number,
): Promise<ChangelogPage> {
  const query: Record<string, string> = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/changelog`, query);
}

// getIssueComments returns the comments for an issue.
export async function getIssueComments(
  client: JiraClient,
  issueIdOrKey: string,
  startAt: number,
  maxResults: number,
): Promise<CommentPage> {
  const query: Record<string, string> = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/comment`, query);
}

// addIssueComment adds a comment to an issue.
export async function addIssueComment(
  client: JiraClient,
  issueIdOrKey: string,
  body: ADFNode,
  visibility?: Visibility | undefined,
): Promise<Comment> {
  const reqBody: { body: ADFNode; visibility?: Visibility } = { body };
  if (visibility != null) {
    reqBody.visibility = visibility;
  }
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/comment`, reqBody);
}

// getIssueComment retrieves a single comment on an issue.
export async function getIssueComment(
  client: JiraClient,
  issueIdOrKey: string,
  commentId: string,
): Promise<Comment> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`);
}

// updateIssueComment updates an existing comment on an issue.
export async function updateIssueComment(
  client: JiraClient,
  issueIdOrKey: string,
  commentId: string,
  body: ADFNode,
  visibility?: Visibility | undefined,
): Promise<Comment> {
  const reqBody: { body: ADFNode; visibility?: Visibility } = { body };
  if (visibility != null) {
    reqBody.visibility = visibility;
  }
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`, reqBody);
}

// deleteIssueComment deletes a comment from an issue.
export async function deleteIssueComment(
  client: JiraClient,
  issueIdOrKey: string,
  commentId: string,
): Promise<void> {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`);
}

// getIssueTransitions returns the available transitions for an issue.
export async function getIssueTransitions(client: JiraClient, issueIdOrKey: string): Promise<TransitionsResponse> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/transitions`);
}

// doIssueTransition performs a transition on an issue.
export async function doIssueTransition(
  client: JiraClient,
  issueIdOrKey: string,
  details: IssueUpdateDetails,
): Promise<void> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/transitions`, details);
}

// getIssueVotes returns the votes for an issue.
export async function getIssueVotes(client: JiraClient, issueIdOrKey: string): Promise<Votes> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/votes`);
}

// addIssueVote adds a vote to an issue for the authenticated user.
export async function addIssueVote(client: JiraClient, issueIdOrKey: string): Promise<void> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/votes`, null);
}

// removeIssueVote removes a vote from an issue for the authenticated user.
export async function removeIssueVote(client: JiraClient, issueIdOrKey: string): Promise<void> {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/votes`);
}

// getIssueWatchers returns the watchers for an issue.
export async function getIssueWatchers(client: JiraClient, issueIdOrKey: string): Promise<Watches> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`);
}

// addIssueWatcher adds a watcher to an issue.
// The Jira API expects the account ID as a quoted JSON string in the body.
export async function addIssueWatcher(
  client: JiraClient,
  issueIdOrKey: string,
  accountID: string,
): Promise<void> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`, JSON.stringify(accountID));
}

// removeIssueWatcher removes a watcher from an issue.
export async function removeIssueWatcher(
  client: JiraClient,
  issueIdOrKey: string,
  accountID: string,
): Promise<void> {
  const query: Record<string, string> = { accountId: accountID };
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`, query);
}

// getIssueWorklogs returns the worklogs for an issue.
export async function getIssueWorklogs(
  client: JiraClient,
  issueIdOrKey: string,
  startAt: number,
  maxResults: number,
): Promise<WorklogPage> {
  const query: Record<string, string> = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/worklog`, query);
}

// addIssueWorklog adds a worklog to an issue.
export async function addIssueWorklog(
  client: JiraClient,
  issueIdOrKey: string,
  worklog: Partial<Worklog>,
): Promise<Worklog> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/worklog`, worklog);
}

// getIssueWorklog retrieves a single worklog entry.
export async function getIssueWorklog(
  client: JiraClient,
  issueIdOrKey: string,
  worklogId: string,
): Promise<Worklog> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`);
}

// updateIssueWorklog updates a worklog entry.
export async function updateIssueWorklog(
  client: JiraClient,
  issueIdOrKey: string,
  worklogId: string,
  worklog: Partial<Worklog>,
): Promise<Worklog> {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`, worklog);
}

// deleteIssueWorklog deletes a worklog entry.
export async function deleteIssueWorklog(
  client: JiraClient,
  issueIdOrKey: string,
  worklogId: string,
): Promise<void> {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`);
}

// addIssueAttachment uploads an attachment to an issue.
export async function addIssueAttachment(
  client: JiraClient,
  issueIdOrKey: string,
  filePath: string,
): Promise<Attachment[]> {
  return await uploadFile(client, `/rest/api/3/issue/${issueIdOrKey}/attachments`, 'file', filePath);
}

// getIssueRemoteLinks returns the remote links for an issue.
export async function getIssueRemoteLinks(client: JiraClient, issueIdOrKey: string): Promise<RemoteLink[]> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink`);
}

// createIssueRemoteLink creates a remote link on an issue.
export async function createIssueRemoteLink(
  client: JiraClient,
  issueIdOrKey: string,
  link: RemoteLink,
): Promise<RemoteLink> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink`, link);
}

// getIssueRemoteLink retrieves a single remote link on an issue.
export async function getIssueRemoteLink(
  client: JiraClient,
  issueIdOrKey: string,
  linkId: string,
): Promise<RemoteLink> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`);
}

// updateIssueRemoteLink updates a remote link on an issue.
export async function updateIssueRemoteLink(
  client: JiraClient,
  issueIdOrKey: string,
  linkId: string,
  link: RemoteLink,
): Promise<RemoteLink> {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`, link);
}

// deleteIssueRemoteLink deletes a remote link on an issue.
export async function deleteIssueRemoteLink(
  client: JiraClient,
  issueIdOrKey: string,
  linkId: string,
): Promise<void> {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`);
}

// notifyIssue sends a notification for an issue.
export async function notifyIssue(
  client: JiraClient,
  issueIdOrKey: string,
  notify: IssueNotifyRequest,
): Promise<void> {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/notify`, notify);
}

// getIssueEditMeta returns the edit metadata for an issue.
export async function getIssueEditMeta(client: JiraClient, issueIdOrKey: string): Promise<IssueEditMeta> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/editmeta`);
}

export interface IssueEditMeta {
  fields?: Record<string, IssueEditMetaField>;
}

export interface IssueEditMetaField {
  required?: boolean;
  schema?: { type?: string; system?: string; custom?: string; customId?: number };
  name?: string;
  fieldId?: string;
  operations?: string[];
  allowedValues?: JsonValue[];
}

// getCreateMeta returns the create metadata for issues.
export async function getCreateMeta(
  client: JiraClient,
  projectKeys?: string[],
  expand?: string[],
): Promise<CreateMeta> {
  const query: Record<string, string> = {};
  if (projectKeys && projectKeys.length > 0) {
    query.projectKeys = projectKeys.join(',');
  }
  if (expand && expand.length > 0) {
    query.expand = expand.join(',');
  }
  return await get(client, '/rest/api/3/issue/createmeta', query);
}

// getIssueProperties returns all property keys for an issue.
export async function getIssueProperties(client: JiraClient, issueIdOrKey: string): Promise<EntityProperty[]> {
  const wrapper = await get<{ keys: EntityProperty[] }>(client, `/rest/api/3/issue/${issueIdOrKey}/properties`);
  return wrapper.keys;
}

// getIssueProperty retrieves a single property of an issue.
export async function getIssueProperty(
  client: JiraClient,
  issueIdOrKey: string,
  propertyKey: string,
): Promise<EntityProperty> {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`);
}

// setIssueProperty sets a property on an issue.
export async function setIssueProperty(
  client: JiraClient,
  issueIdOrKey: string,
  propertyKey: string,
  value: JsonBody,
): Promise<void> {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`, value);
}

// deleteIssueProperty deletes a property from an issue.
export async function deleteIssueProperty(
  client: JiraClient,
  issueIdOrKey: string,
  propertyKey: string,
): Promise<void> {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`);
}

// bulkFetchIssues fetches multiple issues by their IDs.
export async function bulkFetchIssues(
  client: JiraClient,
  issueIDs: string[],
  fields?: string[],
): Promise<{ issues: IssueDetailed[] }> {
  const body: Record<string, unknown> = { issueIds: issueIDs };
  if (fields && fields.length > 0) {
    body.fields = fields;
  }
  return await post(client, '/rest/api/3/issue/bulkfetch', body);
}

import { get, post, put, del, uploadFile } from './client.js';

// createIssue creates a new Jira issue.
export async function createIssue(client, details) {
  return await post(client, '/rest/api/3/issue', details);
}

// bulkCreateIssues creates multiple issues in a single request.
export async function bulkCreateIssues(client, req) {
  return await post(client, '/rest/api/3/issue/bulk', req);
}

// getIssue retrieves a single issue by ID or key.
export async function getIssue(client, issueIdOrKey, fields, expand) {
  const query = {};
  if (fields && fields.length > 0) {
    query.fields = fields.join(',');
  }
  if (expand && expand.length > 0) {
    query.expand = expand.join(',');
  }
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}`, query);
}

// editIssue edits an existing issue.
export async function editIssue(client, issueIdOrKey, details, notifyUsers) {
  const query = {};
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
export async function deleteIssue(client, issueIdOrKey, deleteSubtasks) {
  const query = {};
  if (deleteSubtasks) {
    query.deleteSubtasks = 'true';
  }
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}`, query);
}

// assignIssue assigns an issue to a user.
export async function assignIssue(client, issueIdOrKey, accountID) {
  const body = { accountId: accountID };
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/assignee`, body);
}

// getIssueChangelog returns the changelog for an issue.
export async function getIssueChangelog(client, issueIdOrKey, startAt, maxResults) {
  const query = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/changelog`, query);
}

// getIssueComments returns the comments for an issue.
export async function getIssueComments(client, issueIdOrKey, startAt, maxResults) {
  const query = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/comment`, query);
}

// addIssueComment adds a comment to an issue.
export async function addIssueComment(client, issueIdOrKey, body, visibility) {
  const reqBody = { body };
  if (visibility != null) {
    reqBody.visibility = visibility;
  }
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/comment`, reqBody);
}

// getIssueComment retrieves a single comment on an issue.
export async function getIssueComment(client, issueIdOrKey, commentId) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`, null);
}

// updateIssueComment updates an existing comment on an issue.
export async function updateIssueComment(client, issueIdOrKey, commentId, body, visibility) {
  const reqBody = { body };
  if (visibility != null) {
    reqBody.visibility = visibility;
  }
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`, reqBody);
}

// deleteIssueComment deletes a comment from an issue.
export async function deleteIssueComment(client, issueIdOrKey, commentId) {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/comment/${commentId}`, null);
}

// getIssueTransitions returns the available transitions for an issue.
export async function getIssueTransitions(client, issueIdOrKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/transitions`, null);
}

// doIssueTransition performs a transition on an issue.
export async function doIssueTransition(client, issueIdOrKey, details) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/transitions`, details);
}

// getIssueVotes returns the votes for an issue.
export async function getIssueVotes(client, issueIdOrKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/votes`, null);
}

// addIssueVote adds a vote to an issue for the authenticated user.
export async function addIssueVote(client, issueIdOrKey) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/votes`, null);
}

// removeIssueVote removes a vote from an issue for the authenticated user.
export async function removeIssueVote(client, issueIdOrKey) {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/votes`, null);
}

// getIssueWatchers returns the watchers for an issue.
export async function getIssueWatchers(client, issueIdOrKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`, null);
}

// addIssueWatcher adds a watcher to an issue.
// The Jira API expects the account ID as a quoted JSON string in the body.
export async function addIssueWatcher(client, issueIdOrKey, accountID) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`, JSON.stringify(accountID));
}

// removeIssueWatcher removes a watcher from an issue.
export async function removeIssueWatcher(client, issueIdOrKey, accountID) {
  const query = { accountId: accountID };
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/watchers`, query);
}

// getIssueWorklogs returns the worklogs for an issue.
export async function getIssueWorklogs(client, issueIdOrKey, startAt, maxResults) {
  const query = {
    startAt: String(startAt),
    maxResults: String(maxResults),
  };
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/worklog`, query);
}

// addIssueWorklog adds a worklog to an issue.
export async function addIssueWorklog(client, issueIdOrKey, worklog) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/worklog`, worklog);
}

// getIssueWorklog retrieves a single worklog entry.
export async function getIssueWorklog(client, issueIdOrKey, worklogId) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`, null);
}

// updateIssueWorklog updates a worklog entry.
export async function updateIssueWorklog(client, issueIdOrKey, worklogId, worklog) {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`, worklog);
}

// deleteIssueWorklog deletes a worklog entry.
export async function deleteIssueWorklog(client, issueIdOrKey, worklogId) {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/worklog/${worklogId}`, null);
}

// addIssueAttachment uploads an attachment to an issue.
export async function addIssueAttachment(client, issueIdOrKey, filePath) {
  return await uploadFile(client, `/rest/api/3/issue/${issueIdOrKey}/attachments`, 'file', filePath);
}

// getIssueRemoteLinks returns the remote links for an issue.
export async function getIssueRemoteLinks(client, issueIdOrKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink`, null);
}

// createIssueRemoteLink creates a remote link on an issue.
export async function createIssueRemoteLink(client, issueIdOrKey, link) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink`, link);
}

// getIssueRemoteLink retrieves a single remote link on an issue.
export async function getIssueRemoteLink(client, issueIdOrKey, linkId) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`, null);
}

// updateIssueRemoteLink updates a remote link on an issue.
export async function updateIssueRemoteLink(client, issueIdOrKey, linkId, link) {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`, link);
}

// deleteIssueRemoteLink deletes a remote link on an issue.
export async function deleteIssueRemoteLink(client, issueIdOrKey, linkId) {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/remotelink/${linkId}`, null);
}

// notifyIssue sends a notification for an issue.
export async function notifyIssue(client, issueIdOrKey, notify) {
  return await post(client, `/rest/api/3/issue/${issueIdOrKey}/notify`, notify);
}

// getIssueEditMeta returns the edit metadata for an issue.
export async function getIssueEditMeta(client, issueIdOrKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/editmeta`, null);
}

// getCreateMeta returns the create metadata for issues.
export async function getCreateMeta(client, projectKeys, expand) {
  const query = {};
  if (projectKeys && projectKeys.length > 0) {
    query.projectKeys = projectKeys.join(',');
  }
  if (expand && expand.length > 0) {
    query.expand = expand.join(',');
  }
  return await get(client, '/rest/api/3/issue/createmeta', query);
}

// getIssueProperties returns all property keys for an issue.
export async function getIssueProperties(client, issueIdOrKey) {
  const wrapper = await get(client, `/rest/api/3/issue/${issueIdOrKey}/properties`, null);
  return wrapper.keys;
}

// getIssueProperty retrieves a single property of an issue.
export async function getIssueProperty(client, issueIdOrKey, propertyKey) {
  return await get(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`, null);
}

// setIssueProperty sets a property on an issue.
export async function setIssueProperty(client, issueIdOrKey, propertyKey, value) {
  return await put(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`, value);
}

// deleteIssueProperty deletes a property from an issue.
export async function deleteIssueProperty(client, issueIdOrKey, propertyKey) {
  return await del(client, `/rest/api/3/issue/${issueIdOrKey}/properties/${propertyKey}`, null);
}

// bulkFetchIssues fetches multiple issues by their IDs.
export async function bulkFetchIssues(client, issueIDs, fields) {
  const body = { issueIds: issueIDs };
  if (fields && fields.length > 0) {
    body.fields = fields;
  }
  return await post(client, '/rest/api/3/issue/bulkfetch', body);
}

// Public Jira API surface for library consumers.
//
// Usage:
//   import { createClient, issues, projects } from 'acli/jira';
//   const client = createClient(profile);
//   const issue = await issues.getIssue(client, 'PROJ-123');

export { createClient } from './internal/jira/client.js';
export type { JiraClient, APIError } from './internal/jira/client.js';

export * as client from './internal/jira/client.js';
export * as issues from './internal/jira/issues.js';
export * as projects from './internal/jira/projects.js';
export * as agile from './internal/jira/agile.js';
export * as admin from './internal/jira/admin.js';
export * as schemes from './internal/jira/schemes.js';

export type * from './internal/jira/types.js';

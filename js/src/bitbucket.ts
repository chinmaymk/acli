// Public Bitbucket API surface for library consumers.
//
// Usage:
//   import { createClient, repos, pullrequests } from 'acli/bitbucket';
//   const client = createClient(profile);
//   const list = await repos.listRepositories(client, 'my-workspace');

export { createClient, get, post, put, deleteNoContent, getRaw, getAll, addPaginationParams } from './internal/bitbucket/client.js';
export type { BitbucketClient, PaginationOptions, PaginatedResponse } from './internal/bitbucket/client.js';

export * as branches from './internal/bitbucket/branches.js';
export * as branchrestrictions from './internal/bitbucket/branchrestrictions.js';
export * as commits from './internal/bitbucket/commits.js';
export * as deploykeys from './internal/bitbucket/deploykeys.js';
export * as deployments from './internal/bitbucket/deployments.js';
export * as downloads from './internal/bitbucket/downloads.js';
export * as environments from './internal/bitbucket/environments.js';
export * as issues from './internal/bitbucket/issues.js';
export * as pipelines from './internal/bitbucket/pipelines.js';
export * as projects from './internal/bitbucket/projects.js';
export * as pullrequests from './internal/bitbucket/pullrequests.js';
export * as repos from './internal/bitbucket/repos.js';
export * as search from './internal/bitbucket/search.js';
export * as snippets from './internal/bitbucket/snippets.js';
export * as user from './internal/bitbucket/user.js';
export * as webhooks from './internal/bitbucket/webhooks.js';
export * as workspaces from './internal/bitbucket/workspaces.js';

export type * from './internal/bitbucket/types.js';

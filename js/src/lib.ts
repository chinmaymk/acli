// Library entry point. Re-exports the typed API surface so consumers can
// use the same Atlassian clients that the CLI uses, e.g.:
//
//   import { jira, bitbucket, confluence, loadConfig } from 'acli';
//   const profile = loadConfig.getProfile(loadConfig.load(), 'work');
//   const client = jira.createClient(profile);
//   const issues = await jira.issues.searchJQL(client, { jql: 'project = FOO' });
//
// Or, more granularly:
//
//   import { createClient } from 'acli/jira';
//   import { listRepositories } from 'acli/bitbucket';
//
// The CLI binary lives in `index.ts`; this file is purely a programmatic
// surface and contains no side effects.

// --- Shared primitives ---
export type {
  JsonValue,
  JsonObject,
  JsonArray,
  JsonPrimitive,
  JsonBody,
  ADFNode,
  ADFMark,
  ADFDocument,
  BitbucketPaginatedResponse,
} from './internal/types.js';

// --- Config ---
export * as config from './internal/config/config.js';
export type { Profile, Config, Defaults } from './internal/config/config.js';

// --- ADF rendering ---
export { render as renderADF } from './internal/adf/render.js';

// --- Jira ---
export * as jira from './jira.js';
export type * as jiraTypes from './internal/jira/types.js';

// --- Bitbucket ---
export * as bitbucket from './bitbucket.js';
export type * as bitbucketTypes from './internal/bitbucket/types.js';

// --- Confluence ---
export * as confluence from './confluence.js';

import { load, getProfile as getProfileFromConfig } from '../internal/config/config.js';
import { createClient as createJiraClient } from '../internal/jira/client.js';
import { createClient as createBBClient } from '../internal/bitbucket/client.js';
import { createClient as createConfClient } from '../internal/api/client.js';
import type { Profile } from '../internal/config/config.js';
import type { ConfluenceClient } from '../internal/api/client.js';
import type { JsonBody } from '../internal/types.js';

/**
 * Value type for arbitrary yargs argv entries. Yargs only ever produces
 * primitives, arrays of primitives, or undefined for parsed options, so
 * we model that explicitly instead of leaning on `any`/`unknown`.
 */
export type ArgvValue =
  | string
  | number
  | boolean
  | (string | number | boolean)[]
  | undefined;

// Use a minimal interface for the common flags. We intentionally omit a
// string index signature so narrowly-typed handler argvs (e.g. JiraArgv,
// BitbucketArgv) can be passed here without needing to include one
// themselves. The per-command argv types from the CLI files carry the
// additional option/positional fields that each handler needs.
export interface BaseArgv {
  profile?: string;
  output?: string;
  json?: boolean;
  project?: string;
  page?: number;
  pagelen?: number;
  all?: boolean;
}

// Inferred return types from the client factory functions
export type JiraClient = ReturnType<typeof createJiraClient>;
export type BitbucketClient = ReturnType<typeof createBBClient>;

/**
 * Reads --profile flag from argv, loads config, resolves profile.
 */
export function getProfile(argv: BaseArgv): Profile {
  const cfg = load();
  return getProfileFromConfig(cfg, argv.profile || undefined);
}

/**
 * Gets the active profile and creates a Jira API client.
 */
export function getJiraClient(argv: BaseArgv): JiraClient {
  const profile = getProfile(argv);
  return createJiraClient(profile);
}

/**
 * Gets the active profile and creates a Bitbucket API client.
 */
export function getBitbucketClient(argv: BaseArgv): BitbucketClient {
  const profile = getProfile(argv);
  return createBBClient(profile);
}

/**
 * Gets the active profile and creates a Confluence API client.
 */
export function getConfluenceClient(argv: BaseArgv): ConfluenceClient {
  const profile = getProfile(argv);
  return createConfClient(profile.atlassian_url, profile.email, profile.api_token);
}

/**
 * Returns --project flag value if set, otherwise falls back to profile default.
 */
export function defaultProject(argv: BaseArgv): string {
  if (argv.project) {
    return argv.project;
  }
  const profile = getProfile(argv);
  return profile.defaults?.project || '';
}

/**
 * Returns the positional arg at argIndex if present, otherwise falls back to profile workspace default.
 */
export function defaultWorkspace(argv: BaseArgv, args: string[], argIndex: number): string {
  if (argIndex < args.length) {
    return args[argIndex];
  }
  const profile = getProfile(argv);
  if (profile.defaults?.workspace) {
    return profile.defaults.workspace;
  }
  throw new Error("workspace is required: provide it as an argument or set a default with 'acli config set-defaults'");
}

/**
 * Returns --project flag value if set, otherwise falls back to profile BB project default.
 */
export function defaultBBProject(argv: BaseArgv): string {
  if (argv.project) {
    return argv.project;
  }
  const profile = getProfile(argv);
  return profile.defaults?.bb_project || '';
}

/**
 * Handles the common [workspace] <repo> args pattern.
 * With 2+ args: workspace=args[0], repo=args[1].
 * With 1 arg: workspace from profile default, repo=args[0].
 */
export function resolveWorkspaceAndRepo(argv: BaseArgv, args: string[]): [string, string] {
  if (args.length >= 2) {
    return [args[0], args[1]];
  }
  const workspace = defaultWorkspace(argv, [], 0);
  return [workspace, args[0]];
}

/**
 * Handles the [workspace] <repo> <id> args pattern.
 * With 3+ args: workspace=args[0], repo=args[1], id=args[2].
 * With 2 args: workspace from profile default, repo=args[0], id=args[1].
 */
export function resolveWorkspaceRepoAndID(argv: BaseArgv, args: string[]): [string, string, string] {
  if (args.length >= 3) {
    return [args[0], args[1], args[2]];
  }
  const workspace = defaultWorkspace(argv, [], 0);
  return [workspace, args[0], args[1]];
}

/**
 * Returns the first line of a string.
 */
export function firstLine(s: string): string {
  const idx = s.search(/[\n\r]/);
  if (idx === -1) return s;
  return s.slice(0, idx);
}

/**
 * Truncates a string to maxLen characters, appending "..." if truncated.
 */
export function truncate(s: string, maxLen: number): string {
  if (s.length <= maxLen) return s;
  if (maxLen <= 3) return s.slice(0, maxLen);
  return s.slice(0, maxLen - 3) + '...';
}

// Minimal shape of a Jira issue used for display helpers
interface JiraIssueFields {
  issuetype?: { name?: string };
  IssueType?: { name?: string };
  status?: { name?: string };
  Status?: { name?: string };
  priority?: { name?: string };
  Priority?: { name?: string };
  assignee?: { displayName?: string };
  Assignee?: { displayName?: string };
  summary?: string;
  Summary?: string;
}

interface JiraIssue {
  key?: string;
  fields?: JiraIssueFields;
}

/**
 * Prints a formatted issue row using tab-separated columns.
 * Columns: key, type, status, priority, assignee, summary
 */
export function printIssueRow(issue: JiraIssue): void {
  const key = issue.key || '';
  const issueType = issue.fields?.issuetype?.name || issue.fields?.IssueType?.name || '';
  const status = issue.fields?.status?.name || issue.fields?.Status?.name || '';
  const priority = issue.fields?.priority?.name || issue.fields?.Priority?.name || '';
  const assignee = issue.fields?.assignee?.displayName || issue.fields?.Assignee?.displayName || '';
  const summary = issue.fields?.summary || issue.fields?.Summary || '';

  const cols = [key, issueType, status, priority, assignee, summary];
  console.log(cols.join('\t'));
}

/**
 * Formats issue rows as an aligned table using padEnd for column widths.
 */
export function printIssueTable(issues: JiraIssue[]): void {
  if (issues.length === 0) return;

  const rows = issues.map((issue) => [
    issue.key || '',
    issue.fields?.issuetype?.name || issue.fields?.IssueType?.name || '',
    issue.fields?.status?.name || issue.fields?.Status?.name || '',
    issue.fields?.priority?.name || issue.fields?.Priority?.name || '',
    issue.fields?.assignee?.displayName || issue.fields?.Assignee?.displayName || '',
    issue.fields?.summary || issue.fields?.Summary || '',
  ]);

  // Compute column widths (all except last)
  const colCount = rows[0].length - 1;
  const widths: number[] = Array(colCount).fill(0);
  for (const row of rows) {
    for (let i = 0; i < colCount; i++) {
      if (row[i].length > widths[i]) widths[i] = row[i].length;
    }
  }

  for (const row of rows) {
    const padded = row.slice(0, colCount).map((cell, i) => cell.padEnd(widths[i] + 2));
    console.log(padded.join('') + row[colCount]);
  }
}

/**
 * Returns true if argv.output === 'json'.
 */
export function isJSONOutput(argv: BaseArgv): boolean {
  if (argv.output === 'json') return true;
  // Backward compat: check per-command --json flag if it exists
  if (argv.json === true) return true;
  return false;
}

/**
 * Pretty-prints v as indented JSON to stdout.
 */
export function outputJSON(v: JsonBody): void {
  console.log(JSON.stringify(v, null, 2));
}

/**
 * Outputs a structured result for mutation operations.
 * In text mode, prints the message. In JSON mode, outputs a structured envelope.
 */
export function outputResult(argv: BaseArgv, action: string, key: string, message: string, data: JsonBody): void {
  if (isJSONOutput(argv)) {
    outputJSON({ status: 'ok', action, key, message, data });
    return;
  }
  console.log(message);
}

/**
 * Prints a hint about fetching more results if there are more pages.
 */
export function printPaginationHint(shown: number, total: number): void {
  if (total <= 0 || shown >= total) {
    console.log(`\nShowing ${shown} results`);
    return;
  }
  console.log(`\nShowing ${shown} of ${total} results (use --all to fetch all)`);
}

/**
 * Reads page/pagelen/all from argv for Bitbucket pagination.
 */
export function getBBPaginationOpts(argv: BaseArgv): { page: number; pageLen: number; all: boolean } {
  return {
    page: argv.page || 0,
    pageLen: argv.pagelen || 0,
    all: argv.all || false,
  };
}

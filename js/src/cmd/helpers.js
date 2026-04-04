import { load, getProfile as getProfileFromConfig } from '../internal/config/config.js';
import { createClient as createJiraClient } from '../internal/jira/client.js';
import { createClient as createBBClient } from '../internal/bitbucket/client.js';
import { createClient as createConfClient } from '../internal/api/client.js';

/**
 * Reads --profile flag from argv, loads config, resolves profile.
 * @param {object} argv - yargs parsed argv
 * @returns {import('../internal/config/config.js').Profile}
 */
export function getProfile(argv) {
  const cfg = load();
  return getProfileFromConfig(cfg, argv.profile || null);
}

/**
 * Gets the active profile and creates a Jira API client.
 * @param {object} argv
 * @returns {object} Jira client
 */
export function getJiraClient(argv) {
  const profile = getProfile(argv);
  return createJiraClient(profile);
}

/**
 * Gets the active profile and creates a Bitbucket API client.
 * @param {object} argv
 * @returns {object} Bitbucket client
 */
export function getBitbucketClient(argv) {
  const profile = getProfile(argv);
  return createBBClient(profile);
}

/**
 * Gets the active profile and creates a Confluence API client.
 * @param {object} argv
 * @returns {object} Confluence client
 */
export function getConfluenceClient(argv) {
  const profile = getProfile(argv);
  return createConfClient(profile.atlassian_url, profile.email, profile.api_token);
}

/**
 * Returns --project flag value if set, otherwise falls back to profile default.
 * @param {object} argv
 * @returns {string}
 */
export function defaultProject(argv) {
  if (argv.project) {
    return argv.project;
  }
  const profile = getProfile(argv);
  return profile.defaults?.project || '';
}

/**
 * Returns the positional arg at argIndex if present, otherwise falls back to profile workspace default.
 * @param {object} argv
 * @param {string[]} args - positional arguments array
 * @param {number} argIndex - index of the workspace arg
 * @returns {string}
 */
export function defaultWorkspace(argv, args, argIndex) {
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
 * @param {object} argv
 * @returns {string}
 */
export function defaultBBProject(argv) {
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
 * @param {object} argv
 * @param {string[]} args
 * @returns {[string, string]} [workspace, repo]
 */
export function resolveWorkspaceAndRepo(argv, args) {
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
 * @param {object} argv
 * @param {string[]} args
 * @returns {[string, string, string]} [workspace, repo, id]
 */
export function resolveWorkspaceRepoAndID(argv, args) {
  if (args.length >= 3) {
    return [args[0], args[1], args[2]];
  }
  const workspace = defaultWorkspace(argv, [], 0);
  return [workspace, args[0], args[1]];
}

/**
 * Returns the first line of a string.
 * @param {string} s
 * @returns {string}
 */
export function firstLine(s) {
  const idx = s.search(/[\n\r]/);
  if (idx === -1) return s;
  return s.slice(0, idx);
}

/**
 * Truncates a string to maxLen characters, appending "..." if truncated.
 * @param {string} s
 * @param {number} maxLen
 * @returns {string}
 */
export function truncate(s, maxLen) {
  if (s.length <= maxLen) return s;
  if (maxLen <= 3) return s.slice(0, maxLen);
  return s.slice(0, maxLen - 3) + '...';
}

/**
 * Prints a formatted issue row using tab-separated columns.
 * Columns: key, type, status, priority, assignee, summary
 * @param {object} issue - Jira IssueDetailed object
 */
export function printIssueRow(issue) {
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
 * @param {object[]} issues - array of Jira IssueDetailed objects
 */
export function printIssueTable(issues) {
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
  const widths = Array(colCount).fill(0);
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
 * @param {object} argv
 * @returns {boolean}
 */
export function isJSONOutput(argv) {
  if (argv.output === 'json') return true;
  // Backward compat: check per-command --json flag if it exists
  if (argv.json === true) return true;
  return false;
}

/**
 * Pretty-prints v as indented JSON to stdout.
 * @param {*} v
 */
export function outputJSON(v) {
  console.log(JSON.stringify(v, null, 2));
}

/**
 * Outputs a structured result for mutation operations.
 * In text mode, prints the message. In JSON mode, outputs a structured envelope.
 * @param {object} argv
 * @param {string} action
 * @param {string} key
 * @param {string} message
 * @param {*} data
 */
export function outputResult(argv, action, key, message, data) {
  if (isJSONOutput(argv)) {
    outputJSON({ status: 'ok', action, key, message, data });
    return;
  }
  console.log(message);
}

/**
 * Prints a hint about fetching more results if there are more pages.
 * @param {number} shown
 * @param {number} total
 */
export function printPaginationHint(shown, total) {
  if (total <= 0 || shown >= total) {
    console.log(`\nShowing ${shown} results`);
    return;
  }
  console.log(`\nShowing ${shown} of ${total} results (use --all to fetch all)`);
}

/**
 * Reads page/pagelen/all from argv for Bitbucket pagination.
 * @param {object} argv
 * @returns {{ page: number, pageLen: number, all: boolean }}
 */
export function getBBPaginationOpts(argv) {
  return {
    page: argv.page || 0,
    pageLen: argv.pagelen || 0,
    all: argv.all || false,
  };
}

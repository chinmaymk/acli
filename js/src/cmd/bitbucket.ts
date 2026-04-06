import {
  getBitbucketClient,
  defaultWorkspace,
  defaultBBProject,
  resolveWorkspaceAndRepo,
  resolveWorkspaceRepoAndID,
  isJSONOutput,
  outputJSON,
  outputResult,
  getBBPaginationOpts,
  firstLine,
  truncate,
} from './helpers.js';
import * as bbRepos from '../internal/bitbucket/repos.js';
import * as bbPRs from '../internal/bitbucket/pullrequests.js';
import * as bbPipelines from '../internal/bitbucket/pipelines.js';
import * as bbBranches from '../internal/bitbucket/branches.js';
import * as bbCommits from '../internal/bitbucket/commits.js';
import * as bbWorkspaces from '../internal/bitbucket/workspaces.js';
import * as bbProjects from '../internal/bitbucket/projects.js';
import * as bbWebhooks from '../internal/bitbucket/webhooks.js';
import * as bbEnvironments from '../internal/bitbucket/environments.js';
import * as bbDeployKeys from '../internal/bitbucket/deploykeys.js';
import * as bbDownloads from '../internal/bitbucket/downloads.js';
import * as bbSnippets from '../internal/bitbucket/snippets.js';
import * as bbIssues from '../internal/bitbucket/issues.js';
import * as bbSearch from '../internal/bitbucket/search.js';
import * as bbDeployments from '../internal/bitbucket/deployments.js';
import * as bbBranchRestrictions from '../internal/bitbucket/branchrestrictions.js';
import * as bbUser from '../internal/bitbucket/user.js';
import type {
  CreateRepoRequest,
  ForkRepoRequest,
  UpdatePRRequest,
  MergePRRequest,
  CreatePRTaskRequest,
  UpdatePRTaskRequest,
  CreateIssueRequest,
  CreateWebhookRequest,
  InlineCommentParams,
} from '../internal/bitbucket/types.js';
import type { Argv } from 'yargs';

/**
 * Pagination options as yargs option definitions.
 */
const paginationOptions = {
  page: { type: 'number' as const, description: 'Page number', default: 0 },
  pagelen: { type: 'number' as const, description: 'Page length', default: 0 },
  all: { type: 'boolean' as const, description: 'Fetch all pages', default: false },
};

/**
 * Formats an array of rows as an aligned table using padEnd.
 */
type Cell = string | number | boolean | null | undefined;
function printTable(headers: string[], rows: Cell[][]): void {
  const allRows: Cell[][] = [headers, ...rows];
  const colCount = headers.length;
  const widths = Array(colCount).fill(0);
  for (const row of allRows) {
    for (let i = 0; i < colCount; i++) {
      const cell = String(row[i] ?? '');
      if (cell.length > widths[i]) widths[i] = cell.length;
    }
  }
  for (const row of allRows) {
    const parts = row.map((cell, i) => {
      const s = String(cell ?? '');
      // Last column not padded
      return i < colCount - 1 ? s.padEnd(widths[i] + 2) : s;
    });
    console.log(parts.join(''));
  }
}

/**
 * Resolves workspace/repo/prID/taskID from positional args.
 * With 4 args: workspace=args[0], repo=args[1], prID=args[2], taskID=args[3].
 * With 3 args: workspace from profile default, repo=args[0], prID=args[1], taskID=args[2].
 */
function resolveWorkspaceRepoIDAndTaskID(argv: any, args: string[]): [string, string, string, string] {
  if (args.length >= 4) {
    return [args[0], args[1], args[2], args[3]];
  }
  const workspace = defaultWorkspace(argv, [], 0);
  return [workspace, args[0], args[1], args[2]];
}

/**
 * Registers all Bitbucket subcommands on yargs.
 */
export function registerBitbucketCommands(yargs: Argv): Argv {
  return yargs
        // ─── whoami ──────────────────────────────────────────────────────
        .command(
          'whoami',
          'Show the current authenticated Bitbucket user',
          () => {},
          async (argv: any) => {
            const client = getBitbucketClient(argv);
            const user = await bbUser.getCurrentUser(client);
            if (isJSONOutput(argv)) { outputJSON(user); return; }
            console.log(`Display name: ${user.display_name}`);
            console.log(`Username:     ${user.nickname}`);
            console.log(`UUID:         ${user.uuid}`);
            console.log(`Account ID:   ${user.account_id}`);
          },
        )

        // ─── repo ─────────────────────────────────────────────────────────
        .command(
          ['repo', 'r'],
          'Manage repositories',
          (repoCmd: Argv) => {
            repoCmd
              .command(
                ['list [workspace]', 'ls [workspace]'],
                'List repositories in a workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string', description: 'Workspace slug' })
                  .option('role', { type: 'string', description: 'Filter by role (admin, contributor, member, owner)' })
                  .option('query', { type: 'string', description: 'Filter with query (Bitbucket query syntax)' })
                  .option('sort', { type: 'string', description: 'Sort field (e.g. -updated_on)' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const repos = await bbRepos.listRepositories(client, workspace, {
                    role: argv.role,
                    q: argv.query,
                    sort: argv.sort,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(repos); return; }
                  printTable(
                    ['NAME', 'SLUG', 'LANGUAGE', 'PRIVATE', 'UPDATED'],
                    repos.map((r) => [r.full_name, r.slug, r.language, r.is_private, r.updated_on]),
                  );
                },
              )
              .command(
                'get [workspace] <repo>',
                'Get repository details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string', description: 'Repository slug' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const repo = await bbRepos.getRepository(client, workspace, repoSlug);
                  if (isJSONOutput(argv)) { outputJSON(repo); return; }
                  const mainBranch = repo.mainbranch?.name ?? 'N/A';
                  console.log(`Name:         ${repo.full_name}`);
                  console.log(`Description:  ${repo.description}`);
                  console.log(`Language:     ${repo.language}`);
                  console.log(`SCM:          ${repo.scm}`);
                  console.log(`Private:      ${repo.is_private}`);
                  console.log(`Main Branch:  ${mainBranch}`);
                  console.log(`Created:      ${repo.created_on}`);
                  console.log(`Updated:      ${repo.updated_on}`);
                  console.log(`URL:          ${repo.links?.html?.href}`);
                },
              )
              .command(
                'create [workspace]',
                'Create a new repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .option('name', { type: 'string', description: 'Repository name (required)' })
                  .option('slug', { type: 'string', description: 'Repository slug for the URL (defaults to name)' })
                  .option('description', { type: 'string', description: 'Repository description' })
                  .option('language', { type: 'string', description: 'Programming language' })
                  .option('private', { type: 'boolean', description: 'Make repository private', default: true })
                  .option('project', { type: 'string', description: 'Project key to assign to' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  if (!argv.name) throw new Error('--name is required');
                  const client = getBitbucketClient(argv);
                  const projectKey = defaultBBProject(argv);
                  const req: CreateRepoRequest = {
                    scm: 'git',
                    name: argv.name,
                    is_private: argv.private,
                    description: argv.description,
                    language: argv.language,
                  };
                  if (projectKey) req.project = { key: projectKey };
                  const repo = await bbRepos.createRepository(client, workspace, req);
                  outputResult(argv, 'created', repo.full_name, `Created repository: ${repo.full_name}`, repo);
                },
              )
              .command(
                'delete [workspace] <repo>',
                'Delete a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbRepos.deleteRepository(client, workspace, repoSlug);
                  outputResult(argv, 'deleted', `${workspace}/${repoSlug}`, `Deleted repository: ${workspace}/${repoSlug}`, null);
                },
              )
              .command(
                'fork [workspace] <repo>',
                'Fork a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('name', { type: 'string', description: 'Name for the forked repo' })
                  .option('workspace', { type: 'string', description: 'Target workspace for the fork' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const req: ForkRepoRequest = { name: argv.name };
                  if (argv['target-workspace']) req.workspace = { slug: argv['target-workspace'] };
                  const repo = await bbRepos.forkRepository(client, workspace, repoSlug, req);
                  outputResult(argv, 'forked', repo.full_name, `Forked repository: ${repo.full_name}`, repo);
                },
              )
              .command(
                'forks [workspace] <repo>',
                'List forks of a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const forks = await bbRepos.listForks(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(forks); return; }
                  printTable(
                    ['NAME', 'SLUG', 'OWNER', 'PRIVATE'],
                    forks.map((r) => [r.full_name, r.slug, r.owner?.display_name, r.is_private]),
                  );
                },
              );
          },
        )

        // ─── pr ───────────────────────────────────────────────────────────
        .command(
          ['pr', 'pull-request'],
          'Manage pull requests',
          (prCmd: Argv) => {
            prCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List pull requests',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('state', { type: 'string', description: 'Filter by state (OPEN, MERGED, DECLINED, SUPERSEDED)', default: 'OPEN' })
                  .option('author', { type: 'string', description: "Filter by author Bitbucket username or UUID" })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const prs = await bbPRs.listPullRequests(client, workspace, repoSlug, {
                    state: argv.state,
                    author: argv.author,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(prs); return; }
                  printTable(
                    ['ID', 'TITLE', 'STATE', 'AUTHOR', 'SOURCE', 'DESTINATION'],
                    prs.map((pr) => [pr.id, pr.title, pr.state, pr.author?.display_name, pr.source?.branch?.name, pr.destination?.branch?.name]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <pr-id>',
                'Get pull request details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const pr = await bbPRs.getPullRequest(client, workspace, repoSlug, prID);
                  if (isJSONOutput(argv)) { outputJSON(pr); return; }
                  console.log(`ID:           ${pr.id}`);
                  console.log(`Title:        ${pr.title}`);
                  console.log(`State:        ${pr.state}`);
                  console.log(`Author:       ${pr.author?.display_name}`);
                  console.log(`Source:       ${pr.source?.branch?.name}`);
                  console.log(`Destination:  ${pr.destination?.branch?.name}`);
                  console.log(`Comments:     ${pr.comment_count}`);
                  console.log(`Tasks:        ${pr.task_count}`);
                  console.log(`Created:      ${pr.created_on}`);
                  console.log(`Updated:      ${pr.updated_on}`);
                  console.log(`URL:          ${pr.links?.html?.href}`);
                  if (pr.description) console.log(`\nDescription:\n${pr.description}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('title', { type: 'string', description: 'Pull request title (required)' })
                  .option('source', { type: 'string', description: 'Source branch name (required)' })
                  .option('destination', { type: 'string', description: 'Destination branch name (defaults to main branch)' })
                  .option('description', { type: 'string', description: 'Pull request description' })
                  .option('close-source-branch', { type: 'boolean', description: 'Close source branch after merge', default: false }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.title || !argv.source) throw new Error('--title and --source are required');
                  const client = getBitbucketClient(argv);
                  const pr = await bbPRs.createPullRequest(client, workspace, repoSlug, {
                    title: argv.title,
                    description: argv.description,
                    source_branch: argv.source,
                    destination_branch: argv.destination,
                    close_source_branch: argv['close-source-branch'],
                  });
                  outputResult(argv, 'created', String(pr.id), `Created PR #${pr.id}: ${pr.title}`, pr);
                },
              )
              .command(
                'update [workspace] <repo> <pr-id>',
                'Update a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .option('title', { type: 'string', description: 'New title' })
                  .option('description', { type: 'string', description: 'New description' })
                  .option('destination', { type: 'string', description: 'New destination branch' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const req: UpdatePRRequest & { destination?: { branch: { name: string } } } = {};
                  if (argv.title) req.title = argv.title;
                  if (argv.description !== undefined) req.description = argv.description;
                  if (argv.destination) req.destination = { branch: { name: argv.destination } };
                  const pr = await bbPRs.updatePullRequest(client, workspace, repoSlug, prID, req);
                  outputResult(argv, 'updated', String(pr.id), `Updated PR #${pr.id}: ${pr.title}`, pr);
                },
              )
              .command(
                'approve [workspace] <repo> <pr-id>',
                'Approve a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbPRs.approvePullRequest(client, workspace, repoSlug, prID);
                  outputResult(argv, 'approved', String(prID), `Approved PR #${prID}`, null);
                },
              )
              .command(
                'unapprove [workspace] <repo> <pr-id>',
                'Remove approval from a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbPRs.unapprovePullRequest(client, workspace, repoSlug, prID);
                  outputResult(argv, 'unapproved', String(prID), `Removed approval from PR #${prID}`, null);
                },
              )
              .command(
                'decline [workspace] <repo> <pr-id>',
                'Decline a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbPRs.declinePullRequest(client, workspace, repoSlug, prID);
                  outputResult(argv, 'declined', String(prID), `Declined PR #${prID}`, null);
                },
              )
              .command(
                'merge [workspace] <repo> <pr-id>',
                'Merge a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .option('strategy', { type: 'string', description: 'Merge strategy (merge_commit, squash, fast_forward)' })
                  .option('message', { type: 'string', description: 'Merge commit message' })
                  .option('close-source-branch', { type: 'boolean', description: 'Close source branch after merge' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const req: MergePRRequest = {};
                  if (argv.strategy) req.merge_strategy = argv.strategy;
                  if (argv.message) req.message = argv.message;
                  if (argv['close-source-branch'] !== undefined) req.close_source_branch = argv['close-source-branch'];
                  const pr = await bbPRs.mergePullRequest(client, workspace, repoSlug, prID, req);
                  outputResult(argv, 'merged', String(pr.id), `Merged PR #${pr.id}: ${pr.title}`, pr);
                },
              )
              .command(
                'request-changes [workspace] <repo> <pr-id>',
                'Request changes on a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbPRs.requestChangesPullRequest(client, workspace, repoSlug, prID);
                  outputResult(argv, 'changes_requested', String(prID), `Requested changes on PR #${prID}`, null);
                },
              )
              .command(
                'diff [workspace] <repo> <pr-id>',
                'Get the diff of a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const diff = await bbPRs.getPRDiff(client, workspace, repoSlug, prID);
                  process.stdout.write(diff);
                },
              )
              .command(
                'comments [workspace] <repo> <pr-id>',
                'List comments on a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const comments = await bbPRs.listPRComments(client, workspace, repoSlug, prID, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(comments); return; }
                  for (const c of comments) {
                    console.log(`#${c.id} by ${c.user?.display_name} (${c.created_on})`);
                    if (c.inline) console.log(`  File: ${c.inline.path}`);
                    console.log(`  ${c.content?.raw}\n`);
                  }
                },
              )
              .command(
                'comment [workspace] <repo> <pr-id>',
                'Add a comment to a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .option('body', { type: 'string', description: 'Comment body (required)' })
                  .option('file', { type: 'string', description: 'File path for an inline comment' })
                  .option('line', { type: 'number', description: 'Line number in the new version of the file (requires --file)', default: 0 }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  if (!argv.body) throw new Error('--body is required');
                  if (argv.file && !argv.line) throw new Error('--line is required when --file is specified');
                  if (argv.line && !argv.file) throw new Error('--file is required when --line is specified');
                  const client = getBitbucketClient(argv);
                  let inline: InlineCommentParams | undefined;
                  if (argv.file) inline = { path: argv.file, to: argv.line };
                  const comment = await bbPRs.createPRComment(client, workspace, repoSlug, prID, argv.body, inline);
                  outputResult(argv, 'created', String(comment.id), `Added comment #${comment.id} to PR #${prID}`, comment);
                },
              )
              .command(
                ['commits [workspace] <repo> <pr-id>', 'log [workspace] <repo> <pr-id>'],
                'List commits on a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const commits = await bbPRs.listPRCommits(client, workspace, repoSlug, prID, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(commits); return; }
                  printTable(
                    ['HASH', 'MESSAGE', 'AUTHOR', 'DATE'],
                    commits.map((c) => {
                      const hash = c.hash?.slice(0, 12) ?? '';
                      return [hash, truncate(firstLine(c.message ?? ''), 60), c.author?.raw, c.date];
                    }),
                  );
                },
              )
              .command(
                'diffstat [workspace] <repo> <pr-id>',
                'Get file-level diff statistics for a pull request',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const stats = await bbPRs.getPRDiffStat(client, workspace, repoSlug, prID, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(stats); return; }
                  printTable(
                    ['STATUS', 'FILE', 'ADDED', 'REMOVED'],
                    stats.map((s) => {
                      let filePath = '';
                      if (s.status === 'renamed' && s.old && s.new) {
                        filePath = `${s.old.path} → ${s.new.path}`;
                      } else if (s.new) {
                        filePath = s.new.path;
                      } else if (s.old) {
                        filePath = s.old.path;
                      }
                      return [s.status, filePath, `+${s.lines_added}`, `-${s.lines_removed}`];
                    }),
                  );
                },
              )
              .command(
                'activity [workspace] <repo> <pr-id>',
                'List activity on a pull request (comments, approvals, updates)',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pr-id', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const prID = parseInt(idStr, 10);
                  if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const activities = await bbPRs.listPRActivity(client, workspace, repoSlug, prID, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(activities); return; }
                  for (const a of activities as Array<Record<string, any>>) {
                    if (a.approval) {
                      console.log(`[APPROVAL] ${a.approval.user?.display_name} approved on ${a.approval.date}`);
                    } else if (a.update) {
                      console.log(`[UPDATE] ${a.update.author?.display_name} changed state to ${a.update.state} on ${a.update.date}`);
                      if (a.update.reason) console.log(`  Reason: ${a.update.reason}`);
                    } else if (a.comment) {
                      console.log(`[COMMENT] #${a.comment.id} by ${a.comment.user?.display_name} (${a.comment.created_on})`);
                      if (a.comment.inline) console.log(`  File: ${a.comment.inline.path}`);
                      console.log(`  ${truncate(a.comment.content?.raw ?? '', 120)}`);
                    } else {
                      console.log('[UNKNOWN] unrecognized activity type');
                    }
                    console.log();
                  }
                },
              )
              // ─── pr task ────────────────────────────────────────────────
              .command(
                ['task', 'tasks'],
                'Manage tasks on a pull request',
                (taskCmd: Argv) => {
                  taskCmd
                    .command(
                      ['list [workspace] <repo> <pr-id>', 'ls [workspace] <repo> <pr-id>'],
                      'List tasks on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .options(paginationOptions),
                      async (argv: any) => {
                        const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                        const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                        const prID = parseInt(idStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                        const client = getBitbucketClient(argv);
                        const tasks = await bbPRs.listPRTasks(client, workspace, repoSlug, prID, getBBPaginationOpts(argv));
                        if (isJSONOutput(argv)) { outputJSON(tasks); return; }
                        printTable(
                          ['ID', 'STATE', 'CONTENT', 'CREATOR', 'CREATED'],
                          tasks.map((t) => [t.id, t.state, truncate(t.content?.raw ?? '', 60), t.creator?.display_name, t.created_on]),
                        );
                      },
                    )
                    .command(
                      'get [workspace] <repo> <pr-id> <task-id>',
                      'Get a task on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .positional('task-id', { type: 'string' }),
                      async (argv: any) => {
                        const rawArgs = argv.workspace
                          ? [argv.workspace, argv.repo, argv['pr-id'], argv['task-id']]
                          : [argv.repo, argv['pr-id'], argv['task-id']];
                        const [workspace, repoSlug, prIDStr, taskIDStr] = resolveWorkspaceRepoIDAndTaskID(argv, rawArgs);
                        const prID = parseInt(prIDStr, 10);
                        const taskID = parseInt(taskIDStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${prIDStr}`);
                        if (isNaN(taskID)) throw new Error(`invalid task ID: ${taskIDStr}`);
                        const client = getBitbucketClient(argv);
                        const task = await bbPRs.getPRTask(client, workspace, repoSlug, prID, taskID);
                        if (isJSONOutput(argv)) { outputJSON(task); return; }
                        console.log(`ID:        ${task.id}`);
                        console.log(`State:     ${task.state}`);
                        console.log(`Content:   ${task.content?.raw}`);
                        console.log(`Creator:   ${task.creator?.display_name}`);
                        console.log(`Created:   ${task.created_on}`);
                        console.log(`Updated:   ${task.updated_on}`);
                        if (task.resolved_on) console.log(`Resolved:  ${task.resolved_on}`);
                        if (task.comment) console.log(`Comment:   #${task.comment.id}`);
                      },
                    )
                    .command(
                      'create [workspace] <repo> <pr-id>',
                      'Create a task on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .option('body', { type: 'string', description: 'Task content (required)' })
                        .option('comment-id', { type: 'number', description: 'Associate task with a comment' })
                        .option('file', { type: 'string', description: 'File path to attach the task to a specific line (creates an inline comment)' })
                        .option('line', { type: 'number', description: 'Line number in the new version of the file (requires --file)', default: 0 }),
                      async (argv: any) => {
                        const args = argv.workspace ? [argv.workspace, argv.repo, argv['pr-id']] : [argv.repo, argv['pr-id']];
                        const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                        const prID = parseInt(idStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${idStr}`);
                        if (!argv.body) throw new Error('--body is required');
                        if (argv.file && !argv.line) throw new Error('--line is required when --file is specified');
                        if (argv.line && !argv.file) throw new Error('--file is required when --line is specified');
                        if (argv.file && argv['comment-id']) throw new Error('--comment-id cannot be used with --file/--line');
                        const client = getBitbucketClient(argv);
                        const req: CreatePRTaskRequest = { content: argv.body };
                        if (argv.file) {
                          const comment = await bbPRs.createPRComment(client, workspace, repoSlug, prID, argv.body, { path: argv.file, to: argv.line });
                          req.comment_id = comment.id;
                        } else if (argv['comment-id']) {
                          req.comment_id = argv['comment-id'];
                        }
                        const task = await bbPRs.createPRTask(client, workspace, repoSlug, prID, req);
                        outputResult(argv, 'created', String(task.id), `Created task #${task.id} on PR #${prID}`, task);
                      },
                    )
                    .command(
                      'update [workspace] <repo> <pr-id> <task-id>',
                      'Update a task on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .positional('task-id', { type: 'string' })
                        .option('body', { type: 'string', description: 'Updated task content' })
                        .option('state', { type: 'string', description: 'Task state (RESOLVED, UNRESOLVED)' }),
                      async (argv: any) => {
                        const rawArgs = argv.workspace
                          ? [argv.workspace, argv.repo, argv['pr-id'], argv['task-id']]
                          : [argv.repo, argv['pr-id'], argv['task-id']];
                        const [workspace, repoSlug, prIDStr, taskIDStr] = resolveWorkspaceRepoIDAndTaskID(argv, rawArgs);
                        const prID = parseInt(prIDStr, 10);
                        const taskID = parseInt(taskIDStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${prIDStr}`);
                        if (isNaN(taskID)) throw new Error(`invalid task ID: ${taskIDStr}`);
                        if (!argv.body && !argv.state) throw new Error('at least one of --body or --state is required');
                        const client = getBitbucketClient(argv);
                        const req: UpdatePRTaskRequest = {};
                        if (argv.body) req.content = argv.body;
                        if (argv.state) req.state = argv.state;
                        const task = await bbPRs.updatePRTask(client, workspace, repoSlug, prID, taskID, req);
                        outputResult(argv, 'updated', String(task.id), `Updated task #${task.id} on PR #${prID}`, task);
                      },
                    )
                    .command(
                      'resolve [workspace] <repo> <pr-id> <task-id>',
                      'Resolve a task on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .positional('task-id', { type: 'string' }),
                      async (argv: any) => {
                        const rawArgs = argv.workspace
                          ? [argv.workspace, argv.repo, argv['pr-id'], argv['task-id']]
                          : [argv.repo, argv['pr-id'], argv['task-id']];
                        const [workspace, repoSlug, prIDStr, taskIDStr] = resolveWorkspaceRepoIDAndTaskID(argv, rawArgs);
                        const prID = parseInt(prIDStr, 10);
                        const taskID = parseInt(taskIDStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${prIDStr}`);
                        if (isNaN(taskID)) throw new Error(`invalid task ID: ${taskIDStr}`);
                        const client = getBitbucketClient(argv);
                        const task = await bbPRs.updatePRTask(client, workspace, repoSlug, prID, taskID, { state: 'RESOLVED' });
                        outputResult(argv, 'resolved', String(task.id), `Resolved task #${task.id} on PR #${prID}`, task);
                      },
                    )
                    .command(
                      'delete [workspace] <repo> <pr-id> <task-id>',
                      'Delete a task on a pull request',
                      (y: Argv) => y
                        .positional('workspace', { type: 'string' })
                        .positional('repo', { type: 'string' })
                        .positional('pr-id', { type: 'string' })
                        .positional('task-id', { type: 'string' }),
                      async (argv: any) => {
                        const rawArgs = argv.workspace
                          ? [argv.workspace, argv.repo, argv['pr-id'], argv['task-id']]
                          : [argv.repo, argv['pr-id'], argv['task-id']];
                        const [workspace, repoSlug, prIDStr, taskIDStr] = resolveWorkspaceRepoIDAndTaskID(argv, rawArgs);
                        const prID = parseInt(prIDStr, 10);
                        const taskID = parseInt(taskIDStr, 10);
                        if (isNaN(prID)) throw new Error(`invalid PR ID: ${prIDStr}`);
                        if (isNaN(taskID)) throw new Error(`invalid task ID: ${taskIDStr}`);
                        const client = getBitbucketClient(argv);
                        await bbPRs.deletePRTask(client, workspace, repoSlug, prID, taskID);
                        outputResult(argv, 'deleted', String(taskID), `Deleted task #${taskID} from PR #${prID}`, null);
                      },
                    );
                },
              );
          },
        )

        // ─── pipeline ─────────────────────────────────────────────────────
        .command(
          ['pipeline', 'pipe'],
          'Manage pipelines',
          (plCmd: Argv) => {
            plCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List pipelines',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('status', { type: 'string', description: 'Filter by status (PENDING, BUILDING, PASSED, FAILED, etc.)' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const pipelines = await bbPipelines.listPipelines(client, workspace, repoSlug, {
                    status: argv.status,
                    ...getBBPaginationOpts(argv),
                  });
                  printTable(
                    ['BUILD#', 'STATUS', 'TRIGGER', 'TARGET', 'CREATED'],
                    pipelines.map((p) => {
                      const status = p.state?.result?.name ?? p.state?.name ?? '';
                      const target = p.target?.ref_name ?? '';
                      return [p.build_number, status, p.trigger?.name, target, p.created_on];
                    }),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <pipeline-uuid>',
                'Get pipeline details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pipeline-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pipeline-uuid']] : [argv.repo, argv['pipeline-uuid']];
                  const [workspace, repoSlug, pipelineUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const pipeline = await bbPipelines.getPipeline(client, workspace, repoSlug, pipelineUUID);
                  const status = pipeline.state?.result?.name ?? pipeline.state?.name ?? '';
                  console.log(`Build #:      ${pipeline.build_number}`);
                  console.log(`UUID:         ${pipeline.uuid}`);
                  console.log(`Status:       ${status}`);
                  console.log(`Trigger:      ${pipeline.trigger?.name}`);
                  console.log(`Target:       ${pipeline.target?.ref_name}`);
                  console.log(`Commit:       ${pipeline.target?.commit?.hash}`);
                  console.log(`Creator:      ${pipeline.creator?.display_name}`);
                  console.log(`Created:      ${pipeline.created_on}`);
                  if (pipeline.completed_on) console.log(`Completed:    ${pipeline.completed_on}`);
                  if (pipeline.build_seconds_used > 0) console.log(`Build Time:   ${pipeline.build_seconds_used}s`);
                },
              )
              .command(
                'run [workspace] <repo>',
                'Run a pipeline',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('branch', { type: 'string', description: 'Branch to run pipeline on (required)' })
                  .option('custom', { type: 'string', description: 'Custom pipeline pattern to run' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.branch) throw new Error('--branch is required');
                  const client = getBitbucketClient(argv);
                  const req = argv.custom
                    ? bbPipelines.newCustomPipelineRequest(argv.branch, argv.custom)
                    : bbPipelines.newBranchPipelineRequest(argv.branch);
                  const pipeline = await bbPipelines.runPipeline(client, workspace, repoSlug, req);
                  console.log(`Started pipeline build #${pipeline.build_number}`);
                  console.log(`UUID: ${pipeline.uuid}`);
                },
              )
              .command(
                'stop [workspace] <repo> <pipeline-uuid>',
                'Stop a running pipeline',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pipeline-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pipeline-uuid']] : [argv.repo, argv['pipeline-uuid']];
                  const [workspace, repoSlug, pipelineUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbPipelines.stopPipeline(client, workspace, repoSlug, pipelineUUID);
                  console.log(`Stopped pipeline ${pipelineUUID}`);
                },
              )
              .command(
                'steps [workspace] <repo> <pipeline-uuid>',
                'List steps for a pipeline',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pipeline-uuid', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['pipeline-uuid']] : [argv.repo, argv['pipeline-uuid']];
                  const [workspace, repoSlug, pipelineUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const steps = await bbPipelines.listPipelineSteps(client, workspace, repoSlug, pipelineUUID, getBBPaginationOpts(argv));
                  printTable(
                    ['UUID', 'NAME', 'STATUS', 'DURATION', 'IMAGE'],
                    steps.map((s) => {
                      const status = s.state?.result?.name ?? s.state?.name ?? '';
                      return [s.uuid, s.name, status, `${s.duration_in_seconds}s`, s.image?.name];
                    }),
                  );
                },
              )
              .command(
                'step [workspace] <repo> <pipeline-uuid> <step-uuid>',
                'Get details for a pipeline step',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pipeline-uuid', { type: 'string' })
                  .positional('step-uuid', { type: 'string' }),
                async (argv: any) => {
                  const client = getBitbucketClient(argv);
                  const step = await bbPipelines.getPipelineStep(client, argv.workspace ?? argv.repo, argv.workspace ? argv.repo : argv['pipeline-uuid'], argv.workspace ? argv['pipeline-uuid'] : argv['step-uuid'], argv.workspace ? argv['step-uuid'] : undefined);
                  if (isJSONOutput(argv)) { outputJSON(step); return; }
                  const status = step.state?.result?.name ?? step.state?.name ?? '';
                  console.log(`UUID:      ${step.uuid}`);
                  console.log(`Name:      ${step.name}`);
                  console.log(`Status:    ${status}`);
                  console.log(`Duration:  ${step.duration_in_seconds}s`);
                  console.log(`Image:     ${step.image?.name}`);
                },
              )
              .command(
                'log <workspace> <repo> <pipeline-uuid> <step-uuid>',
                'Get log output for a pipeline step',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('pipeline-uuid', { type: 'string' })
                  .positional('step-uuid', { type: 'string' }),
                async (argv: any) => {
                  const client = getBitbucketClient(argv);
                  const log = await bbPipelines.getStepLog(client, argv.workspace, argv.repo, argv['pipeline-uuid'], argv['step-uuid']);
                  process.stdout.write(log);
                },
              )
              .command(
                ['variables [workspace] <repo>', 'vars [workspace] <repo>'],
                'List pipeline variables',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const vars = await bbPipelines.listPipelineVariables(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  printTable(
                    ['UUID', 'KEY', 'VALUE', 'SECURED'],
                    vars.map((v) => [v.uuid, v.key, v.secured ? '***' : v.value, v.secured]),
                  );
                },
              )
              .command(
                'add-variable [workspace] <repo>',
                'Create a pipeline variable',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('key', { type: 'string', description: 'Variable key (required)' })
                  .option('value', { type: 'string', description: 'Variable value (required)' })
                  .option('secured', { type: 'boolean', description: 'Mark variable as secured', default: false }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.key || !argv.value) throw new Error('--key and --value are required');
                  const client = getBitbucketClient(argv);
                  const v = await bbPipelines.createPipelineVariable(client, workspace, repoSlug, argv.key, argv.value, argv.secured);
                  console.log(`Created variable: ${v.key} (UUID: ${v.uuid})`);
                },
              )
              .command(
                'delete-variable [workspace] <repo> <variable-uuid>',
                'Delete a pipeline variable',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('variable-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['variable-uuid']] : [argv.repo, argv['variable-uuid']];
                  const [workspace, repoSlug, varUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbPipelines.deletePipelineVariable(client, workspace, repoSlug, varUUID);
                  console.log(`Deleted variable ${varUUID}`);
                },
              );
          },
        )

        // ─── branch ───────────────────────────────────────────────────────
        .command(
          ['branch', 'br'],
          'Manage branches',
          (brCmd: Argv) => {
            brCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List branches',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('query', { type: 'string', description: 'Filter branches (e.g. name ~ "feature")' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const branches = await bbBranches.listBranches(client, workspace, repoSlug, {
                    q: argv.query,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(branches); return; }
                  printTable(
                    ['NAME', 'HASH', 'DATE', 'AUTHOR'],
                    branches.map((b) => [b.name, b.target?.hash?.slice(0, 12) ?? '', b.target?.date, b.target?.author?.raw]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <branch-name>',
                'Get branch details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('branch-name', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['branch-name']] : [argv.repo, argv['branch-name']];
                  const [workspace, repoSlug, branchName] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const branch = await bbBranches.getBranch(client, workspace, repoSlug, branchName);
                  console.log(`Name:     ${branch.name}`);
                  console.log(`Hash:     ${branch.target?.hash}`);
                  console.log(`Date:     ${branch.target?.date}`);
                  console.log(`Author:   ${branch.target?.author?.raw}`);
                  console.log(`Message:  ${branch.target?.message}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a branch',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('name', { type: 'string', description: 'Branch name (required)' })
                  .option('target', { type: 'string', description: 'Target commit hash (required)' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.name || !argv.target) throw new Error('--name and --target are required');
                  const client = getBitbucketClient(argv);
                  const branch = await bbBranches.createBranch(client, workspace, repoSlug, {
                    name: argv.name,
                    target: { hash: argv.target },
                  });
                  console.log(`Created branch: ${branch.name} (${branch.target?.hash})`);
                },
              )
              .command(
                'delete [workspace] <repo> <branch-name>',
                'Delete a branch',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('branch-name', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['branch-name']] : [argv.repo, argv['branch-name']];
                  const [workspace, repoSlug, branchName] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbBranches.deleteBranch(client, workspace, repoSlug, branchName);
                  console.log(`Deleted branch: ${branchName}`);
                },
              );
          },
        )

        // ─── tag ──────────────────────────────────────────────────────────
        .command(
          'tag',
          'Manage tags',
          (tagCmd: Argv) => {
            tagCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List tags',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('query', { type: 'string', description: 'Filter tags (e.g. name ~ "v1")' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const tags = await bbBranches.listTags(client, workspace, repoSlug, {
                    q: argv.query,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(tags); return; }
                  printTable(
                    ['NAME', 'HASH', 'DATE', 'MESSAGE'],
                    tags.map((t) => {
                      const hash = t.target?.hash?.slice(0, 12) ?? '';
                      const msg = truncate(t.message ?? '', 60);
                      return [t.name, hash, t.target?.date, msg];
                    }),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <tag-name>',
                'Get tag details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('tag-name', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['tag-name']] : [argv.repo, argv['tag-name']];
                  const [workspace, repoSlug, tagName] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const tag = await bbBranches.getTag(client, workspace, repoSlug, tagName);
                  console.log(`Name:     ${tag.name}`);
                  console.log(`Hash:     ${tag.target?.hash}`);
                  console.log(`Date:     ${tag.target?.date}`);
                  if (tag.message) console.log(`Message:  ${tag.message}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a tag',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('name', { type: 'string', description: 'Tag name (required)' })
                  .option('target', { type: 'string', description: 'Target commit hash (required)' })
                  .option('message', { type: 'string', description: 'Tag message' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.name || !argv.target) throw new Error('--name and --target are required');
                  const client = getBitbucketClient(argv);
                  const tag = await bbBranches.createTag(client, workspace, repoSlug, {
                    name: argv.name,
                    message: argv.message,
                    target: { hash: argv.target },
                  });
                  console.log(`Created tag: ${tag.name} (${tag.target?.hash})`);
                },
              )
              .command(
                'delete [workspace] <repo> <tag-name>',
                'Delete a tag',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('tag-name', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['tag-name']] : [argv.repo, argv['tag-name']];
                  const [workspace, repoSlug, tagName] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbBranches.deleteTag(client, workspace, repoSlug, tagName);
                  console.log(`Deleted tag: ${tagName}`);
                },
              );
          },
        )

        // ─── commit ───────────────────────────────────────────────────────
        .command(
          ['commit', 'cm'],
          'Manage commits',
          (cmCmd: Argv) => {
            cmCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List commits',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('include', { type: 'string', description: 'Include commits reachable from this ref' })
                  .option('exclude', { type: 'string', description: 'Exclude commits reachable from this ref' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const commits = await bbCommits.listCommits(client, workspace, repoSlug, {
                    include: argv.include,
                    exclude: argv.exclude,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(commits); return; }
                  printTable(
                    ['HASH', 'DATE', 'AUTHOR', 'MESSAGE'],
                    commits.map((c) => {
                      const hash = c.hash?.slice(0, 12) ?? '';
                      const msg = truncate(firstLine(c.message ?? ''), 60);
                      return [hash, c.date, c.author?.raw, msg];
                    }),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <commit-hash>',
                'Get commit details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('commit-hash', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['commit-hash']] : [argv.repo, argv['commit-hash']];
                  const [workspace, repoSlug, commitHash] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const commit = await bbCommits.getCommit(client, workspace, repoSlug, commitHash);
                  console.log(`Hash:     ${commit.hash}`);
                  console.log(`Date:     ${commit.date}`);
                  console.log(`Author:   ${commit.author?.raw}`);
                  console.log(`Message:  ${commit.message}`);
                  if (commit.parents?.length > 0) {
                    const parents = commit.parents.map((p) => p.hash?.slice(0, 12)).join(', ');
                    console.log(`Parents:  ${parents}`);
                  }
                  console.log(`URL:      ${commit.links?.html?.href}`);
                },
              )
              .command(
                'statuses [workspace] <repo> <commit-hash>',
                'List build statuses for a commit',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('commit-hash', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['commit-hash']] : [argv.repo, argv['commit-hash']];
                  const [workspace, repoSlug, commitHash] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const statuses = await bbCommits.listCommitStatuses(client, workspace, repoSlug, commitHash, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(statuses); return; }
                  printTable(
                    ['KEY', 'STATE', 'NAME', 'DESCRIPTION', 'UPDATED'],
                    statuses.map((s) => [s.key, s.state, s.name, s.description, s.updated_on]),
                  );
                },
              )
              .command(
                'diff [workspace] <repo> <spec>',
                'Get diff between two commits (e.g. commit1..commit2)',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('spec', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv.spec] : [argv.repo, argv.spec];
                  const [workspace, repoSlug, spec] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const diff = await bbCommits.getDiff(client, workspace, repoSlug, spec);
                  process.stdout.write(diff);
                },
              );
          },
        )

        // ─── workspace ────────────────────────────────────────────────────
        .command(
          ['workspace', 'ws'],
          'Manage workspaces',
          (wsCmd: Argv) => {
            wsCmd
              .command(
                ['list', 'ls'],
                'List workspaces',
                (y: Argv) => y.options(paginationOptions),
                async (argv: any) => {
                  const client = getBitbucketClient(argv);
                  const workspaces = await bbWorkspaces.listWorkspaces(client, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(workspaces); return; }
                  printTable(
                    ['SLUG', 'NAME', 'UUID'],
                    workspaces.map((ws) => [ws.slug, ws.name, ws.uuid]),
                  );
                },
              )
              .command(
                'get [workspace]',
                'Get workspace details',
                (y: Argv) => y.positional('workspace', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const ws = await bbWorkspaces.getWorkspace(client, workspace);
                  console.log(`Name:     ${ws.name}`);
                  console.log(`Slug:     ${ws.slug}`);
                  console.log(`UUID:     ${ws.uuid}`);
                  console.log(`Private:  ${ws.is_private}`);
                  console.log(`Created:  ${ws.created_on}`);
                  console.log(`URL:      ${ws.links?.html?.href}`);
                },
              )
              .command(
                'members [workspace]',
                'List workspace members',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const members = await bbWorkspaces.listWorkspaceMembers(client, workspace, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(members); return; }
                  printTable(
                    ['DISPLAY NAME', 'NICKNAME', 'UUID'],
                    members.map((m) => [m.user?.display_name, m.user?.nickname, m.user?.uuid]),
                  );
                },
              )
              .command(
                'permissions [workspace]',
                'List user permissions in workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const perms = await bbWorkspaces.listWorkspacePermissions(client, workspace, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(perms); return; }
                  printTable(
                    ['USER', 'PERMISSION'],
                    perms.map((p) => [p.user?.display_name, p.permission]),
                  );
                },
              );
          },
        )

        // ─── project ──────────────────────────────────────────────────────
        .command(
          ['project', 'proj'],
          'Manage projects',
          (projCmd: Argv) => {
            projCmd
              .command(
                ['list [workspace]', 'ls [workspace]'],
                'List projects in a workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const projects = await bbProjects.listProjects(client, workspace, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(projects); return; }
                  printTable(
                    ['KEY', 'NAME', 'PRIVATE', 'UPDATED'],
                    projects.map((p) => [p.key, p.name, p.is_private, p.updated_on]),
                  );
                },
              )
              .command(
                'get [workspace] <project-key>',
                'Get project details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('project-key', { type: 'string' }),
                async (argv: any) => {
                  let workspace: string, projectKey: string;
                  if (argv.workspace && argv['project-key']) {
                    workspace = argv.workspace;
                    projectKey = argv['project-key'];
                  } else {
                    workspace = defaultWorkspace(argv, [], 0);
                    projectKey = argv.workspace ?? argv['project-key'];
                  }
                  const client = getBitbucketClient(argv);
                  const project = await bbProjects.getProject(client, workspace, projectKey);
                  console.log(`Key:          ${project.key}`);
                  console.log(`Name:         ${project.name}`);
                  console.log(`Description:  ${project.description}`);
                  console.log(`Private:      ${project.is_private}`);
                  console.log(`Created:      ${project.created_on}`);
                  console.log(`Updated:      ${project.updated_on}`);
                  console.log(`URL:          ${project.links?.html?.href}`);
                },
              )
              .command(
                'create [workspace]',
                'Create a project',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .option('name', { type: 'string', description: 'Project name (required)' })
                  .option('key', { type: 'string', description: 'Project key (required)' })
                  .option('description', { type: 'string', description: 'Project description' })
                  .option('private', { type: 'boolean', description: 'Make project private', default: true }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  if (!argv.name || !argv.key) throw new Error('--name and --key are required');
                  const client = getBitbucketClient(argv);
                  const project = await bbProjects.createProject(client, workspace, {
                    name: argv.name,
                    key: argv.key,
                    description: argv.description,
                    is_private: argv.private,
                  });
                  console.log(`Created project: ${project.name} (${project.key})`);
                  console.log(`URL: ${project.links?.html?.href}`);
                },
              )
              .command(
                'delete [workspace] <project-key>',
                'Delete a project',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('project-key', { type: 'string' }),
                async (argv: any) => {
                  let workspace: string, projectKey: string;
                  if (argv.workspace && argv['project-key']) {
                    workspace = argv.workspace;
                    projectKey = argv['project-key'];
                  } else {
                    workspace = defaultWorkspace(argv, [], 0);
                    projectKey = argv.workspace ?? argv['project-key'];
                  }
                  const client = getBitbucketClient(argv);
                  await bbProjects.deleteProject(client, workspace, projectKey);
                  console.log(`Deleted project: ${projectKey}`);
                },
              );
          },
        )

        // ─── webhook ──────────────────────────────────────────────────────
        .command(
          ['webhook', 'hook'],
          'Manage webhooks',
          (whCmd: Argv) => {
            whCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List webhooks for a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const hooks = await bbWebhooks.listRepoWebhooks(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(hooks); return; }
                  printTable(
                    ['UUID', 'DESCRIPTION', 'URL', 'ACTIVE', 'EVENTS'],
                    hooks.map((h) => [h.uuid, h.description, h.url, h.active, (h.events ?? []).join(',')]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <webhook-uuid>',
                'Get webhook details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('webhook-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['webhook-uuid']] : [argv.repo, argv['webhook-uuid']];
                  const [workspace, repoSlug, hookUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const hook = await bbWebhooks.getRepoWebhook(client, workspace, repoSlug, hookUUID);
                  console.log(`UUID:         ${hook.uuid}`);
                  console.log(`Description:  ${hook.description}`);
                  console.log(`URL:          ${hook.url}`);
                  console.log(`Active:       ${hook.active}`);
                  console.log(`Events:       ${(hook.events ?? []).join(', ')}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a webhook for a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('url', { type: 'string', description: 'Webhook URL (required)' })
                  .option('description', { type: 'string', description: 'Webhook description' })
                  .option('events', { type: 'array', description: 'Events to subscribe to (required, e.g. repo:push pullrequest:created)' })
                  .option('active', { type: 'boolean', description: 'Whether the webhook is active', default: true }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.url || !argv.events?.length) throw new Error('--url and --events are required');
                  const client = getBitbucketClient(argv);
                  const hook = await bbWebhooks.createRepoWebhook(client, workspace, repoSlug, {
                    description: argv.description,
                    url: argv.url,
                    active: argv.active,
                    events: argv.events,
                  });
                  console.log(`Created webhook: ${hook.uuid}`);
                  console.log(`URL: ${hook.url}`);
                },
              )
              .command(
                'update [workspace] <repo> <webhook-uuid>',
                'Update a webhook',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('webhook-uuid', { type: 'string' })
                  .option('url', { type: 'string', description: 'Webhook URL' })
                  .option('description', { type: 'string', description: 'Webhook description' })
                  .option('events', { type: 'array', description: 'Events to subscribe to' })
                  .option('active', { type: 'boolean', description: 'Whether the webhook is active' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['webhook-uuid']] : [argv.repo, argv['webhook-uuid']];
                  const [workspace, repoSlug, hookUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const req: Partial<CreateWebhookRequest> = {};
                  if (argv.url) req.url = argv.url;
                  if (argv.description !== undefined) req.description = argv.description;
                  if (argv.events) req.events = argv.events;
                  if (argv.active !== undefined) req.active = argv.active;
                  const hook = await bbWebhooks.updateRepoWebhook(client, workspace, repoSlug, hookUUID, req as CreateWebhookRequest);
                  outputResult(argv, 'updated', hook.uuid, `Updated webhook: ${hook.uuid}`, hook);
                },
              )
              .command(
                'delete [workspace] <repo> <webhook-uuid>',
                'Delete a webhook',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('webhook-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['webhook-uuid']] : [argv.repo, argv['webhook-uuid']];
                  const [workspace, repoSlug, hookUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbWebhooks.deleteRepoWebhook(client, workspace, repoSlug, hookUUID);
                  console.log(`Deleted webhook: ${hookUUID}`);
                },
              )
              .command(
                'list-workspace [workspace]',
                'List webhooks for a workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const hooks = await bbWebhooks.listWorkspaceWebhooks(client, workspace, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(hooks); return; }
                  printTable(
                    ['UUID', 'DESCRIPTION', 'URL', 'ACTIVE', 'EVENTS'],
                    hooks.map((h) => [h.uuid, h.description, h.url, h.active, (h.events ?? []).join(',')]),
                  );
                },
              )
              .command(
                'create-workspace [workspace]',
                'Create a webhook for a workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .option('url', { type: 'string', description: 'Webhook URL (required)' })
                  .option('description', { type: 'string', description: 'Webhook description' })
                  .option('events', { type: 'array', description: 'Events to subscribe to (required)' })
                  .option('active', { type: 'boolean', description: 'Whether the webhook is active', default: true }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  if (!argv.url || !argv.events?.length) throw new Error('--url and --events are required');
                  const client = getBitbucketClient(argv);
                  const hook = await bbWebhooks.createWorkspaceWebhook(client, workspace, {
                    description: argv.description,
                    url: argv.url,
                    active: argv.active,
                    events: argv.events,
                  });
                  console.log(`Created workspace webhook: ${hook.uuid}`);
                },
              )
              .command(
                'delete-workspace [workspace] <webhook-uuid>',
                'Delete a workspace webhook',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('webhook-uuid', { type: 'string' }),
                async (argv: any) => {
                  let workspace: string, hookUUID: string;
                  if (argv.workspace && argv['webhook-uuid']) {
                    workspace = argv.workspace;
                    hookUUID = argv['webhook-uuid'];
                  } else {
                    workspace = defaultWorkspace(argv, [], 0);
                    hookUUID = argv.workspace ?? argv['webhook-uuid'];
                  }
                  const client = getBitbucketClient(argv);
                  await bbWebhooks.deleteWorkspaceWebhook(client, workspace, hookUUID);
                  console.log(`Deleted workspace webhook: ${hookUUID}`);
                },
              );
          },
        )

        // ─── environment ──────────────────────────────────────────────────
        .command(
          ['environment', 'env'],
          'Manage deployment environments',
          (envCmd: Argv) => {
            envCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List deployment environments',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const envs = await bbEnvironments.listEnvironments(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(envs); return; }
                  printTable(
                    ['UUID', 'NAME', 'TYPE', 'RANK'],
                    envs.map((e) => [e.uuid, e.name, e.environment_type?.name, e.rank]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <environment-uuid>',
                'Get environment details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('environment-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['environment-uuid']] : [argv.repo, argv['environment-uuid']];
                  const [workspace, repoSlug, envUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const env = await bbEnvironments.getEnvironment(client, workspace, repoSlug, envUUID);
                  console.log(`UUID:  ${env.uuid}`);
                  console.log(`Name:  ${env.name}`);
                  console.log(`Type:  ${env.environment_type?.name}`);
                  console.log(`Rank:  ${env.rank}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a deployment environment',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('name', { type: 'string', description: 'Environment name (required)' })
                  .option('type', { type: 'string', description: 'Environment type (Test, Staging, Production)', default: 'Test' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.name) throw new Error('--name is required');
                  const envType = argv.type || 'Test';
                  const rankMap: Record<string, number> = { Production: 2, Staging: 1 };
                  const rank = rankMap[envType] ?? 0;
                  const client = getBitbucketClient(argv);
                  const env = await bbEnvironments.createEnvironment(client, workspace, repoSlug, {
                    name: argv.name,
                    environment_type: { name: envType, rank },
                  });
                  console.log(`Created environment: ${env.name} (UUID: ${env.uuid})`);
                },
              )
              .command(
                'delete [workspace] <repo> <environment-uuid>',
                'Delete a deployment environment',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('environment-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['environment-uuid']] : [argv.repo, argv['environment-uuid']];
                  const [workspace, repoSlug, envUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbEnvironments.deleteEnvironment(client, workspace, repoSlug, envUUID);
                  console.log(`Deleted environment: ${envUUID}`);
                },
              );
          },
        )

        // ─── deploy-key ───────────────────────────────────────────────────
        .command(
          ['deploy-key', 'dk'],
          'Manage deploy keys',
          (dkCmd: Argv) => {
            dkCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List deploy keys',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const keys = await bbDeployKeys.listDeployKeys(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(keys); return; }
                  printTable(
                    ['ID', 'LABEL', 'COMMENT', 'CREATED'],
                    keys.map((k) => [k.id, k.label, k.comment, k.created_on]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <key-id>',
                'Get deploy key details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('key-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['key-id']] : [argv.repo, argv['key-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const keyID = parseInt(idStr, 10);
                  if (isNaN(keyID)) throw new Error(`invalid key ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const key = await bbDeployKeys.getDeployKey(client, workspace, repoSlug, keyID);
                  console.log(`ID:       ${key.id}`);
                  console.log(`Label:    ${key.label}`);
                  console.log(`Comment:  ${key.comment}`);
                  console.log(`Created:  ${key.created_on}`);
                  console.log(`Key:      ${key.key}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Add a deploy key to a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('key', { type: 'string', description: 'SSH public key content (required)' })
                  .option('label', { type: 'string', description: 'Label for the key (required)' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.key || !argv.label) throw new Error('--key and --label are required');
                  const client = getBitbucketClient(argv);
                  const key = await bbDeployKeys.createDeployKey(client, workspace, repoSlug, {
                    key: argv.key,
                    label: argv.label,
                  });
                  console.log(`Created deploy key: ${key.label} (ID: ${key.id})`);
                },
              )
              .command(
                'delete [workspace] <repo> <key-id>',
                'Delete a deploy key',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('key-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['key-id']] : [argv.repo, argv['key-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const keyID = parseInt(idStr, 10);
                  if (isNaN(keyID)) throw new Error(`invalid key ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbDeployKeys.deleteDeployKey(client, workspace, repoSlug, keyID);
                  console.log(`Deleted deploy key: ${idStr}`);
                },
              );
          },
        )

        // ─── download ─────────────────────────────────────────────────────
        .command(
          ['download', 'dl'],
          'Manage repository downloads',
          (dlCmd: Argv) => {
            dlCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List downloads for a repository',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const downloads = await bbDownloads.listDownloads(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(downloads); return; }
                  printTable(
                    ['NAME', 'SIZE', 'DOWNLOADS', 'CREATED'],
                    downloads.map((d) => [d.name, d.size, d.downloads, d.created_on]),
                  );
                },
              )
              .command(
                'delete [workspace] <repo> <filename>',
                'Delete a download artifact',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('filename', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv.filename] : [argv.repo, argv.filename];
                  const [workspace, repoSlug, filename] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  await bbDownloads.deleteDownload(client, workspace, repoSlug, filename);
                  console.log(`Deleted download: ${filename}`);
                },
              );
          },
        )

        // ─── snippet ──────────────────────────────────────────────────────
        .command(
          ['snippet', 'snip'],
          'Manage snippets',
          (snCmd: Argv) => {
            snCmd
              .command(
                ['list [workspace]', 'ls [workspace]'],
                'List snippets in a workspace',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  const client = getBitbucketClient(argv);
                  const snippets = await bbSnippets.listSnippets(client, workspace, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(snippets); return; }
                  printTable(
                    ['ID', 'TITLE', 'PRIVATE', 'CREATED', 'OWNER'],
                    snippets.map((s) => [s.id, s.title, s.is_private, s.created_on, s.owner?.display_name]),
                  );
                },
              )
              .command(
                'get [workspace] <snippet-id>',
                'Get snippet details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('snippet-id', { type: 'string' }),
                async (argv: any) => {
                  let workspace: string, snippetID: string;
                  if (argv.workspace && argv['snippet-id']) {
                    workspace = argv.workspace;
                    snippetID = argv['snippet-id'];
                  } else {
                    workspace = defaultWorkspace(argv, [], 0);
                    snippetID = argv.workspace ?? argv['snippet-id'];
                  }
                  const client = getBitbucketClient(argv);
                  const snippet = await bbSnippets.getSnippet(client, workspace, snippetID);
                  console.log(`ID:       ${snippet.id}`);
                  console.log(`Title:    ${snippet.title}`);
                  console.log(`Private:  ${snippet.is_private}`);
                  console.log(`Owner:    ${snippet.owner?.display_name}`);
                  console.log(`Created:  ${snippet.created_on}`);
                  console.log(`Updated:  ${snippet.updated_on}`);
                  console.log(`URL:      ${snippet.links?.html?.href}`);
                },
              )
              .command(
                'create [workspace]',
                'Create a snippet',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .option('title', { type: 'string', description: 'Snippet title (required)' })
                  .option('private', { type: 'boolean', description: 'Make snippet private', default: false }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace] : [];
                  const workspace = defaultWorkspace(argv, args, 0);
                  if (!argv.title) throw new Error('--title is required');
                  const client = getBitbucketClient(argv);
                  const snippet = await bbSnippets.createSnippet(client, workspace, {
                    title: argv.title,
                    is_private: argv.private,
                  });
                  console.log(`Created snippet: ${snippet.title} (ID: ${snippet.id})`);
                  console.log(`URL: ${snippet.links?.html?.href}`);
                },
              )
              .command(
                'delete [workspace] <snippet-id>',
                'Delete a snippet',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('snippet-id', { type: 'string' }),
                async (argv: any) => {
                  let workspace: string, snippetID: string;
                  if (argv.workspace && argv['snippet-id']) {
                    workspace = argv.workspace;
                    snippetID = argv['snippet-id'];
                  } else {
                    workspace = defaultWorkspace(argv, [], 0);
                    snippetID = argv.workspace ?? argv['snippet-id'];
                  }
                  const client = getBitbucketClient(argv);
                  await bbSnippets.deleteSnippet(client, workspace, snippetID);
                  console.log(`Deleted snippet: ${snippetID}`);
                },
              );
          },
        )

        // ─── issue ────────────────────────────────────────────────────────
        .command(
          'issue',
          'Manage repository issues (Bitbucket issue tracker)',
          (issCmd: Argv) => {
            issCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List issues',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('query', { type: 'string', description: 'Filter issues (Bitbucket query syntax)' })
                  .option('sort', { type: 'string', description: 'Sort field (e.g. -priority)' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const issues = await bbIssues.listIssues(client, workspace, repoSlug, {
                    q: argv.query,
                    sort: argv.sort,
                    ...getBBPaginationOpts(argv),
                  });
                  if (isJSONOutput(argv)) { outputJSON(issues); return; }
                  printTable(
                    ['ID', 'TITLE', 'STATE', 'PRIORITY', 'KIND', 'ASSIGNEE'],
                    issues.map((issue) => [
                      issue.id,
                      issue.title,
                      issue.state,
                      issue.priority,
                      issue.kind,
                      issue.assignee?.display_name ?? '',
                    ]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <issue-id>',
                'Get issue details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('issue-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['issue-id']] : [argv.repo, argv['issue-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const issueID = parseInt(idStr, 10);
                  if (isNaN(issueID)) throw new Error(`invalid issue ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const issue = await bbIssues.getIssue(client, workspace, repoSlug, issueID);
                  if (isJSONOutput(argv)) { outputJSON(issue); return; }
                  console.log(`ID:        ${issue.id}`);
                  console.log(`Title:     ${issue.title}`);
                  console.log(`State:     ${issue.state}`);
                  console.log(`Priority:  ${issue.priority}`);
                  console.log(`Kind:      ${issue.kind}`);
                  console.log(`Reporter:  ${issue.reporter?.display_name}`);
                  if (issue.assignee) console.log(`Assignee:  ${issue.assignee.display_name}`);
                  if (issue.component) console.log(`Component: ${issue.component.name}`);
                  if (issue.milestone) console.log(`Milestone: ${issue.milestone.name}`);
                  console.log(`Votes:     ${issue.votes}`);
                  console.log(`Created:   ${issue.created_on}`);
                  console.log(`Updated:   ${issue.updated_on}`);
                  console.log(`URL:       ${issue.links?.html?.href}`);
                  if (issue.content?.raw) console.log(`\nContent:\n${issue.content.raw}`);
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create an issue',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('title', { type: 'string', description: 'Issue title (required)' })
                  .option('content', { type: 'string', description: 'Issue content/description' })
                  .option('kind', { type: 'string', description: 'Issue kind (bug, enhancement, proposal, task)', default: 'bug' })
                  .option('priority', { type: 'string', description: 'Issue priority (trivial, minor, major, critical, blocker)', default: 'major' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.title) throw new Error('--title is required');
                  const client = getBitbucketClient(argv);
                  const req: CreateIssueRequest = {
                    title: argv.title,
                    kind: argv.kind,
                    priority: argv.priority,
                  };
                  if (argv.content) req.content = { raw: argv.content };
                  const issue = await bbIssues.createIssue(client, workspace, repoSlug, req);
                  console.log(`Created issue #${issue.id}: ${issue.title}`);
                  console.log(`URL: ${issue.links?.html?.href}`);
                },
              )
              .command(
                'update [workspace] <repo> <issue-id>',
                'Update an issue',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('issue-id', { type: 'string' })
                  .option('title', { type: 'string', description: 'New title' })
                  .option('state', { type: 'string', description: 'New state (new, open, resolved, on hold, invalid, duplicate, wontfix, closed)' })
                  .option('kind', { type: 'string', description: 'New kind' })
                  .option('priority', { type: 'string', description: 'New priority' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['issue-id']] : [argv.repo, argv['issue-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const issueID = parseInt(idStr, 10);
                  if (isNaN(issueID)) throw new Error(`invalid issue ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const req = {
                    title: argv.title,
                    state: argv.state,
                    kind: argv.kind,
                    priority: argv.priority,
                  };
                  const issue = await bbIssues.updateIssue(client, workspace, repoSlug, issueID, req);
                  console.log(`Updated issue #${issue.id}: ${issue.title}`);
                },
              )
              .command(
                'delete [workspace] <repo> <issue-id>',
                'Delete an issue',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('issue-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['issue-id']] : [argv.repo, argv['issue-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const issueID = parseInt(idStr, 10);
                  if (isNaN(issueID)) throw new Error(`invalid issue ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbIssues.deleteIssue(client, workspace, repoSlug, issueID);
                  console.log(`Deleted issue #${issueID}`);
                },
              )
              .command(
                'comments [workspace] <repo> <issue-id>',
                'List comments on an issue',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('issue-id', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['issue-id']] : [argv.repo, argv['issue-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const issueID = parseInt(idStr, 10);
                  if (isNaN(issueID)) throw new Error(`invalid issue ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const comments = await bbIssues.listIssueComments(client, workspace, repoSlug, issueID, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(comments); return; }
                  for (const c of comments) {
                    console.log(`#${c.id} by ${c.user?.display_name} (${c.created_on})`);
                    console.log(`  ${c.content?.raw}\n`);
                  }
                },
              )
              .command(
                'comment [workspace] <repo> <issue-id>',
                'Add a comment to an issue',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('issue-id', { type: 'string' })
                  .option('body', { type: 'string', description: 'Comment body (required)' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['issue-id']] : [argv.repo, argv['issue-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const issueID = parseInt(idStr, 10);
                  if (isNaN(issueID)) throw new Error(`invalid issue ID: ${idStr}`);
                  if (!argv.body) throw new Error('--body is required');
                  const client = getBitbucketClient(argv);
                  const comment = await bbIssues.createIssueComment(client, workspace, repoSlug, issueID, argv.body);
                  console.log(`Added comment #${comment.id} to issue #${issueID}`);
                },
              );
          },
        )

        // ─── search ───────────────────────────────────────────────────────
        .command(
          'search',
          'Search code in a workspace',
          (searchCmd: Argv) => {
            searchCmd.command(
              'code [workspace]',
              'Search for code in a workspace',
              (y: Argv) => y
                .positional('workspace', { type: 'string' })
                .option('query', { type: 'string', description: 'Search query (required)' })
                .options(paginationOptions),
              async (argv: any) => {
                const args = argv.workspace ? [argv.workspace] : [];
                const workspace = defaultWorkspace(argv, args, 0);
                if (!argv.query) throw new Error('--query is required');
                const client = getBitbucketClient(argv);
                const results = await bbSearch.searchCode(client, workspace, argv.query, getBBPaginationOpts(argv));
                if (isJSONOutput(argv)) { outputJSON(results); return; }
                console.log(`Found ${results.size ?? 0} results\n`);
                for (const r of results.values ?? []) {
                  console.log(`File: ${r.file?.path} (${r.content_match_count} matches)`);
                  for (const m of r.content_matches ?? []) {
                    for (const line of m.lines ?? []) {
                      const parts = (line.segments ?? []).map((seg) =>
                        seg.match ? `[${seg.text}]` : seg.text,
                      );
                      console.log(`  ${line.line}: ${parts.join('')}`);
                    }
                  }
                  console.log();
                }
              },
            );
          },
        )

        // ─── deployment ───────────────────────────────────────────────────
        .command(
          ['deployment', 'deploy'],
          'Manage deployments',
          (depCmd: Argv) => {
            depCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List deployments',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const deployments = await bbDeployments.listDeployments(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(deployments); return; }
                  printTable(
                    ['UUID', 'ENVIRONMENT', 'STATUS', 'RELEASE', 'COMMIT', 'CREATED'],
                    deployments.map((d) => {
                      const hash = d.release?.commit?.hash?.slice(0, 12) ?? '';
                      return [d.uuid, d.environment?.name, d.state?.status?.name, d.release?.name, hash, d.release?.created_on];
                    }),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <deployment-uuid>',
                'Get deployment details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('deployment-uuid', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['deployment-uuid']] : [argv.repo, argv['deployment-uuid']];
                  const [workspace, repoSlug, deployUUID] = resolveWorkspaceRepoAndID(argv, args);
                  const client = getBitbucketClient(argv);
                  const d = await bbDeployments.getDeployment(client, workspace, repoSlug, deployUUID);
                  console.log(`UUID:         ${d.uuid}`);
                  console.log(`State:        ${d.state?.name}`);
                  console.log(`Status:       ${d.state?.status?.name}`);
                  console.log(`Environment:  ${d.environment?.name}`);
                  console.log(`Release:      ${d.release?.name}`);
                  console.log(`Commit:       ${d.release?.commit?.hash}`);
                  if (d.release?.commit?.message) console.log(`Message:      ${d.release.commit.message}`);
                  console.log(`Created:      ${d.release?.created_on}`);
                },
              );
          },
        )

        // ─── branch-restriction ───────────────────────────────────────────
        .command(
          ['branch-restriction', 'restriction'],
          'Manage branch restrictions',
          (brCmd: Argv) => {
            brCmd
              .command(
                ['list [workspace] <repo>', 'ls [workspace] <repo>'],
                'List branch restrictions',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .options(paginationOptions),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  const client = getBitbucketClient(argv);
                  const restrictions = await bbBranchRestrictions.listBranchRestrictions(client, workspace, repoSlug, getBBPaginationOpts(argv));
                  if (isJSONOutput(argv)) { outputJSON(restrictions); return; }
                  printTable(
                    ['ID', 'KIND', 'PATTERN'],
                    restrictions.map((r) => [r.id, r.kind, r.pattern]),
                  );
                },
              )
              .command(
                'get [workspace] <repo> <restriction-id>',
                'Get branch restriction details',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('restriction-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['restriction-id']] : [argv.repo, argv['restriction-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const id = parseInt(idStr, 10);
                  if (isNaN(id)) throw new Error(`invalid restriction ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  const r = await bbBranchRestrictions.getBranchRestriction(client, workspace, repoSlug, id);
                  console.log(`ID:       ${r.id}`);
                  console.log(`Kind:     ${r.kind}`);
                  console.log(`Pattern:  ${r.pattern}`);
                  if (r.value != null) console.log(`Value:    ${r.value}`);
                  if (r.users?.length > 0) {
                    console.log('Users:');
                    for (const u of r.users) console.log(`  - ${u.display_name}`);
                  }
                  if (r.groups?.length > 0) {
                    console.log('Groups:');
                    for (const g of r.groups) console.log(`  - ${g.name} (${g.slug})`);
                  }
                },
              )
              .command(
                'create [workspace] <repo>',
                'Create a branch restriction',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .option('kind', { type: 'string', description: 'Restriction kind (e.g. push, force, delete, restrict_merges, require_approvals_to_merge)' })
                  .option('pattern', { type: 'string', description: 'Branch pattern (e.g. main, release/*)' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo] : [argv.repo];
                  const [workspace, repoSlug] = resolveWorkspaceAndRepo(argv, args);
                  if (!argv.kind || !argv.pattern) throw new Error('--kind and --pattern are required');
                  const client = getBitbucketClient(argv);
                  const r = await bbBranchRestrictions.createBranchRestriction(client, workspace, repoSlug, {
                    kind: argv.kind,
                    pattern: argv.pattern,
                  });
                  console.log(`Created branch restriction: ${r.id} (${r.kind} on ${r.pattern})`);
                },
              )
              .command(
                'delete [workspace] <repo> <restriction-id>',
                'Delete a branch restriction',
                (y: Argv) => y
                  .positional('workspace', { type: 'string' })
                  .positional('repo', { type: 'string' })
                  .positional('restriction-id', { type: 'string' }),
                async (argv: any) => {
                  const args = argv.workspace ? [argv.workspace, argv.repo, argv['restriction-id']] : [argv.repo, argv['restriction-id']];
                  const [workspace, repoSlug, idStr] = resolveWorkspaceRepoAndID(argv, args);
                  const id = parseInt(idStr, 10);
                  if (isNaN(id)) throw new Error(`invalid restriction ID: ${idStr}`);
                  const client = getBitbucketClient(argv);
                  await bbBranchRestrictions.deleteBranchRestriction(client, workspace, repoSlug, id);
                  console.log(`Deleted branch restriction: ${id}`);
                },
              );
          },
        );
}

import {
  getJiraClient,
  defaultProject,
  isJSONOutput,
  outputJSON,
  outputResult,
  truncate,
  printIssueRow,
  printPaginationHint,
} from './helpers.js';
import * as jiraIssues from '../internal/jira/issues.js';
import * as jiraProjects from '../internal/jira/projects.js';
import * as jiraAgile from '../internal/jira/agile.js';
import * as jiraAdmin from '../internal/jira/admin.js';
import * as jiraSchemes from '../internal/jira/schemes.js';
import { render as renderADF } from '../internal/adf/render.js';
import type { Argv } from 'yargs';
import type { ADFDocument, JsonBody } from '../internal/types.js';

// ============================================================================
// JiraArgv — a single, narrow-typed argv shape that covers every option and
// positional used by the jira subcommands. This replaces the older loose
// handler signatures without losing runtime behaviour.
//
// Design notes:
//   - All fields are optional: yargs only populates the keys that the active
//     command declares, so per-handler unions would be noisy. Handlers simply
//     never touch fields they didn't declare.
//   - Required positionals / options-with-defaults / demandOption options
//     still arrive as populated at runtime, so accessing e.g. `argv.jql` in
//     a handler that declared `--jql` as required is safe.
//   - BaseArgv's common flags (profile, output, json, page, pagelen) are
//     mirrored so `getJiraClient`/`defaultProject`/etc. accept us.
//   - Runtime values for array options arrive as `string[]`, so we use that
//     narrow type (yargs `type: 'array'` widens to `(string|number)[]` in
//     the type inference but every usage here treats them as string[]).
//   - `project` is kept as `string`. The one array-typed variant (issue
//     createmeta) narrows locally with its own interface.
// ============================================================================

interface JiraArgv {
  // BaseArgv-compatible common flags
  profile?: string;
  output?: string;
  json?: boolean;
  page?: number;
  pagelen?: number;

  // --- common options (string). Non-optional at the type level: handlers
  //     only touch options they declared, and declared options either have
  //     a default or are demanded, so at runtime they are populated. The
  //     double cast via Partial<JiraArgv> at the handler entry relaxes
  //     the required-field overlap check. ---
  jql: string;
  assignee: string;
  status: string;
  summary: string;
  description: string;
  priority: string;
  type: string;
  body: string;
  query: string;
  name: string;
  email: string;
  lead: string;
  goal: string;
  state: string;
  key: string;
  // `id` is typed as string here even though ~10 subcommands declare the
  // positional as `{ type: 'number' }`. Those handlers coerce locally with
  // `Number(argv.id)` where the API client needs a numeric id.
  id: string;
  title: string;
  uri: string;
  url: string;
  issue: string;
  inward: string;
  outward: string;
  subject: string;
  message: string;
  started: string;

  // camelCase aliases for kebab-case options
  accountId: string;
  assigneeType: string;
  displayName: string;
  endDate: string;
  issueKey: string;
  moduleKey: string;
  releaseDate: string;
  searchKey: string;
  startDate: string;
  statusColor: string;
  textBody: string;
  timeSpent: string;

  // --- common options (boolean) ---
  all: boolean;
  archived: boolean;
  custom: boolean;
  dismissible: boolean;
  enabled: boolean;
  favourite: boolean;
  favourites: boolean;
  mine: boolean;
  released: boolean;
  deleteSubtasks: boolean;

  // --- common options (number) ---
  maxResults: number;
  startAt: number;
  boardId: number;
  projectId: number;

  // --- common options (array-of-string) ---
  labels: string[];
  components: string[];
  fields: string[];

  // --- positionals (kebab-case access). Always present at runtime when the
  //     handler is invoked, so these are non-optional despite yargs typing
  //     them as `| undefined` in its inferred types. ---
  'issue-key': string;
  'issue-keys': string[];
  'comment-id': string;
  'worklog-id': string;
  'link-id': string;
  'project-key': string;
  'board-id': number;
  'sprint-id': number;
  'epic-id-or-key': string;
  'filter-id': string;
  'account-id': string;
  'dashboard-id': string;
  'gadget-id': string;
  'group-name': string;
  'role-id': number;
  'screen-id': number;
  'tab-id': number;
  'task-id': string;
  'file-path': string;
  'id-or-name': string;

  // --- kebab-case option keys used without camelCase aliasing ---
  'max-results': number;
  'start-at': number;
  'text-body': string;
  'time-spent': string;
  'display-name': string;
  'project-id': number;
  'inward-issue': string;
  'outward-issue': string;
  'module-key': string;
  'status-color': string;
  'search-key': string;
  'release-date': string;
  'start-date': string;
  'end-date': string;
  'assignee-type': string;
  'delete-subtasks': boolean;

  // project — single key (most commands). The createmeta handler uses a
  // local variant that accepts string[].
  project: string;
}

// Variant for `issue createmeta` where --project is declared as an array.
interface JiraCreateMetaArgv extends Omit<JiraArgv, 'project'> {
  project?: string[];
}
import type {
  AnnouncementBanner,
  Dashboard,
  DashboardGadget,
  Field,
  Filter,
  IssueLinkType,
  IssueSecurityScheme,
  IssueType,
  IssueTypeScheme,
  NotificationScheme,
  PermissionScheme,
  Priority,
  Project,
  ProjectCategory,
  ProjectComponent,
  ProjectRole,
  Resolution,
  Screen,
  Sprint,
  UserPermission,
  Version,
  Worklog,
  WorkflowScheme,
  FieldConfiguration,
} from '../internal/jira/types.js';

// ============================================================================
// Helper: aligned table printer (same pattern as bitbucket.js)
// ============================================================================

function printTable(headers: string[], rows: string[][]): void {
  const allRows = [headers, ...rows];
  const colCount = headers.length;
  const widths: number[] = Array(colCount).fill(0);
  for (const row of allRows) {
    for (let i = 0; i < colCount; i++) {
      const cell = String(row[i] ?? '');
      if (cell.length > widths[i]) widths[i] = cell.length;
    }
  }
  for (const row of allRows) {
    const parts = row.map((cell, i) => {
      const s = String(cell ?? '');
      return i < colCount - 1 ? s.padEnd(widths[i] + 2) : s;
    });
    console.log(parts.join(''));
  }
}

// ============================================================================
// ADF body builder (mirrors the Go pattern)
// ============================================================================

function adfDoc(text: string): ADFDocument {
  return {
    type: 'doc',
    version: 1,
    content: [
      {
        type: 'paragraph',
        content: [{ type: 'text', text }],
      },
    ],
  };
}

// ============================================================================
// registerJiraCommands
// ============================================================================

export function registerJiraCommands(yargs: Argv): Argv {
  return yargs

    // =========================================================================
    // issue (alias: i)
    // =========================================================================
    .command(
      ['issue', 'i'],
      'Manage Jira issues',
      (y) =>
        y
          .demandCommand(1)

          // issue list
          .command(
            ['list', 'ls'],
            'List issues using JQL search',
            (y2) =>
              y2
                .option('jql', { type: 'string', description: 'JQL query string (overrides convenience flags)', default: '' })
                .option('project', { type: 'string', description: 'Filter by project key (uses profile default if not set)' })
                .option('assignee', { type: 'string', description: 'Filter by assignee', default: '' })
                .option('status', { type: 'string', description: 'Filter by status', default: '' })
                .option('max-results', { type: 'number', description: 'Maximum number of results per page', default: 50 })
                .option('start-at', { type: 'number', description: 'Index of the first result', default: 0 })
                .option('all', { type: 'boolean', description: 'Fetch all pages', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let jql = argv.jql;
              const project = defaultProject(argv);

              if (!jql) {
                const clauses = [];
                if (project) clauses.push(`project = ${project}`);
                if (argv.assignee) clauses.push(`assignee = "${argv.assignee}"`);
                if (argv.status) clauses.push(`status = "${argv.status}"`);
                jql = clauses.length === 0
                  ? 'created >= -30d order by created DESC'
                  : clauses.join(' AND ') + ' order by created DESC';
              }

              const fields = ['summary', 'issuetype', 'status', 'priority', 'assignee'];
              let results = await jiraProjects.searchJQL(client, jql, argv.startAt, argv.maxResults, fields, undefined);

              if (argv.all) {
                let allIssues = results.issues;
                while (allIssues.length < results.total) {
                  const next = await jiraProjects.searchJQL(client, jql, argv.startAt + allIssues.length, argv.maxResults, fields, undefined);
                  if (!next.issues || next.issues.length === 0) break;
                  allIssues = allIssues.concat(next.issues);
                }
                results.issues = allIssues;
              }

              if (isJSONOutput(argv)) { outputJSON(results); return; }

              const rows = (results.issues || []).map((issue) => [
                issue.key || '',
                issue.fields?.issuetype?.name || '',
                issue.fields?.status?.name || '',
                issue.fields?.priority?.name || '',
                issue.fields?.assignee?.displayName || '',
                issue.fields?.summary || '',
              ]);
              printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
              printPaginationHint((results.issues || []).length, results.total || 0);
            }
          )

          // issue get
          .command(
            'get <issue-key>',
            'Get issue details',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const issue = await jiraIssues.getIssue(client, argv['issue-key'], undefined, undefined);

              if (isJSONOutput(argv)) { outputJSON(issue); return; }

              const f = issue.fields || {};
              const description = renderADF(f.description);
              const labels = (f.labels || []).join(', ');
              const components = (f.components || []).map((c) => c.name).join(', ');
              const fixVersions = (f.fixVersions || []).map((v) => v.name).join(', ');

              console.log(`Key:          ${issue.key}`);
              console.log(`Summary:      ${f.summary || ''}`);
              console.log(`Status:       ${f.status?.name || ''}`);
              console.log(`Type:         ${f.issuetype?.name || ''}`);
              console.log(`Priority:     ${f.priority?.name || ''}`);
              console.log(`Assignee:     ${f.assignee?.displayName || ''}`);
              console.log(`Reporter:     ${f.reporter?.displayName || ''}`);
              console.log(`Created:      ${f.created || ''}`);
              console.log(`Updated:      ${f.updated || ''}`);
              console.log(`Labels:       ${labels}`);
              console.log(`Components:   ${components}`);
              console.log(`Fix Versions: ${fixVersions}`);
              console.log(`Description:\n${description}`);
            }
          )

          // issue create
          .command(
            'create',
            'Create a new issue',
            (y2) =>
              y2
                .option('project', { type: 'string', description: 'Project key (uses profile default if not set)' })
                .option('type', { type: 'string', description: 'Issue type', default: 'Task' })
                .option('summary', { type: 'string', description: 'Issue summary (required)' })
                .option('description', { type: 'string', description: 'Issue description', default: '' })
                .option('assignee', { type: 'string', description: 'Assignee account ID', default: '' })
                .option('priority', { type: 'string', description: 'Priority name', default: '' })
                .option('labels', { type: 'array', description: 'Labels', default: [] })
                .option('components', { type: 'array', description: 'Component names', default: [] })
                .demandOption(['summary']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = defaultProject(argv);
              if (!project) throw new Error("--project is required (or set a default with 'acli config set-defaults')");

              const fields: Record<string, JsonBody> = {
                project: { key: project },
                issuetype: { name: argv.type },
                summary: argv.summary,
              };

              if (argv.description) fields.description = adfDoc(argv.description);
              if (argv.assignee) fields.assignee = { accountId: argv.assignee };
              if (argv.priority) fields.priority = { name: argv.priority };
              if (argv.labels && argv.labels.length > 0) fields.labels = argv.labels;
              if (argv.components && argv.components.length > 0) {
                fields.components = argv.components.map((c: string) => ({ name: c }));
              }

              const created = await jiraIssues.createIssue(client, { fields });
              outputResult(argv, 'created', created.key, `Created issue: ${created.key}`, created);
            }
          )

          // issue edit (update)
          .command(
            'edit <issue-key>',
            'Edit an existing issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('summary', { type: 'string', description: 'Issue summary' })
                .option('description', { type: 'string', description: 'Issue description' })
                .option('assignee', { type: 'string', description: 'Assignee account ID' })
                .option('priority', { type: 'string', description: 'Priority name' })
                .option('labels', { type: 'array', description: 'Labels' })
                .option('components', { type: 'array', description: 'Component names' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const fields: Record<string, JsonBody> = {};

              if (argv.summary !== undefined) fields.summary = argv.summary;
              if (argv.description !== undefined) fields.description = adfDoc(argv.description);
              if (argv.assignee !== undefined) fields.assignee = { accountId: argv.assignee };
              if (argv.priority !== undefined) fields.priority = { name: argv.priority };
              if (argv.labels !== undefined) fields.labels = argv.labels;
              if (argv.components !== undefined) {
                fields.components = argv.components.map((c: string) => ({ name: c }));
              }

              if (Object.keys(fields).length === 0) throw new Error('no fields specified to update');

              await jiraIssues.editIssue(client, argv['issue-key'], { fields }, true);
              outputResult(argv, 'updated', argv['issue-key'], `Issue ${argv['issue-key']} updated successfully`, null);
            }
          )

          // issue delete
          .command(
            'delete <issue-key>',
            'Delete an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('delete-subtasks', { type: 'boolean', description: 'Also delete subtasks', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.deleteIssue(client, argv['issue-key'], argv.deleteSubtasks);
              outputResult(argv, 'deleted', argv['issue-key'], `Issue ${argv['issue-key']} deleted`, null);
            }
          )

          // issue assign
          .command(
            'assign <issue-key> <account-id>',
            "Assign an issue to a user (use '-1' or 'none' to unassign)",
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .positional('account-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let accountID = argv['account-id'];
              if (accountID === '-1' || accountID === 'none') accountID = '';

              await jiraIssues.assignIssue(client, argv['issue-key'], accountID);
              const msg = accountID
                ? `Issue ${argv['issue-key']} assigned to ${argv['account-id']}`
                : `Issue ${argv['issue-key']} unassigned`;
              outputResult(argv, 'assigned', argv['issue-key'], msg, null);
            }
          )

          // issue transition
          .command(
            'transition <issue-key>',
            'Transition an issue to a new status',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('id', { type: 'string', description: 'Transition ID', default: '' })
                .option('status', { type: 'string', description: 'Target status name (case-insensitive)', default: '' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let transitionID: string | undefined = argv.id;
              const statusName = argv.status;

              if (!transitionID && !statusName) throw new Error('either --id or --status must be provided');

              if (!transitionID) {
                const resp = await jiraIssues.getIssueTransitions(client, argv['issue-key']);
                for (const t of resp.transitions || []) {
                  if (
                    t.name?.toLowerCase() === statusName.toLowerCase() ||
                    (t.to && t.to.name?.toLowerCase() === statusName.toLowerCase())
                  ) {
                    transitionID = t.id;
                    break;
                  }
                }
                if (!transitionID) {
                  const available = (resp.transitions || []).map((t) =>
                    t.to && t.to.name !== t.name ? `${t.name} (-> ${t.to.name})` : t.name
                  );
                  throw new Error(`no transition found matching status "${statusName}"; available transitions: ${available.join(', ')}`);
                }
              }

              await jiraIssues.doIssueTransition(client, argv['issue-key'], { transition: { id: transitionID } });
              outputResult(argv, 'transitioned', argv['issue-key'], `Issue ${argv['issue-key']} transitioned successfully`, null);
            }
          )

          // issue transitions (list available)
          .command(
            'transitions <issue-key>',
            'List available transitions for an issue',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const resp = await jiraIssues.getIssueTransitions(client, argv['issue-key']);

              if (isJSONOutput(argv)) { outputJSON(resp); return; }

              const rows = (resp.transitions || []).map((t) => [
                t.id || '',
                t.name || '',
                t.to?.name || '',
              ]);
              printTable(['ID', 'NAME', 'TO STATUS'], rows);
            }
          )

          // issue comment (group)
          .command(
            'comment',
            'Manage issue comments',
            (y2) =>
              y2
                .demandCommand(1)

                .command(
                  ['list <issue-key>', 'ls <issue-key>'],
                  'List comments on an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .option('max-results', { type: 'number', default: 50 })
                      .option('start-at', { type: 'number', default: 0 })
                      .option('all', { type: 'boolean', default: false }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    let page = await jiraIssues.getIssueComments(client, argv['issue-key'], argv.startAt, argv.maxResults);
                    if (argv.all) {
                      let all = page.comments || [];
                      while (all.length < page.total) {
                        const next = await jiraIssues.getIssueComments(client, argv['issue-key'], argv.startAt + all.length, argv.maxResults);
                        if (!next.comments || next.comments.length === 0) break;
                        all = all.concat(next.comments);
                      }
                      page.comments = all;
                    }
                    if (isJSONOutput(argv)) { outputJSON(page); return; }

                    const rows = (page.comments || []).map((c) => {
                      const author = c.author?.displayName || '';
                      let body = renderADF(c.body);
                      if (body.length > 60) body = body.slice(0, 57) + '...';
                      return [c.id || '', author, c.created || '', body];
                    });
                    printTable(['ID', 'AUTHOR', 'CREATED', 'BODY'], rows);
                    printPaginationHint((page.comments || []).length, page.total || 0);
                  }
                )

                .command(
                  'add <issue-key>',
                  'Add a comment to an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .option('body', { type: 'string', description: 'Comment body text (required)' })
                      .demandOption(['body']),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const comment = await jiraIssues.addIssueComment(client, argv['issue-key'], adfDoc(argv.body), undefined);
                    outputResult(argv, 'created', comment.id!, `Comment ${comment.id} added to ${argv['issue-key']}`, comment);
                  }
                )

                .command(
                  'get <issue-key> <comment-id>',
                  'Get a specific comment',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .positional('comment-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const comment = await jiraIssues.getIssueComment(client, argv['issue-key'], argv['comment-id']);
                    if (isJSONOutput(argv)) { outputJSON(comment); return; }
                    const author = comment.author?.displayName || '';
                    console.log(`ID:      ${comment.id}`);
                    console.log(`Author:  ${author}`);
                    console.log(`Created: ${comment.created || ''}`);
                    console.log(`Body:\n${renderADF(comment.body)}`);
                  }
                )

                .command(
                  'delete <issue-key> <comment-id>',
                  'Delete a comment',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .positional('comment-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    await jiraIssues.deleteIssueComment(client, argv['issue-key'], argv['comment-id']);
                    outputResult(argv, 'deleted', argv['comment-id'], `Comment ${argv['comment-id']} deleted from ${argv['issue-key']}`, null);
                  }
                ),
            () => {}
          )

          // issue worklog (group)
          .command(
            'worklog',
            'Manage issue worklogs',
            (y2) =>
              y2
                .demandCommand(1)

                .command(
                  ['list <issue-key>', 'ls <issue-key>'],
                  'List worklogs on an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .option('max-results', { type: 'number', default: 50 })
                      .option('start-at', { type: 'number', default: 0 })
                      .option('all', { type: 'boolean', default: false }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    let page = await jiraIssues.getIssueWorklogs(client, argv['issue-key'], argv.startAt, argv.maxResults);
                    if (argv.all) {
                      let all = page.worklogs || [];
                      while (all.length < page.total) {
                        const next = await jiraIssues.getIssueWorklogs(client, argv['issue-key'], argv.startAt + all.length, argv.maxResults);
                        if (!next.worklogs || next.worklogs.length === 0) break;
                        all = all.concat(next.worklogs);
                      }
                      page.worklogs = all;
                    }
                    if (isJSONOutput(argv)) { outputJSON(page); return; }

                    const rows = (page.worklogs || []).map((wl) => [
                      wl.id || '',
                      wl.author?.displayName || '',
                      wl.timeSpent || '',
                      wl.started || '',
                    ]);
                    printTable(['ID', 'AUTHOR', 'TIME SPENT', 'STARTED'], rows);
                    printPaginationHint((page.worklogs || []).length, page.total || 0);
                  }
                )

                .command(
                  'add <issue-key>',
                  'Add a worklog entry to an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .option('time-spent', { type: 'string', description: "Time spent (e.g. '2h', '30m') (required)" })
                      .option('started', { type: 'string', description: 'Start time (ISO datetime)', default: '' })
                      .demandOption(['time-spent']),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const worklog: Partial<Worklog> = { timeSpent: argv['time-spent'] };
                    if (argv.started) worklog.started = argv.started;
                    const result = await jiraIssues.addIssueWorklog(client, argv['issue-key'], worklog);
                    outputResult(argv, 'created', result.id!, `Worklog ${result.id} added to ${argv['issue-key']}`, result);
                  }
                )

                .command(
                  'delete <issue-key> <worklog-id>',
                  'Delete a worklog entry',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .positional('worklog-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    await jiraIssues.deleteIssueWorklog(client, argv['issue-key'], argv['worklog-id']);
                    outputResult(argv, 'deleted', argv['worklog-id'], `Worklog ${argv['worklog-id']} deleted from ${argv['issue-key']}`, null);
                  }
                ),
            () => {}
          )

          // issue attach
          .command(
            'attach <issue-key> <file-path>',
            'Upload an attachment to an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .positional('file-path', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const attachments = await jiraIssues.addIssueAttachment(client, argv['issue-key'], argv['file-path']);
              if (isJSONOutput(argv)) { outputJSON(attachments); return; }
              for (const a of attachments || []) {
                console.log(`Attached: ${a.filename} (id: ${a.id})`);
              }
            }
          )

          // issue vote
          .command(
            'vote <issue-key>',
            'Add your vote to an issue',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.addIssueVote(client, argv['issue-key']);
              outputResult(argv, 'voted', argv['issue-key'], `Vote added to ${argv['issue-key']}`, null);
            }
          )

          // issue unvote
          .command(
            'unvote <issue-key>',
            'Remove your vote from an issue',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.removeIssueVote(client, argv['issue-key']);
              outputResult(argv, 'unvoted', argv['issue-key'], `Vote removed from ${argv['issue-key']}`, null);
            }
          )

          // issue watch
          .command(
            'watch <issue-key>',
            'Add a watcher to an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('account-id', { type: 'string', description: 'Account ID of user to add as watcher (default: self)', default: '' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.addIssueWatcher(client, argv['issue-key'], argv.accountId);
              const msg = argv.accountId
                ? `Watcher ${argv.accountId} added to ${argv['issue-key']}`
                : `You are now watching ${argv['issue-key']}`;
              outputResult(argv, 'watch_added', argv['issue-key'], msg, null);
            }
          )

          // issue unwatch
          .command(
            'unwatch <issue-key>',
            'Remove a watcher from an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('account-id', { type: 'string', description: 'Account ID of user to remove as watcher (required)' })
                .demandOption(['account-id']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.removeIssueWatcher(client, argv['issue-key'], argv.accountId);
              outputResult(argv, 'watch_removed', argv['issue-key'], `Watcher ${argv.accountId} removed from ${argv['issue-key']}`, null);
            }
          )

          // issue watchers
          .command(
            'watchers <issue-key>',
            'List watchers of an issue',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const watches = await jiraIssues.getIssueWatchers(client, argv['issue-key']);
              if (isJSONOutput(argv)) { outputJSON(watches); return; }

              const rows = (watches.watchers || []).map((w) => [w.accountId || '', w.displayName || '']);
              printTable(['ACCOUNT ID', 'DISPLAY NAME'], rows);
              console.log(`\nTotal watchers: ${watches.watchCount || 0}`);
            }
          )

          // issue changelog
          .command(
            'changelog <issue-key>',
            'List changelog for an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let page = await jiraIssues.getIssueChangelog(client, argv['issue-key'], argv.startAt, argv.maxResults);
              if (argv.all) {
                let allValues = page.values || [];
                while (allValues.length < page.total) {
                  const next = await jiraIssues.getIssueChangelog(client, argv['issue-key'], argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                }
                page.values = allValues;
              }
              if (isJSONOutput(argv)) { outputJSON(page); return; }

              const histories = page.values && page.values.length > 0 ? page.values : (page.histories || []);
              const rows: string[][] = [];
              for (const entry of histories) {
                const author = entry.author?.displayName || '';
                for (const item of entry.items || []) {
                  rows.push([entry.created || '', author, item.field || '', item.fromString || '', item.toString || '']);
                }
              }
              printTable(['DATE', 'AUTHOR', 'FIELD', 'FROM', 'TO'], rows);
              printPaginationHint(histories.length, page.total || 0);
            }
          )

          // issue link (group — remote links)
          .command(
            'link',
            'Manage remote links on issues',
            (y2) =>
              y2
                .demandCommand(1)

                .command(
                  ['list <issue-key>', 'ls <issue-key>'],
                  'List remote links on an issue',
                  (y3) => y3.positional('issue-key', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const links = await jiraIssues.getIssueRemoteLinks(client, argv['issue-key']);
                    if (isJSONOutput(argv)) { outputJSON(links); return; }

                    const rows = (links || []).map((l) => [
                      String(l.id || ''),
                      l.object?.title || '',
                      l.object?.url || '',
                    ]);
                    printTable(['ID', 'TITLE', 'URL'], rows);
                  }
                )

                .command(
                  'create <issue-key>',
                  'Create a remote link on an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .option('url', { type: 'string', description: 'URL of the remote link', default: '' })
                      .option('title', { type: 'string', description: 'Title of the remote link', default: '' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const created = await jiraIssues.createIssueRemoteLink(client, argv['issue-key'], {
                      object: { url: argv.url, title: argv.title },
                    });
                    outputResult(argv, 'created', String(created.id), `Remote link ${created.id} created on ${argv['issue-key']}`, created);
                  }
                )

                .command(
                  'delete <issue-key> <link-id>',
                  'Delete a remote link from an issue',
                  (y3) =>
                    y3
                      .positional('issue-key', { type: 'string' })
                      .positional('link-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    await jiraIssues.deleteIssueRemoteLink(client, argv['issue-key'], argv['link-id']);
                    outputResult(argv, 'deleted', argv['link-id'], `Remote link ${argv['link-id']} deleted from ${argv['issue-key']}`, null);
                  }
                ),
            () => {}
          )

          // issue notify
          .command(
            'notify <issue-key>',
            'Send a notification for an issue',
            (y2) =>
              y2
                .positional('issue-key', { type: 'string' })
                .option('subject', { type: 'string', description: 'Notification subject', default: '' })
                .option('text-body', { type: 'string', description: 'Notification text body', default: '' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraIssues.notifyIssue(client, argv['issue-key'], { subject: argv.subject, textBody: argv['text-body'] });
              outputResult(argv, 'notified', argv['issue-key'], `Notification sent for ${argv['issue-key']}`, null);
            }
          )

          // issue createmeta
          .command(
            'createmeta',
            'Get issue create metadata',
            (y2) => y2.option('project', { type: 'array', description: 'Project keys to filter', default: [] }),
            async (_a) => {
              const argv = _a as Partial<JiraCreateMetaArgv> as JiraCreateMetaArgv;
              // getJiraClient only reads BaseArgv fields (profile/output/json);
              // createmeta's project is an array, not the single-project
              // string BaseArgv/JiraArgv use, so we hand it a trimmed view.
              const client = getJiraClient({ profile: argv.profile, output: argv.output, json: argv.json });
              const projects: string[] = argv.project || [];
              const meta = await jiraIssues.getCreateMeta(client, projects, ['projects.issuetypes.fields']);
              outputJSON(meta);
            }
          )

          // issue editmeta
          .command(
            'editmeta <issue-key>',
            'Get issue edit metadata',
            (y2) => y2.positional('issue-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const meta = await jiraIssues.getIssueEditMeta(client, argv['issue-key']);
              outputJSON(meta);
            }
          ),
      () => {}
    )

    // =========================================================================
    // project (alias: p)
    // =========================================================================
    .command(
      ['project', 'p'],
      'Manage projects',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List or search projects',
            (y2) =>
              y2
                .option('query', { type: 'string', description: 'Search query to filter projects', default: '' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraProjects.searchProjects(client, argv.query, argv.startAt, argv.maxResults, '');

              if (argv.all) {
                while (!result.isLast && (result.values || []).length < result.total) {
                  const next = await jiraProjects.searchProjects(client, argv.query, argv.startAt + result.values.length, argv.maxResults, '');
                  if (!next.values || next.values.length === 0) break;
                  result.values = result.values.concat(next.values);
                  result.isLast = next.isLast;
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.values || []).map((p) => [
                p.key || '',
                p.name || '',
                p.projectTypeKey || '',
                p.lead?.displayName || '',
              ]);
              printTable(['KEY', 'NAME', 'TYPE', 'LEAD'], rows);
              printPaginationHint((result.values || []).length, result.total || 0);
            }
          )

          .command(
            'get <project-key>',
            'Get project details',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = await jiraProjects.getProject(client, argv['project-key'], '');
              if (isJSONOutput(argv)) { outputJSON(project); return; }

              const lead = project.lead?.displayName || 'N/A';
              const category = project.projectCategory?.name || 'N/A';
              console.log(`Key:          ${project.key}`);
              console.log(`Name:         ${project.name}`);
              console.log(`ID:           ${project.id}`);
              console.log(`Type:         ${project.projectTypeKey}`);
              console.log(`Lead:         ${lead}`);
              console.log(`Description:  ${project.description || ''}`);
              console.log(`Category:     ${category}`);
              console.log(`Style:        ${project.style || ''}`);
              console.log(`Simplified:   ${project.simplified}`);
              console.log(`Archived:     ${project.archived}`);
              console.log(`URL:          ${project.url || project.self || ''}`);
            }
          )

          .command(
            'create',
            'Create a project',
            (y2) =>
              y2
                .option('key', { type: 'string', description: 'Project key (required)' })
                .option('name', { type: 'string', description: 'Project name (required)' })
                .option('type', { type: 'string', description: 'Project type key', default: 'software' })
                .option('lead', { type: 'string', description: 'Lead account ID', default: '' })
                .option('description', { type: 'string', description: 'Project description', default: '' })
                .demandOption(['key', 'name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Project> & { leadAccountId?: string } = { key: argv.key, name: argv.name, projectTypeKey: argv.type };
              if (argv.lead) body.leadAccountId = argv.lead;
              if (argv.description) body.description = argv.description;

              const project = await jiraProjects.createProject(client, body);
              outputResult(argv, 'created', project.key, `Created project ${project.key} (ID: ${project.id})`, project);
            }
          )

          .command(
            'update <project-key>',
            'Update a project',
            (y2) =>
              y2
                .positional('project-key', { type: 'string' })
                .option('name', { type: 'string', description: 'Project name' })
                .option('description', { type: 'string', description: 'Project description' })
                .option('lead', { type: 'string', description: 'Lead account ID' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Project> & { leadAccountId?: string } = {};
              if (argv.name !== undefined) body.name = argv.name;
              if (argv.description !== undefined) body.description = argv.description;
              if (argv.lead !== undefined) body.leadAccountId = argv.lead;

              if (Object.keys(body).length === 0) throw new Error('at least one flag (--name, --description, --lead) must be provided');

              const project = await jiraProjects.updateProject(client, argv['project-key'], body);
              outputResult(argv, 'updated', project.key!, `Updated project ${project.key}`, project);
            }
          )

          .command(
            'delete <project-key>',
            'Delete a project',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraProjects.deleteProject(client, argv['project-key']);
              outputResult(argv, 'deleted', argv['project-key'], `Deleted project ${argv['project-key']}`, null);
            }
          )

          .command(
            'components <project-key>',
            'List project components',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const components = await jiraProjects.getProjectComponents(client, argv['project-key']);
              const rows = (components || []).map((c) => [
                c.id || '',
                c.name || '',
                c.lead?.displayName || '',
                c.assigneeType || '',
              ]);
              printTable(['ID', 'NAME', 'LEAD', 'ASSIGNEE TYPE'], rows);
            }
          )

          .command(
            'versions <project-key>',
            'List project versions',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const versions = await jiraProjects.getProjectVersions(client, argv['project-key']);
              const rows = (versions || []).map((v) => {
                let status = 'Unreleased';
                if (v.released) status = 'Released';
                else if (v.archived) status = 'Archived';
                return [v.id || '', v.name || '', status, v.releaseDate || ''];
              });
              printTable(['ID', 'NAME', 'STATUS', 'RELEASE DATE'], rows);
            }
          )

          .command(
            'statuses <project-key>',
            'List project statuses by issue type',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const issueTypes = await jiraProjects.getProjectStatuses(client, argv['project-key']);
              const rows: string[][] = [];
              for (const it of issueTypes || []) {
                for (const s of it.statuses || []) {
                  rows.push([it.name || '', s.name || '', s.statusCategory?.name || '']);
                }
              }
              printTable(['ISSUE TYPE', 'STATUS NAME', 'CATEGORY'], rows);
            }
          )

          .command(
            'roles <project-key>',
            'List project roles',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const roles = await jiraProjects.getProjectRoles(client, argv['project-key']);
              const rows = Object.entries(roles || {}).map(([name, url]) => [name, url as string]);
              printTable(['ROLE NAME', 'URL'], rows);
            }
          )

          .command(
            'archive <project-key>',
            'Archive a project',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraProjects.archiveProject(client, argv['project-key']);
              outputResult(argv, 'archived', argv['project-key'], `Archived project ${argv['project-key']}`, null);
            }
          )

          .command(
            'restore <project-key>',
            'Restore an archived project',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraProjects.restoreProject(client, argv['project-key']);
              outputResult(argv, 'restored', argv['project-key'], `Restored project ${argv['project-key']}`, null);
            }
          )

          .command(
            'features <project-key>',
            'List project features',
            (y2) => y2.positional('project-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const resp = await jiraProjects.getProjectFeatures(client, argv['project-key']);
              const rows = (resp.features || []).map((f) => [
                f.localisedName || f.feature || '',
                f.state || '',
              ]);
              printTable(['FEATURE', 'STATE'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // board (alias: b)
    // =========================================================================
    .command(
      ['board', 'b'],
      'Manage boards (Agile API)',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List boards',
            (y2) =>
              y2
                .option('project', { type: 'string', description: 'Filter by project key or ID (uses profile default if not set)' })
                .option('type', { type: 'string', description: 'Filter by board type (scrum, kanban)', default: '' })
                .option('name', { type: 'string', description: 'Filter by board name', default: '' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = defaultProject(argv);
              let result = await jiraAgile.getBoards(client, argv.startAt, argv.maxResults, project, argv.type, argv.name);

              if (argv.all) {
                while (!result.isLast && (result.values || []).length < result.total) {
                  const next = await jiraAgile.getBoards(client, argv.startAt + result.values.length, argv.maxResults, project, argv.type, argv.name);
                  if (!next.values || next.values.length === 0) break;
                  result.values = result.values.concat(next.values);
                  result.isLast = next.isLast;
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.values || []).map((b) => [
                String(b.id || ''),
                b.name || '',
                b.type || '',
                b.location?.projectKey || '',
              ]);
              printTable(['ID', 'NAME', 'TYPE', 'PROJECT'], rows);
              printPaginationHint((result.values || []).length, result.total || 0);
            }
          )

          .command(
            'get <board-id>',
            'Get board details',
            (y2) => y2.positional('board-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const board = await jiraAgile.getBoard(client, argv['board-id']);
              if (isJSONOutput(argv)) { outputJSON(board); return; }

              console.log(`ID:      ${board.id}`);
              console.log(`Name:    ${board.name}`);
              console.log(`Type:    ${board.type}`);
              if (board.location) {
                console.log(`Project: ${board.location.projectName} (${board.location.projectKey})`);
              }
            }
          )

          .command(
            ['config <board-id>', 'configuration <board-id>'],
            'Get board configuration',
            (y2) => y2.positional('board-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const config = await jiraAgile.getBoardConfiguration(client, argv['board-id']);
              outputJSON(config);
            }
          )

          .command(
            'issues <board-id>',
            'List issues on a board',
            (y2) =>
              y2
                .positional('board-id', { type: 'number' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('jql', { type: 'string', description: 'JQL filter', default: '' })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getBoardIssues(client, argv['board-id'], argv.startAt, argv.maxResults, argv.jql);

              if (argv.all) {
                while ((result.issues || []).length < result.total) {
                  const next = await jiraAgile.getBoardIssues(client, argv['board-id'], argv.startAt + result.issues.length, argv.maxResults, argv.jql);
                  if (!next.issues || next.issues.length === 0) break;
                  result.issues = result.issues.concat(next.issues);
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.issues || []).map((issue) => [
                issue.key || '',
                issue.fields?.issuetype?.name || '',
                issue.fields?.status?.name || '',
                issue.fields?.priority?.name || '',
                issue.fields?.assignee?.displayName || '',
                issue.fields?.summary || '',
              ]);
              printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
              printPaginationHint((result.issues || []).length, result.total || 0);
            }
          )

          .command(
            'backlog <board-id>',
            'List backlog issues for a board',
            (y2) =>
              y2
                .positional('board-id', { type: 'number' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('jql', { type: 'string', description: 'JQL filter', default: '' })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getBoardBacklog(client, argv['board-id'], argv.startAt, argv.maxResults, argv.jql);

              if (argv.all) {
                while ((result.issues || []).length < result.total) {
                  const next = await jiraAgile.getBoardBacklog(client, argv['board-id'], argv.startAt + result.issues.length, argv.maxResults, argv.jql);
                  if (!next.issues || next.issues.length === 0) break;
                  result.issues = result.issues.concat(next.issues);
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.issues || []).map((issue) => [
                issue.key || '',
                issue.fields?.issuetype?.name || '',
                issue.fields?.status?.name || '',
                issue.fields?.priority?.name || '',
                issue.fields?.assignee?.displayName || '',
                issue.fields?.summary || '',
              ]);
              printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
              printPaginationHint((result.issues || []).length, result.total || 0);
            }
          )

          .command(
            'sprints <board-id>',
            'List sprints for a board',
            (y2) =>
              y2
                .positional('board-id', { type: 'number' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('state', { type: 'string', description: 'Filter by state (active, closed, future)', default: '' })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getBoardSprints(client, argv['board-id'], argv.startAt, argv.maxResults, argv.state);

              if (argv.all) {
                while (!result.isLast && (result.values || []).length < result.total) {
                  const next = await jiraAgile.getBoardSprints(client, argv['board-id'], argv.startAt + result.values.length, argv.maxResults, argv.state);
                  if (!next.values || next.values.length === 0) break;
                  result.values = result.values.concat(next.values);
                  result.isLast = next.isLast;
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.values || []).map((s) => [
                String(s.id || ''),
                s.name || '',
                s.state || '',
                s.startDate || '',
                s.endDate || '',
              ]);
              printTable(['ID', 'NAME', 'STATE', 'START DATE', 'END DATE'], rows);
              printPaginationHint((result.values || []).length, result.total || 0);
            }
          )

          .command(
            'epics <board-id>',
            'List epics for a board',
            (y2) =>
              y2
                .positional('board-id', { type: 'number' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getBoardEpics(client, argv['board-id'], argv.startAt, argv.maxResults);

              if (argv.all) {
                while (!result.isLast && (result.values || []).length < result.total) {
                  const next = await jiraAgile.getBoardEpics(client, argv['board-id'], argv.startAt + result.values.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  result.values = result.values.concat(next.values);
                  result.isLast = next.isLast;
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.values || []).map((e) => [
                String(e.id || ''),
                e.key || '',
                e.name || '',
                String(e.done),
                e.summary || '',
              ]);
              printTable(['ID', 'KEY', 'NAME', 'DONE', 'SUMMARY'], rows);
              printPaginationHint((result.values || []).length, result.total || 0);
            }
          ),
      () => {}
    )

    // =========================================================================
    // sprint (alias: sp)
    // =========================================================================
    .command(
      ['sprint', 'sp'],
      'Manage sprints (Agile API)',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <sprint-id>',
            'Get sprint details',
            (y2) => y2.positional('sprint-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const sprint = await jiraAgile.getSprint(client, argv['sprint-id']);
              if (isJSONOutput(argv)) { outputJSON(sprint); return; }

              console.log(`ID:         ${sprint.id}`);
              console.log(`Name:       ${sprint.name}`);
              console.log(`State:      ${sprint.state}`);
              console.log(`Start Date: ${sprint.startDate || ''}`);
              console.log(`End Date:   ${sprint.endDate || ''}`);
              if (sprint.completeDate) console.log(`Completed:  ${sprint.completeDate}`);
              console.log(`Board ID:   ${sprint.originBoardId}`);
              if (sprint.goal) console.log(`Goal:       ${sprint.goal}`);
            }
          )

          .command(
            'create',
            'Create a sprint',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Sprint name (required)' })
                .option('board-id', { type: 'number', description: 'Origin board ID (required)' })
                .option('start-date', { type: 'string', description: 'Start date (ISO 8601)', default: '' })
                .option('end-date', { type: 'string', description: 'End date (ISO 8601)', default: '' })
                .option('goal', { type: 'string', description: 'Sprint goal', default: '' })
                .demandOption(['name', 'board-id']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Sprint> = { name: argv.name, originBoardId: argv.boardId };
              if (argv.startDate) body.startDate = argv.startDate;
              if (argv.endDate) body.endDate = argv.endDate;
              if (argv.goal) body.goal = argv.goal;

              const sprint = await jiraAgile.createSprint(client, body);
              console.log(`Sprint created: ${sprint.name} (ID: ${sprint.id})`);
            }
          )

          .command(
            'update <sprint-id>',
            'Update a sprint',
            (y2) =>
              y2
                .positional('sprint-id', { type: 'number' })
                .option('name', { type: 'string', description: 'Sprint name' })
                .option('state', { type: 'string', description: 'Sprint state (active, closed, future)' })
                .option('start-date', { type: 'string', description: 'Start date (ISO 8601)' })
                .option('end-date', { type: 'string', description: 'End date (ISO 8601)' })
                .option('goal', { type: 'string', description: 'Sprint goal' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Sprint> = {};
              if (argv.name !== undefined) body.name = argv.name;
              if (argv.state !== undefined) body.state = argv.state;
              if (argv.startDate !== undefined) body.startDate = argv.startDate;
              if (argv.endDate !== undefined) body.endDate = argv.endDate;
              if (argv.goal !== undefined) body.goal = argv.goal;

              const sprint = await jiraAgile.updateSprint(client, argv['sprint-id'], body);
              console.log(`Sprint updated: ${sprint.name} (ID: ${sprint.id})`);
            }
          )

          .command(
            'delete <sprint-id>',
            'Delete a sprint',
            (y2) => y2.positional('sprint-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAgile.deleteSprint(client, argv['sprint-id']);
              console.log(`Sprint ${argv['sprint-id']} deleted.`);
            }
          )

          .command(
            'issues <sprint-id>',
            'List issues in a sprint',
            (y2) =>
              y2
                .positional('sprint-id', { type: 'number' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('jql', { type: 'string', description: 'JQL filter', default: '' })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getSprintIssues(client, argv['sprint-id'], argv.startAt, argv.maxResults, argv.jql);

              if (argv.all) {
                while ((result.issues || []).length < result.total) {
                  const next = await jiraAgile.getSprintIssues(client, argv['sprint-id'], argv.startAt + result.issues.length, argv.maxResults, argv.jql);
                  if (!next.issues || next.issues.length === 0) break;
                  result.issues = result.issues.concat(next.issues);
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.issues || []).map((issue) => [
                issue.key || '',
                issue.fields?.issuetype?.name || '',
                issue.fields?.status?.name || '',
                issue.fields?.priority?.name || '',
                issue.fields?.assignee?.displayName || '',
                issue.fields?.summary || '',
              ]);
              printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
              printPaginationHint((result.issues || []).length, result.total || 0);
            }
          )

          .command(
            'move <sprint-id> [issue-keys..]',
            'Move issues to a sprint',
            (y2) =>
              y2
                .positional('sprint-id', { type: 'number' })
                .positional('issue-keys', { type: 'string', array: true }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const issueKeys = argv['issue-keys'] || [];
              if (issueKeys.length === 0) throw new Error('at least one issue key is required');
              await jiraAgile.moveIssuesToSprint(client, argv['sprint-id'], issueKeys);
              console.log(`Moved ${issueKeys.join(', ')} to sprint ${argv['sprint-id']}`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // epic (alias: e)
    // =========================================================================
    .command(
      ['epic', 'e'],
      'Manage epics (Agile API)',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <epic-id-or-key>',
            'Get epic details',
            (y2) => y2.positional('epic-id-or-key', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const epic = await jiraAgile.getEpic(client, argv['epic-id-or-key']);
              if (isJSONOutput(argv)) { outputJSON(epic); return; }

              console.log(`ID:      ${epic.id}`);
              console.log(`Key:     ${epic.key}`);
              console.log(`Name:    ${epic.name}`);
              console.log(`Summary: ${epic.summary}`);
              console.log(`Done:    ${epic.done}`);
            }
          )

          .command(
            'issues <epic-id-or-key>',
            'List issues in an epic',
            (y2) =>
              y2
                .positional('epic-id-or-key', { type: 'string' })
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('jql', { type: 'string', description: 'JQL filter', default: '' })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraAgile.getEpicIssues(client, argv['epic-id-or-key'], argv.startAt, argv.maxResults, argv.jql);

              if (argv.all) {
                while ((result.issues || []).length < result.total) {
                  const next = await jiraAgile.getEpicIssues(client, argv['epic-id-or-key'], argv.startAt + result.issues.length, argv.maxResults, argv.jql);
                  if (!next.issues || next.issues.length === 0) break;
                  result.issues = result.issues.concat(next.issues);
                }
              }

              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.issues || []).map((issue) => [
                issue.key || '',
                issue.fields?.issuetype?.name || '',
                issue.fields?.status?.name || '',
                issue.fields?.priority?.name || '',
                issue.fields?.assignee?.displayName || '',
                issue.fields?.summary || '',
              ]);
              printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
              printPaginationHint((result.issues || []).length, result.total || 0);
            }
          )

          .command(
            'move <epic-id-or-key> [issue-keys..]',
            'Move issues to an epic',
            (y2) =>
              y2
                .positional('epic-id-or-key', { type: 'string' })
                .positional('issue-keys', { type: 'string', array: true }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const issueKeys = argv['issue-keys'] || [];
              if (issueKeys.length === 0) throw new Error('at least one issue key is required');
              await jiraAgile.moveIssuesToEpic(client, argv['epic-id-or-key'], issueKeys);
              console.log(`Moved ${issueKeys.join(', ')} to epic ${argv['epic-id-or-key']}`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // backlog
    // =========================================================================
    .command(
      'backlog <issue-keys..>',
      'Move issues to the backlog',
      (y) => y.positional('issue-keys', { type: 'string', array: true }),
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        const issueKeys = argv['issue-keys'] || [];
        if (issueKeys.length === 0) throw new Error('at least one issue key is required');
        await jiraAgile.moveIssuesToBacklog(client, issueKeys);
        console.log(`Moved ${issueKeys.join(', ')} to backlog`);
      }
    )

    // =========================================================================
    // search (alias: s)
    // =========================================================================
    .command(
      ['search', 's'],
      'Search issues using JQL',
      (y) =>
        y
          .option('jql', { type: 'string', description: 'JQL query (required)' })
          .option('max-results', { type: 'number', default: 50 })
          .option('start-at', { type: 'number', default: 0 })
          .option('fields', { type: 'array', description: 'Fields to return', default: [] })
          .option('all', { type: 'boolean', default: false })
          .demandOption(['jql']),
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        let results = await jiraProjects.searchJQL(client, argv.jql, argv.startAt, argv.maxResults, argv.fields || [], undefined);

        if (argv.all) {
          while ((results.issues || []).length < results.total) {
            const next = await jiraProjects.searchJQL(client, argv.jql, argv.startAt + results.issues.length, argv.maxResults, argv.fields || [], undefined);
            if (!next.issues || next.issues.length === 0) break;
            results.issues = results.issues.concat(next.issues);
          }
        }

        if (isJSONOutput(argv)) { outputJSON(results); return; }

        const rows = (results.issues || []).map((issue) => [
          issue.key || '',
          issue.fields?.issuetype?.name || '',
          issue.fields?.status?.name || '',
          issue.fields?.priority?.name || '',
          issue.fields?.assignee?.displayName || '',
          issue.fields?.summary || '',
        ]);
        printTable(['KEY', 'TYPE', 'STATUS', 'PRIORITY', 'ASSIGNEE', 'SUMMARY'], rows);
        printPaginationHint((results.issues || []).length, results.total || 0);
      }
    )

    // =========================================================================
    // filter (alias: fl)
    // =========================================================================
    .command(
      ['filter', 'f'],
      'Manage filters',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List or search filters',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Filter by name', default: '' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('favourites', { type: 'boolean', description: 'Show favourite filters', default: false })
                .option('mine', { type: 'boolean', description: 'Show my filters', default: false })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let filters: Filter[] = [];
              let total = 0;

              if (argv.favourites) {
                filters = await jiraProjects.getFavouriteFilters(client);
                total = filters.length;
              } else if (argv.mine) {
                filters = await jiraProjects.getMyFilters(client);
                total = filters.length;
              } else {
                let page = await jiraProjects.searchFilters(client, argv.name, argv.startAt, argv.maxResults);
                filters = page.values || [];
                total = page.total || 0;
                if (argv.all) {
                  while (!page.isLast && filters.length < page.total) {
                    const next = await jiraProjects.searchFilters(client, argv.name, argv.startAt + filters.length, argv.maxResults);
                    if (!next.values || next.values.length === 0) break;
                    filters = filters.concat(next.values);
                    page.isLast = next.isLast;
                  }
                }
              }

              if (isJSONOutput(argv)) { outputJSON(filters); return; }

              const rows = filters.map((f) => [
                f.id || '',
                f.name || '',
                f.owner?.displayName || '',
                truncate(f.jql || '', 60),
              ]);
              printTable(['ID', 'NAME', 'OWNER', 'JQL'], rows);
              printPaginationHint(filters.length, total);
            }
          )

          .command(
            'get <filter-id>',
            'Get filter details',
            (y2) => y2.positional('filter-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const filter = await jiraProjects.getFilter(client, argv['filter-id']);
              if (isJSONOutput(argv)) { outputJSON(filter); return; }

              console.log(`ID:          ${filter.id}`);
              console.log(`Name:        ${filter.name}`);
              console.log(`Description: ${filter.description || ''}`);
              if (filter.owner) console.log(`Owner:       ${filter.owner.displayName}`);
              console.log(`JQL:         ${filter.jql || ''}`);
              console.log(`Favourite:   ${filter.favourite}`);
              console.log(`View URL:    ${filter.viewUrl || ''}`);
            }
          )

          .command(
            'create',
            'Create a filter',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Filter name (required)' })
                .option('jql', { type: 'string', description: 'JQL query (required)' })
                .option('description', { type: 'string', description: 'Filter description', default: '' })
                .option('favourite', { type: 'boolean', description: 'Mark as favourite', default: false })
                .demandOption(['name', 'jql']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const created = await jiraProjects.createFilter(client, {
                name: argv.name,
                jql: argv.jql,
                description: argv.description,
                favourite: argv.favourite,
              });
              outputResult(argv, 'created', created.id!, `Filter created: ${created.name} (ID: ${created.id})`, created);
            }
          )

          .command(
            'update <filter-id>',
            'Update a filter',
            (y2) =>
              y2
                .positional('filter-id', { type: 'string' })
                .option('name', { type: 'string', description: 'Filter name' })
                .option('jql', { type: 'string', description: 'JQL query' })
                .option('description', { type: 'string', description: 'Filter description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const filter: Partial<Filter> = {};
              if (argv.name !== undefined) filter.name = argv.name;
              if (argv.jql !== undefined) filter.jql = argv.jql;
              if (argv.description !== undefined) filter.description = argv.description;

              const updated = await jiraProjects.updateFilter(client, argv['filter-id'], filter);
              outputResult(argv, 'updated', updated.id!, `Filter updated: ${updated.name} (ID: ${updated.id})`, updated);
            }
          )

          .command(
            'delete <filter-id>',
            'Delete a filter',
            (y2) => y2.positional('filter-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraProjects.deleteFilter(client, argv['filter-id']);
              outputResult(argv, 'deleted', argv['filter-id'], `Filter ${argv['filter-id']} deleted`, null);
            }
          ),
      () => {}
    )

    // =========================================================================
    // user (alias: u)
    // =========================================================================
    .command(
      ['user', 'u'],
      'Manage users',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <account-id>',
            'Get user details',
            (y2) => y2.positional('account-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const user = await jiraAdmin.getUser(client, argv['account-id']);
              if (isJSONOutput(argv)) { outputJSON(user); return; }

              console.log(`Account ID:   ${user.accountId}`);
              console.log(`Display Name: ${user.displayName}`);
              console.log(`Email:        ${user.emailAddress || ''}`);
              console.log(`Active:       ${user.active}`);
              console.log(`Time Zone:    ${user.timeZone || ''}`);
              console.log(`Account Type: ${user.accountType || ''}`);
            }
          )

          .command(
            'search',
            'Search users',
            (y2) =>
              y2
                .option('query', { type: 'string', description: 'Search query (required)' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false })
                .demandOption(['query']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let users = await jiraAdmin.findUsers(client, argv.query, argv.startAt, argv.maxResults);
              if (argv.all) {
                while (users.length >= argv.maxResults) {
                  const next = await jiraAdmin.findUsers(client, argv.query, argv.startAt + users.length, argv.maxResults);
                  if (!next || next.length === 0) break;
                  users = users.concat(next);
                }
              }
              if (isJSONOutput(argv)) { outputJSON(users); return; }

              const rows = users.map((u) => [u.accountId || '', u.displayName || '', u.emailAddress || '', String(u.active)]);
              printTable(['ACCOUNT_ID', 'DISPLAY_NAME', 'EMAIL', 'ACTIVE'], rows);
            }
          )

          .command(
            'assignable',
            'Find assignable users',
            (y2) =>
              y2
                .option('query', { type: 'string', description: 'Search query', default: '' })
                .option('project', { type: 'string', description: 'Project key (uses profile default if not set)' })
                .option('issue-key', { type: 'string', description: 'Issue key', default: '' })
                .option('max-results', { type: 'number', default: 50 }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = defaultProject(argv);
              const users = await jiraAdmin.findUsersAssignable(client, argv.query, project, argv.issueKey, 0, argv.maxResults);
              const rows = (users || []).map((u) => [u.accountId || '', u.displayName || '', u.emailAddress || '', String(u.active)]);
              printTable(['ACCOUNT_ID', 'DISPLAY_NAME', 'EMAIL', 'ACTIVE'], rows);
            }
          )

          .command(
            'me',
            'Show current user',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const user = await jiraAdmin.getCurrentUser(client);
              console.log(`Account ID:   ${user.accountId}`);
              console.log(`Display Name: ${user.displayName}`);
              console.log(`Email:        ${user.emailAddress || ''}`);
              console.log(`Active:       ${user.active}`);
              console.log(`Time Zone:    ${user.timeZone || ''}`);
            }
          )

          .command(
            ['list', 'ls'],
            'List all users',
            (y2) =>
              y2
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let users = await jiraAdmin.getAllUsers(client, argv.startAt, argv.maxResults);
              if (argv.all) {
                while (users.length >= argv.maxResults) {
                  const next = await jiraAdmin.getAllUsers(client, argv.startAt + users.length, argv.maxResults);
                  if (!next || next.length === 0) break;
                  users = users.concat(next);
                }
              }
              const rows = (users || []).map((u) => [u.accountId || '', u.displayName || '', u.emailAddress || '', String(u.active)]);
              printTable(['ACCOUNT_ID', 'DISPLAY_NAME', 'EMAIL', 'ACTIVE'], rows);
            }
          )

          .command(
            'create',
            'Create a user',
            (y2) =>
              y2
                .option('email', { type: 'string', description: 'User email address (required)' })
                .option('display-name', { type: 'string', description: 'User display name (required)' })
                .demandOption(['email', 'display-name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const user = await jiraAdmin.createUser(client, { emailAddress: argv.email, displayName: argv['display-name'] });
              console.log(`User created: ${user.displayName} (Account ID: ${user.accountId})`);
            }
          )

          .command(
            'delete <account-id>',
            'Delete a user',
            (y2) => y2.positional('account-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteUser(client, argv['account-id']);
              console.log(`User ${argv['account-id']} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // group (alias: g)
    // =========================================================================
    .command(
      ['group', 'g'],
      'Manage groups',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List groups',
            (y2) =>
              y2
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let page = await jiraAdmin.getBulkGroups(client, argv.startAt, argv.maxResults);
              let allGroups = page.values || [];
              if (argv.all) {
                while (!page.isLast && allGroups.length < page.total) {
                  const next = await jiraAdmin.getBulkGroups(client, argv.startAt + allGroups.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allGroups = allGroups.concat(next.values);
                  page = next;
                }
              }
              if (isJSONOutput(argv)) { outputJSON(allGroups); return; }

              const rows = allGroups.map((g) => [g.groupId || '', g.name || '']);
              printTable(['GROUP_ID', 'NAME'], rows);
            }
          )

          .command(
            'get <group-name>',
            'Get group details',
            (y2) => y2.positional('group-name', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const group = await jiraAdmin.getGroup(client, argv['group-name']);
              if (isJSONOutput(argv)) { outputJSON(group); return; }

              console.log(`Group ID: ${group.groupId}`);
              console.log(`Name:     ${group.name}`);
            }
          )

          .command(
            'create',
            'Create a group',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Group name (required)' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const group = await jiraAdmin.createGroup(client, argv.name);
              console.log(`Group created: ${group.name} (ID: ${group.groupId})`);
            }
          )

          .command(
            'delete <group-name>',
            'Delete a group',
            (y2) => y2.positional('group-name', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteGroup(client, argv['group-name']);
              console.log(`Group ${argv['group-name']} deleted`);
            }
          )

          .command(
            'members <group-name>',
            'List group members',
            (y2) =>
              y2
                .positional('group-name', { type: 'string' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let members = await jiraAdmin.getGroupMembers(client, argv['group-name'], argv.startAt, argv.maxResults);
              let allMembers = members.values || [];
              if (argv.all) {
                while (!members.isLast && allMembers.length < members.total) {
                  const next = await jiraAdmin.getGroupMembers(client, argv['group-name'], argv.startAt + allMembers.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allMembers = allMembers.concat(next.values);
                  members = next;
                }
              }
              const rows = allMembers.map((u) => [u.accountId || '', u.displayName || '', u.emailAddress || '', String(u.active)]);
              printTable(['ACCOUNT_ID', 'DISPLAY_NAME', 'EMAIL', 'ACTIVE'], rows);
            }
          )

          .command(
            'add-user <group-name>',
            'Add a user to a group',
            (y2) =>
              y2
                .positional('group-name', { type: 'string' })
                .option('account-id', { type: 'string', description: 'User account ID (required)' })
                .demandOption(['account-id']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.addUserToGroup(client, argv['group-name'], argv.accountId);
              console.log(`User ${argv.accountId} added to group ${argv['group-name']}`);
            }
          )

          .command(
            'remove-user <group-name>',
            'Remove a user from a group',
            (y2) =>
              y2
                .positional('group-name', { type: 'string' })
                .option('account-id', { type: 'string', description: 'User account ID (required)' })
                .demandOption(['account-id']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.removeUserFromGroup(client, argv['group-name'], argv.accountId);
              console.log(`User ${argv.accountId} removed from group ${argv['group-name']}`);
            }
          )

          .command(
            'search',
            'Search groups',
            (y2) =>
              y2
                .option('query', { type: 'string', description: 'Search query', default: '' })
                .option('max-results', { type: 'number', default: 50 }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const found = await jiraAdmin.findGroups(client, argv.query, argv.maxResults);
              const rows = (found.groups || []).map((g) => [g.groupId || '', g.name || '']);
              printTable(['GROUP_ID', 'NAME'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // dashboard (alias: dash)
    // =========================================================================
    .command(
      ['dashboard', 'dash'],
      'Manage dashboards',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List dashboards',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Search dashboards by name', default: '' })
                .option('max-results', { type: 'number', default: 50 })
                .option('start-at', { type: 'number', default: 0 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);

              if (argv.name) {
                let result = await jiraSchemes.searchDashboards(client, argv.name, argv.startAt, argv.maxResults);
                if (argv.all) {
                  while (!result.isLast && (result.values || []).length < result.total) {
                    const next = await jiraSchemes.searchDashboards(client, argv.name, argv.startAt + result.values.length, argv.maxResults);
                    if (!next.values || next.values.length === 0) break;
                    result.values = result.values.concat(next.values);
                    result.isLast = next.isLast;
                  }
                }
                if (isJSONOutput(argv)) { outputJSON(result); return; }

                const rows = (result.values || []).map((d) => [d.id || '', d.name || '', d.owner?.displayName || '']);
                printTable(['ID', 'NAME', 'OWNER'], rows);
                printPaginationHint((result.values || []).length, result.total || 0);
                return;
              }

              let result = await jiraSchemes.getDashboards(client, argv.startAt, argv.maxResults);
              if (argv.all) {
                while ((result.dashboards || []).length < result.total) {
                  const next = await jiraSchemes.getDashboards(client, argv.startAt + result.dashboards.length, argv.maxResults);
                  if (!next.dashboards || next.dashboards.length === 0) break;
                  result.dashboards = result.dashboards.concat(next.dashboards);
                }
              }
              if (isJSONOutput(argv)) { outputJSON(result); return; }

              const rows = (result.dashboards || []).map((d) => [d.id || '', d.name || '', d.owner?.displayName || '']);
              printTable(['ID', 'NAME', 'OWNER'], rows);
              printPaginationHint((result.dashboards || []).length, result.total || 0);
            }
          )

          .command(
            'get <id>',
            'Get a dashboard by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const d = await jiraSchemes.getDashboard(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(d); return; }

              console.log(`ID          ${d.id}`);
              console.log(`Name        ${d.name}`);
              console.log(`Description ${d.description || ''}`);
              if (d.owner) console.log(`Owner       ${d.owner.displayName}`);
            }
          )

          .command(
            'create',
            'Create a dashboard',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Dashboard name (required)' })
                .option('description', { type: 'string', description: 'Dashboard description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Dashboard> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const d = await jiraSchemes.createDashboard(client, body);
              console.log(`Dashboard created: ${d.name} (ID: ${d.id})`);
            }
          )

          .command(
            'update <id>',
            'Update a dashboard',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Dashboard name' })
                .option('description', { type: 'string', description: 'Dashboard description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Dashboard> = {};
              if (argv.name) body.name = argv.name;
              if (argv.description) body.description = argv.description;
              const d = await jiraSchemes.updateDashboard(client, argv.id, body);
              console.log(`Dashboard updated: ${d.name} (ID: ${d.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a dashboard',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteDashboard(client, argv.id);
              console.log(`Dashboard ${argv.id} deleted`);
            }
          )

          .command(
            'copy <id>',
            'Copy a dashboard',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Name for the copy (required)' })
                .option('description', { type: 'string', description: 'Description for the copy', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Dashboard> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const d = await jiraSchemes.copyDashboard(client, argv.id, body);
              console.log(`Dashboard copied: ${d.name} (ID: ${d.id})`);
            }
          )

          // gadget subgroup
          .command(
            'gadget',
            'Manage dashboard gadgets',
            (y2) =>
              y2
                .demandCommand(1)

                .command(
                  ['list <dashboard-id>', 'ls <dashboard-id>'],
                  'List gadgets on a dashboard',
                  (y3) => y3.positional('dashboard-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const gadgets = await jiraSchemes.getDashboardGadgets(client, argv['dashboard-id']);
                    if (isJSONOutput(argv)) { outputJSON(gadgets); return; }

                    const rows = (gadgets.gadgets || []).map((g) => [
                      String(g.id || ''),
                      g.title || '',
                      g.moduleKey || '',
                      g.uri || '',
                    ]);
                    printTable(['ID', 'TITLE', 'MODULE KEY', 'URI'], rows);
                  }
                )

                .command(
                  'add <dashboard-id>',
                  'Add a gadget to a dashboard',
                  (y3) =>
                    y3
                      .positional('dashboard-id', { type: 'string' })
                      .option('module-key', { type: 'string', description: 'Gadget module key', default: '' })
                      .option('uri', { type: 'string', description: 'Gadget URI', default: '' })
                      .option('title', { type: 'string', description: 'Gadget title', default: '' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const body: Partial<DashboardGadget> = {};
                    if (argv.moduleKey) body.moduleKey = argv.moduleKey;
                    if (argv.uri) body.uri = argv.uri;
                    if (argv.title) body.title = argv.title;
                    const gadget = await jiraSchemes.addDashboardGadget(client, argv['dashboard-id'], body);
                    console.log(`Gadget added: ${gadget.title} (ID: ${gadget.id})`);
                  }
                )

                .command(
                  'update <dashboard-id> <gadget-id>',
                  'Update a gadget on a dashboard',
                  (y3) =>
                    y3
                      .positional('dashboard-id', { type: 'string' })
                      .positional('gadget-id', { type: 'string' })
                      .option('title', { type: 'string', description: 'Gadget title', default: '' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    const body: Partial<DashboardGadget> = {};
                    if (argv.title) body.title = argv.title;
                    await jiraSchemes.updateDashboardGadget(client, argv['dashboard-id'], argv['gadget-id'], body);
                    console.log(`Gadget ${argv['gadget-id']} updated on dashboard ${argv['dashboard-id']}`);
                  }
                )

                .command(
                  'remove <dashboard-id> <gadget-id>',
                  'Remove a gadget from a dashboard',
                  (y3) =>
                    y3
                      .positional('dashboard-id', { type: 'string' })
                      .positional('gadget-id', { type: 'string' }),
                  async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
                    const client = getJiraClient(argv);
                    await jiraSchemes.removeDashboardGadget(client, argv['dashboard-id'], argv['gadget-id']);
                    console.log(`Gadget ${argv['gadget-id']} removed from dashboard ${argv['dashboard-id']}`);
                  }
                ),
            () => {}
          ),
      () => {}
    )

    // =========================================================================
    // role (admin)
    // =========================================================================
    .command(
      'role',
      'Manage project roles',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all project roles',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const roles = await jiraSchemes.getAllRoles(client);
              const rows = (roles || []).map((r) => [String(r.id || ''), r.name || '', r.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <role-id>',
            'Get a project role',
            (y2) => y2.positional('role-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const role = await jiraSchemes.getRole(client, argv['role-id']);
              outputJSON(role);
            }
          )

          .command(
            'create',
            'Create a project role',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Role name (required)' })
                .option('description', { type: 'string', description: 'Role description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<ProjectRole> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const role = await jiraSchemes.createRole(client, body);
              outputJSON(role);
            }
          )

          .command(
            'delete <role-id>',
            'Delete a project role',
            (y2) => y2.positional('role-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteRole(client, argv['role-id']);
              console.log(`Role ${argv['role-id']} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // issuelink (alias: il)
    // =========================================================================
    .command(
      ['issuelink', 'il'],
      'Manage issue links',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'create',
            'Create an issue link',
            (y2) =>
              y2
                .option('inward-issue', { type: 'string', description: 'Inward issue key (required)' })
                .option('outward-issue', { type: 'string', description: 'Outward issue key (required)' })
                .option('type', { type: 'string', description: 'Link type name (required)' })
                .demandOption(['inward-issue', 'outward-issue', 'type']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.createIssueLink(client, {
                type: { name: argv.type },
                inwardIssue: { key: argv['inward-issue'] },
                outwardIssue: { key: argv['outward-issue'] },
              });
              console.log('Issue link created.');
            }
          )

          .command(
            'get <link-id>',
            'Get an issue link',
            (y2) => y2.positional('link-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const link = await jiraAdmin.getIssueLink(client, argv['link-id']);
              outputJSON(link);
            }
          )

          .command(
            'delete <link-id>',
            'Delete an issue link',
            (y2) => y2.positional('link-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteIssueLink(client, argv['link-id']);
              console.log(`Issue link ${argv['link-id']} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // issuelinktype (alias: ilt)
    // =========================================================================
    .command(
      ['issuelinktype', 'ilt'],
      'Manage issue link types',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all issue link types',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const types = await jiraAdmin.getIssueLinkTypes(client);
              const rows = (types || []).map((t) => [t.id || '', t.name || '', t.inward || '', t.outward || '']);
              printTable(['ID', 'NAME', 'INWARD', 'OUTWARD'], rows);
            }
          )

          .command(
            'get <id>',
            'Get an issue link type',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const lt = await jiraAdmin.getIssueLinkType(client, argv.id);
              outputJSON(lt);
            }
          )

          .command(
            'create',
            'Create an issue link type',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Link type name (required)' })
                .option('inward', { type: 'string', description: 'Inward description (required)' })
                .option('outward', { type: 'string', description: 'Outward description (required)' })
                .demandOption(['name', 'inward', 'outward']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const lt = await jiraAdmin.createIssueLinkType(client, { name: argv.name, inward: argv.inward, outward: argv.outward });
              outputJSON(lt);
            }
          )

          .command(
            'update <id>',
            'Update an issue link type',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Link type name' })
                .option('inward', { type: 'string', description: 'Inward description' })
                .option('outward', { type: 'string', description: 'Outward description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const lt: Partial<IssueLinkType> = {};
              if (argv.name !== undefined) lt.name = argv.name;
              if (argv.inward !== undefined) lt.inward = argv.inward;
              if (argv.outward !== undefined) lt.outward = argv.outward;
              const result = await jiraAdmin.updateIssueLinkType(client, argv.id, lt);
              outputJSON(result);
            }
          )

          .command(
            'delete <id>',
            'Delete an issue link type',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteIssueLinkType(client, argv.id);
              console.log(`Issue link type ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // screen
    // =========================================================================
    .command(
      'screen',
      'Manage screens',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List screens',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getScreens(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getScreens(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((s) => [String(s.id || ''), s.name || '', s.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'create',
            'Create a screen',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Screen name (required)' })
                .option('description', { type: 'string', description: 'Screen description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Screen> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const screen = await jiraSchemes.createScreen(client, body);
              outputJSON(screen);
            }
          )

          .command(
            'delete <screen-id>',
            'Delete a screen',
            (y2) => y2.positional('screen-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteScreen(client, argv['screen-id']);
              console.log(`Screen ${argv['screen-id']} deleted.`);
            }
          )

          .command(
            'tabs <screen-id>',
            'List tabs for a screen',
            (y2) => y2.positional('screen-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const tabs = await jiraSchemes.getScreenTabs(client, argv['screen-id']);
              const rows = (tabs || []).map((t) => [String(t.id || ''), t.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          )

          .command(
            'tab-fields <screen-id> <tab-id>',
            'List fields for a screen tab',
            (y2) =>
              y2
                .positional('screen-id', { type: 'number' })
                .positional('tab-id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const fields = await jiraSchemes.getScreenTabFields(client, argv['screen-id'], argv['tab-id']);
              const rows = (fields || []).map((f) => [f.id || '', f.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // workflow (alias: wf)
    // =========================================================================
    .command(
      ['workflow', 'wf'],
      'Manage workflows',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all workflows',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const workflows = await jiraSchemes.getWorkflows(client);
              const rows = (workflows || []).map((wf) => [wf.name || '', wf.description || '', String(wf.isDefault)]);
              printTable(['NAME', 'DESCRIPTION', 'DEFAULT'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // workflowscheme (alias: wfs)
    // =========================================================================
    .command(
      ['workflowscheme', 'wfs'],
      'Manage workflow schemes',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List workflow schemes',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getWorkflowSchemes(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getWorkflowSchemes(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((s) => [String(s.id || ''), s.name || '', s.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a workflow scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const scheme = await jiraSchemes.getWorkflowScheme(client, Number(argv.id));
              outputJSON(scheme);
            }
          )

          .command(
            'create',
            'Create a workflow scheme',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Scheme name (required)' })
                .option('description', { type: 'string', description: 'Scheme description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<WorkflowScheme> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const scheme = await jiraSchemes.createWorkflowScheme(client, body);
              outputJSON(scheme);
            }
          )

          .command(
            'update <id>',
            'Update a workflow scheme',
            (y2) =>
              y2
                .positional('id', { type: 'number' })
                .option('name', { type: 'string', description: 'Scheme name' })
                .option('description', { type: 'string', description: 'Scheme description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<WorkflowScheme> = {};
              if (argv.name !== undefined) body.name = argv.name;
              if (argv.description !== undefined) body.description = argv.description;
              const scheme = await jiraSchemes.updateWorkflowScheme(client, Number(argv.id), body);
              outputJSON(scheme);
            }
          )

          .command(
            'delete <id>',
            'Delete a workflow scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteWorkflowScheme(client, Number(argv.id));
              console.log(`Workflow scheme ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // permissionscheme (alias: ps)
    // =========================================================================
    .command(
      ['permissionscheme', 'ps'],
      'Manage permission schemes',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List permission schemes',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const schemes = await jiraSchemes.getPermissionSchemes(client);
              const rows = (schemes || []).map((s) => [String(s.id || ''), s.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a permission scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const scheme = await jiraSchemes.getPermissionScheme(client, Number(argv.id));
              outputJSON(scheme);
            }
          )

          .command(
            'create',
            'Create a permission scheme',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Scheme name (required)' })
                .option('description', { type: 'string', description: 'Scheme description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<PermissionScheme> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const scheme = await jiraSchemes.createPermissionScheme(client, body);
              outputJSON(scheme);
            }
          )

          .command(
            'delete <id>',
            'Delete a permission scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deletePermissionScheme(client, Number(argv.id));
              console.log(`Permission scheme ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // notificationscheme (alias: ns)
    // =========================================================================
    .command(
      ['notificationscheme', 'ns'],
      'Manage notification schemes',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List notification schemes',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getNotificationSchemes(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getNotificationSchemes(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((s) => [String(s.id || ''), s.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a notification scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const scheme = await jiraSchemes.getNotificationScheme(client, Number(argv.id));
              outputJSON(scheme);
            }
          )

          .command(
            'create',
            'Create a notification scheme',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Scheme name (required)' })
                .option('description', { type: 'string', description: 'Scheme description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<NotificationScheme> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const scheme = await jiraSchemes.createNotificationScheme(client, body);
              outputJSON(scheme);
            }
          )

          .command(
            'delete <id>',
            'Delete a notification scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteNotificationScheme(client, Number(argv.id));
              console.log(`Notification scheme ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // issuesecurityscheme (alias: iss)
    // =========================================================================
    .command(
      ['issuesecurityscheme', 'iss'],
      'Manage issue security schemes',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List issue security schemes',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const schemes = await jiraSchemes.getIssueSecuritySchemes(client);
              const rows = (schemes || []).map((s) => [String(s.id || ''), s.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          )

          .command(
            'get <id>',
            'Get an issue security scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const scheme = await jiraSchemes.getIssueSecurityScheme(client, Number(argv.id));
              outputJSON(scheme);
            }
          )

          .command(
            'create',
            'Create an issue security scheme',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Scheme name (required)' })
                .option('description', { type: 'string', description: 'Scheme description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<IssueSecurityScheme> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const scheme = await jiraSchemes.createIssueSecurityScheme(client, body);
              outputJSON(scheme);
            }
          )

          .command(
            'delete <id>',
            'Delete an issue security scheme',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteIssueSecurityScheme(client, Number(argv.id));
              console.log(`Issue security scheme ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // fieldconfig (alias: fc)
    // =========================================================================
    .command(
      ['fieldconfig', 'fc'],
      'Manage field configurations',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List field configurations',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getFieldConfigurations(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getFieldConfigurations(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((fc) => [String(fc.id || ''), fc.name || '', String(fc.isDefault)]);
              printTable(['ID', 'NAME', 'DEFAULT'], rows);
            }
          )

          .command(
            'create',
            'Create a field configuration',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Configuration name (required)' })
                .option('description', { type: 'string', description: 'Configuration description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<FieldConfiguration> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const fc = await jiraSchemes.createFieldConfiguration(client, body);
              outputJSON(fc);
            }
          )

          .command(
            'delete <id>',
            'Delete a field configuration',
            (y2) => y2.positional('id', { type: 'number' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteFieldConfiguration(client, Number(argv.id));
              console.log(`Field configuration ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // issuetypescheme (alias: its)
    // =========================================================================
    .command(
      ['issuetypescheme', 'its'],
      'Manage issue type schemes',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List issue type schemes',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getIssueTypeSchemes(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getIssueTypeSchemes(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((s) => [String(s.id || ''), s.name || '']);
              printTable(['ID', 'NAME'], rows);
            }
          )

          .command(
            'create',
            'Create an issue type scheme',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Scheme name (required)' })
                .option('description', { type: 'string', description: 'Scheme description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<IssueTypeScheme> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const scheme = await jiraSchemes.createIssueTypeScheme(client, body);
              outputJSON(scheme);
            }
          )

          .command(
            'delete <id>',
            'Delete an issue type scheme',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteIssueTypeScheme(client, argv.id);
              console.log(`Issue type scheme ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // serverinfo (alias: si)
    // =========================================================================
    .command(
      ['serverinfo', 'si'],
      'Show Jira server information',
      () => {},
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        const info = await jiraAdmin.getServerInfo(client);
        outputJSON(info);
      }
    )

    // =========================================================================
    // webhook (alias: wh)
    // =========================================================================
    .command(
      ['webhook', 'wh'],
      'Manage webhooks',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List webhooks',
            (y2) =>
              y2
                .option('start-at', { type: 'number', default: 0 })
                .option('max-results', { type: 'number', default: 50 })
                .option('all', { type: 'boolean', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              let result = await jiraSchemes.getWebhooks(client, argv.startAt, argv.maxResults);
              let allValues = result.values || [];
              if (argv.all) {
                while (!result.isLast && allValues.length < result.total) {
                  const next = await jiraSchemes.getWebhooks(client, argv.startAt + allValues.length, argv.maxResults);
                  if (!next.values || next.values.length === 0) break;
                  allValues = allValues.concat(next.values);
                  result = next;
                }
              }
              const rows = allValues.map((wh) => [
                String(wh.id || ''),
                wh.jqlFilter || '',
                (wh.events || []).join(', '),
              ]);
              printTable(['ID', 'JQL', 'EVENTS'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // attachment (alias: att)
    // =========================================================================
    .command(
      ['attachment', 'att'],
      'Manage attachments',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <id>',
            'Get an attachment',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const att = await jiraAdmin.getAttachment(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(att); return; }

              const rows = [[att.id || '', att.filename || '', String(att.size || ''), att.mimeType || '', att.created || '']];
              printTable(['ID', 'FILENAME', 'SIZE', 'MIME TYPE', 'CREATED'], rows);
            }
          )

          .command(
            'delete <id>',
            'Delete an attachment',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteAttachment(client, argv.id);
              console.log(`Attachment ${argv.id} deleted.`);
            }
          )

          .command(
            'meta',
            'Show attachment settings',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const meta = await jiraAdmin.getAttachmentMeta(client);
              outputJSON(meta);
            }
          ),
      () => {}
    )

    // =========================================================================
    // audit
    // =========================================================================
    .command(
      'audit',
      'Show audit records',
      (y) =>
        y
          .option('start-at', { type: 'number', default: 0 })
          .option('max-results', { type: 'number', default: 50 })
          .option('all', { type: 'boolean', default: false }),
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        let records = await jiraAdmin.getAuditRecords(client, argv.startAt, argv.maxResults);
        let allRecords = records.records || [];
        if (argv.all) {
          while (allRecords.length < (records.total ?? 0)) {
            const next = await jiraAdmin.getAuditRecords(client, argv.startAt + allRecords.length, argv.maxResults);
            if (!next.records || next.records.length === 0) break;
            allRecords = allRecords.concat(next.records);
          }
        }
        const rows = allRecords.map((r) => [String(r.id || ''), r.summary || '', r.created || '', r.category || '']);
        printTable(['ID', 'SUMMARY', 'CREATED', 'CATEGORY'], rows);
      }
    )

    // =========================================================================
    // banner
    // =========================================================================
    .command(
      'banner',
      'Manage announcement banner',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get',
            'Get the announcement banner',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const banner = await jiraAdmin.getAnnouncementBanner(client);
              outputJSON(banner);
            }
          )

          .command(
            'set',
            'Set the announcement banner',
            (y2) =>
              y2
                .option('message', { type: 'string', description: 'Banner message' })
                .option('enabled', { type: 'boolean', description: 'Enable the banner' })
                .option('dismissible', { type: 'boolean', description: 'Allow dismissing the banner' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const banner: Partial<AnnouncementBanner> = {};
              if (argv.message !== undefined) banner.message = argv.message;
              if (argv.enabled !== undefined) banner.isEnabled = argv.enabled;
              if (argv.dismissible !== undefined) banner.isDismissible = argv.dismissible;
              await jiraAdmin.setAnnouncementBanner(client, banner);
              console.log('Announcement banner updated.');
            }
          ),
      () => {}
    )

    // =========================================================================
    // configuration (alias: config)
    // =========================================================================
    .command(
      ['configuration', 'config'],
      'Show Jira configuration',
      () => {},
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        const config = await jiraAdmin.getConfiguration(client);
        outputJSON(config);
      }
    )

    // =========================================================================
    // permission (alias: perm)
    // =========================================================================
    .command(
      ['permission', 'perm'],
      'Manage permissions',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'mine',
            'Show my permissions',
            (y2) =>
              y2
                .option('project', { type: 'string', description: 'Project key (uses profile default if not set)' })
                .option('issue', { type: 'string', description: 'Issue key', default: '' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = defaultProject(argv);
              const perms = await jiraAdmin.getMyPermissions(client, project, argv.issue);
              const rows = Object.values(perms || {}).map((p: UserPermission) => [p.key || '', p.name || '', String(p.havePermission)]);
              printTable(['KEY', 'NAME', 'HAVE_PERMISSION'], rows);
            }
          )

          .command(
            'all',
            'List all permissions',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const perms = await jiraAdmin.getAllPermissions(client);
              const rows = Object.values(perms || {}).map((p: UserPermission) => [p.key || '', p.name || '', String(p.havePermission)]);
              printTable(['KEY', 'NAME', 'HAVE_PERMISSION'], rows);
            }
          ),
      () => {}
    )

    // =========================================================================
    // task
    // =========================================================================
    .command(
      'task',
      'Manage async tasks',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <task-id>',
            'Get a task',
            (y2) => y2.positional('task-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const task = await jiraAdmin.getTask(client, argv['task-id']);
              outputJSON(task);
            }
          )

          .command(
            'cancel <task-id>',
            'Cancel a task',
            (y2) => y2.positional('task-id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.cancelTask(client, argv['task-id']);
              console.log(`Task ${argv['task-id']} cancelled.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // projectcategory (alias: pc)
    // =========================================================================
    .command(
      ['projectcategory', 'pc'],
      'Manage project categories',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List project categories',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const cats = await jiraProjects.getProjectCategories(client);
              const rows = (cats || []).map((c) => [c.id || '', c.name || '', c.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a project category',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const cat = await jiraProjects.getProjectCategory(client, argv.id);
              outputJSON(cat);
            }
          )

          .command(
            'create',
            'Create a project category',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Category name (required)' })
                .option('description', { type: 'string', description: 'Category description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const cat: Partial<ProjectCategory> = { name: argv.name };
              if (argv.description) cat.description = argv.description;
              const result = await jiraProjects.createProjectCategory(client, cat);
              outputJSON(result);
            }
          )

          .command(
            'update <id>',
            'Update a project category',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Category name' })
                .option('description', { type: 'string', description: 'Category description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const cat: Partial<ProjectCategory> = {};
              if (argv.name !== undefined) cat.name = argv.name;
              if (argv.description !== undefined) cat.description = argv.description;
              const result = await jiraProjects.updateProjectCategory(client, argv.id, cat);
              outputJSON(result);
            }
          )

          .command(
            'delete <id>',
            'Delete a project category',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraProjects.deleteProjectCategory(client, argv.id);
              console.log(`Project category ${argv.id} deleted.`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // component (alias: comp) — from jira_resources.go
    // =========================================================================
    .command(
      ['component', 'comp'],
      'Manage project components',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <id>',
            'Get a component by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const comp = await jiraAdmin.getComponent(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(comp); return; }

              console.log(`ID            ${comp.id}`);
              console.log(`Name          ${comp.name}`);
              console.log(`Description   ${comp.description || ''}`);
              console.log(`Project       ${comp.project || ''}`);
              console.log(`Assignee Type ${comp.assigneeType || ''}`);
              if (comp.lead) console.log(`Lead          ${comp.lead.displayName}`);
            }
          )

          .command(
            'create',
            'Create a component',
            (y2) =>
              y2
                .option('project', { type: 'string', description: 'Project key (uses profile default if not set)' })
                .option('name', { type: 'string', description: 'Component name (required)' })
                .option('description', { type: 'string', description: 'Component description', default: '' })
                .option('lead', { type: 'string', description: 'Lead account ID', default: '' })
                .option('assignee-type', { type: 'string', description: 'Assignee type', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const project = defaultProject(argv);
              const body: Partial<ProjectComponent> & { leadAccountId?: string } = { project, name: argv.name };
              if (argv.description) body.description = argv.description;
              if (argv.lead) body.leadAccountId = argv.lead;
              if (argv.assigneeType) body.assigneeType = argv.assigneeType;
              const comp = await jiraAdmin.createComponent(client, body);
              console.log(`Component created: ${comp.name} (ID: ${comp.id})`);
            }
          )

          .command(
            'update <id>',
            'Update a component',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Component name' })
                .option('description', { type: 'string', description: 'Component description' })
                .option('lead', { type: 'string', description: 'Lead account ID' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<ProjectComponent> & { leadAccountId?: string } = {};
              if (argv.name) body.name = argv.name;
              if (argv.description) body.description = argv.description;
              if (argv.lead) body.leadAccountId = argv.lead;
              const comp = await jiraAdmin.updateComponent(client, argv.id, body);
              console.log(`Component updated: ${comp.name} (ID: ${comp.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a component',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteComponent(client, argv.id);
              console.log(`Component ${argv.id} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // version (alias: ver) — from jira_resources.go
    // =========================================================================
    .command(
      ['version', 'ver'],
      'Manage project versions',
      (y) =>
        y
          .demandCommand(1)

          .command(
            'get <id>',
            'Get a version by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const ver = await jiraAdmin.getVersion(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(ver); return; }

              console.log(`ID           ${ver.id}`);
              console.log(`Name         ${ver.name}`);
              console.log(`Description  ${ver.description || ''}`);
              console.log(`Released     ${ver.released}`);
              console.log(`Archived     ${ver.archived}`);
              console.log(`Start Date   ${ver.startDate || ''}`);
              console.log(`Release Date ${ver.releaseDate || ''}`);
            }
          )

          .command(
            'create',
            'Create a version',
            (y2) =>
              y2
                .option('project-id', { type: 'number', description: 'Project ID (required)' })
                .option('name', { type: 'string', description: 'Version name (required)' })
                .option('description', { type: 'string', description: 'Version description', default: '' })
                .option('start-date', { type: 'string', description: 'Start date (YYYY-MM-DD)', default: '' })
                .option('release-date', { type: 'string', description: 'Release date (YYYY-MM-DD)', default: '' })
                .option('released', { type: 'boolean', description: 'Whether the version is released', default: false })
                .option('archived', { type: 'boolean', description: 'Whether the version is archived', default: false })
                .demandOption(['project-id', 'name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const ver: Partial<Version> = { projectId: argv.projectId, name: argv.name };
              if (argv.description) ver.description = argv.description;
              if (argv.startDate) ver.startDate = argv.startDate;
              if (argv.releaseDate) ver.releaseDate = argv.releaseDate;
              if (argv.released !== undefined) ver.released = argv.released;
              if (argv.archived !== undefined) ver.archived = argv.archived;
              const result = await jiraAdmin.createVersion(client, ver);
              console.log(`Version created: ${result.name} (ID: ${result.id})`);
            }
          )

          .command(
            'update <id>',
            'Update a version',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Version name' })
                .option('description', { type: 'string', description: 'Version description' })
                .option('start-date', { type: 'string', description: 'Start date (YYYY-MM-DD)' })
                .option('release-date', { type: 'string', description: 'Release date (YYYY-MM-DD)' })
                .option('released', { type: 'boolean', description: 'Whether the version is released' })
                .option('archived', { type: 'boolean', description: 'Whether the version is archived' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const ver = await jiraAdmin.getVersion(client, argv.id);
              if (argv.name) ver.name = argv.name;
              if (argv.description) ver.description = argv.description;
              if (argv.startDate) ver.startDate = argv.startDate;
              if (argv.releaseDate) ver.releaseDate = argv.releaseDate;
              if (argv.released !== undefined) ver.released = argv.released;
              if (argv.archived !== undefined) ver.archived = argv.archived;
              const result = await jiraAdmin.updateVersion(client, argv.id, ver);
              console.log(`Version updated: ${result.name} (ID: ${result.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a version',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteVersion(client, argv.id);
              console.log(`Version ${argv.id} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // field — from jira_resources.go
    // =========================================================================
    .command(
      'field',
      'Manage fields',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all fields',
            (y2) =>
              y2
                .option('custom', { type: 'boolean', description: 'Show only custom fields', default: false }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const fields = await jiraSchemes.getFields(client);
              if (isJSONOutput(argv)) {
                if (argv.custom) {
                  outputJSON((fields || []).filter((f) => f.custom));
                } else {
                  outputJSON(fields);
                }
                return;
              }

              const rows = (fields || [])
                .filter((f) => !argv.custom || f.custom)
                .map((f) => [f.id || '', f.name || '', f.schema?.type || '', String(f.custom)]);
              printTable(['ID', 'NAME', 'TYPE', 'CUSTOM'], rows);
            }
          )

          .command(
            'create',
            'Create a custom field',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Field name (required)' })
                .option('type', { type: 'string', description: 'Field type (required)' })
                .option('description', { type: 'string', description: 'Field description', default: '' })
                .option('search-key', { type: 'string', description: 'Searcher key', default: '' })
                .demandOption(['name', 'type']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Field> & { searcherKey?: string; type?: string; description?: string } = { name: argv.name, type: argv.type };
              if (argv.description) body.description = argv.description;
              if (argv.searchKey) body.searcherKey = argv.searchKey;
              const field = await jiraSchemes.createCustomField(client, body);
              console.log(`Field created: ${field.name} (ID: ${field.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a custom field',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.deleteCustomField(client, argv.id);
              console.log(`Field ${argv.id} deleted`);
            }
          )

          .command(
            'trash <id>',
            'Move a custom field to trash',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.trashCustomField(client, argv.id);
              console.log(`Field ${argv.id} moved to trash`);
            }
          )

          .command(
            'restore <id>',
            'Restore a custom field from trash',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraSchemes.restoreCustomField(client, argv.id);
              console.log(`Field ${argv.id} restored`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // label — from jira_resources.go
    // =========================================================================
    .command(
      'label',
      'List labels',
      (y) =>
        y
          .option('start-at', { type: 'number', default: 0 })
          .option('max-results', { type: 'number', default: 50 })
          .option('all', { type: 'boolean', default: false }),
      async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
        const client = getJiraClient(argv);
        let result = await jiraAdmin.getLabels(client, argv.startAt, argv.maxResults);
        let allLabels = result.values || [];
        if (argv.all) {
          while (!result.isLast && allLabels.length < result.total) {
            const next = await jiraAdmin.getLabels(client, argv.startAt + allLabels.length, argv.maxResults);
            if (!next.values || next.values.length === 0) break;
            allLabels = allLabels.concat(next.values);
            result = next;
          }
        }
        if (isJSONOutput(argv)) { outputJSON(allLabels); return; }

        const rows = allLabels.map((label: string) => [label]);
        printTable(['LABEL'], rows);
      }
    )

    // =========================================================================
    // issuetype (alias: it) — from jira_resources.go
    // =========================================================================
    .command(
      ['issuetype', 'it'],
      'Manage issue types',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all issue types',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const types = await jiraAdmin.getAllIssueTypes(client);
              const rows = (types || []).map((t) => [t.id || '', t.name || '', String(t.subtask), t.description || '']);
              printTable(['ID', 'NAME', 'SUBTASK', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <id>',
            'Get an issue type by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const it = await jiraAdmin.getIssueType(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(it); return; }

              console.log(`ID          ${it.id}`);
              console.log(`Name        ${it.name}`);
              console.log(`Subtask     ${it.subtask}`);
              console.log(`Description ${it.description || ''}`);
            }
          )

          .command(
            'create',
            'Create an issue type',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Issue type name (required)' })
                .option('description', { type: 'string', description: 'Issue type description', default: '' })
                .option('type', { type: 'string', description: 'Issue type: standard or subtask', default: 'standard' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<IssueType> & { type?: string } = { name: argv.name, type: argv.type };
              if (argv.description) body.description = argv.description;
              const it = await jiraAdmin.createIssueType(client, body);
              console.log(`Issue type created: ${it.name} (ID: ${it.id})`);
            }
          )

          .command(
            'update <id>',
            'Update an issue type',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Issue type name' })
                .option('description', { type: 'string', description: 'Issue type description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<IssueType> = {};
              if (argv.name) body.name = argv.name;
              if (argv.description) body.description = argv.description;
              const it = await jiraAdmin.updateIssueType(client, argv.id, body);
              console.log(`Issue type updated: ${it.name} (ID: ${it.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete an issue type',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteIssueType(client, argv.id);
              console.log(`Issue type ${argv.id} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // priority (alias: pri) — from jira_resources.go
    // =========================================================================
    .command(
      ['priority', 'pri'],
      'Manage priorities',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all priorities',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const priorities = await jiraAdmin.getAllPriorities(client);
              const rows = (priorities || []).map((p) => [p.id || '', p.name || '', p.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a priority by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const p = await jiraAdmin.getPriority(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(p); return; }

              console.log(`ID           ${p.id}`);
              console.log(`Name         ${p.name}`);
              console.log(`Description  ${p.description || ''}`);
              console.log(`Status Color ${p.statusColor || ''}`);
            }
          )

          .command(
            'create',
            'Create a priority',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Priority name (required)' })
                .option('description', { type: 'string', description: 'Priority description', default: '' })
                .option('status-color', { type: 'string', description: 'Status color hex', default: '#ffffff' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Priority> = { name: argv.name, statusColor: argv.statusColor };
              if (argv.description) body.description = argv.description;
              const p = await jiraAdmin.createPriority(client, body);
              console.log(`Priority created: ${p.name} (ID: ${p.id})`);
            }
          )

          .command(
            'update <id>',
            'Update a priority',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Priority name' })
                .option('description', { type: 'string', description: 'Priority description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Priority> = {};
              if (argv.name) body.name = argv.name;
              if (argv.description) body.description = argv.description;
              const p = await jiraAdmin.updatePriority(client, argv.id, body);
              console.log(`Priority updated: ${p.name} (ID: ${p.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a priority',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deletePriority(client, argv.id);
              console.log(`Priority ${argv.id} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // resolution (alias: res) — from jira_resources.go
    // =========================================================================
    .command(
      ['resolution', 'res'],
      'Manage resolutions',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all resolutions',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const resolutions = await jiraAdmin.getAllResolutions(client);
              const rows = (resolutions || []).map((r) => [r.id || '', r.name || '', r.description || '']);
              printTable(['ID', 'NAME', 'DESCRIPTION'], rows);
            }
          )

          .command(
            'get <id>',
            'Get a resolution by ID',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const r = await jiraAdmin.getResolution(client, argv.id);
              if (isJSONOutput(argv)) { outputJSON(r); return; }

              console.log(`ID          ${r.id}`);
              console.log(`Name        ${r.name}`);
              console.log(`Description ${r.description || ''}`);
            }
          )

          .command(
            'create',
            'Create a resolution',
            (y2) =>
              y2
                .option('name', { type: 'string', description: 'Resolution name (required)' })
                .option('description', { type: 'string', description: 'Resolution description', default: '' })
                .demandOption(['name']),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Resolution> = { name: argv.name };
              if (argv.description) body.description = argv.description;
              const r = await jiraAdmin.createResolution(client, body);
              console.log(`Resolution created: ${r.name} (ID: ${r.id})`);
            }
          )

          .command(
            'update <id>',
            'Update a resolution',
            (y2) =>
              y2
                .positional('id', { type: 'string' })
                .option('name', { type: 'string', description: 'Resolution name' })
                .option('description', { type: 'string', description: 'Resolution description' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const body: Partial<Resolution> = {};
              if (argv.name) body.name = argv.name;
              if (argv.description) body.description = argv.description;
              const r = await jiraAdmin.updateResolution(client, argv.id, body);
              console.log(`Resolution updated: ${r.name} (ID: ${r.id})`);
            }
          )

          .command(
            'delete <id>',
            'Delete a resolution',
            (y2) => y2.positional('id', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              await jiraAdmin.deleteResolution(client, argv.id);
              console.log(`Resolution ${argv.id} deleted`);
            }
          ),
      () => {}
    )

    // =========================================================================
    // status (alias: st) — from jira_resources.go
    // =========================================================================
    .command(
      ['status', 'st'],
      'Manage statuses',
      (y) =>
        y
          .demandCommand(1)

          .command(
            ['list', 'ls'],
            'List all statuses',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const statuses = await jiraAdmin.getAllStatuses(client);
              const rows = (statuses || []).map((s) => [s.id || '', s.name || '', s.statusCategory?.name || '']);
              printTable(['ID', 'NAME', 'CATEGORY'], rows);
            }
          )

          .command(
            'get <id-or-name>',
            'Get a status by ID or name',
            (y2) => y2.positional('id-or-name', { type: 'string' }),
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const s = await jiraAdmin.getStatus(client, argv['id-or-name']);
              if (isJSONOutput(argv)) { outputJSON(s); return; }

              console.log(`ID          ${s.id}`);
              console.log(`Name        ${s.name}`);
              console.log(`Description ${s.description || ''}`);
              console.log(`Category    ${s.statusCategory?.name || ''}`);
            }
          )

          .command(
            'categories',
            'List status categories',
            () => {},
            async (_a) => {
              const argv = _a as Partial<JiraArgv> as JiraArgv;
              const client = getJiraClient(argv);
              const categories = await jiraAdmin.getStatusCategories(client);
              const rows = (categories || []).map((c) => [
                String(c.id || ''),
                c.key || '',
                c.name || '',
                c.colorName || '',
              ]);
              printTable(['ID', 'KEY', 'NAME', 'COLOR'], rows);
            }
          ),
      () => {}
    );
}

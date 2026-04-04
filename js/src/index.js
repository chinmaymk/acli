#!/usr/bin/env node
import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';
import { readFileSync } from 'fs';
import { createRequire } from 'module';
import { fileURLToPath } from 'url';
import path from 'path';

import { registerJiraCommands } from './cmd/jira.js';
import { registerConfluenceCommands } from './cmd/confluence.js';
import { registerBitbucketCommands } from './cmd/bitbucket.js';
import { registerConfigCommands } from './cmd/config.js';
import { registerCommandsCommand, buildCommandTreeFromList } from './cmd/commands.js';
import { outputJSON } from './cmd/helpers.js';

// Read version info injected at build time via package.json or environment
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

let version = 'dev';
let commit = 'none';
let buildDate = 'unknown';

try {
  const require = createRequire(import.meta.url);
  const pkg = require(path.join(__dirname, '..', 'package.json'));
  if (pkg.version) version = pkg.version;
} catch {
  // ignore
}

// Allow build tooling to inject values via environment variables
if (process.env.ACLI_VERSION) version = process.env.ACLI_VERSION;
if (process.env.ACLI_COMMIT) commit = process.env.ACLI_COMMIT;
if (process.env.ACLI_DATE) buildDate = process.env.ACLI_DATE;

// Handle stdin args mode: if argv[2] === '-', read stdin as JSON array of args
async function resolveArgs() {
  const rawArgs = hideBin(process.argv);
  if (rawArgs.length === 1 && rawArgs[0] === '-') {
    return new Promise((resolve, reject) => {
      let data = '';
      process.stdin.setEncoding('utf8');
      process.stdin.on('data', (chunk) => { data += chunk; });
      process.stdin.on('end', () => {
        const trimmed = data.trim();
        if (!trimmed) {
          reject(new Error('stdin argument mode expects a JSON array, got empty input'));
          return;
        }
        let args;
        try {
          args = JSON.parse(trimmed);
        } catch (err) {
          reject(new Error(`invalid stdin argument JSON array: ${err.message}`));
          return;
        }
        resolve(args);
      });
      process.stdin.on('error', reject);
    });
  }
  return rawArgs;
}

const args = await resolveArgs().catch((err) => {
  console.error(`acli: ${err.message}`);
  process.exit(1);
});

// Build the command tree for the `commands` subcommand
function getCommandTree() {
  return buildCommandTreeFromList([
    { command: 'jira', describe: 'Manage Jira issues, projects, boards, and sprints', parent: 'acli', aliases: ['j'] },
    { command: 'jira issue', describe: 'Manage Jira issues', parent: 'acli', aliases: ['i'] },
    { command: 'jira issue list', describe: 'List issues', parent: 'acli jira' },
    { command: 'jira issue get <issue-key>', describe: 'Get issue details', parent: 'acli jira' },
    { command: 'jira issue create', describe: 'Create a new issue', parent: 'acli jira' },
    { command: 'jira issue update <issue-key>', describe: 'Update an existing issue', parent: 'acli jira' },
    { command: 'jira issue comment <issue-key>', describe: 'Add a comment to an issue', parent: 'acli jira' },
    { command: 'jira issue transition <issue-key> <transition>', describe: 'Transition an issue to a new status', parent: 'acli jira' },
    { command: 'jira issue assign <issue-key> [assignee]', describe: 'Assign an issue to a user', parent: 'acli jira' },
    { command: 'jira issue delete <issue-key>', describe: 'Delete an issue', parent: 'acli jira' },
    { command: 'jira project', describe: 'Manage Jira projects', parent: 'acli jira', aliases: ['p'] },
    { command: 'jira project list', describe: 'List projects', parent: 'acli jira' },
    { command: 'jira board', describe: 'Manage Jira boards', parent: 'acli jira', aliases: ['b'] },
    { command: 'jira board list', describe: 'List boards', parent: 'acli jira' },
    { command: 'jira board issues <board-id>', describe: 'List issues on a board', parent: 'acli jira' },
    { command: 'jira sprint', describe: 'Manage Jira sprints', parent: 'acli jira', aliases: ['s'] },
    { command: 'jira sprint list <board-id>', describe: 'List sprints for a board', parent: 'acli jira' },
    { command: 'jira sprint issues <sprint-id>', describe: 'List issues in a sprint', parent: 'acli jira' },
    { command: 'confluence', describe: 'Manage Confluence spaces and pages', parent: 'acli', aliases: ['c', 'conf'] },
    { command: 'confluence space', describe: 'Manage Confluence spaces', parent: 'acli confluence', aliases: ['s'] },
    { command: 'confluence space list', describe: 'List spaces', parent: 'acli confluence' },
    { command: 'confluence page', describe: 'Manage Confluence pages', parent: 'acli confluence', aliases: ['p'] },
    { command: 'confluence page list', describe: 'List pages in a space', parent: 'acli confluence' },
    { command: 'confluence page get <page-id>', describe: 'Get page details', parent: 'acli confluence' },
    { command: 'bitbucket', describe: 'Manage Bitbucket repositories, pull requests, and pipelines', parent: 'acli', aliases: ['bb'] },
    { command: 'bitbucket repo', describe: 'Manage Bitbucket repositories', parent: 'acli bitbucket', aliases: ['r'] },
    { command: 'bitbucket repo list [workspace]', describe: 'List repositories', parent: 'acli bitbucket' },
    { command: 'bitbucket repo get [workspace] <repo>', describe: 'Get repository details', parent: 'acli bitbucket' },
    { command: 'bitbucket pr', describe: 'Manage Bitbucket pull requests', parent: 'acli bitbucket', aliases: ['p'] },
    { command: 'bitbucket pr list [workspace] <repo>', describe: 'List pull requests', parent: 'acli bitbucket' },
    { command: 'bitbucket pr get [workspace] <repo> <id>', describe: 'Get pull request details', parent: 'acli bitbucket' },
    { command: 'bitbucket pr create [workspace] <repo>', describe: 'Create a pull request', parent: 'acli bitbucket' },
    { command: 'bitbucket pr approve [workspace] <repo> <id>', describe: 'Approve a pull request', parent: 'acli bitbucket' },
    { command: 'bitbucket pr merge [workspace] <repo> <id>', describe: 'Merge a pull request', parent: 'acli bitbucket' },
    { command: 'bitbucket pr decline [workspace] <repo> <id>', describe: 'Decline a pull request', parent: 'acli bitbucket' },
    { command: 'bitbucket pipeline', describe: 'Manage Bitbucket pipelines', parent: 'acli bitbucket', aliases: ['pipe'] },
    { command: 'bitbucket pipeline list [workspace] <repo>', describe: 'List pipelines', parent: 'acli bitbucket' },
    { command: 'bitbucket pipeline get [workspace] <repo> <id>', describe: 'Get pipeline details', parent: 'acli bitbucket' },
    { command: 'config', describe: 'Manage acli configuration profiles', parent: 'acli' },
    { command: 'config setup <profile-name>', describe: 'Create or update a configuration profile', parent: 'acli config' },
    { command: 'config list', describe: 'List all configured profiles', parent: 'acli config', aliases: ['ls'] },
    { command: 'config show [profile-name]', describe: 'Show details for a profile', parent: 'acli config' },
    { command: 'config delete <profile-name>', describe: 'Delete a configuration profile', parent: 'acli config', aliases: ['rm'] },
    { command: 'config set-default <profile-name>', describe: 'Set the default profile', parent: 'acli config' },
    { command: 'config set-defaults [profile-name]', describe: 'Set default project and workspace for a profile', parent: 'acli config' },
    { command: 'version', describe: 'Print version information', parent: 'acli' },
    { command: 'commands', describe: 'List all commands in machine-readable JSON (for agents and scripts)', parent: 'acli' },
  ]);
}

const cli = yargs(args)
  .scriptName('acli')
  .usage('acli - Atlassian CLI\n\nUsage: $0 <command> [options]')
  .option('profile', {
    alias: 'p',
    type: 'string',
    describe: 'configuration profile to use (defaults to the default profile)',
  })
  .option('output', {
    alias: 'o',
    type: 'string',
    default: 'text',
    describe: 'output format: text or json (json is recommended for programmatic/agent use)',
  })
  .command(
    ['jira', 'j'],
    'Manage Jira issues, projects, boards, and sprints',
    (yargs) => registerJiraCommands(yargs).demandCommand(1, 'Specify a jira subcommand').help(),
  )
  .command(
    ['confluence', 'c', 'conf'],
    'Manage Confluence spaces and pages',
    (yargs) => registerConfluenceCommands(yargs).demandCommand(1, 'Specify a confluence subcommand').help(),
  )
  .command(
    ['bitbucket', 'bb'],
    'Manage Bitbucket repositories, pull requests, and pipelines',
    (yargs) => registerBitbucketCommands(yargs).demandCommand(1, 'Specify a bitbucket subcommand').help(),
  )
  .command(
    ['config', 'cfg'],
    'Manage acli configuration profiles',
    (yargs) => registerConfigCommands(yargs).demandCommand(1, 'Specify a config subcommand').help(),
  )
  .command(
    'version',
    'Print version information',
    () => {},
    () => {
      console.log(`acli ${version} (commit: ${commit}, built: ${buildDate})`);
    },
  );

registerCommandsCommand(cli, getCommandTree);

cli
  .demandCommand(1, 'Specify a command. Use --help to list available commands.')
  .strict()
  .help()
  .alias('help', 'h')
  .parse();

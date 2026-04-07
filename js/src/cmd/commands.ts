import { outputJSON } from './helpers.js';
import type { Argv } from 'yargs';

export interface FlagInfo {
  name: string;
  shorthand?: string;
  type: string;
  default?: string;
  required?: boolean;
  usage: string;
}

export interface CommandInfo {
  name: string;
  fullCommand: string;
  description: string;
  aliases: string[];
  usage: string;
  args?: string;
  flags?: FlagInfo[];
  subcommands?: CommandInfo[];
}

/**
 * Registers the `commands` subcommand onto a yargs instance.
 * The command outputs the full command tree as JSON.
 */
export function registerCommandsCommand(yargsInstance: Argv, getTree: () => CommandInfo): Argv {
  return yargsInstance.command(
    'commands',
    'List all commands in machine-readable JSON (for agents and scripts)',
    () => {},
    () => {
      const tree = getTree ? getTree() : { name: 'acli', fullCommand: 'acli', description: '', aliases: [], usage: 'acli', subcommands: [] };
      outputJSON(tree);
    },
  );
}

interface CommandListEntry {
  command: string;
  describe: string;
  parent?: string;
  aliases?: string[];
  flags?: FlagInfo[];
}

/**
 * Builds a CommandInfo tree from a flat list of command descriptors.
 * This is the primary way to generate the command tree in yargs (since yargs
 * does not expose the same introspection as cobra).
 */
export function buildCommandTreeFromList(commandList: CommandListEntry[]): CommandInfo {
  const root: CommandInfo = {
    name: 'acli',
    fullCommand: 'acli',
    description: 'Atlassian CLI - manage Jira, Confluence, and Bitbucket from the terminal',
    aliases: [],
    usage: 'acli',
    flags: [
      { name: 'profile', shorthand: 'p', type: 'string', default: '', usage: 'configuration profile to use (defaults to the default profile)' },
      { name: 'output', shorthand: 'o', type: 'string', default: 'text', usage: 'output format: text or json (json is recommended for programmatic/agent use)' },
    ],
    subcommands: [],
  };

  // Build a map by full command path
  const byPath: Record<string, CommandInfo> = { acli: root };

  for (const entry of commandList) {
    const parent = entry.parent || 'acli';
    const parts = entry.command.split(' ');
    const name = parts[0];
    const args = parts.slice(1).join(' ');
    const fullCommand = `${parent} ${name}`;

    const info: CommandInfo = {
      name,
      fullCommand,
      description: entry.describe || '',
      aliases: entry.aliases || [],
      usage: args ? `${fullCommand} ${args}` : fullCommand,
      args: args || undefined,
      flags: entry.flags || [],
      subcommands: [],
    };

    byPath[fullCommand] = info;

    const parentNode = byPath[parent];
    if (parentNode) {
      parentNode.subcommands!.push(info);
    }
  }

  return root;
}

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
 * Extracts the argument portion from a yargs usage string.
 * e.g. "get <issue-key>" -> "<issue-key>"
 */
function extractArgsFromUse(use: string): string {
  const idx = use.indexOf(' ');
  if (idx === -1) return '';
  return use.slice(idx + 1);
}

// Shape of yargs internal options object (partial)
interface YargsOptionsInternal {
  describe?: Record<string, string | undefined>;
  string?: Record<string, unknown>;
  default?: Record<string, unknown>;
  alias?: Record<string, string | string[]>;
  demandedOptions?: Record<string, unknown>;
}

// Minimal shape of yargs instance internals we access
interface YargsInternals {
  usage?: string;
  parsed?: unknown;
  _yargsOptionsSpecified?: Record<string, unknown>;
  _getOptionsAccess?: () => YargsOptionsInternal;
  options?: YargsOptionsInternal;
  getContext?: () => unknown;
}

/**
 * Builds a command tree from a yargs instance.
 * Walks the internal command registry to gather CommandInfo for each command.
 */
function buildCommandTree(yargsInstance: YargsInternals, parentPath?: string): CommandInfo {
  const usage = yargsInstance.usage || 'acli';
  const name = parentPath ? parentPath.split(' ').pop()! : 'acli';
  const fullCommand = parentPath || 'acli';

  const info: CommandInfo = {
    name,
    fullCommand,
    description: '',
    aliases: [],
    usage: fullCommand,
    args: extractArgsFromUse(fullCommand) || undefined,
    flags: [],
    subcommands: [],
  };

  // Access yargs internals to read registered options
  try {
    yargsInstance.getContext?.();

    // Read options from yargs internal state
    const yargsOptions: YargsOptionsInternal = yargsInstance._getOptionsAccess
      ? yargsInstance._getOptionsAccess()
      : (yargsInstance.options || {});

    const describe = yargsOptions.describe || {};
    const defaults = yargsOptions.default || {};
    const alias = yargsOptions.alias || {};
    const demandedOptions = yargsOptions.demandedOptions || {};

    for (const [flagName, flagDesc] of Object.entries(describe)) {
      if (flagName === '_' || flagName === '$0') continue;

      // Find shorthand from alias map
      let shorthand = '';
      for (const [aliasKey, aliases] of Object.entries(alias)) {
        if (Array.isArray(aliases) && aliases.includes(flagName)) {
          if (aliasKey.length === 1) shorthand = aliasKey;
        } else if (aliasKey === flagName) {
          const found = (Array.isArray(aliases) ? aliases : [aliases]).find((a) => a.length === 1);
          if (found) shorthand = found;
        }
      }

      info.flags!.push({
        name: flagName,
        shorthand: shorthand || undefined,
        type: 'string',
        default: defaults[flagName] != null ? String(defaults[flagName]) : undefined,
        required: !!demandedOptions[flagName],
        usage: typeof flagDesc === 'string' ? flagDesc : '',
      });
    }
  } catch {
    // If yargs internals are not accessible, skip flag collection
  }

  // Suppress unused variable warning — usage is read for potential future use
  void usage;

  return info;
}

/**
 * Collects command info from a registered yargs command module.
 */
function commandInfoFromModule(
  command: string,
  describe: string,
  parentPath: string,
  _builder?: (yargs: Argv) => Argv,
): CommandInfo {
  const parts = command.split(' ');
  const name = parts[0];
  const args = parts.slice(1).join(' ');
  const fullCommand = parentPath ? `${parentPath} ${name}` : name;

  const info: CommandInfo = {
    name,
    fullCommand,
    description: describe || '',
    aliases: [],
    usage: args ? `${fullCommand} ${args}` : fullCommand,
    args: args || undefined,
    flags: [],
    subcommands: [],
  };

  return info;
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
      const tree = getTree ? getTree() : { name: 'acli', subcommands: [] };
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

// Suppress unused function warnings for helpers only used internally
void buildCommandTree;
void commandInfoFromModule;

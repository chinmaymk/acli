import { outputJSON } from './helpers.js';

/**
 * @typedef {{ name: string, fullCommand: string, description: string, aliases?: string[], usage: string, args?: string, flags?: FlagInfo[], subcommands?: CommandInfo[] }} CommandInfo
 * @typedef {{ name: string, shorthand?: string, type: string, default?: string, required?: boolean, usage: string }} FlagInfo
 */

/**
 * Extracts the argument portion from a yargs usage string.
 * e.g. "get <issue-key>" -> "<issue-key>"
 * @param {string} use
 * @returns {string}
 */
function extractArgsFromUse(use) {
  const idx = use.indexOf(' ');
  if (idx === -1) return '';
  return use.slice(idx + 1);
}

/**
 * Builds a command tree from a yargs instance.
 * Walks the internal command registry to gather CommandInfo for each command.
 * @param {object} yargsInstance - a yargs instance
 * @param {string} [parentPath] - full command path prefix
 * @returns {CommandInfo}
 */
function buildCommandTree(yargsInstance, parentPath) {
  const usage = yargsInstance.usage || 'acli';
  const name = parentPath ? parentPath.split(' ').pop() : 'acli';
  const fullCommand = parentPath || 'acli';

  /** @type {CommandInfo} */
  const info = {
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
    const context = yargsInstance.getContext ? yargsInstance.getContext() : {};
    const parsed = yargsInstance.parsed || {};
    const options = yargsInstance._yargsOptionsSpecified || {};

    // Read options from yargs internal state
    const yargsOptions = yargsInstance._getOptionsAccess
      ? yargsInstance._getOptionsAccess()
      : (yargsInstance.options || {});

    const describe = yargsOptions.describe || {};
    const type = yargsOptions.string || {};
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

      info.flags.push({
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

  return info;
}

/**
 * Collects command info from a registered yargs command module.
 * @param {string} command - command use string (e.g. "jira [subcommand]")
 * @param {string} describe - short description
 * @param {string} parentPath - parent command path
 * @param {Function|undefined} builder - builder function
 * @returns {CommandInfo}
 */
function commandInfoFromModule(command, describe, parentPath, builder) {
  const parts = command.split(' ');
  const name = parts[0];
  const args = parts.slice(1).join(' ');
  const fullCommand = parentPath ? `${parentPath} ${name}` : name;

  /** @type {CommandInfo} */
  const info = {
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
 * @param {import('yargs').Argv} yargsInstance
 * @param {() => CommandInfo} getTree - callback that returns the full command tree
 * @returns {import('yargs').Argv}
 */
export function registerCommandsCommand(yargsInstance, getTree) {
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

/**
 * Builds a CommandInfo tree from a flat list of command descriptors.
 * This is the primary way to generate the command tree in yargs (since yargs
 * does not expose the same introspection as cobra).
 *
 * @param {Array<{ command: string, describe: string, parent?: string, aliases?: string[] }>} commandList
 * @returns {CommandInfo}
 */
export function buildCommandTreeFromList(commandList) {
  /** @type {CommandInfo} */
  const root = {
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
  const byPath = { acli: root };

  for (const entry of commandList) {
    const parent = entry.parent || 'acli';
    const parts = entry.command.split(' ');
    const name = parts[0];
    const args = parts.slice(1).join(' ');
    const fullCommand = `${parent} ${name}`;

    /** @type {CommandInfo} */
    const info = {
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
      parentNode.subcommands.push(info);
    }
  }

  return root;
}

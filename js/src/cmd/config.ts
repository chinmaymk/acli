import readline from 'readline';
import type { Argv, ArgumentsCamelCase } from 'yargs';
import * as config from '../internal/config/config.js';

type ConfigArgv = ArgumentsCamelCase<{ 'profile-name'?: string }>;

/**
 * Prompt for input with an optional default value display.
 * If current is set, shows "[current]" and returns current on empty input.
 * If placeholder is set (and no current), shows "(placeholder)" and returns placeholder on empty input.
 */
async function promptWithDefault(
  rl: readline.Interface,
  label: string,
  current: string,
  placeholder: string,
): Promise<string> {
  let prompt: string;
  if (current !== '') {
    prompt = `  ${label} [${current}]: `;
  } else if (placeholder !== '') {
    prompt = `  ${label} (${placeholder}): `;
  } else {
    prompt = `  ${label}: `;
  }

  return new Promise((resolve) => {
    rl.question(prompt, (input) => {
      input = input.trim();
      if (input === '') {
        if (current !== '') {
          resolve(current);
        } else {
          resolve(placeholder);
        }
      } else {
        resolve(input);
      }
    });
  });
}

/**
 * Mask an API token showing first 4 and last 4 characters.
 */
function maskToken(token: string): string {
  if (!token) return '';
  if (token.length <= 8) return '****';
  return token.slice(0, 4) + '****' + token.slice(-4);
}

/**
 * Register all config subcommands onto the given yargs instance.
 */
export function registerConfigCommands(yargs: Argv): Argv {
  return yargs
    .command(
      'setup <profile-name>',
      'Create or update a configuration profile',
      (yargs: Argv) => {
        return yargs.positional('profile-name', {
          type: 'string',
          describe: 'Name of the profile to create or update',
        });
      },
      async (argv: ConfigArgv) => {
        const profileName = argv['profile-name'] as string;

        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        const existing: Partial<config.Profile> = (cfg.profiles && cfg.profiles[profileName]) || {};
        if (Object.keys(existing).length > 0) {
          console.log(`Updating profile "${profileName}" (press Enter to keep current value)\n`);
        } else {
          console.log(`Creating profile "${profileName}"\n`);
        }

        const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

        let atlassianURL: string, email: string, apiToken: string;
        try {
          atlassianURL = await promptWithDefault(rl, 'Atlassian URL', existing.atlassian_url || '', 'https://your-instance.atlassian.net');
          email = await promptWithDefault(rl, 'Email', existing.email || '', '');
          const maskedExisting = maskToken(existing.api_token || '');
          apiToken = await promptWithDefault(rl, 'API Token', maskedExisting, '');
          // If user just pressed enter and there was an existing token, keep it
          if (apiToken === maskedExisting && existing.api_token) {
            apiToken = existing.api_token;
          }
        } finally {
          rl.close();
        }

        const profile: config.Profile = {
          name: profileName,
          atlassian_url: atlassianURL!.replace(/\/+$/, ''),
          email: email!,
          api_token: apiToken!,
          defaults: existing.defaults || { project: '', workspace: '', bb_project: '' },
        };

        if (!cfg.profiles) {
          cfg.profiles = {};
        }

        const isFirst = Object.keys(cfg.profiles).length === 0;
        cfg.profiles[profileName] = profile;
        if (isFirst || !cfg.default_profile) {
          cfg.default_profile = profileName;
        }

        try {
          config.save(cfg);
        } catch (err) {
          throw new Error(`saving config: ${(err as Error).message}`);
        }

        console.log(`\nProfile "${profileName}" saved to ~/.config/acli/config.json`);
        if (isFirst || cfg.default_profile === profileName) {
          console.log(`Profile "${profileName}" is the default profile`);
        }
      },
    )
    .command(
      ['list', 'ls'],
      'List all configured profiles',
      () => {},
      () => {
        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        if (!cfg.profiles || Object.keys(cfg.profiles).length === 0) {
          console.log("No profiles configured. Run 'acli config setup' to create one.");
          return;
        }

        // Build rows for tabular output
        const rows: string[][] = [['PROFILE', 'DEFAULT', 'URL', 'EMAIL']];
        for (const [name, p] of Object.entries(cfg.profiles)) {
          const def = name === cfg.default_profile ? '*' : '';
          rows.push([name, def, p.atlassian_url || '', p.email || '']);
        }

        // Compute column widths
        const colWidths = rows[0].map((_, ci) =>
          Math.max(...rows.map((r) => (r[ci] || '').length)),
        );

        for (const row of rows) {
          const line = row
            .map((cell, ci) => (cell || '').padEnd(colWidths[ci]))
            .join('  ');
          console.log(line);
        }
      },
    )
    .command(
      'show [profile-name]',
      'Show details for a profile (tokens are masked)',
      (yargs: Argv) => {
        return yargs.positional('profile-name', {
          type: 'string',
          describe: 'Name of the profile to show (defaults to current default profile)',
        });
      },
      (argv: ConfigArgv) => {
        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        let p: config.Profile;
        try {
          p = config.getProfile(cfg, argv['profile-name'] || '');
        } catch (err) {
          throw err;
        }

        console.log(`Profile: ${p.name}`);
        console.log(`  Atlassian URL:    ${p.atlassian_url || ''}`);
        console.log(`  Email:            ${p.email || ''}`);
        console.log(`  API Token:        ${maskToken(p.api_token || '')}`);

        const d = p.defaults || ({} as Partial<config.Defaults>);
        if (d.project || d.workspace || d.bb_project) {
          console.log('  Defaults:');
          if (d.project) {
            console.log(`    Project:        ${d.project}`);
          }
          if (d.workspace) {
            console.log(`    Workspace:      ${d.workspace}`);
          }
          if (d.bb_project) {
            console.log(`    BB Project:     ${d.bb_project}`);
          }
        }
      },
    )
    .command(
      ['delete <profile-name>', 'rm <profile-name>'],
      'Delete a configuration profile',
      (yargs: Argv) => {
        return yargs.positional('profile-name', {
          type: 'string',
          describe: 'Name of the profile to delete',
        });
      },
      (argv: ConfigArgv) => {
        const profileName = argv['profile-name'] as string;

        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        if (!cfg.profiles || !cfg.profiles[profileName]) {
          throw new Error(`profile "${profileName}" not found`);
        }

        delete cfg.profiles[profileName];

        if (cfg.default_profile === profileName) {
          cfg.default_profile = '';
          // If one profile remains, make it the new default
          const remaining = Object.keys(cfg.profiles);
          if (remaining.length > 0) {
            cfg.default_profile = remaining[0];
          }
        }

        try {
          config.save(cfg);
        } catch (err) {
          throw new Error(`saving config: ${(err as Error).message}`);
        }

        console.log(`Profile "${profileName}" deleted`);
      },
    )
    .command(
      'set-default <profile-name>',
      'Set the default profile',
      (yargs: Argv) => {
        return yargs.positional('profile-name', {
          type: 'string',
          describe: 'Name of the profile to set as default',
        });
      },
      (argv: ConfigArgv) => {
        const profileName = argv['profile-name'] as string;

        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        if (!cfg.profiles || !cfg.profiles[profileName]) {
          throw new Error(`profile "${profileName}" not found`);
        }

        cfg.default_profile = profileName;

        try {
          config.save(cfg);
        } catch (err) {
          throw new Error(`saving config: ${(err as Error).message}`);
        }

        console.log(`Default profile set to "${profileName}"`);
      },
    )
    .command(
      'set-defaults [profile-name]',
      'Set default project and workspace for a profile',
      (yargs: Argv) => {
        return yargs.positional('profile-name', {
          type: 'string',
          describe: 'Name of the profile to set defaults for (uses current default if omitted)',
        });
      },
      async (argv: ConfigArgv) => {
        let cfg: config.Config;
        try {
          cfg = config.load();
        } catch (err) {
          throw new Error(`loading config: ${(err as Error).message}`);
        }

        let profile: config.Profile;
        try {
          profile = config.getProfile(cfg, argv['profile-name'] || '');
        } catch (err) {
          throw err;
        }

        const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

        console.log(`Setting defaults for profile "${profile.name}"\n`);

        const d = profile.defaults || ({} as Partial<config.Defaults>);
        let project: string, workspace: string, bbProject: string;
        try {
          project = await promptWithDefault(rl, 'Default Jira project key', d.project || '', '');
          workspace = await promptWithDefault(rl, 'Default Bitbucket workspace', d.workspace || '', '');
          bbProject = await promptWithDefault(rl, 'Default Bitbucket project key', d.bb_project || '', '');
        } finally {
          rl.close();
        }

        profile.defaults = {
          project: project!,
          workspace: workspace!,
          bb_project: bbProject!,
        };

        cfg.profiles[profile.name] = profile;

        try {
          config.save(cfg);
        } catch (err) {
          throw new Error(`saving config: ${(err as Error).message}`);
        }

        console.log(`\nDefaults saved for profile "${profile.name}"`);
        if (project!) {
          console.log(`  Default project:   ${project}`);
        }
        if (workspace!) {
          console.log(`  Default workspace: ${workspace}`);
        }
        if (bbProject!) {
          console.log(`  Default BB project: ${bbProject}`);
        }
      },
    );
}

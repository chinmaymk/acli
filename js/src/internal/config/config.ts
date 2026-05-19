import fs from 'fs';
import path from 'path';
import os from 'os';

export interface Defaults {
  project: string;
  workspace: string;
  bb_project: string;
}

export interface Profile {
  name: string;
  atlassian_url: string;
  email: string;
  api_token: string;
  defaults: Defaults;
}

export interface Config {
  default_profile: string;
  profiles: Record<string, Profile>;
}

export function configDir(): string {
  return path.join(os.homedir(), '.config', 'acli');
}

export function load(): Config {
  const dir = configDir();
  const filePath = path.join(dir, 'config.json');

  let data: string;
  try {
    data = fs.readFileSync(filePath, 'utf8');
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
      return { default_profile: '', profiles: {} };
    }
    throw new Error(`reading config: ${(err as Error).message}`);
  }

  let cfg: Config;
  try {
    cfg = JSON.parse(data) as Config;
  } catch (err) {
    throw new Error(`parsing config: ${(err as Error).message}`);
  }
  applyEnvOverrides(cfg);
  return cfg;
}

// applyEnvOverrides overlays ACLI_<PROFILE>_* environment variables onto each
// profile. Variable names use uppercased profile names, e.g. for profile "jira":
//
//   ACLI_JIRA_API_TOKEN
//   ACLI_JIRA_EMAIL
//   ACLI_JIRA_ATLASSIAN_URL
//
// Env vars take precedence over values in config.json.
function applyEnvOverrides(cfg: Config): void {
  for (const name of Object.keys(cfg.profiles)) {
    const prefix = `ACLI_${name.toUpperCase()}_`;
    const p = cfg.profiles[name];
    const token = process.env[`${prefix}API_TOKEN`];
    const email = process.env[`${prefix}EMAIL`];
    const url = process.env[`${prefix}ATLASSIAN_URL`];
    if (token) p.api_token = token;
    if (email) p.email = email;
    if (url) p.atlassian_url = url;
  }
}

export function save(config: Config): void {
  const dir = configDir();
  fs.mkdirSync(dir, { recursive: true, mode: 0o700 });

  const filePath = path.join(dir, 'config.json');
  const data = JSON.stringify(config, null, 2);
  fs.writeFileSync(filePath, data, { mode: 0o600 });
}

export function getProfile(config: Config, name?: string): Profile {
  if (name) {
    const p = config.profiles[name];
    if (!p) {
      throw new Error(`profile "${name}" not found`);
    }
    return p;
  }

  if (config.default_profile) {
    const p = config.profiles[config.default_profile];
    if (p) {
      return p;
    }
  }

  const keys = Object.keys(config.profiles);
  if (keys.length === 1) {
    return config.profiles[keys[0]];
  }
  if (keys.length === 0) {
    throw new Error("no profiles configured, run 'acli config setup' to create one");
  }
  throw new Error("multiple profiles configured, specify one with -p or set a default with 'acli config set-default <name>'");
}

import fs from 'fs';
import path from 'path';
import os from 'os';

/**
 * @typedef {{ project?: string, workspace?: string, bb_project?: string }} Defaults
 * @typedef {{ name: string, atlassian_url: string, email: string, api_token: string, defaults?: Defaults }} Profile
 * @typedef {{ default_profile?: string, profiles: Record<string, Profile> }} Config
 */

export function configDir() {
  return path.join(os.homedir(), '.config', 'acli');
}

export function load() {
  const dir = configDir();
  const filePath = path.join(dir, 'config.json');

  let data;
  try {
    data = fs.readFileSync(filePath, 'utf8');
  } catch (err) {
    if (err.code === 'ENOENT') {
      return { profiles: {} };
    }
    throw new Error(`reading config: ${err.message}`);
  }

  let cfg;
  try {
    cfg = JSON.parse(data);
  } catch (err) {
    throw new Error(`parsing config: ${err.message}`);
  }
  return cfg;
}

export function save(config) {
  const dir = configDir();
  fs.mkdirSync(dir, { recursive: true, mode: 0o700 });

  const filePath = path.join(dir, 'config.json');
  const data = JSON.stringify(config, null, 2);
  fs.writeFileSync(filePath, data, { mode: 0o600 });
}

export function getProfile(config, name) {
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

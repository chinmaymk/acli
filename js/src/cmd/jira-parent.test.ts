import assert from 'node:assert/strict';
import { execFile } from 'node:child_process';
import { mkdir, mkdtemp, rm, writeFile } from 'node:fs/promises';
import { createServer } from 'node:http';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { test } from 'node:test';
import { fileURLToPath } from 'node:url';
import { promisify } from 'node:util';

test('create a parent, attach multiple children, and move an existing child', async (t) => {
  const issues = new Map<string, { key: string; fields: { parent?: { key: string }; summary: string } }>();
  const server = createServer(async (req, res) => {
    const url = new URL(req.url!, 'http://localhost');
    res.setHeader('Content-Type', 'application/json');
    if (req.method === 'GET' && url.pathname === '/rest/api/3/search/jql') {
      const children = [...issues.values()].filter(issue =>
        issue.fields.parent && url.searchParams.get('jql') === `parent = ${issue.fields.parent.key}`);
      res.end(JSON.stringify({ issues: children, total: children.length }));
      return;
    }

    let body = '';
    for await (const chunk of req) body += chunk;
    if (req.method === 'POST' && url.pathname === '/rest/api/3/issue') {
      const key = `TEST-${issues.size + 10}`;
      issues.set(key, { key, fields: JSON.parse(body).fields });
      res.writeHead(201).end(JSON.stringify({ key }));
    } else if (req.method === 'PUT' && url.pathname.startsWith('/rest/api/3/issue/')) {
      const issue = issues.get(url.pathname.slice('/rest/api/3/issue/'.length));
      if (!issue) {
        res.writeHead(404).end();
        return;
      }
      Object.assign(issue.fields, JSON.parse(body).fields);
      res.writeHead(204).end();
    } else {
      res.writeHead(404).end();
    }
  });
  await new Promise<void>(resolve => server.listen(0, '127.0.0.1', resolve));
  t.after(() => server.close());
  const address = server.address();
  assert.ok(address && typeof address === 'object');

  const home = await mkdtemp(join(tmpdir(), 'acli-parent-'));
  t.after(() => rm(home, { recursive: true, force: true }));
  const configDir = join(home, '.config', 'acli');
  await mkdir(configDir, { recursive: true });
  await writeFile(join(configDir, 'config.json'), JSON.stringify({
    profiles: { test: { atlassian_url: `http://127.0.0.1:${address.port}`, api_token: 'test' } },
  }));

  // Set ACLI_TEST_BIN to run the same workflow against the compiled Go CLI.
  const binary = process.env.ACLI_TEST_BIN;
  async function run(...args: string[]) {
    const cliArgs = binary ? args : [fileURLToPath(new URL('../index.js', import.meta.url)), ...args];
    const { stdout } = await promisify(execFile)(binary || process.execPath, cliArgs, {
      env: { ...process.env, HOME: home, USERPROFILE: home },
    });
    return JSON.parse(stdout);
  }
  function create(summary: string, ...args: string[]) {
    return run('jira', 'issue', 'create', '--project', 'TEST', '--summary', summary, '-o', 'json', ...args);
  }

  await assert.rejects(create('Invalid child', '--parent', ' '), /--parent requires an issue key/);
  assert.equal(issues.size, 0);
  const parent = await create('Release notes');
  assert.ok(issues.has(parent.key));
  assert.equal(issues.get(parent.key)?.fields.parent, undefined);
  const first = await create('Draft notes', '--type', 'Sub-task', '--parent', parent.key);
  const second = await create('Review notes', '--type', 'Sub-task', '--parent', parent.key);
  for (const child of [first, second]) {
    assert.deepEqual(issues.get(child.key)?.fields.parent, { key: parent.key });
  }

  const search = () => run('jira', 'search', '--jql', `parent = ${parent.key}`, '-o', 'json');
  const childKeys = async () => (await search()).issues.map((issue: { key: string }) => issue.key);
  assert.deepEqual(await childKeys(), [first.key, second.key]);

  const newParent = await create('Next release');
  await run('jira', 'issue', 'edit', first.key, '--parent', newParent.key, '-o', 'json');
  assert.deepEqual(issues.get(first.key)?.fields.parent, { key: newParent.key });
  assert.deepEqual(await childKeys(), [second.key]);
});

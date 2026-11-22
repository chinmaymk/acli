import assert from 'node:assert/strict';
import { test } from 'node:test';
import { searchJQL } from '../internal/jira/projects.js';

const client = { baseURL: 'https://jira.invalid', email: '', apiToken: 'test' };

test('search follows opaque cursors across empty pages and retains query options', async (t) => {
  const tokens: (string | null)[] = [];
  const pages = [
    { issues: [{ key: 'TEST-1' }], nextPageToken: 'a+/=', isLast: false },
    { issues: [], nextPageToken: 'b' },
    { issues: [{ key: 'TEST-2' }], isLast: true },
  ];
  t.mock.method(globalThis, 'fetch', async (url: string) => {
    const { pathname, searchParams } = new URL(url);
    assert.equal(pathname, '/rest/api/3/search/jql');
    assert.equal(searchParams.get('jql'), 'project = TEST');
    assert.equal(searchParams.get('fields'), 'summary,status');
    assert.equal(searchParams.get('maxResults'), '1');
    assert.equal(searchParams.has('startAt'), false);
    tokens.push(searchParams.get('nextPageToken'));
    return Response.json(pages.shift());
  });

  const result = await searchJQL(client, 'project = TEST', 0, 1, ['summary', 'status'], undefined, true);
  assert.deepEqual(tokens, [null, 'a+/=', 'b']);
  assert.deepEqual(result.issues.map(i => i.key), ['TEST-1', 'TEST-2']);
  assert.equal(result.isLast, true);
  assert.equal(result.nextPageToken, undefined);
});

test('single-page search exposes the cursor and --all rejects a non-advancing cursor', async (t) => {
  const mock = t.mock.method(globalThis, 'fetch', async () =>
    Response.json({ issues: [], nextPageToken: 'same' }));
  const page = await searchJQL(client, 'project = TEST');
  assert.equal(mock.mock.callCount(), 1);
  assert.equal(page.nextPageToken, 'same');

  await assert.rejects(searchJQL(client, 'project = TEST', 0, 1, [], undefined, true), /repeated nextPageToken/);
  assert.equal(mock.mock.callCount(), 3);
});

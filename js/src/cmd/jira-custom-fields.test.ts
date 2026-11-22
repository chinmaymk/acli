import assert from 'node:assert/strict';
import { test } from 'node:test';
import { getIssue } from '../internal/jira/issues.js';
import { searchJQL } from '../internal/jira/projects.js';

test('issue get and search retain Jira custom fields in JSON', async (t) => {
  const fields = { summary: 'Test', customfield_1: null, customfield_2: [{ value: 'A' }] };
  const issue = { key: 'TEST-1', fields };
  t.mock.method(globalThis, 'fetch', async (url: string) =>
    Response.json(url.includes('/search/jql') ? { issues: [issue] } : issue));
  const client = { baseURL: 'https://jira.invalid', email: '', apiToken: 'test' };
  const found = await getIssue(client, issue.key);
  const search = await searchJQL(client, 'project = TEST', 0, 1, Object.keys(fields));

  for (const result of [found, search.issues[0]]) {
    assert.equal(result.fields.customfield_1, null);
    assert.deepEqual(JSON.parse(JSON.stringify(result)).fields, fields);
  }
});

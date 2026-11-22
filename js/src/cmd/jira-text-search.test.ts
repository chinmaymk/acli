import assert from 'node:assert/strict';
import { test } from 'node:test';
import { textSearchJQL, issuePreviews } from './jira-text-search.js';

test('combine literal words and filters, or resolve an issue key', () => {
  assert.equal(textSearchJQL('release notes', { project: 'DEV', mine: true, updatedSince: '30d' }),
    String.raw`text ~ "\"release\"" AND text ~ "\"notes\"" AND project = "DEV" AND assignee = currentUser() AND updated >= "-30d"`);
  assert.equal(textSearchJQL('dev-12', {}), 'key = "DEV-12"');
  assert.equal(textSearchJQL('', { assignee: 'a"b' }), String.raw`assignee = "a\"b"`);
  assert.throws(() => textSearchJQL('', {}));
  assert.throws(() => textSearchJQL('...', {}));
  assert.throws(() => textSearchJQL('', { mine: true, assignee: 'a' }));
  assert.throws(() => textSearchJQL('test', { updatedSince: 'yesterday' }));
});

test('preview title matches first and include text around a description match', () => {
  const description = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Review release notes' }] }] };
  const issues = [
    { key: 'DEV-1', fields: { summary: 'Prepare documentation', description } },
    { key: 'DEV-2', fields: { summary: 'Release notes' } },
  ];
  const previews = issuePreviews(issues, 'release notes', 'https://example.test');
  assert.deepEqual(previews.map(item => item.key), ['DEV-2', 'DEV-1']);
  assert.equal(previews[0].excerpt, undefined);
  assert.equal(previews[1].excerpt, 'Review release notes');
  assert.equal(previews[1].url, 'https://example.test/browse/DEV-1');
});

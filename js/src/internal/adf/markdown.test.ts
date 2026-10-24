import { test } from 'node:test';
import assert from 'node:assert/strict';
import { markdownToADF } from './markdown.js';

test('table regex does not misparse a bullet list after a pipe-containing line', () => {
  const md = ['The list of options |is below|:', '- foo', '- bar'].join('\n');
  const doc = markdownToADF(md);
  const types = doc.content.map((b) => b.type);
  assert.ok(!types.includes('table'), `expected no table, got types: ${types.join(',')}`);
  assert.ok(types.includes('bulletList'), `expected bulletList, got types: ${types.join(',')}`);
});

test('valid table with separator row still parses as a table', () => {
  const md = ['| Col A | Col B |', '|-------|-------|', '| v1    | v2    |'].join('\n');
  const doc = markdownToADF(md);
  const types = doc.content.map((b) => b.type);
  assert.deepEqual(types, ['table']);
});

test('two-row pipe-containing block (no separator) parses as a table with header row', () => {
  const md = ['| Col A | Col B |', '| v1    | v2    |'].join('\n');
  const doc = markdownToADF(md);
  const table = doc.content[0];
  assert.equal(table.type, 'table');
  const rows = (table.content ?? []) as { content?: { type: string }[] }[];
  assert.equal(rows.length, 2);
  const headerCellType = rows[0].content?.[0]?.type;
  assert.equal(headerCellType, 'tableHeader');
});

test('headings h1 through h6', () => {
  for (let level = 1; level <= 6; level++) {
    const md = '#'.repeat(level) + ' heading ' + level;
    const doc = markdownToADF(md);
    const node = doc.content[0];
    assert.equal(node.type, 'heading');
    assert.equal((node.attrs as { level: number }).level, level);
  }
});

test('underscore italic does not mangle snake_case identifiers', () => {
  const doc = markdownToADF('Set MY_LOG_LEVEL=1 to disable');
  const para = doc.content[0];
  assert.equal(para.type, 'paragraph');
  const inline = (para.content ?? []) as { type: string; text?: string; marks?: { type: string }[] }[];
  const hasItalicOnIdentifier = inline.some((n) => n.text?.includes('LOG') && n.marks?.some((m) => m.type === 'em'));
  assert.equal(hasItalicOnIdentifier, false, `identifier was italicized: ${JSON.stringify(inline)}`);
});

test('GitHub-style admonition produces a panel node', () => {
  const md = '> [!warning] heads up';
  const doc = markdownToADF(md);
  const node = doc.content[0];
  assert.equal(node.type, 'panel');
  assert.equal((node.attrs as { panelType: string }).panelType, 'warning');
});

test('@mention(id, name) produces an ADF mention node', () => {
  const doc = markdownToADF('hi @mention(abc-123, Jane Doe)!');
  const para = doc.content[0];
  const inline = (para.content ?? []) as { type: string; attrs?: { id: string; text: string } }[];
  const mention = inline.find((n) => n.type === 'mention');
  assert.ok(mention, `expected mention node, got: ${JSON.stringify(inline)}`);
  assert.equal(mention?.attrs?.id, 'abc-123');
  assert.equal(mention?.attrs?.text, '@Jane Doe');
});

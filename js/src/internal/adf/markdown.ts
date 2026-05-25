import type { ADFDocument, ADFNode, ADFMark } from '../types.js';

const PANEL_TYPES = new Set(['info', 'warning', 'note', 'success', 'error']);

export function markdownToADF(md: string): ADFDocument {
  const lines = md.split('\n');
  const content: ADFNode[] = [];
  let i = 0;

  while (i < lines.length) {
    const line = lines[i];

    if (line.startsWith('```')) {
      const lang = line.slice(3).trim();
      const codeLines: string[] = [];
      i++;
      while (i < lines.length && !lines[i].startsWith('```')) {
        codeLines.push(lines[i]);
        i++;
      }
      i++;
      const node: ADFNode = {
        type: 'codeBlock',
        content: [{ type: 'text', text: codeLines.join('\n') }],
      };
      if (lang) node.attrs = { language: lang };
      content.push(node);
      continue;
    }

    const headingMatch = line.match(/^(#{1,6})\s+(.*)/);
    if (headingMatch) {
      const level = headingMatch[1].length;
      content.push({
        type: 'heading',
        attrs: { level },
        content: parseInline(headingMatch[2]),
      });
      i++;
      continue;
    }

    const panelMatch = line.match(/^>\s*\[!(\w+)\]\s*(.*)/);
    if (panelMatch) {
      const rawType = panelMatch[1].toLowerCase();
      const panelType = PANEL_TYPES.has(rawType) ? rawType : 'info';
      const panelLines: string[] = [];
      const firstLine = panelMatch[2];
      if (firstLine.trim() !== '') panelLines.push(firstLine);
      i++;
      while (i < lines.length && lines[i].startsWith('>')) {
        panelLines.push(lines[i].replace(/^>\s?/, ''));
        i++;
      }
      content.push({
        type: 'panel',
        attrs: { panelType },
        content: [{ type: 'paragraph', content: parseInline(panelLines.join(' ')) }],
      });
      continue;
    }

    if (line.includes('|') && line.trim() !== '') {
      const next = lines[i + 1] ?? '';
      // A real Markdown table separator row must contain a `|` AND consist only of
      // pipes/dashes/colons/whitespace. Without the `|` requirement, a bullet list
      // line like `- foo` after a `|`-containing line gets misparsed as a table.
      const hasSeparator = next.includes('|') && /^[\s|:-]+$/.test(next);
      const nextIsTableRow = next.includes('|') && next.trim() !== '';
      if (hasSeparator || nextIsTableRow) {
        const tableRows: ADFNode[] = [];
        const headerCells = splitTableRow(line);
        tableRows.push({
          type: 'tableRow',
          content: headerCells.map((cell) => ({
            type: 'tableHeader',
            content: [{ type: 'paragraph', content: parseInline(cell.trim()) }],
          })),
        });
        i += hasSeparator ? 2 : 1;
        while (i < lines.length && lines[i].includes('|') && lines[i].trim() !== '') {
          const cells = splitTableRow(lines[i]);
          tableRows.push({
            type: 'tableRow',
            content: cells.map((cell) => ({
              type: 'tableCell',
              content: [{ type: 'paragraph', content: parseInline(cell.trim()) }],
            })),
          });
          i++;
        }
        content.push({ type: 'table', content: tableRows });
        continue;
      }
    }

    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line.trim())) {
      content.push({ type: 'rule' });
      i++;
      continue;
    }

    if (/^\s*[-*]\s+/.test(line)) {
      const items: ADFNode[] = [];
      while (i < lines.length && /^\s*[-*]\s+/.test(lines[i])) {
        const text = lines[i].replace(/^\s*[-*]\s+/, '');
        items.push({
          type: 'listItem',
          content: [{ type: 'paragraph', content: parseInline(text) }],
        });
        i++;
      }
      content.push({ type: 'bulletList', content: items });
      continue;
    }

    if (/^\s*\d+\.\s+/.test(line)) {
      const items: ADFNode[] = [];
      while (i < lines.length && /^\s*\d+\.\s+/.test(lines[i])) {
        const text = lines[i].replace(/^\s*\d+\.\s+/, '');
        items.push({
          type: 'listItem',
          content: [{ type: 'paragraph', content: parseInline(text) }],
        });
        i++;
      }
      content.push({ type: 'orderedList', content: items });
      continue;
    }

    if (line.trim() === '') {
      i++;
      continue;
    }

    const paraLines: string[] = [line];
    i++;
    while (
      i < lines.length &&
      lines[i].trim() !== '' &&
      !lines[i].startsWith('```') &&
      !/^#{1,6}\s/.test(lines[i]) &&
      !/^\s*[-*]\s+/.test(lines[i]) &&
      !/^\s*\d+\.\s+/.test(lines[i]) &&
      !/^>/.test(lines[i]) &&
      !/^(-{3,}|\*{3,}|_{3,})$/.test(lines[i].trim()) &&
      !lines[i].includes('|')
    ) {
      paraLines.push(lines[i]);
      i++;
    }
    const paraContent: ADFNode[] = [];
    paraLines.forEach((pl, idx) => {
      if (idx > 0) paraContent.push({ type: 'hardBreak' });
      paraContent.push(...parseInline(pl));
    });
    content.push({ type: 'paragraph', content: paraContent });
  }

  if (content.length === 0) {
    content.push({ type: 'paragraph', content: [{ type: 'text', text: '' }] });
  }

  return { type: 'doc', version: 1, content };
}

function splitTableRow(row: string): string[] {
  let s = row.trim();
  if (s.startsWith('|')) s = s.slice(1);
  if (s.endsWith('|')) s = s.slice(0, -1);
  return s.split('|');
}

function isWordChar(ch: string | undefined): boolean {
  return ch !== undefined && /[A-Za-z0-9]/.test(ch);
}

function parseInline(text: string): ADFNode[] {
  const nodes: ADFNode[] = [];
  let i = 0;
  const len = text.length;

  while (i < len) {
    if (text[i] === '`') {
      const end = text.indexOf('`', i + 1);
      if (end !== -1) {
        nodes.push({ type: 'text', text: text.slice(i + 1, end), marks: [{ type: 'code' }] });
        i = end + 1;
        continue;
      }
    }

    if (text.slice(i, i + 2) === '**') {
      const end = text.indexOf('**', i + 2);
      if (end !== -1 && end > i + 2) {
        nodes.push({ type: 'text', text: text.slice(i + 2, end), marks: [{ type: 'strong' }] });
        i = end + 2;
        continue;
      }
    }

    if (text.slice(i, i + 2) === '__') {
      const before = text[i - 1];
      if (!isWordChar(before)) {
        const end = findDoubleUnderscoreEnd(text, i + 2);
        if (end !== -1 && end > i + 2) {
          nodes.push({ type: 'text', text: text.slice(i + 2, end), marks: [{ type: 'strong' }] });
          i = end + 2;
          continue;
        }
      }
    }

    if (text[i] === '*' && text[i + 1] !== '*' && text[i - 1] !== '*') {
      const end = findItalicEndAsterisk(text, i + 1);
      if (end !== -1 && end > i + 1) {
        nodes.push({ type: 'text', text: text.slice(i + 1, end), marks: [{ type: 'em' }] });
        i = end + 1;
        continue;
      }
    }

    if (text[i] === '_' && text[i + 1] !== '_' && !isWordChar(text[i - 1])) {
      const end = findItalicEndUnderscore(text, i + 1);
      if (end !== -1 && end > i + 1) {
        nodes.push({ type: 'text', text: text.slice(i + 1, end), marks: [{ type: 'em' }] });
        i = end + 1;
        continue;
      }
    }

    if (text.slice(i, i + 9) === '@mention(') {
      const end = text.indexOf(')', i + 9);
      if (end !== -1) {
        const inner = text.slice(i + 9, end);
        const commaIdx = inner.indexOf(',');
        if (commaIdx !== -1) {
          const id = inner.slice(0, commaIdx).trim();
          const name = inner.slice(commaIdx + 1).trim();
          nodes.push({
            type: 'mention',
            attrs: { id, text: '@' + name },
          });
          i = end + 1;
          continue;
        }
      }
    }

    if (text[i] === '[') {
      const closeBracket = text.indexOf(']', i + 1);
      if (closeBracket !== -1 && text[closeBracket + 1] === '(') {
        const closeParen = text.indexOf(')', closeBracket + 2);
        if (closeParen !== -1) {
          const linkText = text.slice(i + 1, closeBracket);
          const href = text.slice(closeBracket + 2, closeParen);
          const marks: ADFMark[] = [{ type: 'link', attrs: { href } }];
          nodes.push({ type: 'text', text: linkText, marks });
          i = closeParen + 1;
          continue;
        }
      }
    }

    let end = i + 1;
    while (end < len) {
      const ch = text[end];
      if (ch === '`' || ch === '[' || ch === '@') break;
      if (text.slice(end, end + 2) === '**') break;
      if (text.slice(end, end + 2) === '__' && !isWordChar(text[end - 1])) break;
      if (ch === '*' && text[end + 1] !== '*' && text[end - 1] !== '*') break;
      if (ch === '_' && text[end + 1] !== '_' && !isWordChar(text[end - 1])) break;
      end++;
    }
    nodes.push({ type: 'text', text: text.slice(i, end) });
    i = end;
  }

  if (nodes.length === 0) {
    nodes.push({ type: 'text', text: '' });
  }

  return nodes;
}

function findItalicEndAsterisk(text: string, start: number): number {
  for (let j = start; j < text.length; j++) {
    if (text[j] === '*' && text[j + 1] !== '*' && text[j - 1] !== '*') {
      return j;
    }
  }
  return -1;
}

function findItalicEndUnderscore(text: string, start: number): number {
  for (let j = start; j < text.length; j++) {
    if (text[j] === '_' && text[j + 1] !== '_' && !isWordChar(text[j + 1])) {
      return j;
    }
  }
  return -1;
}

function findDoubleUnderscoreEnd(text: string, start: number): number {
  for (let j = start; j < text.length - 1; j++) {
    if (text[j] === '_' && text[j + 1] === '_' && !isWordChar(text[j + 2])) {
      return j;
    }
  }
  return -1;
}

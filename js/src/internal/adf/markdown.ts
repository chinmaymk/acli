import type { ADFDocument, ADFNode, ADFMark } from '../types.js';

/**
 * Converts a Markdown string to an Atlassian Document Format (ADF) document.
 * Handles headings, code blocks, tables, bullet/ordered lists, horizontal
 * rules, and inline marks (bold, code, links).
 */
export function markdownToADF(md: string): ADFDocument {
  const lines = md.split('\n');
  const content: ADFNode[] = [];
  let i = 0;

  while (i < lines.length) {
    const line = lines[i];

    // Fenced code block
    if (line.startsWith('```')) {
      const lang = line.slice(3).trim();
      const codeLines: string[] = [];
      i++;
      while (i < lines.length && !lines[i].startsWith('```')) {
        codeLines.push(lines[i]);
        i++;
      }
      i++; // skip closing ```
      const attrs: Record<string, string> = {};
      if (lang) attrs.language = lang;
      content.push({
        type: 'codeBlock',
        attrs,
        content: [{ type: 'text', text: codeLines.join('\n') }],
      });
      continue;
    }

    // Heading
    const headingMatch = line.match(/^(#{1,4})\s+(.*)/);
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

    // Table
    if (line.includes('|') && i + 1 < lines.length && /^\|?\s*[-:]+/.test(lines[i + 1])) {
      const tableRows: ADFNode[] = [];
      // Header row
      const headerCells = splitTableRow(line);
      tableRows.push({
        type: 'tableRow',
        content: headerCells.map((cell) => ({
          type: 'tableHeader',
          content: [{ type: 'paragraph', content: parseInline(cell.trim()) }],
        })),
      });
      i += 2; // skip header + separator
      // Data rows
      while (i < lines.length && lines[i].includes('|')) {
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

    // Horizontal rule
    if (/^---+$/.test(line.trim())) {
      content.push({ type: 'rule' });
      i++;
      continue;
    }

    // Bullet list
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

    // Ordered list
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

    // Empty line - skip
    if (line.trim() === '') {
      i++;
      continue;
    }

    // Paragraph
    content.push({ type: 'paragraph', content: parseInline(line) });
    i++;
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

/**
 * Parses inline markdown formatting into ADF text nodes with marks.
 */
function parseInline(text: string): ADFNode[] {
  const nodes: ADFNode[] = [];
  let i = 0;

  while (i < text.length) {
    // Inline code
    if (text[i] === '`') {
      const end = text.indexOf('`', i + 1);
      if (end !== -1) {
        const code = text.slice(i + 1, end);
        nodes.push({ type: 'text', text: code, marks: [{ type: 'code' }] });
        i = end + 1;
        continue;
      }
    }

    // Bold (**text**)
    if (text.slice(i, i + 2) === '**') {
      const end = text.indexOf('**', i + 2);
      if (end !== -1) {
        const bold = text.slice(i + 2, end);
        nodes.push({ type: 'text', text: bold, marks: [{ type: 'strong' }] });
        i = end + 2;
        continue;
      }
    }

    // Link [text](url)
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

    // Plain text - accumulate until next special char
    let end = i + 1;
    while (end < text.length && text[end] !== '`' && text[end] !== '[' &&
           !(text.slice(end, end + 2) === '**')) {
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

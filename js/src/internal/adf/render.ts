import type { ADFNode } from '../types.js';

export type { ADFNode };

/**
 * Renders an Atlassian Document Format tree (or string) to plain text.
 * Accepts either an ADF tree, a plain string (passed through), or null/undefined.
 */
export function render(doc: ADFNode | string | null | undefined): string {
  if (doc == null) {
    return '';
  }
  if (typeof doc === 'string') {
    return doc;
  }
  return renderNode(doc).replace(/\n+$/, '');
}

function renderNode(node: ADFNode): string {
  const nodeType = node.type;

  switch (nodeType) {
    case 'doc':
    case 'blockquote':
    case 'panel':
    case 'expand':
      return renderChildren(node);

    case 'paragraph':
    case 'heading':
    case 'codeBlock':
      return renderChildren(node) + '\n';

    case 'bulletList':
      return renderBulletList(node);

    case 'orderedList':
      return renderOrderedList(node);

    case 'listItem':
      return renderListItem(node);

    case 'text':
      return node.text ?? '';

    case 'hardBreak':
      return '\n';

    case 'rule':
      return '---\n';

    case 'table':
      return renderChildren(node);

    case 'tableRow':
      return renderTableRow(node);

    case 'tableHeader':
    case 'tableCell':
      return renderChildren(node).replace(/\n+$/, '');

    case 'mention':
      return nodeAttr(node, 'text') ?? '@unknown';

    case 'inlineCard':
      return nodeAttr(node, 'url') ?? '';

    case 'emoji': {
      const text = nodeAttr(node, 'text');
      if (text) return text;
      return nodeAttr(node, 'shortName') ?? '';
    }

    case 'date':
      return nodeAttr(node, 'timestamp') ?? '';

    case 'status':
      return nodeAttr(node, 'text') ?? '';

    case 'mediaSingle':
    case 'mediaInline':
    case 'media':
    case 'mediaGroup':
      return '';

    default:
      return renderChildren(node);
  }
}

function renderChildren(node: ADFNode): string {
  const content = node.content;
  if (!Array.isArray(content)) {
    return '';
  }
  return content
    .filter((child): child is ADFNode => child != null && typeof child === 'object')
    .map(child => renderNode(child))
    .join('');
}

function renderBulletList(node: ADFNode): string {
  const content = node.content;
  if (!Array.isArray(content)) {
    return '';
  }
  return content
    .filter((child): child is ADFNode => child != null && typeof child === 'object')
    .map(child => renderNode(child).replace(/\n+$/, '') + '\n')
    .join('');
}

function renderOrderedList(node: ADFNode): string {
  const content = node.content;
  if (!Array.isArray(content)) {
    return '';
  }
  return content
    .filter((child): child is ADFNode => child != null && typeof child === 'object')
    .map((child, i) => `${i + 1}. ${renderNode(child).replace(/\n+$/, '')}\n`)
    .join('');
}

function renderListItem(node: ADFNode): string {
  const content = node.content;
  if (!Array.isArray(content)) {
    return '';
  }
  const parts: string[] = [];
  for (const child of content) {
    if (child != null) {
      const text = renderNode(child).replace(/\n+$/, '');
      if (text) {
        parts.push(text);
      }
    }
  }
  return parts.join('\n');
}

function renderTableRow(node: ADFNode): string {
  const content = node.content;
  if (!Array.isArray(content)) {
    return '';
  }
  const cells = content
    .filter((child): child is ADFNode => child != null && typeof child === 'object')
    .map(child => renderNode(child));
  return cells.join('\t') + '\n';
}

function nodeAttr(node: ADFNode, key: string): string | undefined {
  if (!node.attrs) {
    return undefined;
  }
  const val = node.attrs[key];
  return typeof val === 'string' ? val : undefined;
}

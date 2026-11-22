import { render } from '../internal/adf/render.js';
import type { IssueDetailed } from '../internal/jira/types.js';

interface SearchFilters {
  project?: string;
  assignee?: string;
  mine?: boolean;
  updatedSince?: string;
}

function searchWords(value: string): string[] {
  return value.match(/[\p{L}\p{N}][\p{L}\p{M}\p{N}]*/gu) || [];
}

export function textSearchJQL(text: string, filters: SearchFilters): string {
  for (const [flag, value] of [['project', filters.project], ['assignee', filters.assignee], ['updated-since', filters.updatedSince]]) {
    if (value !== undefined && (typeof value !== 'string' || !value.trim())) {
      throw new Error(`--${flag} requires one nonempty value`);
    }
  }
  const clauses: string[] = [];
  if (text) {
    if (/^[a-z][a-z0-9_]*-[0-9]+$/i.test(text)) {
      clauses.push('key = ' + JSON.stringify(text.toUpperCase()));
    } else {
      const words = searchWords(text);
      if (!words.length) throw new Error('search text must contain a letter or number');
      for (const word of words) {
        // Quote the word for both Jira's text parser and JQL.
        clauses.push('text ~ ' + JSON.stringify(JSON.stringify(word)));
      }
    }
  }
  if (filters.project) clauses.push('project = ' + JSON.stringify(filters.project));
  if (filters.mine && filters.assignee) throw new Error('choose either --mine or --assignee');
  if (filters.mine) clauses.push('assignee = currentUser()');
  if (filters.assignee) clauses.push('assignee = ' + JSON.stringify(filters.assignee));
  if (filters.updatedSince) {
    if (!/^[1-9][0-9]*[dw]$/.test(filters.updatedSince)) {
      throw new Error('--updated-since requires a positive number of days or weeks (e.g. 30d or 2w)');
    }
    clauses.push('updated >= ' + JSON.stringify('-' + filters.updatedSince));
  }
  if (!clauses.length) throw new Error('provide search text, a filter, or --jql');
  return clauses.join(' AND ');
}

export function issuePreviews(issues: IssueDetailed[], text: string, baseURL: string) {
  const terms = searchWords(text.toLowerCase());
  const matches = (value: string) => terms.some(term => value.toLowerCase().includes(term));
  const previews = issues.map(issue => {
    const body = render(issue.fields?.description).split(/\s+/u).filter(Boolean);
    const index = body.findIndex(matches);
    const excerpt = index < 0 ? undefined : Array.from(
      body.slice(Math.max(0, index - 5), index + 14).join(' ')
    ).slice(0, 160).join('');
    const key = issue.key || '';
    return { key, summary: issue.fields?.summary || '', status: issue.fields?.status?.name || '',
      url: baseURL + '/browse/' + encodeURIComponent(key), excerpt };
  });
  return previews.sort((a, b) => Number(matches(b.summary)) - Number(matches(a.summary)));
}

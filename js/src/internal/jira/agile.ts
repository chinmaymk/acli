import { get, post, put, del } from './client.js';
import type { JiraClient } from './client.js';

// --- Boards ---

/**
 * Returns a paginated list of boards.
 */
export async function getBoards(
  client: JiraClient,
  startAt: number,
  maxResults: number,
  projectKeyOrID?: string,
  boardType?: string,
  name?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (projectKeyOrID) query.projectKeyOrId = projectKeyOrID;
  if (boardType) query.type = boardType;
  if (name) query.name = name;
  return get(client, '/rest/agile/1.0/board', query as Record<string, string>);
}

/**
 * Returns a board by ID.
 */
export async function getBoard(client: JiraClient, boardID: number): Promise<any> {
  return get(client, `/rest/agile/1.0/board/${boardID}`);
}

/**
 * Returns a board's configuration.
 */
export async function getBoardConfiguration(client: JiraClient, boardID: number): Promise<any> {
  return get(client, `/rest/agile/1.0/board/${boardID}/configuration`);
}

/**
 * Returns issues for a board.
 */
export async function getBoardIssues(
  client: JiraClient,
  boardID: number,
  startAt: number,
  maxResults: number,
  jql?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/board/${boardID}/issue`, query as Record<string, string>);
}

/**
 * Returns backlog issues for a board.
 */
export async function getBoardBacklog(
  client: JiraClient,
  boardID: number,
  startAt: number,
  maxResults: number,
  jql?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/board/${boardID}/backlog`, query as Record<string, string>);
}

/**
 * Returns sprints for a board.
 */
export async function getBoardSprints(
  client: JiraClient,
  boardID: number,
  startAt: number,
  maxResults: number,
  state?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (state) query.state = state;
  return get(client, `/rest/agile/1.0/board/${boardID}/sprint`, query as Record<string, string>);
}

/**
 * Returns epics for a board.
 */
export async function getBoardEpics(
  client: JiraClient,
  boardID: number,
  startAt: number,
  maxResults: number,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  return get(client, `/rest/agile/1.0/board/${boardID}/epic`, query as Record<string, string>);
}

// --- Sprints ---

/**
 * Returns a sprint by ID.
 */
export async function getSprint(client: JiraClient, sprintID: number): Promise<any> {
  return get(client, `/rest/agile/1.0/sprint/${sprintID}`);
}

/**
 * Creates a new sprint.
 */
export async function createSprint(client: JiraClient, sprint: unknown): Promise<any> {
  return post(client, '/rest/agile/1.0/sprint', sprint);
}

/**
 * Updates a sprint.
 */
export async function updateSprint(
  client: JiraClient,
  sprintID: number,
  sprint: unknown,
): Promise<any> {
  return put(client, `/rest/agile/1.0/sprint/${sprintID}`, sprint);
}

/**
 * Partially updates a sprint.
 */
export async function partialUpdateSprint(
  client: JiraClient,
  sprintID: number,
  sprint: unknown,
): Promise<any> {
  return post(client, `/rest/agile/1.0/sprint/${sprintID}`, sprint);
}

/**
 * Deletes a sprint.
 */
export async function deleteSprint(client: JiraClient, sprintID: number): Promise<void> {
  return del(client, `/rest/agile/1.0/sprint/${sprintID}`);
}

/**
 * Returns issues in a sprint.
 */
export async function getSprintIssues(
  client: JiraClient,
  sprintID: number,
  startAt: number,
  maxResults: number,
  jql?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/sprint/${sprintID}/issue`, query as Record<string, string>);
}

/**
 * Moves issues to a sprint.
 */
export async function moveIssuesToSprint(
  client: JiraClient,
  sprintID: number,
  issueKeys: string[],
): Promise<any> {
  return post(client, `/rest/agile/1.0/sprint/${sprintID}/issue`, { issues: issueKeys });
}

/**
 * Moves issues to the backlog.
 */
export async function moveIssuesToBacklog(
  client: JiraClient,
  issueKeys: string[],
): Promise<any> {
  return post(client, '/rest/agile/1.0/backlog/issue', { issues: issueKeys });
}

// --- Epics ---

/**
 * Returns an epic by ID or key.
 */
export async function getEpic(client: JiraClient, epicIdOrKey: string): Promise<any> {
  return get(client, `/rest/agile/1.0/epic/${epicIdOrKey}`);
}

/**
 * Moves issues to an epic.
 */
export async function moveIssuesToEpic(
  client: JiraClient,
  epicIdOrKey: string,
  issueKeys: string[],
): Promise<any> {
  return post(client, `/rest/agile/1.0/epic/${epicIdOrKey}/issue`, { issues: issueKeys });
}

/**
 * Returns issues belonging to an epic.
 */
export async function getEpicIssues(
  client: JiraClient,
  epicIdOrKey: string,
  startAt: number,
  maxResults: number,
  jql?: string,
): Promise<any> {
  const query: Record<string, string | number> = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/epic/${epicIdOrKey}/issue`, query as Record<string, string>);
}

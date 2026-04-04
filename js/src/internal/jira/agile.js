import { get, post, put, del } from './client.js';

// --- Boards ---

/**
 * Returns a paginated list of boards.
 * @param {object} client
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [projectKeyOrID]
 * @param {string} [boardType]
 * @param {string} [name]
 * @returns {Promise<object>}
 */
export async function getBoards(client, startAt, maxResults, projectKeyOrID, boardType, name) {
  const query = { startAt, maxResults };
  if (projectKeyOrID) query.projectKeyOrId = projectKeyOrID;
  if (boardType) query.type = boardType;
  if (name) query.name = name;
  return get(client, '/rest/agile/1.0/board', query);
}

/**
 * Returns a board by ID.
 * @param {object} client
 * @param {number} boardID
 * @returns {Promise<object>}
 */
export async function getBoard(client, boardID) {
  return get(client, `/rest/agile/1.0/board/${boardID}`);
}

/**
 * Returns a board's configuration.
 * @param {object} client
 * @param {number} boardID
 * @returns {Promise<object>}
 */
export async function getBoardConfiguration(client, boardID) {
  return get(client, `/rest/agile/1.0/board/${boardID}/configuration`);
}

/**
 * Returns issues for a board.
 * @param {object} client
 * @param {number} boardID
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [jql]
 * @returns {Promise<object>}
 */
export async function getBoardIssues(client, boardID, startAt, maxResults, jql) {
  const query = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/board/${boardID}/issue`, query);
}

/**
 * Returns backlog issues for a board.
 * @param {object} client
 * @param {number} boardID
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [jql]
 * @returns {Promise<object>}
 */
export async function getBoardBacklog(client, boardID, startAt, maxResults, jql) {
  const query = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/board/${boardID}/backlog`, query);
}

/**
 * Returns sprints for a board.
 * @param {object} client
 * @param {number} boardID
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [state]
 * @returns {Promise<object>}
 */
export async function getBoardSprints(client, boardID, startAt, maxResults, state) {
  const query = { startAt, maxResults };
  if (state) query.state = state;
  return get(client, `/rest/agile/1.0/board/${boardID}/sprint`, query);
}

/**
 * Returns epics for a board.
 * @param {object} client
 * @param {number} boardID
 * @param {number} startAt
 * @param {number} maxResults
 * @returns {Promise<object>}
 */
export async function getBoardEpics(client, boardID, startAt, maxResults) {
  const query = { startAt, maxResults };
  return get(client, `/rest/agile/1.0/board/${boardID}/epic`, query);
}

// --- Sprints ---

/**
 * Returns a sprint by ID.
 * @param {object} client
 * @param {number} sprintID
 * @returns {Promise<object>}
 */
export async function getSprint(client, sprintID) {
  return get(client, `/rest/agile/1.0/sprint/${sprintID}`);
}

/**
 * Creates a new sprint.
 * @param {object} client
 * @param {object} sprint
 * @returns {Promise<object>}
 */
export async function createSprint(client, sprint) {
  return post(client, '/rest/agile/1.0/sprint', sprint);
}

/**
 * Updates a sprint.
 * @param {object} client
 * @param {number} sprintID
 * @param {object} sprint
 * @returns {Promise<object>}
 */
export async function updateSprint(client, sprintID, sprint) {
  return put(client, `/rest/agile/1.0/sprint/${sprintID}`, sprint);
}

/**
 * Partially updates a sprint.
 * @param {object} client
 * @param {number} sprintID
 * @param {object} sprint
 * @returns {Promise<object>}
 */
export async function partialUpdateSprint(client, sprintID, sprint) {
  return post(client, `/rest/agile/1.0/sprint/${sprintID}`, sprint);
}

/**
 * Deletes a sprint.
 * @param {object} client
 * @param {number} sprintID
 * @returns {Promise<null>}
 */
export async function deleteSprint(client, sprintID) {
  return del(client, `/rest/agile/1.0/sprint/${sprintID}`);
}

/**
 * Returns issues in a sprint.
 * @param {object} client
 * @param {number} sprintID
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [jql]
 * @returns {Promise<object>}
 */
export async function getSprintIssues(client, sprintID, startAt, maxResults, jql) {
  const query = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/sprint/${sprintID}/issue`, query);
}

/**
 * Moves issues to a sprint.
 * @param {object} client
 * @param {number} sprintID
 * @param {string[]} issueKeys
 * @returns {Promise<null>}
 */
export async function moveIssuesToSprint(client, sprintID, issueKeys) {
  return post(client, `/rest/agile/1.0/sprint/${sprintID}/issue`, { issues: issueKeys });
}

/**
 * Moves issues to the backlog.
 * @param {object} client
 * @param {string[]} issueKeys
 * @returns {Promise<null>}
 */
export async function moveIssuesToBacklog(client, issueKeys) {
  return post(client, '/rest/agile/1.0/backlog/issue', { issues: issueKeys });
}

// --- Epics ---

/**
 * Returns an epic by ID or key.
 * @param {object} client
 * @param {string} epicIdOrKey
 * @returns {Promise<object>}
 */
export async function getEpic(client, epicIdOrKey) {
  return get(client, `/rest/agile/1.0/epic/${epicIdOrKey}`);
}

/**
 * Moves issues to an epic.
 * @param {object} client
 * @param {string} epicIdOrKey
 * @param {string[]} issueKeys
 * @returns {Promise<null>}
 */
export async function moveIssuesToEpic(client, epicIdOrKey, issueKeys) {
  return post(client, `/rest/agile/1.0/epic/${epicIdOrKey}/issue`, { issues: issueKeys });
}

/**
 * Returns issues belonging to an epic.
 * @param {object} client
 * @param {string} epicIdOrKey
 * @param {number} startAt
 * @param {number} maxResults
 * @param {string} [jql]
 * @returns {Promise<object>}
 */
export async function getEpicIssues(client, epicIdOrKey, startAt, maxResults, jql) {
  const query = { startAt, maxResults };
  if (jql) query.jql = jql;
  return get(client, `/rest/agile/1.0/epic/${epicIdOrKey}/issue`, query);
}

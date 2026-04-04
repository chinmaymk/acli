import { getConfluenceClient } from './helpers.js';
import { confluenceV2 } from '../internal/api/client.js';

// ---------------------------------------------------------------------------
// Pagination helpers
// ---------------------------------------------------------------------------

const DEFAULT_LIMIT = 50;

/**
 * Build a query object from common pagination/sort/format flags in argv.
 * @param {object} argv
 * @returns {object}
 */
function getPaginationQuery(argv) {
  const q = {};
  if (argv.limit != null && argv.limit > 0) q.limit = String(argv.limit);
  if (argv.cursor && !argv.all) q.cursor = argv.cursor;
  if (argv.sort) q.sort = argv.sort;
  if (argv['body-format']) q['body-format'] = argv['body-format'];
  if (argv.status) {
    const statuses = Array.isArray(argv.status) ? argv.status : [argv.status];
    if (statuses.length > 0) q.status = statuses;
  }
  return q;
}

/**
 * Fetch a single page of results from the Confluence v2 API.
 */
async function confGet(client, path, query) {
  return confluenceV2(client, 'GET', path, query, null);
}

/**
 * POST to the Confluence v2 API.
 */
async function confPost(client, path, query, body) {
  return confluenceV2(client, 'POST', path, query, body);
}

/**
 * PUT to the Confluence v2 API.
 */
async function confPut(client, path, query, body) {
  return confluenceV2(client, 'PUT', path, query, body);
}

/**
 * DELETE from the Confluence v2 API.
 */
async function confDelete(client, path, query) {
  return confluenceV2(client, 'DELETE', path, query, null);
}

/**
 * Fetch paginated results, following cursor links when all=true.
 * @param {object} client
 * @param {string} path
 * @param {object} query
 * @param {boolean} all
 * @returns {Promise<object>}
 */
async function confGetPaginated(client, path, query, all) {
  if (!all) {
    return confGet(client, path, query);
  }

  const allResults = [];
  let currentPath = path;
  let currentQuery = query;

  for (;;) {
    const data = await confluenceV2(client, 'GET', currentPath, currentQuery, null);
    if (!data || !data.results) {
      return data;
    }
    allResults.push(...data.results);
    if (!data._links || !data._links.next) break;

    // Strip /wiki/api/v2 prefix since confluenceV2 adds it
    let nextURL = data._links.next;
    const prefix = '/wiki/api/v2';
    const idx = nextURL.indexOf(prefix);
    if (idx >= 0) nextURL = nextURL.slice(idx + prefix.length);

    const qIdx = nextURL.indexOf('?');
    if (qIdx >= 0) {
      const params = new URLSearchParams(nextURL.slice(qIdx + 1));
      currentQuery = Object.fromEntries(params.entries());
      currentPath = nextURL.slice(0, qIdx);
    } else {
      currentPath = nextURL;
      currentQuery = {};
    }
  }

  return { results: allResults };
}

/**
 * Pretty-print a JSON value to stdout.
 */
function printJSON(data) {
  console.log(JSON.stringify(data, null, 2));
}

/**
 * Parse a string as JSON; if it fails treat it as a plain string.
 */
function parseJSONOrString(s) {
  try {
    return JSON.parse(s);
  } catch {
    return s;
  }
}

// ---------------------------------------------------------------------------
// Shared option builders
// ---------------------------------------------------------------------------

function addPaginationOptions(yargs) {
  return yargs
    .option('limit', { type: 'number', default: DEFAULT_LIMIT, describe: 'Maximum number of results to return' })
    .option('cursor', { type: 'string', describe: 'Pagination cursor' })
    .option('all', { type: 'boolean', default: false, describe: 'Fetch all pages of results (follows pagination cursors)' });
}

function addSortOption(yargs) {
  return yargs.option('sort', { type: 'string', describe: 'Sort field' });
}

function addBodyFormatOption(yargs) {
  return yargs.option('body-format', {
    type: 'string',
    describe: 'Body format (storage, atlas_doc_format, view, export_view, anonymous_export_view, styled_view, editor)',
  });
}

function addStatusOption(yargs) {
  return yargs.option('status', { type: 'array', describe: 'Filter by status' });
}

function addIncludeTreeOptions(yargs) {
  return yargs
    .option('include-collaborators', { type: 'boolean', default: false, describe: 'Include collaborators' })
    .option('include-direct-children', { type: 'boolean', default: false, describe: 'Include direct children' })
    .option('include-operations', { type: 'boolean', default: false, describe: 'Include operations' })
    .option('include-properties', { type: 'boolean', default: false, describe: 'Include properties' });
}

// ---------------------------------------------------------------------------
// addTreeSubResources: ancestors, descendants, direct-children, operations, properties
// Used by whiteboard, database, folder, smart-link
// ---------------------------------------------------------------------------

function addTreeSubResources(yargs, pathPrefix, resourceName) {
  return yargs
    .command(
      'ancestors <id>',
      `Get all ancestors of ${resourceName}`,
      (y) =>
        y
          .positional('id', { type: 'string' })
          .option('limit', { type: 'number', default: DEFAULT_LIMIT, describe: 'Maximum number of results' }),
      async (argv) => {
        const client = getConfluenceClient(argv);
        const q = {};
        if (argv.limit > 0) q.limit = String(argv.limit);
        const data = await confGet(client, `${pathPrefix}/${argv.id}/ancestors`, q);
        printJSON(data);
      }
    )
    .command(
      'descendants <id>',
      `Get descendants of a ${resourceName}`,
      (y) => {
        y = addPaginationOptions(y.positional('id', { type: 'string' }));
        y = addSortOption(y);
        return y.option('depth', { type: 'number', default: 0, describe: 'Maximum depth of descendants' });
      },
      async (argv) => {
        const client = getConfluenceClient(argv);
        const q = getPaginationQuery(argv);
        if (argv.depth > 0) q.depth = String(argv.depth);
        const data = await confGetPaginated(client, `${pathPrefix}/${argv.id}/descendants`, q, argv.all);
        printJSON(data);
      }
    )
    .command(
      'direct-children <id>',
      `Get direct children of a ${resourceName}`,
      (y) => addSortOption(addPaginationOptions(y.positional('id', { type: 'string' }))),
      async (argv) => {
        const client = getConfluenceClient(argv);
        const q = getPaginationQuery(argv);
        const data = await confGetPaginated(client, `${pathPrefix}/${argv.id}/direct-children`, q, argv.all);
        printJSON(data);
      }
    )
    .command(
      'operations <id>',
      `Get permitted operations for ${resourceName}`,
      (y) => y.positional('id', { type: 'string' }),
      async (argv) => {
        const client = getConfluenceClient(argv);
        const data = await confGet(client, `${pathPrefix}/${argv.id}/operations`, null);
        printJSON(data);
      }
    )
    .command(
      'properties <id>',
      `Get content properties for ${resourceName}`,
      (y) => {
        y = addPaginationOptions(y.positional('id', { type: 'string' }));
        y = addSortOption(y);
        return y.option('key', { type: 'string', describe: 'Filter by property key' });
      },
      async (argv) => {
        const client = getConfluenceClient(argv);
        const q = getPaginationQuery(argv);
        if (argv.key) q.key = argv.key;
        const data = await confGetPaginated(client, `${pathPrefix}/${argv.id}/properties`, q, argv.all);
        printJSON(data);
      }
    );
}

// ---------------------------------------------------------------------------
// registerConfluenceCommands
// ---------------------------------------------------------------------------

export function registerConfluenceCommands(yargs) {
  return yargs
        // ----------------------------------------------------------------
        // space
        // ----------------------------------------------------------------
        .command(
          ['space', 's'],
          'Manage spaces',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List spaces',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return y
                    .option('ids', { type: 'array', describe: 'Filter by space IDs' })
                    .option('keys', { type: 'array', describe: 'Filter by space keys' })
                    .option('type', { type: 'string', describe: 'Filter by type (global, personal)' })
                    .option('labels', { type: 'array', describe: 'Filter by labels' })
                    .option('favorited-by', { type: 'string', describe: 'Filter by favorited-by user account ID' })
                    .option('not-favorited-by', { type: 'string', describe: 'Filter by not-favorited-by user account ID' })
                    .option('description-format', { type: 'string', describe: 'Description format (plain, view)' })
                    .option('include-icon', { type: 'boolean', default: false, describe: 'Include space icon' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.ids && argv.ids.length > 0) q.ids = argv.ids;
                  if (argv.keys && argv.keys.length > 0) q.keys = argv.keys;
                  if (argv.type) q.type = argv.type;
                  if (argv.labels && argv.labels.length > 0) q.labels = argv.labels;
                  if (argv['favorited-by']) q['favorited-by'] = argv['favorited-by'];
                  if (argv['not-favorited-by']) q['not-favorited-by'] = argv['not-favorited-by'];
                  if (argv['description-format']) q['description-format'] = argv['description-format'];
                  if (argv['include-icon']) q['include-icon'] = 'true';
                  const data = await confGetPaginated(client, '/spaces', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <space-id>',
                'Get space by ID',
                (y) =>
                  y
                    .positional('space-id', { type: 'string' })
                    .option('description-format', { type: 'string', describe: 'Description format' })
                    .option('include-icon', { type: 'boolean', default: false, describe: 'Include icon' })
                    .option('include-operations', { type: 'boolean', default: false, describe: 'Include operations' })
                    .option('include-properties', { type: 'boolean', default: false, describe: 'Include properties' })
                    .option('include-permissions', { type: 'boolean', default: false, describe: 'Include permissions' })
                    .option('include-role-assignments', { type: 'boolean', default: false, describe: 'Include role assignments' })
                    .option('include-labels', { type: 'boolean', default: false, describe: 'Include labels' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv['description-format']) q['description-format'] = argv['description-format'];
                  if (argv['include-icon']) q['include-icon'] = 'true';
                  if (argv['include-operations']) q['include-operations'] = 'true';
                  if (argv['include-properties']) q['include-properties'] = 'true';
                  if (argv['include-permissions']) q['include-permissions'] = 'true';
                  if (argv['include-role-assignments']) q['include-role-assignments'] = 'true';
                  if (argv['include-labels']) q['include-labels'] = 'true';
                  const data = await confGet(client, `/spaces/${argv['space-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'create',
                'Create a space',
                (y) =>
                  y
                    .option('name', { type: 'string', demandOption: true, describe: 'Space name' })
                    .option('key', { type: 'string', describe: 'Space key' })
                    .option('alias', { type: 'string', describe: 'Space alias' })
                    .option('description', { type: 'string', describe: 'Space description' })
                    .option('private', { type: 'boolean', default: false, describe: 'Create as private space' })
                    .option('template-key', { type: 'string', describe: 'Template key' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = { name: argv.name };
                  if (argv.key) body.key = argv.key;
                  if (argv.alias) body.alias = argv.alias;
                  if (argv.description) body.description = { representation: 'plain', value: argv.description };
                  if (argv.private) body.createPrivateSpace = true;
                  if (argv['template-key']) body.templateKey = argv['template-key'];
                  const data = await confPost(client, '/spaces', null, body);
                  printJSON(data);
                }
              )
              .command(
                'update <space-id>',
                'Update a space',
                (y) =>
                  y
                    .positional('space-id', { type: 'string' })
                    .option('name', { type: 'string', describe: 'Space name' })
                    .option('description', { type: 'string', describe: 'Space description' })
                    .option('homepage-id', { type: 'string', describe: 'Homepage page ID' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = {};
                  if (argv.name) body.name = argv.name;
                  if (argv.description) body.description = { representation: 'plain', value: argv.description };
                  if (argv['homepage-id']) body.homepageId = argv['homepage-id'];
                  const data = await confPut(client, `/spaces/${argv['space-id']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <space-id>',
                'Delete a space',
                (y) => y.positional('space-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/spaces/${argv['space-id']}`, null);
                  console.log('Space deleted successfully.');
                }
              )
              .command(
                'pages <space-id>',
                'List pages in a space',
                (y) => {
                  y = addPaginationOptions(y.positional('space-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y
                    .option('depth', { type: 'string', describe: 'Filter by depth (root, all)' })
                    .option('title', { type: 'string', describe: 'Filter by title' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.depth) q.depth = argv.depth;
                  if (argv.title) q.title = argv.title;
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/pages`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'blogposts <space-id>',
                'List blog posts in a space',
                (y) => {
                  y = addPaginationOptions(y.positional('space-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('title', { type: 'string', describe: 'Filter by title' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  q['space-id'] = argv['space-id'];
                  if (argv.title) q.title = argv.title;
                  const data = await confGetPaginated(client, '/blogposts', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'labels <space-id>',
                'Get labels for a space',
                (y) => addPaginationOptions(y.positional('space-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'content-labels <space-id>',
                'Get labels for space content',
                (y) =>
                  addPaginationOptions(y.positional('space-id', { type: 'string' }))
                    .option('prefix', { type: 'string', describe: 'Filter by prefix' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.prefix) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/content/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'custom-content <space-id>',
                'Get custom content by type in space',
                (y) => {
                  y = addPaginationOptions(y.positional('space-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('type', { type: 'string', describe: 'Custom content type' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.type) q.type = argv.type;
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/custom-content`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'operations <space-id>',
                'Get permitted operations for space',
                (y) => y.positional('space-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/spaces/${argv['space-id']}/operations`, null);
                  printJSON(data);
                }
              )
              .command(
                'permissions <space-id>',
                'Get space permissions assignments',
                (y) => addPaginationOptions(y.positional('space-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.cursor) q.cursor = argv.cursor;
                  if (argv.limit > 0) q.limit = String(argv.limit);
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/permissions`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'role-assignments <space-id>',
                'Get space role assignments',
                (y) =>
                  addPaginationOptions(y.positional('space-id', { type: 'string' }))
                    .option('role-id', { type: 'string', describe: 'Filter by role ID' })
                    .option('role-type', { type: 'string', describe: 'Filter by role type' })
                    .option('principal-id', { type: 'string', describe: 'Filter by principal ID' })
                    .option('principal-type', { type: 'string', describe: 'Filter by principal type' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.cursor) q.cursor = argv.cursor;
                  if (argv.limit > 0) q.limit = String(argv.limit);
                  if (argv['role-id']) q['role-id'] = argv['role-id'];
                  if (argv['role-type']) q['role-type'] = argv['role-type'];
                  if (argv['principal-id']) q['principal-id'] = argv['principal-id'];
                  if (argv['principal-type']) q['principal-type'] = argv['principal-type'];
                  const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/role-assignments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'set-role-assignments <space-id>',
                'Set space role assignments (provide JSON body via --body flag)',
                (y) =>
                  y
                    .positional('space-id', { type: 'string' })
                    .option('body', { type: 'string', demandOption: true, describe: 'JSON body for role assignments' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPost(client, `/spaces/${argv['space-id']}/role-assignments`, null, body);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a space subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // page
        // ----------------------------------------------------------------
        .command(
          ['page', 'p'],
          'Manage pages',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List pages',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y
                    .option('id', { type: 'array', describe: 'Filter by page IDs' })
                    .option('space-id', { type: 'array', describe: 'Filter by space IDs' })
                    .option('title', { type: 'string', describe: 'Filter by title' })
                    .option('subtype', { type: 'string', describe: 'Filter by subtype' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.id && argv.id.length > 0) q.id = argv.id;
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  if (argv.title) q.title = argv.title;
                  if (argv.subtype) q.subtype = argv.subtype;
                  const data = await confGetPaginated(client, '/pages', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <page-id>',
                'Get page by ID',
                (y) => {
                  y = addBodyFormatOption(y.positional('page-id', { type: 'string' }));
                  y = addStatusOption(y);
                  return y
                    .option('get-draft', { type: 'boolean', default: false, describe: 'Retrieve draft version' })
                    .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                    .option('include-labels', { type: 'boolean', default: false, describe: 'Include labels' })
                    .option('include-properties', { type: 'boolean', default: false, describe: 'Include properties' })
                    .option('include-operations', { type: 'boolean', default: false, describe: 'Include operations' })
                    .option('include-likes', { type: 'boolean', default: false, describe: 'Include likes' })
                    .option('include-versions', { type: 'boolean', default: false, describe: 'Include versions' })
                    .option('include-version', { type: 'boolean', default: false, describe: 'Include current version' })
                    .option('include-favorited-by-current-user-status', { type: 'boolean', default: false, describe: 'Include favorited status' })
                    .option('include-webresources', { type: 'boolean', default: false, describe: 'Include web resources' })
                    .option('include-collaborators', { type: 'boolean', default: false, describe: 'Include collaborators' })
                    .option('include-direct-children', { type: 'boolean', default: false, describe: 'Include direct children' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv['body-format']) q['body-format'] = argv['body-format'];
                  if (argv['get-draft']) q['get-draft'] = 'true';
                  if (argv.version > 0) q.version = String(argv.version);
                  for (const flag of [
                    'include-labels', 'include-properties', 'include-operations',
                    'include-likes', 'include-versions', 'include-version',
                    'include-favorited-by-current-user-status', 'include-webresources',
                    'include-collaborators', 'include-direct-children',
                  ]) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/pages/${argv['page-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'create',
                'Create a page',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Page title' })
                    .option('parent-id', { type: 'string', describe: 'Parent page ID' })
                    .option('status', { type: 'string', describe: 'Page status (current, draft)' })
                    .option('body', { type: 'string', describe: 'Page body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                    .option('subtype', { type: 'string', describe: 'Page subtype' })
                    .option('embedded', { type: 'boolean', default: false, describe: 'Create as embedded content' })
                    .option('private', { type: 'boolean', default: false, describe: 'Create as private page' })
                    .option('root-level', { type: 'boolean', default: false, describe: 'Create at root level of space' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.embedded) q.embedded = 'true';
                  if (argv.private) q.private = 'true';
                  if (argv['root-level']) q['root-level'] = 'true';
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv.status) body.status = argv.status;
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  if (argv.subtype) body.subtype = argv.subtype;
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  const data = await confPost(client, '/pages', q, body);
                  printJSON(data);
                }
              )
              .command(
                'update <page-id>',
                'Update a page',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .option('title', { type: 'string', demandOption: true, describe: 'Page title' })
                    .option('status', { type: 'string', default: 'current', describe: 'Page status' })
                    .option('body', { type: 'string', describe: 'Page body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                    .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                    .option('version-message', { type: 'string', describe: 'Version message' })
                    .option('space-id', { type: 'string', describe: 'Space ID' })
                    .option('parent-id', { type: 'string', describe: 'Parent page ID' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = {
                    id: argv['page-id'],
                    status: argv.status,
                    title: argv.title,
                    version: { number: argv['version-number'], message: argv['version-message'] || '' },
                  };
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  if (argv['space-id']) body.spaceId = argv['space-id'];
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  const data = await confPut(client, `/pages/${argv['page-id']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'update-title <page-id>',
                'Update page title only',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .option('title', { type: 'string', demandOption: true, describe: 'New title' })
                    .option('status', { type: 'string', default: 'current', describe: 'Page status' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = { title: argv.title, status: argv.status };
                  const data = await confPut(client, `/pages/${argv['page-id']}/title`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <page-id>',
                'Delete a page',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .option('purge', { type: 'boolean', default: false, describe: 'Purge the page' })
                    .option('draft', { type: 'boolean', default: false, describe: 'Delete a draft page' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.purge) q.purge = 'true';
                  if (argv.draft) q.draft = 'true';
                  await confDelete(client, `/pages/${argv['page-id']}`, q);
                  console.log('Page deleted successfully.');
                }
              )
              .command(
                'children <page-id>',
                'Get child pages',
                (y) => addSortOption(addPaginationOptions(y.positional('page-id', { type: 'string' }))),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/children`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'direct-children <page-id>',
                'Get direct children of a page',
                (y) => addSortOption(addPaginationOptions(y.positional('page-id', { type: 'string' }))),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/direct-children`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'ancestors <page-id>',
                'Get all ancestors of page',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .option('limit', { type: 'number', default: DEFAULT_LIMIT }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.limit > 0) q.limit = String(argv.limit);
                  const data = await confGet(client, `/pages/${argv['page-id']}/ancestors`, q);
                  printJSON(data);
                }
              )
              .command(
                'descendants <page-id>',
                'Get descendants of page',
                (y) =>
                  addPaginationOptions(y.positional('page-id', { type: 'string' }))
                    .option('depth', { type: 'number', default: 0, describe: 'Maximum depth of descendants' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.depth > 0) q.depth = String(argv.depth);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/descendants`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'versions <page-id>',
                'Get page versions',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/versions`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'version-details <page-id> <version-number>',
                'Get version details for page version',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .positional('version-number', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/pages/${argv['page-id']}/versions/${argv['version-number']}`, null);
                  printJSON(data);
                }
              )
              .command(
                'labels <page-id>',
                'Get labels for page',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  return y.option('prefix', { type: 'string', describe: 'Filter by prefix' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.prefix) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'attachments <page-id>',
                'Get attachments for page',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return y
                    .option('media-type', { type: 'string', describe: 'Filter by media type' })
                    .option('filename', { type: 'string', describe: 'Filter by filename' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['media-type']) q.mediaType = argv['media-type'];
                  if (argv.filename) q.filename = argv.filename;
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/attachments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'footer-comments <page-id>',
                'Get footer comments for page',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/footer-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'inline-comments <page-id>',
                'Get inline comments for page',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('resolution-status', { type: 'array', describe: 'Filter by resolution status' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['resolution-status'] && argv['resolution-status'].length > 0) {
                    q['resolution-status'] = argv['resolution-status'];
                  }
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/inline-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'custom-content <page-id>',
                'Get custom content by type in page',
                (y) => {
                  y = addPaginationOptions(y.positional('page-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('type', { type: 'string', describe: 'Custom content type' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.type) q.type = argv.type;
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/custom-content`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'operations <page-id>',
                'Get permitted operations for page',
                (y) => y.positional('page-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/pages/${argv['page-id']}/operations`, null);
                  printJSON(data);
                }
              )
              .command(
                'likes-count <page-id>',
                'Get like count for page',
                (y) => y.positional('page-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/pages/${argv['page-id']}/likes/count`, null);
                  printJSON(data);
                }
              )
              .command(
                'likes-users <page-id>',
                'Get account IDs of likes for page',
                (y) => addPaginationOptions(y.positional('page-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/pages/${argv['page-id']}/likes/users`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'redact <page-id>',
                'Redact content in a Confluence page',
                (y) =>
                  y
                    .positional('page-id', { type: 'string' })
                    .option('body', { type: 'string', demandOption: true, describe: 'JSON redaction request body' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPost(client, `/pages/${argv['page-id']}/redact`, null, body);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a page subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // blogpost
        // ----------------------------------------------------------------
        .command(
          ['blogpost', 'blog', 'bp'],
          'Manage blog posts',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List blog posts',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y
                    .option('id', { type: 'array', describe: 'Filter by blog post IDs' })
                    .option('space-id', { type: 'array', describe: 'Filter by space IDs' })
                    .option('title', { type: 'string', describe: 'Filter by title' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.id && argv.id.length > 0) q.id = argv.id;
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  if (argv.title) q.title = argv.title;
                  const data = await confGetPaginated(client, '/blogposts', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <blogpost-id>',
                'Get blog post by ID',
                (y) => {
                  y = addBodyFormatOption(y.positional('blogpost-id', { type: 'string' }));
                  y = addStatusOption(y);
                  return y
                    .option('get-draft', { type: 'boolean', default: false, describe: 'Retrieve draft version' })
                    .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                    .option('include-labels', { type: 'boolean', default: false, describe: 'Include labels' })
                    .option('include-properties', { type: 'boolean', default: false, describe: 'Include properties' })
                    .option('include-operations', { type: 'boolean', default: false, describe: 'Include operations' })
                    .option('include-likes', { type: 'boolean', default: false, describe: 'Include likes' })
                    .option('include-versions', { type: 'boolean', default: false, describe: 'Include versions' })
                    .option('include-version', { type: 'boolean', default: false, describe: 'Include current version' })
                    .option('include-favorited-by-current-user-status', { type: 'boolean', default: false, describe: 'Include favorited status' })
                    .option('include-webresources', { type: 'boolean', default: false, describe: 'Include web resources' })
                    .option('include-collaborators', { type: 'boolean', default: false, describe: 'Include collaborators' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv['body-format']) q['body-format'] = argv['body-format'];
                  if (argv['get-draft']) q['get-draft'] = 'true';
                  if (argv.version > 0) q.version = String(argv.version);
                  for (const flag of [
                    'include-labels', 'include-properties', 'include-operations',
                    'include-likes', 'include-versions', 'include-version',
                    'include-favorited-by-current-user-status', 'include-webresources',
                    'include-collaborators',
                  ]) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/blogposts/${argv['blogpost-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'create',
                'Create a blog post',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Blog post title' })
                    .option('status', { type: 'string', describe: 'Blog post status' })
                    .option('body', { type: 'string', describe: 'Blog post body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                    .option('private', { type: 'boolean', default: false, describe: 'Create as private' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.private) q.private = 'true';
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv.status) body.status = argv.status;
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  const data = await confPost(client, '/blogposts', q, body);
                  printJSON(data);
                }
              )
              .command(
                'update <blogpost-id>',
                'Update a blog post',
                (y) =>
                  y
                    .positional('blogpost-id', { type: 'string' })
                    .option('title', { type: 'string', demandOption: true, describe: 'Blog post title' })
                    .option('status', { type: 'string', default: 'current', describe: 'Blog post status' })
                    .option('body', { type: 'string', describe: 'Blog post body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                    .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                    .option('version-message', { type: 'string', describe: 'Version message' })
                    .option('space-id', { type: 'string', describe: 'Space ID' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = {
                    id: argv['blogpost-id'],
                    status: argv.status,
                    title: argv.title,
                    version: { number: argv['version-number'], message: argv['version-message'] || '' },
                  };
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  if (argv['space-id']) body.spaceId = argv['space-id'];
                  const data = await confPut(client, `/blogposts/${argv['blogpost-id']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <blogpost-id>',
                'Delete a blog post',
                (y) =>
                  y
                    .positional('blogpost-id', { type: 'string' })
                    .option('purge', { type: 'boolean', default: false, describe: 'Purge the blog post' })
                    .option('draft', { type: 'boolean', default: false, describe: 'Delete a draft blog post' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.purge) q.purge = 'true';
                  if (argv.draft) q.draft = 'true';
                  await confDelete(client, `/blogposts/${argv['blogpost-id']}`, q);
                  console.log('Blog post deleted successfully.');
                }
              )
              .command(
                'attachments <blogpost-id>',
                'Get attachments for blog post',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return y
                    .option('media-type', { type: 'string', describe: 'Filter by media type' })
                    .option('filename', { type: 'string', describe: 'Filter by filename' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['media-type']) q.mediaType = argv['media-type'];
                  if (argv.filename) q.filename = argv.filename;
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/attachments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'labels <blogpost-id>',
                'Get labels for blog post',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  return y.option('prefix', { type: 'string', describe: 'Filter by prefix' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.prefix) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'footer-comments <blogpost-id>',
                'Get footer comments for blog post',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/footer-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'inline-comments <blogpost-id>',
                'Get inline comments for blog post',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('resolution-status', { type: 'array', describe: 'Filter by resolution status' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['resolution-status'] && argv['resolution-status'].length > 0) {
                    q['resolution-status'] = argv['resolution-status'];
                  }
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/inline-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'custom-content <blogpost-id>',
                'Get custom content in blog post',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('type', { type: 'string', describe: 'Custom content type' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.type) q.type = argv.type;
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/custom-content`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'operations <blogpost-id>',
                'Get permitted operations for blog post',
                (y) => y.positional('blogpost-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/blogposts/${argv['blogpost-id']}/operations`, null);
                  printJSON(data);
                }
              )
              .command(
                'versions <blogpost-id>',
                'Get blog post versions',
                (y) => {
                  y = addPaginationOptions(y.positional('blogpost-id', { type: 'string' }));
                  y = addSortOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/versions`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'version-details <blogpost-id> <version-number>',
                'Get version details for blog post version',
                (y) =>
                  y
                    .positional('blogpost-id', { type: 'string' })
                    .positional('version-number', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/blogposts/${argv['blogpost-id']}/versions/${argv['version-number']}`, null);
                  printJSON(data);
                }
              )
              .command(
                'likes-count <blogpost-id>',
                'Get like count for blog post',
                (y) => y.positional('blogpost-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/blogposts/${argv['blogpost-id']}/likes/count`, null);
                  printJSON(data);
                }
              )
              .command(
                'likes-users <blogpost-id>',
                'Get account IDs of likes for blog post',
                (y) => addPaginationOptions(y.positional('blogpost-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/blogposts/${argv['blogpost-id']}/likes/users`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'redact <blogpost-id>',
                'Redact content in a Confluence blog post',
                (y) =>
                  y
                    .positional('blogpost-id', { type: 'string' })
                    .option('body', { type: 'string', demandOption: true, describe: 'JSON redaction request body' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPost(client, `/blogposts/${argv['blogpost-id']}/redact`, null, body);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a blogpost subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // comment
        // ----------------------------------------------------------------
        .command(
          ['comment', 'cm'],
          'Manage comments (footer and inline)',
          (yargs) => {
            return yargs
              .command(
                ['footer', 'fc'],
                'Manage footer comments',
                (yargs) => {
                  return yargs
                    .command(
                      ['list', 'ls'],
                      'List footer comments',
                      (y) => {
                        y = addPaginationOptions(y);
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, '/footer-comments', q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'get <comment-id>',
                      'Get footer comment by ID',
                      (y) => {
                        y = addBodyFormatOption(y.positional('comment-id', { type: 'string' }));
                        return y
                          .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                          .option('include-properties', { type: 'boolean', default: false })
                          .option('include-operations', { type: 'boolean', default: false })
                          .option('include-likes', { type: 'boolean', default: false })
                          .option('include-versions', { type: 'boolean', default: false })
                          .option('include-version', { type: 'boolean', default: false });
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = {};
                        if (argv['body-format']) q['body-format'] = argv['body-format'];
                        if (argv.version > 0) q.version = String(argv.version);
                        for (const flag of ['include-properties', 'include-operations', 'include-likes', 'include-versions', 'include-version']) {
                          if (argv[flag]) q[flag] = 'true';
                        }
                        const data = await confGet(client, `/footer-comments/${argv['comment-id']}`, q);
                        printJSON(data);
                      }
                    )
                    .command(
                      'create',
                      'Create a footer comment',
                      (y) =>
                        y
                          .option('page-id', { type: 'string', describe: 'Page ID to comment on' })
                          .option('blogpost-id', { type: 'string', describe: 'Blog post ID to comment on' })
                          .option('attachment-id', { type: 'string', describe: 'Attachment ID to comment on' })
                          .option('custom-content-id', { type: 'string', describe: 'Custom content ID to comment on' })
                          .option('parent-comment-id', { type: 'string', describe: 'Parent comment ID (for replies)' })
                          .option('body', { type: 'string', describe: 'Comment body content' })
                          .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = {};
                        if (argv['page-id']) body.pageId = argv['page-id'];
                        if (argv['blogpost-id']) body.blogPostId = argv['blogpost-id'];
                        if (argv['attachment-id']) body.attachmentId = argv['attachment-id'];
                        if (argv['custom-content-id']) body.customContentId = argv['custom-content-id'];
                        if (argv['parent-comment-id']) body.parentCommentId = argv['parent-comment-id'];
                        if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                        const data = await confPost(client, '/footer-comments', null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      'update <comment-id>',
                      'Update a footer comment',
                      (y) =>
                        y
                          .positional('comment-id', { type: 'string' })
                          .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                          .option('version-message', { type: 'string', describe: 'Version message' })
                          .option('body', { type: 'string', describe: 'Comment body content' })
                          .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = {
                          version: { number: argv['version-number'], message: argv['version-message'] || '' },
                        };
                        if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                        const data = await confPut(client, `/footer-comments/${argv['comment-id']}`, null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      'delete <comment-id>',
                      'Delete a footer comment',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        await confDelete(client, `/footer-comments/${argv['comment-id']}`, null);
                        console.log('Footer comment deleted successfully.');
                      }
                    )
                    .command(
                      'children <comment-id>',
                      'Get children footer comments',
                      (y) => {
                        y = addPaginationOptions(y.positional('comment-id', { type: 'string' }));
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/footer-comments/${argv['comment-id']}/children`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'operations <comment-id>',
                      'Get permitted operations',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/footer-comments/${argv['comment-id']}/operations`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'versions <comment-id>',
                      'Get footer comment versions',
                      (y) => {
                        y = addPaginationOptions(y.positional('comment-id', { type: 'string' }));
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/footer-comments/${argv['comment-id']}/versions`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'version-details <comment-id> <version-number>',
                      'Get version details for footer comment version',
                      (y) =>
                        y
                          .positional('comment-id', { type: 'string' })
                          .positional('version-number', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/footer-comments/${argv['comment-id']}/versions/${argv['version-number']}`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'likes-count <comment-id>',
                      'Get like count',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/footer-comments/${argv['comment-id']}/likes/count`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'likes-users <comment-id>',
                      'Get like users',
                      (y) => addPaginationOptions(y.positional('comment-id', { type: 'string' })),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/footer-comments/${argv['comment-id']}/likes/users`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .demandCommand(1, 'Specify a footer-comment subcommand');
                },
                () => {}
              )
              .command(
                ['inline', 'ic'],
                'Manage inline comments',
                (yargs) => {
                  return yargs
                    .command(
                      ['list', 'ls'],
                      'List inline comments',
                      (y) => {
                        y = addPaginationOptions(y);
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, '/inline-comments', q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'get <comment-id>',
                      'Get inline comment by ID',
                      (y) => {
                        y = addBodyFormatOption(y.positional('comment-id', { type: 'string' }));
                        return y
                          .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                          .option('include-properties', { type: 'boolean', default: false })
                          .option('include-operations', { type: 'boolean', default: false })
                          .option('include-likes', { type: 'boolean', default: false })
                          .option('include-versions', { type: 'boolean', default: false })
                          .option('include-version', { type: 'boolean', default: false });
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = {};
                        if (argv['body-format']) q['body-format'] = argv['body-format'];
                        if (argv.version > 0) q.version = String(argv.version);
                        for (const flag of ['include-properties', 'include-operations', 'include-likes', 'include-versions', 'include-version']) {
                          if (argv[flag]) q[flag] = 'true';
                        }
                        const data = await confGet(client, `/inline-comments/${argv['comment-id']}`, q);
                        printJSON(data);
                      }
                    )
                    .command(
                      'create',
                      'Create an inline comment',
                      (y) =>
                        y
                          .option('page-id', { type: 'string', describe: 'Page ID to comment on' })
                          .option('blogpost-id', { type: 'string', describe: 'Blog post ID to comment on' })
                          .option('parent-comment-id', { type: 'string', describe: 'Parent comment ID (for replies)' })
                          .option('body', { type: 'string', describe: 'Comment body content' })
                          .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                          .option('inline-comment-properties', { type: 'string', describe: 'JSON inline comment properties' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = {};
                        if (argv['page-id']) body.pageId = argv['page-id'];
                        if (argv['blogpost-id']) body.blogPostId = argv['blogpost-id'];
                        if (argv['parent-comment-id']) body.parentCommentId = argv['parent-comment-id'];
                        if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                        if (argv['inline-comment-properties']) {
                          body.inlineCommentProperties = JSON.parse(argv['inline-comment-properties']);
                        }
                        const data = await confPost(client, '/inline-comments', null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      'update <comment-id>',
                      'Update an inline comment',
                      (y) =>
                        y
                          .positional('comment-id', { type: 'string' })
                          .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                          .option('version-message', { type: 'string', describe: 'Version message' })
                          .option('body', { type: 'string', describe: 'Comment body content' })
                          .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                          .option('resolved', { type: 'boolean', describe: 'Resolved state' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = {
                          version: { number: argv['version-number'], message: argv['version-message'] || '' },
                        };
                        if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                        if (argv.resolved !== undefined) body.resolved = argv.resolved;
                        const data = await confPut(client, `/inline-comments/${argv['comment-id']}`, null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      'delete <comment-id>',
                      'Delete an inline comment',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        await confDelete(client, `/inline-comments/${argv['comment-id']}`, null);
                        console.log('Inline comment deleted successfully.');
                      }
                    )
                    .command(
                      'children <comment-id>',
                      'Get children inline comments',
                      (y) => {
                        y = addPaginationOptions(y.positional('comment-id', { type: 'string' }));
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/inline-comments/${argv['comment-id']}/children`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'operations <comment-id>',
                      'Get permitted operations',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/inline-comments/${argv['comment-id']}/operations`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'versions <comment-id>',
                      'Get inline comment versions',
                      (y) => {
                        y = addPaginationOptions(y.positional('comment-id', { type: 'string' }));
                        y = addSortOption(y);
                        return addBodyFormatOption(y);
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/inline-comments/${argv['comment-id']}/versions`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      'version-details <comment-id> <version-number>',
                      'Get version details for inline comment version',
                      (y) =>
                        y
                          .positional('comment-id', { type: 'string' })
                          .positional('version-number', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/inline-comments/${argv['comment-id']}/versions/${argv['version-number']}`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'likes-count <comment-id>',
                      'Get like count',
                      (y) => y.positional('comment-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `/inline-comments/${argv['comment-id']}/likes/count`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'likes-users <comment-id>',
                      'Get like users',
                      (y) => addPaginationOptions(y.positional('comment-id', { type: 'string' })),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        const data = await confGetPaginated(client, `/inline-comments/${argv['comment-id']}/likes/users`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .demandCommand(1, 'Specify an inline-comment subcommand');
                },
                () => {}
              )
              .demandCommand(1, 'Specify a comment subcommand (footer or inline)');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // label
        // ----------------------------------------------------------------
        .command(
          ['label', 'l'],
          'Manage labels',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List labels',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  return y
                    .option('label-id', { type: 'array', describe: 'Filter by label IDs' })
                    .option('prefix', { type: 'array', describe: 'Filter by prefix' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['label-id'] && argv['label-id'].length > 0) q['label-id'] = argv['label-id'];
                  if (argv.prefix && argv.prefix.length > 0) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, '/labels', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'pages <label-id>',
                'Get pages for label',
                (y) => {
                  y = addPaginationOptions(y.positional('label-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('space-id', { type: 'array', describe: 'Filter by space IDs' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  const data = await confGetPaginated(client, `/labels/${argv['label-id']}/pages`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'blogposts <label-id>',
                'Get blog posts for label',
                (y) => {
                  y = addPaginationOptions(y.positional('label-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('space-id', { type: 'array', describe: 'Filter by space IDs' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  const data = await confGetPaginated(client, `/labels/${argv['label-id']}/blogposts`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'attachments <label-id>',
                'Get attachments for label',
                (y) => addSortOption(addPaginationOptions(y.positional('label-id', { type: 'string' }))),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/labels/${argv['label-id']}/attachments`, q, argv.all);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a label subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // attachment
        // ----------------------------------------------------------------
        .command(
          ['attachment', 'att', 'a'],
          'Manage attachments',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List attachments',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return y
                    .option('media-type', { type: 'string', describe: 'Filter by media type' })
                    .option('filename', { type: 'string', describe: 'Filter by filename' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['media-type']) q.mediaType = argv['media-type'];
                  if (argv.filename) q.filename = argv.filename;
                  const data = await confGetPaginated(client, '/attachments', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <attachment-id>',
                'Get attachment by ID',
                (y) =>
                  y
                    .positional('attachment-id', { type: 'string' })
                    .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                    .option('include-labels', { type: 'boolean', default: false })
                    .option('include-properties', { type: 'boolean', default: false })
                    .option('include-operations', { type: 'boolean', default: false })
                    .option('include-versions', { type: 'boolean', default: false })
                    .option('include-version', { type: 'boolean', default: false })
                    .option('include-collaborators', { type: 'boolean', default: false }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.version > 0) q.version = String(argv.version);
                  for (const flag of ['include-labels', 'include-properties', 'include-operations', 'include-versions', 'include-version', 'include-collaborators']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/attachments/${argv['attachment-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'delete <attachment-id>',
                'Delete attachment',
                (y) =>
                  y
                    .positional('attachment-id', { type: 'string' })
                    .option('purge', { type: 'boolean', default: false, describe: 'Purge the attachment' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.purge) q.purge = 'true';
                  await confDelete(client, `/attachments/${argv['attachment-id']}`, q);
                  console.log('Attachment deleted successfully.');
                }
              )
              .command(
                'labels <attachment-id>',
                'Get labels for attachment',
                (y) => {
                  y = addPaginationOptions(y.positional('attachment-id', { type: 'string' }));
                  y = addSortOption(y);
                  return y.option('prefix', { type: 'string', describe: 'Filter by prefix' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.prefix) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, `/attachments/${argv['attachment-id']}/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'comments <attachment-id>',
                'Get footer comments for attachment',
                (y) => {
                  y = addPaginationOptions(y.positional('attachment-id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y.option('version', { type: 'number', default: 0, describe: 'Filter by version' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.version > 0) q.version = String(argv.version);
                  const data = await confGetPaginated(client, `/attachments/${argv['attachment-id']}/footer-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'operations <attachment-id>',
                'Get permitted operations',
                (y) => y.positional('attachment-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/attachments/${argv['attachment-id']}/operations`, null);
                  printJSON(data);
                }
              )
              .command(
                'versions <attachment-id>',
                'Get attachment versions',
                (y) => addSortOption(addPaginationOptions(y.positional('attachment-id', { type: 'string' }))),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/attachments/${argv['attachment-id']}/versions`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'version-details <attachment-id> <version-number>',
                'Get version details for attachment version',
                (y) =>
                  y
                    .positional('attachment-id', { type: 'string' })
                    .positional('version-number', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/attachments/${argv['attachment-id']}/versions/${argv['version-number']}`, null);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify an attachment subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // task
        // ----------------------------------------------------------------
        .command(
          ['task', 't'],
          'Manage tasks',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List tasks',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addBodyFormatOption(y);
                  return y
                    .option('task-status', { type: 'string', describe: 'Filter by task status (complete, incomplete)' })
                    .option('task-id', { type: 'array', describe: 'Filter by task IDs' })
                    .option('space-id', { type: 'array', describe: 'Filter by space IDs' })
                    .option('page-id', { type: 'array', describe: 'Filter by page IDs' })
                    .option('blogpost-id', { type: 'array', describe: 'Filter by blog post IDs' })
                    .option('created-by', { type: 'array', describe: 'Filter by creator account IDs' })
                    .option('assigned-to', { type: 'array', describe: 'Filter by assignee account IDs' })
                    .option('completed-by', { type: 'array', describe: 'Filter by completer account IDs' })
                    .option('include-blank-tasks', { type: 'boolean', default: false, describe: 'Include blank tasks' })
                    .option('created-at-from', { type: 'string', describe: 'Filter by creation date start (epoch ms)' })
                    .option('created-at-to', { type: 'string', describe: 'Filter by creation date end (epoch ms)' })
                    .option('due-at-from', { type: 'string', describe: 'Filter by due date start (epoch ms)' })
                    .option('due-at-to', { type: 'string', describe: 'Filter by due date end (epoch ms)' })
                    .option('completed-at-from', { type: 'string', describe: 'Filter by completion date start (epoch ms)' })
                    .option('completed-at-to', { type: 'string', describe: 'Filter by completion date end (epoch ms)' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['task-status']) q.status = argv['task-status'];
                  if (argv['task-id'] && argv['task-id'].length > 0) q['task-id'] = argv['task-id'];
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  if (argv['page-id'] && argv['page-id'].length > 0) q['page-id'] = argv['page-id'];
                  if (argv['blogpost-id'] && argv['blogpost-id'].length > 0) q['blogpost-id'] = argv['blogpost-id'];
                  if (argv['created-by'] && argv['created-by'].length > 0) q['created-by'] = argv['created-by'];
                  if (argv['assigned-to'] && argv['assigned-to'].length > 0) q['assigned-to'] = argv['assigned-to'];
                  if (argv['completed-by'] && argv['completed-by'].length > 0) q['completed-by'] = argv['completed-by'];
                  if (argv['include-blank-tasks']) q['include-blank-tasks'] = 'true';
                  for (const f of ['created-at-from', 'created-at-to', 'due-at-from', 'due-at-to', 'completed-at-from', 'completed-at-to']) {
                    if (argv[f]) q[f] = argv[f];
                  }
                  const data = await confGetPaginated(client, '/tasks', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <task-id>',
                'Get task by ID',
                (y) => addBodyFormatOption(y.positional('task-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGet(client, `/tasks/${argv['task-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'update <task-id>',
                'Update a task',
                (y) => {
                  y = addBodyFormatOption(y.positional('task-id', { type: 'string' }));
                  return y
                    .option('task-status', { type: 'string', demandOption: true, describe: 'Task status (complete, incomplete)' })
                    .option('task-id-field', { type: 'string', describe: 'Task ID field' })
                    .option('assigned-to', { type: 'string', describe: 'Assignee account ID' })
                    .option('due-at', { type: 'string', describe: 'Due date' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const body = { status: argv['task-status'] };
                  if (argv['task-id-field']) body.id = argv['task-id-field'];
                  if (argv['assigned-to']) body.assignedTo = argv['assigned-to'];
                  if (argv['due-at']) body.dueAt = argv['due-at'];
                  const data = await confPut(client, `/tasks/${argv['task-id']}`, q, body);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a task subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // custom-content
        // ----------------------------------------------------------------
        .command(
          ['custom-content', 'cc'],
          'Manage custom content',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'List custom content by type',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  y = addBodyFormatOption(y);
                  return y
                    .option('type', { type: 'string', describe: 'Custom content type' })
                    .option('id', { type: 'array', describe: 'Filter by IDs' })
                    .option('space-id', { type: 'array', describe: 'Filter by space IDs' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.type) q.type = argv.type;
                  if (argv.id && argv.id.length > 0) q.id = argv.id;
                  if (argv['space-id'] && argv['space-id'].length > 0) q['space-id'] = argv['space-id'];
                  const data = await confGetPaginated(client, '/custom-content', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <custom-content-id>',
                'Get custom content by ID',
                (y) => {
                  y = addBodyFormatOption(y.positional('custom-content-id', { type: 'string' }));
                  return y
                    .option('version', { type: 'number', default: 0, describe: 'Retrieve a specific version' })
                    .option('include-labels', { type: 'boolean', default: false })
                    .option('include-properties', { type: 'boolean', default: false })
                    .option('include-operations', { type: 'boolean', default: false })
                    .option('include-versions', { type: 'boolean', default: false })
                    .option('include-version', { type: 'boolean', default: false })
                    .option('include-collaborators', { type: 'boolean', default: false });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv['body-format']) q['body-format'] = argv['body-format'];
                  if (argv.version > 0) q.version = String(argv.version);
                  for (const flag of ['include-labels', 'include-properties', 'include-operations', 'include-versions', 'include-version', 'include-collaborators']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/custom-content/${argv['custom-content-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'create',
                'Create custom content',
                (y) =>
                  y
                    .option('type', { type: 'string', demandOption: true, describe: 'Custom content type' })
                    .option('title', { type: 'string', demandOption: true, describe: 'Title' })
                    .option('status', { type: 'string', describe: 'Status' })
                    .option('space-id', { type: 'string', describe: 'Space ID' })
                    .option('page-id', { type: 'string', describe: 'Page ID' })
                    .option('blogpost-id', { type: 'string', describe: 'Blog post ID' })
                    .option('custom-content-id', { type: 'string', describe: 'Parent custom content ID' })
                    .option('body', { type: 'string', describe: 'Body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = { type: argv.type, title: argv.title };
                  if (argv.status) body.status = argv.status;
                  if (argv['space-id']) body.spaceId = argv['space-id'];
                  if (argv['page-id']) body.pageId = argv['page-id'];
                  if (argv['blogpost-id']) body.blogPostId = argv['blogpost-id'];
                  if (argv['custom-content-id']) body.customContentId = argv['custom-content-id'];
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  const data = await confPost(client, '/custom-content', null, body);
                  printJSON(data);
                }
              )
              .command(
                'update <custom-content-id>',
                'Update custom content',
                (y) =>
                  y
                    .positional('custom-content-id', { type: 'string' })
                    .option('type', { type: 'string', demandOption: true, describe: 'Custom content type' })
                    .option('title', { type: 'string', demandOption: true, describe: 'Title' })
                    .option('status', { type: 'string', default: 'current', describe: 'Status' })
                    .option('body', { type: 'string', describe: 'Body content' })
                    .option('body-format', { type: 'string', default: 'storage', describe: 'Body format' })
                    .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                    .option('version-message', { type: 'string', describe: 'Version message' })
                    .option('space-id', { type: 'string', describe: 'Space ID' })
                    .option('page-id', { type: 'string', describe: 'Page ID' })
                    .option('blogpost-id', { type: 'string', describe: 'Blog post ID' })
                    .option('custom-content-id-parent', { type: 'string', describe: 'Parent custom content ID' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = {
                    id: argv['custom-content-id'],
                    type: argv.type,
                    status: argv.status,
                    title: argv.title,
                    version: { number: argv['version-number'], message: argv['version-message'] || '' },
                  };
                  if (argv.body) body.body = { representation: argv['body-format'], value: argv.body };
                  if (argv['space-id']) body.spaceId = argv['space-id'];
                  if (argv['page-id']) body.pageId = argv['page-id'];
                  if (argv['blogpost-id']) body.blogPostId = argv['blogpost-id'];
                  if (argv['custom-content-id-parent']) body.customContentId = argv['custom-content-id-parent'];
                  const data = await confPut(client, `/custom-content/${argv['custom-content-id']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <custom-content-id>',
                'Delete custom content',
                (y) =>
                  y
                    .positional('custom-content-id', { type: 'string' })
                    .option('purge', { type: 'boolean', default: false, describe: 'Purge the custom content' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.purge) q.purge = 'true';
                  await confDelete(client, `/custom-content/${argv['custom-content-id']}`, q);
                  console.log('Custom content deleted successfully.');
                }
              )
              .command(
                'attachments <id>',
                'Get attachments for custom content',
                (y) => {
                  y = addPaginationOptions(y.positional('id', { type: 'string' }));
                  y = addSortOption(y);
                  y = addStatusOption(y);
                  return y
                    .option('media-type', { type: 'string', describe: 'Filter by media type' })
                    .option('filename', { type: 'string', describe: 'Filter by filename' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv['media-type']) q.mediaType = argv['media-type'];
                  if (argv.filename) q.filename = argv.filename;
                  const data = await confGetPaginated(client, `/custom-content/${argv.id}/attachments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'children <id>',
                'Get child custom content',
                (y) => addSortOption(addPaginationOptions(y.positional('id', { type: 'string' }))),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/custom-content/${argv.id}/children`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'labels <id>',
                'Get labels for custom content',
                (y) => {
                  y = addPaginationOptions(y.positional('id', { type: 'string' }));
                  y = addSortOption(y);
                  return y.option('prefix', { type: 'string', describe: 'Filter by prefix' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  if (argv.prefix) q.prefix = argv.prefix;
                  const data = await confGetPaginated(client, `/custom-content/${argv.id}/labels`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'comments <id>',
                'Get footer comments for custom content',
                (y) => {
                  y = addPaginationOptions(y.positional('id', { type: 'string' }));
                  y = addSortOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/custom-content/${argv.id}/footer-comments`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'operations <id>',
                'Get permitted operations',
                (y) => y.positional('id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/custom-content/${argv.id}/operations`, null);
                  printJSON(data);
                }
              )
              .command(
                'versions <id>',
                'Get custom content versions',
                (y) => {
                  y = addPaginationOptions(y.positional('id', { type: 'string' }));
                  y = addSortOption(y);
                  return addBodyFormatOption(y);
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, `/custom-content/${argv.id}/versions`, q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'version-details <id> <version-number>',
                'Get version details for custom content version',
                (y) =>
                  y
                    .positional('id', { type: 'string' })
                    .positional('version-number', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/custom-content/${argv.id}/versions/${argv['version-number']}`, null);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a custom-content subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // whiteboard
        // ----------------------------------------------------------------
        .command(
          ['whiteboard', 'wb'],
          'Manage whiteboards',
          (yargs) => {
            yargs = yargs
              .command(
                'create',
                'Create a whiteboard',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Whiteboard title' })
                    .option('parent-id', { type: 'string', describe: 'Parent ID' })
                    .option('template-key', { type: 'string', describe: 'Template key' })
                    .option('locale', { type: 'string', describe: 'Locale' })
                    .option('private', { type: 'boolean', default: false, describe: 'Create as private' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.private) q.private = 'true';
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  if (argv['template-key']) body.templateKey = argv['template-key'];
                  if (argv.locale) body.locale = argv.locale;
                  const data = await confPost(client, '/whiteboards', q, body);
                  printJSON(data);
                }
              )
              .command(
                'get <whiteboard-id>',
                'Get whiteboard by ID',
                (y) => addIncludeTreeOptions(y.positional('whiteboard-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  for (const flag of ['include-collaborators', 'include-direct-children', 'include-operations', 'include-properties']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/whiteboards/${argv['whiteboard-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'delete <whiteboard-id>',
                'Delete a whiteboard',
                (y) => y.positional('whiteboard-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/whiteboards/${argv['whiteboard-id']}`, null);
                  console.log('Whiteboard deleted successfully.');
                }
              );
            return addTreeSubResources(yargs, '/whiteboards', 'whiteboard')
              .demandCommand(1, 'Specify a whiteboard subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // database
        // ----------------------------------------------------------------
        .command(
          ['database', 'db'],
          'Manage databases',
          (yargs) => {
            yargs = yargs
              .command(
                'create',
                'Create a database',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Database title' })
                    .option('parent-id', { type: 'string', describe: 'Parent ID' })
                    .option('private', { type: 'boolean', default: false, describe: 'Create as private' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.private) q.private = 'true';
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  const data = await confPost(client, '/databases', q, body);
                  printJSON(data);
                }
              )
              .command(
                'get <database-id>',
                'Get database by ID',
                (y) => addIncludeTreeOptions(y.positional('database-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  for (const flag of ['include-collaborators', 'include-direct-children', 'include-operations', 'include-properties']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/databases/${argv['database-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'delete <database-id>',
                'Delete a database',
                (y) => y.positional('database-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/databases/${argv['database-id']}`, null);
                  console.log('Database deleted successfully.');
                }
              );
            return addTreeSubResources(yargs, '/databases', 'database')
              .demandCommand(1, 'Specify a database subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // folder
        // ----------------------------------------------------------------
        .command(
          ['folder', 'f'],
          'Manage folders',
          (yargs) => {
            yargs = yargs
              .command(
                'create',
                'Create a folder',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Folder title' })
                    .option('parent-id', { type: 'string', describe: 'Parent ID' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  const data = await confPost(client, '/folders', null, body);
                  printJSON(data);
                }
              )
              .command(
                'get <folder-id>',
                'Get folder by ID',
                (y) => addIncludeTreeOptions(y.positional('folder-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  for (const flag of ['include-collaborators', 'include-direct-children', 'include-operations', 'include-properties']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/folders/${argv['folder-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'delete <folder-id>',
                'Delete a folder',
                (y) => y.positional('folder-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/folders/${argv['folder-id']}`, null);
                  console.log('Folder deleted successfully.');
                }
              );
            return addTreeSubResources(yargs, '/folders', 'folder')
              .demandCommand(1, 'Specify a folder subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // smart-link
        // ----------------------------------------------------------------
        .command(
          ['smart-link', 'sl', 'embed'],
          'Manage smart links (embeds)',
          (yargs) => {
            yargs = yargs
              .command(
                'create',
                'Create a smart link in the content tree',
                (y) =>
                  y
                    .option('space-id', { type: 'string', demandOption: true, describe: 'Space ID' })
                    .option('title', { type: 'string', describe: 'Smart link title' })
                    .option('parent-id', { type: 'string', describe: 'Parent ID' })
                    .option('embed-url', { type: 'string', describe: 'Embed URL' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = { spaceId: argv['space-id'] };
                  if (argv.title) body.title = argv.title;
                  if (argv['parent-id']) body.parentId = argv['parent-id'];
                  if (argv['embed-url']) body.embedUrl = argv['embed-url'];
                  const data = await confPost(client, '/embeds', null, body);
                  printJSON(data);
                }
              )
              .command(
                'get <smart-link-id>',
                'Get smart link by ID',
                (y) => addIncludeTreeOptions(y.positional('smart-link-id', { type: 'string' })),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  for (const flag of ['include-collaborators', 'include-direct-children', 'include-operations', 'include-properties']) {
                    if (argv[flag]) q[flag] = 'true';
                  }
                  const data = await confGet(client, `/embeds/${argv['smart-link-id']}`, q);
                  printJSON(data);
                }
              )
              .command(
                'delete <smart-link-id>',
                'Delete a smart link',
                (y) => y.positional('smart-link-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/embeds/${argv['smart-link-id']}`, null);
                  console.log('Smart link deleted successfully.');
                }
              );
            return addTreeSubResources(yargs, '/embeds', 'smart link')
              .demandCommand(1, 'Specify a smart-link subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // property
        // ----------------------------------------------------------------
        .command(
          ['property', 'prop'],
          'Manage content properties',
          (yargs) => {
            for (const res of [
              { name: 'page', pathPrefix: '/pages', idParam: 'page-id' },
              { name: 'blogpost', pathPrefix: '/blogposts', idParam: 'blogpost-id' },
              { name: 'comment', pathPrefix: '/comments', idParam: 'comment-id' },
              { name: 'attachment', pathPrefix: '/attachments', idParam: 'attachment-id' },
              { name: 'custom-content', pathPrefix: '/custom-content', idParam: 'custom-content-id' },
              { name: 'whiteboard', pathPrefix: '/whiteboards', idParam: 'whiteboard-id' },
              { name: 'database', pathPrefix: '/databases', idParam: 'database-id' },
              { name: 'folder', pathPrefix: '/folders', idParam: 'folder-id' },
              { name: 'smart-link', pathPrefix: '/embeds', idParam: 'embed-id' },
            ]) {
              const { name, pathPrefix, idParam } = res;
              yargs = yargs.command(
                name,
                `Manage content properties for ${name}`,
                (yargs) => {
                  return yargs
                    .command(
                      [`list <${idParam}>`, `ls <${idParam}>`],
                      `List content properties for ${name}`,
                      (y) => {
                        y = addPaginationOptions(y.positional(idParam, { type: 'string' }));
                        y = addSortOption(y);
                        return y.option('key', { type: 'string', describe: 'Filter by property key' });
                      },
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const q = getPaginationQuery(argv);
                        if (argv.key) q.key = argv.key;
                        const data = await confGetPaginated(client, `${pathPrefix}/${argv[idParam]}/properties`, q, argv.all);
                        printJSON(data);
                      }
                    )
                    .command(
                      `get <${idParam}> <property-id>`,
                      `Get content property for ${name} by ID`,
                      (y) =>
                        y
                          .positional(idParam, { type: 'string' })
                          .positional('property-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `${pathPrefix}/${argv[idParam]}/properties/${argv['property-id']}`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      `create <${idParam}>`,
                      `Create content property for ${name}`,
                      (y) =>
                        y
                          .positional(idParam, { type: 'string' })
                          .option('key', { type: 'string', demandOption: true, describe: 'Property key' })
                          .option('value', { type: 'string', describe: 'Property value (JSON)' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = { key: argv.key };
                        if (argv.value) body.value = parseJSONOrString(argv.value);
                        const data = await confPost(client, `${pathPrefix}/${argv[idParam]}/properties`, null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      `update <${idParam}> <property-id>`,
                      `Update content property for ${name}`,
                      (y) =>
                        y
                          .positional(idParam, { type: 'string' })
                          .positional('property-id', { type: 'string' })
                          .option('key', { type: 'string', demandOption: true, describe: 'Property key' })
                          .option('value', { type: 'string', describe: 'Property value (JSON)' })
                          .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                          .option('version-message', { type: 'string', describe: 'Version message' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = {
                          key: argv.key,
                          version: { number: argv['version-number'], message: argv['version-message'] || '' },
                        };
                        if (argv.value) body.value = parseJSONOrString(argv.value);
                        const data = await confPut(client, `${pathPrefix}/${argv[idParam]}/properties/${argv['property-id']}`, null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      `delete <${idParam}> <property-id>`,
                      `Delete content property for ${name}`,
                      (y) =>
                        y
                          .positional(idParam, { type: 'string' })
                          .positional('property-id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        await confDelete(client, `${pathPrefix}/${argv[idParam]}/properties/${argv['property-id']}`, null);
                        console.log('Property deleted successfully.');
                      }
                    )
                    .demandCommand(1, `Specify a property subcommand for ${name}`);
                },
                () => {}
              );
            }

            yargs = yargs.command(
              'space',
              'Manage space properties',
              (yargs) => {
                return yargs
                  .command(
                    ['list <space-id>', 'ls <space-id>'],
                    'List space properties',
                    (y) => {
                      y = addPaginationOptions(y.positional('space-id', { type: 'string' }));
                      y = addSortOption(y);
                      return y.option('key', { type: 'string', describe: 'Filter by property key' });
                    },
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const q = getPaginationQuery(argv);
                      if (argv.key) q.key = argv.key;
                      const data = await confGetPaginated(client, `/spaces/${argv['space-id']}/properties`, q, argv.all);
                      printJSON(data);
                    }
                  )
                  .command(
                    'get <space-id> <property-id>',
                    'Get space property by ID',
                    (y) =>
                      y
                        .positional('space-id', { type: 'string' })
                        .positional('property-id', { type: 'string' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const data = await confGet(client, `/spaces/${argv['space-id']}/properties/${argv['property-id']}`, null);
                      printJSON(data);
                    }
                  )
                  .command(
                    'create <space-id>',
                    'Create space property',
                    (y) =>
                      y
                        .positional('space-id', { type: 'string' })
                        .option('key', { type: 'string', demandOption: true, describe: 'Property key' })
                        .option('value', { type: 'string', describe: 'Property value (JSON)' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const body = { key: argv.key };
                      if (argv.value) body.value = parseJSONOrString(argv.value);
                      const data = await confPost(client, `/spaces/${argv['space-id']}/properties`, null, body);
                      printJSON(data);
                    }
                  )
                  .command(
                    'update <space-id> <property-id>',
                    'Update space property',
                    (y) =>
                      y
                        .positional('space-id', { type: 'string' })
                        .positional('property-id', { type: 'string' })
                        .option('key', { type: 'string', demandOption: true, describe: 'Property key' })
                        .option('value', { type: 'string', describe: 'Property value (JSON)' })
                        .option('version-number', { type: 'number', demandOption: true, describe: 'Version number' })
                        .option('version-message', { type: 'string', describe: 'Version message' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const body = {
                        key: argv.key,
                        version: { number: argv['version-number'], message: argv['version-message'] || '' },
                      };
                      if (argv.value) body.value = parseJSONOrString(argv.value);
                      const data = await confPut(client, `/spaces/${argv['space-id']}/properties/${argv['property-id']}`, null, body);
                      printJSON(data);
                    }
                  )
                  .command(
                    'delete <space-id> <property-id>',
                    'Delete space property',
                    (y) =>
                      y
                        .positional('space-id', { type: 'string' })
                        .positional('property-id', { type: 'string' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      await confDelete(client, `/spaces/${argv['space-id']}/properties/${argv['property-id']}`, null);
                      console.log('Space property deleted successfully.');
                    }
                  )
                  .demandCommand(1, 'Specify a space property subcommand');
              },
              () => {}
            );

            return yargs.demandCommand(1, 'Specify a property resource type');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // space-permission
        // ----------------------------------------------------------------
        .command(
          ['space-permission', 'sp'],
          'Manage space permissions and roles',
          (yargs) => {
            return yargs
              .command(
                'available',
                'Get available space permissions',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, '/space-permissions', null);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a space-permission subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // admin-key
        // ----------------------------------------------------------------
        .command(
          ['admin-key', 'ak'],
          'Manage admin key',
          (yargs) => {
            return yargs
              .command(
                'get',
                'Get admin key status',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, '/admin-key', null);
                  printJSON(data);
                }
              )
              .command(
                'enable',
                'Enable admin key',
                (y) => y.option('duration', { type: 'number', default: 0, describe: 'Duration in minutes' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = {};
                  if (argv.duration > 0) body.durationInMinutes = argv.duration;
                  const data = await confPost(client, '/admin-key', null, body);
                  printJSON(data);
                }
              )
              .command(
                'disable',
                'Disable admin key',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, '/admin-key', null);
                  console.log('Admin key disabled successfully.');
                }
              )
              .demandCommand(1, 'Specify an admin-key subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // data-policy
        // ----------------------------------------------------------------
        .command(
          ['data-policy', 'dp'],
          'Manage data policies',
          (yargs) => {
            return yargs
              .command(
                'metadata',
                'Get data policy metadata for the workspace',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, '/data-policies/metadata', null);
                  printJSON(data);
                }
              )
              .command(
                'spaces',
                'Get spaces with data policies',
                (y) => {
                  y = addPaginationOptions(y);
                  y = addSortOption(y);
                  return y
                    .option('ids', { type: 'array', describe: 'Filter by space IDs' })
                    .option('keys', { type: 'array', describe: 'Filter by space keys' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv.cursor) q.cursor = argv.cursor;
                  if (argv.limit > 0) q.limit = String(argv.limit);
                  if (argv.sort) q.sort = argv.sort;
                  if (argv.ids && argv.ids.length > 0) q.ids = argv.ids;
                  if (argv.keys && argv.keys.length > 0) q.keys = argv.keys;
                  const data = await confGetPaginated(client, '/data-policies/spaces', q, argv.all);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a data-policy subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // classification
        // ----------------------------------------------------------------
        .command(
          ['classification', 'cl'],
          'Manage classification levels',
          (yargs) => {
            yargs = yargs
              .command(
                ['list', 'ls'],
                'Get list of classification levels',
                (y) => addPaginationOptions(y),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = getPaginationQuery(argv);
                  const data = await confGetPaginated(client, '/classification-levels', q, argv.all);
                  printJSON(data);
                }
              );

            for (const res of [
              { name: 'page', pathPrefix: '/pages' },
              { name: 'blogpost', pathPrefix: '/blogposts' },
              { name: 'database', pathPrefix: '/databases' },
              { name: 'whiteboard', pathPrefix: '/whiteboards' },
            ]) {
              const { name, pathPrefix } = res;
              yargs = yargs.command(
                name,
                `Manage classification level for ${name}`,
                (yargs) => {
                  return yargs
                    .command(
                      'get <id>',
                      `Get ${name} classification level`,
                      (y) => y.positional('id', { type: 'string' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const data = await confGet(client, `${pathPrefix}/${argv.id}/classification-level`, null);
                        printJSON(data);
                      }
                    )
                    .command(
                      'set <id>',
                      `Update ${name} classification level`,
                      (y) =>
                        y
                          .positional('id', { type: 'string' })
                          .option('classification-id', { type: 'string', demandOption: true, describe: 'Classification level ID' })
                          .option('status', { type: 'string', demandOption: true, describe: 'Status' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = { id: argv['classification-id'], status: argv.status };
                        const data = await confPut(client, `${pathPrefix}/${argv.id}/classification-level`, null, body);
                        printJSON(data);
                      }
                    )
                    .command(
                      'reset <id>',
                      `Reset ${name} classification level`,
                      (y) =>
                        y
                          .positional('id', { type: 'string' })
                          .option('status', { type: 'string', demandOption: true, describe: 'Status' }),
                      async (argv) => {
                        const client = getConfluenceClient(argv);
                        const body = { status: argv.status };
                        const data = await confPost(client, `${pathPrefix}/${argv.id}/classification-level/reset`, null, body);
                        printJSON(data);
                      }
                    )
                    .demandCommand(1, `Specify a classification subcommand for ${name}`);
                },
                () => {}
              );
            }

            yargs = yargs.command(
              'space',
              'Manage space default classification level',
              (yargs) => {
                return yargs
                  .command(
                    'get <space-id>',
                    'Get space default classification level',
                    (y) => y.positional('space-id', { type: 'string' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const data = await confGet(client, `/spaces/${argv['space-id']}/classification-level/default`, null);
                      printJSON(data);
                    }
                  )
                  .command(
                    'set <space-id>',
                    'Update space default classification level',
                    (y) =>
                      y
                        .positional('space-id', { type: 'string' })
                        .option('classification-id', { type: 'string', demandOption: true, describe: 'Classification level ID' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      const body = { id: argv['classification-id'] };
                      const data = await confPut(client, `/spaces/${argv['space-id']}/classification-level/default`, null, body);
                      printJSON(data);
                    }
                  )
                  .command(
                    'delete <space-id>',
                    'Delete space default classification level',
                    (y) => y.positional('space-id', { type: 'string' }),
                    async (argv) => {
                      const client = getConfluenceClient(argv);
                      await confDelete(client, `/spaces/${argv['space-id']}/classification-level/default`, null);
                      console.log('Space default classification level deleted successfully.');
                    }
                  )
                  .demandCommand(1, 'Specify a space classification subcommand');
              },
              () => {}
            );

            return yargs.demandCommand(1, 'Specify a classification subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // user
        // ----------------------------------------------------------------
        .command(
          ['user', 'u'],
          'Manage user access',
          (yargs) => {
            return yargs
              .command(
                'bulk-lookup [account-ids..]',
                'Create bulk user lookup using IDs',
                (y) => y.positional('account-ids', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const ids = argv['account-ids'] || [];
                  const data = await confPost(client, '/users-bulk', null, { accountIds: ids });
                  printJSON(data);
                }
              )
              .command(
                'check-access [emails..]',
                'Check site access for a list of emails',
                (y) => y.positional('emails', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const emails = argv.emails || [];
                  const data = await confPost(client, '/user/access/check-access-by-email', null, { emails });
                  printJSON(data);
                }
              )
              .command(
                'invite [emails..]',
                'Invite a list of emails to the site',
                (y) => y.positional('emails', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const emails = argv.emails || [];
                  const data = await confPost(client, '/user/access/invite-by-email', null, { emails });
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a user subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // space-role
        // ----------------------------------------------------------------
        .command(
          ['space-role', 'sr'],
          'Manage space roles',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'Get available space roles',
                (y) => {
                  y = addPaginationOptions(y);
                  return y
                    .option('space-id', { type: 'string', describe: 'Filter by space ID' })
                    .option('role-type', { type: 'string', describe: 'Filter by role type' })
                    .option('principal-id', { type: 'string', describe: 'Filter by principal ID' })
                    .option('principal-type', { type: 'string', describe: 'Filter by principal type' });
                },
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const q = {};
                  if (argv['space-id']) q['space-id'] = argv['space-id'];
                  if (argv['role-type']) q['role-type'] = argv['role-type'];
                  if (argv['principal-id']) q['principal-id'] = argv['principal-id'];
                  if (argv['principal-type']) q['principal-type'] = argv['principal-type'];
                  if (argv.cursor) q.cursor = argv.cursor;
                  if (argv.limit > 0) q.limit = String(argv.limit);
                  const data = await confGetPaginated(client, '/space-roles', q, argv.all);
                  printJSON(data);
                }
              )
              .command(
                'get <role-id>',
                'Get space role by ID',
                (y) => y.positional('role-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/space-roles/${argv['role-id']}`, null);
                  printJSON(data);
                }
              )
              .command(
                'create',
                'Create a space role',
                (y) => y.option('body', { type: 'string', demandOption: true, describe: 'JSON space role definition' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPost(client, '/space-roles', null, body);
                  printJSON(data);
                }
              )
              .command(
                'update <role-id>',
                'Update a space role',
                (y) =>
                  y
                    .positional('role-id', { type: 'string' })
                    .option('body', { type: 'string', demandOption: true, describe: 'JSON space role definition' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPut(client, `/space-roles/${argv['role-id']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <role-id>',
                'Delete a space role',
                (y) => y.positional('role-id', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/space-roles/${argv['role-id']}`, null);
                  console.log('Space role deleted successfully.');
                }
              )
              .command(
                'mode',
                'Get space role mode',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, '/space-role-mode', null);
                  printJSON(data);
                }
              )
              .demandCommand(1, 'Specify a space-role subcommand');
          },
          () => {}
        )

        // ----------------------------------------------------------------
        // convert-ids (misc top-level confluence command)
        // ----------------------------------------------------------------
        .command(
          'convert-ids [content-ids..]',
          'Convert content IDs to content types',
          (y) => y.positional('content-ids', { type: 'string' }),
          async (argv) => {
            const client = getConfluenceClient(argv);
            const ids = argv['content-ids'] || [];
            const data = await confPost(client, '/content/convert-ids-to-types', null, { contentIds: ids });
            printJSON(data);
          }
        )

        // ----------------------------------------------------------------
        // app-property
        // ----------------------------------------------------------------
        .command(
          ['app-property', 'ap'],
          'Manage Forge app properties',
          (yargs) => {
            return yargs
              .command(
                ['list', 'ls'],
                'Get Forge app properties',
                () => {},
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, '/app/properties', null);
                  printJSON(data);
                }
              )
              .command(
                'get <property-key>',
                'Get a Forge app property by key',
                (y) => y.positional('property-key', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const data = await confGet(client, `/app/properties/${argv['property-key']}`, null);
                  printJSON(data);
                }
              )
              .command(
                'set <property-key>',
                'Create or update a Forge app property',
                (y) =>
                  y
                    .positional('property-key', { type: 'string' })
                    .option('body', { type: 'string', demandOption: true, describe: 'JSON property value' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  const body = JSON.parse(argv.body);
                  const data = await confPut(client, `/app/properties/${argv['property-key']}`, null, body);
                  printJSON(data);
                }
              )
              .command(
                'delete <property-key>',
                'Delete a Forge app property',
                (y) => y.positional('property-key', { type: 'string' }),
                async (argv) => {
                  const client = getConfluenceClient(argv);
                  await confDelete(client, `/app/properties/${argv['property-key']}`, null);
                  console.log('App property deleted successfully.');
                }
              )
              .demandCommand(1, 'Specify an app-property subcommand');
          },
          () => {}
        )

        .demandCommand(1, 'Specify a confluence subcommand');
}

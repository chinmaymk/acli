import { get } from './client.js';
import type { BitbucketClient } from './client.js';
import type { CurrentUser } from './types.js';

export async function getCurrentUser(client: BitbucketClient): Promise<CurrentUser> {
  return get(client, '/user');
}

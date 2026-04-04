import { get } from './client.js';
import type { BitbucketClient } from './client.js';

export interface CurrentUser {
  uuid: string;
  nickname: string;
  display_name: string;
  account_id: string;
}

export async function getCurrentUser(client: BitbucketClient): Promise<CurrentUser> {
  return get(client, '/user');
}

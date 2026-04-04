import { get } from './client.js';

/**
 * @param {{ email: string, token: string }} client
 * @returns {Promise<{ uuid: string, nickname: string, display_name: string, account_id: string }>}
 */
export async function getCurrentUser(client) {
  return get(client, '/user');
}

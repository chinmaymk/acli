export interface ConfluenceClient {
  baseURL: string;
  email: string;
  apiToken: string;
}

export function createClient(baseURL: string, email: string, apiToken: string): ConfluenceClient {
  return {
    baseURL: baseURL.replace(/\/+$/, ''),
    email,
    apiToken,
  };
}

export async function confluenceV2(
  client: ConfluenceClient,
  method: string,
  path: string,
  query?: Record<string, string> | null,
  body?: unknown,
): Promise<any> {
  let endpoint = client.baseURL + '/wiki/api/v2' + path;

  if (query && Object.keys(query).length > 0) {
    const params = new URLSearchParams(query);
    endpoint += '?' + params.toString();
  }

  const headers: Record<string, string> = {
    Accept: 'application/json',
  };

  if (client.email) {
    const credentials = Buffer.from(`${client.email}:${client.apiToken}`).toString('base64');
    headers['Authorization'] = 'Basic ' + credentials;
  } else {
    headers['Authorization'] = 'Bearer ' + client.apiToken;
  }

  let bodyStr: string | undefined;
  if (body != null) {
    bodyStr = JSON.stringify(body);
    headers['Content-Type'] = 'application/json';
  }

  const resp = await fetch(endpoint, {
    method,
    headers,
    body: bodyStr,
  });

  const text = await resp.text();

  if (resp.status < 200 || resp.status >= 300) {
    throw new Error(`API error ${resp.status}: ${text}`);
  }

  if (!text) {
    return null;
  }

  return JSON.parse(text);
}

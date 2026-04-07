// Shared low-level type primitives used across the internal API clients
// and CLI handlers. We avoid `any` and `unknown` everywhere; instead, anything
// that is structurally arbitrary JSON uses `JsonValue`, and anything that is
// an Atlassian Document Format (ADF) tree uses `ADFNode`.

export type JsonPrimitive = string | number | boolean | null;
export interface JsonObject {
  [key: string]: JsonValue;
}
export type JsonArray = JsonValue[];
export type JsonValue = JsonPrimitive | JsonObject | JsonArray;

/**
 * A JSON-serializable value used as an HTTP request body. Accepts both
 * structural domain interfaces (without explicit index signatures) and
 * raw JsonValue trees, while still excluding `any`/`unknown` and the
 * non-serializable primitives (`undefined`, `symbol`, `bigint`, functions).
 */
export type JsonBody = JsonValue | object;

// ----------------------------------------------------------------------------
// Atlassian Document Format (ADF)
// ----------------------------------------------------------------------------
//
// ADF is a recursive tree of nodes. Each node has a `type`, optional `content`
// (children), optional `text` (leaf text nodes), optional `attrs` and optional
// `marks`. We model the union loosely but precisely enough for the renderer
// and command builders. Specific node shapes (paragraph, heading, mention,
// emoji, etc.) are all subtypes of `ADFNode`.

export interface ADFMark {
  type: string;
  attrs?: JsonObject;
}

export interface ADFNode {
  type: string;
  text?: string;
  content?: ADFNode[];
  attrs?: JsonObject;
  marks?: ADFMark[];
}

export interface ADFDocument extends ADFNode {
  type: 'doc';
  version: number;
  content: ADFNode[];
}

// ----------------------------------------------------------------------------
// Pagination envelopes
// ----------------------------------------------------------------------------

export interface BitbucketPaginatedResponse<T> {
  values: T[];
  size?: number;
  page?: number;
  pagelen?: number;
  next?: string;
  previous?: string;
}

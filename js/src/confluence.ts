// Public Confluence API surface for library consumers.
//
// The Confluence v2 API is currently consumed via the low-level
// `confluenceV2` helper plus structural typing at the call site, since
// no dedicated `internal/confluence/` client module exists yet (matching
// the Go port). For now we re-export the v2 helper and the shared client
// constructor; richer typed wrappers can be added incrementally.

export { createClient, confluenceV2 } from './internal/api/client.js';
export type { ConfluenceClient } from './internal/api/client.js';

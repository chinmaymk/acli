// Bitbucket API type definitions, mirroring the Go types in internal/bitbucket/.

// --- Common ---

export interface BBUser {
  display_name: string;
  uuid: string;
  nickname?: string;
  account_id?: string;
  type?: string;
}

export interface BBLink {
  href: string;
}

export interface BBContent {
  raw: string;
  markup: string;
  html: string;
}

// --- Workspace ---

export interface Workspace {
  uuid: string;
  name: string;
  slug: string;
  is_private: boolean;
  created_on: string;
  updated_on: string;
  links: {
    html: BBLink;
    avatar: BBLink;
  };
}

export interface WorkspaceMember {
  user: BBUser;
  workspace: {
    slug: string;
    uuid: string;
  };
}

export interface WorkspacePermission {
  permission: string;
  user: BBUser;
}

// --- Repository ---

export interface Repository {
  uuid: string;
  name: string;
  slug: string;
  full_name: string;
  description: string;
  is_private: boolean;
  language: string;
  created_on: string;
  updated_on: string;
  scm: string;
  mainbranch?: { name: string };
  links: {
    html: BBLink;
    clone?: { href: string; name: string }[];
  };
  owner: BBUser;
}

export interface CreateRepoRequest {
  scm: string;
  name: string;
  is_private: boolean;
  description?: string;
  language?: string;
  has_issues?: boolean;
  has_wiki?: boolean;
  fork_policy?: string;
  project?: { key: string };
}

export interface ForkRepoRequest {
  name?: string;
  workspace?: { slug: string };
  is_private?: boolean;
  description?: string;
  language?: string;
}

// --- Branch / Tag ---

export interface Branch {
  name: string;
  target: {
    hash: string;
    date: string;
    message: string;
    author: {
      raw: string;
      user: BBUser;
    };
  };
  links: { html: BBLink };
}

export interface Tag {
  name: string;
  target: {
    hash: string;
    date: string;
    message: string;
    author: {
      raw: string;
      user: BBUser;
    };
  };
  message: string;
  tagger: {
    raw: string;
    user: BBUser;
  };
  links: { html: BBLink };
}

export interface CreateBranchRequest {
  name: string;
  target: { hash: string };
}

export interface CreateTagRequest {
  name: string;
  target: { hash: string };
  message?: string;
}

// --- Pull Requests ---

export interface PullRequest {
  id: number;
  title: string;
  description: string;
  state: string;
  created_on: string;
  updated_on: string;
  author: BBUser;
  source: {
    branch: { name: string };
    repository: { full_name: string };
  };
  destination: {
    branch: { name: string };
    repository: { full_name: string };
  };
  close_source_branch: boolean;
  comment_count: number;
  task_count: number;
  links: { html: BBLink };
}

export interface CreatePRRequest {
  title: string;
  description?: string;
  source_branch: string;
  destination_branch?: string;
  close_source_branch?: boolean;
}

export interface UpdatePRRequest {
  title?: string;
  description?: string;
  close_source_branch?: boolean;
}

export interface MergePRRequest {
  merge_strategy?: string;
  close_source_branch?: boolean;
  message?: string;
}

export interface Participant {
  user: BBUser;
  role: string;
  approved: boolean;
  state: string;
}

export interface PRComment {
  id: number;
  type?: string;
  content: BBContent;
  created_on: string;
  updated_on: string;
  user: BBUser;
  deleted?: boolean;
  inline?: {
    path: string;
    from?: number | null;
    to?: number | null;
    start_from?: number | null;
    start_to?: number | null;
  };
  parent?: { id: number };
  resolution?: CommentResolution;
  pending?: boolean;
}

export interface CommentResolution {
  type: string;
  user: BBUser;
  created_on: string;
}

export interface InlineCommentParams {
  path: string;
  to: number;
}

export interface PRTask {
  id: number;
  state: string;
  content: BBContent;
  creator: BBUser;
  created_on: string;
  updated_on: string;
  resolved_on?: string;
  resolved_by?: BBUser;
  comment?: { id: number };
}

export interface CreatePRTaskRequest {
  content: string;
  comment_id?: number;
}

export interface UpdatePRTaskRequest {
  content?: string;
  state?: string;
}

export interface DiffStat {
  type: string;
  status: string;
  old: CommitFile | null;
  new: CommitFile | null;
  lines_added: number;
  lines_removed: number;
}

export interface CommitFile {
  type: string;
  path: string;
  commit?: Commit;
  attributes?: string;
  escaped_path?: string;
}

// --- Commits ---

export interface Commit {
  hash: string;
  date: string;
  message: string;
  author: {
    raw: string;
    user: BBUser;
  };
  summary?: BBContent;
  parents: { hash: string }[];
  links: { html: BBLink };
  repository: {
    full_name: string;
    uuid: string;
  };
}

export interface CommitStatus {
  uuid: string;
  key: string;
  state: string;
  name: string;
  description: string;
  url: string;
  created_on: string;
  updated_on: string;
  refname: string;
}

// --- Pipelines ---

export interface Pipeline {
  uuid: string;
  build_number: number;
  created_on: string;
  completed_on: string;
  build_seconds_used: number;
  creator: BBUser;
  target: {
    type: string;
    ref_type: string;
    ref_name: string;
    selector: { type: string; pattern: string };
    commit: { hash: string };
  };
  trigger: { type: string; name: string };
  state: {
    name: string;
    result?: { name: string };
    stage?: { name: string };
  };
}

export interface RunPipelineRequest {
  target: {
    type: string;
    ref_type: string;
    ref_name: string;
    selector?: { type: string; pattern: string };
  };
}

export interface PipelineStep {
  uuid: string;
  name: string;
  started_on: string;
  completed_on: string;
  duration_in_seconds: number;
  build_seconds_used: number;
  run_number: number;
  max_time: number;
  state: {
    name: string;
    result?: { name: string };
  };
  setup_commands: { name: string; command: string }[];
  script_commands: { name: string; command: string }[];
  image: { name: string };
}

export interface PipelineVariable {
  uuid: string;
  key: string;
  value: string;
  secured: boolean;
}

// --- Issues ---

export interface BBIssue {
  id: number;
  title: string;
  state: string;
  priority: string;
  kind: string;
  content: BBContent;
  created_on: string;
  updated_on: string;
  reporter: BBUser;
  assignee?: BBUser;
  component?: { name: string };
  milestone?: { name: string };
  version?: { name: string };
  votes: number;
  links: { html: BBLink };
}

export interface CreateIssueRequest {
  title: string;
  content?: { raw: string };
  kind?: string;
  priority?: string;
  state?: string;
  assignee?: { uuid: string };
  component?: { name: string };
  milestone?: { name: string };
  version?: { name: string };
}

export interface UpdateIssueRequest {
  title?: string;
  content?: { raw: string };
  kind?: string;
  priority?: string;
  state?: string;
  assignee?: { uuid: string };
}

export interface IssueComment {
  id: number;
  content: BBContent;
  created_on: string;
  updated_on: string;
  user: BBUser;
}

// --- Projects ---

export interface BBProject {
  uuid: string;
  key: string;
  name: string;
  description: string;
  is_private: boolean;
  created_on: string;
  updated_on: string;
  owner: BBUser;
  links: {
    html: BBLink;
    avatar: BBLink;
  };
}

export interface CreateProjectRequest {
  name: string;
  key: string;
  description?: string;
  is_private?: boolean;
}

// --- Webhooks ---

export interface BBWebhook {
  uuid: string;
  url: string;
  description: string;
  active: boolean;
  created_on: string;
  events: string[];
  subject: { type: string; full_name: string };
  links: { self: BBLink };
}

export interface CreateWebhookRequest {
  description: string;
  url: string;
  active: boolean;
  events: string[];
}

// --- Environments ---

export interface Environment {
  uuid: string;
  name: string;
  slug: string;
  environment_type: { name: string; rank: number };
  rank: number;
  deployment_gate: { name: string };
  lock: { name: string };
  restrictions: { admin_only: boolean; type: string };
}

export interface CreateEnvironmentRequest {
  name: string;
  environment_type: { name: string; rank: number };
}

// --- Deploy Keys ---

export interface DeployKey {
  id: number;
  key: string;
  label: string;
  comment: string;
  created_on: string;
  last_used: string;
  owner: BBUser;
  repository: { full_name: string; uuid: string };
  links: { self: BBLink };
}

export interface CreateDeployKeyRequest {
  key: string;
  label: string;
}

// --- Deployments ---

export interface Deployment {
  uuid: string;
  state: {
    name: string;
    status: { name: string };
  };
  environment: { uuid: string; name: string };
  release: {
    uuid: string;
    name: string;
    url: string;
    commit: { hash: string; message: string };
    created_on: string;
  };
  step: { uuid: string };
}

// --- Downloads ---

export interface Download {
  name: string;
  size: number;
  created_on: string;
  downloads: number;
  links: { self: BBLink };
}

// --- Branch Restrictions ---

export interface BranchRestriction {
  id: number;
  kind: string;
  pattern: string;
  value?: number;
  users: BBUser[];
  groups: { name: string; slug: string }[];
  links: { self: BBLink };
}

export interface CreateBranchRestrictionRequest {
  kind: string;
  pattern: string;
  value?: number;
}

// --- Snippets ---

export interface Snippet {
  id: number;
  title: string;
  scm: string;
  is_private: boolean;
  created_on: string;
  updated_on: string;
  owner: BBUser;
  creator: BBUser;
  links: { html: BBLink };
}

export interface CreateSnippetRequest {
  title: string;
  is_private?: boolean;
  scm?: string;
  files?: Record<string, { content: string }>;
}

// --- Search ---

export interface SearchResult {
  content_match_count: number;
  content_matches: {
    lines: {
      line: number;
      segments: { text: string; match: boolean }[];
    }[];
  }[];
  path_matches: { text: string; match: boolean }[];
  file: { path: string };
}

export interface SearchResponse {
  size: number;
  page: number;
  pagelen: number;
  next: string;
  values: SearchResult[];
}

// --- Pull request activity feed ---

export interface PRActivityUser {
  display_name: string;
  uuid: string;
  nickname: string;
  type: string;
  account_id: string;
}

export interface PRActivityApproval {
  date: string;
  user: PRActivityUser;
}

export interface PRActivityBranchRef {
  branch: { name: string };
  commit: { hash: string; type: string };
  repository: { full_name: string; uuid: string };
}

export interface PRActivityUpdate {
  date: string;
  title: string;
  description: string;
  state: string; // OPEN | MERGED | DECLINED
  reason: string;
  author: PRActivityUser;
  source: PRActivityBranchRef;
  destination: PRActivityBranchRef;
}

export interface PRActivityPullRequestRef {
  type: string;
  id: number;
  title: string;
  links: { html: { href: string } };
}

export interface PRActivity {
  approval?: PRActivityApproval;
  update?: PRActivityUpdate;
  comment?: PRComment;
  pull_request?: PRActivityPullRequestRef;
}

// --- User ---

export interface CurrentUser {
  uuid: string;
  nickname: string;
  display_name: string;
  account_id: string;
}

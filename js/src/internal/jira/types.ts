// Jira API type definitions, mirroring the Go types in internal/jira/types.go.

// --- Pagination ---

export interface Pagination {
  startAt: number;
  maxResults: number;
  total: number;
}

export interface PageBean<T> extends Pagination {
  isLast: boolean;
  values: T[];
}

// --- Core entities ---

export interface AvatarURLs {
  '16x16'?: string;
  '24x24'?: string;
  '32x32'?: string;
  '48x48'?: string;
}

export interface UserDetails {
  self?: string;
  accountId?: string;
  accountType?: string;
  displayName?: string;
  emailAddress?: string;
  active?: boolean;
  avatarUrls?: AvatarURLs;
  timeZone?: string;
  locale?: string;
}

export interface Group {
  name?: string;
  groupId?: string;
  self?: string;
}

export interface StatusCategory {
  self?: string;
  id?: number;
  key?: string;
  colorName?: string;
  name?: string;
}

export interface StatusDetails {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  iconUrl?: string;
  statusCategory?: StatusCategory;
}

export interface Priority {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  iconUrl?: string;
  statusColor?: string;
  isDefault?: boolean;
}

export interface Resolution {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
}

export interface IssueType {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  iconUrl?: string;
  subtask?: boolean;
  avatarId?: number;
  hierarchyLevel?: number;
  statuses?: StatusDetails[];
}

export interface ProjectCategory {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
}

export interface ProjectComponent {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  lead?: UserDetails;
  leadAccountId?: string;
  assigneeType?: string;
  assignee?: UserDetails;
  realAssigneeType?: string;
  realAssignee?: UserDetails;
  isAssigneeTypeValid?: boolean;
  project?: string;
  projectId?: number;
}

export interface Version {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  archived?: boolean;
  released?: boolean;
  overdue?: boolean;
  startDate?: string;
  releaseDate?: string;
  userStartDate?: string;
  userReleaseDate?: string;
  projectId?: number;
}

export interface Project {
  self?: string;
  id?: string;
  key?: string;
  name?: string;
  description?: string;
  lead?: UserDetails;
  components?: ProjectComponent[];
  issueTypes?: IssueType[];
  url?: string;
  email?: string;
  assigneeType?: string;
  versions?: Version[];
  archived?: boolean;
  deleted?: boolean;
  projectTypeKey?: string;
  simplified?: boolean;
  style?: string;
  favourite?: boolean;
  isPrivate?: boolean;
  projectCategory?: ProjectCategory;
  avatarUrls?: AvatarURLs;
}

// --- Issues ---

export interface TimeTracking {
  originalEstimate?: string;
  remainingEstimate?: string;
  timeSpent?: string;
  originalEstimateSeconds?: number;
  remainingEstimateSeconds?: number;
  timeSpentSeconds?: number;
}

export interface IssueLink {
  id?: string;
  self?: string;
  type?: IssueLinkType;
  inwardIssue?: Issue;
  outwardIssue?: Issue;
}

export interface IssueLinkType {
  id?: string;
  name?: string;
  inward?: string;
  outward?: string;
  self?: string;
}

export interface Attachment {
  self?: string;
  id?: string;
  filename?: string;
  author?: UserDetails;
  created?: string;
  size?: number;
  mimeType?: string;
  content?: string;
  thumbnail?: string;
}

export interface Visibility {
  type?: string;
  value?: string;
  identifier?: string;
}

export interface Comment {
  self?: string;
  id?: string;
  author?: UserDetails;
  body?: unknown;
  updateAuthor?: UserDetails;
  created?: string;
  updated?: string;
  visibility?: Visibility;
  jsdPublic?: boolean;
}

export interface CommentPage extends Pagination {
  comments: Comment[];
}

export interface Worklog {
  self?: string;
  id?: string;
  author?: UserDetails;
  updateAuthor?: UserDetails;
  comment?: unknown;
  created?: string;
  updated?: string;
  started?: string;
  timeSpent?: string;
  timeSpentSeconds?: number;
  issueId?: string;
  visibility?: Visibility;
}

export interface WorklogPage extends Pagination {
  worklogs: Worklog[];
}

export interface Votes {
  self?: string;
  votes?: number;
  hasVoted?: boolean;
  voters?: UserDetails[];
}

export interface Watches {
  self?: string;
  watchCount?: number;
  isWatching?: boolean;
  watchers?: UserDetails[];
}

export interface IssueFields {
  summary?: string;
  description?: unknown;
  issuetype?: IssueType;
  project?: Project;
  status?: StatusDetails;
  priority?: Priority;
  assignee?: UserDetails;
  reporter?: UserDetails;
  creator?: UserDetails;
  resolution?: Resolution;
  labels?: string[];
  components?: ProjectComponent[];
  fixVersions?: Version[];
  versions?: Version[];
  created?: string;
  updated?: string;
  duedate?: string;
  parent?: Issue;
  subtasks?: Issue[];
  comment?: CommentPage;
  worklog?: WorklogPage;
  attachment?: Attachment[];
  watches?: Watches;
  votes?: Votes;
  timetracking?: TimeTracking;
  issuelinks?: IssueLink[];
}

export interface Issue {
  id?: string;
  key?: string;
  self?: string;
  fields?: Record<string, unknown>;
}

export interface IssueDetailed {
  id?: string;
  key?: string;
  self?: string;
  fields: IssueFields;
}

export interface CreatedIssue {
  id: string;
  key: string;
  self: string;
}

export interface IssueUpdateDetails {
  fields?: Record<string, unknown>;
  update?: Record<string, FieldUpdateOperation[]>;
  transition?: IssueTransition;
}

export interface FieldUpdateOperation {
  add?: unknown;
  set?: unknown;
  remove?: unknown;
  edit?: unknown;
}

export interface IssueTransition {
  id?: string;
  name?: string;
  to?: StatusDetails;
  hasScreen?: boolean;
  isGlobal?: boolean;
  isInitial?: boolean;
  isAvailable?: boolean;
  isConditional?: boolean;
  isLooped?: boolean;
}

export interface TransitionsResponse {
  transitions: IssueTransition[];
}

export interface SearchResults {
  startAt: number;
  maxResults: number;
  total: number;
  issues: IssueDetailed[];
  warningMessages?: string[];
}

export interface SearchRequest {
  jql?: string;
  startAt?: number;
  maxResults?: number;
  fields?: string[];
  expand?: string[];
  validateQuery?: string;
}

// --- Changelog ---

export interface ChangelogItem {
  field?: string;
  fieldtype?: string;
  fieldId?: string;
  from?: string;
  fromString?: string;
  to?: string;
  toString?: string;
}

export interface Changelog {
  id?: string;
  author?: UserDetails;
  created?: string;
  items?: ChangelogItem[];
}

export interface ChangelogPage extends Pagination {
  histories?: Changelog[];
  values?: Changelog[];
}

// --- Remote Links ---

export interface RemoteLink {
  id?: number;
  self?: string;
  globalId?: string;
  application?: { type?: string; name?: string };
  relationship?: string;
  object?: {
    url?: string;
    title?: string;
    summary?: string;
    icon?: { url16x16?: string; title?: string };
    status?: { resolved?: boolean; icon?: { url16x16?: string; title?: string; link?: string } };
  };
}

// --- Entity Property ---

export interface EntityProperty {
  key?: string;
  value?: unknown;
}

// --- Create Meta ---

export interface CreateMeta {
  projects?: CreateMetaProject[];
}

export interface CreateMetaProject {
  self?: string;
  id?: string;
  key?: string;
  name?: string;
  issuetypes?: CreateMetaIssueType[];
}

export interface CreateMetaIssueType {
  self?: string;
  id?: string;
  description?: string;
  iconUrl?: string;
  name?: string;
  subtask?: boolean;
  fields?: Record<string, unknown>;
}

// --- Notification ---

export interface IssueNotifyRequest {
  htmlBody?: string;
  subject?: string;
  textBody?: string;
  to?: NotifyTo;
  restrict?: NotifyRestrict;
}

export interface NotifyTo {
  reporter?: boolean;
  assignee?: boolean;
  watchers?: boolean;
  voters?: boolean;
  users?: UserDetails[];
  groups?: Group[];
}

export interface NotifyRestrict {
  groups?: Group[];
  permissions?: Permission[];
}

export interface Permission {
  id?: string;
  key?: string;
}

// --- Bulk ---

export interface BulkIssueCreateRequest {
  issueUpdates: IssueUpdateDetails[];
}

export interface BulkIssueCreateResponse {
  issues: CreatedIssue[];
  errors: unknown[];
}

// --- Admin entities ---

export interface ComponentIssueCount {
  self?: string;
  issueCount?: number;
}

export interface VersionIssueCounts {
  self?: string;
  issuesFixedCount?: number;
  issuesAffectedCount?: number;
}

export interface VersionUnresolvedIssueCount {
  self?: string;
  issuesUnresolvedCount?: number;
  issuesCount?: number;
}

export interface GroupMembers extends Pagination {
  isLast?: boolean;
  values: UserDetails[];
}

export interface FoundGroups {
  header?: string;
  total?: number;
  groups?: FoundGroup[];
}

export interface FoundGroup {
  name?: string;
  html?: string;
  labels?: GroupLabel[];
  groupId?: string;
}

export interface GroupLabel {
  text?: string;
  title?: string;
  type?: string;
}

export interface ServerInfo {
  baseUrl?: string;
  version?: string;
  versionNumbers?: number[];
  buildNumber?: number;
  deploymentType?: string;
  displayName?: string;
  scmInfo?: string;
  serverTitle?: string;
}

export interface Configuration {
  votingEnabled?: boolean;
  watchingEnabled?: boolean;
  unassignedIssuesAllowed?: boolean;
  subTasksEnabled?: boolean;
  issueLinkingEnabled?: boolean;
  timeTrackingEnabled?: boolean;
  attachmentsEnabled?: boolean;
}

export interface AnnouncementBanner {
  message?: string;
  isDismissible?: boolean;
  isEnabled?: boolean;
  hashCode?: string;
  visibility?: string;
}

export interface AuditRecords {
  offset?: number;
  limit?: number;
  total?: number;
  records?: AuditRecord[];
}

export interface AuditRecord {
  id?: number;
  summary?: string;
  remoteAddress?: string;
  created?: string;
  category?: string;
  eventSource?: string;
  description?: string;
}

export interface ApplicationRole {
  key?: string;
  name?: string;
  selectedByDefault?: boolean;
  defined?: boolean;
  numberOfSeats?: number;
  remainingSeats?: number;
  userCount?: number;
  hasUnlimitedSeats?: boolean;
  platform?: boolean;
}

export interface TaskResult {
  self?: string;
  id?: string;
  description?: string;
  status?: string;
  result?: string;
  progress?: number;
  elapsedRuntime?: number;
  submitted?: number;
  started?: number;
  finished?: number;
  lastUpdate?: number;
}

export interface AttachmentMeta {
  enabled?: boolean;
  uploadLimit?: number;
}

export interface UserPermission {
  id?: string;
  key?: string;
  name?: string;
  type?: string;
  description?: string;
  havePermission?: boolean;
}

// --- Schemes ---

export interface SharePermission {
  id?: number;
  type?: string;
  project?: Project;
  role?: ProjectRole;
  group?: Group;
}

export interface Dashboard {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  owner?: UserDetails;
  isFavourite?: boolean;
  popularity?: number;
  rank?: number;
  view?: string;
  sharePermissions?: SharePermission[];
  editPermissions?: SharePermission[];
}

export interface DashboardList extends Pagination {
  dashboards: Dashboard[];
}

export interface DashboardGadget {
  id?: number;
  color?: string;
  position?: { row: number; column: number };
  title?: string;
  moduleKey?: string;
  uri?: string;
}

export interface DashboardGadgetList {
  gadgets: DashboardGadget[];
}

export interface Field {
  id?: string;
  key?: string;
  name?: string;
  custom?: boolean;
  orderable?: boolean;
  navigable?: boolean;
  searchable?: boolean;
  clauseNames?: string[];
  schema?: FieldSchema;
}

export interface FieldSchema {
  type?: string;
  items?: string;
  system?: string;
  custom?: string;
  customId?: number;
}

export interface FieldContext {
  id?: string;
  name?: string;
  description?: string;
  isGlobalContext?: boolean;
  isAnyIssueType?: boolean;
}

export interface Workflow {
  id?: string;
  name?: string;
  description?: string;
  isDefault?: boolean;
}

export interface WorkflowScheme {
  id?: number;
  name?: string;
  description?: string;
  defaultWorkflow?: string;
  issueTypeMappings?: Record<string, string>;
  draft?: WorkflowScheme;
  self?: string;
}

export interface Screen {
  id?: number;
  name?: string;
  description?: string;
}

export interface ScreenScheme {
  id?: number;
  name?: string;
  description?: string;
  screens?: ScreenTypes;
}

export interface ScreenTypes {
  create?: number;
  default?: number;
  view?: number;
  edit?: number;
}

export interface ScreenTab {
  id?: number;
  name?: string;
}

export interface ScreenField {
  id?: string;
  name?: string;
}

export interface PermissionScheme {
  id?: number;
  self?: string;
  name?: string;
  description?: string;
  permissions?: PermissionGrant[];
}

export interface PermissionGrant {
  id?: number;
  self?: string;
  holder?: PermissionHolder;
  permission?: string;
}

export interface PermissionHolder {
  type?: string;
  parameter?: string;
}

export interface NotificationScheme {
  id?: number;
  self?: string;
  name?: string;
  description?: string;
}

export interface IssueSecurityScheme {
  self?: string;
  id?: number;
  name?: string;
  description?: string;
  defaultSecurityLevelId?: number;
}

export interface FieldConfiguration {
  id?: number;
  name?: string;
  description?: string;
  isDefault?: boolean;
}

export interface FieldConfigurationScheme {
  id?: string;
  name?: string;
  description?: string;
}

export interface IssueTypeScheme {
  id?: string;
  name?: string;
  description?: string;
  defaultIssueTypeId?: string;
  isDefault?: boolean;
}

export interface IssueTypeScreenScheme {
  id?: string;
  name?: string;
  description?: string;
}

export interface ProjectRole {
  self?: string;
  id?: number;
  name?: string;
  description?: string;
  actors?: RoleActor[];
}

export interface RoleActor {
  id?: number;
  displayName?: string;
  type?: string;
  name?: string;
  actorUser?: UserDetails;
  actorGroup?: Group;
}

export interface Webhook {
  id?: number;
  jqlFilter?: string;
  fieldIdsFilter?: string[];
  issuePropertyKeysFilter?: string[];
  events?: string[];
  expirationDate?: number;
}

// --- Agile ---

export interface Board {
  id: number;
  self?: string;
  name?: string;
  type?: string;
  location?: BoardLocation;
}

export interface BoardLocation {
  projectId?: number;
  displayName?: string;
  projectName?: string;
  projectKey?: string;
  projectTypeKey?: string;
  name?: string;
}

export interface BoardList {
  maxResults: number;
  startAt: number;
  total: number;
  isLast: boolean;
  values: Board[];
}

export interface BoardConfiguration {
  id: number;
  name?: string;
  type?: string;
  self?: string;
  location?: BoardLocation;
  filter?: BoardFilter;
  columnConfig?: BoardColumnConfig;
  ranking?: Record<string, unknown>;
}

export interface BoardFilter {
  id?: string;
  self?: string;
}

export interface BoardColumnConfig {
  columns?: BoardColumn[];
  constraintType?: string;
}

export interface BoardColumn {
  name?: string;
  statuses?: BoardStatus[];
  min?: number;
  max?: number;
}

export interface BoardStatus {
  id?: string;
  self?: string;
}

export interface Sprint {
  id: number;
  self?: string;
  state?: string;
  name?: string;
  startDate?: string;
  endDate?: string;
  completeDate?: string;
  originBoardId?: number;
  goal?: string;
}

export interface SprintList {
  maxResults: number;
  startAt: number;
  total: number;
  isLast: boolean;
  values: Sprint[];
}

export interface SprintIssuesResponse {
  maxResults: number;
  startAt: number;
  total: number;
  issues: IssueDetailed[];
}

export interface Epic {
  id: number;
  key?: string;
  self?: string;
  name?: string;
  summary?: string;
  done?: boolean;
}

export interface EpicList {
  maxResults: number;
  startAt: number;
  total: number;
  isLast: boolean;
  values: Epic[];
}

// --- Filters ---

export interface Filter {
  self?: string;
  id?: string;
  name?: string;
  description?: string;
  owner?: UserDetails;
  jql?: string;
  viewUrl?: string;
  searchUrl?: string;
  favourite?: boolean;
  favouritedCount?: number;
  sharePermissions?: SharePermission[];
  editPermissions?: SharePermission[];
}

// --- Project Features ---

export interface ProjectFeature {
  projectId?: number;
  state?: string;
  toggleLocked?: boolean;
  feature?: string;
  imageUri?: string;
  localisedName?: string;
  localisedDescription?: string;
}

export interface ProjectFeaturesResponse {
  features?: ProjectFeature[];
}

// --- Project Types ---

export interface ProjectType {
  key?: string;
  formattedKey?: string;
  descriptionI18nKey?: string;
  icon?: string;
  color?: string;
}

// --- Project Validation ---

export interface ProjectKeyValidation {
  errorMessages?: string[];
  errors?: Record<string, string>;
}

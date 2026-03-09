package jira

import "encoding/json"

// Pagination is the common pagination fields in Jira responses.
type Pagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

// PageBean is a generic paginated response.
type PageBean[T any] struct {
	Pagination
	IsLast bool `json:"isLast"`
	Values []T  `json:"values"`
}

// AvatarURLs contains avatar URLs at various sizes.
type AvatarURLs struct {
	Size16 string `json:"16x16,omitempty"`
	Size24 string `json:"24x24,omitempty"`
	Size32 string `json:"32x32,omitempty"`
	Size48 string `json:"48x48,omitempty"`
}

// UserDetails represents a Jira user.
type UserDetails struct {
	Self         string     `json:"self,omitempty"`
	AccountID    string     `json:"accountId,omitempty"`
	AccountType  string     `json:"accountType,omitempty"`
	DisplayName  string     `json:"displayName,omitempty"`
	EmailAddress string     `json:"emailAddress,omitempty"`
	Active       bool       `json:"active,omitempty"`
	AvatarURLs   AvatarURLs `json:"avatarUrls,omitempty"`
	TimeZone     string     `json:"timeZone,omitempty"`
	Locale       string     `json:"locale,omitempty"`
}

// Group represents a Jira group.
type Group struct {
	Name    string `json:"name,omitempty"`
	GroupID string `json:"groupId,omitempty"`
	Self    string `json:"self,omitempty"`
}

// StatusCategory represents a status category.
type StatusCategory struct {
	Self      string `json:"self,omitempty"`
	ID        int    `json:"id,omitempty"`
	Key       string `json:"key,omitempty"`
	ColorName string `json:"colorName,omitempty"`
	Name      string `json:"name,omitempty"`
}

// StatusDetails represents an issue status.
type StatusDetails struct {
	Self           string         `json:"self,omitempty"`
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Description    string         `json:"description,omitempty"`
	IconURL        string         `json:"iconUrl,omitempty"`
	StatusCategory StatusCategory `json:"statusCategory,omitempty"`
}

// Priority represents an issue priority.
type Priority struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	StatusColor string `json:"statusColor,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

// Resolution represents an issue resolution.
type Resolution struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// IssueType represents an issue type.
type IssueType struct {
	Self           string `json:"self,omitempty"`
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	IconURL        string `json:"iconUrl,omitempty"`
	Subtask        bool   `json:"subtask,omitempty"`
	AvatarID       int    `json:"avatarId,omitempty"`
	HierarchyLevel int    `json:"hierarchyLevel,omitempty"`
}

// ProjectCategory represents a project category.
type ProjectCategory struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ProjectComponent represents a project component.
type ProjectComponent struct {
	Self                string       `json:"self,omitempty"`
	ID                  string       `json:"id,omitempty"`
	Name                string       `json:"name,omitempty"`
	Description         string       `json:"description,omitempty"`
	Lead                *UserDetails `json:"lead,omitempty"`
	LeadAccountID       string       `json:"leadAccountId,omitempty"`
	AssigneeType        string       `json:"assigneeType,omitempty"`
	Assignee            *UserDetails `json:"assignee,omitempty"`
	RealAssigneeType    string       `json:"realAssigneeType,omitempty"`
	RealAssignee        *UserDetails `json:"realAssignee,omitempty"`
	IsAssigneeTypeValid bool         `json:"isAssigneeTypeValid,omitempty"`
	Project             string       `json:"project,omitempty"`
	ProjectID           int          `json:"projectId,omitempty"`
}

// Version represents a project version.
type Version struct {
	Self            string `json:"self,omitempty"`
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Archived        bool   `json:"archived,omitempty"`
	Released        bool   `json:"released,omitempty"`
	Overdue         bool   `json:"overdue,omitempty"`
	StartDate       string `json:"startDate,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty"`
	UserStartDate   string `json:"userStartDate,omitempty"`
	UserReleaseDate string `json:"userReleaseDate,omitempty"`
	ProjectID       int    `json:"projectId,omitempty"`
}

// Project represents a Jira project.
type Project struct {
	Self            string             `json:"self,omitempty"`
	ID              string             `json:"id,omitempty"`
	Key             string             `json:"key,omitempty"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	Lead            *UserDetails       `json:"lead,omitempty"`
	Components      []ProjectComponent `json:"components,omitempty"`
	IssueTypes      []IssueType        `json:"issueTypes,omitempty"`
	URL             string             `json:"url,omitempty"`
	Email           string             `json:"email,omitempty"`
	AssigneeType    string             `json:"assigneeType,omitempty"`
	Versions        []Version          `json:"versions,omitempty"`
	Archived        bool               `json:"archived,omitempty"`
	Deleted         bool               `json:"deleted,omitempty"`
	ProjectTypeKey  string             `json:"projectTypeKey,omitempty"`
	Simplified      bool               `json:"simplified,omitempty"`
	Style           string             `json:"style,omitempty"`
	Favourite       bool               `json:"favourite,omitempty"`
	IsPrivate       bool               `json:"isPrivate,omitempty"`
	ProjectCategory *ProjectCategory   `json:"projectCategory,omitempty"`
	AvatarURLs      AvatarURLs         `json:"avatarUrls,omitempty"`
}

// IssueFields contains the standard fields of a Jira issue.
type IssueFields struct {
	Summary      string             `json:"summary,omitempty"`
	Description  interface{}        `json:"description,omitempty"`
	IssueType    *IssueType         `json:"issuetype,omitempty"`
	Project      *Project           `json:"project,omitempty"`
	Status       *StatusDetails     `json:"status,omitempty"`
	Priority     *Priority          `json:"priority,omitempty"`
	Assignee     *UserDetails       `json:"assignee,omitempty"`
	Reporter     *UserDetails       `json:"reporter,omitempty"`
	Creator      *UserDetails       `json:"creator,omitempty"`
	Resolution   *Resolution        `json:"resolution,omitempty"`
	Labels       []string           `json:"labels,omitempty"`
	Components   []ProjectComponent `json:"components,omitempty"`
	FixVersions  []Version          `json:"fixVersions,omitempty"`
	Versions     []Version          `json:"versions,omitempty"`
	Created      string             `json:"created,omitempty"`
	Updated      string             `json:"updated,omitempty"`
	DueDate      string             `json:"duedate,omitempty"`
	Parent       *Issue             `json:"parent,omitempty"`
	Subtasks     []Issue            `json:"subtasks,omitempty"`
	Comment      *CommentPage       `json:"comment,omitempty"`
	Worklog      *WorklogPage       `json:"worklog,omitempty"`
	Attachment   []Attachment       `json:"attachment,omitempty"`
	Watches      *Watches           `json:"watches,omitempty"`
	Votes        *Votes             `json:"votes,omitempty"`
	TimeTracking *TimeTracking      `json:"timetracking,omitempty"`
	IssueLinks   []IssueLink        `json:"issuelinks,omitempty"`
}

// Issue represents a Jira issue.
type Issue struct {
	ID     string          `json:"id,omitempty"`
	Key    string          `json:"key,omitempty"`
	Self   string          `json:"self,omitempty"`
	Fields json.RawMessage `json:"fields,omitempty"`
}

// IssueDetailed is an Issue with parsed standard fields.
type IssueDetailed struct {
	ID     string      `json:"id,omitempty"`
	Key    string      `json:"key,omitempty"`
	Self   string      `json:"self,omitempty"`
	Fields IssueFields `json:"fields,omitempty"`
}

// CreatedIssue represents the response after creating an issue.
type CreatedIssue struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// IssueUpdateDetails is the request body for creating/updating issues.
type IssueUpdateDetails struct {
	Fields     map[string]interface{}            `json:"fields,omitempty"`
	Update     map[string][]FieldUpdateOperation `json:"update,omitempty"`
	Transition *IssueTransition                  `json:"transition,omitempty"`
}

// FieldUpdateOperation represents a field update operation.
type FieldUpdateOperation struct {
	Add    interface{} `json:"add,omitempty"`
	Set    interface{} `json:"set,omitempty"`
	Remove interface{} `json:"remove,omitempty"`
	Edit   interface{} `json:"edit,omitempty"`
}

// IssueTransition represents a transition to perform on an issue.
type IssueTransition struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	To            *StatusDetails `json:"to,omitempty"`
	HasScreen     bool           `json:"hasScreen,omitempty"`
	IsGlobal      bool           `json:"isGlobal,omitempty"`
	IsInitial     bool           `json:"isInitial,omitempty"`
	IsAvailable   bool           `json:"isAvailable,omitempty"`
	IsConditional bool           `json:"isConditional,omitempty"`
	IsLooped      bool           `json:"isLooped,omitempty"`
}

// TransitionsResponse is the response from getting issue transitions.
type TransitionsResponse struct {
	Transitions []IssueTransition `json:"transitions"`
}

// SearchResults represents search results.
type SearchResults struct {
	StartAt         int             `json:"startAt"`
	MaxResults      int             `json:"maxResults"`
	Total           int             `json:"total"`
	Issues          []IssueDetailed `json:"issues"`
	WarningMessages []string        `json:"warningMessages,omitempty"`
}

// SearchRequest is the request body for JQL search.
type SearchRequest struct {
	JQL           string   `json:"jql,omitempty"`
	StartAt       int      `json:"startAt,omitempty"`
	MaxResults    int      `json:"maxResults,omitempty"`
	Fields        []string `json:"fields,omitempty"`
	Expand        []string `json:"expand,omitempty"`
	ValidateQuery string   `json:"validateQuery,omitempty"`
}

// Comment represents a Jira issue comment.
type Comment struct {
	Self         string       `json:"self,omitempty"`
	ID           string       `json:"id,omitempty"`
	Author       *UserDetails `json:"author,omitempty"`
	Body         interface{}  `json:"body,omitempty"`
	UpdateAuthor *UserDetails `json:"updateAuthor,omitempty"`
	Created      string       `json:"created,omitempty"`
	Updated      string       `json:"updated,omitempty"`
	Visibility   *Visibility  `json:"visibility,omitempty"`
	JsdPublic    bool         `json:"jsdPublic,omitempty"`
}

// CommentPage is a paginated list of comments.
type CommentPage struct {
	Pagination
	Comments []Comment `json:"comments"`
}

// Visibility represents comment/worklog visibility restrictions.
type Visibility struct {
	Type       string `json:"type,omitempty"`
	Value      string `json:"value,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// Worklog represents a worklog entry.
type Worklog struct {
	Self             string       `json:"self,omitempty"`
	ID               string       `json:"id,omitempty"`
	Author           *UserDetails `json:"author,omitempty"`
	UpdateAuthor     *UserDetails `json:"updateAuthor,omitempty"`
	Comment          interface{}  `json:"comment,omitempty"`
	Created          string       `json:"created,omitempty"`
	Updated          string       `json:"updated,omitempty"`
	Started          string       `json:"started,omitempty"`
	TimeSpent        string       `json:"timeSpent,omitempty"`
	TimeSpentSeconds int64        `json:"timeSpentSeconds,omitempty"`
	IssueID          string       `json:"issueId,omitempty"`
	Visibility       *Visibility  `json:"visibility,omitempty"`
}

// WorklogPage is a paginated list of worklogs.
type WorklogPage struct {
	Pagination
	Worklogs []Worklog `json:"worklogs"`
}

// Attachment represents a file attachment.
type Attachment struct {
	Self      string       `json:"self,omitempty"`
	ID        string       `json:"id,omitempty"`
	Filename  string       `json:"filename,omitempty"`
	Author    *UserDetails `json:"author,omitempty"`
	Created   string       `json:"created,omitempty"`
	Size      int64        `json:"size,omitempty"`
	MimeType  string       `json:"mimeType,omitempty"`
	Content   string       `json:"content,omitempty"`
	Thumbnail string       `json:"thumbnail,omitempty"`
}

// Watches represents issue watchers.
type Watches struct {
	Self       string        `json:"self,omitempty"`
	WatchCount int           `json:"watchCount,omitempty"`
	IsWatching bool          `json:"isWatching,omitempty"`
	Watchers   []UserDetails `json:"watchers,omitempty"`
}

// Votes represents issue votes.
type Votes struct {
	Self     string        `json:"self,omitempty"`
	Votes    int           `json:"votes,omitempty"`
	HasVoted bool          `json:"hasVoted,omitempty"`
	Voters   []UserDetails `json:"voters,omitempty"`
}

// TimeTracking represents time tracking fields.
type TimeTracking struct {
	OriginalEstimate         string `json:"originalEstimate,omitempty"`
	RemainingEstimate        string `json:"remainingEstimate,omitempty"`
	TimeSpent                string `json:"timeSpent,omitempty"`
	OriginalEstimateSeconds  int64  `json:"originalEstimateSeconds,omitempty"`
	RemainingEstimateSeconds int64  `json:"remainingEstimateSeconds,omitempty"`
	TimeSpentSeconds         int64  `json:"timeSpentSeconds,omitempty"`
}

// IssueLink represents a link between two issues.
type IssueLink struct {
	ID           string         `json:"id,omitempty"`
	Self         string         `json:"self,omitempty"`
	Type         *IssueLinkType `json:"type,omitempty"`
	InwardIssue  *Issue         `json:"inwardIssue,omitempty"`
	OutwardIssue *Issue         `json:"outwardIssue,omitempty"`
}

// IssueLinkType represents a type of link between issues.
type IssueLinkType struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Inward  string `json:"inward,omitempty"`
	Outward string `json:"outward,omitempty"`
	Self    string `json:"self,omitempty"`
}

// RemoteIssueLink represents a remote link on an issue.
type RemoteIssueLink struct {
	ID           int          `json:"id,omitempty"`
	Self         string       `json:"self,omitempty"`
	GlobalID     string       `json:"globalId,omitempty"`
	Application  *Application `json:"application,omitempty"`
	Relationship string       `json:"relationship,omitempty"`
	Object       *RemoteObject `json:"object,omitempty"`
}

// Application represents a remote application.
type Application struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// RemoteObject represents a remote object in a remote link.
type RemoteObject struct {
	URL     string        `json:"url,omitempty"`
	Title   string        `json:"title,omitempty"`
	Summary string        `json:"summary,omitempty"`
	Icon    *Icon         `json:"icon,omitempty"`
	Status  *RemoteStatus `json:"status,omitempty"`
}

// Icon represents an icon.
type Icon struct {
	URL16x16 string `json:"url16x16,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
}

// RemoteStatus represents a remote status.
type RemoteStatus struct {
	Resolved bool `json:"resolved,omitempty"`
	Icon     *Icon `json:"icon,omitempty"`
}

// Filter represents a saved Jira filter.
type Filter struct {
	Self             string            `json:"self,omitempty"`
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Owner            *UserDetails      `json:"owner,omitempty"`
	JQL              string            `json:"jql,omitempty"`
	ViewURL          string            `json:"viewUrl,omitempty"`
	SearchURL        string            `json:"searchUrl,omitempty"`
	Favourite        bool              `json:"favourite,omitempty"`
	FavouritedCount  int               `json:"favouritedCount,omitempty"`
	SharePermissions []SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []SharePermission `json:"editPermissions,omitempty"`
}

// SharePermission represents a sharing permission.
type SharePermission struct {
	ID      int          `json:"id,omitempty"`
	Type    string       `json:"type,omitempty"`
	Project *Project     `json:"project,omitempty"`
	Role    *ProjectRole `json:"role,omitempty"`
	Group   *Group       `json:"group,omitempty"`
	User    *UserDetails `json:"user,omitempty"`
}

// ProjectRole represents a project role.
type ProjectRole struct {
	Self        string        `json:"self,omitempty"`
	ID          int           `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Actors      []RoleActor   `json:"actors,omitempty"`
}

// RoleActor represents an actor in a project role.
type RoleActor struct {
	ID          int          `json:"id,omitempty"`
	DisplayName string       `json:"displayName,omitempty"`
	Type        string       `json:"type,omitempty"`
	Name        string       `json:"name,omitempty"`
	ActorUser   *UserDetails `json:"actorUser,omitempty"`
	ActorGroup  *Group       `json:"actorGroup,omitempty"`
}

// Dashboard represents a Jira dashboard.
type Dashboard struct {
	Self             string            `json:"self,omitempty"`
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Owner            *UserDetails      `json:"owner,omitempty"`
	IsFavourite      bool              `json:"isFavourite,omitempty"`
	Popularity       int               `json:"popularity,omitempty"`
	Rank             int               `json:"rank,omitempty"`
	View             string            `json:"view,omitempty"`
	SharePermissions []SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []SharePermission `json:"editPermissions,omitempty"`
}

// DashboardList represents a paginated list of dashboards.
type DashboardList struct {
	Pagination
	Dashboards []Dashboard `json:"dashboards"`
}

// DashboardGadget represents a dashboard gadget.
type DashboardGadget struct {
	ID        int                 `json:"id,omitempty"`
	Color     string              `json:"color,omitempty"`
	Position  *DashboardGadgetPos `json:"position,omitempty"`
	Title     string              `json:"title,omitempty"`
	ModuleKey string              `json:"moduleKey,omitempty"`
	URI       string              `json:"uri,omitempty"`
}

// DashboardGadgetPos represents gadget position.
type DashboardGadgetPos struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// DashboardGadgetList is a list of gadgets.
type DashboardGadgetList struct {
	Gadgets []DashboardGadget `json:"gadgets"`
}

// Field represents a Jira field.
type Field struct {
	ID          string       `json:"id,omitempty"`
	Key         string       `json:"key,omitempty"`
	Name        string       `json:"name,omitempty"`
	Custom      bool         `json:"custom,omitempty"`
	Orderable   bool         `json:"orderable,omitempty"`
	Navigable   bool         `json:"navigable,omitempty"`
	Searchable  bool         `json:"searchable,omitempty"`
	ClauseNames []string     `json:"clauseNames,omitempty"`
	Schema      *FieldSchema `json:"schema,omitempty"`
}

// FieldSchema represents a field's schema.
type FieldSchema struct {
	Type     string `json:"type,omitempty"`
	Items    string `json:"items,omitempty"`
	System   string `json:"system,omitempty"`
	Custom   string `json:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty"`
}

// Workflow represents a Jira workflow.
type Workflow struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

// WorkflowScheme represents a workflow scheme.
type WorkflowScheme struct {
	ID                int               `json:"id,omitempty"`
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	DefaultWorkflow   string            `json:"defaultWorkflow,omitempty"`
	IssueTypeMappings map[string]string `json:"issueTypeMappings,omitempty"`
	Draft             *WorkflowScheme   `json:"draft,omitempty"`
	Self              string            `json:"self,omitempty"`
}

// Screen represents a Jira screen.
type Screen struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ScreenScheme represents a screen scheme.
type ScreenScheme struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Screens     *ScreenTypes `json:"screens,omitempty"`
}

// ScreenTypes maps operation types to screen IDs.
type ScreenTypes struct {
	Create  int `json:"create,omitempty"`
	Default int `json:"default,omitempty"`
	View    int `json:"view,omitempty"`
	Edit    int `json:"edit,omitempty"`
}

// ScreenTab represents a tab on a screen.
type ScreenTab struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ScreenField represents a field on a screen tab.
type ScreenField struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// PermissionScheme represents a permission scheme.
type PermissionScheme struct {
	ID          int               `json:"id,omitempty"`
	Self        string            `json:"self,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Permissions []PermissionGrant `json:"permissions,omitempty"`
}

// PermissionGrant represents a single permission grant.
type PermissionGrant struct {
	ID         int               `json:"id,omitempty"`
	Self       string            `json:"self,omitempty"`
	Holder     *PermissionHolder `json:"holder,omitempty"`
	Permission string            `json:"permission,omitempty"`
}

// PermissionHolder represents who holds a permission.
type PermissionHolder struct {
	Type      string `json:"type,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// NotificationScheme represents a notification scheme.
type NotificationScheme struct {
	ID          int    `json:"id,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// IssueSecurityScheme represents an issue security scheme.
type IssueSecurityScheme struct {
	Self                   string `json:"self,omitempty"`
	ID                     int    `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	DefaultSecurityLevelID int    `json:"defaultSecurityLevelId,omitempty"`
}

// FieldConfiguration represents a field configuration.
type FieldConfiguration struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

// FieldConfigurationScheme represents a field configuration scheme.
type FieldConfigurationScheme struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Webhook represents a Jira webhook.
type Webhook struct {
	ID                      int      `json:"id,omitempty"`
	JqlFilter               string   `json:"jqlFilter,omitempty"`
	FieldIDsFilter          []string `json:"fieldIdsFilter,omitempty"`
	IssuePropertyKeysFilter []string `json:"issuePropertyKeysFilter,omitempty"`
	Events                  []string `json:"events,omitempty"`
	ExpirationDate          int64    `json:"expirationDate,omitempty"`
}

// ServerInfo represents Jira server information.
type ServerInfo struct {
	BaseURL        string `json:"baseUrl,omitempty"`
	Version        string `json:"version,omitempty"`
	VersionNumbers []int  `json:"versionNumbers,omitempty"`
	BuildNumber    int    `json:"buildNumber,omitempty"`
	DeploymentType string `json:"deploymentType,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	ScmInfo        string `json:"scmInfo,omitempty"`
	ServerTitle    string `json:"serverTitle,omitempty"`
}

// BulkIssueCreateRequest represents a bulk issue creation request.
type BulkIssueCreateRequest struct {
	IssueUpdates []IssueUpdateDetails `json:"issueUpdates"`
}

// BulkIssueCreateResponse represents the response from bulk issue creation.
type BulkIssueCreateResponse struct {
	Issues []CreatedIssue `json:"issues"`
	Errors []interface{}  `json:"errors"`
}

// IssueNotifyRequest represents a notification request for an issue.
type IssueNotifyRequest struct {
	HTMLBody string          `json:"htmlBody,omitempty"`
	Subject  string          `json:"subject,omitempty"`
	TextBody string          `json:"textBody,omitempty"`
	To       *NotifyTo       `json:"to,omitempty"`
	Restrict *NotifyRestrict `json:"restrict,omitempty"`
}

// NotifyTo represents notification recipients.
type NotifyTo struct {
	Reporter bool          `json:"reporter,omitempty"`
	Assignee bool          `json:"assignee,omitempty"`
	Watchers bool          `json:"watchers,omitempty"`
	Voters   bool          `json:"voters,omitempty"`
	Users    []UserDetails `json:"users,omitempty"`
	Groups   []Group       `json:"groups,omitempty"`
}

// NotifyRestrict represents notification restrictions.
type NotifyRestrict struct {
	Groups      []Group      `json:"groups,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Permission represents a permission.
type Permission struct {
	ID   string `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
}

// EntityProperty represents an entity property.
type EntityProperty struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// CreateMeta represents issue create metadata.
type CreateMeta struct {
	Projects []CreateMetaProject `json:"projects,omitempty"`
}

// CreateMetaProject represents a project in create metadata.
type CreateMetaProject struct {
	Self       string                `json:"self,omitempty"`
	ID         string                `json:"id,omitempty"`
	Key        string                `json:"key,omitempty"`
	Name       string                `json:"name,omitempty"`
	IssueTypes []CreateMetaIssueType `json:"issuetypes,omitempty"`
}

// CreateMetaIssueType represents an issue type in create metadata.
type CreateMetaIssueType struct {
	Self        string                 `json:"self,omitempty"`
	ID          string                 `json:"id,omitempty"`
	Description string                 `json:"description,omitempty"`
	IconURL     string                 `json:"iconUrl,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Subtask     bool                   `json:"subtask,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// Changelog represents an issue changelog.
type Changelog struct {
	ID      string          `json:"id,omitempty"`
	Author  *UserDetails    `json:"author,omitempty"`
	Created string          `json:"created,omitempty"`
	Items   []ChangelogItem `json:"items,omitempty"`
}

// ChangelogItem represents a single change in a changelog.
type ChangelogItem struct {
	Field      string `json:"field,omitempty"`
	FieldType  string `json:"fieldtype,omitempty"`
	FieldID    string `json:"fieldId,omitempty"`
	From       string `json:"from,omitempty"`
	FromString string `json:"fromString,omitempty"`
	To         string `json:"to,omitempty"`
	ToString   string `json:"toString,omitempty"`
}

// ChangelogPage is a paginated list of changelogs.
type ChangelogPage struct {
	Pagination
	Histories []Changelog `json:"histories,omitempty"`
	Values    []Changelog `json:"values,omitempty"`
}

// ApplicationRole represents an application role.
type ApplicationRole struct {
	Key               string `json:"key,omitempty"`
	Name              string `json:"name,omitempty"`
	SelectedByDefault bool   `json:"selectedByDefault,omitempty"`
	Defined           bool   `json:"defined,omitempty"`
	NumberOfSeats     int    `json:"numberOfSeats,omitempty"`
	RemainingSeats    int    `json:"remainingSeats,omitempty"`
	UserCount         int    `json:"userCount,omitempty"`
	HasUnlimitedSeats bool   `json:"hasUnlimitedSeats,omitempty"`
	Platform          bool   `json:"platform,omitempty"`
}

// TaskResult represents an async task result.
type TaskResult struct {
	Self           string `json:"self,omitempty"`
	ID             string `json:"id,omitempty"`
	Description    string `json:"description,omitempty"`
	Status         string `json:"status,omitempty"`
	Result         string `json:"result,omitempty"`
	Progress       int    `json:"progress,omitempty"`
	ElapsedRuntime int    `json:"elapsedRuntime,omitempty"`
	Submitted      int64  `json:"submitted,omitempty"`
	Started        int64  `json:"started,omitempty"`
	Finished       int64  `json:"finished,omitempty"`
	LastUpdate     int64  `json:"lastUpdate,omitempty"`
}

// AttachmentMeta represents attachment metadata/settings.
type AttachmentMeta struct {
	Enabled     bool `json:"enabled,omitempty"`
	UploadLimit int  `json:"uploadLimit,omitempty"`
}

// VersionIssueCounts represents issue counts for a version.
type VersionIssueCounts struct {
	Self                string `json:"self,omitempty"`
	IssuesFixedCount    int    `json:"issuesFixedCount,omitempty"`
	IssuesAffectedCount int    `json:"issuesAffectedCount,omitempty"`
}

// VersionUnresolvedIssueCount represents unresolved issue count for a version.
type VersionUnresolvedIssueCount struct {
	Self                  string `json:"self,omitempty"`
	IssuesUnresolvedCount int    `json:"issuesUnresolvedCount,omitempty"`
	IssuesCount           int    `json:"issuesCount,omitempty"`
}

// ComponentIssueCount represents issue counts related to a component.
type ComponentIssueCount struct {
	Self       string `json:"self,omitempty"`
	IssueCount int    `json:"issueCount,omitempty"`
}

// GroupMembers represents members of a group.
type GroupMembers struct {
	Pagination
	IsLast bool          `json:"isLast,omitempty"`
	Values []UserDetails `json:"values,omitempty"`
}

// FoundGroups represents group picker results.
type FoundGroups struct {
	Header string       `json:"header,omitempty"`
	Total  int          `json:"total,omitempty"`
	Groups []FoundGroup `json:"groups,omitempty"`
}

// FoundGroup represents a group in picker results.
type FoundGroup struct {
	Name    string       `json:"name,omitempty"`
	HTML    string       `json:"html,omitempty"`
	Labels  []GroupLabel `json:"labels,omitempty"`
	GroupID string       `json:"groupId,omitempty"`
}

// GroupLabel represents a group label.
type GroupLabel struct {
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

// AuditRecords represents audit records.
type AuditRecords struct {
	Offset  int           `json:"offset,omitempty"`
	Limit   int           `json:"limit,omitempty"`
	Total   int           `json:"total,omitempty"`
	Records []AuditRecord `json:"records,omitempty"`
}

// AuditRecord represents a single audit record.
type AuditRecord struct {
	ID            int                   `json:"id,omitempty"`
	Summary       string                `json:"summary,omitempty"`
	RemoteAddress string                `json:"remoteAddress,omitempty"`
	Created       string                `json:"created,omitempty"`
	Category      string                `json:"category,omitempty"`
	EventSource   string                `json:"eventSource,omitempty"`
	Description   string                `json:"description,omitempty"`
}

// AnnouncementBanner represents the announcement banner settings.
type AnnouncementBanner struct {
	Message       string `json:"message,omitempty"`
	IsDismissible bool   `json:"isDismissible,omitempty"`
	IsEnabled     bool   `json:"isEnabled,omitempty"`
	HashCode      string `json:"hashCode,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
}

// Configuration represents Jira configuration.
type Configuration struct {
	VotingEnabled           bool `json:"votingEnabled,omitempty"`
	WatchingEnabled         bool `json:"watchingEnabled,omitempty"`
	UnassignedIssuesAllowed bool `json:"unassignedIssuesAllowed,omitempty"`
	SubTasksEnabled         bool `json:"subTasksEnabled,omitempty"`
	IssueLinkingEnabled     bool `json:"issueLinkingEnabled,omitempty"`
	TimeTrackingEnabled     bool `json:"timeTrackingEnabled,omitempty"`
	AttachmentsEnabled      bool `json:"attachmentsEnabled,omitempty"`
}

// TimeTrackingConfiguration represents time tracking configuration.
type TimeTrackingConfiguration struct {
	WorkingHoursPerDay float64 `json:"workingHoursPerDay,omitempty"`
	WorkingDaysPerWeek float64 `json:"workingDaysPerWeek,omitempty"`
	TimeFormat         string  `json:"timeFormat,omitempty"`
	DefaultUnit        string  `json:"defaultUnit,omitempty"`
}

// TimeTrackingProvider represents a time tracking provider.
type TimeTrackingProvider struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// IssueTypeScheme represents an issue type scheme.
type IssueTypeScheme struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	DefaultIssueTypeID string `json:"defaultIssueTypeId,omitempty"`
	IsDefault          bool   `json:"isDefault,omitempty"`
}

// IssueTypeScreenScheme represents an issue type screen scheme.
type IssueTypeScreenScheme struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Plan represents a Jira plan (Advanced Roadmaps).
type Plan struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	IsArchived    bool   `json:"isArchived,omitempty"`
	LeadAccountID string `json:"leadAccountId,omitempty"`
}

// Status represents a workflow status (for CRUD operations).
type Status struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	StatusCategory string `json:"statusCategory,omitempty"`
}

// UIModification represents a UI modification.
type UIModification struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

// ColumnItem represents a column configuration item.
type ColumnItem struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

// UserPermission represents a user's permissions.
type UserPermission struct {
	ID             string `json:"id,omitempty"`
	Key            string `json:"key,omitempty"`
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	Description    string `json:"description,omitempty"`
	HavePermission bool   `json:"havePermission,omitempty"`
}

// BulkPermissionsRequest represents a bulk permissions check request.
type BulkPermissionsRequest struct {
	ProjectPermissions []BulkProjectPermissions `json:"projectPermissions,omitempty"`
	GlobalPermissions  []string                 `json:"globalPermissions,omitempty"`
	AccountID          string                   `json:"accountId,omitempty"`
}

// BulkProjectPermissions represents project permissions in a bulk check.
type BulkProjectPermissions struct {
	Issues      []int    `json:"issues,omitempty"`
	Projects    []int    `json:"projects,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// SecurityLevel represents an issue security level.
type SecurityLevel struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

// ProjectFeature represents a project feature.
type ProjectFeature struct {
	ProjectID            int    `json:"projectId,omitempty"`
	State                string `json:"state,omitempty"`
	ToggleLocked         bool   `json:"toggleLocked,omitempty"`
	Feature              string `json:"feature,omitempty"`
	ImageURI             string `json:"imageUri,omitempty"`
	LocalisedName        string `json:"localisedName,omitempty"`
	LocalisedDescription string `json:"localisedDescription,omitempty"`
}

// ProjectFeaturesResponse is the response for project features.
type ProjectFeaturesResponse struct {
	Features []ProjectFeature `json:"features,omitempty"`
}

// FieldConfigScheme represents a field scheme (not field configuration scheme).
type FieldConfigScheme struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// --- Agile API types ---

// Board represents a Jira Agile board.
type Board struct {
	ID       int          `json:"id"`
	Self     string       `json:"self,omitempty"`
	Name     string       `json:"name,omitempty"`
	Type     string       `json:"type,omitempty"`
	Location *BoardLocation `json:"location,omitempty"`
}

// BoardLocation represents the project location of a board.
type BoardLocation struct {
	ProjectID      int    `json:"projectId,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	ProjectName    string `json:"projectName,omitempty"`
	ProjectKey     string `json:"projectKey,omitempty"`
	ProjectTypeKey string `json:"projectTypeKey,omitempty"`
	Name           string `json:"name,omitempty"`
}

// BoardList is the paginated response for boards.
type BoardList struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

// BoardConfiguration represents a board's configuration.
type BoardConfiguration struct {
	ID            int                    `json:"id"`
	Name          string                 `json:"name,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Self          string                 `json:"self,omitempty"`
	Location      *BoardLocation         `json:"location,omitempty"`
	Filter        *BoardFilter           `json:"filter,omitempty"`
	ColumnConfig  *BoardColumnConfig     `json:"columnConfig,omitempty"`
	Ranking       map[string]interface{} `json:"ranking,omitempty"`
}

// BoardFilter represents the filter associated with a board.
type BoardFilter struct {
	ID   string `json:"id,omitempty"`
	Self string `json:"self,omitempty"`
}

// BoardColumnConfig represents board column configuration.
type BoardColumnConfig struct {
	Columns    []BoardColumn `json:"columns,omitempty"`
	ConstraintType string   `json:"constraintType,omitempty"`
}

// BoardColumn represents a single column on a board.
type BoardColumn struct {
	Name     string          `json:"name,omitempty"`
	Statuses []BoardStatus   `json:"statuses,omitempty"`
	Min      int             `json:"min,omitempty"`
	Max      int             `json:"max,omitempty"`
}

// BoardStatus represents a status in a board column.
type BoardStatus struct {
	ID   string `json:"id,omitempty"`
	Self string `json:"self,omitempty"`
}

// Sprint represents a Jira Agile sprint.
type Sprint struct {
	ID            int    `json:"id"`
	Self          string `json:"self,omitempty"`
	State         string `json:"state,omitempty"`
	Name          string `json:"name,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompleteDate  string `json:"completeDate,omitempty"`
	OriginBoardID int    `json:"originBoardId,omitempty"`
	Goal          string `json:"goal,omitempty"`
}

// SprintList is the paginated response for sprints.
type SprintList struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

// SprintIssuesResponse is the response for issues in a sprint/board.
type SprintIssuesResponse struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// Epic represents a Jira Agile epic.
type Epic struct {
	ID      int    `json:"id"`
	Key     string `json:"key,omitempty"`
	Self    string `json:"self,omitempty"`
	Name    string `json:"name,omitempty"`
	Summary string `json:"summary,omitempty"`
	Done    bool   `json:"done,omitempty"`
}

// EpicList is the paginated response for epics.
type EpicList struct {
	MaxResults int    `json:"maxResults"`
	StartAt    int    `json:"startAt"`
	Total      int    `json:"total"`
	IsLast     bool   `json:"isLast"`
	Values     []Epic `json:"values"`
}

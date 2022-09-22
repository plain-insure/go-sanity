package sanity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ProjectsService is a client for the Sanity Projects API.
//
// Refer to https://www.sanity.io/docs/projects-api for more information.
type ProjectsService service

// -----------------------------------------------------------------------------
// Projects

// A Project is a Sanity project that appears in your Sanity account.
type Project struct {
	// Id is the unique identifier for the project.
	Id string `json:"id"`

	// DisplayName is the user-friendly name for the project.
	// This is the name presented on the Sanity dashboard.
	DisplayName string `json:"displayName"`

	// StudioHost is the hostname for a studio deployment.
	// A complete url has the form `https://<hostname>.sanity.studio/`, but note
	// that this field is just the hostname. This field may be empty if the studio
	// has not been deployed.
	StudioHost string `json:"studioHost,omitempty"`

	// OrganizationId is the id of the organization that owns the project.
	OrganizationId string `json:"organizationId,omitempty"`

	// Metadata about the project.
	//
	// May include the following fields:
	//   `color`: a hex string that describes the color of the project logo shown on the Sanity dashboard.
	//   `externalStudioHost`: the URL of the Sanity studio if it is deployed outside of Sanity
	Metadata map[string]string `json:"metadata"`

	// MaxRetentionDays is the amount of time revisions are stored before they are
	// deleted.
	//
	// See also: https://www.sanity.io/docs/history-experience
	MaxRetentionDays int `json:"maxRetentionDays,omitempty"`

	DataClass string `json:"dataClass,omitempty"`

	IsBlocked bool `json:"isBlocked"`

	IsDisabled bool `json:"isDisabled"`

	// IsDisabledByUser indicates whether the project is archived.
	IsDisabledByUser bool `json:"isDisabledByUser"`

	// ActivityFeedEnabled indicates whether changes to the project are reflected
	// on the Sanity dashboard.
	ActivityFeedEnabled bool `json:"activityFeedEnabled"`

	// CreatedAt is the creation time of the project.
	CreatedAt time.Time `json:"createdAt"`

	// Members contains information about the project members and their roles.
	Members []Member `json:"members"`

	// Features is a list of feature names that are enabled for the project.
	Features []string `json:"features,omitempty"`

	// PendingInvites is the number of outstanding invitations for people to join
	// the project as members.
	PendingInvites int `json:"pendingInvites,omitempty"`
}

// A Member is an account that may access a project in some capacity.
type Member struct {
	// Id is the unique identifier for the member.
	Id string `json:"id"`

	// CreatedAt is the creation time of the member.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is the last time the member was updated.
	UpdatedAt time.Time `json:"updatedAt"`

	IsCurrentUser bool `json:"isCurrentUser"`

	// IsRobot indicates whether the member is a robot user.
	IsRobot bool `json:"isRobot"`

	// Roles
	Roles []Role `json:"roles"`
}

// A Role describes the type of access for members assigned the role.
type Role struct {
	// The Name of the role (e.g., `administrator`).
	Name string `json:"name"`

	// The Title is the display-friendly name of the role (e.g., `Administrator`).
	Title string `json:"title"`

	// Description is a display-friendly short text that explains the roles
	// capabilities.
	Description string `json:"description,omitempty"`
}

// List fetches and returns all the projects.
func (s *ProjectsService) List(ctx context.Context) ([]Project, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects", s.client.baseURL)

	var projects []Project
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &projects)

	return projects, err
}

type CreateProjectRequest struct {
	// DisplayName is the user-friendly name for the project.
	// This is the name presented on the Sanity dashboard.
	DisplayName string `json:"displayName"`

	// OrganizationId is the id of the organization that owns the project. If left
	// blank, the project will be created in the personal account of the
	// authenticated user.
	OrganizationId string `json:"organizationId,omitempty"`
}

// Create generates a new project in Sanity.
func (s *ProjectsService) Create(ctx context.Context, r *CreateProjectRequest) (*Project, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects", s.client.baseURL)

	var project Project
	err := do(ctx, s.client.client, url, http.MethodPost, r, &project)

	return &project, err
}

// Get fetches a project by its unique identifier.
func (s *ProjectsService) Get(ctx context.Context, projectId string) (*Project, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s", s.client.baseURL, projectId)

	var project Project
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &project)

	return &project, err
}

type UpdateProjectRequest struct {
	// DisplayName is the user-friendly name for the project.
	// This is the name presented on the Sanity dashboard.
	DisplayName string

	// StudioHost is the hostname for a studio deployment.
	// A complete url has the form `https://<hostname>.sanity.studio/`, but note
	// that this field is just the hostname.
	//
	// Important note: This is a one-time operation. Once the StudioHost value has
	// been set, further attempts to change it will fail.
	StudioHost string

	// Color is a hex string that describes the color of the project logo shown on
	// the Sanity dashboard.
	Color string

	// ExternalStudioHost is the URL of the Sanity studio if it is deployed
	// outside of Sanity.
	ExternalStudioHost string

	// IsDisabledByUser indicates whether the project is archived.
	IsDisabledByUser *bool

	// ActivityFeedEnabled indicates whether changes to the project are reflected
	// on the Sanity dashboard.
	ActivityFeedEnabled *bool
}

func (r *UpdateProjectRequest) MarshalJSON() ([]byte, error) {
	type request struct {
		DisplayName         string            `json:"displayName,omitempty"`
		StudioHost          string            `json:"studioHost,omitempty"`
		Metadata            map[string]string `json:"metadata,omitempty"`
		IsDisabledByUser    *bool             `json:"isDisabledByUser,omitempty"`
		ActivityFeedEnabled *bool             `json:"activityFeedEnabled,omitempty"`
	}

	req := &request{
		DisplayName:         r.DisplayName,
		StudioHost:          r.StudioHost,
		Metadata:            make(map[string]string),
		IsDisabledByUser:    r.IsDisabledByUser,
		ActivityFeedEnabled: r.ActivityFeedEnabled,
	}
	if r.Color != "" {
		req.Metadata["color"] = strings.ToLower(r.Color) // if upper case, API returns a 400
	}
	if r.ExternalStudioHost != "" {
		req.Metadata["externalHost"] = r.ExternalStudioHost
	}

	return json.Marshal(req)
}

// Update applies the requested changes to the specified project.
//
// Note that zero valeus in the update request are ignored.
func (s *ProjectsService) Update(ctx context.Context, projectId string, r *UpdateProjectRequest) (*Project, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s", s.client.baseURL, projectId)

	var project Project
	err := do(ctx, s.client.client, url, http.MethodPatch, r, &project)

	return &project, err
}

// Delete destroys the project without additional prompt.
func (s *ProjectsService) Delete(ctx context.Context, projectId string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s", s.client.baseURL, projectId)

	type response struct {
		Deleted bool `json:"deleted"`
	}

	var resp response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &resp)
	return resp.Deleted, err
}

// -----------------------------------------------------------------------------
// CORS

// A CORSEntry represents an allowed CORS origin for a project.
type CORSEntry struct {
	// Id is the unique idenifiter for the entry.
	Id int64 `json:"id"`

	// Origin is the full URL for the CORS entry, e.g., `http://localhost:3333`.
	// Supports wildcards with `*`.
	Origin string `json:"origin"`

	// AllowCredentials indicates whether the origin may make authenticated
	// requests with a token. This is required if hosting a studio instance at the
	// origin.
	AllowCredentials bool `json:"allowCredentials"`

	// CreatedAt is the time the entry was created.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is the time the entry was last updated.
	UpdatedAt time.Time `json:"updatedAt"`

	// ProjectId is the identifier of the project this entry belongs to.
	ProjectId string `json:"projectId"`
}

// ListCORSEntries fetches and returns all CORS entries for the specified project.
func (s *ProjectsService) ListCORSEntries(ctx context.Context, projectId string) ([]CORSEntry, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/cors", s.client.baseURL, projectId)

	var entries []CORSEntry
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &entries)

	return entries, err
}

type CreateCORSEntryRequest struct {
	// Origin is the full URL for the CORS entry, e.g., `http://localhost:3333`.
	// Supports wildcards with `*`.
	Origin string `json:"origin"`

	// AllowCredentials indicates whether the origin may make authenticated
	// requests with a token. This is required if hosting a studio instance at the
	// origin.
	AllowCredentials *bool `json:"allowCredentials,omitempty"`
}

// CreateCORSEntry will add a new CORS entry to the specified Sanity project.
func (s *ProjectsService) CreateCORSEntry(ctx context.Context, projectId string, r *CreateCORSEntryRequest) (*CORSEntry, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/cors", s.client.baseURL, projectId)

	var entry CORSEntry
	err := do(ctx, s.client.client, url, http.MethodPost, r, &entry)

	return &entry, err
}

// DeleteCORSEntry removes the specified entry from the project.
func (s *ProjectsService) DeleteCORSEntry(ctx context.Context, projectId string, entryId int64) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/cors/%d", s.client.baseURL, projectId, entryId)

	type response struct {
		Id      int64 `json:"id"`
		Deleted bool  `json:"deleted"`
	}

	var res response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &res)

	return res.Deleted, err
}

// -----------------------------------------------------------------------------
// Datasets

const (
	AclModePublic  = "public"
	AclModePrivate = "private"
)

// A Dataset represents a collection of documents and assets within a project.
type Dataset struct {
	// Name is the name of the dataset and serves as the unique identifier for
	// this dataset in the project.
	Name string `json:"name"`

	// AclMode describes whether the dataset is accessible publicly or privately.
	// If available privately, the data in the dataset is only accessible via a
	// token.
	AclMode string `json:"aclMode"`
}

// ListDatasets fetches and returns all the datasets in the specified project.
func (s *ProjectsService) ListDatasets(ctx context.Context, projectId string) ([]Dataset, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets", s.client.baseURL, projectId)

	var datasets []Dataset
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &datasets)

	return datasets, err
}

type CreateDatasetRequest struct {
	// Name is the name of the dataset and serves as the unique identifier for
	// this dataset in the project.
	Name string `json:"-"`

	// AclMode describes whether the dataset is accessible publicly or privately.
	// If available privately, the data in the dataset is only accessible via a
	// token.
	AclMode string `json:"aclMode,omitempty"`
}

// CreateDataset adds a new dataset to the Sanity project.
func (s *ProjectsService) CreateDataset(ctx context.Context, projectId string, r *CreateDatasetRequest) (*Dataset, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s", s.client.baseURL, projectId, r.Name)

	if strings.Contains(r.Name, " ") {
		return nil, errors.New("name cannot contain spaces")
	}

	type response struct {
		Name    string `json:"datasetName"`
		AclMode string `json:"aclMode"`
	}

	var resp response
	err := do(ctx, s.client.client, url, http.MethodPut, r, &resp)

	if err != nil {
		return nil, err
	}

	return &Dataset{Name: resp.Name, AclMode: resp.AclMode}, nil
}

type CopyDatasetRequest struct {
	// SourceDataset is the name of the dataset to be copied from.
	SourceDataset string `json:"-"`

	// TargetDataset is the name of the dataset to be copied to.
	TargetDataset string `json:"targetDataset"`
}

type CopyDatasetResponse struct {
	Name    string `json:"datasetName"`
	Message string `json:"message"`
	AclMode string `json:"aclMode"`
	JobId   string `json:"jobId"`
}

// CopyDataset copies data from one dataset into another.
//
// NOTE: This is enterprise feature and is only available for business and
// enterprise plans.
func (s *ProjectsService) CopyDataset(ctx context.Context, projectId string, r *CopyDatasetRequest) (*CopyDatasetResponse, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s/copy", s.client.baseURL, projectId, r.SourceDataset)

	var response CopyDatasetResponse
	err := do(ctx, s.client.client, url, http.MethodPut, r, &response)

	return &response, err
}

// DeleteDataset removes the specified dataset from the project without prompt.
func (s *ProjectsService) DeleteDataset(ctx context.Context, projectId string, datasetName string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s", s.client.baseURL, projectId, datasetName)

	type response struct {
		Deleted bool `json:"deleted"`
	}

	var res response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &res)

	return res.Deleted, err
}

// -----------------------------------------------------------------------------
// Features

// ListActiveFeatures fetches and returns a list of all active features on the
// specified project.
func (s *ProjectsService) ListActiveFeatures(ctx context.Context, projectId string) ([]string, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/features", s.client.baseURL, projectId)

	var features []string
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &features)

	return features, err
}

// CheckFeatureActive accepts a project id and a feature name and returns a
// value indicating whether that feature is active on the specified project.
//
// Currently works with features named `privateDataset` and `thirdPartyLogin`.
func (s *ProjectsService) CheckFeatureActive(ctx context.Context, projectId string, featureName string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/features/%s", s.client.baseURL, projectId, featureName)

	active := false
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &active)

	return active, err
}

// -----------------------------------------------------------------------------
// Users and roles

// ListPermissions returns a list of permissions that the authenticated user
// has for the specified project.
func (s *ProjectsService) ListPermissions(ctx context.Context, projectId string) ([]string, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/permissions", s.client.baseURL, projectId)

	var permissions []string
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &permissions)

	return permissions, err
}

type User struct {
	// Id is the unique identifier for the user.
	Id string `json:"id"`

	SanityUserId string `json:"sanityUserId"`

	// ProjectId is the unique ID for the project to which the specified user is
	// assigned.
	ProjectId string `json:"projectId"`

	// DisplayName is the user's full name.
	DisplayName string `json:"displayName"`

	// Family name is the user's last name.
	FamilyName string `json:"familyName"`

	// GivenName is the user's first name.
	GivenName string `json:"givenName"`

	// MiddleName is the user's middle name.
	MiddleName string `json:"middleName,omitempty"`

	// ImageURL is a url pointing to an image for the user.
	ImageURL string `json:"imageUrl,omitempty"`

	// CreatedAt is the time the user was created for the project.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is the time the user was last updated.
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetUser fetches and returns information about a user on a project.
func (s *ProjectsService) GetUser(ctx context.Context, projectId string, userId string) (*User, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/users/%s", s.client.baseURL, projectId, userId)

	var user User
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &user)

	return &user, err
}

type ProjectRole struct {
	// Id is the identifier for the role. This may be an empty string if the role
	// is one of the default roles created by Sanity, such as the `administrator`,
	// `deploy-studio`, `editor`, and `viewer`.
	Id string `json:"id"`

	// Name is the name of the role.
	Name string `json:"name"`

	// Description explains the permissions associated with the role.
	Description string `json:"description"`

	IsRootRole bool `json:"isRootRole"`

	ReadOnly bool `json:"readOnly"`

	IsListed bool `json:"isListed"`

	Permissions []string `json:"permissions"`
}

// ListProjectRoles fetches and returns the roles associated with the specified
// project.
func (s *ProjectsService) ListProjectRoles(ctx context.Context, projectId string) ([]ProjectRole, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/roles", s.client.baseURL, projectId)

	var roles []ProjectRole
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &roles)

	return roles, err
}

// -----------------------------------------------------------------------------
// Tokens

type ProjectToken struct {
	// Id is the unique identifier for the token.
	Id string `json:"id"`

	// Label is a descriptive name for the token.
	Label string `json:"label"`

	// ProjectUserId is the id of the user for whom the token was created.
	ProjectUserId string `json:"projectUserId"`

	// CreatedAt is the time the token was created.
	CreatedAt time.Time `json:"createdAt"`

	// Roles describe the various roles associated with the token.
	Roles []Role `json:"roles"`
}

// ListProjectTokens fetches and returns all access tokens associated with the
// specified project.
func (s *ProjectsService) ListProjectTokens(ctx context.Context, projectId string) ([]ProjectToken, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tokens", s.client.baseURL, projectId)

	var tokens []ProjectToken
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &tokens)

	return tokens, err
}

type CreateProjectTokenRequest struct {
	// Label is a descriptive name for the token.
	Label string `json:"label"`

	// The name of the role to assign to the token. On a free plan, it must be
	// one of the following values: `viewer`, `editor`, or `deploy-studio`.
	RoleName string `json:"roleName"`
}

type CreateProjectTokenResponse struct {
	ProjectToken

	// Key is the access token. This value can only be returned once from the API
	// and should be treated as a secret value.
	Key string `json:"key"`
}

// CreateProjectToken creates a new token for the specified project. It is
// important to note that the `Key` value in the response can only be returned
// from the API once, and the value should be treated as a secret value.
func (s *ProjectsService) CreateProjectToken(ctx context.Context, projectId string, r *CreateProjectTokenRequest) (*CreateProjectTokenResponse, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tokens", s.client.baseURL, projectId)

	var response CreateProjectTokenResponse
	err := do(ctx, s.client.client, url, http.MethodPost, r, &response)

	return &response, err
}

// DeleteProjectToken deletes the specified token without prompt.
func (s *ProjectsService) DeleteProjectToken(ctx context.Context, projectId string, tokenId string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tokens/%s", s.client.baseURL, projectId, tokenId)

	type response struct {
		Id          string            `json:"id"`
		Deleted     bool              `json:"deleted"`
		Description string            `json:"description,omitempty"`
		Metadata    map[string]string `json:"metadata,omitempty"`
	}

	var resp response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &resp)

	return resp.Deleted, err
}

// -----------------------------------------------------------------------------
// Dataset tags

type DatasetTag struct {
	// Name is the name of the tag and also serves as the tag's unique identifier.
	Name string `json:"name"`

	// Title is a display-friendly label for the tag.
	Title string `json:"title"`
}

// ListDatasetTags gets a list of all tags associated with the specified dataset.
func (s *ProjectsService) ListsDatasetTags(ctx context.Context, projectId, datasetName string) ([]DatasetTag, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s/tags", s.client.baseURL, projectId, datasetName)

	var tags []DatasetTag
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &tags)

	return tags, err
}

const (
	ToneDefault     = "default"
	TonePrimary     = "primary"
	TonePositive    = "positive"
	ToneCaution     = "caution"
	ToneCritical    = "critical"
	ToneTransparent = "transparent"
)

type CreateDatasetTagRequest struct {
	// Name is the name of the tag and also serves as the tag's unique identifier.
	Name string

	// Title is a display-friendly label for the tag.
	Title string

	// Description is a short descriptive text describing the tag.
	Description string

	// Tone is the color of the tag. Valid values are represented as the `Tone*`
	// constants in this package.
	Tone string
}

func (r *CreateDatasetTagRequest) MarshalJSON() ([]byte, error) {
	if r.Name == "" {
		return nil, errors.New("name is required")
	}
	if r.Title == "" {
		return nil, errors.New("title is required")
	}

	type request struct {
		Name        string            `json:"name"`
		Title       string            `json:"title"`
		Description string            `json:"description,omitempty"`
		Metadata    map[string]string `json:"metadata,omitempty"`
	}

	req := &request{
		Name:        r.Name,
		Title:       r.Title,
		Description: r.Description,
		Metadata:    make(map[string]string),
	}
	if r.Tone != "" {
		req.Metadata["tone"] = r.Tone
	}

	return json.Marshal(req)
}

// CreateDatasetTag creates and returns a new tag.
func (s *ProjectsService) CreateDatasetTag(ctx context.Context, projectId string, r *CreateDatasetTagRequest) (*DatasetTag, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tags", s.client.baseURL, projectId)

	var tag DatasetTag
	err := do(ctx, s.client.client, url, http.MethodPost, r, &tag)

	return &tag, err
}

type EditDatasetTagRequest struct {
	// Name is the name of the tag and also serves as the tag's unique identifier.
	Name string

	// Title is a display-friendly label for the tag.
	Title string

	// Description is a short descriptive text describing the tag.
	Description string

	// Tone is the color of the tag. Valid values are represented as the `Tone*`
	// constants in this package.
	Tone string
}

func (r *EditDatasetTagRequest) MarshalJSON() ([]byte, error) {
	type request struct {
		Name        string            `json:"name"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Metadata    map[string]string `json:"metadata,omitempty"`
	}

	req := &request{
		Name:        r.Name,
		Title:       r.Title,
		Description: r.Description,
		Metadata:    make(map[string]string),
	}
	if r.Tone != "" {
		req.Metadata["tone"] = r.Tone
	}

	return json.Marshal(req)
}

// EditDatasetTag updates and returns the specified tag.
func (s *ProjectsService) EditDatasetTag(ctx context.Context, projectId, tagIdentifier string, r *EditDatasetTagRequest) (*DatasetTag, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tags/%s", s.client.baseURL, projectId, tagIdentifier)

	var tag DatasetTag
	err := do(ctx, s.client.client, url, http.MethodPut, r, &tag)

	return &tag, err
}

// AssignDatasetTag assigns the specified tag to the dataset.
func (s *ProjectsService) AssignDatasetTag(ctx context.Context, projectId, datasetName, tagIdentifier string) error {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s/tags/%s", s.client.baseURL, projectId, datasetName, tagIdentifier)

	var x any
	return do(ctx, s.client.client, url, http.MethodPut, nil, &x)
}

// AssignDatasetTag removes the specified tag from the dataset.
func (s *ProjectsService) UnassignDatasetTag(ctx context.Context, projectId, datasetName, tagIdentifier string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/datasets/%s/tags/%s", s.client.baseURL, projectId, datasetName, tagIdentifier)

	type response struct {
		Deleted bool `json:"deleted"`
	}
	var resp response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &resp)

	return resp.Deleted, err
}

// DeleteDatasetTag destroys the tag without prompt. In order for this operation
// to be successful, the tag must first be removed from all datasets.
func (s *ProjectsService) DeleteDatasetTag(ctx context.Context, projectId, tagIdentifier string) (bool, error) {
	url := fmt.Sprintf("%s/v2021-06-07/projects/%s/tags/%s", s.client.baseURL, projectId, tagIdentifier)

	type response struct {
		Deleted bool `json:"deleted"`
	}
	var resp response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &resp)

	return resp.Deleted, err
}

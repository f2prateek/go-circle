package circle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CircleCI interface {
	// Provides information about the authenticated user.
	// https://circleci.com/docs/api#user
	// https://circleci.com/api/v1/me
	Me() (Me, error)

	// Provides information about projects followed by the authenticated user.
	//
	// https://circleci.com/docs/api#projects
	// https://circleci.com/api/v1/projects
	Projects() ([]Project, error)

	// Provides Build summary for each of the last 30 recent builds, ordered by BuildNum.
	//
	// https://circleci.com/docs/api#recent-builds
	// https://circleci.com/api/v1/recent-builds
	RecentBuilds() ([]BuildSummary, error)

	// Provides build summary for each of the last 30 builds for a single git repo.
	//
	// https://circleci.com/docs/api#recent-builds-project
	// https://circleci.com/api/v1/project/{username}/{project}
	RecentBuildsForProject(username, project string) ([]BuildSummary, error)

	// Provides build summary for each of the last 30 builds for a single branch of a
	// github branch.
	//
	// https://circleci.com/docs/api#recent-builds-project
	// https://circleci.com/api/v1/project/{username}/{project}
	RecentBuildsForProjectBranch(username, project, branch string, opts RecentBuildsOptions) ([]BuildSummary, error)

	// Provides a detailed build summary for the given build for the project.
	//
	// https://circleci.com/docs/api#build
	// https://circleci.com/api/v1/project/{username}/{project}/{num}
	BuildSummary(username, project string, num int) (DetailedBuildSummary, error)

	// List the artifacts produced by the given build.
	//
	// https://circleci.com/docs/api#build-artifacts
	// https://circleci.com/api/v1/project/{username}/{project}/{num}/artifacts
	Artifacts(username, project string, num int) ([]Artifact, error)

	// Retries the build and returns a summary of the new build.
	//
	// https://circleci.com/docs/api#retry-build
	// https://circleci.com/api/v1/project/{username}/{project}/{num}/retry
	Retry(username, project string, num int) (Build, error)

	// Cancels the build and returns a summary of the build.
	//
	// https://circleci.com/docs/api#cancel-build
	// https://circleci.com/api/v1/project/{username}/{project}/{num}/cancel
	Cancel(username, project string, num int) (Build, error)

	// Triggers a new build and returns a summary of the build.
	//
	// https://circleci.com/docs/api#new-build
	// https://circleci.com/api/v1/project/{username}/{project}/tree/{branch}
	Build(username, project, branch string) (Build, error)

	// Clears the cache for a project
	//
	// https://circleci.com/docs/api#clear-cache
	// https://circleci.com/api/v1/project/{username}/{project}/build-cache
	ClearCache(username, project string) (ClearCacheResponse, error)
}

type client struct {
	token string
	http  *http.Client
}

// New returns a Client for the given `token`.
func New(token string) CircleCI {
	return &client{token, http.DefaultClient}
}

func (c *client) endpoint(endpoint string) string {
	return fmt.Sprintf("https://circleci.com/api/v1%s?circle-token=%s", endpoint, c.token)
}

// Information about the authenticated user.
type Me struct {
	Admin               bool        `json:"admin"`
	Emails              []string    `json:"all_emails"`
	AvatarURL           string      `json:"avatar_url"`
	BasicEmailPrefs     string      `json:"basic_email_prefs"`
	Containers          int         `json:"containers"`
	CreatedAt           string      `json:"created_at"`
	DaysLeftInTrial     int         `json:"days_left_in_trial"`
	DevAdmin            bool        `json:"dev_admin"`
	GithubID            int         `json:"github_id"`
	GithubOauthScopes   []string    `json:"github_oauth_scopes"`
	GravatarID          interface{} `json:"gravatar_id"`
	HerokuAPIKey        interface{} `json:"heroku_api_key"`
	LastViewedChangelog string      `json:"last_viewed_changelog"`
	Login               string      `json:"login"`
	Name                string      `json:"name"`
	Parallelism         int         `json:"parallelism"`
	Plan                interface{} `json:"plan"`
	Projects            map[string]struct {
		Emails      string `json:"emails"`
		OnDashboard bool   `json:"on_dashboard"`
	} `json:"projects"`
	SelectedEmail string `json:"selected_email"`
	SignInCount   int    `json:"sign_in_count"`
	TrialEnd      string `json:"trial_end"`
}

func (c *client) Me() (Me, error) {
	url := c.endpoint("/me")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Me{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return Me{}, err
	}

	var m Me
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		return Me{}, err
	}

	return m, nil
}

// Information about a project.
type Project struct {
	AWS struct {
		KeyPair interface{} `json:"keypair"`
	} `json:"aws"`
	Branches map[string]struct {
		AddedAt     string `json:"added_at"`
		BuildNum    int    `json:"build_num"`
		Outcome     string `json:"outcome"`
		PushedAt    string `json:"pushed_at"`
		Status      string `json:"status"`
		VcsRevision string `json:"vcs_revision"`
	} `json:"branches"`
	CampfireNotifyPrefs interface{}     `json:"campfire_notify_prefs"`
	CampfireRoom        interface{}     `json:"campfire_room"`
	CampfireSubdomain   interface{}     `json:"campfire_subdomain"`
	CampfireToken       interface{}     `json:"campfire_token"`
	Compile             string          `json:"compile"`
	DefaultBranch       string          `json:"default_branch"`
	Dependencies        string          `json:"dependencies"`
	Extra               string          `json:"extra"`
	FeatureFlags        map[string]bool `json:"feature_flags"`
	FlowdockAPIToken    interface{}     `json:"flowdock_api_token"`
	Followed            bool            `json:"followed"`
	HasUsableKey        bool            `json:"has_usable_key"`
	HerokuDeployUser    interface{}     `json:"heroku_deploy_user"`
	HipChatAPIToken     interface{}     `json:"hipchat_api_token"`
	HipChatNotify       interface{}     `json:"hipchat_notify"`
	HipChatNotifyPrefs  interface{}     `json:"hipchat_notify_prefs"`
	HipChatRoom         interface{}     `json:"hipchat_room"`
	IRCChannel          interface{}     `json:"irc_channel"`
	IRCKeyword          interface{}     `json:"irc_keyword"`
	IRCNotifyPrefs      interface{}     `json:"irc_notify_prefs"`
	IRCPassword         interface{}     `json:"irc_password"`
	IRCServer           interface{}     `json:"irc_server"`
	IRCUsername         interface{}     `json:"irc_username"`
	Parallel            int             `json:"parallel"`
	Reponame            string          `json:"reponame"`
	Scopes              []string        `json:"scopes"`
	Setup               string          `json:"setup"`
	SlackAPIToken       interface{}     `json:"slack_api_token"`
	SlackChannel        interface{}     `json:"slack_channel"`
	SlackNotifyPrefs    interface{}     `json:"slack_notify_prefs"`
	SlackSubdomain      interface{}     `json:"slack_subdomain"`
	SlackWebhookURL     string          `json:"slack_webhook_url"`
	SSHKeys             []interface{}   `json:"ssh_keys"`
	Test                string          `json:"test"`
	Username            string          `json:"username"`
	VCSURL              string          `json:"vcs_url"`
}

func (c *client) Projects() ([]Project, error) {
	url := c.endpoint("/projects")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]Project, 0), err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return make([]Project, 0), err
	}

	var p []Project
	err = json.NewDecoder(response.Body).Decode(&p)
	if err != nil {
		return make([]Project, 0), err
	}

	return p, nil
}

// Summary of a build.
type BuildSummary struct {
	CommitDetails []struct {
		AuthorDate     string `json:"author_date"`
		AuthorEmail    string `json:"author_email"`
		AuthorLogin    string `json:"author_login"`
		AuthorName     string `json:"author_name"`
		Body           string `json:"body"`
		Branch         string `json:"branch"`
		Commit         string `json:"commit"`
		CommitURL      string `json:"commit_url"`
		CommitterDate  string `json:"committer_date"`
		CommitterEmail string `json:"committer_email"`
		CommitterLogin string `json:"committer_login"`
		CommitterName  string `json:"committer_name"`
		Subject        string `json:"subject"`
	} `json:"all_commit_details"`
	AuthorDate      string      `json:"author_date"`
	AuthorEmail     string      `json:"author_email"`
	AuthorName      string      `json:"author_name"`
	Body            string      `json:"body"`
	Branch          string      `json:"branch"`
	BuildNum        int         `json:"build_num"`
	BuildParameters interface{} `json:"build_parameters"`
	BuildTimeMillis int         `json:"build_time_millis"`
	BuildURL        string      `json:"build_url"`
	Canceled        bool        `json:"canceled"`
	Canceler        interface{} `json:"canceler"`
	CircleYml       struct {
		String string `json:"string"`
	} `json:"circle_yml"`
	CommitterDate      string        `json:"committer_date"`
	CommitterEmail     string        `json:"committer_email"`
	CommitterName      string        `json:"committer_name"`
	Compare            string        `json:"compare"`
	DontBuild          interface{}   `json:"dont_build"`
	Failed             interface{}   `json:"failed"`
	FeatureFlags       struct{}      `json:"feature_flags"`
	HasArtifacts       bool          `json:"has_artifacts"`
	InfrastructureFail bool          `json:"infrastructure_fail"`
	IsFirstGreenBuild  bool          `json:"is_first_green_build"`
	JobName            interface{}   `json:"job_name"`
	Lifecycle          string        `json:"lifecycle"`
	Messages           []interface{} `json:"messages"`
	Node               []struct {
		ImageID      string      `json:"image_id"`
		Port         int         `json:"port"`
		PublicIPAddr string      `json:"public_ip_addr"`
		SSHEnabled   interface{} `json:"ssh_enabled"`
		Username     string      `json:"username"`
	} `json:"node"`
	Oss      bool   `json:"oss"`
	Outcome  string `json:"outcome"`
	Parallel int    `json:"parallel"`
	Previous struct {
		BuildNum        int    `json:"build_num"`
		BuildTimeMillis int    `json:"build_time_millis"`
		Status          string `json:"status"`
	} `json:"previous"`
	PreviousSuccessfulBuild struct {
		BuildNum        int    `json:"build_num"`
		BuildTimeMillis int    `json:"build_time_millis"`
		Status          string `json:"status"`
	} `json:"previous_successful_build"`
	QueuedAt      string        `json:"queued_at"`
	Reponame      string        `json:"reponame"`
	Retries       interface{}   `json:"retries"`
	RetryOf       int           `json:"retry_of"`
	SSHEnabled    interface{}   `json:"ssh_enabled"`
	SSHUsers      []interface{} `json:"ssh_users"`
	StartTime     string        `json:"start_time"`
	Status        string        `json:"status"`
	StopTime      string        `json:"stop_time"`
	Subject       string        `json:"subject"`
	Timedout      bool          `json:"timedout"`
	UsageQueuedAt string        `json:"usage_queued_at"`
	User          struct {
		Email  string `json:"email"`
		IsUser bool   `json:"is_user"`
		Login  string `json:"login"`
		Name   string `json:"name"`
	} `json:"user"`
	Username    string `json:"username"`
	VCSRevision string `json:"vcs_revision"`
	VCSURL      string `json:"vcs_url"`
	Why         string `json:"why"`
}

func (c *client) RecentBuilds() ([]BuildSummary, error) {
	url := c.endpoint("/recent-builds")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	var b []BuildSummary
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	return b, nil
}

func (c *client) RecentBuildsForProject(username, project string) ([]BuildSummary, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s", username, project))

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	var b []BuildSummary
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	return b, nil
}

type RecentBuildsOptions struct {
	Limit  *int
	Offset *int
	Filter *string
}

func (c *client) RecentBuildsForProjectBranch(username, project, branch string, options RecentBuildsOptions) ([]BuildSummary, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/tree/%s", username, project, branch))
	if options.Limit != nil {
		url = fmt.Sprintf("%s&limit=%d", url, *options.Limit)
	}
	if options.Offset != nil {
		url = fmt.Sprintf("%s&offset=%d", url, *options.Offset)
	}
	if options.Filter != nil {
		url = fmt.Sprintf("%s&filter=%s", url, *options.Filter)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	var b []BuildSummary
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return make([]BuildSummary, 0), err
	}

	return b, nil
}

// Detailed summary of a build.
type DetailedBuildSummary struct {
	BuildSummary
	Owners          []string      `json:"owners"`
	PullRequestUrls []interface{} `json:"pull_request_urls"`
	Steps           []struct {
		Actions []struct {
			BashCommand        interface{}   `json:"bash_command"`
			Canceled           interface{}   `json:"canceled"`
			Command            string        `json:"command"`
			Continue           interface{}   `json:"continue"`
			EndTime            string        `json:"end_time"`
			ExitCode           interface{}   `json:"exit_code"`
			Failed             interface{}   `json:"failed"`
			HasOutput          bool          `json:"has_output"`
			Index              int           `json:"index"`
			InfrastructureFail interface{}   `json:"infrastructure_fail"`
			Messages           []interface{} `json:"messages"`
			Name               string        `json:"name"`
			Parallel           bool          `json:"parallel"`
			RunTimeMillis      int           `json:"run_time_millis"`
			StartTime          string        `json:"start_time"`
			Status             string        `json:"status"`
			Step               int           `json:"step"`
			Timedout           interface{}   `json:"timedout"`
			Truncated          bool          `json:"truncated"`
			Type               string        `json:"type"`
		} `json:"actions"`
		Name string `json:"name"`
	} `json:"steps"`
}

func (c *client) BuildSummary(username, project string, num int) (DetailedBuildSummary, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/%d", username, project, num))

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DetailedBuildSummary{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return DetailedBuildSummary{}, err
	}

	var b DetailedBuildSummary
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return DetailedBuildSummary{}, err
	}

	return b, nil
}

// Artifact created by a build.
type Artifact struct {
	NodeIndex  int    `json:"node_index"`
	Path       string `json:"path"`
	PrettyPath string `json:"pretty_path"`
	URL        string `json:"url"`
}

func (c *client) Artifacts(username, project string, num int) ([]Artifact, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/%d/artifacts", username, project, num))

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]Artifact, 0), err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return make([]Artifact, 0), err
	}

	var a []Artifact
	err = json.NewDecoder(response.Body).Decode(&a)
	if err != nil {
		return make([]Artifact, 0), err
	}

	return a, nil
}

// Information about a build.
type Build struct {
	Body            string      `json:"body"`
	Branch          string      `json:"branch"`
	BuildNum        int         `json:"build_num"`
	BuildTimeMillis int         `json:"build_time_millis"`
	BuildURL        string      `json:"build_url"`
	CommitterEmail  string      `json:"committer_email"`
	CommitterName   string      `json:"committer_name"`
	DontBuild       string      `json:"dont_build"`
	Lifecycle       string      `json:"lifecycle"`
	Outcome         interface{} `json:"outcome"`
	Previous        struct {
		BuildNum int    `json:"build_num"`
		Status   string `json:"status"`
	} `json:"previous"`
	QueuedAt    string `json:"queued_at"`
	Reponame    string `json:"reponame"`
	RetryOf     int    `json:"retry_of"`
	StartTime   string `json:"start_time"`
	Status      string `json:"status"`
	StopTime    string `json:"stop_time"`
	Subject     string `json:"subject"`
	Username    string `json:"username"`
	VCSRevision string `json:"vcs_revision"`
	VCSURL      string `json:"vcs_url"`
	Why         string `json:"why"`
}

func (c *client) Retry(username, project string, num int) (Build, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/%d/retry", username, project, num))

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Build{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return Build{}, err
	}

	var b Build
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return Build{}, err
	}

	return b, nil
}

// Information about a build.
func (c *client) Cancel(username, project string, num int) (Build, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/%d/cancel", username, project, num))

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Build{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return Build{}, err
	}

	var b Build
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return Build{}, err
	}

	return b, nil
}

func (c *client) Build(username, project, branch string) (Build, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/tree/%s", username, project, branch))

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Build{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return Build{}, err
	}

	var b Build
	err = json.NewDecoder(response.Body).Decode(&b)
	if err != nil {
		return Build{}, err
	}

	return b, nil
}

// Response type indicating the status of clearing the cache.
type ClearCacheResponse struct {
	Status string `json:"status"`
}

func (c *client) ClearCache(username, project string) (ClearCacheResponse, error) {
	url := c.endpoint(fmt.Sprintf("/project/%s/%s/build-cache", username, project))

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return ClearCacheResponse{}, err
	}

	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return ClearCacheResponse{}, err
	}

	var res ClearCacheResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return ClearCacheResponse{}, err
	}

	return res, nil
}

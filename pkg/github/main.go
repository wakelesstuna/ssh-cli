package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var GITHUB_BASE_URL = "https://api.github.com"

func GetAllReposForUser(username string) ([]Repo, error) {

	token := os.Getenv("GITHUB_TOKEN")

	var allRepos []Repo
	page := 1
	for {
		url := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=100", GITHUB_BASE_URL, username, page)
		req, _ := http.NewRequest("GET", url, nil)
		// Error handling omitted for brevity

		req.Header.Add("Accept", "application/vnd.github+json")
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("Error fetching repositories:", err)
		}

		defer resp.Body.Close()
		fmt.Printf("StatusCode: %d\n", resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Error reading response body:", err)
		}

		// Convert the byte slice to a string and log it.
		bodyString := string(bodyBytes)
		// log.Println("Raw JSON response:", bodyString)

		var repos []Repo

		err = json.Unmarshal([]byte(bodyString), &repos)
		if err != nil {
			fmt.Println("Error:", err)
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < 100 {
			break // This was the last page
		}
		page++ // Prepare to fetch the next page
	}

	return allRepos, nil
}

type PostBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	IsTemplate  bool   `json:"is_template"`
}

type CreateRepoResponse struct {
	Name   string `json:"name"`
	SshUrl string `json:"ssh_url"`
	Url    string `json:"clone_url"`
}

func CreateRepository(name string, description string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	url := fmt.Sprintf("%s/user/repos", GITHUB_BASE_URL)

	data := PostBody{
		Name:        name,
		Description: description,
		Private:     false,
		IsTemplate:  false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("s: %d\n", resp.StatusCode)
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("could not create repo:statuscode: %d %s", resp.StatusCode, string(respBody))
	}

	var response CreateRepoResponse
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return "", err
	}

	return response.Url, nil
}

type Owner struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Permissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type SecurityAndAnalysis struct {
	AdvancedSecurity struct {
		Status string `json:"status"`
	} `json:"advanced_security"`
	SecretScanning struct {
		Status string `json:"status"`
	} `json:"secret_scanning"`
	SecretScanningPushProtection struct {
		Status string `json:"status"`
	} `json:"secret_scanning_push_protection"`
}

type Repo struct {
	ID                  int                 `json:"id"`
	NodeID              string              `json:"node_id"`
	Name                string              `json:"name"`
	FullName            string              `json:"full_name"`
	Owner               Owner               `json:"owner"`
	Private             bool                `json:"private"`
	HTMLURL             string              `json:"html_url"`
	Description         string              `json:"description"`
	Fork                bool                `json:"fork"`
	URL                 string              `json:"url"`
	ArchiveURL          string              `json:"archive_url"`
	AssigneesURL        string              `json:"assignees_url"`
	BlobsURL            string              `json:"blobs_url"`
	BranchesURL         string              `json:"branches_url"`
	CollaboratorsURL    string              `json:"collaborators_url"`
	CommentsURL         string              `json:"comments_url"`
	CommitsURL          string              `json:"commits_url"`
	CompareURL          string              `json:"compare_url"`
	ContentsURL         string              `json:"contents_url"`
	ContributorsURL     string              `json:"contributors_url"`
	DeploymentsURL      string              `json:"deployments_url"`
	DownloadsURL        string              `json:"downloads_url"`
	EventsURL           string              `json:"events_url"`
	ForksURL            string              `json:"forks_url"`
	GitCommitsURL       string              `json:"git_commits_url"`
	GitRefsURL          string              `json:"git_refs_url"`
	GitTagsURL          string              `json:"git_tags_url"`
	GitURL              string              `json:"git_url"`
	IssueCommentURL     string              `json:"issue_comment_url"`
	IssueEventsURL      string              `json:"issue_events_url"`
	IssuesURL           string              `json:"issues_url"`
	KeysURL             string              `json:"keys_url"`
	LabelsURL           string              `json:"labels_url"`
	LanguagesURL        string              `json:"languages_url"`
	MergesURL           string              `json:"merges_url"`
	MilestonesURL       string              `json:"milestones_url"`
	NotificationsURL    string              `json:"notifications_url"`
	PullsURL            string              `json:"pulls_url"`
	ReleasesURL         string              `json:"releases_url"`
	SSHURL              string              `json:"ssh_url"`
	StargazersURL       string              `json:"stargazers_url"`
	StatusesURL         string              `json:"statuses_url"`
	SubscribersURL      string              `json:"subscribers_url"`
	SubscriptionURL     string              `json:"subscription_url"`
	TagsURL             string              `json:"tags_url"`
	TeamsURL            string              `json:"teams_url"`
	TreesURL            string              `json:"trees_url"`
	CloneURL            string              `json:"clone_url"`
	MirrorURL           string              `json:"mirror_url"`
	HooksURL            string              `json:"hooks_url"`
	SVNURL              string              `json:"svn_url"`
	Homepage            string              `json:"homepage"`
	Language            interface{}         `json:"language"`
	ForksCount          int                 `json:"forks_count"`
	StargazersCount     int                 `json:"stargazers_count"`
	WatchersCount       int                 `json:"watchers_count"`
	Size                int                 `json:"size"`
	DefaultBranch       string              `json:"default_branch"`
	OpenIssuesCount     int                 `json:"open_issues_count"`
	IsTemplate          bool                `json:"is_template"`
	Topics              []string            `json:"topics"`
	HasIssues           bool                `json:"has_issues"`
	HasProjects         bool                `json:"has_projects"`
	HasWiki             bool                `json:"has_wiki"`
	HasPages            bool                `json:"has_pages"`
	HasDownloads        bool                `json:"has_downloads"`
	HasDiscussions      bool                `json:"has_discussions"`
	Archived            bool                `json:"archived"`
	Disabled            bool                `json:"disabled"`
	Visibility          string              `json:"visibility"`
	PushedAt            time.Time           `json:"pushed_at"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	Permissions         Permissions         `json:"permissions"`
	SecurityAndAnalysis SecurityAndAnalysis `json:"security_and_analysis"`
}

/**
 * @Time: 2019-06-29 09:31
 * @Author: solacowa@gmail.com
 * @File: githooks
 * @Software: GoLand
 */

package public

import "time"

type gitlabHook struct {
	Name        string
	Namespace   string
	KeyWord     string
	Branch      string
	Token       string
	After       string `json:"after"`
	Before      string `json:"before"`
	CheckoutSha string `json:"checkout_sha"`
	Commits     []struct {
		Added  []string `json:"added"`
		Author struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"author"`
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Modified  []string  `json:"modified"`
		Removed   []string  `json:"removed"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
	} `json:"commits"`
	EventName  string      `json:"event_name"`
	Message    interface{} `json:"message"`
	ObjectKind string      `json:"object_kind"`
	Project    struct {
		AvatarURL         string `json:"avatar_url"`
		DefaultBranch     string `json:"default_branch"`
		Description       string `json:"description"`
		GitHTTPURL        string `json:"git_http_url"`
		GitSSHURL         string `json:"git_ssh_url"`
		Homepage          string `json:"homepage"`
		HTTPURL           string `json:"http_url"`
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		PathWithNamespace string `json:"path_with_namespace"`
		SSHURL            string `json:"ssh_url"`
		URL               string `json:"url"`
		VisibilityLevel   int    `json:"visibility_level"`
		WebURL            string `json:"web_url"`
	} `json:"project"`
	ProjectID  int    `json:"project_id"`
	Ref        string `json:"ref"`
	Repository struct {
		Description     string `json:"description"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		Homepage        string `json:"homepage"`
		Name            string `json:"name"`
		URL             string `json:"url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
	TotalCommitsCount int    `json:"total_commits_count"`
	UserAvatar        string `json:"user_avatar"`
	UserEmail         string `json:"user_email"`
	UserID            int    `json:"user_id"`
	UserName          string `json:"user_name"`
}

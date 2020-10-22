/**
 * @Time : 2019-07-08 18:16
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package git_repo

import (
	"context"
	"github.com/google/go-github/v26/github"
	"github.com/icowan/config"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`

	// The following fields are only populated by Webhook events.
	Login *string `json:"username,omitempty"` // Renamed for go-github consistency.
}

type CommitStats struct {
	Additions *int `json:"additions,omitempty"`
	Deletions *int `json:"deletions,omitempty"`
	Total     *int `json:"total,omitempty"`
}

type Commit struct {
	SHA       *string       `json:"sha,omitempty"`
	Author    *CommitAuthor `json:"author,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Message   *string       `json:"message,omitempty"`
	//Tree         *Tree                  `json:"tree,omitempty"`
	Parents []Commit     `json:"parents,omitempty"`
	Stats   *CommitStats `json:"stats,omitempty"`
	HTMLURL *string      `json:"html_url,omitempty"`
	URL     *string      `json:"url,omitempty"`
	//Verification *SignatureVerification `json:"verification,omitempty"`
	NodeID *string `json:"node_id,omitempty"`

	// CommentCount is the number of GitHub comments on the commit. This
	// is only populated for requests that fetch GitHub data like
	// Pulls.ListCommits, Repositories.ListCommits, etc.
	CommentCount *int `json:"comment_count,omitempty"`

	// SigningKey denotes a key to sign the commit with. If not nil this key will
	// be used to sign the commit. The private key must be present and already
	// decrypted. Ignored if Verification.Signature is defined.
	SigningKey *openpgp.Entity `json:"-"`
}

type TagResponse struct {
	Name   string  `json:"name"`
	Commit *Commit `json:"commit"`
}

type BranchResponse struct {
	Name      string  `json:"name,omitempty"`
	Commit    *Commit `json:"commit,omitempty"`
	Protected *bool   `json:"protected,omitempty"`
}

type Repo interface {
	ListTags(ctx context.Context, owner, repo string, opts *ListOptions) (res []*TagResponse, err error)
	ListBranch(ctx context.Context, owner, repo string, opts *ListOptions) (res []*BranchResponse, err error)
	GetFile(ctx context.Context, owner, repo, branch, fileName string) (string, error)
}

func NewClient(cf *config.Config) Repo {
	var transport http.RoundTripper
	if cf.GetString("server", "http_proxy") != "" {
		dialer := &net.Dialer{
			Timeout:   time.Duration(30 * time.Second),
			KeepAlive: time.Duration(30 * time.Second),
		}
		transport = &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(cf.GetString("server", "http_proxy"))
			},
			DialContext: dialer.DialContext,
		}
	} else {
		transport = http.DefaultTransport
	}

	gitType := cf.GetString("git", "git_type")
	token := cf.GetString("git", "token")
	gitAddr := cf.GetString("git", "git_addr")

	if gitType == "github" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		)
		tc := oauth2.NewClient(ctx, ts)
		tc.Transport = transport
		client := github.NewClient(tc)
		return NewGithub(client)
	} else if gitType == "gitlab" {
		client := gitlab.NewClient(&http.Client{
			Transport: transport,
		}, token)
		_ = client.SetBaseURL(gitAddr)

		return NewGitlab(client)
	}
	return nil
}

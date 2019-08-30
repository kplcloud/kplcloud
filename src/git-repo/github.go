/**
 * @Time : 2019-07-08 18:59
 * @Author : solacowa@gmail.com
 * @File : github
 * @Software: GoLand
 */

package git_repo

import (
	"context"
	"encoding/base64"
	"github.com/google/go-github/v26/github"
)

type GithubRepository interface {
}

type githubRepository struct {
	client *github.Client
}

func NewGithub(client *github.Client) Repo {
	return &githubRepository{
		client: client,
	}
}

func (c *githubRepository) ListTags(ctx context.Context, owner, repo string, opts *ListOptions) (res []*TagResponse, err error) {
	var page, perPage int
	if opts != nil {
		page = opts.Page
		perPage = opts.PerPage
	}

	tags, _, err := c.client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	})

	if err != nil {
		return
	}

	for _, tag := range tags {
		res = append(res, &TagResponse{
			Name: tag.GetName(),
			Commit: &Commit{
				SHA: tag.GetCommit().SHA,
			},
		})
	}

	return
}

func (c *githubRepository) ListBranch(ctx context.Context, owner, repo string, opts *ListOptions) (res []*BranchResponse, err error) {
	var page, perPage int
	if opts != nil {
		page = opts.Page
		perPage = opts.PerPage
	}

	branches, _, err := c.client.Repositories.ListBranches(ctx, owner, repo, &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	})

	if err != nil {
		return
	}

	for _, branch := range branches {
		res = append(res, &BranchResponse{
			Name: branch.GetName(),
			Commit: &Commit{
				SHA: branch.Commit.SHA,
			},
		})
	}

	return
}

func (c *githubRepository) GetFile(ctx context.Context, owner, repo, branch, fileName string) (string, error) {
	fileContent, _, _, err := c.client.Repositories.GetContents(ctx, owner, repo, fileName, &github.RepositoryContentGetOptions{
		Ref: branch,
	})
	if err != nil {
		return "", err
	}
	var body []byte
	if *fileContent.Encoding == "base64" {
		if body, err = base64.StdEncoding.DecodeString(*fileContent.Content); err != nil {
			return "", err
		}
	} else {
		body = []byte(*fileContent.Content)
	}

	return string(body), nil
}

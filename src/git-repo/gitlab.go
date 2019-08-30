/**
 * @Time : 2019-07-08 18:59
 * @Author : solacowa@gmail.com
 * @File : gitlab
 * @Software: GoLand
 */

package git_repo

import (
	"context"
	"encoding/base64"
	"github.com/xanzy/go-gitlab"
)

type gitlabRepository struct {
	client *gitlab.Client
}

func NewGitlab(client *gitlab.Client) Repo {
	return &gitlabRepository{client: client}
}

func (c *gitlabRepository) ListTags(ctx context.Context, owner, repo string, opts *ListOptions) (res []*TagResponse, err error) {
	var page, perPage int
	if opts != nil {
		page = opts.Page
		perPage = opts.PerPage
	}

	tags, _, err := c.client.Tags.ListTags(owner+"/"+repo, &gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	})

	if err != nil {
		return
	}

	for _, tag := range tags {
		res = append(res, &TagResponse{
			Name: tag.Name,
			Commit: &Commit{
				SHA: &tag.Commit.ID,
			},
		})
	}

	return
}

func (c *gitlabRepository) ListBranch(ctx context.Context, owner, repo string, opts *ListOptions) (res []*BranchResponse, err error) {

	var page, perPage int
	if opts != nil {
		page = opts.Page
		perPage = opts.PerPage
	}

	// /api/v4/projects/solacowa/hello/repository/tags
	branches, _, err := c.client.Branches.ListBranches(owner+"/"+repo, &gitlab.ListBranchesOptions{
		Page:    page,
		PerPage: perPage,
	})

	if err != nil {
		return
	}

	for _, branch := range branches {
		res = append(res, &BranchResponse{
			Name: branch.Name,
			Commit: &Commit{
				SHA: &branch.Commit.ID,
			},
		})
	}

	return
}

func (c *gitlabRepository) GetFile(ctx context.Context, owner, repo, branch, fileName string) (string, error) {
	file, _, err := c.client.RepositoryFiles.GetFile(owner+"/"+repo, fileName, &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	})
	if err != nil {
		return "", err
	}

	var content []byte

	if file.Encoding == "base64" {
		content, err = base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return "", err
		}
	} else {
		content = []byte(file.Content)
	}

	return string(content), nil
}

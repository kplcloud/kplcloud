package git

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/git-repo"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

var (
	ErrProjectTemplateGet  = errors.New("项目基础模版获取错误")
	ErrProjectTagsGet      = errors.New("项目的Tags获取错误")
	ErrProjectTagsBranches = errors.New("项目的Branches获取错误")
)

type Service interface {
	// git tags 列表
	// desc: 需要区分是gitlab还是github
	Tags(ctx context.Context) (res []string, err error)

	// 获取git的分支
	Branches(ctx context.Context) (res []string, err error)

	// 获取项目的Dockerfile
	GetDockerfile(ctx context.Context, fileName string) (res string, err error)

	// 根据Git地址 获取tag列表
	TagsByGitPath(ctx context.Context, gitPath string) (res []string, err error)

	// 根据Git地址 获取branch列表
	BranchesByGitPath(ctx context.Context, gitPath string) (res []string, err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	gitClient  git_repo.Repo
	repository repository.Repository
}

func NewService(logger log.Logger, config *config.Config,
	gitClient git_repo.Repo,
	store repository.Repository) Service {
	return &service{logger,
		config,
		gitClient,
		store}
}

func (c *service) GetDockerfile(ctx context.Context, fileName string) (res string, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return "", ErrProjectTemplateGet
	}

	owner, repo := parseGitAddr(projectTemplate.FieldStruct.GitAddr)

	// todo 这里需要处理
	projectTemplate.FieldStruct.GitType = "master"

	return c.gitClient.GetFile(ctx, owner, repo, projectTemplate.FieldStruct.GitType, fileName)
}

func (c *service) Tags(ctx context.Context) (res []string, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return nil, ErrProjectTemplateGet
	}

	owner, repo := parseGitAddr(projectTemplate.FieldStruct.GitAddr)

	tags, err := c.gitClient.ListTags(ctx, owner, repo, nil)
	if err != nil {
		_ = level.Error(c.logger).Log("gitClient", "ListTags", "err", err.Error())
		return nil, errors.New(ErrProjectTagsGet.Error() + ":" + err.Error())
	}

	for _, tag := range tags {
		res = append(res, tag.Name)
	}

	return
}

func (c *service) Branches(ctx context.Context) (res []string, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	projectTemplate, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("projectTemplateRepository", "FindByProjectId", "err", err.Error())
		return nil, ErrProjectTemplateGet
	}

	owner, repo := parseGitAddr(projectTemplate.FieldStruct.GitAddr)

	branches, err := c.gitClient.ListBranch(ctx, owner, repo, nil)
	if err != nil {
		_ = level.Error(c.logger).Log("gitClient", "ListBranch", "err", err.Error())
		return nil, errors.New(ErrProjectTagsBranches.Error() + ":" + err.Error())
	}

	for _, branch := range branches {
		res = append(res, branch.Name)
	}

	return
}

func (c *service) TagsByGitPath(ctx context.Context, gitPath string) (res []string, err error) {
	_ = level.Debug(c.logger).Log("gitpaht", gitPath)
	owner, repo := parseGitAddr(gitPath)
	tags, err := c.gitClient.ListTags(ctx, owner, repo, nil)
	if err != nil {
		_ = level.Error(c.logger).Log("gitClient", "ListTags", "err", err.Error())
		return nil, errors.New(ErrProjectTagsGet.Error() + ":" + err.Error())
	}

	for _, tag := range tags {
		res = append(res, tag.Name)
	}

	return
}

func (c *service) BranchesByGitPath(ctx context.Context, gitPath string) (res []string, err error) {
	owner, repo := parseGitAddr(gitPath)
	branches, err := c.gitClient.ListBranch(ctx, owner, repo, nil)
	if err != nil {
		_ = level.Error(c.logger).Log("gitClient", "ListBranch", "err", err.Error())
		return nil, errors.New(ErrProjectTagsBranches.Error() + ":" + err.Error())
	}
	for _, branch := range branches {
		res = append(res, branch.Name)
	}
	return
}

func parseGitAddr(gitAddr string) (string, string) {
	gitAddr = strings.Replace(gitAddr, ".git", "", -1)
	addr := strings.Split(gitAddr, ":")
	names := strings.Split(addr[len(addr)-1], "/")
	owner := names[len(names)-2]
	repo := names[len(names)-1]

	return owner, repo
}

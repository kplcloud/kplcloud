/**
 * @Time : 6/8/21 3:45 PM
 * @Author : solacowa@gmail.com
 * @File : coding
 * @Software: GoLand
 */

package git_repo

import "context"

type codingRepository struct {
}

func (s *codingRepository) ListTags(ctx context.Context, owner, repo string, opts *ListOptions) (res []*TagResponse, err error) {
	panic("implement me")
}

func (s *codingRepository) ListBranch(ctx context.Context, owner, repo string, opts *ListOptions) (res []*BranchResponse, err error) {
	panic("implement me")
}

func (s *codingRepository) GetFile(ctx context.Context, owner, repo, branch, fileName string) (string, error) {
	panic("implement me")
}

func NewCoding() Repo {
	return &codingRepository{}
}

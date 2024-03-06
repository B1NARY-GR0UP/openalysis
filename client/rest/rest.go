package rest

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

var GlobalV3Client *github.Client

func Init() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.GlobalConfig.Backend.Token,
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	GlobalV3Client = github.NewClient(httpClient)
}

// GetContributorsByRepo return contributors, contributor count according to the provided repo
func GetContributorsByRepo(ctx context.Context, owner, name string) ([]*github.Contributor, int, error) {
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	res := make([]*github.Contributor, 0)
	for {
		contributors, resp, err := GlobalV3Client.Repositories.ListContributors(ctx, owner, name, opts)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, contributors...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return res, len(res), nil
}

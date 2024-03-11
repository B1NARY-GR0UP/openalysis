package rest

import (
	"context"
	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

var GlobalV3Client *github.Client

// Init go-github rest client
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
func GetContributorsByRepo(ctx context.Context, owner, name, repoNodeID string) ([]*model.Contributor, int, error) {
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	var cs []*github.Contributor
	for {
		contributors, resp, err := GlobalV3Client.Repositories.ListContributors(ctx, owner, name, opts)
		if err != nil {
			return nil, 0, err
		}
		cs = append(cs, contributors...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	var contributorData []*model.Contributor
	for _, c := range cs {
		user, err := graphql.QueryUserInfo(ctx, c.GetNodeID())
		if err != nil {
			return nil, 0, err
		}
		contributorData = append(contributorData, &model.Contributor{
			Login:         c.GetLogin(),
			NodeID:        c.GetNodeID(),
			Company:       user.Company,
			Location:      user.Location,
			AvatarURL:     c.GetAvatarURL(),
			RepoNodeID:    repoNodeID,
			Contributions: 0,
		})
	}
	return contributorData, len(contributorData), nil
}

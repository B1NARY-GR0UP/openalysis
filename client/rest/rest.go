// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"context"

	"github.com/B1NARY-GR0UP/openalysis/client/graphql"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"github.com/google/go-github/v60/github"
)

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
		user, err := graphql.QuerySingleUser(ctx, c.GetNodeID())
		if err != nil {
			return nil, 0, err
		}
		contributorData = append(contributorData, &model.Contributor{
			Login:         c.GetLogin(),
			NodeID:        c.GetNodeID(),
			Company:       user.Company,
			Location:      user.Location,
			AvatarURL:     c.GetAvatarURL(),
			RepoOwner:     owner,
			RepoName:      name,
			RepoNodeID:    repoNodeID,
			Contributions: c.GetContributions(),
		})
	}
	return contributorData, len(contributorData), nil
}

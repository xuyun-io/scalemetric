/*
Provider™ is a cloud native CaaS platform.
Copyright (C) 2019  Xuyun Authors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/xuyun-io/cloudprovider/types"
)

// RepoClient define repo client
type RepoClient struct {
	AccessKey       string
	SecretAccessKey string
	Region          string
}

// NewRepoClient return repo client
func NewRepoClient(region, accessKey, secretAccessKey string) types.RepoInterface {
	return &RepoClient{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}
}

// Repositories return aws ecr repos
func (rc *RepoClient) Repositories(ctx context.Context) ([]*types.Repository, error) {
	ecrClient, err := rc.eCRClient()
	if err != nil {
		return nil, err
	}
	out, err := ecrClient.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		return nil, err
	}
	if len(out.Repositories) <= 0 {
		return make([]*types.Repository, 0), nil
	}
	return toRepositories(out.Repositories), nil
}

var (
	maxResult int64 = 1000
)

// TagList return aws ecr images tags
func (rc *RepoClient) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	ecrClient, err := rc.eCRClient()
	if err != nil {
		return nil, err
	}
	output, err := ecrClient.DescribeImages(&ecr.DescribeImagesInput{
		MaxResults:     &maxResult,
		Filter:         &ecr.DescribeImagesFilter{TagStatus: aws.String("TAGGED")},
		RepositoryName: aws.String(repoName),
	})
	if err != nil {
		return nil, err
	}
	if len(output.ImageDetails) <= 0 {
		return make([]*types.ImageDetail, 0), nil
	}
	return toImageDetails(output.ImageDetails), nil
}

// ImageDescribe return image details
func (rc *RepoClient) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	ecrClient, err := rc.eCRClient()
	if err != nil {
		return nil, err
	}
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(imageName),
	}
	output, err := ecrClient.DescribeImages(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func toRepositories(repos []*ecr.Repository) []*types.Repository {
	typeRepos := make([]*types.Repository, 0)
	for i := range repos {
		typeRepos = append(typeRepos, &types.Repository{

			CreatedAt:      repos[i].CreatedAt,
			RegistryID:     repos[i].RegistryId,
			RepositoryName: repos[i].RepositoryName,
		})
	}
	return typeRepos
}

// ECRClient return aws ecr client
func (rc *RepoClient) eCRClient() (*ecr.ECR, error) {
	cfg := getAWSConfig(rc.Region, rc.AccessKey, rc.SecretAccessKey)
	ss, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return ecr.New(ss), nil
}

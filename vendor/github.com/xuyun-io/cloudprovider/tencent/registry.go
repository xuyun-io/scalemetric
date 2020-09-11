package tencent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/xuyun-io/cloudprovider/types"
)

// RepoClientset return aws repo clientset
func (tec *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	return tec
}

// Repositories return huawei swr repos
func (tec *Client) Repositories(ctx context.Context) ([]*types.Repository, error) {
	client, err := tec.NewTCRClient()
	if err != nil {
		return nil, err
	}
	request := tcr.NewDescribeRepositoriesRequest()
	request.Limit = common.Int64Ptr(1000)
	request.Offset = common.Int64Ptr(0)
	request.SortBy = common.StringPtr("-update_time")
	request.NamespaceName = common.StringPtr(tec.Auth.TCRAuth.Namespace)
	request.RegistryId = &tec.Auth.TCRAuth.RegistryID
	resp, err := client.DescribeRepositories(request)
	if err != nil {
		return nil, err
	}
	rs := make([]*types.Repository, 0)
	if resp.Response != nil && len(resp.Response.RepositoryList) > 0 {
		for i := range resp.Response.RepositoryList {
			t := *resp.Response.RepositoryList[i].Name
			if len(t) <= 0 {
				continue
			}
			// tag := strings.Join([]string{tec.Auth.TCRAuth.Host, t}, "/")
			r := types.Repository{RepositoryName: &t}
			rs = append(rs, &r)
		}
	}
	return rs, nil
}

// TagList return repo tag list
func (tec *Client) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	client, err := tec.NewTCRClient()
	if err != nil {
		return nil, err
	}
	repos := strings.Split(repoName, "/")
	l := len(repos)
	if l <= 1 {
		return nil, fmt.Errorf("repo is empty")
	}
	rpName := repos[l-1]
	if strings.Contains(rpName, ":") {
		repoName = strings.Split(rpName, ":")[0]
	}
	request := tcr.NewDescribeImagesRequest()
	request.Limit = common.Int64Ptr(1000)
	request.Offset = common.Int64Ptr(0)
	request.NamespaceName = common.StringPtr(tec.Auth.TCRAuth.Namespace)
	request.RegistryId = &tec.Auth.TCRAuth.RegistryID
	request.RepositoryName = &rpName
	resp, err := client.DescribeImages(request)
	if err != nil {
		return nil, err
	}
	imageDetails := make([]*types.ImageDetail, 0)
	if resp.Response != nil && len(resp.Response.ImageInfoList) > 0 {
		for i := range resp.Response.ImageInfoList {
			tag := resp.Response.ImageInfoList[i].ImageVersion
			imageDetails = append(imageDetails, &types.ImageDetail{
				ImageDigest:      resp.Response.ImageInfoList[i].Digest,
				ImageSizeInBytes: resp.Response.ImageInfoList[i].Size,
				ImageTags:        []*string{tag},
			})
		}
	}
	return imageDetails, nil
}

// ImageDescribe return image describe
func (tec *Client) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	return nil, errors.New("tcr image warehouse does not support external access")
}

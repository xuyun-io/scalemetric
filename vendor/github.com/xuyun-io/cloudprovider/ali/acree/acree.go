package acree

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/xuyun-io/cloudprovider/types"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	ListRepository = "ListRepository"
	ListRepoTag    = "ListRepoTag"
	GetRepository  = "GetRepository"
)

var (
	errRepor = errors.New("alibaba cloud image warehouse does not support external access")
)

type Client struct {
	RegionID        string
	AccessKey       string
	SecretAccessKey string
	RepoEEParam     RepoEEParams
}

func NewClient(regionID, accessKey, secretAccessKey string, params RepoEEParams) *Client {
	return &Client{
		RegionID:        regionID,
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		RepoEEParam:     params,
	}
}
func (ali *Client) ACREEClient() (*sdk.Client, error) {
	return sdk.NewClientWithAccessKey(ali.RegionID, ali.AccessKey, ali.SecretAccessKey)
}

// Repositories return repository list
func (c *Client) Repositories(ctx context.Context) ([]*types.Repository, error) {
	return c.listRepository()
}

// TagList return repo tag list
func (c *Client) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	return c.tagList(ctx, repoName)
}

// ImageDescribe return image describe
func (c *Client) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	return nil, errRepor
}

func (c *Client) tagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	rb, err := c.GetRepo(repoName)
	if err != nil {
		return nil, err
	}
	if len(rb.RepoID) <= 0 {
		return nil, fmt.Errorf("repo id is empty")
	}
	reqQuery := c.NewTagListReqQuery(rb.RepoID)
	response, err := c.reqACR(reqQuery)
	if err != nil {
		return nil, err
	}
	imList := &ImageTagList{}
	if err := json.Unmarshal(response.GetHttpContentBytes(), imList); err != nil {
		return nil, fmt.Errorf("acr body get failed, %v", err)
	}
	if !imList.IsSuccess {
		return nil, fmt.Errorf("acr requset failed, %v", imList.Message)
	}

	return imList.toImageDetail()

}

func (c *Client) reqACR(reqQuery *requests.CommonRequest) (*responses.CommonResponse, error) {
	acr, err := c.ACREEClient()
	if err != nil {
		return nil, fmt.Errorf("acr client get failed, %v", err)
	}
	response, err := acr.ProcessCommonRequest(reqQuery)
	if err != nil {
		return nil, fmt.Errorf("requset acr failed, %v", err)
	}
	if !response.IsSuccess() {
		return nil, fmt.Errorf("acr api http response code is incorrect(required: 200-300)")
	}
	return response, nil
}

func (c *Client) GetRepo(repoName string) (*RepositoryDetail, error) {
	repos := strings.Split(repoName, "/")
	if len(repos) <= 1 {
		return nil, fmt.Errorf("repo is empty")
	}
	l := len(repos)
	rpName := repos[l-1]
	if strings.Contains(rpName, ":") {
		repoName = strings.Split(rpName, ":")[0]
	}
	repoNamespaceName := repos[l-2]
	requestQuery := c.NewRepositoryReqQuery(repoNamespaceName, rpName)
	response, err := c.reqACR(requestQuery)
	if err != nil {
		return nil, err
	}
	rb := &RepositoryDetail{}
	if err := json.Unmarshal(response.GetHttpContentBytes(), rb); err != nil {
		return nil, fmt.Errorf("acr body get failed, %v", err)
	}
	if !rb.IsSuccess {
		return nil, fmt.Errorf("acr requset failed, %v", rb.Message)
	}
	return rb, nil
}

func (c *Client) listRepository() ([]*types.Repository, error) {
	reqQuery := c.NewListRepositoryReqQuery()
	response, err := c.reqACR(reqQuery)
	if err != nil {
		return nil, err
	}
	rb := &RepositoriesBody{}
	if err := json.Unmarshal(response.GetHttpContentBytes(), rb); err != nil {
		return nil, fmt.Errorf("acr body get failed, %v", err)
	}
	if !rb.IsSuccess {
		return nil, fmt.Errorf("acr requset failed, %v", rb.Message)
	}
	return rb.ToRepositories(c.RepoEEParam.Host), nil

}

// NewRepositoryReqQuery return repository requeset query.
func (c *Client) NewRepositoryReqQuery(repoNamespaceName, repoName string) *requests.CommonRequest {
	reqQuery := c.newRepoCommon()
	reqQuery.QueryParams["RepoNamespaceName"] = repoNamespaceName
	reqQuery.QueryParams["RepoName"] = repoName
	reqQuery.ApiName = GetRepository
	reqQuery.QueryParams["Action"] = GetRepository
	return reqQuery

}

// NewTagListReqQuery return tag list requery.
func (c *Client) NewTagListReqQuery(repoID string) *requests.CommonRequest {
	reqQuery := c.newRepoCommon()
	reqQuery.ApiName = ListRepoTag
	reqQuery.QueryParams["Action"] = ListRepoTag
	reqQuery.QueryParams["PageSize"] = "10000"
	reqQuery.QueryParams["RepoId"] = repoID
	return reqQuery
}

func (c *Client) newRepoCommon() *requests.CommonRequest {
	reqQuery := requests.NewCommonRequest()
	reqQuery.Method = "POST"
	reqQuery.Scheme = c.RepoEEParam.Scheme
	reqQuery.Version = "2018-12-01"
	reqQuery.QueryParams["RegionId"] = c.RegionID
	reqQuery.QueryParams["InstanceId"] = c.RepoEEParam.InstanceId
	reqQuery.QueryParams["PageSize"] = "10000"
	reqQuery.QueryParams["RepoName"] = c.RepoEEParam.RepoName
	reqQuery.QueryParams["RepoNamespaceName"] = c.RepoEEParam.RepoNamespaceName
	return reqQuery
}

// NewListRepositoryReqQuery return  RepositoryReqQuery.
func (c *Client) NewListRepositoryReqQuery() *requests.CommonRequest {
	reqQuery := c.newRepoCommon()
	reqQuery.ApiName = ListRepository
	reqQuery.QueryParams["Action"] = ListRepository
	reqQuery.QueryParams["RepoStatus"] = "ALL"
	return reqQuery
}

func (ims *ImageTagList) toImageDetail() ([]*types.ImageDetail, error) {
	imageDetails := make([]*types.ImageDetail, 0)
	if len(ims.Images) <= 0 {
		return imageDetails, nil
	}
	for i := range ims.Images {
		tagptr := &ims.Images[i].Tag
		tags := []*string{tagptr}
		t := int64ToTime(ims.Images[i].ImageUpdate)
		imageDetails = append(imageDetails, &types.ImageDetail{
			ImageDigest:      &ims.Images[i].Digest,
			ImagePushedAt:    &t,
			ImageSizeInBytes: &ims.Images[i].ImageSize,
			ImageTags:        tags,
			RegistryID:       &ims.Images[i].ImageID,
		})
	}
	return imageDetails, nil
}

// ToRepositories to repositories
func (rsb *RepositoriesBody) ToRepositories(host string) []*types.Repository {
	if len(rsb.Repositories) <= 0 {
		return make([]*types.Repository, 0)
	}
	rs := make([]*types.Repository, 0)
	for i := range rsb.Repositories {
		reposlice := []string{rsb.Repositories[i].RepoNamespaceName, rsb.Repositories[i].RepoName}
		repoName := strings.Join(reposlice, "/")
		rs = append(rs, &types.Repository{RepositoryName: &repoName, RegistryID: &rsb.Repositories[i].RepoID})
	}
	return rs
}

func int64ToTime(timeInt64 int64) time.Time {
	return time.Unix(timeInt64, 0)
}

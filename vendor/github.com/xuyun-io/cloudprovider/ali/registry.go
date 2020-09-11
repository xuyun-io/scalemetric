package ali

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/xuyun-io/cloudprovider/types"
)

var (
	errRepor = errors.New("alibaba cloud image warehouse does not support external access")
)

type RepoParams struct {
	// required, docker host
	Host string `json:"host,omitempty"`
	// required
	RepoNamespace string `json:"repoNamespace,omitempty"`
	// required default sbG93JsZA363489123
	// APIPassword       string `json:"apiPassword,omitempty"`
	PullImageUser     string `json:"pullImageUser,omitempty"`
	PullImagePassword string `json:"pullImagePassword,omitempty"`
}

// Repositories return repository list
func (c *Client) Repositories(ctx context.Context) ([]*types.Repository, error) {
	return c.repositories()
}

// TagList return repo tag list
func (c *Client) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	return c.tagList(repoName)
}

// ImageDescribe return image describe
func (c *Client) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	return nil, errRepor
}

func (c *Client) tagList(repoName string) ([]*types.ImageDetail, error) {
	repos := strings.Split(repoName, "/")
	l := len(repos)
	if l <= 1 {
		return nil, fmt.Errorf("repo is empty")
	}
	rpName := repos[l-1]
	if strings.Contains(rpName, ":") {
		repoName = strings.Split(rpName, ":")[0]
	}
	repoNamespaceName := repos[l-2]
	acrAPI, err := c.NewACRClient()
	if err != nil {
		return nil, err
	}
	r := cr.CreateGetRepoTagsRequest()
	r.PathParams["RepoNamespace"] = repoNamespaceName
	r.PathParams["RepoName"] = rpName
	r.PageSize = requests.NewInteger(100)
	resp, err := acrAPI.GetRepoTags(r)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("ali api get failed")
	}
	tb := &TagsBody{}
	if err := json.Unmarshal(resp.GetHttpContentBytes(), tb); err != nil {
		return nil, err
	}
	if len(tb.Message) > 0 {
		return nil, fmt.Errorf("ali api failed: %v-> %s", tb.Code, tb.Message)
	}
	return tb.toImageDetail()
}

type TagsBody struct {
	Data    TagData `json:"data"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
}

func (tb *TagsBody) toImageDetail() ([]*types.ImageDetail, error) {
	imageDetails := make([]*types.ImageDetail, 0)
	if len(tb.Data.Tags) <= 0 {
		return imageDetails, nil
	}
	for i := range tb.Data.Tags {
		tagptr := &tb.Data.Tags[i].Tag
		tags := []*string{tagptr}
		imageDetails = append(imageDetails, &types.ImageDetail{
			ImageDigest:      &tb.Data.Tags[i].Digest,
			ImageSizeInBytes: &tb.Data.Tags[i].ImageSize,
			ImageTags:        tags,
			RegistryID:       &tb.Data.Tags[i].ImageID,
		})
	}
	return imageDetails, nil
}

func int64ToTime(timeInt64 int64) time.Time {
	return time.Unix(timeInt64, 0)
}

type TagData struct {
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	Tags     []TagMessage `json:"tags"`
	PageSize int          `json:"pageSize"`
}

type TagMessage struct {
	Status      string `json:"status"`
	Digest      string `json:"digest"`
	ImageCreate int64  `json:"imageCreate"`
	ImageID     string `json:"imageId"`
	ImageUpdate int64  `json:"imageUpdate"`
	Tag         string `json:"tag"`
	ImageSize   int64  `json:"imageSize"`
}

var (
	content = `{
	"User": {
		"Password": "%s"
	}
}`
)

func (c *Client) getContent() []byte {
	return []byte(fmt.Sprintf(content, c.RepoParam.PullImagePassword))
}

func (c *Client) applyUserPassword() (bool, error) {
	acrAPI, err := c.NewACRClient()
	if err != nil {
		return false, err
	}
	contentByts := c.getContent()
	r := cr.CreateUpdateUserInfoRequest()
	r.SetContent(contentByts)
	resp, err := acrAPI.UpdateUserInfo(r)
	if err != nil {
		return false, err
	}
	if resp.IsSuccess() {
		return true, nil
	}
	rr := cr.CreateCreateUserInfoRequest()
	rr.SetContent(contentByts)
	cresp, err := acrAPI.CreateUserInfo(rr)
	if err != nil {
		return false, err
	}
	if cresp.IsSuccess() {
		return true, nil
	}
	return false, fmt.Errorf("apply user password failed")
}

func (c *Client) repositories() ([]*types.Repository, error) {
	ok, err := c.applyUserPassword()
	if err != nil {
		return nil, fmt.Errorf("ali user: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("apply user message failed")
	}
	return c.getRepoListByNamespace(c.RepoParam.Host, c.RepoParam.RepoNamespace)
}

func (c *Client) getRepoListByNamespace(host, repoNamespace string) ([]*types.Repository, error) {
	acrAPI, err := c.NewACRClient()
	if err != nil {
		return nil, fmt.Errorf("ali user info apply failed, %v", err)
	}
	r := cr.CreateGetRepoListByNamespaceRequest()
	r.RegionId = c.RegionID
	r.PathParams["RepoNamespace"] = repoNamespace
	r.Status = "ALL"
	r.PageSize = requests.NewInteger(100)
	resp, err := acrAPI.GetRepoListByNamespace(r)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("get repolist by namespace failed")
	}

	rb := &RepoBody{}
	if err := json.Unmarshal(resp.GetHttpContentBytes(), rb); err != nil {
		return nil, err
	}
	if len(rb.Message) > 0 {
		return nil, fmt.Errorf("ali api failed: %v-> %s", rb.Code, rb.Message)
	}
	return rb.ToToRepositories(host), nil
}

type RepoBody struct {
	Data    RepoData `json:"data"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
}

func (rb *RepoBody) ToToRepositories(host string) []*types.Repository {
	rs := make([]*types.Repository, 0)
	if len(rb.Data.Repos) <= 0 {
		return rs
	}
	for i := range rb.Data.Repos {
		reposlice := []string{rb.Data.Repos[i].RepoNamespace, rb.Data.Repos[i].RepoName}
		repoName := strings.Join(reposlice, "/")
		repoStr := strconv.Itoa(rb.Data.Repos[i].RepoID)
		rs = append(rs, &types.Repository{RepositoryName: &repoName, RegistryID: &repoStr})
	}
	return rs
}

type RepoData struct {
	Page     int    `json:"page"`
	Total    int    `json:"total"`
	Repos    []Repo `json:"repos"`
	PageSize int    `json:"pageSize"`
}

type Repo struct {
	Summary        string `json:"summary"`
	RepoStatus     string `json:"repoStatus"`
	RepoID         int    `json:"repoId"`
	RepoOriginType string `json:"repoOriginType"`
	RepoBuildType  string `json:"repoBuildType"`
	Logo           string `json:"logo"`
	Stars          int    `json:"stars"`
	RegionID       string `json:"regionId"`
	RepoDomainList struct {
		Internal string `json:"internal"`
		Public   string `json:"public"`
		Vpc      string `json:"vpc"`
	} `json:"repoDomainList"`
	RepoAuthorizeType string `json:"repoAuthorizeType"`
	Downloads         int    `json:"downloads"`
	GmtCreate         int64  `json:"gmtCreate"`
	RepoType          string `json:"repoType"`
	RepoNamespace     string `json:"repoNamespace"`
	RepoName          string `json:"repoName"`
	GmtModified       int64  `json:"gmtModified"`
}

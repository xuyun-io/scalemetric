package huawei

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/imroc/req"
	"github.com/xuyun-io/cloudprovider/kubeconfig"
	"github.com/xuyun-io/cloudprovider/types"
)

func generateTokenHeader(token string) req.Header {
	return req.Header{
		"X-Auth-Token": token,
		"Content-Type": "application/json",
	}
}

func (hw *Client) generateSWRRepos(region, namespace string) string {
	url := fmt.Sprintf("https://swr-api.%s.myhuaweicloud.com/v2/manage/repos?filter=namespace::%s|offset::0|limit::1000|order_column::updated_time|order_type::desc", region, namespace)
	return url
}

func (hw *Client) generateSWRTags(region, namespace, repository string) string {
	url := fmt.Sprintf("https://swr-api.%s.myhuaweicloud.com/v2/manage/namespaces/%s/repos/%s/tags?filter=offset::0|limit::1000|order_column::updated_at|order_type::desc", region, namespace, repository)
	return url

}

func (hw *Client) getCluster(clusterid string) (*types.DescribeCluster, error) {
	token, err := hw.getCCEToken()
	if err != nil {
		return nil, err
	}
	url := generateClusterURL(hw.Auth.Base.Region, hw.Auth.CCE.ProjectID, clusterid)
	header := generateTokenHeader(token)
	resp, err := req.Get(url, header)
	if err != nil {
		return nil, err
	}
	cluster := &Cluster{}
	if err := resp.ToJSON(cluster); err != nil {
		return nil, err
	}
	return &types.DescribeCluster{
		ClusterMeta: types.ClusterMeta{
			UseUnique: cluster.Metadata.UID,
			RegionID:  hw.Auth.Base.Region,
			Name:      cluster.Metadata.Name,
			State:     cluster.Status.Phase,
			ClusterID: cluster.Metadata.UID,
		},
	}, nil
}

func (hw *Client) getClusters() (*types.ClusterList, error) {
	token, err := hw.getCCEToken()
	if err != nil {
		return nil, err
	}
	url := generateClustersURL(hw.Auth.Base.Region, hw.Auth.CCE.ProjectID)
	header := generateTokenHeader(token)
	resp, err := req.Get(url, header)
	if err != nil {
		return nil, err
	}
	cList := &ClusterList{Items: make([]Cluster, 0)}
	if err := resp.ToJSON(cList); err != nil {
		return nil, err
	}
	clusterList := &types.ClusterList{Item: make([]types.DescribeCluster, 0)}
	for i := range cList.Items {
		clusterList.Item = append(clusterList.Item, types.DescribeCluster{
			ClusterMeta: types.ClusterMeta{
				UseUnique: cList.Items[i].Metadata.UID,
				RegionID:  hw.Auth.Base.Region,
				Name:      cList.Items[i].Metadata.Name,
				State:     cList.Items[i].Status.Phase,
				ClusterID: cList.Items[i].Metadata.UID,
			},
		})
	}
	return clusterList, nil
}

func (hw *Client) kubeconfig(cluster string) (types.Client, error) {
	cfg, err := hw.getkubeconfig(cluster)
	if err != nil {
		return nil, err
	}
	return kubeconfig.NewClientByRegion(hw.Auth.Base.Region, cfg)
}

func (hw *Client) getRepoListByNamespace(region, namespace string) ([]*types.Repository, error) {
	token, err := hw.getSWRToken()
	if err != nil {
		return nil, err
	}
	url := hw.generateSWRRepos(region, namespace)
	resp, err := req.Get(url, generateTokenHeader(token))
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("repositories get failed form huawei, http code is:%d", resp.Response().StatusCode)
	}
	rbs := make([]*RepoBody, 0)
	if err := resp.ToJSON(&rbs); err != nil {
		return nil, fmt.Errorf("repo list body get failed, %v", err)
	}
	if len(rbs) <= 0 {
		return nil, nil
	}
	host := repoHost(hw.Auth.SWR.ProjectName) + "/"
	rs := make([]*types.Repository, 0)
	for i := range rbs {
		repo := rbs[i].Path
		if strings.HasPrefix(repo, host) {
			repo = strings.TrimPrefix(repo, host)
		}
		rs = append(rs, &types.Repository{
			RepositoryName: &repo,
		})
	}
	return rs, nil
}

func (hw *Client) tags(ctx context.Context, region, repoName string) ([]*types.ImageDetail, error) {
	token, err := hw.getSWRToken()
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
	repoNamespaceName := repos[l-2]
	url := hw.generateSWRTags(region, repoNamespaceName, rpName)
	resp, err := req.Get(url, generateTokenHeader(token))
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tags get failed form huawei, http code is:%d", resp.Response().StatusCode)
	}
	tags := make([]*TagBody, 0)
	if err := resp.ToJSON(&tags); err != nil {
		return nil, fmt.Errorf("tag list body get failed, %v", err)
	}
	if len(tags) <= 0 {
		return nil, nil
	}
	imageDetails := make([]*types.ImageDetail, 0)
	for i := range tags {
		ts := []*string{&tags[i].Tag}
		imageDetails = append(imageDetails, &types.ImageDetail{
			ImageDigest:      &tags[i].Digest,
			ImageSizeInBytes: &tags[i].Size,
			ImageTags:        ts,
			RegistryID:       &tags[i].ImageID,
			ImagePushedAt:    &tags[i].Updated,
		})
	}
	return imageDetails, nil
}

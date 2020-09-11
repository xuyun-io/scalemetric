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
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/xuyun-io/cloudprovider/types"
)

// Client define aws provider
type Client struct {
	Region string
	Auth   Auth
}

var (
	regionList = RegionList()
)

// Auth define aws client auth
type Auth struct {
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
}

// GetConfig return aws client config
func (aws *Auth) GetConfig() types.Auth {
	return aws
}

// DataCheck check aws client auth info
func (aws *Auth) DataCheck() error {
	if len(aws.AccessKey) <= 0 {
		return fmt.Errorf("AccessKey cann't be empty")
	}
	if len(aws.SecretAccessKey) <= 0 {
		return fmt.Errorf("SecretAccessKey cann't be empty")
	}
	return nil
}

// Regions return regions for aws
func Regions() *types.RegionList {
	return regionList
}

// NewClient return aws provider
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderAWS {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("aws auth config get failed")
	}
	return newClient(providercfg.GetSpec().GetRegionID(), auth.AccessKey, auth.SecretAccessKey), nil

}

// K8SClientset return k8s clientset
func (a *Client) K8SClientset() types.KubernetesClientset {
	return a
}

// Clusterset return aws provider cluster set
func (a *Client) Clusterset() types.Clusterset {
	return a
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (a *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := a.RestConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("restconfig create error %v", err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, nil, fmt.Errorf("client create error %v", err)
	}
	return client, restConf, nil
}

// RepoClientset return aws repo clientset
func (a *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	// use default repo
	return a.NewRepoClient()
}

// EKS return aws eks client
func (a *Client) EKS() (*eks.EKS, error) {
	cfg := a.getAWSConfig(a.Region)
	ss, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("session create failed. %v", err)
	}
	return eks.New(ss), nil
}

// ECR return aws ecr client
func (a *Client) ECR() (*ecr.ECR, error) {
	cfg := a.getAWSConfig(a.Region)
	ss, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return ecr.New(ss), nil
}

// NewRepoClient return repo client
func (a *Client) NewRepoClient() *RepoClient {
	return &RepoClient{
		AccessKey:       a.Auth.AccessKey,
		SecretAccessKey: a.Auth.SecretAccessKey,
		Region:          a.Region,
	}
}

// Clusters return aws cluster on region
func (a *Client) Clusters() (*types.ClusterList, error) {
	eksClient, err := a.EKS()
	if err != nil {
		return nil, err
	}
	input := &eks.ListClustersInput{}
	clustersOutput, err := eksClient.ListClusters(input)
	if err != nil {
		return nil, fmt.Errorf("cluster list failed. %v", err)
	}
	if len(clustersOutput.Clusters) <= 0 {
		return nil, fmt.Errorf("cluster not found in region: %s", a.Region)
	}
	clusters := []string{}
	for i := range clustersOutput.Clusters {
		cluster := aws.StringValue(clustersOutput.Clusters[i])
		if cluster != "" {
			clusters = append(clusters, cluster)
		}
	}
	sort.Strings(clusters)
	clusterList := &types.ClusterList{Item: make([]types.DescribeCluster, 0)}
	for i := range clusters {
		clusterList.Item = append(clusterList.Item, types.DescribeCluster{ClusterMeta: types.ClusterMeta{UseUnique: clusters[i], Name: clusters[i]}})
	}
	return clusterList, nil
}

// ConfigClientset return kubentes clientset, rest.config and error
func (a *Client) ConfigClientset(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConfig, err := a.RestConfig(cluster)
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	return clientset, restConfig, nil
}

// Cluster reutn cluster describe
func (a *Client) Cluster(cluster string) (*types.DescribeCluster, error) {
	cfg := a.getAWSConfig(a.Region)
	ss, err := session.NewSessionWithOptions(session.Options{
		Config:                  *cfg,
		AssumeRoleTokenProvider: token.StdinStderrTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("new session failed. %v", err)
	}
	eksclient := eks.New(ss, cfg)
	input := &eks.DescribeClusterInput{Name: aws.String(cluster)}
	output, err := eksclient.DescribeCluster(input)
	if err != nil {
		return nil, err
	}
	describe := &types.DescribeCluster{
		ClusterMeta: types.ClusterMeta{
			Name:      *output.Cluster.Name,
			UseUnique: *output.Cluster.Name,
			State:     *output.Cluster.Status,
		},
	}
	return describe, nil
}

// RestConfig return kubernetes rest config
func (a *Client) RestConfig(cluster string) (*rest.Config, error) {
	cfg := a.getAWSConfig(a.Region)
	tg, err := token.NewGenerator(false, false)
	if err != nil {
		return nil, err
	}
	ss, err := session.NewSessionWithOptions(session.Options{
		Config:                  *cfg,
		AssumeRoleTokenProvider: token.StdinStderrTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	ts, err := tg.GetWithRoleForSession(cluster, "", ss)
	if err != nil {
		return nil, err
	}
	_, err = token.NewVerifier(cluster).Verify(ts.Token)
	if err != nil {
		return nil, err
	}
	eksclient := eks.New(ss, cfg)
	input := &eks.DescribeClusterInput{Name: aws.String(cluster)}
	output, err := eksclient.DescribeCluster(input)
	if err != nil {
		return nil, err
	}
	capem, err := base64.StdEncoding.DecodeString(*output.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}
	config := &rest.Config{
		Host: aws.StringValue(output.Cluster.Endpoint),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: capem,
		},
		BearerToken: ts.Token,
	}
	return config, nil
}

// Clientset return kubernetes clientset
func (a *Client) Clientset(cluter string) (kubernetes.Interface, error) {
	restConfig, err := a.RestConfig(cluter)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes dynamic kubernetes clientset
func (a *Client) DynamicClientset(cluter string) (dynamic.Interface, error) {
	restConfig, err := a.RestConfig(cluter)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

// ZoneList return all zones on aws
func ZoneList() types.ZoneList {
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()
	if len(partitions) <= 0 {
		return types.ZoneList{}
	}

	zoneList := types.ZoneList{Item: make([]types.Zone, 0)}
	for i := range partitions {
		rs := partitions[i].Regions()
		if len(rs) <= 0 {
			continue
		}
		for _, zone := range rs {
			zoneList.Item = append(zoneList.Item, types.Zone{ID: zone.ID(), Name: zone.Description()})
		}
	}
	return zoneList
}

// RegionList return all regions on aws
func RegionList() *types.RegionList {
	got := ZoneList()
	if len(got.Item) <= 0 {
		return &types.RegionList{}
	}
	rl := types.NewRegionList()
	for i := range got.Item {
		area := strings.Split(got.Item[i].Name, " ")
		r := types.Region{RegionID: got.Item[i].ID, Name: got.Item[i].Name}
		z := types.Zone{ID: got.Item[i].ID}
		rl.AddRegion(area[0], r)
		rl.AddZone(got.Item[i].ID, z)
	}
	return rl
}

// NewClient return aws provider
func newClient(region, accessKey, secretAccessKey string) *Client {
	return &Client{
		Region: region,
		Auth: Auth{
			AccessKey:       accessKey,
			SecretAccessKey: secretAccessKey,
		},
	}
}

// getAWSConfig return aws config
func (a *Client) getAWSConfig(region string) *aws.Config {
	return getAWSConfig(region, a.Auth.AccessKey, a.Auth.SecretAccessKey)
}
func getAWSConfig(region, accessKey, secretAccessKey string) *aws.Config {
	credential := credentials.NewStaticCredentials(accessKey, secretAccessKey, "")
	cfg := &aws.Config{Credentials: credential}
	cfg.WithRegion(region)
	return cfg

}
func toImageDetails(images []*ecr.ImageDetail) []*types.ImageDetail {
	imageDetails := make([]*types.ImageDetail, 0)
	for i := range images {
		imageDetails = append(imageDetails, &types.ImageDetail{
			ImageDigest:      images[i].ImageDigest,
			ImagePushedAt:    images[i].ImagePushedAt,
			ImageSizeInBytes: images[i].ImageSizeInBytes,
			ImageTags:        images[i].ImageTags,
			RegistryID:       images[i].RegistryId,
			RepositoryName:   images[i].RepositoryName,
		})
	}
	return imageDetails
}

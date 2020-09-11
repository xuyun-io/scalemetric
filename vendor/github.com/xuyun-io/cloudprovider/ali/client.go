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

package ali

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/xuyun-io/cloudprovider/ali/acree"
	"github.com/xuyun-io/cloudprovider/types"
	"github.com/xuyun-io/cloudprovider/utils"
)

// Client define ali cloud provider
type Client struct {
	RegionID  string
	Auth      Auth
	RepoParam RepoParams
}

// Auth define client auth config
type Auth struct {
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
}

const (
	// ClusterTypeKubernetes define cluster type is kubernetes
	ClusterTypeKubernetes = "Kubernetes"
)

var _ types.Client = &Client{}

// DataCheck check client auth data
func (ali *Auth) DataCheck() error {
	if len(ali.AccessKey) <= 0 {
		return fmt.Errorf("AccessKey cann't be empty")
	}
	if len(ali.SecretAccessKey) <= 0 {
		return fmt.Errorf("SecretAccessKey cann't be empty")
	}
	return nil
}

// GetConfig return client config
func (ali *Auth) GetConfig() types.Auth {
	return ali
}

// NewClient return ali cloud client
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderALiCloud {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("ali auth config get failed")
	}
	return newClient(providercfg.GetSpec().GetRegionID(), auth.AccessKey, auth.SecretAccessKey), nil
}

// NewRepoEEClient return acree api
func (ali *Client) NewRepoEEClient(params acree.RepoEEParams) types.RepoInterface {
	return NewRepoEEClient(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey, params)
}

// NewRepoClient return acr client.
func (ali *Client) NewRepoClient(params RepoParams) types.RepoInterface {
	return NewRepoClient(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey, params)
}

// NewRepoClient return repo client.
func NewRepoClient(regionID, accessKey, secretAccessKey string, params RepoParams) types.RepoInterface {
	client := newClient(regionID, accessKey, secretAccessKey)
	client.RepoParam = params
	return client
}

// NewACRClient return acr client.
func (ali *Client) NewACRClient() (*cr.Client, error) {
	return cr.NewClientWithAccessKey(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey)
}

// NewRepoEEClient return ali docker registry client.
func NewRepoEEClient(regionID, accessKey, secretAccessKey string, params acree.RepoEEParams) types.RepoInterface {
	return acree.NewClient(regionID, accessKey, secretAccessKey, params)
}

// K8SClientset return k8s clientset
func (ali *Client) K8SClientset() types.KubernetesClientset {
	return ali
}

// Clusterset return aws provider cluster set
func (ali *Client) Clusterset() types.Clusterset {
	return ali
}

// RepoClientset return aws repo clientset
func (ali *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	return ali
}

// Regions return regions of ali cloud
func Regions() *types.RegionList {
	return regions()
}

// Cluster return cluster details
func (ali *Client) Cluster(cluserID string) (*types.DescribeCluster, error) {
	client, err := ali.CSClient()
	if err != nil {
		return nil, err
	}
	r := cs.CreateDescribeClusterDetailRequest()
	r.ClusterId = cluserID
	resp, err := client.DescribeClusterDetail(r)
	if err = isJSONUnmarshalErrorCode(err); err != nil {
		return nil, err
	}
	cluster := &DescribeCluster{}
	if err = json.Unmarshal(resp.GetHttpContentBytes(), cluster); err != nil {
		return nil, err
	}
	cluster.ClusterMeta.UseUnique = cluster.ClusterID
	return &types.DescribeCluster{
		ClusterMeta: cluster.ClusterMeta,
	}, nil
}

// Clusters return kubernetes clusters meta data message in region
func (ali *Client) Clusters() (*types.ClusterList, error) {
	client, err := ali.CSClient()
	if err != nil {
		return nil, err
	}
	r := cs.CreateDescribeClustersRequest()
	r.ClusterType = ClusterTypeKubernetes
	resp, err := client.DescribeClusters(r)
	if err = isJSONUnmarshalErrorCode(err); err != nil {
		return nil, err
	}
	clusterList := &ClusterList{Item: make([]DescribeCluster, 0)}
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &clusterList.Item); err != nil {
		return nil, err
	}
	cl := &types.ClusterList{Item: make([]types.DescribeCluster, 0)}
	if len(clusterList.Item) <= 0 {
		return cl, nil
	}
	for i := range clusterList.Item {
		clusterList.Item[i].UseUnique = clusterList.Item[i].ClusterID
		cl.Item = append(cl.Item, types.DescribeCluster{ClusterMeta: clusterList.Item[i].ClusterMeta})
	}
	return cl, nil
}

// CSClient return container service client
func (ali *Client) CSClient() (*cs.Client, error) {
	return cs.NewClientWithAccessKey(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey)
}

// ACREEClient return acr client.
func (ali *Client) ACREEClient() (*sdk.Client, error) {
	return sdk.NewClientWithAccessKey(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey)
}

// ACRClient2 return acr client.
func (ali *Client) ACRClient2() (*cr.Client, error) {
	return cr.NewClientWithAccessKey(ali.RegionID, ali.Auth.AccessKey, ali.Auth.SecretAccessKey)
}

// KubernetesConfig return api config of clusterID
func (ali *Client) KubernetesConfig(clusterID string) (*api.Config, error) {
	client, err := ali.CSClient()
	if err != nil {
		return nil, err
	}
	r := cs.CreateDescribeClusterUserKubeconfigRequest()
	r.ClusterId = clusterID
	resp, err := client.DescribeClusterUserKubeconfig(r)
	if err != nil {
		return nil, err
	}
	ak8scfg := &k8sConfig{}
	byts := resp.GetHttpContentBytes()
	if err := json.Unmarshal(byts, ak8scfg); err != nil {
		return nil, err
	}
	if ak8scfg.Config == "" {
		return nil, fmt.Errorf("kubeconfig get error, config is empty")
	}
	return utils.K8sV1Config([]byte(ak8scfg.Config))
}

// RestConfig return kubernetes restconfig
func (ali *Client) RestConfig(clusterID string) (*rest.Config, error) {
	config, err := ali.KubernetesConfig(clusterID)
	if err != nil {
		return nil, err
	}
	return utils.RestConfig(config)
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (ali *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := ali.RestConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("restconfig create error %v", err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, nil, fmt.Errorf("client create error %v", err)
	}
	return client, restConf, nil
}

// Clientset return kubernetes clientset
func (ali *Client) Clientset(clusterID string) (kubernetes.Interface, error) {
	restConfig, err := ali.RestConfig(clusterID)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes dynamic kubernetes clientset
func (ali *Client) DynamicClientset(clusterID string) (dynamic.Interface, error) {
	restConfig, err := ali.RestConfig(clusterID)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

func isJSONUnmarshalErrorCode(err error) error {
	if err != nil {
		if clientErr, isClientErr := err.(*errors.ClientError); isClientErr {
			if clientErr.ErrorCode() != errors.JsonUnmarshalErrorCode {
				return err
			}
		}
	}
	return nil
}

func newClient(regionID, accessKey, secretAccessKey string) *Client {
	return &Client{
		RegionID: regionID,
		Auth: Auth{
			AccessKey:       accessKey,
			SecretAccessKey: secretAccessKey,
		},
	}
}

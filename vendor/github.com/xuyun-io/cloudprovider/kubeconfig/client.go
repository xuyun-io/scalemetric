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

package kubeconfig

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/xuyun-io/cloudprovider/registry"
	"github.com/xuyun-io/cloudprovider/types"
	"github.com/xuyun-io/cloudprovider/utils"
)

// Client define kubeconfig client
type Client struct {
	RegionID string
	Auth     Auth
}

// Auth define aws client auth
type Auth struct {
	KubeConfigString string `json:"kubeConfigString"`
}

// DataCheck check aws client auth info
func (kubeconfig *Auth) DataCheck() error {
	if len(kubeconfig.KubeConfigString) <= 0 {
		return fmt.Errorf("KubeConfig cann't be emtpy")
	}
	return nil
}

// GetConfig return aws client config
func (kubeconfig *Auth) GetConfig() types.Auth {
	return kubeconfig
}

// NewClient return kubeconfig client
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderKubeconfig {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("auth config get failed")
	}
	return newClient(providercfg.GetSpec().GetRegionID(), auth.KubeConfigString), nil
}

// NewClientByRegion return clientset by region and kubeconfig.
func NewClientByRegion(region, kubeconfig string) (types.Client, error) {
	return newClient(region, kubeconfig), nil
}

// K8SClientset return k8s clientset
func (kc *Client) K8SClientset() types.KubernetesClientset {
	return kc
}

// Clusterset return aws provider cluster set
func (kc *Client) Clusterset() types.Clusterset {
	return kc
}

// RepoClientset return aws repo clientset
func (kc *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	// use dockerhub
	return registry.NewClient("")
}

// Cluster return cluster describe
func (kc *Client) Cluster(cluster string) (*types.DescribeCluster, error) {
	clientCfg, err := utils.K8sV1Config([]byte(kc.Auth.KubeConfigString))
	if err != nil {
		return nil, err
	}
	for name := range clientCfg.Contexts {
		if name == cluster {
			return &types.DescribeCluster{
				ClusterMeta: types.ClusterMeta{
					Name:      name,
					UseUnique: name,
				},
			}, nil
		}
	}
	return &types.DescribeCluster{}, nil
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (kc *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := kc.RestConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("restconfig create error %v", err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, nil, fmt.Errorf("client create error %v", err)
	}
	return client, restConf, nil
}

// Clusters return aws cluster on region
func (kc *Client) Clusters() (*types.ClusterList, error) {
	clientCfg, err := utils.K8sV1Config([]byte(kc.Auth.KubeConfigString))
	if err != nil {
		return nil, err
	}
	clusterList := &types.ClusterList{Item: make([]types.DescribeCluster, 0)}
	for name := range clientCfg.Contexts {
		clusterList.Item = append(clusterList.Item, types.DescribeCluster{
			ClusterMeta: types.ClusterMeta{
				Name:      name,
				UseUnique: name,
			},
		})
	}
	return clusterList, nil
}

// RestConfig return rest config
func (kc *Client) RestConfig(cluster string) (*rest.Config, error) {
	clientCfg, err := utils.K8sV1Config([]byte(kc.Auth.KubeConfigString))
	if err != nil {
		return nil, err
	}
	return utils.RestConfig(clientCfg)
}

// Clientset return kubernetes clientset
func (kc *Client) Clientset(cluter string) (kubernetes.Interface, error) {
	restConfig, err := kc.RestConfig(cluter)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes dynamic kubernetes clientset
func (kc *Client) DynamicClientset(cluster string) (dynamic.Interface, error) {
	restConfig, err := kc.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

func newClient(regoinID string, kubecfg string) *Client {
	return &Client{
		RegionID: regoinID,
		Auth: Auth{
			KubeConfigString: kubecfg,
		},
	}
}

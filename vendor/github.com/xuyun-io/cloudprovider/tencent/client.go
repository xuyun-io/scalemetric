package tencent

import (
	"errors"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/xuyun-io/cloudprovider/types"
	"github.com/xuyun-io/cloudprovider/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client define tencent provider
type Client struct {
	// look to : https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/common/regions/regions.go
	Region string
	Auth   *Auth
}

// TCRAuth define tcr auth.
type TCRAuth struct {
	RegistryID string
	Namespace  string
	Host       string
	DockerAuth DockerAuth
}

// RegistryNameToHost generate registry host
// func RegistryNameToHost(registryName string) string {
// 	return fmt.Sprintf("%s.tencentcloudcr.com", registryName)
// }

// DockerAuth docker login auth.
type DockerAuth struct {
	// Host     string
	Username string
	Password string
}

// Auth define aws client auth
type Auth struct {
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
	TCRAuth         TCRAuth
}

// GetConfig return aws client config
func (tec *Auth) GetConfig() types.Auth {
	return tec
}

// DataCheck check aws client auth info
func (tec *Auth) DataCheck() error {
	if len(tec.AccessKey) <= 0 {
		return errors.New("AccessKey cann't be empty")
	}
	if len(tec.SecretAccessKey) <= 0 {
		return errors.New("SecretAccessKey cann't be empty")
	}
	return nil
}

// NewClient return tencent provider
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderTencent {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("huawei auth config get failed")
	}
	return NewClientByRegion(providercfg.GetSpec().GetRegionID(), auth), nil
}

// NewClientByRegion return tencent provider
func NewClientByRegion(region string, auth *Auth) *Client {
	return &Client{
		Region: region,
		Auth:   auth,
	}
}

// K8SClientset return k8s clientset
func (tec *Client) K8SClientset() types.KubernetesClientset {
	return tec
}

// Clusterset return aws provider cluster set
func (tec *Client) Clusterset() types.Clusterset {
	return tec
}

// NewTKEClient create tke client.
func (tec *Client) NewTKEClient() (*tke.Client, error) {
	credential := common.NewCredential(tec.Auth.AccessKey, tec.Auth.SecretAccessKey)
	return tke.NewClient(credential, tec.Region, profile.NewClientProfile())
}

// NewTCRClient create tcr client.
func (tec *Client) NewTCRClient() (*tcr.Client, error) {
	credential := common.NewCredential(tec.Auth.AccessKey, tec.Auth.SecretAccessKey)
	return tcr.NewClient(credential, tec.Region, profile.NewClientProfile())
}

// Cluster return cluster describe
func (tec *Client) Cluster(cluster string) (*types.DescribeCluster, error) {
	clusterList, err := tec.Clusters()
	if err != nil {
		return nil, err
	}
	for i := range clusterList.Item {
		if clusterList.Item[i].UseUnique == cluster {
			return &clusterList.Item[i], nil
		}
	}
	return nil, errors.New("cluster not found")
}

// Clusters return aws cluster on region
func (tec *Client) Clusters() (*types.ClusterList, error) {
	tkec, err := tec.NewTKEClient()
	if err != nil {
		return nil, err
	}
	request := tke.NewDescribeClustersRequest()
	resp, err := tkec.DescribeClusters(request)
	if err != nil {
		return nil, err
	}
	if resp.Response != nil && len(resp.Response.Clusters) > 0 {
		clusterList := &types.ClusterList{Item: make([]types.DescribeCluster, 0)}
		for i := range resp.Response.Clusters {
			clusterList.Item = append(clusterList.Item, types.DescribeCluster{
				ClusterMeta: types.ClusterMeta{
					UseUnique: *resp.Response.Clusters[i].ClusterId,
					Name:      *resp.Response.Clusters[i].ClusterName,
					State:     *resp.Response.Clusters[i].ClusterStatus,
					ClusterID: *resp.Response.Clusters[i].ClusterId,
					RegionID:  tec.Region,
				},
			})
		}
		return clusterList, nil
	}
	return nil, errors.New("clusters get failed from tke")
}

// RestConfig return kubernetes rest config
func (tec *Client) RestConfig(cluster string) (*rest.Config, error) {
	tkec, err := tec.NewTKEClient()
	if err != nil {
		return nil, err
	}
	request := tke.NewDescribeClusterSecurityRequest()
	request.ClusterId = common.StringPtr(cluster)
	resp, err := tkec.DescribeClusterSecurity(request)
	if err != nil {
		return nil, err
	}
	if resp.Response != nil && resp.Response.Kubeconfig != nil {
		cfg := *resp.Response.Kubeconfig
		if len(cfg) <= 0 {
			return nil, errors.New("kubeconfig not found")
		}
		clientCfg, err := utils.K8sV1Config([]byte(cfg))
		if err != nil {
			return nil, err
		}
		return utils.RestConfig(clientCfg)
	}
	return nil, errors.New("kubeconfig get failed from tke")
}

// Clientset return kubernetes clientset
func (tec *Client) Clientset(cluster string) (kubernetes.Interface, error) {
	restConfig, err := tec.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes dynamic kubernetes clientset
func (tec *Client) DynamicClientset(cluster string) (dynamic.Interface, error) {
	restConfig, err := tec.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (tec *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := tec.RestConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("restconfig create error %v", err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, nil, fmt.Errorf("client create error %v", err)
	}
	return client, restConf, nil
}

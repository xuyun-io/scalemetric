package huawei

import (
	"errors"
	"fmt"

	"github.com/imroc/req"
	"github.com/xuyun-io/cloudprovider/types"
	"github.com/xuyun-io/cloudprovider/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// 1. Unauthorized

// https://support.huaweicloud.com/usermanual-cce/cce_01_0188.html
// CCE Administrator

// 2. Forbidden

// AuthType define auth type
type AuthType string

// auth
const (
	AuthCCE AuthType = "AuthCCE"
	AuthSWR AuthType = "AuthSWR"
	AuthALL AuthType = "AuthALL"
)

// Auth define huawei client auth
type Auth struct {
	AuthType AuthType
	Base     BaseAuth
	CCE      CCEAuth
	SWR      SWRAuth
}

// BaseAuth define base auth by sub account.
type BaseAuth struct {
	Region     string
	DomainName string
	Account    string
	Password   string
}

// CCEAuth define cce auth
type CCEAuth struct {
	ProjectID   string
	ProjectName string
}

// SWRAuth define swr auth
type SWRAuth struct {
	Namespace       string
	ProjectName     string
	AccessKey       string
	SecretAccessKey string
}

// GetConfig define huawei client auth
func (hw *Auth) GetConfig() types.Auth {
	return hw
}

func (hw *Auth) cceAuth() error {
	if len(hw.CCE.ProjectID) <= 0 {
		return errors.New("cce projectid cann't be empty")
	}
	if len(hw.CCE.ProjectName) <= 0 {
		return errors.New("cce ProjectName cann't be empty")
	}
	return nil
}

func (hw *Auth) swrAuth() error {
	if len(hw.SWR.Namespace) <= 0 {
		return errors.New("swr namespace cann't be empty")
	}
	if len(hw.SWR.ProjectName) <= 0 {
		return errors.New("swr ProjectName cann't be empty")
	}
	if len(hw.SWR.AccessKey) <= 0 {
		return errors.New("swr AccessKey cann't be empty")
	}
	if len(hw.SWR.SecretAccessKey) <= 0 {
		return errors.New("swr SecretAccessKey cann't be empty")
	}
	return nil
}

// DataCheck return huawei client config
func (hw *Auth) DataCheck() error {
	if len(hw.Base.Region) <= 0 {
		return errors.New("region cann't be empty")
	}
	if len(hw.Base.DomainName) <= 0 {
		return errors.New("domian cann't be empty")
	}
	if len(hw.Base.Account) <= 0 {
		return errors.New("sub account cann't be empty")
	}
	if len(hw.Base.Password) <= 0 {
		return errors.New("password cann't be empty")
	}
	switch hw.AuthType {
	case AuthCCE:
		return hw.cceAuth()
	case AuthSWR:
		return hw.swrAuth()
	default:
		if err := hw.cceAuth(); err != nil {
			return err
		}
		if err := hw.swrAuth(); err != nil {
			return err
		}
	}
	return nil
}

// Client define huawei provider
type Client struct {
	Auth Auth
}

// NewClient return huawei provider
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderHuawei {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("huawei auth config get failed")
	}
	return newClientByRegion(providercfg.GetSpec().GetRegionID(), *auth)
}

// NewRepoClient return swr client.
func NewRepoClient(regionID string, baseAuth BaseAuth, swrAuth SWRAuth) types.RepoInterface {
	auth := Auth{
		AuthType: AuthSWR,
		Base:     baseAuth,
		SWR:      swrAuth,
	}
	auth.Base.Region = regionID
	c, _ := newClient(auth)
	return c
}

func newClientByRegion(regionID string, auth Auth) (types.Client, error) {
	auth.Base.Region = regionID
	return &Client{
		Auth: auth,
	}, nil
}

func newClient(auth Auth) (*Client, error) {
	return &Client{
		Auth: auth,
	}, nil
}

func (hw *Client) getCCEToken() (string, error) {
	return getCCEToken(hw.Auth)
}

func (hw *Client) getSWRToken() (string, error) {
	return getSWRToken(hw.Auth)
}

func (hw *Client) getkubeconfig(cluster string) (string, error) {
	token, err := hw.getCCEToken()
	if err != nil {
		return "", err
	}
	return hw.getkubeconfigByToken(token, cluster)
}

func (hw *Client) getkubeconfigByToken(token, cluster string) (string, error) {
	url := generateKubeconfigURL(hw.Auth.Base.Region, hw.Auth.CCE.ProjectID, cluster)
	header := generateTokenHeader(token)
	resp, err := req.Get(url, header)
	if err != nil {
		return "", err
	}
	if resp.Response().StatusCode >= 500 {
		return "", fmt.Errorf("huawei auth failed, http code != %d", resp.Response().StatusCode)
	}

	kubecfg, err := resp.ToString()
	if err != nil {
		return "", err
	}
	return kubecfg, nil
}

// K8SClientset return k8s clientset
func (hw *Client) K8SClientset() types.KubernetesClientset {
	return hw
}

// Clusterset return aws provider cluster set
func (hw *Client) Clusterset() types.Clusterset {
	return hw
}

// Clusters return aws cluster on region
func (hw *Client) Clusters() (*types.ClusterList, error) {
	return hw.getClusters()
}

// Cluster return cluster describe
func (hw *Client) Cluster(cluster string) (*types.DescribeCluster, error) {
	return hw.getCluster(cluster)
}

func server(region string) string {
	return fmt.Sprintf("https://cce.%s.myhuaweicloud.com/", region)
}

// RestConfig return rest config
func (hw *Client) RestConfig(cluster string) (*rest.Config, error) {
	token, err := hw.getCCEToken()
	if err != nil {
		return nil, err
	}
	cfg, err := hw.getkubeconfigByToken(token, cluster)
	if err != nil {
		return nil, err
	}
	clientCfg, err := utils.K8sV1Config([]byte(cfg))
	if err != nil {
		return nil, err
	}
	return utils.RestConfig(clientCfg)
}

// Clientset return kubernetes clientset
func (hw *Client) Clientset(cluster string) (kubernetes.Interface, error) {
	restConfig, err := hw.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes dynamic kubernetes clientset
func (hw *Client) DynamicClientset(cluster string) (dynamic.Interface, error) {
	restConfig, err := hw.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (hw *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := hw.RestConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("restconfig create error %v", err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, nil, fmt.Errorf("client create error %v", err)
	}
	return client, restConf, nil
}

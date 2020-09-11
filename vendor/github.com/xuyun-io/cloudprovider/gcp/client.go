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

package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	container "google.golang.org/api/container/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/xuyun-io/cloudprovider/registry"
	"github.com/xuyun-io/cloudprovider/types"
)

// Client define gcp config
type Client struct {
	ZoneID string `json:"zone_id"`
	Auth   Auth
	Repo   types.RepoInterface
}

// Auth define gcp client auth
type Auth struct {
	Config `json:"config"`
}

// GetConfig return gcp client config
func (g *Auth) GetConfig() types.Auth {
	return g
}

// DataCheck check gcp client auth info
func (g *Auth) DataCheck() error {
	return nil
}

// Config define fcp connect config
type Config struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// Regions returns all regions supported by gcp
func Regions() *types.RegionList {
	return regionList
}

// JSONConfig return gcp connection config
func (gcp *Client) JSONConfig() ([]byte, error) {
	return json.Marshal(gcp.Auth.Config)
}

// K8SClientset return k8s clientset
func (gcp *Client) K8SClientset() types.KubernetesClientset {
	return gcp
}

// Clusterset return aws provider cluster set
func (gcp *Client) Clusterset() types.Clusterset {
	return gcp
}

// RepoClientset return aws repo clientse
func (gcp *Client) RepoClientset(repo ...types.RepoInterface) types.RepoInterface {
	if len(repo) > 0 {
		return repo[0]
	}
	return gcp.NewRepoClient("gcr.io")
}

// NewClient return aws provider
func NewClient(providercfg types.ProviderConfig) (types.Client, error) {
	if providercfg.GetSpec().GetProvider() != types.ProviderGCP {
		return nil, fmt.Errorf("provider and auth do not match on ProviderConfig")
	}
	auth, ok := providercfg.GetSpec().GetAuth().(*Auth)
	if !ok {
		return nil, fmt.Errorf("ali auth config get failed")
	}

	return newClient(providercfg.GetSpec().GetRegionID(), auth.Config), nil

}

// NewClient return gcp client
func newClient(zoneID string, cfg Config) *Client {
	cfg.AuthURI = "https://accounts.google.com/o/oauth2/auth"
	cfg.TokenURI = "https://accounts.google.com/o/oauth2/token"
	cfg.AuthProviderX509CertURL = "https://www.googleapis.com/oauth2/v1/certs"
	return &Client{
		ZoneID: zoneID,
		Auth: Auth{
			Config: cfg,
		},
	}
}

// NewRepoClient return gcp client
func (gcp *Client) NewRepoClient(host string) types.RepoInterface {
	ctx := context.Background()
	token, err := gcp.token(ctx)
	if err != nil {
		return &RepoClient{Client: registry.Client{Error: err}}
	}
	opt := registry.NewBearerTokenAuth(token)
	return &RepoClient{Client: *newRepoClient(host, opt)}
}

// Clusters return gcp clusters on region
func (gcp *Client) Clusters() (*types.ClusterList, error) {
	ctx := context.Background()
	cs, err := gcp.ContainerClient(ctx)
	if err != nil {
		return nil, err
	}
	listClusters, err := cs.Projects.Zones.Clusters.List(gcp.Auth.ProjectID, gcp.ZoneID).Do()
	if err != nil {
		return nil, err
	}
	if len(listClusters.Clusters) <= 0 {
		return nil, err
	}
	clusterList := &types.ClusterList{
		Item: make([]types.DescribeCluster, 0),
	}
	for i := range listClusters.Clusters {
		clusterList.Item = append(clusterList.Item, types.DescribeCluster{
			ClusterMeta: types.ClusterMeta{
				Name:      listClusters.Clusters[i].Name,
				UseUnique: listClusters.Clusters[i].Name,
				State:     listClusters.Clusters[i].Status,
				ZoneID:    listClusters.Clusters[i].Zone,
			},
		})
	}
	return clusterList, nil
}

// Cluster return gcp cluster on gcp
func (gcp *Client) Cluster(cluster string) (*types.DescribeCluster, error) {
	ctx := context.Background()
	cs, err := gcp.ContainerClient(ctx)
	if err != nil {
		return nil, err
	}
	clusterDescribe, err := cs.Projects.Zones.Clusters.Get(gcp.Auth.ProjectID, gcp.ZoneID, cluster).Do()
	if err != nil {
		return nil, err
	}
	// clusterDescribe
	return &types.DescribeCluster{
		ClusterMeta: types.ClusterMeta{
			Name:      clusterDescribe.Name,
			UseUnique: clusterDescribe.Name,
			State:     clusterDescribe.Status,
			ZoneID:    clusterDescribe.Zone,
		},
	}, nil

}

// Credentials return credenails
func (gcp *Client) Credentials(ctx context.Context) (*google.Credentials, error) {
	cfgByts, err := gcp.JSONConfig()
	if err != nil {
		return nil, err
	}
	return google.CredentialsFromJSON(ctx, cfgByts, container.CloudPlatformScope)

}

// ContainerClient return gcp container client
func (gcp *Client) ContainerClient(ctx context.Context) (*container.Service, error) {
	creds, err := gcp.Credentials(ctx)
	if err != nil {
		return nil, err
	}
	c := oauth2.NewClient(ctx, creds.TokenSource)
	cs, err := container.New(c)
	return cs, err
}

// RestConfig return kubernetes client rest config
func (gcp *Client) RestConfig(cluster string) (*rest.Config, error) {
	ctx := context.Background()
	creds, err := gcp.Credentials(ctx)
	if err != nil {
		return nil, err
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, err
	}
	bearerToken := token.AccessToken
	c := oauth2.NewClient(ctx, creds.TokenSource)
	cs, err := container.New(c)
	describeCluster, err := cs.Projects.Zones.Clusters.Get(gcp.Auth.ProjectID, gcp.ZoneID, cluster).Do()
	if err != nil {
		return nil, err
	}
	clusterCaCertificate, err := base64.StdEncoding.DecodeString(describeCluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}
	return &rest.Config{
		Host: describeCluster.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: clusterCaCertificate,
		},
		BearerToken: bearerToken,
	}, nil
}

// ClientsetAndRestConifg return kubernetes clientset and restconfig
func (gcp *Client) ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error) {
	restConf, err := gcp.RestConfig(cluster)
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
func (gcp *Client) Clientset(cluster string) (kubernetes.Interface, error) {
	restConfig, err := gcp.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// DynamicClientset return kubernetes  dynamic kubernetes clientset
func (gcp *Client) DynamicClientset(cluster string) (dynamic.Interface, error) {
	restConfig, err := gcp.RestConfig(cluster)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}
func (gcp *Client) token(ctx context.Context) (string, error) {
	cfgByts, err := gcp.JSONConfig()
	if err != nil {
		return "", err
	}
	creds, err := google.CredentialsFromJSON(ctx, cfgByts, container.CloudPlatformScope)
	if err != nil {
		return "", err
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (gcp *Client) httpClient(ctx context.Context) (*http.Client, error) {
	creds, err := gcp.Credentials(ctx)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, creds.TokenSource), nil
}

// BearToken return oauth2 token
func (gcp *Client) bearToken(ctx context.Context) (*oauth2.Token, error) {
	creds, err := gcp.Credentials(ctx)
	if err != nil {
		return nil, err
	}
	return creds.TokenSource.Token()
}

// ImageDescribe define image detail
type ImageDescribe struct {
	Manifest map[string]interface{} `json:"manifest"`
	Name     string                 `json:"name"`
	Tags     []string               `json:"tags"`
}

/*
link:

eks:
	https://developers.google.com/accounts/docs/application-default-credentials
	https://console.developers.google.com/apis/api/container.googleapis.com/overview?project=686192924235

Container Register:
	https://cloud.google.com/container-registry/docs/pushing-and-pulling?hl=zh-cn
	https://cloud.google.com/container-registry/docs/advanced-authentication
	https://cloud.google.com/container-registry/docs/apis?authuser=1
*/

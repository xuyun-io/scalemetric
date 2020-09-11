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

package cloudprovider

import (
	"context"

	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/xuyun-io/cloudprovider/ali"
	"github.com/xuyun-io/cloudprovider/ali/acree"
	"github.com/xuyun-io/cloudprovider/aws"
	"github.com/xuyun-io/cloudprovider/gcp"
	"github.com/xuyun-io/cloudprovider/huawei"
	"github.com/xuyun-io/cloudprovider/kubeconfig"
	"github.com/xuyun-io/cloudprovider/registry"
	"github.com/xuyun-io/cloudprovider/tencent"
	"github.com/xuyun-io/cloudprovider/types"
)

var (
	regionMaps = map[types.Provider]*types.RegionList{
		types.ProviderAWS:      aws.Regions(),
		types.ProviderGCP:      gcp.Regions(),
		types.ProviderALiCloud: ali.Regions(),
	}
)

// ProviderList return provider list
func ProviderList() map[types.Provider]bool {
	return types.SupportProviders()
}

// RegionList return provider region list
func RegionList(provider types.Provider) *types.RegionList {
	return regionMaps[provider]
}

// NewDockerRegistry return docker registry client
func NewDockerRegistry(host string, options ...remote.Option) types.RepoInterface {
	return registry.NewClient(host, options...)
}

// NewGCPRepoClient return gcp repository client
func NewGCPRepoClient(ctx context.Context, host string, cfg gcp.Config) types.RepoInterface {
	return gcp.NewRepoClient(ctx, host, cfg)
}

// NewAWSRepoClient return aws repository client
func NewAWSRepoClient(regionID, accessKey, secretAccessKey string) types.RepoInterface {
	return aws.NewRepoClient(regionID, accessKey, secretAccessKey)
}

// NewAliRepoEEClient return ali container register ee client
func NewAliRepoEEClient(regionID, accessKey, secretAccessKey string, params acree.RepoEEParams) types.RepoInterface {
	return ali.NewRepoEEClient(regionID, accessKey, secretAccessKey, params)
}

// NewAliRepoClient return ali acr.
func NewAliRepoClient(regionID, accessKey, secretAccessKey string, params ali.RepoParams) types.RepoInterface {
	return ali.NewRepoClient(regionID, accessKey, secretAccessKey, params)
}

// NewHuaweiRepoClient return huawei swr.
func NewHuaweiRepoClient(regionID string, baseAuth huawei.BaseAuth, swrAuth huawei.SWRAuth) types.RepoInterface {
	return huawei.NewRepoClient(regionID, baseAuth, swrAuth)
}

// NewTencentRepoClient return tencent tcr.
func NewTencentRepoClient(regionID string, auth *tencent.Auth) types.RepoInterface {
	return tencent.NewClientByRegion(regionID, auth)
}

// NewHuaweiProviderConfig new huawei provider config
func NewHuaweiProviderConfig(regionID string, conf huawei.Auth) types.ProviderConfig {
	cfg := NewProviderConfig()
	cfg.SetAuth(types.ProviderHuawei, regionID, &conf)
	return cfg
}

// NewAliProviderConfig new ali provider config
func NewAliProviderConfig(regionID, accessKey, secretAccessKey string) types.ProviderConfig {
	cfg := NewProviderConfig()
	auth := &ali.Auth{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
	}
	cfg.SetAuth(types.ProviderALiCloud, regionID, auth)
	return cfg
}

// NewAWSProviderConfig new aws provider config
func NewAWSProviderConfig(regionID, accessKey, secretAccessKey string) types.ProviderConfig {
	auth := &aws.Auth{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
	}
	cfg := NewProviderConfig()
	cfg.SetAuth(types.ProviderAWS, regionID, auth)
	return cfg
}

// NewTencentProviderConfig new tencent provider config
func NewTencentProviderConfig(regionID string, auth *tencent.Auth) types.ProviderConfig {
	cfg := NewProviderConfig()
	cfg.SetAuth(types.ProviderTencent, regionID, auth)
	return cfg
}

// NewGCPProviderConfig new gcp provider config
func NewGCPProviderConfig(zoneID string, conf gcp.Config) types.ProviderConfig {
	auth := &gcp.Auth{
		Config: conf,
	}
	cfg := NewProviderConfig()
	cfg.SetAuth(types.ProviderGCP, zoneID, auth)
	return cfg
}

// NewKubeConfigProviderConfig new kubeconfig provider config
func NewKubeConfigProviderConfig(kubecfg string) types.ProviderConfig {
	auth := &kubeconfig.Auth{
		KubeConfigString: kubecfg,
	}
	cfg := NewProviderConfig()
	cfg.SetAuth(types.ProviderKubeconfig, "-", auth)
	return cfg
}

// NewClient return provider client
func NewClient(providerCfg types.ProviderConfig) (types.Client, error) {
	return providerCfg.NewClient()
}

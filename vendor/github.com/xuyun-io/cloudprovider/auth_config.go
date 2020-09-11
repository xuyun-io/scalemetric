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
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xuyun-io/cloudprovider/ali"
	"github.com/xuyun-io/cloudprovider/aws"
	"github.com/xuyun-io/cloudprovider/gcp"
	"github.com/xuyun-io/cloudprovider/huawei"
	"github.com/xuyun-io/cloudprovider/kubeconfig"
	"github.com/xuyun-io/cloudprovider/tencent"
	"github.com/xuyun-io/cloudprovider/types"
)

// Config define provider config
type Config struct {
	metav1.TypeMeta `json:",inline"`
	Spec            ConfigSpec
	Status          ConfigStatus
}

// ConfigSpec define provider config spec
type ConfigSpec struct {
	Name     string         `json:"name"`
	RegionID string         `json:"regionId"`
	Provider types.Provider `json:"provider"`
	Auth     types.Auth     `json:"auth"`
}

type getClient func(providerConfig types.ProviderConfig) (types.Client, error)

// ConfigStatus define provider config status
type ConfigStatus struct{}

// define provider config meta data
const (
	KindProviderConfig = "ProviderConfig"
	APIVersion         = "v1beta1"
)

var (
	_                  types.ProviderConfig = &Config{}
	providerGetClients                      = map[types.Provider]getClient{
		// kubecofig  get client func
		types.ProviderKubeconfig: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return kubeconfig.NewClient(providerConfig)
		},
		// ali get client func
		types.ProviderALiCloud: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return ali.NewClient(providerConfig)
		},
		// aws get client func
		types.ProviderAWS: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return aws.NewClient(providerConfig)
		},
		// gcp get client func
		types.ProviderGCP: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return gcp.NewClient(providerConfig)
		},
		types.ProviderHuawei: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return huawei.NewClient(providerConfig)
		},
		types.ProviderTencent: func(providerConfig types.ProviderConfig) (types.Client, error) {
			return tencent.NewClient(providerConfig)
		},
	}
)

// NewProviderConfig return providerconfig
func NewProviderConfig() types.ProviderConfig {
	return &Config{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindProviderConfig,
			APIVersion: APIVersion,
		},
	}
}

// NewClient return provider client by provider config
func (cfg *Config) NewClient() (types.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	clientFunc, ok := providerGetClients[cfg.GetProvider()]
	if !ok {
		return nil, fmt.Errorf("provider client create err")
	}
	return clientFunc(cfg)
}

// GetSpec return provider config spec
func (cfg *Config) GetSpec() types.Spec {
	return cfg.Spec
}

// Validate validate provider info
func (cfg *Config) Validate() error {
	if reflect.DeepEqual(cfg.Spec, ConfigSpec{}) {
		return fmt.Errorf("spec cann't emtpy")
	}
	return cfg.Spec.Validate()
}

// SetAuth set provider auth config
func (cfg *Config) SetAuth(provider types.Provider, regionID string, auth types.Auth) {
	spec := ConfigSpec{
		RegionID: regionID,
		Provider: provider,
		Auth:     auth,
	}
	cfg.Spec = spec
}

// GetTypeMeta return provider meta data
func (cfg *Config) GetTypeMeta() metav1.TypeMeta {
	return cfg.TypeMeta
}

// GetProvider return provider
func (cfg *Config) GetProvider() types.Provider {
	return cfg.Spec.Provider
}

// GetRegionID return provider region
func (cfg *Config) GetRegionID() string {
	return cfg.Spec.RegionID
}

// GetAuth return provider auth info
func (cfg *Config) GetAuth() types.Auth {
	return cfg.Spec.Auth
}

// Validate validate provider auth info
func (cs ConfigSpec) Validate() error {
	if _, ok := types.SupportProviders()[cs.Provider]; !ok {
		return fmt.Errorf("provider: %s not supported, support list: %v", cs.Provider, types.SupportProviders())
	}
	if cs.Provider != types.ProviderKubeconfig {
		if len(cs.RegionID) <= 0 {
			return fmt.Errorf("RegionID or ZoneID cann't be emtpy")
		}
	}
	if cs.Auth != nil {
		if err := cs.Auth.DataCheck(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("auth info cann't be emtpy")
}

// GetProvider return provider
func (cs ConfigSpec) GetProvider() types.Provider {
	return cs.Provider
}

// GetRegionID return provider regionid
func (cs ConfigSpec) GetRegionID() string {
	return cs.RegionID
}

// GetAuth return provider config spec
func (cs ConfigSpec) GetAuth() types.Auth {
	return cs.Auth
}

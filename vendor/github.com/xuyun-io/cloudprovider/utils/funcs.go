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

package utils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
)

// K8sV1Config return config from kubeconfig
// params: cfgByts is kubeconfig bytes
func K8sV1Config(cfgByts []byte) (*clientcmdapi.Config, error) {
	cfg := clientcmdapi.NewConfig()
	decoded, _, err := clientcmdlatest.Codec.Decode(cfgByts, &schema.GroupVersionKind{Version: clientcmdlatest.Version, Kind: "Config"}, cfg)
	if err != nil {
		return nil, err
	}
	return decoded.(*api.Config), nil
}

// ServerEndPoint return kubernetes apiserver endpoint
func ServerEndPoint(config *clientcmdapi.Config) (string, error) {
	if config == nil {
		return "", fmt.Errorf("server not found, config is nil")
	}
	if len(config.Contexts) <= 0 {
		return "", fmt.Errorf("server not found, contexts is empty")
	}
	currentCluster := ""
	for name, v := range config.Contexts {
		if config.CurrentContext == name {
			currentCluster = v.Cluster
			break
		}
	}
	if len(currentCluster) <= 0 {
		return "", fmt.Errorf("server not found, currentCluster is empty")
	}
	for k, v := range config.Clusters {
		if k == currentCluster {
			return v.Server, nil
		}
	}
	return "", fmt.Errorf("server not found, currentCluster and contexts do not match")
}

// RestConfig return rest config
func RestConfig(config *clientcmdapi.Config) (*rest.Config, error) {
	getter := func() (*clientcmdapi.Config, error) { return config, nil }
	masterURL, err := ServerEndPoint(config)
	if err != nil {
		return nil, err
	}
	return clientcmd.BuildConfigFromKubeconfigGetter(masterURL, getter)
}

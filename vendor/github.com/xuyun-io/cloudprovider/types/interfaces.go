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

package types

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client define provider clientset
type Client interface {
	K8SClientset() KubernetesClientset
	Clusterset() Clusterset
	RepoClientset(repo ...RepoInterface) RepoInterface
}

// Clusterset define provider clusterset
type Clusterset interface {
	Cluster(cluster string) (*DescribeCluster, error)
	Clusters() (*ClusterList, error)
}

// KubernetesClientset define kubernetes clientset
type KubernetesClientset interface {
	RestConfig(cluster string) (*rest.Config, error)
	Clientset(cluster string) (kubernetes.Interface, error)
	ClientsetAndRestConifg(cluster string) (kubernetes.Interface, *rest.Config, error)
	DynamicClientset(cluster string) (dynamic.Interface, error)
}

// RepoInterface define repository clientset
type RepoInterface interface {
	Repositories(ctx context.Context) ([]*Repository, error)
	TagList(ctx context.Context, repoName string) ([]*ImageDetail, error)
	ImageDescribe(ctx context.Context, imageName string) (interface{}, error)
}

// ProviderConfig define provider config interface
type ProviderConfig interface {
	NewClient() (Client, error)
	Validate() error
	GetTypeMeta() metav1.TypeMeta
	GetSpec() Spec
	SetAuth(provider Provider, regionID string, auth Auth)
}

// Spec define provicer config spec interface
type Spec interface {
	GetProvider() Provider
	GetRegionID() string
	GetAuth() Auth
}

// Auth define provider auth interface
type Auth interface {
	GetConfig() Auth
	DataCheck() error
}

package huawei

import "time"

// DockerUser define docker login message.
type DockerUser struct {
	User     string
	Password string
	Host     string
}

// ClusterList define clusterList
type ClusterList struct {
	Kind       string    `json:"kind"`
	APIVersion string    `json:"apiVersion"`
	Items      []Cluster `json:"items"`
}

// Cluster define cluster message.
type Cluster struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   `json:"metadata"`
	Spec       `json:"spec"`
	Status     `json:"status"`
}

// Metadata define metadata message.
type Metadata struct {
	Name              string `json:"name"`
	UID               string `json:"uid"`
	CreationTimestamp string `json:"creationTimestamp"`
	UpdateTimestamp   string `json:"updateTimestamp"`
}

// Spec define cluster spec.
type Spec struct {
	Type                 string `json:"type"`
	Flavor               string `json:"flavor"`
	Version              string `json:"version"`
	Ipv6Enable           bool   `json:"ipv6enable"`
	HostNetwork          `json:"hostNetwork"`
	ContainerNetwork     `json:"containerNetwork"`
	EniNetwork           `json:"eniNetwork"`
	Authentication       `json:"authentication"`
	BillingMode          int    `json:"billingMode"`
	KubernetesSvcIPRange string `json:"kubernetesSvcIpRange"`
	KubeProxyMode        string `json:"kubeProxyMode"`
	Az                   string `json:"az"`
	ExtendParam          `json:"extendParam"`
	SupportIstio         bool `json:"supportIstio"`
}

// HostNetwork define cluster HostNetwork
type HostNetwork struct {
	Vpc           string `json:"vpc"`
	Subnet        string `json:"subnet"`
	SecurityGroup string `json:"SecurityGroup"`
}

// ContainerNetwork deine container network.
type ContainerNetwork struct {
	Mode string `json:"mode"`
	Cidr string `json:"cidr"`
}

// EniNetwork define eni
type EniNetwork struct {
}

// Authentication define cluster auth.
type Authentication struct {
	Mode                string   `json:"mode"`
	AuthenticatingProxy struct{} `json:"authenticatingProxy"`
}

// ExtendParam define auth cluster
type ExtendParam struct {
	AlphaCceFixPoolMask          string `json:"alpha.cce/fixPoolMask"`
	KubernetesIoCPUManagerPolicy string `json:"kubernetes.io/cpuManagerPolicy"`
	Upgradefrom                  string `json:"upgradefrom"`
}

// Status define cluster status.
type Status struct {
	Phase     string `json:"phase"`
	Endpoints []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"endpoints"`
}

// RepoBody define repository body.
type RepoBody struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Size         int      `json:"size"`
	IsPublic     bool     `json:"is_public"`
	NumImages    int      `json:"num_images"`
	NumDownload  int      `json:"num_download"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	Logo         string   `json:"logo"`
	URL          string   `json:"url"`
	Path         string   `json:"path"`
	InternalPath string   `json:"internal_path"`
	DomainName   string   `json:"domain_name"`
	Namespace    string   `json:"namespace"`
	Tags         []string `json:"tags"`
	Status       bool     `json:"status"`
}

// TagBody define swr tag body.
type TagBody struct {
	ID           int         `json:"id"`
	RepoID       int         `json:"repo_id"`
	Tag          string      `json:"Tag"`
	ImageID      string      `json:"image_id"`
	Manifest     string      `json:"manifest"`
	Digest       string      `json:"digest"`
	Schema       int         `json:"schema"`
	Path         string      `json:"path"`
	InternalPath string      `json:"internal_path"`
	Size         int64       `json:"size"`
	IsTrusted    bool        `json:"is_trusted"`
	Created      time.Time   `json:"created"`
	Updated      time.Time   `json:"updated"`
	Deleted      interface{} `json:"deleted"`
}

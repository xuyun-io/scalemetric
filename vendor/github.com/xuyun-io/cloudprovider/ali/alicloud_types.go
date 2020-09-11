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
	"time"

	"github.com/xuyun-io/cloudprovider/types"
)

type k8sConfig struct {
	Config string `json:"config,omitempty"`
}

// ClusterList define ali cloud cluster list
type ClusterList struct {
	Item []DescribeCluster `json:"item,omitempty"`
}

// Region define ali region struct
type Region struct {
	RegionID       string `json:"RegionId"`
	RegionEndpoint string `json:"RegionEndpoint"`
	LocalName      string `json:"LocalName"`
}

// DescribeCluster define cluster describe in ali cloud
type DescribeCluster struct {
	types.ClusterMeta
	InstanceType           string      `json:"instance_type"`
	VpcID                  string      `json:"vpc_id"`
	VswitchID              string      `json:"vswitch_id"`
	VswitchCidr            string      `json:"vswitch_cidr"`
	DataDiskSize           int         `json:"data_disk_size"`
	DataDiskCategory       string      `json:"data_disk_category"`
	SecurityGroupID        string      `json:"security_group_id"`
	Tags                   []Tag       `json:"tags"`
	Size                   int         `json:"size"`
	NetworkMode            string      `json:"network_mode"`
	SubnetCidr             string      `json:"subnet_cidr"`
	MasterURL              string      `json:"master_url"`
	ExternalLoadbalancerID string      `json:"external_loadbalancer_id"`
	Created                time.Time   `json:"created"`
	Updated                time.Time   `json:"updated"`
	Port                   int         `json:"port"`
	NodeStatus             string      `json:"node_status"`
	ClusterHealthy         string      `json:"cluster_healthy"`
	DockerVersion          string      `json:"docker_version"`
	ClusterType            string      `json:"cluster_type"`
	SwarmMode              bool        `json:"swarm_mode"`
	InitVersion            string      `json:"init_version"`
	CurrentVersion         string      `json:"current_version"`
	MetaData               string      `json:"meta_data"`
	GwBridge               string      `json:"gw_bridge"`
	ResourceGroupID        string      `json:"resource_group_id"`
	PrivateZone            bool        `json:"private_zone"`
	Profile                string      `json:"profile"`
	DeletionProtection     bool        `json:"deletion_protection"`
	Capabilities           interface{} `json:"capabilities"`
	EnabledMigration       bool        `json:"enabled_migration"`
	NeedUpdateAgent        bool        `json:"need_update_agent"`
	Outputs                []Output    `json:"outputs"`
	UpgradeComponents      `json:"upgrade_components"`
	Parameters             `json:"parameters"`
}

// Tag define cluster tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// UpgradeComponents define cluster component types
type UpgradeComponents struct {
	Kubernetes `json:"Kubernetes"`
}

// Kubernetes define kubernetes component
type Kubernetes struct {
	ComponentName  string      `json:"component_name"`
	Version        string      `json:"version"`
	NextVersion    string      `json:"next_version"`
	Changed        string      `json:"changed"`
	CanUpgrade     bool        `json:"can_upgrade"`
	Force          bool        `json:"force"`
	Policy         string      `json:"policy"`
	ExtraVars      interface{} `json:"ExtraVars"`
	ReadyToUpgrade string      `json:"ready_to_upgrade"`
	Message        string      `json:"message"`
	Exist          bool        `json:"exist"`
	Required       bool        `json:"required"`
}

// Output define output message
type Output struct {
	Description string      `json:"Description"`
	OutputKey   string      `json:"OutputKey"`
	OutputValue interface{} `json:"OutputValue"`
}

// Parameters define ali cloud parameters
type Parameters struct {
	ALIYUNAccountID          string `json:"ALIYUN::AccountId"`
	ALIYUNNoValue            string `json:"ALIYUN::NoValue"`
	ALIYUNRegion             string `json:"ALIYUN::Region"`
	ALIYUNStackID            string `json:"ALIYUN::StackId"`
	ALIYUNStackName          string `json:"ALIYUN::StackName"`
	ALIYUNTenantID           string `json:"ALIYUN::TenantId"`
	AdjustmentType           string `json:"AdjustmentType"`
	AuditFlags               string `json:"AuditFlags"`
	BetaVersion              string `json:"BetaVersion"`
	CA                       string `json:"CA"`
	ClientCA                 string `json:"ClientCA"`
	CloudMonitorFlags        string `json:"CloudMonitorFlags"`
	CloudMonitorVersion      string `json:"CloudMonitorVersion"`
	ContainerCIDR            string `json:"ContainerCIDR"`
	DisableAddons            string `json:"DisableAddons"`
	DockerVersion            string `json:"DockerVersion"`
	ESSDeletionProtection    string `json:"ESSDeletionProtection"`
	Eip                      string `json:"Eip"`
	EipAddress               string `json:"EipAddress"`
	ElasticSearchHost        string `json:"ElasticSearchHost"`
	ElasticSearchPass        string `json:"ElasticSearchPass"`
	ElasticSearchPort        string `json:"ElasticSearchPort"`
	ElasticSearchUser        string `json:"ElasticSearchUser"`
	EtcdVersion              string `json:"EtcdVersion"`
	ExecuteVersion           string `json:"ExecuteVersion"`
	GPUFlags                 string `json:"GPUFlags"`
	HealthCheckType          string `json:"HealthCheckType"`
	ImageID                  string `json:"ImageId"`
	K8SMasterPolicyDocument  string `json:"K8SMasterPolicyDocument"`
	K8SWorkerPolicyDocument  string `json:"K8sWorkerPolicyDocument"`
	Key                      string `json:"Key"`
	KeyPair                  string `json:"KeyPair"`
	KubernetesVersion        string `json:"KubernetesVersion"`
	LoggingType              string `json:"LoggingType"`
	MasterAmounts            string `json:"MasterAmounts"`
	MasterAutoRenew          string `json:"MasterAutoRenew"`
	MasterAutoRenewPeriod    string `json:"MasterAutoRenewPeriod"`
	MasterCount              string `json:"MasterCount"`
	MasterDataDisk           string `json:"MasterDataDisk"`
	MasterDataDisks          string `json:"MasterDataDisks"`
	MasterDeletionProtection string `json:"MasterDeletionProtection"`
	MasterDeploymentSetID    string `json:"MasterDeploymentSetId"`
	MasterHpcClusterID       string `json:"MasterHpcClusterId"`
	MasterImageID            string `json:"MasterImageId"`
	MasterInstanceChargeType string `json:"MasterInstanceChargeType"`
	MasterInstanceTypes      string `json:"MasterInstanceTypes"`
	MasterKeyPair            string `json:"MasterKeyPair"`
	MasterLoginPassword      string `json:"MasterLoginPassword"`
	MasterPeriod             string `json:"MasterPeriod"`
	MasterPeriodUnit         string `json:"MasterPeriodUnit"`
	MasterSlbSSHHealthCheck  string `json:"MasterSlbSShHealthCheck"`
	MasterSnapshotPolicyID   string `json:"MasterSnapshotPolicyId"`
	MasterSystemDiskCategory string `json:"MasterSystemDiskCategory"`
	MasterSystemDiskSize     string `json:"MasterSystemDiskSize"`
	MasterVSwitchIds         string `json:"MasterVSwitchIds"`
	NatGateway               string `json:"NatGateway"`
	NatGatewayID             string `json:"NatGatewayId"`
	Network                  string `json:"Network"`
	NodeCIDRMask             string `json:"NodeCIDRMask"`
	NodeNameMode             string `json:"NodeNameMode"`
	NumOfNodes               string `json:"NumOfNodes"`
	Password                 string `json:"Password"`
	PodVswitchIds            string `json:"PodVswitchIds"`
	ProtectedInstances       string `json:"ProtectedInstances"`
	ProxyMode                string `json:"ProxyMode"`
	PublicSLB                string `json:"PublicSLB"`
	RemoveInstanceIds        string `json:"RemoveInstanceIds"`
	SLBDeletionProtection    string `json:"SLBDeletionProtection"`
	SLSProjectName           string `json:"SLSProjectName"`
	SNatEntry                string `json:"SNatEntry"`
	SSHFlags                 string `json:"SSHFlags"`
	SecurityGroupID          string `json:"SecurityGroupId"`
	ServiceCIDR              string `json:"ServiceCIDR"`
	SetUpArgs                string `json:"SetUpArgs"`
	SnatTableID              string `json:"SnatTableId"`
	Tags                     string `json:"Tags"`
	UserCA                   string `json:"UserCA"`
	UserData                 string `json:"UserData"`
	VpcID                    string `json:"VpcId"`
	WillReplace              string `json:"WillReplace"`
	WorkerAutoRenew          string `json:"WorkerAutoRenew"`
	WorkerAutoRenewPeriod    string `json:"WorkerAutoRenewPeriod"`
	WorkerDataDisk           string `json:"WorkerDataDisk"`
	WorkerDataDisks          string `json:"WorkerDataDisks"`
	WorkerDeletionProtection string `json:"WorkerDeletionProtection"`
	WorkerDeploymentSetID    string `json:"WorkerDeploymentSetId"`
	WorkerHpcClusterID       string `json:"WorkerHpcClusterId"`
	WorkerImageID            string `json:"WorkerImageId"`
	WorkerInstanceChargeType string `json:"WorkerInstanceChargeType"`
	WorkerInstanceTypes      string `json:"WorkerInstanceTypes"`
	WorkerKeyPair            string `json:"WorkerKeyPair"`
	WorkerLoginPassword      string `json:"WorkerLoginPassword"`
	WorkerPeriod             string `json:"WorkerPeriod"`
	WorkerPeriodUnit         string `json:"WorkerPeriodUnit"`
	WorkerSnapshotPolicyID   string `json:"WorkerSnapshotPolicyId"`
	WorkerSystemDiskCategory string `json:"WorkerSystemDiskCategory"`
	WorkerSystemDiskSize     string `json:"WorkerSystemDiskSize"`
	WorkerVSwitchIds         string `json:"WorkerVSwitchIds"`
	ZoneID                   string `json:"ZoneId"`
}

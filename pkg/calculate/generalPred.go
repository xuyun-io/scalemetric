package calculate

import (
	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
	pluginhelper "k8s.io/kubernetes/pkg/scheduler/framework/plugins/helper"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/nodeaffinity"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/nodename"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/nodeports"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/noderesources"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

func generalNodeInfo(node *v1.Node, nodePods []*v1.Pod) *schedulernodeinfo.NodeInfo {
	var (
		nodeInfo = &schedulernodeinfo.NodeInfo{}
	)
	if len(nodePods) <= 0 {
		nodeInfo = schedulernodeinfo.NewNodeInfo(&v1.Pod{})
	} else {
		nodeInfo = schedulernodeinfo.NewNodeInfo(nodePods...)
	}
	nodeInfo.SetNode(node)
	return nodeInfo
}

// GeneralPred return node resource scheduling status.
func GeneralPred(node *v1.Node, nodePods []*v1.Pod, pod *v1.Pod) ([]types.PredicateFailureReason, error) {
	nodeInfo := generalNodeInfo(node, nodePods)
	var reasons []types.PredicateFailureReason
	for _, r := range noderesources.Fits(pod, nodeInfo, nil) {
		reasons = append(reasons, &InsufficientResourceError{
			ResourceName: r.ResourceName,
			Requested:    r.Requested,
			Used:         r.Used,
			Capacity:     r.Capacity,
		})
	}
	if !pluginhelper.PodMatchesNodeSelectorAndAffinityTerms(pod, nodeInfo.Node()) {
		reasons = append(reasons, &PredicateFailureError{nodeaffinity.Name, nodeaffinity.ErrReason})
	}
	if !nodename.Fits(pod, nodeInfo) {
		reasons = append(reasons, &PredicateFailureError{nodename.Name, nodename.ErrReason})
	}
	if !nodeports.Fits(pod, nodeInfo) {
		reasons = append(reasons, &PredicateFailureError{nodeports.Name, nodeports.ErrReason})
	}
	return reasons, nil

}

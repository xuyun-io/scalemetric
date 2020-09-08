package types

import (
	"github.com/xuyun-io/scalemetric/pkg/calculate"
	v1 "k8s.io/api/core/v1"
)

// ClusterScheduling define cluster scheduling status.
type ClusterScheduling struct {
	SchedulingStatus []PodRequestScheduling
}

// PodRequestScheduling define pod scheduling status.
type PodRequestScheduling struct {
	Pod                    *v1.Pod
	PredMaxschedulingCount int64
	NodeScheduling         []NodeScheduling
}

// NodeScheduling define node scheduling status.
type NodeScheduling struct {
	Node                   *v1.Node
	PredMaxschedulingCount int64
	LastReason             []calculate.PredicateFailureReason
}

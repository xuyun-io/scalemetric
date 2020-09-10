package types

import (
	v1 "k8s.io/api/core/v1"
)

// ClusterScheduling define cluster scheduling status.
type ClusterScheduling struct {
	SchedulingStatus []PodRequestScheduling `json:"schedulingStatus"`
}

// PodRequestScheduling define pod scheduling status.
type PodRequestScheduling struct {
	Pod                    *v1.Pod          `json:"pod"`
	PredMaxschedulingCount int64            `json:"predMaxschedulingCount"`
	NodeScheduling         []NodeScheduling `json:"nodeScheduling"`
}

// NodeScheduling define node scheduling status.
type NodeScheduling struct {
	Node                   *v1.Node                 `json:"node,omitempty"`
	PredMaxschedulingCount int64                    `json:"predMaxschedulingCount"`
	LastReason             []PredicateFailureReason `json:"lastReason,omitempty"`
}

// PredicateFailureReason interface represents the failure reason of a predicate.
type PredicateFailureReason interface {
	GetReason() string
	Error() string
}

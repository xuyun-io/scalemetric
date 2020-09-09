package calculate

import (
	"errors"
	"fmt"

	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// CheckNodeConditionPred check node ready status.
func CheckNodeConditionPred(node *v1.Node) ([]types.PredicateFailureReason, error) {
	var reasons []types.PredicateFailureReason
	predicate := getNodeConditionPredicate()
	if ok, err := predicate(node); !ok {
		reasons = append(reasons, &NodeConditionError{
			PredicateName: "CheckNodeConditionPred",
			Reason:        err.Error(),
		})
	}
	return reasons, nil

}

func getNodeConditionPredicate() NodeConditionPredicate {
	return func(node *v1.Node) (bool, error) {
		// We add the master to the node list, but its unschedulable.  So we use this to filter
		// the master.
		if node.Spec.Unschedulable {
			return false, errors.New("Unschedulable")
		}

		// if utilfeature.DefaultFeatureGate.Enabled(legacyNodeRoleBehaviorFeature) {
		// 	// As of 1.6, we will taint the master, but not necessarily mark it unschedulable.
		// 	// Recognize nodes labeled as master, and filter them also, as we were doing previously.
		// 	if _, hasMasterRoleLabel := node.Labels[labelNodeRoleMaster]; hasMasterRoleLabel {
		// 		return false
		// 	}
		// }
		// if utilfeature.DefaultFeatureGate.Enabled(serviceNodeExclusionFeature) {
		// 	// Will be removed in 1.18
		// 	if _, hasExcludeBalancerLabel := node.Labels[labelAlphaNodeRoleExcludeBalancer]; hasExcludeBalancerLabel {
		// 		return false
		// 	}
		// 	if _, hasExcludeBalancerLabel := node.Labels[labelNodeRoleExcludeBalancer]; hasExcludeBalancerLabel {
		// 		return false
		// 	}
		// }

		// If we have no info, don't accept
		if len(node.Status.Conditions) == 0 {
			return false, errors.New("conditions is empty, don't accept")
		}
		for _, cond := range node.Status.Conditions {
			// We consider the node for load balancing only when its NodeReady condition status
			// is ConditionTrue
			if cond.Type == v1.NodeReady && cond.Status != v1.ConditionTrue {
				// klog.V(4).Infof("Ignoring node %v with %v condition status %v", node.Name, cond.Type, cond.Status)
				return false, fmt.Errorf("%s, %s", cond.Reason, cond.Message)
			}
		}
		return true, nil
	}
}

// NodeConditionPredicate is a function that indicates whether the given node's conditions meet
// some set of criteria defined by the function.
type NodeConditionPredicate func(node *v1.Node) (bool, error)

// // listWithPredicate gets nodes that matches predicate function.
// func listWithPredicate(nodeLister corelisters.NodeLister, predicate NodeConditionPredicate) ([]*v1.Node, error) {
// 	nodes, err := nodeLister.List(labels.Everything())
// 	if err != nil {
// 		return nil, err
// 	}

// 	var filtered []*v1.Node
// 	for i := range nodes {
// 		if predicate(nodes[i]) {
// 			filtered = append(filtered, nodes[i])
// 		}
// 	}

// 	return filtered, nil
// }

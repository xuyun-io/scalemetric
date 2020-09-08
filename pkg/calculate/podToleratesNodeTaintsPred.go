package calculate

import  (
	"fmt"
	v1 "k8s.io/api/core/v1"
v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
// schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

func PodToleratesNodeTaintsPred(node *v1.Node,pod *v1.Pod) ([]PredicateFailureReason) {
	// if nodeInfo == nil || nodeInfo.Node() == nil {
	// 	return framework.NewStatus(framework.Error, "invalid nodeInfo")
	// }
	nodeInfo:= generalNodeInfo(node,nil)
	var reasons []PredicateFailureReason
	filterPredicate := func(t *v1.Taint) bool {
		// PodToleratesNodeTaints is only interested in NoSchedule and NoExecute taints.
		return t.Effect == v1.TaintEffectNoSchedule || t.Effect == v1.TaintEffectNoExecute
	}
	taint, isUntolerated := v1helper.FindMatchingUntoleratedTaint(nodeInfo.Node().Spec.Taints, pod.Spec.Tolerations, filterPredicate)
	if !isUntolerated {
		return reasons
	}
	errReason := fmt.Sprintf("node(s) had taint {%s: %s}, that the pod didn't tolerate",
	taint.Key, taint.Value)
	reasons = append(reasons, &PodToleratesNodeTaintsError{
		PredicateName: "PodToleratesNodeTaints",
		Reason:errReason,
	})
	return  reasons

}



// PodToleratesNodeTaintsError describes a failure error of predicate.
type PodToleratesNodeTaintsError struct {
	PredicateName string
	Reason        string
}

func (e *PodToleratesNodeTaintsError) Error() string {
	return fmt.Sprintf("Predicate %s failed, reason: %s", e.PredicateName, e.GetReason())
}

// GetReason returns the reason of the NodeConditionError.
func (e *PodToleratesNodeTaintsError) GetReason() string {
	return e.Reason
}

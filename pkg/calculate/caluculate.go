package calculate

import (
	"encoding/json"
	"fmt"

	"github.com/xuyun-io/scalemetric/pkg/rand"
	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	"k8s.io/kubernetes/pkg/features"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// ObjectCopy copy data.
func ObjectCopy(dst, src interface{}) error {
	byts, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(byts, dst)
}

// NodeMaxScheduling return node max schedule.
func nodeMaxScheduling(node *v1.Node, nodePods []*v1.Pod, predPod *v1.Pod) (int64, []types.PredicateFailureReason, error) {
	var (
		sum int64
	)
	conditionReason, _ := CheckNodeConditionPred(node)
	if len(conditionReason) > 0 {
		return sum, conditionReason, nil
	}
	toleratesReason := PodToleratesNodeTaintsPred(node, predPod)
	if len(toleratesReason) > 0 {
		return sum, toleratesReason, nil
	}
	tempPods := make([]*v1.Pod, len(nodePods))
	tempPod := &v1.Pod{}
	copy(tempPods, nodePods)
	if err := ObjectCopy(tempPod, predPod); err != nil {
		return sum, nil, err
	}
	for {
		reasons, _ := nodeSchedulingStatus(node, tempPods, predPod)
		if len(reasons) <= 0 {
			sum = sum + 1
			tempPod.Name = rand.New()
			tempPod.Spec.NodeName = node.GetName()
			tempPods = append(tempPods, tempPod)
			continue
		}
		return sum, reasons, nil
	}
}

// NodeSchedulingStatus return node scheduling status.
func nodeSchedulingStatus(node *v1.Node, nodePods []*v1.Pod, pod *v1.Pod) ([]types.PredicateFailureReason, error) {
	return GeneralPred(node, nodePods, pod)

}

// computePodResourceRequest returns a schedulernodeinfo.Resource that covers the largest
// width in each resource dimension. Because init-containers run sequentially, we collect
// the max in each dimension iteratively. In contrast, we sum the resource vectors for
// regular containers since they run simultaneously.
//
// If Pod Overhead is specified and the feature gate is set, the resources defined for Overhead
// are added to the calculated Resource request sum
//
// Example:
//
// Pod:
//   InitContainers
//     IC1:
//       CPU: 2
//       Memory: 1G
//     IC2:
//       CPU: 2
//       Memory: 3G
//   Containers
//     C1:
//       CPU: 2
//       Memory: 1G
//     C2:
//       CPU: 1
//       Memory: 1G
//
// Result: CPU: 3, Memory: 3G
func computePodResourceRequest(pod *v1.Pod) *preFilterState {
	result := &preFilterState{}
	for _, container := range pod.Spec.Containers {
		result.Add(container.Resources.Requests)
	}

	// take max_resource(sum_pod, any_init_container)
	for _, container := range pod.Spec.InitContainers {
		result.SetMaxResource(container.Resources.Requests)
	}

	// If Overhead is being utilized, add to the total requests for the pod
	if pod.Spec.Overhead != nil && utilfeature.DefaultFeatureGate.Enabled(features.PodOverhead) {
		result.Add(pod.Spec.Overhead)
	}

	return result
}

func pred(nodeList v1.NodeList, otherPods []*v1.Pod) {
	for i := range nodeList.Items {
		node := nodeList.Items[i]
		nodeInfo := &schedulernodeinfo.NodeInfo{}
		nodeInfo.SetNode(&node)

	}

}

// Fits checks if node have enough resources to host the pod.
func Fits(pod *v1.Pod, nodeInfo *schedulernodeinfo.NodeInfo, ignoredExtendedResources sets.String) []InsufficientResource {
	return fitsRequest(computePodResourceRequest(pod), nodeInfo, ignoredExtendedResources)
}

func fitsRequest(podRequest *preFilterState, nodeInfo *schedulernodeinfo.NodeInfo, ignoredExtendedResources sets.String) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	allowedPodNumber := nodeInfo.AllowedPodNumber()
	if len(nodeInfo.Pods())+1 > allowedPodNumber {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourcePods,
			"Too many pods",
			1,
			int64(len(nodeInfo.Pods())),
			int64(allowedPodNumber),
		})
	}

	if ignoredExtendedResources == nil {
		ignoredExtendedResources = sets.NewString()
	}

	if podRequest.MilliCPU == 0 &&
		podRequest.Memory == 0 &&
		podRequest.EphemeralStorage == 0 &&
		len(podRequest.ScalarResources) == 0 {
		return insufficientResources
	}

	allocatable := nodeInfo.AllocatableResource()
	if allocatable.MilliCPU < podRequest.MilliCPU+nodeInfo.RequestedResource().MilliCPU {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceCPU,
			"Insufficient cpu",
			podRequest.MilliCPU,
			nodeInfo.RequestedResource().MilliCPU,
			allocatable.MilliCPU,
		})
	}
	if allocatable.Memory < podRequest.Memory+nodeInfo.RequestedResource().Memory {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceMemory,
			"Insufficient memory",
			podRequest.Memory,
			nodeInfo.RequestedResource().Memory,
			allocatable.Memory,
		})
	}
	if allocatable.EphemeralStorage < podRequest.EphemeralStorage+nodeInfo.RequestedResource().EphemeralStorage {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceEphemeralStorage,
			"Insufficient ephemeral-storage",
			podRequest.EphemeralStorage,
			nodeInfo.RequestedResource().EphemeralStorage,
			allocatable.EphemeralStorage,
		})
	}

	for rName, rQuant := range podRequest.ScalarResources {
		if v1helper.IsExtendedResourceName(rName) {
			// If this resource is one of the extended resources that should be
			// ignored, we will skip checking it.
			if ignoredExtendedResources.Has(string(rName)) {
				continue
			}
		}
		if allocatable.ScalarResources[rName] < rQuant+nodeInfo.RequestedResource().ScalarResources[rName] {
			insufficientResources = append(insufficientResources, InsufficientResource{
				rName,
				fmt.Sprintf("Insufficient %v", rName),
				podRequest.ScalarResources[rName],
				nodeInfo.RequestedResource().ScalarResources[rName],
				allocatable.ScalarResources[rName],
			})
		}
	}

	return insufficientResources
}

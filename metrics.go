package scale

import (
	"fmt"

	"github.com/xuyun-io/scalemetric/pkg/resources"
	"github.com/xuyun-io/scalemetric/testdata"
	v1 "k8s.io/api/core/v1"
)

func Metric(pod *v1.Pod) {
	client := testdata.KubernetesClientset()
	nodes, err := resources.GetNodes(client)
	if err != nil {
		panic(err.Error())
	}
	if len(nodes.Items) <= 0 {
		panic("nodes")
	}
	pods, err := resources.GetPods(client)
	if err != nil {
		panic(err.Error())
	}
	if len(nodes.Items) <= 0 {
		panic("pods")
	}

	schedulerStatus := PodRequestScheduling(pod, nodes, pods)
	cpu := schedulerStatus.Pod.Spec.Containers[0].Resources.Requests.Cpu()
	memory := schedulerStatus.Pod.Spec.Containers[0].Resources.Requests.Memory()
	status := fmt.Sprintf("预计集群还能调度(cpu: %v, memory: %v)Pod数: %d", cpu, memory, schedulerStatus.PredMaxschedulingCount)
	fmt.Println(status)
	fmt.Println("各节点预计未来不可调度原因: ")
	for _, nodesche := range schedulerStatus.NodeScheduling {

		// nodeSchedulingStatus := clusterStatus.SchedulingStatus[0].NodeSchedulingStatus[i]
		lastR := fmt.Sprintf("     节点: %s, 还能调度: %d, 原因: %s", nodesche.Node.GetName(), nodesche.PredMaxschedulingCount, nodesche.LastReason[0])
		fmt.Println(lastR)
	}

	// var (
	// 	count               int64
	// 	podSchedulingStatus = PodSchedulingStatus{}
	// )

	// podSchedulingStatus.Pod = pod

	// for i := range nodes.Items {
	// 	nodePods := resources.FilterNodePods(&nodes.Items[i], pods)
	// 	max, LastReason, err := calculate.NodeMaxScheduling(&nodes.Items[i], nodePods, pod)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	count = count + max

	// 	podSchedulingStatus.NodeSchedulingStatus = append(podSchedulingStatus.NodeSchedulingStatus, NodeSchedulingStatus{
	// 		NodeName:               nodes.Items[i].GetName(),
	// 		PredMaxschedulingCount: max,
	// 		LastReason:             LastReason,
	// 	})

	// }
	// podSchedulingStatus.PredMaxschedulingCount = count
	// clusterStatus := ClusterSchedulingStatus{
	// 	SchedulingStatus: []PodSchedulingStatus{podSchedulingStatus},
	// }
	// status := fmt.Sprintf("预计集群还能调度(1c2g)Pod数: %d", clusterStatus.SchedulingStatus[0].PredMaxschedulingCount)
	// fmt.Println(status)
	// fmt.Println("各节点预计未来不可调度原因: ")
	// for i := range clusterStatus.SchedulingStatus[0].NodeSchedulingStatus {
	// 	nodeSchedulingStatus := clusterStatus.SchedulingStatus[0].NodeSchedulingStatus[i]
	// 	lastR := fmt.Sprintf("     节点: %s, 还能调度: %d, 原因: %s", nodeSchedulingStatus.NodeName, nodeSchedulingStatus.PredMaxschedulingCount, nodeSchedulingStatus.LastReason[0])
	// 	fmt.Println(lastR)
	// }
	panic("finished")
}

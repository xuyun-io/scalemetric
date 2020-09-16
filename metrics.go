package main

import (
	"fmt"

	"github.com/xuyun-io/scalemetric/pkg/calculate"
	"github.com/xuyun-io/scalemetric/pkg/clientset"
	"github.com/xuyun-io/scalemetric/pkg/resources"
	v1 "k8s.io/api/core/v1"
)

// import (
// 	"fmt"

// 	"github.com/xuyun-io/scalemetric/pkg/calculate"
// 	"github.com/xuyun-io/scalemetric/pkg/resources"
// 	"github.com/xuyun-io/scalemetric/testdata"
// 	v1 "k8s.io/api/core/v1"
// )

func Metric(pod *v1.Pod) {
	client := clientset.KubernetesClientset()
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

	schedulerStatus := calculate.PodRequestScheduling(pod, nodes, pods)
	cpu := schedulerStatus.Pod.Spec.Containers[0].Resources.Requests.Cpu()
	memory := schedulerStatus.Pod.Spec.Containers[0].Resources.Requests.Memory()
	status := fmt.Sprintf("预计集群还能调度(cpu: %v, memory: %v)Pod数: %d", cpu, memory, schedulerStatus.PredMaxschedulingCount)
	fmt.Println(status)
	fmt.Println("各节点预计未来不可调度原因: ")
	for _, nodesche := range schedulerStatus.NodeScheduling {
		lastR := fmt.Sprintf("     节点: %s, 还能调度: %d, 原因: %s", nodesche.Node.GetName(), nodesche.PredMaxschedulingCount, nodesche.LastReason[0].Error())
		fmt.Println(lastR)
	}
	panic("test")
}

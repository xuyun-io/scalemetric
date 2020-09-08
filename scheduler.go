package scale

import (
	"sync"

	"github.com/xuyun-io/scalemetric/pkg/calculate"
	"github.com/xuyun-io/scalemetric/pkg/resources"
	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// func ClusterScheduling(predPods []*v1.Pod, nodeList *v1.NodeList, allrunPods *v1.PodList) types.ClusterScheduling {
// 	cs := types.ClusterScheduling{SchedulingStatus: make([]types.PodRequestScheduling, 0)}
// 	for i := range predPods {
// 		podScheduling := PodRequestScheduling(predPods[i], nodeList, allrunPods)
// 		cs.SchedulingStatus = append(cs.SchedulingStatus, podScheduling)
// 	}
// 	return cs
// }

func PodRequestScheduling(predPod *v1.Pod, nodeList *v1.NodeList, allrunPods *v1.PodList) types.PodRequestScheduling {
	var (
		syncNodeScheduling = NewSyncNodeScheduling()
		wg                 sync.WaitGroup
	)
	for _, node := range nodeList.Items {
		wg.Add(1)
		go func(pod *v1.Pod, no v1.Node, allPods *v1.PodList) {
			count, conditions, _ := NodeScheduling(pod, no, allPods)
			syncNodeScheduling.Append(count, types.NodeScheduling{
				Node:                   &no,
				PredMaxschedulingCount: count,
				LastReason:             conditions,
			})
			defer wg.Done()
		}(predPod, node, allrunPods)
	}
	wg.Wait()
	sum, status := syncNodeScheduling.Get()
	return types.PodRequestScheduling{
		Pod:                    predPod,
		PredMaxschedulingCount: sum,
		NodeScheduling:         status,
	}
}

// NodeScheduling node scheduling calculate status.
func NodeScheduling(predPod *v1.Pod, node v1.Node, allrunPods *v1.PodList) (int64, []calculate.PredicateFailureReason, error) {
	nodePods := resources.FilterNodePods(&node, allrunPods)
	return calculate.NodeMaxScheduling(&node, nodePods, predPod)
}

// type SyncInt64 struct {
// 	data int64
// 	lock sync.RWMutex
// }

// func (sInt SyncInt64) Add(num int64) {
// 	sInt.lock.Lock()
// 	defer sInt.lock.Unlock()
// 	sInt.data += num
// }

// func (sInt SyncInt64) Get() int64 {
// 	sInt.lock.RLock()
// 	defer sInt.lock.RUnlock()
// 	return sInt.data
// }

type SyncNodeScheduling struct {
	lock             sync.RWMutex
	schedulingStatus []types.NodeScheduling
	count            int64
}

func NewSyncNodeScheduling() *SyncNodeScheduling {
	return &SyncNodeScheduling{
		schedulingStatus: make([]types.NodeScheduling, 0),
		count:            0,
	}
}

func (sns *SyncNodeScheduling) Append(num int64, ns types.NodeScheduling) {
	sns.lock.Lock()
	defer sns.lock.Unlock()
	sns.schedulingStatus = append(sns.schedulingStatus, ns)
	sns.count = sns.count + num
}

func (sns *SyncNodeScheduling) Get() (int64, []types.NodeScheduling) {
	sns.lock.RLock()
	defer sns.lock.RUnlock()
	return sns.count, sns.schedulingStatus
}

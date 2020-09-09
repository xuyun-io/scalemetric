package calculate

import (
	"sync"

	"github.com/xuyun-io/scalemetric/pkg/resources"
	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
)

func PodRequestScheduling(predPod *v1.Pod, nodeList *v1.NodeList, allrunPods *v1.PodList) types.PodRequestScheduling {
	var (
		syncNodeScheduling = newSyncNodeScheduling()
		wg                 sync.WaitGroup
	)
	for _, node := range nodeList.Items {
		wg.Add(1)
		go func(pod *v1.Pod, no v1.Node, allPods *v1.PodList) {
			count, conditions, _ := nodeScheduling(pod, no, allPods)
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
func nodeScheduling(predPod *v1.Pod, node v1.Node, allrunPods *v1.PodList) (int64, []types.PredicateFailureReason, error) {
	nodePods := resources.FilterNodePods(&node, allrunPods)
	return nodeMaxScheduling(&node, nodePods, predPod)
}

type syncNodeScheduling struct {
	lock             sync.RWMutex
	schedulingStatus []types.NodeScheduling
	count            int64
}

func newSyncNodeScheduling() *syncNodeScheduling {
	return &syncNodeScheduling{
		schedulingStatus: make([]types.NodeScheduling, 0),
		count:            0,
	}
}

func (sns *syncNodeScheduling) Append(num int64, ns types.NodeScheduling) {
	sns.lock.Lock()
	defer sns.lock.Unlock()
	sns.schedulingStatus = append(sns.schedulingStatus, ns)
	sns.count = sns.count + num
}

func (sns *syncNodeScheduling) Get() (int64, []types.NodeScheduling) {
	sns.lock.RLock()
	defer sns.lock.RUnlock()
	return sns.count, sns.schedulingStatus
}

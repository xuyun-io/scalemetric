package calculate

import v1 "k8s.io/api/core/v1"

// GetActivePods returns non-terminal pods
func GetActivePods(allPods []*v1.Pod) []*v1.Pod {
	activePods := filterOutTerminatedPods(allPods)
	return activePods
}

func filterOutTerminatedPods(pods []*v1.Pod) []*v1.Pod {
	filteredPods := make([]*v1.Pod, 0)
	for _, p := range pods {
		if podIsTerminated(p) {
			continue
		}
		filteredPods = append(filteredPods, p)
	}
	return filteredPods
}

func podIsTerminated(pod *v1.Pod) bool {
	_, podWorkerTerminal := podAndContainersAreTerminal(pod)
	return podWorkerTerminal
}

func podAndContainersAreTerminal(pod *v1.Pod) (containersTerminal, podWorkerTerminal bool) {
	status := pod.Status
	containersTerminal = notRunning(status.ContainerStatuses)
	podWorkerTerminal = status.Phase == v1.PodFailed || status.Phase == v1.PodSucceeded || (pod.DeletionTimestamp != nil && containersTerminal)
	return
}

func notRunning(statuses []v1.ContainerStatus) bool {
	for _, status := range statuses {
		if status.State.Terminated == nil && status.State.Waiting == nil {
			return false
		}
	}
	return true
}

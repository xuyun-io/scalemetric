package resources

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPods return pod list.
func GetPods(client kubernetes.Interface, opts ...metav1.ListOptions) (*v1.PodList, error) {
	opt := defaultListOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), opt)
}

// GetNodePods Return the pods on the node.
func GetNodePods(client kubernetes.Interface, nodeName string) (*v1.PodList, error) {
	opt := metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName)}
	return client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), opt)
}

// GetPod return pod descirble.
func GetPod(client kubernetes.Interface, namespace, name string) (*v1.Pod, error) {
	return client.CoreV1().Pods(namespace).Get(context.TODO(), name, defaultGetOptions)
}

// FilterNodePods filter pods.
func FilterNodePods(node *v1.Node, pods *v1.PodList) []*v1.Pod {
	pos := make([]*v1.Pod, 0)
	if pods != nil && len(pods.Items) > 0 {
		for i := range pods.Items {
			if pods.Items[i].Spec.NodeName == node.Name {
				pos = append(pos, &pods.Items[i])
			}
		}
	}
	return pos
}

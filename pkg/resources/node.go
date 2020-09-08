package resources

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	defaultListOptions = metav1.ListOptions{}
	defaultGetOptions  = metav1.GetOptions{}
)

// GetNodes return node list.
func GetNodes(client kubernetes.Interface, opts ...metav1.ListOptions) (*v1.NodeList, error) {
	opt := defaultListOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return client.CoreV1().Nodes().List(context.TODO(), opt)
}

// GetNode return node describe.
func GetNode(client kubernetes.Interface, nodeName string) (*v1.Node, error) {
	return client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
}

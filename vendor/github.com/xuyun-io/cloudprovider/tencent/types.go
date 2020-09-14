package tencent

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) getNode() *v1.NodeList {
	clusters, err := c.Clusters()
	if err != nil {
		panic(err.Error())
	}
	for i := range clusters.Item {
		k8sClient, err := c.K8SClientset().Clientset(clusters.Item[i].UseUnique)
		if err != nil {
			panic(err.Error())
		}
		nodeList, err := k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		return nodeList
	}
	return nil
}

package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xuyun-io/cloudprovider"
	"github.com/xuyun-io/scalemetric/pkg/calculate"
	"github.com/xuyun-io/scalemetric/pkg/resources"
)

func main() {
	// api
	// route := route.InitRoute()
	// if err := route.Run(); err != nil {
	// 	log.Panic(fmt.Errorf("run failed, %v, exit", err.Error()))
	// }

	// lambda
	// lambdaStart()
	lambda.Start(LambdaHandler)

}

func LambdaHandler() {
	// get client
	cfg := &Config{}
	if err := cfg.Get(); err != nil {
		panic(err.Error())
	}
	client, err := cloudprovider.NewAWSProviderConfig(cfg.RegionID, cfg.AccessKey, cfg.SecretAccessKey).NewClient()
	if err != nil {
		panic(err.Error())
	}
	k8sclient, err := client.K8SClientset().Clientset(cfg.ClusterName)
	if err != nil {
		panic(err.Error())
	}
	nodeList, err := resources.GetNodes(k8sclient)
	if err != nil {
		panic(err.Error())
	}
	if len(nodeList.Items) <= 0 {
		panic(fmt.Sprintf("node is empty"))
	}
	podList, err := resources.GetPods(k8sclient)
	if err != nil {
		panic(err.Error())
	}
	pod, err := generatePod(cfg.CPURequest, cfg.MemoryRequest)
	if err != nil {
		panic(err.Error())
	}
	status := calculate.ClusterPodRequestScheduling(pod, nodeList, podList)
	metrics := ClusterSchedulingToAWSMetric(status)
	cw, err := newCloudwatchClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	output, err := PutMetrics(cw, cfg.LambdaNamespace, metrics)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("output: %v\n", output)

}

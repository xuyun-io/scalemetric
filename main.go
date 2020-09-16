package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mlycore/log"
	"github.com/xuyun-io/cloudprovider"
	"github.com/xuyun-io/scalemetric/pkg/calculate"
	pkglambda "github.com/xuyun-io/scalemetric/pkg/lambda"
	"github.com/xuyun-io/scalemetric/pkg/resources"
	"github.com/xuyun-io/scalemetric/pkg/types"
)

const (
	// EnvLogLevel declares environment variables LOGLEVEL
	EnvLogLevel = "LOGLEVEL"
)

func init() {
	log.NewDefaultLogger()
	log.SetFormatter(&log.TextFormatter{})
	if lv := os.Getenv(EnvLogLevel); !strings.EqualFold(lv, "") {
		log.SetLevel(lv)
	}
}

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

// LambdaHandler  execute lambda handler.
func LambdaHandler() {
	// get client
	cfg, err := pkglambda.Get()
	if err != nil {
		log.Fatalf("get config failed, %v", err.Error)
	}
	client, err := cloudprovider.NewAWSProviderConfig(cfg.RegionID, cfg.AccessKey, cfg.SecretAccessKey).NewClient()
	if err != nil {
		log.Fatalf("aws client create failed, %v", err)
	}
	k8sclient, err := client.K8SClientset().Clientset(cfg.ClusterName)
	if err != nil {
		log.Fatalf("kubernetes client create failed, %v", err)
	}
	nodeList, err := resources.GetNodes(k8sclient)
	if err != nil {
		log.Fatalf("node get failed, %v", err)
	}
	if len(nodeList.Items) <= 0 {
		log.Infoln("node is empty,skip calculate")
		return
	}
	podList, err := resources.GetPods(k8sclient)
	if err != nil {
		log.Fatalf("pod get failed, %v", err)
	}
	pod, err := types.GeneratePod(cfg.CPURequest, cfg.MemoryRequest)
	if err != nil {
		log.Fatalf("generate pod  failed, %v", err)
	}
	status := calculate.ClusterPodRequestScheduling(pod, nodeList, podList)
	metrics := pkglambda.ClusterSchedulingToAWSMetric(cfg, status)
	cw, err := pkglambda.NewCloudwatchClient(cfg)
	if err != nil {
		log.Fatalf("create cloudwatch client failed, %v", err)
		return
	}

	output, err := pkglambda.PutMetrics(cw, cfg.LambdaNamespace, metrics)
	if err != nil {
		log.Fatalf("push metric to cloudwatch failed, %v", err)
		return
	}
	log.Infoln(fmt.Sprintf("push metric results: %s", output.String()))

}

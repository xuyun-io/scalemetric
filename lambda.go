package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/xuyun-io/cloudprovider"
	"github.com/xuyun-io/scalemetric/pkg/calculate"
	"github.com/xuyun-io/scalemetric/pkg/resources"
	"github.com/xuyun-io/scalemetric/pkg/types"
	v1 "k8s.io/api/core/v1"
)

const (
	RegionID        = "RegionID"
	AccessKey       = "AccessKey"
	SecretAccessKey = "SecretAccessKey"
	ClusterName     = "ClusterName"
	LambdaNamespace = "LambdaNamespace"
)

type Config struct {
	RegionID        string
	AccessKey       string
	SecretAccessKey string
	ClusterName     string
	LambdaNamespace string
}

// Get get config from env
func (conf *Config) Get() error {
	cfg := &Config{
		RegionID:        os.Getenv(RegionID),
		AccessKey:       os.Getenv(AccessKey),
		SecretAccessKey: os.Getenv(SecretAccessKey),
		ClusterName:     os.Getenv(ClusterName),
		LambdaNamespace: os.Getenv(LambdaNamespace),
	}
	if len(cfg.RegionID) <= 0 {
		return errors.New("RegionID is not allowed to be emtpy")
	}
	if len(cfg.AccessKey) <= 0 {
		return errors.New("AccessKey is not allowed to be emtpy")

	}
	if len(cfg.SecretAccessKey) <= 0 {
		return errors.New("SecretAccessKey is not allowed to be empty")
	}
	if len(cfg.ClusterName) <= 0 {
		return errors.New("ClusterName is not allowed to be empty")
	}
	if len(cfg.LambdaNamespace) <= 0 {
		return errors.New("LambdaNamespace is not allowed to be empty")
	}
	return nil
}

func init() {

}

func lambdaHandler() {
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
	pod := &v1.Pod{}
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

func lambdaStart() {
	lambda.Start(lambdaHandler)
}

// NewCloudwatchClient define cloudwatch
func newCloudwatchClient(cfg *Config) (*cloudwatch.CloudWatch, error) {
	awsConf := getAWSConfig(cfg.RegionID, cfg.AccessKey, cfg.SecretAccessKey)
	ss, err := session.NewSession(awsConf)
	if err != nil {
		return nil, err
	}
	return cloudwatch.New(ss), nil
}

func getAWSConfig(region, accessKey, secretAccessKey string) *aws.Config {
	credential := credentials.NewStaticCredentials(accessKey, secretAccessKey, "")
	cfg := &aws.Config{Credentials: credential}
	cfg.WithRegion(region)
	return cfg

}

// PutMetrics return put metrics.
func PutMetrics(cw *cloudwatch.CloudWatch, namespace string, data []*cloudwatch.MetricDatum) (*cloudwatch.PutMetricDataOutput, error) {
	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(namespace),
		MetricData: data,
	}
	return cw.PutMetricData(input)
}

// ClusterSchedulingToAWSMetric return cloud watch metric data.
func ClusterSchedulingToAWSMetric(scheduling *types.ClusterScheduling) []*cloudwatch.MetricDatum {
	metrics := make([]*cloudwatch.MetricDatum, 0)
	dimension := make([]*cloudwatch.Dimension, 0)
	schedulingStatus := scheduling.SchedulingStatus[0]
	// cpu := sche.Pod.Spec.Containers[0].Resources.Requests.Cpu()
	// memory := sche.Pod.Spec.Containers[0].Resources.Requests.Memory()
	// sche.PredMaxschedulingCount
	for j := range schedulingStatus.NodeScheduling {
		nodeStatus := schedulingStatus.NodeScheduling[j]
		dimension = append(dimension, &cloudwatch.Dimension{
			Name:  aws.String(nodeStatus.Node.GetName()),
			Value: aws.String(fmt.Sprintf("%d", nodeStatus.PredMaxschedulingCount)),
		})
	}

	m := &cloudwatch.MetricDatum{
		MetricName: aws.String("clusterMaxSchedulingPodPred"),
		Unit:       aws.String("num"),
		Value:      aws.Float64(toFloat64(schedulingStatus.PredMaxschedulingCount)),
		Dimensions: dimension,
		Timestamp:  aws.Time(time.Now()),
	}
	metrics = append(metrics, m)
	return metrics
}

func toFloat64(i int64) float64 {
	in := int(i)
	return float64(in)
}

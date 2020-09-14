package lambda

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/xuyun-io/scalemetric/pkg/types"
)

const (
	RegionID        = "RegionID"
	AccessKey       = "AccessKey"
	SecretAccessKey = "SecretAccessKey"
	ClusterName     = "ClusterName"
	LambdaNamespace = "LambdaNamespace"
	CPURequest      = "CPURequest"
	MemoryRequest   = "MemoryRequest"
)

// Config define config.
type Config struct {
	RegionID        string
	AccessKey       string
	SecretAccessKey string
	ClusterName     string
	LambdaNamespace string
	MemoryRequest   string
	CPURequest      string
}

// Get get config from env
func Get() (*Config, error) {
	cfg := &Config{
		RegionID:        os.Getenv(RegionID),
		AccessKey:       os.Getenv(AccessKey),
		SecretAccessKey: os.Getenv(SecretAccessKey),
		ClusterName:     os.Getenv(ClusterName),
		LambdaNamespace: os.Getenv(LambdaNamespace),
		CPURequest:      os.Getenv(CPURequest),
		MemoryRequest:   os.Getenv(MemoryRequest),
	}
	if len(cfg.RegionID) <= 0 {
		return cfg, errors.New("RegionID is not allowed to be emtpy")
	}
	if len(cfg.AccessKey) <= 0 {
		return cfg, errors.New("AccessKey is not allowed to be emtpy")

	}
	if len(cfg.SecretAccessKey) <= 0 {
		return cfg, errors.New("SecretAccessKey is not allowed to be empty")
	}
	if len(cfg.ClusterName) <= 0 {
		return cfg, errors.New("ClusterName is not allowed to be empty")
	}
	if len(cfg.LambdaNamespace) <= 0 {
		return cfg, errors.New("LambdaNamespace is not allowed to be empty")
	}
	// resources
	if len(cfg.MemoryRequest) <= 0 {
		return cfg, errors.New("MemoryRequest is not allowed to be empty")
	}
	if len(cfg.CPURequest) <= 0 {
		return cfg, errors.New("CPURequest is not allowed to be empty")
	}
	return cfg, nil
}

// NewCloudwatchClient define cloudwatch
func NewCloudwatchClient(cfg *Config) (*cloudwatch.CloudWatch, error) {
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

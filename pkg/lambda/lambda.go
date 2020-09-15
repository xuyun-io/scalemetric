package lambda

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/xuyun-io/scalemetric/pkg/types"
)

const (
	RegionID            = "RegionID"
	AccessKey           = "AccessKey"
	SecretAccessKey     = "SecretAccessKey"
	ClusterName         = "ClusterName"
	LambdaNamespace     = "LambdaNamespace"
	CPURequest          = "CPURequest"
	MemoryRequest       = "MemoryRequest"
	AutoScalingGroupKey = "AutoScalingGroupKey"
)

// Get get config from env
func Get() (*Config, error) {
	cfg := &Config{
		RegionID:            os.Getenv(RegionID),
		AccessKey:           os.Getenv(AccessKey),
		SecretAccessKey:     os.Getenv(SecretAccessKey),
		ClusterName:         os.Getenv(ClusterName),
		LambdaNamespace:     os.Getenv(LambdaNamespace),
		CPURequest:          os.Getenv(CPURequest),
		MemoryRequest:       os.Getenv(MemoryRequest),
		AutoScalingGroupKey: os.Getenv(AutoScalingGroupKey),
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
func ClusterSchedulingToAWSMetric(cfg *Config, scheduling *types.ClusterScheduling) []*cloudwatch.MetricDatum {
	metrics := make([]*cloudwatch.MetricDatum, 0)
	dimension := make([]*cloudwatch.Dimension, 0)
	schedulingStatus := scheduling.SchedulingStatus[0]
	dimension = append(dimension, &cloudwatch.Dimension{
		Name:  aws.String("ClusterName"),
		Value: aws.String(cfg.ClusterName),
	})
	m := &cloudwatch.MetricDatum{
		MetricName: aws.String("ClusterMaxSchedulingPodPred"),
		Unit:       aws.String("Count"),
		Value:      aws.Float64(toFloat64(schedulingStatus.PredMaxschedulingCount)),
		Dimensions: dimension,
		Timestamp:  aws.Time(time.Now()),
	}
	metrics = append(metrics, m)
	groupMetrics := filterGroup(cfg.AutoScalingGroupKey, scheduling)
	if len(groupMetrics) > 0 {
		gm := AutoScalingMapToAWSMetric(cfg.ClusterName, cfg.AutoScalingGroupKey, groupMetrics)
		metrics = append(metrics, gm...)
	}
	return metrics
}

func toFloat64(i int64) float64 {
	in := int(i)
	return float64(in)
}

// AutoScalingMapToAWSMetric return auto scaling map to aws metric.
func AutoScalingMapToAWSMetric(clusterName, autoScalingGroupKey string, m AutoScalingMap) []*cloudwatch.MetricDatum {
	dimension := &cloudwatch.Dimension{
		Name:  aws.String("ClusterName"),
		Value: aws.String(clusterName),
	}
	metrics := make([]*cloudwatch.MetricDatum, 0)
	for k, v := range m {
		log.Println(fmt.Sprintf("labels %s= %s  AutoGroupMaxSchedulingPodPred: %d", autoScalingGroupKey, k, v))
		dimensions := []*cloudwatch.Dimension{dimension}
		dimensions = append(dimensions, &cloudwatch.Dimension{
			Name:  aws.String(autoScalingGroupKey),
			Value: aws.String(k),
		})
		metrics = append(metrics, &cloudwatch.MetricDatum{
			MetricName: aws.String("ClusterMaxSchedulingPodPred"),
			Unit:       aws.String("Count"),
			Value:      aws.Float64(toFloat64(v)),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: dimensions,
		})
	}
	return metrics

}

func filterGroup(key string, scheduling *types.ClusterScheduling) AutoScalingMap {
	if len(key) <= 0 {
		return nil
	}
	m := NewAutoScalingMap()
	for i := range scheduling.SchedulingStatus {
		for j := range scheduling.SchedulingStatus[i].NodeScheduling {
			nodeScheduling := scheduling.SchedulingStatus[i].NodeScheduling[j]
			labels := nodeScheduling.Node.GetLabels()
			v, ok := labels[key]
			if ok {
				m.Add(v, nodeScheduling.PredMaxschedulingCount)
			}
		}
	}
	return m
}

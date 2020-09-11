package aws

import (
	"encoding/json"
	"fmt"

	awsutil "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func ELBDescribe() {
	a := newClient("ap-northeast-1", "asdsd", "uYtYTiuasdasd")
	client, err := a.ELB()
	if err != nil {
		panic(err)
	}

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{awsutil.String("b5b41197-kubesystem-defaul-2296")},
	}
	output, err := client.DescribeLoadBalancers(input)
	if err != nil {
		panic(err.Error())
	}
	byts, err := json.Marshal(output)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(byts))
}

func ELBV2DescribeListers() {
	a := newClient("ap-northeast-1", "", "")
	client, err := a.ELBv2()
	if err != nil {
		panic(err)
	}
	// inputgi := &elbv2.DescribeTargetGroupsInput{
	// 	// LoadBalancerArn: awsutil.String("arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:loadbalancer/app/NSS-Tokyo-Gitlab-ALB/b67ff189c69aa9ae"),
	// 	// TargetGroupArns: []*string{awsutil.String("arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:targetgroup/NSS-Tokyo-Gitlab/ce52eec49d5868c5")},
	// 	Names: []*string{aws.String("b5b41197-1bc68682826ac8e62f7")},
	// }
	// output, err := client.DescribeTargetGroups(inputgi)
	input := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: awsutil.String("arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:targetgroup/b5b41197-1bc68682826ac8e62f7/7bfe1b589e249235"),
	}
	output, err := client.DescribeTargetHealth(input)

	// input := &elbv2.DescribeTargetGroupAttributesInput{
	// 	TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:targetgroup/NSS-Tokyo-Gitlab/ce52eec49d5868c5"),
	// }
	// output, err := client.DescribeTargetGroupAttributes(input)
	// client.DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput)
	// input := &elbv2.DescribeListenersInput{
	// LoadBalancerArn: awsutil.String("arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:loadbalancer/app/NSS-Tokyo-Gitlab-ALB/b67ff189c69aa9ae"),
	// }
	// output, err := client.DescribeListeners(input)
	if err != nil {
		panic(err.Error())
	}
	byts, err := json.Marshal(output)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(byts))
}

func getLB(client *elbv2.ELBV2, lbName string) (*elbv2.DescribeLoadBalancersOutput, error) {
	input := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{awsutil.String(lbName)},
	}
	return client.DescribeLoadBalancers(input)
}
func getLBListeners(client *elbv2.ELBV2, loadBalancerArn string) (*elbv2.DescribeListenersOutput, error) {
	input := &elbv2.DescribeListenersInput{LoadBalancerArn: awsutil.String(loadBalancerArn)}
	return client.DescribeListeners(input)
}

func getTargetHealth(client *elbv2.ELBV2, targetGroupArn string) (*elbv2.DescribeTargetHealthOutput, error) {
	input := &elbv2.DescribeTargetHealthInput{TargetGroupArn: awsutil.String(targetGroupArn)}
	return client.DescribeTargetHealth(input)
}

func ELBV2Describe() {
	a := newClient("ap-northeast-1", "adasd", "uYtYTiu/asds")
	client, err := a.ELBv2()
	if err != nil {
		panic(err)
	}
	fmt.Println("============DescribeLoadBalancers==================")
	lbOutput, err := getLB(client, "NSS-Tokyo-Gitlab-ALB")
	if err != nil {
		panic(err.Error())
	}
	lbOutputByts := jsonMarshal(lbOutput)
	fmt.Println(string(lbOutputByts))
	fmt.Println("============DescribeLoadBalancers finished==================")
	listeroutput, err := getLBListeners(client, "arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:loadbalancer/app/NSS-Tokyo-Gitlab-ALB/b67ff189c69aa9ae")
	if err != nil {
		panic(err.Error())
	}
	listeroutputByts := jsonMarshal(listeroutput)
	fmt.Println(string(listeroutputByts))
	fmt.Println("============getLBListeners finished==================")
	// Type: forward, TargetGroupArn
	targetsoutput, err := getTargetHealth(client, "arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:targetgroup/NSS-Tokyo-Gitlab/ce52eec49d5868c5")
	if err != nil {
		panic(err.Error())
	}
	// targetsoutput, err := getTargetHealth(client, "arn:aws:elasticloadbalancing:ap-northeast-1:847553930390:targetgroup/b5b41197-1bc68682826ac8e62f7/7bfe1b589e249235")
	// if err != nil {
	// 	panic(err.Error())
	// }

	targetsoutputByts := jsonMarshal(targetsoutput)
	fmt.Println(string(targetsoutputByts))
	fmt.Println("============getTargetHealth finished==================")

}

func jsonMarshal(data interface{}) []byte {
	byts, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
	return byts

}

func (a *Client) ELB() (*elb.ELB, error) {
	cfg := a.getAWSConfig(a.Region)
	ss, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return elb.New(ss), nil
}

func (a *Client) ELBv2() (*elbv2.ELBV2, error) {
	cfg := a.getAWSConfig(a.Region)
	ss, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return elbv2.New(ss), nil
}

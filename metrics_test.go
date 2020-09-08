package scalemetric

import (
	"testing"

	"github.com/xuyun-io/scalemetric/testdata"
	v1 "k8s.io/api/core/v1"
)

func Test_metric(t *testing.T) {
	type args struct {
		pod *v1.Pod
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				testdata.GetPod(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Metric(tt.args.pod)
		})
	}
}

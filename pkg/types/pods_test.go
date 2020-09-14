package types

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestGeneratePod(t *testing.T) {
	type args struct {
		cpuRequest    string
		memoryRequest string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.Pod
		wantErr bool
	}{
		{
			name: "test cpu ",
			args: args{
				cpuRequest:    "1",
				memoryRequest: "2G",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePod(tt.args.cpuRequest, tt.args.memoryRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePod() = %v, want %v", got, tt.want)
			}
		})
	}
}

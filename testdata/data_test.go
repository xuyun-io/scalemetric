package testdata

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func Test_getPod(t *testing.T) {
	tests := []struct {
		name string
		want *v1.Pod
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPod(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPod() = %v, want %v", got, tt.want)
			}
		})
	}
}

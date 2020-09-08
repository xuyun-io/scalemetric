package main

import (
	scale "github.com/xuyun-io/scalemetric"
	"github.com/xuyun-io/scalemetric/testdata"
)

func main() {
	scale.Metric(testdata.GetPod())
}

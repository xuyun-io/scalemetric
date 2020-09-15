package lambda

type AutoScalingMap map[string]int64

func NewAutoScalingMap() AutoScalingMap {
	return make(map[string]int64)
}

func (m AutoScalingMap) Add(key string, value int64) {
	v, ok := m[key]
	if ok {
		m[key] = v + value
		return
	}
	m[key] = value
}

func (m AutoScalingMap) Get() map[string]int64 {
	return m
}

// Config define config.
type Config struct {
	RegionID            string
	AccessKey           string
	SecretAccessKey     string
	ClusterName         string
	LambdaNamespace     string
	MemoryRequest       string
	CPURequest          string
	AutoScalingGroupKey string
}

package types

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

var podTemp string = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "labels": {
            "app": "testapp",
            "type": "Deployment"
        },
        "name": "testapp-6674965445-bbhzq",
        "namespace": "default"
    },
    "spec": {
        "containers": [
            {
                "image": "docker.io/nginx:1.9.2",
                "imagePullPolicy": "IfNotPresent",
                "name": "nginx",
                "ports": [
                    {
                        "containerPort": 80,
                        "name": "port1",
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "limits": {
                        "cpu": "%s",
                        "memory": "%s"
                    },
                    "requests": {
                        "cpu": "%s",
                        "memory": "%s"
                    }
                }
            }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler"
    }
}`

func getPodString(cpu, memory string) string {
	return fmt.Sprintf(podTemp, cpu, memory, cpu, memory)
}

// GeneratePod return pod and error.
func GeneratePod(cpuRequest, memoryRequest string) (*v1.Pod, error) {
	podString := getPodString(cpuRequest, memoryRequest)
	pod := &v1.Pod{}
	err := json.Unmarshal([]byte(podString), pod)
	return pod, err
}

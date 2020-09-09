package clientset

import (
	"flag"
	"log"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//DefaultClient is global Kubernetes rest client
var DefaultClient kubernetes.Interface
var once sync.Once

// KubernetesClientset return local clientset.
func KubernetesClientset() kubernetes.Interface {
	once.Do(func() {
		var kubeconfig *string
		var err error
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("read config from flags error: %s", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("init client with config error: %s", err)
		}
		DefaultClient = clientset
	})
	return DefaultClient
}

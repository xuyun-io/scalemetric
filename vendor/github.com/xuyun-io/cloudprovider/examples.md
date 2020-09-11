
### 示例


#### 对接AWS的EKS和ECR
```go
func ExampleAWSProvider() {
	//	create provider config
	client, err := NewAWSProviderConfig("regionID", "accessKey", "secretAccessKey").NewClient()
	//  or 	client, err := NewProvider(NewAWSProviderConfig("regionID", "accessKey", "secretAccessKey"))
	if err != nil {
		panic(err)
	}v
	// use provider cluster client
	clusters, err := client.Clusterset().Clusters()
	if err != nil {
		panic(err)
	}
	fmt.Println(clusters)
	// use kubernetes client
	k8sClient, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(nodes)
	// use aws repo clientset
	// if u want to use other repo, u can: (eg: use docker registry)
	// tags,err:=	client.RepoClientset(NewDockerRegistry("dockerhub.io")).TagList(context.Background(), "busybox")
	tags, err := client.RepoClientset().TagList(context.Background(), "busybox")
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}
```
#### 对接GCP的GKE和CSR
```go
func ExampleGCPProvider() {
	client, err := NewGCPProviderConfig("zoneID", gcp.Config{}).NewClient()
	if err != nil {
		panic(err)
	}
	// use cluster
	describeCluster, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	fmt.Println(describeCluster)
	// use kubernetes client
	k8sClient, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(nodes)
	// use aws repo clientset
	// if u want to use other repo, u can: (eg: use docker registry)
	// tags,err:=	client.RepoClientset(NewDockerRegistry("dockerhub.io")).TagList(context.Background(), "busybox")
	tags, err := client.RepoClientset().TagList(context.Background(), "gcr.io/busybox")
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}

```

#### 对接Kubeconfig
```go
func ExampleKubeconfig() {
    client, err := NewKubeConfigProviderConfig("kubeconfig").NewClient()
    if err != nil {
		panic(err)
	}
	// use cluster
	describeCluster, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	fmt.Println(describeCluster)
	// use kubernetes client
	k8sClient, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(nodes)
	// use aws repo clientset
	// if u want to use other repo, u can: (eg: use docker registry)
	// tags,err:=	client.RepoClientset(NewDockerRegistry("dockerhub.io")).TagList(context.Background(), "busybox")
	// Warning: default use dockerhub repo client
	tags, err := client.RepoClientset().TagList(context.Background(), "busybox")
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}
```
#### 对接ALi
```go
func ExampleAli() {
	client, err := NewAWSProviderConfig("regionID", "AccessKey", "SecretAccessKey").NewClient()
	if err != nil {
		panic(err)
	}
	// use cluster
	describeCluster, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	fmt.Println(describeCluster)
	// use kubernetes client
	k8sClient, err := client.K8SClientset().Clientset("cluster")
	if err != nil {
		panic(err)
	}
	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(nodes)
	// use aws repo clientset
	// if u want to use other repo, u can: (eg: use docker registry)
    // tags,err:=	client.RepoClientset(NewDockerRegistry("dockerhub.io")).TagList(context.Background(), "busybox")
    // Warning: ali repository 
	tags, err := client.RepoClientset().TagList(context.Background(), "busybox")
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}

```

#### 对接 华为云CCE
```go
func TestNewHuaweiCCEProviderConfig(t *testing.T) {
	type args struct {
		regionID string
		conf     huawei.Auth
	}
	tests := []struct {
		name string
		args args
		want types.ProviderConfig
	}{
		{
			name: "huewei cce",
			args: args{
				regionID: defaultHuaweiCCEAuth.Region,
				conf:     defaultHuaweiCCEAuth,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewHuaweiCCEProviderConfig(tt.args.regionID, tt.args.conf).NewClient()
			if err != nil {
				t.Error(err)
			}
			clusterList, err := client.Clusterset().Clusters()
			if err != nil {
				t.Error(err)
			}
			for i := range clusterList.Item {
				k8sClient, err := client.K8SClientset().Clientset(clusterList.Item[i].UseUnique)
				if err != nil {
					t.Error(err)
				}
				cm, err := k8sClient.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
				if err != nil {
					t.Error(err)
				}
				t.Log(len(cm.Items))
			}
		})
	}
}

```

### 单独对接各云厂商的镜像库
```go
func ExampleRegisters() {
	// 1. use docker register
	// repoInterface:= NewDockerRegistry("")
	// 2. use aws ecr
	// repoInterface := NewAWSRepoClient("regionID", "accessKey", "secretAccessKey")
	// 3. use gcp csr
	repoInterface := NewGCPRepoClient(context.Background(), "gcr.io", gcp.Config{})
	tags, err := repoInterface.TagList(context.Background(), "gcr.io/busybox")
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}
```
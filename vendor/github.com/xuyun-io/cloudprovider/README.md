# cloudprovider


CloudProvider 是一个多云链接库,支持主流的公有云和私有云

### 安装
```
go get -u github.com/xuyun-io/cloudprovider
```

### 当前状态

| 云提供商 | Kubernetes 服务 | 镜像仓库服务 | 使用文档 |
|--|--|--|--|
| AWS(amazon web Services) | EKS | ECR | [权限配置](https://github.com/xuyun-io/cloudprovider#aws-%E8%AE%A4%E8%AF%81) |
| GoogleCloud(谷歌云) | GKE | ContainerRegistry | [权限配置](https://github.com/xuyun-io/cloudprovider#gcp-%E8%AE%A4%E8%AF%81) |
| Aliyun(阿里云) |  ACK | ACR(并非ACREE) | [权限配置](https://github.com/xuyun-io/cloudprovider#aliyun-%E8%AE%A4%E8%AF%81) |
| Tencent(腾讯云) | TKE | TCR | [权限配置](https://github.com/xuyun-io/cloudprovider#tencent-%E8%AE%A4%E8%AF%81) |
| Huaweicloud(华为云) | CCE | SWR | [权限配置](https://github.com/xuyun-io/cloudprovider#huaweicloud-%E8%AE%A4%E8%AF%81) |
| 私有云, 自建云| kubeconfig | docker,harbor | [权限配置](https://github.com/xuyun-io/cloudprovider#%E7%A7%81%E6%9C%89%E4%BA%91) |

### AWS 认证

认证配置:
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
|Region | string | 地区 |required | 例如: ap-northeast-1 |
|AccessKey | string | AK |required | - |
|SecretAccessKey | string | SK |required | - |

#### 操作EKS所需权限:
1. 对于 EKS 集群所需的最小权限为:
	```json
	{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "VisualEditor0",
				"Effect": "Allow",
				"Action": [
					"eks:ListFargateProfiles",
					"eks:DescribeNodegroup",
					"eks:ListNodegroups",
					"eks:DescribeFargateProfile",
					"eks:ListTagsForResource",
					"eks:ListUpdates",
					"eks:DescribeUpdate",
					"eks:DescribeCluster",
					"eks:ListClusters"
				],
				"Resource": "*"
			}
		]
	}
	```

2. 为对应的AK/SK用户配置 EKS API 访问权限

**注意: 非集群创建者需要配置，集群创建者跳过此操作**
原因:kubestar 调用 kubernetes 的API时，需要EKS API的访问权限。否则会报错: error: You must be logged in to the server (Unauthorized)。

将 user 添加到 configmap(kubectl get configmap -n kube-system aws-auth -o yaml)

```yaml
mapUsers:
 ----
 - userarn: arn:aws:iam::221483597:user/###
   username: admin
   groups:
     - system:masters
```

更多详情查看 AWS 官方文档: https://amazonaws-china.com/cn/premiumsupport/knowledge-center/eks-api-server-unauthorized-error/
#### 操作ECR 配置权限
1.![avatar](/images/aws/step3.png)

### GCP 认证

暂无说明
### Aliyun 认证
认证字段:
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
| RegionID | string | 地区| required | - |
|AccessKey | string | AK |  required | - |
|SecretAccessKey | string | SK |  required | - |
| Host | string | 镜像仓库地址 | 如需操作acr,required | 示例: registry.cn-hangzhou.aliyuncs.com, 格式: registry.{region}.aliyuncs.com|
|RepoNamespace | string |镜像仓库的namespace |如需操作acr,required | - |
|PullImageUser| string | 拉取镜像的账号 | 如需操作acr,required |生成secret,用户名为阿里云登陆账户名|
|PullImagePassword| string | 拉取镜像的密码 | 如需操作acr,required | 生成secret, 密码需要在`容器镜像服务`控制台设置,路径为: 容器镜像服务 -> 默认实例-> 访问凭证 -> 设置固定密码|

#### 权限配置
控制台路径: RAM 访问控制 -> 人员管理 -> 添加权限 -> 系统策略

#### 示例权限配置
```
1. AliyunCSFullAccess (可选)
2. AliyunContainerRegistryFullAccess
3. AliyunContainerRegistryReadOnlyAccess
```

### Tencent 认证

认证字段:
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
|Region | string | 地区 | required|例如: ap-beijing ,参考:[源码](https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/common/regions/regions.go) |
| AccessKey | string | AK | required | -|
|SecretAccessKey| string | SK | required | -|
| RegistryID | string | TCR 所需的实例id | 如需操作TCR,required | -|
| Namespace | string | TCR 组织和命令空间 | 如需操作TCR,required | -|
| Host | string | 镜像仓库host,格式为: {registryName}.tencentcloudcr.com, 注意区分registryName和RegistryID|如需操作TCR,required | -|
| Username | string | 使用docker 命令登陆镜像仓库账号 | optional | 计划快捷为kubernetes 生成pullimagesecret|
| Password | string |使用docker 命令登陆镜像仓库密码 | optional | 计划快捷为kubernetes 生成pullimagesecret|

#### 操作 TKE 所需权限
1. 如果TKE 需要外网或者内网访问,请到对应的TKE设置集群APIServer信息, 位置: 控制台-> 容器服务 -> 集群 -> 点击某集群 -> 基本信息 ->  集群APIServer信息 -> 内网访问/外网访问开启
2. 为AK/SK 配置合适的TKE集群操作权限 https://cloud.tencent.com/document/product/457/46033
3. 测试时为AK/SK 账号配置了`QcloudTKEFullAccess`,`QcloudAccessForTKERole`,`QcloudAccessForTKERoleInCreatingCFSStorageclass`,`QcloudAccessForTKERoleInOpsManagement`,注意: 权限可能过大,仅供参考。

#### 操作 TCR 所需权限
1. 测试是为 AK/SK 账号配置了`QcloudTCRFullAccess`,` 	QcloudAccessForTCRRole`,注意: 权限可能过大,仅供参考。

### Huaweicloud 认证
认证配置:
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
|AuthType | AuthType | 认证类型 | 分为CCE认证,SWR认证和所有配置, 默认情况下会检查所以认证是否正确| 支持三种值:`AuthCCE`,`AuthSWR`,`AuthALL`| 
|Region | string | 地区 | required|基础认证,例如: cn-north-1 |
|Account | string | 子账号 | required | 基础认证 |
|Password | string | 密码 | required | 基础认证 |
|DomainName | string|  父账号,标记所属关系| required |基础认证|
|ProjectID| string| CCE 所需的项目ID | 如果需要操作CCE,required| CCE 认证配置下|
|ProjectName | string| CCE 的项目名字,和 Region 接近|  如需操作CCE,required| CCE 认证配置下|
|Namespace | string|指定SWR的组织或者命名空间| 如需操作SWR,required| SWR 认证配置下|
|ProjectName | string | 指定SWR的项目名字, 和 Region 接近| 如需操作SWR,required| SWR 认证配置下 |
|AccessKey | sting | AK | 如需操作SWR,required| SWR 认证配置下|
| SecretAccessKey | string | SK | 如需操作SWR,required| SWR 认证配置下|
#### 操作 CCE 所需权限
1. 华为云的CCE是使用的子账号认证(Token 认证), 官方文档: https://support.huaweicloud.com/api-iam/iam_30_0001.html
2. 如遇 `Unauthorized`,访问: https://support.huaweicloud.com/usermanual-cce/cce_01_0188.html
3. 如遇 `Timeout(默认30s超时)`, 请遵照执行:如需通过公网使用client-go/kubectl,需要为集群设置公网弹性IP, 设置路径: 华为云控制台-> 进入容器云引擎-> 资源管理-> 集群管理-> 命令行管理 -> kubectl -> 互联网访问-> 绑定弹性ip, 会生成连接信息:`https://ip:5443`

#### 操作 SWR 所需权限
1. 华为云的SWR是使用 Token 认证 + AK/SK认证, Token 认证用于接口调用, AK/SK认证 用户生成pullimagesecret, 此库需要 openssl,od,sed 命令行支持。

***<g-emoji class="g-emoji" alias="warning" fallback-src="https://github.githubassets.com/images/icons/emoji/unicode/26a0.png">⚠️</g-emoji> 注意:***


### 私有云
认证配置:
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
| RegionID | string | 区域 | required | 用户可自定义填写,可以分机房, 机架,地区 填写|
|KubeConfigString | string | kubernetes kubeconfig | required | kubernetes yaml 配置 |

#### 私有镜像库配置(docker,harbor)
| 字段 | 类型 | 含义 | 必填 | 备注 |
|--|--|--|--|--|
|Host | string | 镜像仓库地址|  required|, 如果是harbor的话,需要使用admin|
|Username | string | 镜像仓库访问账号|optional|-|
|Password | string | 镜像仓库密码|optional|-|

示例代码:
```go
auth := registry.NewUsernamePasswordAuth(username, password)
client := NewDockerRegistry(host, auth)
// or 
// client := NewDockerRegistry(host)
tags, _ := client.TagList(context.Background(), "busyboxy")
```

### 关于使用
1. 每个云都有独立的认证配置, 认证配置通常分为两部分, 一部分是 kubernetes 服务的认证, 一部分是镜像仓库的认证。不同的云认证方式不同, 
	1. 有些云使用同一套AK/SK即可对接两个服务, 如AWS;
	2. 有些云对镜像仓库调用时,还需要额外的认证配置,如华为云, kubernetes 服务认证使用的Token认证, 而镜像仓库服务还需要使用AK/SK 认证。

2. 如何快速对接 kubernetes 云服务。
	1. 创建对应云的配置`NewXXXProviderConfig()`, 然后生成对应的云服务商的 Client, 然后操作对应的kubernetes client 或者镜像仓库。
	```go
	 	以AWS 云为示例:
		//	create provider config
		client, err := cloudprovider.NewAWSProviderConfig("regionID", "accessKey", "secretAccessKey").NewClient()
		if err != nil {
			panic(err)
		}
		// use kubernetes client
		k8sClient, err := client.K8SClientset().Clientset("clustername")
		if err != nil {
			panic(err)
		}
		nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		// 快速使用镜像仓库
		// use image client
		tags, err := client.RepoClientset().TagList(context.Background(), "busybox")
		if err != nil {
		panic(err)
	}
	```
3. 如果只想对应镜像仓库服务而不使用 kubernetes 服务, 则可以直接创建对应的 RepoClient,代码类似:`cloudprovider.NewXXXRepoClient()`
	```go
		以AWS 云为示例:
		repoClient := cloudprovider.NewAWSRepoClient(regionID, accessKey, secretAccessKey)
		// image list
		repositories, err := repoClient.Repositories(context.Background())
		if err != nill {
			panic(err.Error())
		}
		// tag list
		tags, err := repoClient.TagList(context.Background(), "imagename")
		if err != nil {
			panic(err.Error())
		}
		// 通常镜像详情接口并未实现
		desc, err := repoClient.ImageDescribe(context.Background(), "imagename")
		if err != nil {
			panic(err.Error())
		}
	```

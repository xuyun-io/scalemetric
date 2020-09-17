# scalemetric


### For lambda

 如果需要把程序发布到 lambda 并把指标推送到cloudwatch里, 则需要进行以下流程。

1. 打包符合lambda 标准的zip 包
```
 执行: make lambda
```
2. 到AWS lambda 平台上传程序包
  1. 可通过zip包
  2. 可通过s3
3. 在AWS 控制台修改 lambda 的基本设置
```
处理程序: main
运行时: Go 1.x
```
4. 为程序设置环境变量

| 字段名 | 类别 | 含义 | 示例|
|-|-|-|-|
|AccessKey | string | AK,aws 访问eks 集群和cloudwatch 的权限, 其中EKS 得有读取权限, cloudwatch 得有推送权限| -|
|SecretAccessKey | string |SK,aws 访问eks 集群和cloudwatch 的权限, 其中EKS 得有读取权限, cloudwatch 得有推送权限| - |
|AutoScalingGroupKey | string | kuberenetes node 上的标签key, 做计算时, 会以此标签做分组, 只能填写一个key, 可不填, 不填意味着不做分组处理。| kubestar.io/autoscaling-group-name|
|CPURequest | string | pod request 的cpu量 | 1|
|MemoryRequest | string | pod requeset 的 memory 量|1G|
|ClusterName  |string | eks 集群名字, 只计算指定集群的预计容量| prod|
|LambdaNamespace | string | 把指标推送到指定cloudwatch 的namesapce 下| EKSScalemetric|
|RegionID | string| 集群和cloudwatch 所在region  | ap-northeast-1

注意: 未做可不填说明，则意味着是必须填写的。

5. 测试 lambda 是否好用。
    当配置完成后,在aws lambda 页面点击测试，检查lambda 是否可以正常运行。
6. 测试完成后, 进入cloudwatch页面检查在指定 namespace下是否存在指定指标, 若存在则部署完成。



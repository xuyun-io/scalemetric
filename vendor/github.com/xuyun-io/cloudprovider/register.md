# 镜像仓库

## ALi ACR

### 必要的参数
```

region 
accessKey
secretAccessKey
host
repoNamespace
apiPassword
pullImageUser
pullImagePassword
```

| 字段| 类型| 含义| 注释|
|--|--|--|--|
| region| string | 可用区| 比如: cn-hangzhou|
|accessKey | string | - |
| secretAccessKey | string |- |
|host | string | 拉取镜像时的host | 格式: registry.{region}.aliyuncs.com|
| repoNamespace | string| 阿里云规定的镜像namespace|
|apiPassword | string | 访问阿里云 CR API使用的密码| 可以随意创建，8位及以上，包含数字，小写字母，大写字母|
| pullImageUser | string | kubelet拉取镜像所需的用户名| 生成secret,用户名为阿里云登陆账户名|
|pullImagePassword | string| kubelet 拉取镜像所需的密码|  生成secret, 密码需要在`容器镜像服务`控制台设置,路径为: 容器镜像服务 -> 默认实例-> 访问凭证 -> 设置固定密码|


### 权限分配
***控制台路径***
```
RAM 访问控制 -> 人员管理 -> 添加权限 -> 系统策略
```
***分配以下权限***

```
1. AliyunCSFullAccess (可选)
2. AliyunContainerRegistryFullAccess
3. AliyunContainerRegistryReadOnlyAccess
```




package testdata

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
)

func GetNode() *v1.Node {
	node := &v1.Node{}
	if err := json.Unmarshal([]byte(testNode), node); err != nil {
		panic(err.Error())
	}
	return node
}

var testNode = `{
    "apiVersion": "v1",
    "kind": "Node",
    "metadata": {
        "annotations": {
            "node.alpha.kubernetes.io/ttl": "0",
            "volumes.kubernetes.io/controller-managed-attach-detach": "true"
        },
        "creationTimestamp": "2019-10-08T08:38:40Z",
        "labels": {
            "alpha.eksctl.io/cluster-name": "prod",
            "alpha.eksctl.io/instance-id": "i-0e8f2a05bd33ecaea",
            "alpha.eksctl.io/nodegroup-name": "ng-1-addons",
            "beta.kubernetes.io/arch": "amd64",
            "beta.kubernetes.io/instance-type": "m5.4xlarge",
            "beta.kubernetes.io/os": "linux",
            "failure-domain.beta.kubernetes.io/region": "ap-northeast-1",
            "failure-domain.beta.kubernetes.io/zone": "ap-northeast-1c",
            "kubernetes.io/arch": "amd64",
            "kubernetes.io/hostname": "ip-192-168-102-67.ap-northeast-1.compute.internal",
            "kubernetes.io/os": "linux",
            "kubernetes.io/role": "ng-1-addons",
            "node-role.kubernetes.io": "ng-1-addons",
            "nodegroup": "ng-1-addons",
            "role": "addons"
        },
        "name": "ip-192-168-102-67.ap-northeast-1.compute.internal",
        "resourceVersion": "113429429",
        "selfLink": "/api/v1/nodes/ip-192-168-102-67.ap-northeast-1.compute.internal",
        "uid": "074ecf3b-e9a7-11e9-8fee-0e393284b1da"
    },
    "spec": {
        "providerID": "aws:///ap-northeast-1c/i-0e8f2a05bd33ecaea"
    },
    "status": {
        "addresses": [
            {
                "address": "192.168.102.67",
                "type": "InternalIP"
            },
            {
                "address": "ip-192-168-102-67.ap-northeast-1.compute.internal",
                "type": "Hostname"
            },
            {
                "address": "ip-192-168-102-67.ap-northeast-1.compute.internal",
                "type": "InternalDNS"
            }
        ],
        "allocatable": {
            "attachable-volumes-aws-ebs": "25",
            "cpu": "16",
            "ephemeral-storage": "19316009748",
            "hugepages-1Gi": "0",
            "hugepages-2Mi": "0",
            "memory": "64359040Ki",
            "pods": "234"
        },
        "capacity": {
            "attachable-volumes-aws-ebs": "25",
            "cpu": "16",
            "ephemeral-storage": "20959212Ki",
            "hugepages-1Gi": "0",
            "hugepages-2Mi": "0",
            "memory": "64461440Ki",
            "pods": "234"
        },
        "conditions": [
            {
                "lastHeartbeatTime": "2020-08-19T08:18:09Z",
                "lastTransitionTime": "2019-10-08T08:38:40Z",
                "message": "kubelet has sufficient memory available",
                "reason": "KubeletHasSufficientMemory",
                "status": "False",
                "type": "MemoryPressure"
            },
            {
                "lastHeartbeatTime": "2020-08-19T08:18:09Z",
                "lastTransitionTime": "2020-06-24T02:47:05Z",
                "message": "kubelet has no disk pressure",
                "reason": "KubeletHasNoDiskPressure",
                "status": "False",
                "type": "DiskPressure"
            },
            {
                "lastHeartbeatTime": "2020-08-19T08:18:09Z",
                "lastTransitionTime": "2019-10-08T08:38:40Z",
                "message": "kubelet has sufficient PID available",
                "reason": "KubeletHasSufficientPID",
                "status": "False",
                "type": "PIDPressure"
            },
            {
                "lastHeartbeatTime": "2020-08-19T08:18:09Z",
                "lastTransitionTime": "2019-10-08T08:39:20Z",
                "message": "kubelet is posting ready status",
                "reason": "KubeletReady",
                "status": "True",
                "type": "Ready"
            }
        ],
        "daemonEndpoints": {
            "kubeletEndpoint": {
                "Port": 10250
            }
        },
        "images": [
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/angelia-backend@sha256:812729fe875971361ea0e241e6725c84e77ecd88a5f6bd49706a355a0de7fe62",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/angelia-backend:v3.2"
                ],
                "sizeBytes": 1462813587
            },
            {
                "names": [
                    "couchbase@sha256:44170af7930af029e6ecc227582a579b927161075fc1ae0458272078905f1ce5",
                    "couchbase:enterprise-6.5.1"
                ],
                "sizeBytes": 1130252189
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/operation-center@sha256:54a0b6fc7373caf93ac7a0abfee3e92aece0dea53a718bfa8f960dc0e09a5f17",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/operation-center:20200619083920-e077d18c"
                ],
                "sizeBytes": 1114218108
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/customer-portal@sha256:40a772c2d479caa354f65ca67cc460f975c3f58043ee7d0b4a018c9fb16aaee6",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/customer-portal:20200202161748-981d4b8d"
                ],
                "sizeBytes": 1101733535
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/payment-center@sha256:9d2eadecaa19911a4f6d85e1e771c7abd76a7e1e4aef4574bb44441636e63b91",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/payment-center:20200212201559-6c3e2b2d"
                ],
                "sizeBytes": 1091312749
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/mobile-front@sha256:56914240ba637802d5a313f681e2ba0b13a02f2f60a361923eea98e7b65887e0",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/xuyun/mobile-front:20200202150158-5d715b7e"
                ],
                "sizeBytes": 1076209060
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:636e4516b33753ffd33c447d40e1c485bfaeebfef6fbf981772aa24ab221f8bb",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202008121932-14f411ef"
                ],
                "sizeBytes": 827931150
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:fe8bc0120c7097a419df56064556663dc94a2a9555e80efdba4549e1ade3a8f9",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202008122013-14f411ef"
                ],
                "sizeBytes": 827931150
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlemasterjob@sha256:7764734bdf994222739a1acc3fc005722a86460e5e373b3437f038a7adafa847",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlemasterjob:202008121841-934071a3"
                ],
                "sizeBytes": 827726431
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/hotrod/hotrod-github@sha256:823b0a98cb6448cbdf82607736778155c1ec6989c1b19a2300e0be12251a8162",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/hotrod/hotrod-github:v0.1.0"
                ],
                "sizeBytes": 826579144
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:2012f1cb03cbfde1d6c34d645d7baaf108fc819c1918cfdc9226436ebb15f2b7",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202005261152"
                ],
                "sizeBytes": 825185174
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:6b93039c7dd454fbf689c84aa9278f9efdfa121a2d88b6f97d2182e2b42a91ee",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202007071210"
                ],
                "sizeBytes": 825168965
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:9e7dd5ee1df7793e571933d40e8c41a50d59bb3ac8e162fcc3faea7406fad065",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202007200620"
                ],
                "sizeBytes": 825168965
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob@sha256:026f188dfe75c584b2adeb667807038d234d0580d3cabb5bab423270ce956de6",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlesubjob:202007221808"
                ],
                "sizeBytes": 825168965
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlemasterjob@sha256:4c5d7782f8599dee09676dfc43b0759471d0253eb49e1f480d3f2a9aa4cf5ce4",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/settlemasterjob:202007222102-ce363340"
                ],
                "sizeBytes": 824981123
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/domaincer-monitor@sha256:d7b88ff30ec986251b4402ebf1e6d80b915ba655380b71ca3ace90434610d0a4",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/domaincer-monitor:202008051553"
                ],
                "sizeBytes": 784083204
            },
            {
                "names": [
                    "datadog/agent@sha256:daef39d359d12cb21b13c3bcf94d1b2cf87caeb17c657e497306df73a321488c",
                    "datadog/agent:7.21.1"
                ],
                "sizeBytes": 695360760
            },
            {
                "names": [
                    "netdata/netdata@sha256:3b27efcdd2d20fc2f86287f5711ba00b165179fca97aecdeca0ec05ce9543f54",
                    "netdata/netdata:v1.21.0"
                ],
                "sizeBytes": 464525164
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/loicsharma/baget@sha256:43f2bc3b045891fb678f57c8c99b0869a8a293bede5307c5dfd719aa030a194c",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/loicsharma/baget:latest"
                ],
                "sizeBytes": 296390147
            },
            {
                "names": [
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/amazon-k8s-cni@sha256:c071dfc45cd957fc6ab2db769ae6374b1f59a08db90b0ff0b9166b8531497a35",
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/amazon-k8s-cni:v1.5.3"
                ],
                "sizeBytes": 290731139
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/dev/starpay/wechatpay@sha256:41922f303771e953fbe8c649eef672156664fb87e106b44d8e095183a49ff687",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/dev/starpay/wechatpay:20191129090818-8b836521"
                ],
                "sizeBytes": 287384169
            },
            {
                "names": [
                    "ccr.ccs.tencentyun.com/mirrors/cadvisor@sha256:46d4d730ef886aaece9e0a65a912564cab0303cf88718d82b3df84d3add6885c",
                    "ccr.ccs.tencentyun.com/mirrors/cadvisor:v0.34.0"
                ],
                "sizeBytes": 185393665
            },
            {
                "names": [
                    "kubestar/namespace@sha256:ac68741c8e2e893e0958d9c5d9902f8434bd8405a08c50a998c62d6c6e2f5b22",
                    "kubestar/namespace:0.0.3"
                ],
                "sizeBytes": 182328314
            },
            {
                "names": [
                    "kubestar/namespace@sha256:c07f52c6da189ae091814721c419f9014edc4586ea2c8cb78bebd9547a39dc06",
                    "kubestar/namespace:0.0.1"
                ],
                "sizeBytes": 182306187
            },
            {
                "names": [
                    "gcr.io/kubecost1/checks@sha256:89f0f0147f144df586c37e8735f3bf7e5ca70e0ef221fef94e273512a4c6e829",
                    "gcr.io/kubecost1/checks:prod-1.61.3"
                ],
                "sizeBytes": 148898443
            },
            {
                "names": [
                    "nginx@sha256:c56e8eaec9ff118688c857e263f744039f5ba9fe3b30efcaeef5342f27534a01",
                    "nginx:1.9.9"
                ],
                "sizeBytes": 133864957
            },
            {
                "names": [
                    "nginx@sha256:3b50ebc3ae6fb29b713a708d4dc5c15f4223bde18ddbf3c8730b228093788a3c",
                    "nginx:1.9.7"
                ],
                "sizeBytes": 132773214
            },
            {
                "names": [
                    "nginx@sha256:a93c8a0b0974c967aebe868a186e5c205f4d3bcb5423a56559f2f9599074bbcd",
                    "nginx:latest"
                ],
                "sizeBytes": 132484492
            },
            {
                "names": [
                    "k8s.gcr.io/kubernetes-dashboard-amd64@sha256:0ae6b69432e78069c5ce2bcde0fe409c5c4d6f0f4d9cd50a17974fea38898747",
                    "k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1"
                ],
                "sizeBytes": 121711221
            },
            {
                "names": [
                    "kubesphere/prometheus@sha256:60c989c93c8097ef7719c1b3b0f4dc54ea61b5e836c222258a5d9512fb3e6181",
                    "kubesphere/prometheus:v2.5.0"
                ],
                "sizeBytes": 99822269
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kvstore/redis@sha256:e73ef998c22f9a98793d9951bb2915cd945d8fa6f9ec1b324e85d19617efc2fd",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kvstore/redis:5.0.7"
                ],
                "sizeBytes": 98204589
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/webserver/nginx@sha256:b1f5935eb2e9e2ae89c0b3e2e148c19068d91ca502e857052f14db230443e4c2",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/webserver/nginx:1.7.9"
                ],
                "sizeBytes": 91664166
            },
            {
                "names": [
                    "mackerel/mackerel-container-agent@sha256:b17fec83ed099b8df5396f1cc27776dc10988056c912eebf4f89eaaca11ae4a8",
                    "mackerel/mackerel-container-agent:latest"
                ],
                "sizeBytes": 90718317
            },
            {
                "names": [
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/eks/kube-proxy@sha256:d3a6122f63202665aa50f3c08644ef504dbe56c76a1e0ab05f8e296328f3a6b4",
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/eks/kube-proxy:v1.14.6"
                ],
                "sizeBytes": 82044796
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:0efc3d0131144c6d8582b7b44687fc45045e6a6de6555f9d3d1d5a2f6afeb60e",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200811-97d9e53c"
                ],
                "sizeBytes": 80868673
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:310bb606a67f5fcfdfb62471e7ac074f287f9a8bd0e717050c350badb338f324",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200804-3591a522"
                ],
                "sizeBytes": 80760365
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:fd6233319cb384eb2af536954e54239438294e127b2b28128925c39c09630a25",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200804-7702290d"
                ],
                "sizeBytes": 80757163
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:c827975fbc6c70fc1eee7c50180b7f39bc0682f9e53a53de9f7e4dd0c56b940b",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200803-809ec804"
                ],
                "sizeBytes": 80661555
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:b354df464fa9f3c2ab0ee54acb1e0fdf49cb8357b17bad1306f802be05ac8d5f",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200803-132e92e4"
                ],
                "sizeBytes": 80656095
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:eab2cb217826670e473bd50ca02fd253dd13d9d387d327be1d175784fbfc3800",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200803-b2400f02"
                ],
                "sizeBytes": 80653915
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:fa7a3a348b3456fe1d95ffb6a16299d2967fe96f92afb465b48f7d0c8de26b1e",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200731100406-00ab047f-ts"
                ],
                "sizeBytes": 80643413
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:1f140b97cdb46923f6f391db335cad6eb6088c983cdcda5bdd3407a421e3480f",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:20200611-f546eaac"
                ],
                "sizeBytes": 78405114
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar@sha256:5be9d36725a1b38cd6ce7aee57651d1cca45d8f16405c37adfe802b8ae46d8dc",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/kubestar:v0.7.0-alpha.1"
                ],
                "sizeBytes": 78395784
            },
            {
                "names": [
                    "ubuntu@sha256:747d2dbbaaee995098c9792d99bd333c6783ce56150d1b11e333bbceed5c54d7",
                    "ubuntu:latest"
                ],
                "sizeBytes": 73852122
            },
            {
                "names": [
                    "894847497797.dkr.ecr.us-west-2.amazonaws.com/aws-alb-ingress-controller@sha256:3c03aaed555bd64b9097d587d8c438c2e2cb032adbb899457362f7ba34c8de19",
                    "894847497797.dkr.ecr.us-west-2.amazonaws.com/aws-alb-ingress-controller:v1.0.0"
                ],
                "sizeBytes": 38286474
            },
            {
                "names": [
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/eks/coredns@sha256:c85954b828a5627b9f3c4540893ab9d8a4be5f8da7513882ad122e08f5c2e60a",
                    "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/eks/coredns:v1.3.1"
                ],
                "sizeBytes": 35174083
            },
            {
                "names": [
                    "kubestar/git@sha256:8600bf1b6124ad32a2ea477828632591f977b68952dc3027eb9f942ee2925d2d",
                    "kubestar/git:test5"
                ],
                "sizeBytes": 27807249
            },
            {
                "names": [
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/jaegertracing/jaeger-agent@sha256:63a55ab001818259b085a69477e1c66b4b35bbb78cbed030223797b71ad8d664",
                    "jaegertracing/jaeger-agent@sha256:63a55ab001818259b085a69477e1c66b4b35bbb78cbed030223797b71ad8d664",
                    "847553930390.dkr.ecr.ap-northeast-1.amazonaws.com/jaegertracing/jaeger-agent:1.14.0",
                    "jaegertracing/jaeger-agent:1.14.0"
                ],
                "sizeBytes": 27153511
            },
            {
                "names": [
                    "prom/node-exporter@sha256:a2f29256e53cc3e0b64d7a472512600b2e9410347d53cdc85b49f659c17e02ee",
                    "prom/node-exporter:v0.18.1"
                ],
                "sizeBytes": 22933477
            },
            {
                "names": [
                    "prom/node-exporter@sha256:b2dd31b0d23fda63588674e40fd8d05010d07c5d4ac37163fc596ba9065ce38d",
                    "prom/node-exporter:v0.18.0"
                ],
                "sizeBytes": 22889868
            }
        ],
        "nodeInfo": {
            "architecture": "amd64",
            "bootID": "44ba3810-3e58-4377-9f3f-82d12eccbbdc",
            "containerRuntimeVersion": "docker://18.6.1",
            "kernelVersion": "4.14.133-113.112.amzn2.x86_64",
            "kubeProxyVersion": "v1.14.6-eks-5047ed",
            "kubeletVersion": "v1.14.6-eks-5047ed",
            "machineID": "ec25719c8de5553a011dde23b50f91c2",
            "operatingSystem": "linux",
            "osImage": "Amazon Linux 2",
            "systemUUID": "EC25719C-8DE5-553A-011D-DE23B50F91C2"
        },
        "volumesAttached": [
            {
                "devicePath": "/dev/xvdbc",
                "name": "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0b86dba3ac16c75c6"
            },
            {
                "devicePath": "/dev/xvdbh",
                "name": "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0c92dde76500b68a8"
            },
            {
                "devicePath": "/dev/xvdbi",
                "name": "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-06635180718acbd2c"
            },
            {
                "devicePath": "/dev/xvdcb",
                "name": "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0a435110684ba04db"
            }
        ],
        "volumesInUse": [
            "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-06635180718acbd2c",
            "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0a435110684ba04db",
            "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0b86dba3ac16c75c6",
            "kubernetes.io/aws-ebs/aws://ap-northeast-1c/vol-0c92dde76500b68a8"
        ]
    }
}
`

func GetPod() *v1.Pod {
	pod := &v1.Pod{}
	if err := json.Unmarshal([]byte(testPod), pod); err != nil {
		panic(err.Error())
	}
	return pod

}

var testPod = `{
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
                        "cpu": "2",
                        "memory": "2G"
                    },
                    "requests": {
                        "cpu": "1",
                        "memory": "1G"
                    }
                }
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30
    }
}
`

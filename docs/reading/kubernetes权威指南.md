## 第一章 Kubernetes 入门
#### 概念与术语
1. **Master**
  - `kube-apiserver` http接口与集群控制的入口
  - `kube-controller-manage` 所有资源对象的自动化控制中心
  - `kube-scheduler` 负责资源(Pod)调度
2. **Node**
  除 Master 以外 k8s 集群中的其余机器
  - `kubelet` 负责 Pod 对应的容器创建、启停等任务，同时与 Master 节点密切协作，实现集群管理的基本功能
  - `kube-proxy` 实现 k8s Service 的通信与负载均衡机制的重要组件
  - `Docker Engine` Docker 引擎，负责本机的容器创建和管理工作
3. **Pod**
  每个 Pod 有个根容器 Pause 容器，代表整个容器组的状态。
4. **Label**
  可以附加到各种资源对象上,如 Node,Pod,Service,RC 等。
  Label 和 Lable Selector 构成了 k8s 系统中最核心的应用模型，使得被管理对象能够被精细地分组管理，同时实现了整个集群的高可用性。
5. **Replication Controller**
  定义了一个期望的场景，包括：
  - Pod 期待的副本数 replicas
  - 用于筛选目标 Pod 的 Label Selector
  - 当 Pod 的副本数量小于预期数量时，用于创建新 Pod 和 Pod 模板
6. **Deployment**
  为了更好地解决 Pod 的编排问题。
  内部使用了 Replica Set 来实现，相当于 RC 的升级
7. **Horizontal Pod Autoscaler**
  HPA 可以有一下两种方式作为 Pod 负载的度量指标
  - CPUUtilizationPercentage
  - 应用程序自定义的度量指标，如TPS/QPS
8. **Service**
  Service 即微服务，系统由多个提供不同业务能力而又彼此独立的微服务单元所组成，服务之间通过 TCP/IP 进行通信，从而形成了强大而灵活的弹性网格。
  服务发现： 通过 Add-On 增值包的方式引入了 DNS 系统，把服务名作为 DNS 域名
  三种IP:
  - `Node IP` Node 节点的物理网卡真实IP
  - `Pod IP` 每个 Pod 的虚拟 IP,是 Docker Engine 根据 docker0 网桥的 IP 地址进行分配的
  - `Service IP` 虚拟IP，仅仅作用于 k8s Service 这个对象，由 k8s 分配和管理，无法被 Ping，结合 Service Port 才能通信
9. **Volume**
  被定义在 Pod 上，被一个 Pod 里的多个容器挂载到具体的文件目录下
  Volume 生命周期和 Pod 相同，不随着容器终止而丢失
  Volume 类型主要有：
  - `emptyDir`
  - `hostPath` Pod 挂载在宿主机上的目录或文件
  - `gcePersistentDisk` Node 必须为 GCE 虚拟机，Google
  - `awsElasticBlockStore` 亚马逊公有云
  - `NFS`
10. **Persistent Volume**
  PV 为 k8s 集群中的某个网络存储中对应的一块存储
11. **Namespace**
  namespace 实现多租户的资源隔离

## Kubernetes 实践指南
1. Pod 的状态
  - `Pending` API Server 已经创建该 Pod，但 Pod 内还有一个或多个容器的镜像没有创建，包括正在下载镜像的过程
  - `Running` Pod 内所有容器均已创建，且至少有一个容器处于运行状态、正在启动状态或正在重启状态
  - `Succedded` Pod 内所有容器均成功执行退出，且不会再重启
  - `Failed` Pod 内所有容器均已退出，但至少有一个容器退出为失败状态
  - `Unknown` 由于某种原因无法获取该 Pod 的状态，可能由于网络通信不畅导致

## Kubernetes 源码分析

#### kube-apiserver

#### kube-controller-manager

#### kube-scheduler

#### kubelet

#### kube-proxy

#### kubectl
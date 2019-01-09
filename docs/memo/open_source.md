# 开源漫谈
*本文仅做科普性的介绍，若有部分观点值得商榷，欢迎联系[个人邮箱](jun.pan@shunwang.com)探讨* 

## 主流开源产品
- Android -> Symbian, IOS
- Linux -> Windows
- MySQL -> Oracle
- Apache
- OpenStack
- Hadoop
- Docker
- Kubernetes

## 中国企业对开源的态度
#### 以前
- 免费，质量不高，高危 （宁肯购买Oracle服务、聘请Oracle专家）
- 很少有人贡献代码，或者贡献的内容质量低

#### 现在
- 软件架构基本由开源项目搭建
- 阿里、华为、百度等公司，已是许多主流开源项目的重要贡献者
- 逐渐涌现个人和小型公司，推出高质量的开源产品

## 为什么开源是一种趋势
#### 闭源软件的局限性
**开发是一个高成本的工作**
```
A simple “Hello world” is a miracle:
- Compiler
- Linker
- VM (maybe)
- OS
```
```
A RPC “Hello world” is a miracle:
- Coordinator (zookeeper, etcd)
- RPC implementation
- Network stack
- Encoding/Decoding library
- Compiler for programming languages or [protocol buffers, avro, msgpack, capn]
```
#### 那么，为什么开源产品能成功呢
- **扩展性** 免费，使用范围广
- **快速迭代** 更容易发现问题，并迅速解决问题
- **安全性** 代码透明，安全可控，二次开发
- **发展性** 市场可引导产品的发展方向
> 成功的开源产品是市场价值的体现(淘汰得更多,例如其余容器技术)

## 开源给我们带来了什么？
#### 软件开发周期大大减少
- **高层语言和框架的出现** PHP，Golang, JS框架
- **分工更加明确** 专注于各自的领域
- **规范性要求高** 有助于实现自动化

#### 企业更专注于核心业务模块
- 目前开源产品主要集中于基础架构层
- 基础架构的稳定，保证可以集中于业务层的开发
- 也可针对开源产品的维护、二次开发，作为业务方向

## 如何参与到开源项目
####follow 项目
访问对应代码托管的地址，例如 [Linux](https://github.com/torvalds/linux) 与 [Kubernetes](https://github.com/kubernetes/kubernetes)
####阅读官方文档
了解官方给出的介绍、操作、常见问题等信息
####与他人交流
因项目而异，会采用 邮件、Blog、社区、Issue、Slack、在线会议 等不同的方式
####使用开源项目
在本地或者测试服务器上部署，体验其功能，将问题反馈到社区
####反馈社区
- 提交Issue，要求有详细的报错和复现过程
- 完善文档，开源项目的最主要交流工具是文档(甚至出现了 Technique Writer 的工作)
- 提交代码，包括 问题修复/新增特性

## 我们从开源能获得什么？
#### 简单来说
- Coding能力
- 出现问题时的定位能力
- 文字表达能力
- 获得社区的帮助(尤其是 Bug Fix)
#### 更深层次
- 了解行业领域的发展方向
- 具备一定的项目管理能力
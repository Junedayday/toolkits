# 编程入坑

## 本文目的: 帮助大家在学习编程语言前，对学习编程语言有个大致的了解

欢迎随时提问

---

## 技术路线图 (以 devops 为例)
[原文链接](https://github.com/kamranahmedse/developer-roadmap)

![Web Developer Roadmap Introduction](../img/intro.png)

![DevOps Roadmap](../img/devops.png)

---

## 首先，也是最重要的: 我为什么学习编程
- **懒** 逐渐引入自动化，减少重复工作，有更多的个人时间
> 用一个 shell 脚本解决几条多个操作，如清理日志; 更进一步，加个定时功能
- **不被忽悠** 理解开发的相关术语和思想
> 无论是 产品经理/项目经理/测试/运维/开发 等等，了解对方开发的工作量
- **提高竞争力**
  - 测试岗: 从 `功能测试` 到 `自动化测试`，再到 `开发自测`
   > **功能测试** 到处点点点，发现问题现象后找研发; **自动化测试** 根据产品或代码，写测试脚本;  **开发自测** 引入测试框架，测试代码量远大于开发代码，如测试驱动开发的模式
  - 运维岗: 从 `传统运维` 到 `devops`
  - 产品/项目经理: 从 `纯业务/客户导向` 到 `了解技术趋势`

## 如何选择适合自己的编程语言?
- **编程语言的分类** 
  1. 前端/后端
  2. 面向过程编程/面向对象编程/函数式编程 
  3. 脚本语言/编译语言
  4. 初学者更关注于难易程度，这里有个观点很重要: `编程语言是有 Level 的`

用 C++ 原生代码起一个 http 服务需要数十行，而 Go 只需要一行(Go里会牺牲很多底层的特性)。我尝试用 MySQL 不完全类比一下:
> 我们在对一个数据库的复杂应用时，需要连接到 MySQL 上，进行相关操作:
**方案一** 执行多个 SQL 语句
**方案二** 执行一个已经创建好的存储过程
两种方式都是一种编程的思想，对使用者来说，了解各自提供的接口(原生SQL语法 和 存储过程的参数意义)。孰优孰劣，无法下定义，前者底层，后者高效。

- 个人推荐
  1. **python** 易入门，适合作为工具，更专注于使用的场景，如自动化、数据分析
  2. **go** 了解计算机底层的基础上，开发中小型项目的利器
  3. **javascript** 兼容性最好的语言、全栈
  4. **java** 开发大型项目的核武器
  5. **C/C++** 底层语言(Unix/Linux 的操作内核就是C写的)
  6. 小众语言: Scala/Haskell/Rust/Lisp...
> **如何选型** 结合 学习难度/应用范围
以 TiDB 为例，内部采用的是 Go/Rust 的组合方式，测试时又引入了大量C++/Python 等等

## 那么，如何开始学习一门编程语言
  - **确定目标** 我要学到什么程度
    1. 读完《七天精通XXX》
    2. 了解主要的语法，会写 Hello World！
    3. 能独立写一款小工具
    > 例如检查某个库表的字段是否合规
    4. 能合作参与项目，做出一定成果并分享
    5. ...
  - **环境准备** 工欲善其事，必先利其器
    1. 代码环境
    2. 学习氛围
    3. 分享交流
  - **边学边用** "功利性"地学语言
    1. 实现一个功能，写一点代码，掌握一些基础
    > 不建议采用学校里的方式，先阅读大量基础教材
    2. 解决问题思路: 搜索引擎>基础教程>寻求他人帮助
  - **实践应用** 实践是最好的老师
    1. 提升个人效率
    2. 应用在部门的自动化平台上
    3. 了解数据库技术的部分底层实现

## 如何让自己持续地获得动力
  - 能实际提升 个人/部门 的效率，更多自由时间
  - 多分享，多交流，获得集体的认同感
  - 在 部门/公司 层展示成果
  - 在自己的成果上打上个人标签/彩蛋(非商用)
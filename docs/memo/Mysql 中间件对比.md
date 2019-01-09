## MySQL 中间件技术选型初探
1. 商业应用 
  - 商用可以在一定程度上证明该技术的成熟度
  - 尤其是在大公司的使用，可以保证该工具适用于高并发、高数据量的场景

2. 开发语言
  - 由于 Go 语言的特性,在应用层开发效率高,大量应用于各类中间件的开发 
  - JAVA 语言对底层、应用层都支持，开发效率相对较低，但生态圈强大,有不少现成的整体框架 

3. 扩展性
  - 中间件的发展前景是否满足我们的应用场景
  - 与其余工具的对接开发

4. 社区活跃度
  - 影响到 bug 的 fix 效率
  - 对新特性的开发速度
  - 和商业应用互相影响

## 常见 MySQL 中间件对比

|中间件名称|商业应用|开发效率|扩展性|社区活跃度|
|:---:|:---:|:---:|:---:|:---:|
|[Cobar](https://github.com/alibaba/cobar)|不维护||||
|DRDS|阿里商用||||
|[DTLE](https://github.com/actiontech/dtle)|上海爱可生|Golang|目前支持的是Mysql的增删改的数据解析|开源周期短，[活跃度](https://github.com/actiontech/dtle/graphs/contributors)暂时不高|
|[Vitess](https://github.com/vitessio/vitess)|京东,Youtube|Golang|偏向于支持 MySQL 的原生版本|[活跃度](https://github.com/vitessio/vitess/graphs/contributors)较高，社区交流偏英文|
|[Mycat](https://github.com/MyCATApache/Mycat-Server)|中国电信，中国联通|JAVA|由于其架构特点,具有高扩展性,详情可见[规划](http://www.mycat.io/)|[活跃度](https://github.com/MyCATApache/Mycat-Server/graphs/contributors)较高，社区交流偏中文|

> DTLE 的局限点：
> - 依赖于很多其余的开源库，本身是一个比较轻量级的工具
> - 对 MySQL 解析没有到源文件阶段 (例如 Vitess 会对 sql 语法的源文件 `sql_yacc.yy` 进行解析)
> - 对表的索引、主键等等，依赖于用 sql 查询 (例如 Vitess 可以直接从磁盘文件里，读出表的字段结构、索引等等)
> - 应用场景只限于 爱可生公司，应用起来风险未知

## 详细对比 Mycat 与 Vitess 各自优势
### Mycat 优势
- 架构支持大规模的 MySQL 集群应用场景
- 官方建议应用于 1000w+ 条数据的单表，读写分离，分库分表等功能都很完善
- 开源社区是中文，方便阅读
- 有大量的文档书籍，门槛较低
- 项目的扩展性好，后续功能可以持续引入，如自动索引，对接 HDFS 等等
- 可以引入不少现成的运维体系

### Vitess 优势
- 内部组件 vttable 是一对一映射到一个 MySQL 实例，做二次开发效率高
- 系统架构相对来说简单，维护代码的成本会比较低
- 适合开发轻量级 MySQL 插件
- 重点关注 MySQL 的扩展性，是 CNCF 中的一员(即后续在云服务、容器这块支持性会很好)

**总体来说，为保证稳定性，选择 Mycat 或者 Vitess 这两个开源技术，都必须有团队了解其开源代码(除非找第三方做商业的技术支持)，Mycat 在这方面需要更多的开发和运维成本。在此基础上，从长远的发展角度来看，Mycat 的扩展性和发展前景，更为适合；而 Vitess 适合敏捷开发、二次开发，对现有数据库的架构改动比较小。**
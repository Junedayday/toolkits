# Differences
- mysql存放在磁盘中，检索涉及到IO
- redis比memcache多了几种数据结构封装，适用业务
- web应用先访问redis，再访问mysql

# Usage
- redis实现数据； 读写，用队列定时将数据写入mysql
- 方案1：读: 读redis->没有，读mysql->把mysql数据写回redis； 写: 写mysql->成功，写redis
- 方案2：redis读取mysql的binlog(对mysql理解不够)

# Why use redis
- 确定是否是数据库瓶颈
- 缓存量大但又不经常变的东西
- 考虑是否需要主从，读写分离，分布式，水平伸缩

# redis数据结构
  - string,hash,list,set,sorted sets，stream
  - 使用持久化时，需要考虑持久化和写性能的配比，也要考虑redis使用的内存大小和硬盘写速率的比例
  - redis是单进程，一个实例用到一个CPU
  - master->master存储到slave的rdb，slave加载到内存
  - 数据一致性：长期运行后，会做周期性地检查全量数据，实时检查增量数据。对于主库未及时同步导致的不一致，称之为延时问题。

# redis缺点
  - 缓存和数据库双写一致性问题
    最终一致性和强一致性，强一致性不能放缓存（采取正确更新策略，先更新数据库，再删缓存。其次，因为可能存在删除缓存失败的问题，提供一个补偿措施即可，例如利用消息队列）
  - 缓存雪崩问题
  - 缓存击穿问题
  - 缓存的并发竞争问题
    准备一个分布式锁，大家去抢锁，抢到锁就做set操作（如果多操作，可能加入时间戳）

# redis为什么这么快
  - 纯内存操作
  - 单线程操作，避免频繁的上下文切换
  - 采用了非阻塞的I/O多路复用机制

# redis数据结构使用
  - string 复杂的计数功能的缓存,SET/INCR
  - hash 单点登录，cookieID作为key
  - list做简单的消息队列的功能，做基于redis的分页功能
  - set做全局去重；利用交集，并集，差级，计算共同还好，全部喜好，独有喜好
  - sorted set 权重参数，做排行榜应用，top N，做延时任务，做范围查找

# redis删除机制：定期删除+惰性删除
  - 定期删除是100ms随机检查，惰性删除是调用key时进行检查
  - 内存淘汰机制：推荐allkeys-lru(移除最近最少使用的key)

# redis怎么实现bgsave
    fork一个子进程，copy on write

# redis实现分布式锁
  SETEX lock.foo timestamp
  为了防止crash，增加个expire

# Source Code
http://bbs.redis.cn/forum.php?mod=viewthread&tid=544 
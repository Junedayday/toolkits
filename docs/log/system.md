# Linux System

- [2019/01/23 Finish Reading Linux-vbird Book](#finish-reading-vbird-linux)
- [2019/01/14 Start reading Linux-vbird Book](#start-reading-vbird-linux)
- [知识点](#key-points)

## Finish Reading Vbird Linux

完成了本书的初次阅读，总体来说难度低，对 Linux 的体系有了更系统的了解

## Start Reading Vbird Linux

尝试了几本 unix 基础的专业书后，由于很难深入，因此先转向《鸟哥的Linux私房菜》，从入门开始
共计 Chap 26

## Key Points

- 鸟哥认为学习的动力在于两点: `兴趣` 和 `成就感`。个人认为 `成就感` 更为重要，一定的成就，容易激发兴趣
- 主分区+扩展分区 最多 4 个，扩展分区最多 1 个，逻辑分区是由扩展分区切割出来的分区
- Linux 下的各文件夹存放都有具体意义，详情可参考 [FHS](http://www.pathname.com/fhs)
- 文件查询 which/whereis/locate/findfind
- Linux FS Ext2 主要包括 superblock/inode/block 三块
- 管道常用命令: `less` `grep` `tee` `tr` `xargs` `awk`
- 正则表达式需要保证语系 `LANG = C`，常用字符 `[]` `^`(2种含义) `$` `\{\}` `.`, 扩展 `+` `?` `|` `()`
- 可用 `crontab` 配置例行性工作调度，`at` 配置单次任务
- 系统资源查看 `free` `uname` `uptime` `netstat` `dmesg` `vmstat`
- stand alone daemon 一般放在 `/etc/init.d/*` 目录下，super daemon 放在 `/etc/xinetd.conf` 和 `/etc/xinetd.d/*`，开机启动可用 `chkconfig`
- Linux 开启 `syslog` 后，常见的日志文件名 `var/log目录`：`cron` `dmesg` `maillog,mail/` `lastlog` `messages` `secure` `wtmp,faillog` `httpd/,news/,samba/` ... CentOS 官方提供了 `logwatch` 功能
> 部分暂不深入，在实际使用中再研究：shell script，账号管理，磁盘配额
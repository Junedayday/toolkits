# Git

---
Git is a free and open source distributed version control system designed to handle everything from small to very large projects with speed and efficiency.

- 分布式
- 分支管理

## Compare to SVN

- Client-Server 模型 checkout-commit
- if server failed？

## Issues for Git

- 学习成本相对变高
- 权限控制方式转变
- 代码泄漏
- 二进制文件支持性差

## Install Git Client

[官方链接](https://git-scm.com/downloads)
> 尤其推荐在 windows 上进行安装
> 另外也可安装 TortoiseGit

## Before Git

- **English**
- Gitlab/Github and Account
- Markdown:
  - IDE: 推荐 vscode,emacs,有道云笔记,简书,[stackedit](https://stackedit.io/app#) 等等
  - 教程: [中文链接](https://www.jianshu.com/p/191d1e21f7ed)
- linux basic operation

## One Basic Conception：Branch

- 分支，对应的主干是 master
- git 的 branch 优势是切换快/合并快
- 不要有太多的 branch

## 5 commands to get start in local

#### 流程：下载代码-选择分支-查看改动-添加改动-保存改动

1. git clone `url`
2. git checkout `branch_name`
  - 如果是新建一个分支，那就在名字前加上 `-b`
  - 尽量把branch_name带有标示性，例如 dev/pj_dev
3. git status
4. git add `filename1 filename2 ...` 或者全部文件 `-A .`
  > 个人建议这个工作放在 IDE 上进行,尤其是文件多的时候
5. git commit -m "message"
  - 最关键的一步，提交的信息写得详细，对他人的帮助很大

## 3 commands for interaction

#### 获取更新-(合并更新)-推送更新

1. git pull
> 所有远端的分支更新，都会被拉下来，自动合并到本地的分支中
2. git merge `target_branch_name`
> 注意先查看当前所在的分支，避免合并错误
3. git push origin `branch_name`
> 注意先查看当前所在的分支，避免推送错误
> 一般操作远程的分支，都会带 origin 这个关键字
---

*合并时，会小概率出现 conflict, 尽量在 IDE 解决后再 commit*

## step further on（with a target）

- 我的库弄乱了，我想回到以前 commit 过的一个版本（阶段性 commit 的重要性）
  - git log 查看日志，每一个 commit 都有一个 commit id/hashcode ，是全局唯一
  > 如果嫌每个 log 显示太多行，可加 `--oneline` 后缀
  - git reset --hard `hashcode` 强制恢复
  > 如果想直接回到最近一次的 commit, 输入 git reset --hard head
- 我的 commit 的记录太多了，想减少点
  - git commit --amend 直接把本次提交附加到上次的提交后，使用同一个 hashcode，也可以加上 message，会自动替换
  - git rebase -i HEAD~X 整理最近X条 commit 记录，在每个 commit 前可以选择 pick/reword/squash/fixup/edit/exec/drop
  > 一般前四个够用:

  rebase 中指令|作用
  ----|----
  pick | 继续使用
  reword | 重新编辑message
  squash | 合并到上条并合并message
  fixup | 合并到上条并抛弃本次的message
- 我在某个 branch 上有个 commit，想移植到另一个 branch 中，但是又不想移植其余的 commit
  - git cherry pick `hashcode`
  > 没错，操作就是这么简单，问题是 branch 之间可能存在差异，可能需要人工改动
- 有个频繁改动的文件，但对于我来说没什么用，每次提交时总要自己去忽略它，很麻烦
  
  - git 项目的根目录上，新建一个 `.gitignore` 文件，可指定提交时忽略

## Git on web page

-  merge request: 最核心的权限控制，protect branch
  - 发起请求
  - 选择 src_branch 和 target_branch, 对比差异
  - 加上合并理由、说明，提交请求
  - 负责人审核代码，有问题的话加上备注，让发起者重新修改
  - 确认无误的话，合并成功 
- diff:对比不同 branch/commit_id 之间的文件差异

## Workflow in github

以 `kubernetes` 官方的github为例. 
[链接](https://github.com/kubernetes/community/blob/master/contributors/guide/github-workflow.md)

## Best Practice?

- [开发案例](https://nvie.com/posts/a-successful-git-branching-model/)

- [version control](https://stackoverflow.com/questions/47883823/version-controlling-with-mysql-databases)

## Why Git?

- 开源的趋势
- 减少对多库的维护
- 与第三方工具集成(文档,Task,Bug,CI,CD...)
...
总而言之，为了提升开发效率，让开发人员有更多的自由时间。x
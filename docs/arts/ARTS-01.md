## Algorithms

[2. 两数相加](https://leetcode-cn.com/problems/add-two-numbers/)

```java
class Solution {
    public ListNode addTwoNumbers(ListNode l1, ListNode l2) {
        ListNode returnList = new ListNode(0);
        ListNode curList = returnList;
        boolean carryOver = false;
        while (l1 != null || l2 != null || carryOver) {
            int re = 0;
            if (l1 != null) {
                re += l1.val;
                l1 = l1.next;
            }
            if (l2 != null) {
                re += l2.val;
                l2 = l2.next;
            }
            if (carryOver) {
                carryOver = false;
                re++;
            }
            if (re >= 10) {
                re -= 10;
                carryOver = true;
            }
            curList.next = new ListNode(re);
            curList = curList.next;
        }
        return returnList.next;
    }
}
```

> 解题时有个初始化问题 `l1=null,l2=null` 时很难处理。
>
> 这里有个思路转变 `不一定直接定义返回结果，可以用其指针next来表示，解决了初始化的难题`

## Review

Unix Programming: Process Relationships

- login 采用了 fork 子进程的技术（本地和网络登录有少量区别）

- process group 组合了多个 process， session又组合了多个 process group

  

## Tip

说明简单的安装过程时，脚本+文档的形式更高效。

有时，我们安装某个工具或者软件，会写一个冗长的文档，一旦用户不仔细阅读，就会出现各种奇怪的问题。这时，采用脚本+文档的组合，会有更高的效率。

- `脚本` 覆盖大部分的情况，可供不求深入了解的人，进行一键安装
- `文档` 说明关键性的内容，如原理、特殊情况等，不宜过长



## Share

最近在恶补计算机基础知识，在这个过程中，发现了很多和当前主流技术的共通之处。窃以为，如果想要在程序员这个行业里深耕，学习基础是一个必要的学习内容，越早投入、收益越大。

高层次的编程技术，很容易因主流技术的变迁而过时，而对新技术的学习周期，很大程度决定于基础的扎实程度。
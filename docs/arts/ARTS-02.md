## Algorithms

[3. 无重复字符的最长子串](https://leetcode-cn.com/problems/longest-substring-without-repeating-characters/)

```java
class Solution {
    public String longestPalindrome(String s) {
        if (s == null || s.length() < 1) return "";
        int start = 0, end = 0;
        for (int i = 0; i < s.length(); i++) {
            int len1 = expandAroundCenter(s, i, i);
            int len2 = expandAroundCenter(s, i, i + 1);
            int len = Math.max(len1, len2);
            if (len > end - start) {
                start = i - (len - 1) / 2;
                end = i + len / 2;
            }
        }
        return s.substring(start, end + 1);
    }

    private int expandAroundCenter(String s, int left, int right) {
        int L = left, R = right;
        while (L >= 0 && R < s.length() && s.charAt(L) == s.charAt(R)) {
            L--;
            R++;
        }
        return R - L - 1;
    }
}
```

> 刚看到回文相关问题，就觉得挺麻烦的，最优解也需要不小的工作量
>
> 这个解法的代码量不大，关键在于 `回文的最大特点是存在一个中心点（区分长度的奇偶性），从这个中心店进行遍历`

## Review

No Silver Bullet   -- Essence and Accident in Software Engineering




## Tip

`Golang` 的模块管理 `go mod` 并不像其余编程语言的包管理软件那么完善，存在不少弊端。

与 go 1.11 版本之前的 vendor 方案进行兼容也需要一定的考量，例如指定 `-mod=vendor` 等，但这不是最佳方案。

期待 go 2.0 能带来新的体验，同时自己有时间，也可投入精力到原先的实现细节中。



## Share

十一刚从澳大利亚旅游一圈回来，由于行程紧张，打乱了原来的学习计划。

澳大利亚之行带给我的最大体会，就是要 `更精致地工作与生活`。简单地概括下：

1. 拥有东西不在于多，而在于精；
2. 事情不在于快速完成，而是享受其中的过程，并有所得；
3. 锻炼身体，修养心态，有条理地进行工作和生活；

接下来一段时间，我会尝试着提高自己的工作和生活的水平（大概率包括消费水平），让自己的职业生涯更长久。

高层次的编程技术，很容易因主流技术的变迁而过时，而对新技术的学习周期，很大程度决定于基础的扎实程度。
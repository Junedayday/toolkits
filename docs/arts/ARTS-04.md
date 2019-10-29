## Algorithms

[9. 回文数](https://leetcode-cn.com/problems/palindrome-number/)

```java
class Solution {
    public boolean isPalindrome(int x) {
        if(x < 0 || (x % 10 == 0 && x != 0)) {
            return false;
        }

        int revertedNumber = 0;
        while(x > revertedNumber) {
            revertedNumber = revertedNumber * 10 + x % 10;
            x /= 10;
        }

        return x == revertedNumber || x == revertedNumber/10;
    }
}
```

> 回文数的主要思想容易想到，但有2个细节容易发生代码冗余（为了方便区分，取名为`被截数`与`截取数`）：
>
> 1. 如何确定`截取数`达到了一半的位数？  对比 `被截数`与`截取数`
> 2. 如何区分奇位数？直接将 被截数/10 

## Review

[How to Read the Right Way: A Complete Guide](https://medium.com/the-mission/how-to-read-the-right-way-a-complete-guide-82042876be2c)

## Tip

用vim编辑时，善用可视模式，能快速地对大量文本进行复制、剪切、黏贴等操作。

## Share

新款的 Macbook 带了 Touch Bar，这对重度 vim 用户的我来说，无疑是一个很糟糕的体验。

最近几年，如果 Macbook Pro 再无明显提升，我个人将大概率抛弃 macbook 系列，转用 thinkpad + ubuntu 的工作环境。


## Algorithms

[6. Z 字形变换](https://leetcode-cn.com/problems/zigzag-conversion/)

```java
class Solution {
    public String convert(String s, int numRows) {
        if (numRows == 1) return s;

        List<StringBuilder> rows = new ArrayList<>();
        for (int i = 0; i < numRows; i++)
            rows.add(new StringBuilder());

        int curRow = 0;
        boolean goingDown = false;

        for (char c : s.toCharArray()) {
            rows.get(curRow).append(c);
            if (curRow == 0 || curRow == numRows - 1) goingDown = !goingDown;
            curRow += goingDown ? 1 : -1;
        }

        StringBuilder ret = new StringBuilder();
        for (StringBuilder row : rows) ret.append(row);
        return ret.toString();
    }
}
```

> Z 字形变换可以直接通过计算得出，但在代码侧不够直观。直接通过N行数据及上下移动的方式，更方便理解。
>
> 

## Review

Kafka: a Distributed Messaging System for Log Processing




## Tip

`Golang` 的标准库很多，许多人都忽略了标准库提供的一些功能，例如 `查找切片中的元素`，`sort` 包中有对应实现(且不论效率)。

加强学习基础库，是提高效率、把代码写得更优美的一个很便捷的实现，尽管最近一直在补基础，这块也得记住。



## Share

之前个人过多地关注于 `读书`，而忽略了 `选择书` 这一过程。一本真正的好书，作者不会仅仅会花时间于内容，而在 标题、目录、序言、索引 等细节上，也会细啄。

之后，在选择书之前，更应关注这些内容，让自己能迅速地分辨出书的优劣。
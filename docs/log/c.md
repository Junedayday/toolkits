# C Learning Log

- [2019/01/16 C Programming Language e-Book finished](#book-finished)
- [2019/01/15 C Programming Language e-Book started](#book-started)
- [知识点](#key-points)

## Book Finished

- Chap 8 均完成
- 对 标准输入输出函数 和 Unix系统接口 这两块未细看，后续结合操作系统研究

## Book Started

- Start to Read 《The C Programming Language》
- 由于有一定基础，总体难度偏低，有部分知识点值得记忆
- 总计 Chap 8，今日阅读至 Chap 4，其中最有价值的是 Chap 4

## Key Points

1. 数字和运算符可用 `逆波兰表示法` 压入栈中，如 `(1 - 2) * (4 + 5)` => `12-45+*`
2. static 静态变量作用域仅为 `当前文件`
3. 宏处理注意格式，如 `#define square( x ) x * x` 作用于 `square(y + 1)`;并且会使用 `if  !defined()` 格式控制重复包含
4. 在数组中, `pa[i]` 与 `*(pa+i)` 等价， 其中 `char*` 也可认为是数组
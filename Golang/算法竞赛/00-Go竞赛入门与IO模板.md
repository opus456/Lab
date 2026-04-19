# Go 算法竞赛入门与 IO 模板

> **目标**：从零搭建 Go 刷题环境，掌握竞赛中最高效的输入输出方式，熟悉 Go 语言语法中与竞赛最相关的特性。

---

## 一、为什么用 Go 打算法竞赛？

| 优势 | 说明 |
|------|------|
| 编译型语言 | 速度接近 C++，远快于 Python |
| 自带大数 | `math/big` 直接处理大整数 |
| 内置排序 | `sort.Slice` 一行搞定自定义排序 |
| GC 管理内存 | 不需要手动 `new/delete` |
| 编译速度快 | 秒级编译，调试体验好 |

| 劣势 | 说明 |
|------|------|
| 无 STL | 没有 C++ 的 set/map/priority_queue 那么方便，需要手写 |
| 泛型较新 | Go 1.18+ 才支持泛型 |
| 部分 OJ 不支持 | 但 Codeforces、LeetCode、AtCoder 均支持 |

---

## 二、环境搭建

### 2.1 安装 Go

```bash
# Windows: 下载 https://go.dev/dl/ 安装包
# 验证安装
go version
# 输出类似: go version go1.22.0 windows/amd64
```

### 2.2 竞赛文件模板

创建一个 `main.go`：

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

var reader *bufio.Reader
var writer *bufio.Writer

func main() {
    reader = bufio.NewReader(os.Stdin)
    writer = bufio.NewWriter(os.Stdout)
    defer writer.Flush() // 程序结束前把缓冲区内容全部输出

    // 你的代码写在这里
}
```

### 2.3 为什么需要 bufio？

Go 的 `fmt.Scan` 和 `fmt.Print` 直接操作标准输入输出，**每次调用产生一次系统调用**。当输入数据量大（如 10^6 个数）时，逐个 `Scan` 会超时。

`bufio.Reader` 和 `bufio.Writer` 在用户空间维护缓冲区，批量读写，性能提升 **10-50 倍**。

```
无缓冲: Scan → 内核 → 返回 → Scan → 内核 → 返回 → ...  (每次都进内核)
有缓冲: Reader.Read → 一次读入大块数据到缓冲区 → 后续从缓冲区取 (极少进内核)
```

---

## 三、输入输出详解

### 3.1 基本输入

```go
// 方式1: fmt.Fscan (推荐竞赛用法)
var n, m int
fmt.Fscan(reader, &n, &m)

// 方式2: 读一整行
var line string
line, _ = reader.ReadString('\n')

// 方式3: 读浮点数
var x float64
fmt.Fscan(reader, &x)

// 方式4: 读字符串
var s string
fmt.Fscan(reader, &s) // 读到空格或换行为止
```

### 3.2 基本输出

```go
// 方式1: fmt.Fprintln (自动换行)
fmt.Fprintln(writer, ans)

// 方式2: fmt.Fprintf (格式化)
fmt.Fprintf(writer, "%d %d\n", x, y)

// 方式3: 输出数组，空格分隔
for i, v := range arr {
    if i > 0 {
        fmt.Fprint(writer, " ")
    }
    fmt.Fprint(writer, v)
}
fmt.Fprintln(writer)
```

### 3.3 完整输入输出模板

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

var reader *bufio.Reader
var writer *bufio.Writer

func main() {
    reader = bufio.NewReader(os.Stdin)
    writer = bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    fmt.Fscan(reader, &t) // 读测试用例数量

    for ; t > 0; t-- {
        solve()
    }
}

func solve() {
    var n int
    fmt.Fscan(reader, &n)

    a := make([]int, n)
    for i := range a {
        fmt.Fscan(reader, &a[i])
    }

    // 处理逻辑...
    ans := 0
    fmt.Fprintln(writer, ans)
}
```

### 3.4 超快速输入（手写 readInt）

当数据量极大（>10^6）时，`fmt.Fscan` 也可能不够快。手写字节级读取：

```go
func readInt() int {
    n := 0
    c, _ := reader.ReadByte()
    neg := false
    // 跳过空白
    for c == ' ' || c == '\n' || c == '\r' {
        c, _ = reader.ReadByte()
    }
    if c == '-' {
        neg = true
        c, _ = reader.ReadByte()
    }
    for c >= '0' && c <= '9' {
        n = n*10 + int(c-'0')
        c, _ = reader.ReadByte()
    }
    if neg {
        return -n
    }
    return n
}
```

**速度对比**（读入 10^6 个整数）：

| 方法 | 耗时 |
|------|------|
| `fmt.Scan` (无缓冲) | ~3000ms |
| `fmt.Fscan` + bufio | ~300ms |
| 手写 `readInt` | ~80ms |

---

## 四、Go 语法速查（竞赛常用）

### 4.1 变量与类型

```go
// 基本类型
var n int        // 0
var f float64    // 0.0
var s string     // ""
var b bool       // false

// 短声明（函数内）
n := 42
s := "hello"

// 类型转换（Go 不自动转换，必须显式）
x := 3.14
y := int(x)     // 3
z := float64(y) // 3.0

// 常用类型大小
// int:    64位平台上为64位
// int32:  -2^31 ~ 2^31-1
// int64:  -2^63 ~ 2^63-1  (竞赛中常用，对应C++的long long)
// 注意：两个 int 相乘可能溢出，必要时用 int64
```

### 4.2 数组与切片

```go
// 固定大小数组
var arr [100]int

// 切片（动态数组，竞赛最常用）
a := make([]int, n)       // 长度为n，全0
b := make([]int, 0, n)    // 长度为0，预分配容量n
c := []int{1, 2, 3}       // 字面量初始化

// 追加
b = append(b, 42)

// 切片操作
sub := a[1:4]  // a[1], a[2], a[3]（不含a[4]）

// 二维切片
grid := make([][]int, n)
for i := range grid {
    grid[i] = make([]int, m)
}
```

### 4.3 Map（哈希表）

```go
// 创建
mp := make(map[string]int)

// 插入/更新
mp["hello"] = 1

// 查找
val, exists := mp["hello"]
if exists {
    // 存在
}

// 删除
delete(mp, "hello")

// 遍历（顺序不确定！）
for key, val := range mp {
    fmt.Println(key, val)
}

// 当计数器用
count := make(map[int]int)
for _, v := range arr {
    count[v]++
}
```

### 4.4 字符串

```go
s := "hello"
// 长度
len(s) // 5 (字节数，非字符数!)

// 遍历字节
for i := 0; i < len(s); i++ {
    fmt.Println(s[i]) // byte 类型
}

// 遍历字符(rune)
for _, ch := range s {
    fmt.Printf("%c", ch) // rune 类型 (int32)
}

// 字符串不可变！要修改需转为 []byte
bs := []byte(s)
bs[0] = 'H'
s = string(bs)

// 拼接（竞赛中大量拼接用 strings.Builder）
var sb strings.Builder
sb.WriteByte('a')
sb.WriteString("bc")
result := sb.String() // "abc"
```

### 4.5 排序

```go
import "sort"

// 对 []int 排序
a := []int{3, 1, 4, 1, 5}
sort.Ints(a)  // [1, 1, 3, 4, 5]

// 自定义排序
sort.Slice(a, func(i, j int) bool {
    return a[i] > a[j]  // 降序
})

// 对结构体切片排序
type Pair struct{ x, y int }
pairs := []Pair{{1, 3}, {2, 1}, {1, 2}}
sort.Slice(pairs, func(i, j int) bool {
    if pairs[i].x != pairs[j].x {
        return pairs[i].x < pairs[j].x
    }
    return pairs[i].y < pairs[j].y
})

// 二分搜索
idx := sort.SearchInts(a, 3) // 第一个 >= 3 的位置
```

### 4.6 数学相关

```go
import "math"

math.Abs(x)        // 绝对值 (float64)
math.Max(a, b)     // 最大值 (float64)
math.Min(a, b)     // 最小值 (float64)
math.Sqrt(x)       // 平方根
math.Pow(a, b)     // a^b (float64)
math.Inf(1)        // +∞
math.MaxInt        // int 最大值 (Go 1.17+)

// 竞赛中整数取 max/min（Go 没有整数版的 math.Max！）
func max(a, b int) int {
    if a > b { return a }
    return b
}
func min(a, b int) int {
    if a < b { return a }
    return b
}
// Go 1.21+ 内置了 min() 和 max()
```

### 4.7 控制流

```go
// if-else
if x > 0 {
    // ...
} else if x == 0 {
    // ...
} else {
    // ...
}

// if 带初始化语句
if val, ok := mp[key]; ok {
    // val 只在这个 if 块内有效
}

// for（Go 只有 for，没有 while）
for i := 0; i < n; i++ { }        // 经典 for
for i, v := range arr { }          // 遍历切片
for key, val := range mp { }       // 遍历 map
for condition { }                   // 相当于 while
for { }                             // 无限循环

// switch
switch x {
case 1:
    // 不需要 break，Go 默认不穿透
case 2, 3:
    // 多值匹配
default:
    // ...
}
```

---

## 五、竞赛常用代码片段

### 5.1 GCD 和 LCM

```go
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func lcm(a, b int) int {
    return a / gcd(a, b) * b // 先除后乘防溢出
}
```

### 5.2 快速幂

```go
func powmod(base, exp, mod int) int {
    result := 1
    base %= mod
    for exp > 0 {
        if exp%2 == 1 {
            result = result * base % mod
        }
        exp /= 2
        base = base * base % mod
    }
    return result
}
```

### 5.3 方向数组（BFS/DFS 常用）

```go
// 四方向
dx := [4]int{-1, 1, 0, 0}
dy := [4]int{0, 0, -1, 1}

// 八方向
dx8 := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
dy8 := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}
```

### 5.4 快速读入整个数组

```go
func readArr(n int) []int {
    a := make([]int, n)
    for i := range a {
        fmt.Fscan(reader, &a[i])
    }
    return a
}
```

### 5.5 取模运算辅助

```go
const MOD = 1_000_000_007

func add(a, b int) int { return (a + b) % MOD }
func mul(a, b int) int { return a % MOD * (b % MOD) % MOD }
func sub(a, b int) int { return ((a-b)%MOD + MOD) % MOD }
```

---

## 六、完整例题演示

### 题目：给定 N 个数，求所有数对 (i,j) 中 i<j 且 a[i]+a[j] == target 的对数。

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

var reader *bufio.Reader
var writer *bufio.Writer

func main() {
    reader = bufio.NewReader(os.Stdin)
    writer = bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, target int
    fmt.Fscan(reader, &n, &target)

    a := make([]int, n)
    for i := range a {
        fmt.Fscan(reader, &a[i])
    }

    // 用 map 统计
    count := make(map[int]int)
    ans := 0
    for _, v := range a {
        // 之前出现过的 target-v 的个数就是能和 v 配对的数量
        ans += count[target-v]
        count[v]++
    }

    fmt.Fprintln(writer, ans)
}
```

**思路解析**：
1. 遍历数组，对每个数 v，查找之前有多少个数等于 `target - v`
2. 用 map 记录每个值出现的次数
3. 时间复杂度 O(N)，空间复杂度 O(N)

---

## 七、刷题平台推荐

| 平台 | 特点 | 适合阶段 |
|------|------|---------|
| [力扣 LeetCode](https://leetcode.cn) | 题目分类清晰，面试向 | 入门 → 中级 |
| [Codeforces](https://codeforces.com) | 比赛多，难度梯度好 | 中级 → 高级 |
| [AtCoder](https://atcoder.jp) | 题目质量高，难度合理 | 中级 → 高级 |
| [洛谷](https://www.luogu.com.cn) | 中文，题目丰富 | 入门 → 高级 |

---

## 速查表

| 操作 | 代码 |
|------|------|
| 快速读 int | `fmt.Fscan(reader, &n)` |
| 快速输出 | `fmt.Fprintln(writer, ans)` |
| 排序 | `sort.Ints(a)` 或 `sort.Slice(...)` |
| 二分 | `sort.SearchInts(a, x)` |
| GCD | `for b!=0 { a,b = b,a%b }` |
| 快速幂 | 见 5.2 节 |
| 取模加 | `(a + b) % MOD` |
| 最值 | `max(a,b)` (Go 1.21+) |

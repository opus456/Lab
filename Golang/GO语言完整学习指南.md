# Go语言完整学习指南

## 目录
1. [Go语言简介](#go语言简介)
2. [基础概念](#基础概念)
3. [变量声明](#变量声明)
4. [常量](#常量)
5. [函数](#函数)
6. [包与初始化](#包与初始化)
7. [指针](#指针)
8. [延迟执行](#延迟执行defer)
9. [数组](#数组)
10. [切片](#切片)
11. [映射](#映射map)
12. [结构体与方法](#结构体与方法)
13. [面向对象编程](#面向对象编程)

---

## Go语言简介

Go（Golang）是由Google开发的一种静态强类型编译型编程语言。

### 特点
- **简洁高效**：语法简洁，编译速度快
- **并发支持**：原生支持并发编程
- **垃圾回收**：自动内存管理
- **快速编译**：编译为机器码，运行速度快
- **跨平台**：支持多种操作系统

### 基本项目结构
```
project/
├── go.mod          # 模块定义文件
├── main.go         # 主程序入口
└── package_name/   # 自定义包
    └── *.go        # 包内的Go源文件
```

---

## 基础概念

### 包（Package）

Go的源代码都必须属于某个包。包是Go语言的命名空间机制。

```go
package main  // 声明这是主包

import "fmt"  // 导入fmt包
```

### Main函数

每个可执行程序都必须有一个main函数作为程序入口。

```go
func main() {
    fmt.Println("Hello, Go!")
}
```

### 导入

- **单个导入**：`import "fmt"`
- **多个导入**：
```go
import (
    "fmt"
    "math"
)
```

---

## 变量声明

### 1. 使用 `var` 声明（初始值为零值）

```go
var a int
fmt.Println(a)  // 输出：0

var c string
fmt.Println(c)  // 输出：空字符串
```

**Go语言的零值规则**：
- 整数类型：0
- 浮点数类型：0.0
- 布尔类型：false
- 字符串类型：空字符串 ""

### 2. 声明时赋值

```go
var b int = 100
fmt.Println(b)  // 输出：100
```

### 3. 类型推导

```go
var c = "hello"
fmt.Printf("type of c is %T\n", c)  // 输出：string
```

### 4. 短声明操作符 `:=`（推荐方式）

```go
d := 100
e := ":="
fmt.Println(d, e)  // 输出：100 :=
```

**注意**：
- `:=` 只能在函数内部使用
- 声明全局变量必须使用 `var`

### 5. 批量声明

```go
var (
    name string = "Tom"
    age  int    = 25
    flag bool   = true
)
```

### 6. 多个变量同时赋值

```go
var a, b = 99, "I wanna study Golang"
fmt.Println(a, b)
```

---

## 常量

常量是在编译时确定的值，运行时不能改变。

### 基本常量声明

```go
const len int = 10
fmt.Println(len)  // 输出：10
```

### 常量的特点

- 必须在编译时赋值
- 不能使用 `:=` 声明
- 可以进行编译时的计算
- 不能像变量那样被修改

### 批量声明常量

```go
const (
    RED    = 0
    GREEN  = 1
    BLUE   = 2
)
```

---

## 函数

### 函数定义

```go
func 函数名(参数 类型) 返回类型 {
    // 函数体
}
```

### 1. 基本函数 - 单个返回值

```go
func foo1(a string, b int) int {
    fmt.Println("a =", a)
    fmt.Println("b =", b)
    
    c := 100
    return c
}
```

使用：
```go
r1 := foo1("hello", 99)
fmt.Println("function return value is", r1)
```

### 2. 多个返回值 - 匿名返回

```go
func foo2(a int, b int) (int, int) {
    return a, b
}

// 使用
aa, bb := foo2(10, 20)
fmt.Println(aa, bb)  // 输出：10 20
```

### 3. 多个返回值 - 命名返回

```go
func foo3(a int, b int) (r1 int, r2 int) {
    // 返回值被命名，相当于声明了局部变量
    r1 = a
    r2 = b
    return  // 直接return，自动返回r1和r2
}

// 使用
aaa, bbb := foo3(10, 20)
fmt.Println(aaa, bbb)
```

### 4. 可变参数函数

```go
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

result := sum(1, 2, 3, 4, 5)
```

### 5. 函数作为参数

```go
func apply(f func(int) int, x int) int {
    return f(x)
}

double := func(x int) int {
    return x * 2
}

result := apply(double, 5)  // 返回10
```

---

## 包与初始化

### 包的导入

```go
import (
    "fmt"
    "Golang/4-init/lib1"     // 导入但不使用，只执行init函数
    mylib2 "Golang/4-init/lib2"  // 给包起别名
)
```

### Init函数

- **自动执行**：在包被导入时自动执行，早于main函数
- **无参数**：不能有参数和返回值
- **用途**：进行包初始化

```go
package lib1

func init() {
    fmt.Println("lib1 is initializing...")
}
```

### 导入包的三种方式

1. **标准导入**：`import "fmt"`
   - 使用：`fmt.Println()`

2. **别名导入**：`import f "fmt"`
   - 使用：`f.Println()`

3. **仅执行init**：`import _ "package"`
   - 导入包但不使用其他功能

---

## 指针

指针是一个变量，存储另一个变量的内存地址。

### 指针的基本操作

```go
a := 10
var p *int        // 声明一个指向int的指针
p = &a            // 取址：获取a的地址
fmt.Println(p)    // 输出：内存地址
fmt.Println(*p)   // 解址：获取指针指向的值，输出：10
```

### 指针与函数参数（按引用传递）

**普通参数传递（值传递）**：
```go
func swap(a int, b int) {
    var temp int
    temp = a
    a = b
    b = temp
}

var a, b = 10, 20
swap(a, b)
fmt.Println(a, b)  // 仍然是 10 20，没有交换
```

**指针参数传递（引用传递）**：
```go
func swap(a *int, b *int) {
    var temp int
    temp = *a
    *a = *b
    *b = temp
}

var a, b = 10, 20
swap(&a, &b)
fmt.Println(a, b)  // 输出：20 10，成功交换！
```

### 二级指针

```go
a := 10
p := &a       // 一级指针
pp := &p      // 二级指针，指向一级指针

fmt.Println(pp)   // 一级指针的地址
fmt.Println(&p)   // 同上
fmt.Println(*pp)  // 得到p的值，即a的地址
fmt.Println(*(*pp)) // 得到a的值，即10
```

### 指针的优点

1. **效率**：避免大数据结构的复制
2. **修改**：可以修改原变量的值
3. **灵活性**：支持复杂的数据结构

---

## 延迟执行（Defer）

`defer`关键字用于延迟函数的执行，直到包含该defer语句的函数返回时才执行。

### 基本使用

```go
func returnAndDefer() int {
    defer fmt.Println("defer function called..")
    fmt.Println("return function called..")
    return 0
}

// 输出：
// return function called..
// defer function called..
```

### 栈的执行顺序

多个defer按照后进先出（LIFO）的顺序执行，就像栈一样。

```go
func test() {
    defer fmt.Println("第1个defer")  // 最后执行
    defer fmt.Println("第2个defer")  // 中间执行
    defer fmt.Println("第3个defer")  // 最先执行
}

// 执行顺序：第3个defer -> 第2个defer -> 第1个defer
```

### Defer的常见用途

1. **资源清理**：关闭文件、关闭数据库连接
```go
file, _ := os.Open("test.txt")
defer file.Close()  // 函数返回时自动关闭
```

2. **Panic恢复**：配合recover()处理异常
```go
defer func() {
    if err := recover(); err != nil {
        fmt.Println("捕获到panic:", err)
    }
}()
```

---

## 数组

数组是固定长度的相同类型元素的集合。

### 数组的声明

```go
// 方式1：声明后初始化
arr := [10]int{}
arr = [10]int{1, 2, 3, 4}

// 方式2：声明时初始化
arr := [5]int{1, 2, 3, 4, 5}

// 方式3：让编译器推断长度
arr := [...]int{1, 2, 3, 4}
```

### 数组的特点

- **长度固定**：一旦声明，长度不能改变
- **类型相同**：所有元素类型必须相同
- **值传递**：函数参数传递是值拷贝

### 遍历数组

```go
arr := [5]int{1, 2, 3, 4, 5}

// 方式1：传统for循环
for i := 0; i < len(arr); i++ {
    fmt.Print(arr[i], " ")
}

// 方式2：for-range循环
for index, value := range arr {
    fmt.Println(index, value)
}

// 方式3：只获取值
for _, v := range arr {
    fmt.Print(v, " ")
}
```

### 数组作为函数参数

```go
func change(a [10]int) {
    for i := 0; i < len(a); i++ {
        fmt.Print(a[i], " ")
    }
    a[0] = 100  // 这里改变无效，因为是值传递
}

arr := [10]int{1, 2, 3}
change(arr)
fmt.Println(arr[0])  // 仍然是1
```

---

## 切片

切片（Slice）是动态数组，是对数组的一个连续片段的引用。

### 切片与数组的区别

| 特性 | 数组 | 切片 |
|------|------|------|
| 长度 | 固定 | 动态 |
| 申明 | `[n]type` | `[]type` |
| 内存 | 分配固定空间 | 动态分配 |
| 传递 | 值传递 | 引用传递 |

### 切片的声明

```go
// 方式1：声明一个nil切片
var arr1 []int
fmt.Println(arr1 == nil)  // true

// 方式2：声明时初始化
arr2 := []int{1, 2, 3, 4}

// 方式3：使用make创建
arr3 := make([]int, 3)      // 长度为3，容量为3
arr4 := make([]int, 3, 5)   // 长度为3，容量为5
```

### 长度（len）vs 容量（cap）

```go
s := make([]int, 3, 5)
fmt.Println("len =", len(s), "cap =", cap(s))  // len = 3 cap = 5

s = append(s, 1, 2)  // 追加两个元素
fmt.Println("len =", len(s), "cap =", cap(s))  // len = 5 cap = 5

s = append(s, 3)     // 容量满了，自动扩容
fmt.Println("len =", len(s), "cap =", cap(s))  // len = 6 cap = 10
```

### 切片的切割操作

```go
s1 := []int{1, 2, 3, 4, 5}
s2 := s1[0:2]    // [1, 2]，s2和s1共享内存

s2[0] = 100
fmt.Println(s1)  // [100, 2, 3, 4, 5]，s1也改变了
fmt.Println(s2)  // [100, 2]
```

### 复制切片

```go
s1 := []int{1, 2, 3, 4, 5}
s_c := make([]int, 3)
copy(s_c, s1)    // 复制s1的前3个元素到s_c

s_c[0] = 100
fmt.Println(s1)  // [1, 2, 3, 4, 5]，s1不变
fmt.Println(s_c) // [100, 2, 3]
```

### 删除切片中的元素

```go
x := []int{1, 2, 3, 4, 5}
i := 2  // 删除第2个位置的元素

// 使用copy将后续元素前移
copy(x[i:], x[i+1:])  // x变为 [1, 2, 4, 5, 5]
x = x[:len(x)-1]      // x变为 [1, 2, 4, 5]

fmt.Println(x)
```

### 切片作为函数参数

```go
func test(a []int) {
    a[0] = 100  // 这里的修改会影响原切片
}

arr := []int{1, 2, 3, 4}
test(arr)
fmt.Println(arr)  // [100, 2, 3, 4]，切片是引用传递
```

---

## 映射（Map）

映射是一个无序的键值对集合，类似于其他语言中的字典或哈希表。

### Map的声明

```go
// 方式1：声明后初始化
var myMap1 map[int]string
if myMap1 == nil {
    fmt.Println("myMap1 is nil")
}
myMap1 = make(map[int]string, 3)  // 分配空间
myMap1[1] = "python"
myMap1[2] = "c++"
myMap1[3] = "javascript"

// 方式2：声明时分配空间
myMap2 := make(map[int]string)

// 方式3：声明时直接赋值
myMap3 := map[int]string{
    1: "python",
    2: "Golang",
    3: "javaScript",
}
```

### Map的基本操作

```go
// 添加或更新
myMap[4] = "Golang"

// 访问
fmt.Println(myMap[1])  // python

// 删除
delete(myMap, 1)

// 检查键是否存在
val, ok := myMap[1]
if ok {
    fmt.Println("键存在，值为:", val)
} else {
    fmt.Println("键不存在")
}
```

### 遍历Map

```go
myMap := map[string]int{
    "apple":  5,
    "banana": 3,
    "cherry": 8,
}

for key, value := range myMap {
    fmt.Println(key, "=>", value)
}

// 只获取键
for key := range myMap {
    fmt.Println(key)
}

// 只获取值
for _, value := range myMap {
    fmt.Println(value)
}
```

### Map的特点

- **无序**：键值对的顺序不确定
- **自动扩容**：超过预分配空间时自动扩容
- **键值类型**：键必须能使用 `==` 比较，值可以是任意类型
- **非线程安全**：在并发中需要加锁

---

## 结构体与方法

### 结构体定义

结构体是一种聚合类型，将多个类型的字段组合在一起。

```go
type Book struct {
    title string
    auth  string
}

type Hero struct {
    Name  string
    Ad    int
    Level int
}
```

### 结构体的初始化

```go
// 方式1：按字段顺序初始化
book1 := Book{"Go语言", "Rob Pike"}

// 方式2：按字段名初始化
hero := Hero{
    Name:  "超人",
    Ad:    100,
    Level: 10,
}

// 方式3：使用var声明
var book2 Book
book2.title = "Golang基础"
book2.auth = "Tom"
```

### 访问结构体字段

```go
book := Book{"Go语言", "Rob Pike"}

fmt.Println(book.title)  // Go语言
fmt.Println(book.auth)   // Rob Pike

book.title = "Go高级编程"  // 修改字段
```

### 结构体作为函数参数

```go
func changeBook(book Book) {
    // 值传递，修改无效
    book.auth = "New Author"
}

func changeBook2(book *Book) {
    // 指针传递，修改有效
    book.auth = "New Author"
}

book := Book{"Go语言", "原作者"}

changeBook(book)
fmt.Println(book.auth)   // 仍然是"原作者"

changeBook2(&book)
fmt.Println(book.auth)   // "New Author"
```

---

## 面向对象编程

### 方法（Method）

方法是与某个接收者相关联的函数。

### 1. 值接收者（值传递）

```go
type Hero struct {
    Name  string
    Ad    int
    Level int
}

func (this Hero) Show() {
    fmt.Println("Name is", this.Name)
    fmt.Println("Ad is", this.Ad)
}

hero := Hero{"ljq", 100, 10}
hero.Show()  // 调用方法
```

**问题**：值接收者无法修改结构体的值。

### 2. 指针接收者（引用传递）- 推荐

```go
func (this *Hero) SetName(newName string) {
    this.Name = newName  // 可以修改
}

func (this *Hero) GetName() string {
    return this.Name
}

hero := Hero{"ljq", 100, 10}
hero.SetName("fengnuan")
fmt.Println(hero.GetName())  // fengnuan
```

**优点**：
- 可以修改结构体的值
- 避免大结构体的复制
- 方法调用更高效

### 继承（嵌入结构体）

Go没有传统的继承，但可以通过结构体嵌入实现继承功能。

```go
type Human struct {
    name string
    sex  string
}

func (this *Human) ShowInfo() {
    fmt.Println("name:", this.name)
    fmt.Println("sex:", this.sex)
}

func (this *Human) Walk() {
    fmt.Println(this.name, "is Walking")
}

type SuperMan struct {
    Human  // 嵌入Human结构体
    level  int
}

func (this *SuperMan) ShowInfo() {
    fmt.Println("name:", this.name)
    fmt.Println("sex:", this.sex)
    fmt.Println("level:", this.level)
}

func (this *SuperMan) LevelUp() {
    this.level += 1
}
```

### 使用继承

```go
superMan := SuperMan{
    Human: Human{name: "beauty", sex: "female"},
    level: 10,
}

superMan.ShowInfo()    // 显示所有信息，包括level
superMan.LevelUp()     // 级别提升
superMan.Walk()        // 继承自Human的方法
superMan.Eat()         // 继承自Human的方法
```

### 多重继承（多层嵌入）

```go
type A struct {
    a int
}

type B struct {
    A     // 嵌入A
    b int
}

type C struct {
    B     // 嵌入B
    c int
}

// C拥有A和B的所有字段和方法
```

---

## 输出和格式化

### Print系列函数

```go
// Print - 不换行，不自动加空格
fmt.Print("hello", "Golang")  // 输出：helloGolang

// Println - 换行，自动加空格
fmt.Println("hello", "Golang")  // 输出：hello Golang\n

// Printf - 格式化输出，不换行
fmt.Printf("a = %v\n", 99)  // 输出：a = 99
```

### 常用格式化字符

| 格式 | 说明 |
|------|------|
| `%v` | 值的默认格式 |
| `%T` | 值的类型 |
| `%s` | 字符串 |
| `%d` | 整数 |
| `%f` | 浮点数 |
| `%e` | 科学计数法 |
| `%x` | 十六进制 |

```go
a := 99
b := "I love Go"

fmt.Printf("a = %v\n", a)           // a = 99
fmt.Printf("type of a is %T\n", a)  // type of a is int
fmt.Printf("b = %s\n", b)           // b = I love Go
```

---

## 常见错误和注意事项

### 1. 忘记初始化切片

```go
// 错误
var s []int
s[0] = 100  // panic: index out of range

// 正确
s := make([]int, 1)
s[0] = 100
```

### 2. 忘记初始化Map

```go
// 错误
var m map[string]int
m["key"] = 100  // panic

// 正确
m := make(map[string]int)
m["key"] = 100
```

### 3. 指针的正确使用

```go
// 错误：多次解址
var p *int
**p = 100  // panic

// 正确
a := 10
p := &a
*p = 100
fmt.Println(a)  // 100
```

### 4. 结构体字段的首字母大小写

```go
// 首字母大写：可导出（可在其他包中使用）
type Person struct {
    Name string  // 可导出
    age  int     // 不可导出
}

// 访问
p := Person{Name: "Tom", age: 25}
fmt.Println(p.Name)  // 可以
fmt.Println(p.age)   // 错误：在其他包中无法访问
```

### 5. Defer中的注意事项

```go
// 注意：defer在函数返回前执行，但参数在defer时就已计算
i := 0
defer fmt.Println(i)  // 立即计算i，输出0
i = 100
// 输出：0
```

---

## 最佳实践

### 1. 变量命名规范

```go
// 好的命名
userName := "Tom"
userAge := 25
calculateTotal := func() {}

// 避免
u := "Tom"
a := 25
f := func() {}
```

### 2. 使用指针接收者

```go
// 推荐：如果方法需要修改数据或结构体很大
func (p *Person) SetName(name string) {
    p.name = name
}

// 也可以：如果方法不需要修改数据且结构体很小
func (p Person) GetName() string {
    return p.name
}
```

### 3. 错误处理

```go
// 好的实践
result, err := someFunction()
if err != nil {
    fmt.Println("Error:", err)
    return
}
// 使用result
```

### 4. 使用切片而不是数组

```go
// 避免：数组长度固定，缺乏灵活性
func process(arr [100]int) {}

// 推荐：切片更灵活
func process(arr []int) {}
```

### 5. 优先使用Blank Identifier

```go
// 当不需要某个返回值时
_, err := someFunction()
if err != nil {
    // 处理错误
}
```

---

## 总结

Go语言学习的核心概念：

1. **变量与常量**：理解声明方式和零值规则
2. **函数**：掌握多返回值和函数作为参数
3. **指针**：理解值传递和引用传递的区别
4. **集合类型**：掌握数组、切片和Map的用法
5. **结构体与方法**：理解OOP的实现方式
6. **包管理**：理解包的导入和初始化

### 推荐学习进度

1. 完成变量和常量的基础学习
2. 学习函数和控制流
3. 掌握指针和内存管理
4. 学习集合类型的使用
5. 学习结构体和方法
6. 项目实战应用

### 后续进阶主题

- 接口（Interface）
- 并发（Goroutine 和 Channel）
- 异常处理（error 和 panic）
- 文件操作（I/O）
- 网络编程（HTTP）
- 数据库操作

---

## 快速参考

### 声明速查表

```go
// 变量
var x int                  // 声明
var x int = 10            // 声明并赋值
var x = 10                // 类型推导
x := 10                   // 短声明（仅函数内）

// 常量
const x = 10              // 声明常量

// 数组
var arr [5]int            // 固定长度数组
arr := [5]int{1,2,3,4,5}  // 数组字面量

// 切片
var s []int               // nil切片
s := []int{1,2,3}         // 切片字面量
s := make([]int, 3, 5)    // 长度3，容量5

// Map
var m map[string]int      // nil map
m := make(map[string]int) // 空map

// 结构体
type Person struct {      // 定义
    Name string
    Age  int
}
p := Person{Name: "Tom"} // 初始化

// 函数
func foo(a int) int {}         // 单返回值
func foo() (int, string) {}    // 多返回值
func (p *Person) foo() {}      // 方法
```

祝你Go语言学习愉快！🚀

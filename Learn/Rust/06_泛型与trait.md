# 06 泛型与 trait

> **难度：⭐ 重点难点** | 泛型和 trait 是 Rust 抽象能力的核心。工程中写库、封装服务、定义可替换组件时都会大量使用。

---

## 1. 泛型：把类型作为参数

**标注：核心高频**

泛型允许你编写与具体类型无关的代码。它存在的原因是：真实项目里经常有相同逻辑作用于不同数据类型，如果为每种类型都写一份，会造成重复和维护困难。

```rust
fn first<T>(items: &[T]) -> Option<&T> {
    items.first()
}

fn main() {
    let nums = vec![1, 2, 3];
    let names = vec!["alice", "bob"];

    println!("{:?}", first(&nums));
    println!("{:?}", first(&names));
}
```

工程中常见于仓储层、缓存层、解析器、通用响应结构。

```rust
#[derive(Debug)]
struct ApiResponse<T> {
    code: u16,
    data: T,
}

fn ok<T>(data: T) -> ApiResponse<T> {
    ApiResponse { code: 200, data }
}
```

**为什么这样设计、这样设计带来的好处：**

Rust 泛型默认通过单态化实现：编译器为实际使用到的类型生成专门代码。这样既保留抽象，又避免运行时动态派发开销。

| 方案 | 优点 | 缺点 |
|------|------|------|
| 复制多份函数 | 简单 | 维护成本高 |
| `Box<dyn Any>` | 灵活 | 失去类型安全 |
| 泛型 | 类型安全、零成本 | 编译产物可能变大 |

---

## 2. trait：定义能力边界

**标注：核心高频**

trait 用来描述类型必须具备的行为。它存在的原因是：工程代码需要依赖抽象，而不是依赖具体实现。

```rust
trait Notifier {
    fn send(&self, message: &str) -> Result<(), String>;
}

struct EmailNotifier;

impl Notifier for EmailNotifier {
    fn send(&self, message: &str) -> Result<(), String> {
        println!("send email: {}", message);
        Ok(())
    }
}

fn notify_user<N: Notifier>(notifier: &N, message: &str) -> Result<(), String> {
    notifier.send(message)
}

fn main() -> Result<(), String> {
    let email = EmailNotifier;
    notify_user(&email, "order paid")?;
    Ok(())
}
```

真实工程里，trait 常用于外部服务抽象、存储抽象和策略模式。

**为什么这样设计、这样设计带来的好处：**

trait 把“某个类型是什么”和“它能做什么”分离，调用方只关心能力，不关心具体类型。这让测试替换、模块解耦和库设计更容易。

---

## 3. trait bound：约束泛型能力

**标注：核心高频**

泛型本身不知道 `T` 能做什么，trait bound 用来告诉编译器：这个泛型参数必须实现某些能力。

```rust
use std::fmt::{Debug, Display};

fn log_value<T>(value: T)
where
    T: Display + Debug,
{
    println!("display={}, debug={:?}", value, value);
}
```

**为什么这样设计、这样设计带来的好处：**

Rust 不做“鸭子类型”推断，而是显式声明约束。好处是错误更早、更精确，库的 API 合约也更清晰。

---

## 4. trait 对象：运行时多态

**标注：重要**

当集合中需要放入多种不同类型，但它们都实现同一个 trait 时，可以使用 trait 对象。

```rust
trait Job {
    fn run(&self);
}

struct BackupJob;
struct ReportJob;

impl Job for BackupJob {
    fn run(&self) {
        println!("backup database");
    }
}

impl Job for ReportJob {
    fn run(&self) {
        println!("generate report");
    }
}

fn main() {
    let jobs: Vec<Box<dyn Job>> = vec![Box::new(BackupJob), Box::new(ReportJob)];

    for job in jobs {
        job.run();
    }
}
```

```text
Box<dyn Job>
  ├─ 数据指针：指向具体对象
  └─ vtable：指向 run 等方法的函数表
```

**为什么这样设计、这样设计带来的好处：**

泛型适合编译期确定类型，trait 对象适合运行时选择实现。Rust 同时提供两者，让你在性能和灵活性之间显式取舍。

---

## 5. 关联类型

**标注：重要**

关联类型把 trait 中某些输出类型交给实现者决定。

```rust
trait Repository {
    type Item;

    fn find(&self, id: u64) -> Option<Self::Item>;
}

#[derive(Debug)]
struct User {
    id: u64,
    name: String,
}

struct UserRepo;

impl Repository for UserRepo {
    type Item = User;

    fn find(&self, id: u64) -> Option<Self::Item> {
        Some(User { id, name: "Alice".to_string() })
    }
}
```

**为什么这样设计、这样设计带来的好处：**

关联类型让 trait 的实现更像一个完整协议：实现者不仅提供方法，还声明这个协议中固定的相关类型，减少泛型参数在调用链中到处传播。

---

## 动手练习

### 练习 1：通用分页响应

为 Web API 设计一个 `Page<T>` 泛型结构，包含列表、页码、每页数量、总数。

**思路提示：** 先定义结构体，再写一个 `new` 构造函数。注意 `T` 不需要额外 trait bound。

### 练习 2：通知服务抽象

定义 `Notifier` trait，并实现 `EmailNotifier` 和 `SmsNotifier`。

**思路提示：** 让业务函数依赖 `impl Notifier`，不要直接依赖具体结构体。

### 练习 3：任务调度器

用 `Vec<Box<dyn Job>>` 存放不同任务，并统一执行。

**思路提示：** 每个任务结构体实现 `Job`，调度器只持有 trait 对象。

### 练习 4：仓储层关联类型

定义 `Repository` trait，为 `UserRepo` 和 `OrderRepo` 分别实现不同的 `Item`。

**思路提示：** 使用 `type Item` 表达每个仓储返回的实体类型。

---

**下一章：** [07_生命周期](./07_生命周期.md)

# 13. Unsafe Rust 与 FFI（外部函数接口）

> **难度**: 高级（不常用）—— 日常开发中很少直接使用，但理解它对于编写高性能库、与 C/C++ 交互至关重要。

## 目录

1. [什么是 Unsafe Rust](#1-什么是-unsafe-rust)
2. [Unsafe 允许的五种操作](#2-unsafe-允许的五种操作)
3. [何时使用 Unsafe](#3-何时使用-unsafe)
4. [Unsafe 最佳实践](#4-unsafe-最佳实践)
5. [FFI 基础：与 C 交互](#5-ffi-基础与-c-交互)
6. [工具链：bindgen 与 cbindgen](#6-工具链bindgen-与-cbindgen)
7. [从其他语言调用 Rust](#7-从其他语言调用-rust)
8. [练习题](#8-练习题)

---

## 1. 什么是 Unsafe Rust

Rust 的安全保证建立在编译器的静态分析之上，但有些操作编译器无法验证其安全性。
`unsafe` 关键字告诉编译器："我（程序员）已经手动验证了这段代码的安全性。"

```rust
// 安全的 Rust —— 编译器保证内存安全
let x = 42;
let r = &x;

// Unsafe Rust —— 程序员承诺安全性
unsafe {
    let raw = &x as *const i32;
    println!("raw pointer value: {}", *raw);
}
```

### 为什么这样设计？

Rust 的设计哲学是 **零成本抽象** + **安全优先**：
- 绝大多数代码在 Safe Rust 中完成，享受编译器保护
- 少量底层操作（硬件交互、性能优化、FFI）需要绕过检查
- `unsafe` 块明确标记了"信任边界"，方便代码审查时重点关注
- 这比 C/C++ 的"处处 unsafe"要好得多——bug 的搜索范围大大缩小

---

## 2. Unsafe 允许的五种操作

`unsafe` 块中可以执行以下五种在 Safe Rust 中被禁止的操作：

### 2.1 解引用裸指针（Raw Pointers）

```rust
fn main() {
    let mut num = 5;

    // 创建裸指针是安全的，解引用才需要 unsafe
    let r1 = &num as *const i32;      // 不可变裸指针
    let r2 = &mut num as *mut i32;    // 可变裸指针

    unsafe {
        println!("r1 = {}", *r1);
        println!("r2 = {}", *r2);
    }
}
```

**裸指针 vs 引用的区别：**
- 允许同时存在可变和不可变指针指向同一位置
- 不保证指向有效内存
- 允许为 null
- 没有自动清理机制

### 2.2 调用 Unsafe 函数

```rust
unsafe fn dangerous() {
    // 这个函数的实现可能违反内存安全
    // 调用者必须确保前置条件满足
}

fn main() {
    unsafe {
        dangerous();
    }
}
```

**真实案例：`slice::from_raw_parts`**

```rust
use std::slice;

/// 将一个指针和长度转换为切片
/// 安全前提：ptr 必须有效，且 ptr..ptr+len 范围内的内存已初始化
fn safe_slice_from_raw(ptr: *const i32, len: usize) -> &'static [i32] {
    unsafe { slice::from_raw_parts(ptr, len) }
}
```

### 2.3 访问或修改可变静态变量

```rust
static mut COUNTER: u32 = 0;

fn add_to_counter(inc: u32) {
    unsafe {
        COUNTER += inc;
    }
}

fn main() {
    add_to_counter(3);
    unsafe {
        println!("COUNTER: {}", COUNTER);
    }
}
```

> **警告**：可变静态变量在多线程环境下极易产生数据竞争，优先使用 `AtomicU32` 或 `Mutex`。

### 2.4 实现 Unsafe Trait

```rust
/// 标记类型可以安全地在线程间传递
/// 编译器无法自动验证，需要程序员保证
unsafe trait MySync {
    // ...
}

/// 实现者承诺：该类型确实可以安全地跨线程共享
unsafe impl MySync for MyType {
    // ...
}
```

**真实案例**：`Send` 和 `Sync` 就是 unsafe trait，因为编译器无法验证自定义类型的线程安全性。

### 2.5 访问 Union 的字段

```rust
#[repr(C)]
union MyUnion {
    i: i32,
    f: f32,
}

fn main() {
    let u = MyUnion { i: 42 };
    // 访问 union 字段需要 unsafe，因为编译器不知道当前存储的是哪个变体
    unsafe {
        println!("as int: {}", u.i);
        println!("as float: {}", u.f); // 可能是无意义的值
    }
}
```

---

## 3. 何时使用 Unsafe

### 合理的使用场景

| 场景 | 示例 |
|------|------|
| 性能关键路径 | 跳过边界检查 `get_unchecked()` |
| FFI 调用 | 调用 C 库函数 |
| 实现底层数据结构 | 链表、B树等需要裸指针 |
| 硬件交互 | 内存映射 I/O |
| 编译器无法证明的安全模式 | 自引用结构体 |

### 不应使用 Unsafe 的场景

- 仅仅为了"方便"绕过借用检查器 —— 重新设计数据结构
- 全局可变状态 —— 使用 `Mutex` 或 `RwLock`
- 不确定是否安全时 —— 如果不能证明安全，就不要用

---

## 4. Unsafe 最佳实践

### 4.1 最小化 unsafe 块

```rust
// 差：整个函数都是 unsafe
unsafe fn bad_example(data: *const u8, len: usize) -> Vec<u8> {
    let slice = std::slice::from_raw_parts(data, len);
    slice.to_vec()
}

// 好：只在必要的地方使用 unsafe
fn good_example(data: *const u8, len: usize) -> Vec<u8> {
    let slice = unsafe { std::slice::from_raw_parts(data, len) };
    slice.to_vec()  // 这行不需要 unsafe
}
```

### 4.2 封装为安全接口

```rust
pub struct SafeBuffer {
    ptr: *mut u8,
    len: usize,
    cap: usize,
}

impl SafeBuffer {
    /// 创建新缓冲区
    pub fn new(cap: usize) -> Self {
        let layout = std::alloc::Layout::array::<u8>(cap).unwrap();
        let ptr = unsafe { std::alloc::alloc(layout) };
        if ptr.is_null() {
            std::alloc::handle_alloc_error(layout);
        }
        SafeBuffer { ptr, len: 0, cap }
    }

    /// 安全的写入方法——内部使用 unsafe，但对外接口是安全的
    pub fn push(&mut self, byte: u8) {
        assert!(self.len < self.cap, "buffer full");
        unsafe {
            self.ptr.add(self.len).write(byte);
        }
        self.len += 1;
    }

    /// 安全的读取方法
    pub fn as_slice(&self) -> &[u8] {
        unsafe { std::slice::from_raw_parts(self.ptr, self.len) }
    }
}

impl Drop for SafeBuffer {
    fn drop(&mut self) {
        let layout = std::alloc::Layout::array::<u8>(self.cap).unwrap();
        unsafe { std::alloc::dealloc(self.ptr, layout); }
    }
}
```

### 4.3 文档化安全不变量

```rust
/// 从裸指针创建字符串切片
///
/// # Safety
///
/// 调用者必须确保：
/// - `ptr` 指向有效的 UTF-8 编码的字节序列
/// - `ptr` 指向的内存在返回的 `&str` 生命周期内保持有效
/// - `len` 不超过分配的内存大小
pub unsafe fn str_from_raw(ptr: *const u8, len: usize) -> &'static str {
    let slice = std::slice::from_raw_parts(ptr, len);
    std::str::from_utf8_unchecked(slice)
}
```

---

## 5. FFI 基础：与 C 交互

### 5.1 从 Rust 调用 C 函数

```rust
// 声明外部 C 函数
extern "C" {
    fn abs(input: i32) -> i32;
    fn strlen(s: *const std::os::raw::c_char) -> usize;
}

fn main() {
    unsafe {
        println!("abs(-5) = {}", abs(-5));
    }
}
```

### 5.2 C 类型映射

| C 类型 | Rust 类型 |
|--------|-----------|
| `int` | `c_int` (即 `i32`) |
| `unsigned int` | `c_uint` (即 `u32`) |
| `char*` | `*const c_char` 或 `*mut c_char` |
| `void*` | `*mut c_void` |
| `size_t` | `usize` |
| `bool` | `bool` |

### 5.3 完整的 C 库调用示例

```rust
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

extern "C" {
    fn getenv(name: *const c_char) -> *const c_char;
}

/// 安全封装：获取环境变量
fn safe_getenv(name: &str) -> Option<String> {
    let c_name = CString::new(name).ok()?;
    unsafe {
        let ptr = getenv(c_name.as_ptr());
        if ptr.is_null() {
            None
        } else {
            Some(CStr::from_ptr(ptr).to_string_lossy().into_owned())
        }
    }
}

fn main() {
    match safe_getenv("PATH") {
        Some(path) => println!("PATH = {}", path),
        None => println!("PATH not set"),
    }
}
```

### 为什么这样设计？

FFI 的设计体现了 Rust 的核心原则：
- **extern "C"** 指定使用 C 的调用约定（ABI），确保二进制兼容
- 所有 FFI 调用都是 `unsafe` 的，因为 Rust 无法验证外部代码的正确性
- `CString`/`CStr` 处理 Rust 字符串（UTF-8, 无 null 终止）和 C 字符串（null 终止）的差异

---

## 6. 工具链：bindgen 与 cbindgen

### 6.1 bindgen：自动生成 C → Rust 绑定

```toml
# Cargo.toml
[build-dependencies]
bindgen = "0.69"
```

```rust
// build.rs
use bindgen;

fn main() {
    let bindings = bindgen::Builder::default()
        .header("wrapper.h")
        .parse_callbacks(Box::new(bindgen::CargoCallbacks::new()))
        .generate()
        .expect("Unable to generate bindings");

    let out_path = std::path::PathBuf::from(std::env::var("OUT_DIR").unwrap());
    bindings
        .write_to_file(out_path.join("bindings.rs"))
        .expect("Couldn't write bindings!");
}
```

```rust
// src/lib.rs
#![allow(non_upper_case_globals)]
#![allow(non_camel_case_types)]
include!(concat!(env!("OUT_DIR"), "/bindings.rs"));
```

### 6.2 cbindgen：生成 Rust → C 头文件

```toml
# cbindgen.toml
language = "C"
header = "/* Generated by cbindgen */"
include_guard = "MY_LIB_H"
```

```bash
# 生成头文件
cbindgen --config cbindgen.toml --crate my_lib --output my_lib.h
```

---

## 7. 从其他语言调用 Rust

### 7.1 暴露 Rust 函数给 C

```rust
// src/lib.rs

/// 计算斐波那契数列第 n 项
/// 使用 #[no_mangle] 防止 Rust 修改函数名
/// 使用 extern "C" 确保 C ABI 兼容
#[no_mangle]
pub extern "C" fn fibonacci(n: u32) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => {
            let (mut a, mut b) = (0u64, 1u64);
            for _ in 2..=n {
                let temp = b;
                b = a + b;
                a = temp;
            }
            b
        }
    }
}

/// 释放 Rust 分配的字符串
#[no_mangle]
pub extern "C" fn free_rust_string(ptr: *mut c_char) {
    if !ptr.is_null() {
        unsafe { drop(CString::from_raw(ptr)); }
    }
}
```

```toml
# Cargo.toml
[lib]
crate-type = ["cdylib"]  # 生成动态链接库 (.so / .dll / .dylib)
```

### 7.2 从 Python 调用 Rust（使用 PyO3）

```toml
# Cargo.toml
[dependencies]
pyo3 = { version = "0.21", features = ["extension-module"] }

[lib]
crate-type = ["cdylib"]
name = "my_rust_lib"
```

```rust
use pyo3::prelude::*;

/// 高性能的字符串处理函数
#[pyfunction]
fn count_words(text: &str) -> usize {
    text.split_whitespace().count()
}

/// 斐波那契计算——比纯 Python 快 100 倍以上
#[pyfunction]
fn fibonacci_py(n: u32) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => {
            let (mut a, mut b) = (0u64, 1u64);
            for _ in 2..=n {
                let temp = b;
                b = a + b;
                a = temp;
            }
            b
        }
    }
}

/// Python 模块定义
#[pymodule]
fn my_rust_lib(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(count_words, m)?)?;
    m.add_function(wrap_pyfunction!(fibonacci_py, m)?)?;
    Ok(())
}
```

```python
# Python 中使用
import my_rust_lib

print(my_rust_lib.count_words("Hello World from Rust"))  # 4
print(my_rust_lib.fibonacci_py(50))  # 12586269025
```

---

## 8. 练习题

### 练习 1：实现安全的 split_at_mut（中等）

标准库的 `split_at_mut` 内部使用了 unsafe。尝试自己实现它：

```rust
fn my_split_at_mut(slice: &mut [i32], mid: usize) -> (&mut [i32], &mut [i32]) {
    let len = slice.len();
    assert!(mid <= len);
    // TODO: 使用 unsafe 实现
    // 提示：使用 slice.as_mut_ptr() 和 slice::from_raw_parts_mut
    todo!()
}
```

<details>
<summary>提示</summary>

```rust
fn my_split_at_mut(slice: &mut [i32], mid: usize) -> (&mut [i32], &mut [i32]) {
    let len = slice.len();
    let ptr = slice.as_mut_ptr();
    assert!(mid <= len);
    unsafe {
        (
            std::slice::from_raw_parts_mut(ptr, mid),
            std::slice::from_raw_parts_mut(ptr.add(mid), len - mid),
        )
    }
}
```
</details>

### 练习 2：调用 C 数学库（高级）

使用 FFI 调用 C 标准库的数学函数，并封装为安全接口：

```rust
// TODO: 声明 extern "C" 块，包含 sin, cos, sqrt
// TODO: 创建安全的封装函数
// TODO: 编写测试验证结果正确性
```

<details>
<summary>提示</summary>

```rust
extern "C" {
    fn sin(x: f64) -> f64;
    fn cos(x: f64) -> f64;
    fn sqrt(x: f64) -> f64;
}

pub fn safe_sqrt(x: f64) -> Option<f64> {
    if x < 0.0 {
        None
    } else {
        Some(unsafe { sqrt(x) })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn test_sqrt() {
        assert_eq!(safe_sqrt(4.0), Some(2.0));
        assert_eq!(safe_sqrt(-1.0), None);
    }
}
```
</details>

### 练习 3：暴露 Rust 函数给 Python（高级）

使用 PyO3 创建一个 Python 模块，提供以下功能：
- `is_palindrome(s: str) -> bool` —— 判断字符串是否为回文
- `prime_sieve(n: int) -> list[int]` —— 返回小于 n 的所有素数

<details>
<summary>提示</summary>

1. 创建项目：`cargo new --lib rust_utils`
2. 在 `Cargo.toml` 中添加 `pyo3` 依赖
3. 设置 `crate-type = ["cdylib"]`
4. 实现函数并用 `#[pyfunction]` 标记
5. 用 `maturin develop` 构建并安装到 Python 环境
</details>

---

## 总结

| 概念 | 关键点 |
|------|--------|
| unsafe 块 | 最小化范围，文档化安全不变量 |
| 裸指针 | 创建安全，解引用需要 unsafe |
| FFI | extern "C" + unsafe，注意类型映射 |
| bindgen/cbindgen | 自动生成绑定，减少手写错误 |
| PyO3 | Rust → Python 的最佳方案 |

> **下一章**: [14_宏编程.md](./14_宏编程.md) —— 学习 Rust 强大的元编程能力

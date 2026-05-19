# Rust Web 开发

> **难度标记：实战** | 需要前置知识：异步编程、错误处理、trait、生命周期

## 为什么用 Rust 做 Web 开发？

- **性能**：接近 C/C++ 的运行速度，远超 Go/Java/Node.js
- **内存安全**：无 GC 停顿，无数据竞争，适合高并发服务
- **类型系统**：编译期捕获大量 bug，重构有信心
- **低资源占用**：单个二进制文件，Docker 镜像可以小到 10MB

真实案例：Cloudflare、Discord、Dropbox 都在生产环境使用 Rust Web 服务。

---

## 1. Axum 框架

Axum 是 Tokio 团队开发的 Web 框架，设计哲学是：**利用 Rust 类型系统，让错误在编译期暴露**。

### 1.1 最小示例

```rust
use axum::{routing::get, Router};

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(|| async { "Hello, World!" }));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}
```

**Cargo.toml 依赖：**

```toml
[dependencies]
axum = "0.7"
tokio = { version = "1", features = ["full"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
```

### 1.2 路由（Routing）

```rust
use axum::{routing::{get, post, put, delete}, Router};

fn app() -> Router {
    Router::new()
        .route("/users", get(list_users).post(create_user))
        .route("/users/:id", get(get_user).put(update_user).delete(delete_user))
        .nest("/api/v1", api_routes()) // 路由嵌套
        .fallback(handler_404)        // 404 处理
}

fn api_routes() -> Router {
    Router::new()
        .route("/health", get(health_check))
        .route("/posts", get(list_posts).post(create_post))
}
```

> **为什么选择 Axum 而不是 Actix-web？**
> Axum 与 Tokio 生态深度集成，使用 Tower 中间件（与 tonic gRPC 共享），
> 类型提取器设计更符合 Rust 惯用风格。Actix-web 性能略高但生态独立。

### 1.3 处理器（Handlers）

Handler 就是一个异步函数，参数是"提取器"，返回值实现了 `IntoResponse`：

```rust
use axum::{extract::Path, http::StatusCode, Json};
use serde::{Deserialize, Serialize};

#[derive(Serialize)]
struct User {
    id: u64,
    name: String,
    email: String,
}

// GET /users/:id
async fn get_user(Path(id): Path<u64>) -> Result<Json<User>, StatusCode> {
    // 模拟数据库查询
    if id == 0 {
        return Err(StatusCode::NOT_FOUND);
    }
    Ok(Json(User {
        id,
        name: "Alice".to_string(),
        email: "alice@example.com".to_string(),
    }))
}

// POST /users
#[derive(Deserialize)]
struct CreateUser {
    name: String,
    email: String,
}

async fn create_user(Json(payload): Json<CreateUser>) -> (StatusCode, Json<User>) {
    let user = User {
        id: 1, // 实际应由数据库生成
        name: payload.name,
        email: payload.email,
    };
    (StatusCode::CREATED, Json(user))
}
```

### 1.4 提取器（Extractors）

提取器是 Axum 的核心设计——**从请求中类型安全地提取数据**：

```rust
use axum::extract::{Path, Query, State, Json};
use std::collections::HashMap;

// 路径参数
async fn get_user(Path(id): Path<u64>) -> String {
    format!("User {}", id)
}

// 多个路径参数
async fn get_post(Path((user_id, post_id)): Path<(u64, u64)>) -> String {
    format!("User {} Post {}", user_id, post_id)
}

// 查询参数: GET /search?q=rust&page=1
#[derive(Deserialize)]
struct SearchParams {
    q: String,
    page: Option<u32>,
}

async fn search(Query(params): Query<SearchParams>) -> String {
    format!("搜索: {}, 页码: {}", params.q, params.page.unwrap_or(1))
}

// 请求头
use axum::http::HeaderMap;
async fn headers(headers: HeaderMap) -> String {
    let ua = headers.get("user-agent")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("unknown");
    format!("User-Agent: {}", ua)
}
```

> **为什么提取器这么好用？**
> 传统框架（Express/Flask）需要手动从 request 对象取值并转换类型，
> Axum 的提取器在编译期就确保了类型正确，运行时自动反序列化并处理错误。

### 1.5 应用状态（State）

```rust
use axum::extract::State;
use std::sync::Arc;
use tokio::sync::RwLock;

// 共享状态
#[derive(Clone)]
struct AppState {
    db: sqlx::PgPool,
    redis: redis::Client,
}

// 或者用 Arc<RwLock<T>> 做可变状态
type SharedState = Arc<RwLock<Vec<User>>>;

async fn list_users(State(state): State<AppState>) -> Json<Vec<User>> {
    let users = sqlx::query_as!(User, "SELECT * FROM users")
        .fetch_all(&state.db)
        .await
        .unwrap();
    Json(users)
}

#[tokio::main]
async fn main() {
    let state = AppState {
        db: sqlx::PgPool::connect("postgres://localhost/mydb").await.unwrap(),
        redis: redis::Client::open("redis://127.0.0.1/").unwrap(),
    };

    let app = Router::new()
        .route("/users", get(list_users))
        .with_state(state);

    // ...启动服务器
}
```

### 1.6 中间件（Middleware）

```rust
use axum::{middleware, extract::Request, http::StatusCode};
use axum::response::Response;

// 自定义中间件：请求计时
async fn timing_middleware(
    req: Request,
    next: middleware::Next,
) -> Response {
    let start = std::time::Instant::now();
    let path = req.uri().path().to_string();

    let response = next.run(req).await;

    let duration = start.elapsed();
    tracing::info!("{} took {:?}", path, duration);
    response
}

// 认证中间件
async fn auth_middleware(
    headers: HeaderMap,
    req: Request,
    next: middleware::Next,
) -> Result<Response, StatusCode> {
    let token = headers.get("Authorization")
        .and_then(|v| v.to_str().ok())
        .ok_or(StatusCode::UNAUTHORIZED)?;

    if !verify_token(token) {
        return Err(StatusCode::UNAUTHORIZED);
    }

    Ok(next.run(req).await)
}

// 应用中间件
let app = Router::new()
    .route("/public", get(public_handler))
    .route("/protected", get(protected_handler))
    .route_layer(middleware::from_fn(auth_middleware)) // 仅对上面的路由生效
    .layer(middleware::from_fn(timing_middleware));    // 对所有路由生效
```

---

## 2. Actix-web 概览

Actix-web 是另一个主流 Rust Web 框架，以极致性能著称：

```rust
use actix_web::{web, App, HttpServer, HttpResponse};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(|| async { HttpResponse::Ok().body("Hello") }))
            .service(
                web::scope("/api")
                    .route("/users", web::get().to(list_users))
            )
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
```

**Axum vs Actix-web 对比：**

| 特性 | Axum | Actix-web |
|------|------|-----------|
| 异步运行时 | Tokio | Tokio (也支持其他) |
| 中间件 | Tower (通用) | 自有系统 |
| 性能 | 极高 | 极高（略胜） |
| 生态集成 | Tokio 全家桶 | 独立生态 |
| 学习曲线 | 中等 | 中等 |
| 类型安全 | 更强 | 强 |

---

## 3. 数据库操作

### 3.1 SQLx —— 编译期检查的 SQL

SQLx 的杀手特性：**在编译时验证你的 SQL 语句是否正确**。

```rust
use sqlx::{PgPool, FromRow};

#[derive(Debug, FromRow, Serialize)]
struct User {
    id: i64,
    name: String,
    email: String,
    created_at: chrono::NaiveDateTime,
}

// 编译期检查：如果表/列不存在，编译失败！
async fn get_user(pool: &PgPool, id: i64) -> Result<User, sqlx::Error> {
    sqlx::query_as!(User, "SELECT id, name, email, created_at FROM users WHERE id = $1", id)
        .fetch_one(pool)
        .await
}

// 插入并返回
async fn create_user(pool: &PgPool, name: &str, email: &str) -> Result<User, sqlx::Error> {
    sqlx::query_as!(
        User,
        r#"INSERT INTO users (name, email) VALUES ($1, $2)
           RETURNING id, name, email, created_at"#,
        name, email
    )
    .fetch_one(pool)
    .await
}

// 事务
async fn transfer(pool: &PgPool, from: i64, to: i64, amount: f64) -> Result<(), sqlx::Error> {
    let mut tx = pool.begin().await?;

    sqlx::query!("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
        .execute(&mut *tx).await?;
    sqlx::query!("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
        .execute(&mut *tx).await?;

    tx.commit().await
}
```

**数据库迁移：**

```bash
# 安装 sqlx-cli
cargo install sqlx-cli

# 创建迁移
sqlx migrate add create_users_table

# 运行迁移
sqlx migrate run

# 离线模式（CI 中无需数据库连接）
cargo sqlx prepare
```

> **为什么 SQLx 的编译期检查如此重要？**
> 传统 ORM 的 SQL 错误只在运行时暴露（可能是上线后）。
> SQLx 在 `cargo build` 时就连接数据库验证 SQL，把运行时错误变成编译错误。

### 3.2 Diesel ORM

Diesel 是类型安全的 ORM，适合喜欢 ActiveRecord 风格的开发者：

```rust
// schema.rs (由 diesel CLI 自动生成)
diesel::table! {
    users (id) {
        id -> Int4,
        name -> Varchar,
        email -> Varchar,
    }
}

// models.rs
#[derive(Queryable, Selectable)]
#[diesel(table_name = users)]
struct User {
    id: i32,
    name: String,
    email: String,
}

#[derive(Insertable)]
#[diesel(table_name = users)]
struct NewUser<'a> {
    name: &'a str,
    email: &'a str,
}
```

---

## 4. API 设计模式

### 4.1 统一错误响应

```rust
use axum::{http::StatusCode, response::{IntoResponse, Response}, Json};
use serde_json::json;

// 自定义错误类型
#[derive(Debug)]
enum AppError {
    NotFound(String),
    BadRequest(String),
    Unauthorized,
    Internal(anyhow::Error),
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg),
            AppError::BadRequest(msg) => (StatusCode::BAD_REQUEST, msg),
            AppError::Unauthorized => (StatusCode::UNAUTHORIZED, "未授权".to_string()),
            AppError::Internal(e) => {
                tracing::error!("内部错误: {:?}", e);
                (StatusCode::INTERNAL_SERVER_ERROR, "服务器内部错误".to_string())
            }
        };

        let body = json!({
            "error": {
                "code": status.as_u16(),
                "message": message,
            }
        });

        (status, Json(body)).into_response()
    }
}

// 让 ? 运算符自动转换错误
impl From<sqlx::Error> for AppError {
    fn from(e: sqlx::Error) -> Self {
        match e {
            sqlx::Error::RowNotFound => AppError::NotFound("资源不存在".to_string()),
            _ => AppError::Internal(e.into()),
        }
    }
}

// 使用：handler 返回 Result<T, AppError>
async fn get_user(
    State(pool): State<PgPool>,
    Path(id): Path<i64>,
) -> Result<Json<User>, AppError> {
    let user = sqlx::query_as!(User, "SELECT * FROM users WHERE id = $1", id)
        .fetch_one(&pool)
        .await?; // 自动转换为 AppError
    Ok(Json(user))
}
```

### 4.2 JWT 认证中间件

```rust
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    sub: String,      // 用户 ID
    exp: usize,       // 过期时间
    role: String,     // 角色
}

// 生成 token
fn create_token(user_id: &str, role: &str, secret: &[u8]) -> Result<String, AppError> {
    let expiration = chrono::Utc::now()
        .checked_add_signed(chrono::Duration::hours(24))
        .unwrap()
        .timestamp() as usize;

    let claims = Claims {
        sub: user_id.to_string(),
        exp: expiration,
        role: role.to_string(),
    };

    encode(&Header::default(), &claims, &EncodingKey::from_secret(secret))
        .map_err(|e| AppError::Internal(e.into()))
}

// 提取器：从请求中提取认证信息
struct AuthUser(Claims);

#[axum::async_trait]
impl<S> axum::extract::FromRequestParts<S> for AuthUser
where S: Send + Sync
{
    type Rejection = AppError;

    async fn from_request_parts(
        parts: &mut axum::http::request::Parts,
        _state: &S,
    ) -> Result<Self, Self::Rejection> {
        let token = parts.headers
            .get("Authorization")
            .and_then(|v| v.to_str().ok())
            .and_then(|v| v.strip_prefix("Bearer "))
            .ok_or(AppError::Unauthorized)?;

        let secret = std::env::var("JWT_SECRET").unwrap_or_else(|_| "secret".into());
        let token_data = decode::<Claims>(
            token,
            &DecodingKey::from_secret(secret.as_bytes()),
            &Validation::default(),
        ).map_err(|_| AppError::Unauthorized)?;

        Ok(AuthUser(token_data.claims))
    }
}

// 使用：只需在 handler 参数中加入 AuthUser
async fn protected_route(AuthUser(claims): AuthUser) -> String {
    format!("欢迎, 用户 {}! 角色: {}", claims.sub, claims.role)
}
```

### 4.3 限流中间件

```rust
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;
use std::time::{Duration, Instant};

#[derive(Clone)]
struct RateLimiter {
    requests: Arc<Mutex<HashMap<String, Vec<Instant>>>>,
    max_requests: usize,
    window: Duration,
}

impl RateLimiter {
    fn new(max_requests: usize, window: Duration) -> Self {
        Self {
            requests: Arc::new(Mutex::new(HashMap::new())),
            max_requests,
            window,
        }
    }

    async fn check(&self, key: &str) -> bool {
        let mut map = self.requests.lock().await;
        let now = Instant::now();
        let entries = map.entry(key.to_string()).or_default();

        // 清除过期记录
        entries.retain(|&t| now.duration_since(t) < self.window);

        if entries.len() >= self.max_requests {
            false
        } else {
            entries.push(now);
            true
        }
    }
}

async fn rate_limit_middleware(
    State(limiter): State<RateLimiter>,
    req: Request,
    next: middleware::Next,
) -> Result<Response, StatusCode> {
    let ip = req.headers()
        .get("x-forwarded-for")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("unknown")
        .to_string();

    if !limiter.check(&ip).await {
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    Ok(next.run(req).await)
}
```

---

## 5. 部署

### 5.1 Docker 多阶段构建

```dockerfile
# 构建阶段
FROM rust:1.77 as builder
WORKDIR /app
COPY . .
RUN cargo build --release

# 运行阶段（极小镜像）
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/myapp /usr/local/bin/
EXPOSE 3000
CMD ["myapp"]
```

### 5.2 配置管理

```rust
use serde::Deserialize;

#[derive(Deserialize)]
struct Config {
    database_url: String,
    jwt_secret: String,
    port: u16,
    #[serde(default = "default_log_level")]
    log_level: String,
}

fn default_log_level() -> String { "info".to_string() }

impl Config {
    fn from_env() -> Result<Self, envy::Error> {
        // 从环境变量加载（支持 .env 文件）
        dotenvy::dotenv().ok();
        envy::from_env()
    }
}
```

### 5.3 健康检查与优雅关闭

```rust
use tokio::signal;

async fn shutdown_signal() {
    let ctrl_c = async { signal::ctrl_c().await.unwrap() };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .unwrap().recv().await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
    tracing::info!("收到关闭信号，正在优雅关闭...");
}

#[tokio::main]
async fn main() {
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}
```

---

## 6. 完整项目结构

```
my-api/
├── Cargo.toml
├── .env
├── migrations/
│   └── 001_create_users.sql
├── src/
│   ├── main.rs          # 入口，启动服务器
│   ├── config.rs        # 配置加载
│   ├── routes/
│   │   ├── mod.rs
│   │   ├── users.rs
│   │   └── auth.rs
│   ├── models/
│   │   ├── mod.rs
│   │   └── user.rs
│   ├── middleware/
│   │   ├── mod.rs
│   │   ├── auth.rs
│   │   └── rate_limit.rs
│   └── errors.rs        # 统一错误处理
└── tests/
    └── api_tests.rs     # 集成测试
```

---

## 练习题

### 练习 1：构建 CRUD REST API（难度：⭐⭐⭐）

构建一个完整的待办事项（Todo）API：
- `GET /todos` - 列出所有待办
- `POST /todos` - 创建待办
- `GET /todos/:id` - 获取单个待办
- `PUT /todos/:id` - 更新待办
- `DELETE /todos/:id` - 删除待办
- 使用 SQLx + PostgreSQL 持久化
- 实现分页（`?page=1&per_page=20`）
- 统一错误响应格式

### 练习 2：添加 JWT 认证（难度：⭐⭐⭐⭐）

在练习 1 基础上：
- `POST /register` - 用户注册（密码用 argon2 哈希）
- `POST /login` - 登录返回 JWT token
- 保护 CRUD 路由，只有登录用户能操作自己的 todo
- 实现 token 刷新机制
- 添加角色系统（admin 可以看所有人的 todo）

### 练习 3：实现限流中间件（难度：⭐⭐⭐⭐）

- 实现滑动窗口限流算法
- 支持按 IP 和按用户两种维度
- 返回 `X-RateLimit-Remaining` 和 `X-RateLimit-Reset` 响应头
- 不同路由可配置不同限流策略
- 用 Redis 替代内存存储（支持多实例部署）

---

## 推荐学习资源

- [Axum 官方示例](https://github.com/tokio-rs/axum/tree/main/examples)
- [Zero To Production In Rust](https://www.zero2prod.com/) - 最佳 Rust Web 开发书籍
- [Shuttle.rs](https://shuttle.rs/) - Rust 云部署平台
- [Loco.rs](https://loco.rs/) - 类 Rails 的 Rust Web 框架（快速开发）

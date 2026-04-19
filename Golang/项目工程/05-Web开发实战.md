# Web 开发实战

> **目标**：从标准库 net/http 到 Gin 框架，掌握 Go Web 开发的完整链路：路由、中间件、请求处理、JSON、文件上传、鉴权。

---

## 一、net/http 标准库

### 1.1 最小服务器

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
    })

    log.Println("server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 1.2 Handler 接口

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// 任何实现了 ServeHTTP 的类型都是 Handler
type apiHandler struct{}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
}

http.Handle("/api", &apiHandler{})
```

### 1.3 ServeMux（路由器）

```go
mux := http.NewServeMux()

// Go 1.22+ 支持方法和路径参数
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)

func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // Go 1.22+
    // ...
}

server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
log.Fatal(server.ListenAndServe())
```

### 1.4 中间件模式

```go
type Middleware func(http.Handler) http.Handler

func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w, "Internal Server Error", 500)
                log.Printf("panic: %v", err)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// 组合中间件
handler := logging(recovery(mux))
http.ListenAndServe(":8080", handler)
```

---

## 二、JSON 处理

### 2.1 结构体标签

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"` // 空值时省略
    Password  string    `json:"-"`                // 永不序列化
    CreatedAt time.Time `json:"created_at"`
}
```

### 2.2 编码解码

```go
// 序列化：struct → JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// 反序列化：JSON → struct
func readJSON(r *http.Request, dst any) error {
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields() // 严格模式
    return decoder.Decode(dst)
}

// 使用
func createUser(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := readJSON(r, &input); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    // 创建用户...
    writeJSON(w, 201, map[string]string{"status": "created"})
}
```

---

## 三、Gin 框架

### 3.1 安装

```bash
go get -u github.com/gin-gonic/gin
```

### 3.2 基本使用

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default() // 自带 Logger + Recovery 中间件

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

### 3.3 路由与分组

```go
r := gin.Default()

// 路由参数
r.GET("/users/:id", getUser)        // /users/123
r.GET("/files/*filepath", getFile)  // /files/path/to/file

// 路由分组
api := r.Group("/api/v1")
{
    api.GET("/users", listUsers)
    api.POST("/users", createUser)
    api.GET("/users/:id", getUser)
    api.PUT("/users/:id", updateUser)
    api.DELETE("/users/:id", deleteUser)
}

// 带中间件的分组
admin := r.Group("/admin", authMiddleware())
{
    admin.GET("/stats", getStats)
}
```

### 3.4 请求参数获取

```go
func getUser(c *gin.Context) {
    // 路径参数
    id := c.Param("id")

    // 查询参数: /users?page=1&size=10
    page := c.DefaultQuery("page", "1")
    size := c.Query("size")

    // POST 表单
    name := c.PostForm("name")

    // 请求头
    token := c.GetHeader("Authorization")
}
```

### 3.5 参数绑定与验证

```go
type CreateUserInput struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
    Password string `json:"password" binding:"required,min=8"`
}

func createUser(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // input 已验证通过
}
```

### 3.6 Gin 中间件

```go
// 自定义中间件
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        userID, err := validateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }
        c.Set("userID", userID) // 存进上下文
        c.Next()                // 继续后续处理
    }
}

// 在处理函数中取值
func getProfile(c *gin.Context) {
    userID := c.GetInt("userID")
    // ...
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

---

## 四、RESTful API 完整示例

### 4.1 项目结构

```
myapi/
├── main.go
├── handler/
│   └── user.go
├── model/
│   └── user.go
├── middleware/
│   └── auth.go
├── service/
│   └── user.go
└── go.mod
```

### 4.2 统一响应格式

```go
// handler/response.go
type Response struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
    c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
    c.JSON(status, Response{Code: -1, Message: msg})
}
```

### 4.3 CRUD Handler

```go
// handler/user.go
type UserHandler struct {
    svc *service.UserService
}

func (h *UserHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    users, total, err := h.svc.List(c.Request.Context(), page, size)
    if err != nil {
        Fail(c, 500, "internal error")
        return
    }
    Success(c, gin.H{
        "users": users,
        "total": total,
        "page":  page,
        "size":  size,
    })
}

func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        Fail(c, 400, "invalid id")
        return
    }
    user, err := h.svc.GetByID(c.Request.Context(), id)
    if err != nil {
        Fail(c, 404, "user not found")
        return
    }
    Success(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        Fail(c, 400, err.Error())
        return
    }
    user, err := h.svc.Create(c.Request.Context(), input)
    if err != nil {
        Fail(c, 500, "create failed")
        return
    }
    c.JSON(201, Response{Code: 0, Message: "created", Data: user})
}
```

---

## 五、文件上传

```go
// 单文件上传
r.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    dst := filepath.Join("./uploads", file.Filename)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": "save failed"})
        return
    }
    c.JSON(200, gin.H{"filename": file.Filename, "size": file.Size})
})

// 多文件上传
r.POST("/uploads", func(c *gin.Context) {
    form, _ := c.MultipartForm()
    files := form.File["files"]
    for _, file := range files {
        dst := filepath.Join("./uploads", file.Filename)
        c.SaveUploadedFile(file, dst)
    }
    c.JSON(200, gin.H{"count": len(files)})
})

// 限制文件大小
r.MaxMultipartMemory = 8 << 20 // 8MB
```

---

## 六、优雅关停

```go
func main() {
    r := gin.Default()
    // ... 路由注册

    srv := &http.Server{Addr: ":8080", Handler: r}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("server forced to shutdown:", err)
    }
    log.Println("server exited")
}
```

---

## 速查表

| 概念 | 要点 |
|------|------|
| Handler | 实现 `ServeHTTP(w, r)` |
| ServeMux | Go 1.22 支持 `METHOD /path/{param}` |
| 中间件 | `func(http.Handler) http.Handler` |
| Gin 路由 | `:id` 精确 / `*path` 通配 |
| 参数绑定 | `ShouldBindJSON` + `binding:"required"` |
| JSON 响应 | `c.JSON(status, data)` |
| 路由分组 | `r.Group("/api", middleware)` |
| 优雅关停 | `signal.Notify` + `srv.Shutdown(ctx)` |

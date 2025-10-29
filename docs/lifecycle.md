# 资源管理详解

在生产环境中，正确管理过滤器的生命周期至关重要。`go-sensitive-word` 提供了完善的资源管理机制，确保优雅关闭和避免资源泄漏。

## 为什么需要资源管理？

过滤器在后台运行 goroutine 处理动态更新（添加/删除词），如果不正确关闭可能导致：

- **Goroutine 泄漏**：后台 goroutine 无法退出
- **资源占用**：内存、CPU 资源持续占用
- **数据不一致**：未完成的异步操作可能丢失

## 基本关闭

### Close

使用 `Close()` 立即关闭过滤器，停止所有 goroutine：

```go
defer filter.Close()
```

**特点：**
- 立即停止所有后台 goroutine
- 不等待未完成的异步操作
- 适合快速关闭场景

**示例：**
```go
func main() {
    filter, _ := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    defer filter.Close() // 程序退出时关闭

    // 使用过滤器...
}
```

## 优雅关闭（推荐）

### Shutdown

使用 `Shutdown(ctx)` 优雅关闭，等待异步处理完成：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
defer filter.Shutdown(ctx)
```

**特点：**
- 等待未完成的异步操作完成
- 支持超时控制
- 适合生产环境

**示例：**
```go
func main() {
    filter, _ := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )

    // 优雅关闭（带超时）
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    defer filter.Shutdown(ctx)

    // 使用过滤器...
}
```

## HTTP 服务集成

在 HTTP 服务中，结合 `graceful shutdown` 使用：

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    // 创建过滤器
    filter, _ := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    defer filter.Close()

    // HTTP 服务器
    mux := http.NewServeMux()
    mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
        text := r.URL.Query().Get("text")
        isSensitive := filter.IsSensitive(text)
        fmt.Fprintf(w, "结果: %v", isSensitive)
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // 启动服务器
    go server.ListenAndServe()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // 优雅关闭 HTTP 服务器
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    server.Shutdown(ctx)
    
    // 优雅关闭过滤器
    filter.Shutdown(ctx)
    
    fmt.Println("服务已关闭")
}
```

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    // 创建过滤器
    filter, err := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    if err != nil {
        log.Fatal(err)
    }

    // 加载词库
    filter.LoadDictEmbed(
        sensitive.DictPolitical,
        sensitive.DictViolence,
    )

    // 启动工作 goroutine
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                text := fmt.Sprintf("测试文本 %d-%d", id, j)
                filter.IsSensitive(text)
            }
        }(i)
    }

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    fmt.Println("收到关闭信号，开始优雅关闭...")

    // 等待工作完成
    wg.Wait()

    // 优雅关闭过滤器
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := filter.Shutdown(ctx); err != nil {
        log.Printf("关闭过滤器时出错: %v", err)
    } else {
        fmt.Println("过滤器已优雅关闭")
    }
}
```

## 最佳实践

### 1. 始终使用 defer 关闭

```go
filter, _ := sensitive.NewFilter(...)
defer filter.Close() // 或 defer filter.Shutdown(ctx)
```

### 2. 生产环境使用 Shutdown

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
defer filter.Shutdown(ctx)
```

### 3. 设置合理的超时时间

```go
// 根据实际场景设置超时（建议 5-30 秒）
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
```

### 4. 结合信号处理

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 执行优雅关闭
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
filter.Shutdown(ctx)
```

### 5. 错误处理

```go
if err := filter.Shutdown(ctx); err != nil {
    if err == context.DeadlineExceeded {
        log.Println("关闭超时，强制关闭")
    } else {
        log.Printf("关闭时出错: %v", err)
    }
}
```

## 常见问题

### Q1: Close 和 Shutdown 有什么区别？

**A:**
- `Close()`: 立即关闭，不等待异步操作
- `Shutdown(ctx)`: 优雅关闭，等待异步操作完成或超时

### Q2: 什么时候使用 Close？

**A:** 
- 测试代码
- 快速关闭场景
- 确定没有异步操作

### Q3: 什么时候使用 Shutdown？

**A:**
- 生产环境（推荐）
- 需要确保数据一致性
- 需要等待异步操作完成

### Q4: 超时时间如何设置？

**A:**
- 默认：10 秒
- 大词库/频繁更新：15-30 秒
- 小词库/简单场景：5-10 秒

## 相关文档

- [API 参考文档](./api-reference.md)
- [最佳实践](./best-practices.md)
- [资源管理示例](../../examples/lifecycle/main.go)

# 最佳实践

本文档提供 `go-sensitive-word` 在生产环境中的最佳实践建议。

## 生产环境部署

### 1. 算法选择

**强烈推荐使用 AC 算法：**

```go
filter, err := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC}, // 推荐
)
```

**原因：**
- 性能更好，特别是长文本和大词库
- 支持窗口合并，动态更新更高效
- 并发安全性更好

### 2. 资源管理

**始终使用优雅关闭：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
defer filter.Shutdown(ctx)
```

**HTTP 服务集成：**
```go
// 接收信号
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 优雅关闭
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
filter.Shutdown(ctx)
```

详见：[资源管理详解](./lifecycle.md)

### 3. 异步处理

**添加/删除词后等待处理完成：**

```go
filter.AddWord("新词")
time.Sleep(100 * time.Millisecond) // 等待异步处理完成
result := filter.IsSensitive("包含新词的文本")
```

**推荐延迟时间：**
- DFA: 100ms
- AC: 100-200ms（窗口合并可能需要更长时间）

### 4. 词库加载

**优先使用回调函数从配置中心或数据库加载：**

```go
filter.LoadDictCallback(func() ([]string, error) {
    // 从配置中心或数据库加载
    return configCenter.GetWords(), nil
}, "config-center")
```

**优点：**
- 集中管理词库
- 支持动态更新
- 便于版本控制

**备用方案：**
```go
// 从文件加载（适合配置文件）
filter.LoadDictPath("/etc/sensitive-words.txt")
```

详见：[词库加载详解](./word-loading.md)

### 5. 监控统计

**定期监控词库状态：**

```go
// 定时任务
ticker := time.NewTicker(5 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    stats := filter.GetStats()
    log.Printf("词库状态: 总词数=%d, 最后更新=%s", 
        stats.TotalWords, stats.LastUpdate)
    
    // 告警：词库为空或异常
    if stats.TotalWords == 0 {
        log.Warn("词库为空，请检查!")
    }
}
```

## 性能优化

### 1. 大词库场景

**使用 AC 算法：**
```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
```

**预加载完整词库：**
```go
// 启动时一次性加载所有词库
filter.LoadDictEmbed(
    sensitive.DictPolitical,
    sensitive.DictViolence,
    // ... 所有需要的词库
)
```

### 2. 高频更新场景

**批量操作优于多次单次操作：**

```go
// ✅ 推荐：批量添加
filter.AddWords([]string{"词1", "词2", "词3", "词4", "词5"})

// ❌ 不推荐：多次单次添加
filter.AddWord("词1")
filter.AddWord("词2")
filter.AddWord("词3")
filter.AddWord("词4")
filter.AddWord("词5")
```

**合并多次更新：**
```go
// 收集一批更新请求
var pendingAdds []string
var pendingDels []string

// 定期批量处理（如每 1 秒）
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for range ticker.C {
    if len(pendingAdds) > 0 {
        filter.AddWords(pendingAdds)
        pendingAdds = nil
    }
    if len(pendingDels) > 0 {
        filter.DelWords(pendingDels)
        pendingDels = nil
    }
}
```

### 3. 长文本检测

**AC 算法单次扫描，性能更优：**

```go
// AC 算法适合长文本
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)

// 长文本检测
longText := strings.Repeat("测试文本", 1000)
result := filter.IsSensitive(longText)
```

### 4. 并发优化

**过滤器是并发安全的，可以安全地在多 goroutine 中使用：**

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        text := fmt.Sprintf("文本 %d", id)
        filter.IsSensitive(text)
    }(i)
}
wg.Wait()
```

## 代码组织

### 1. 单例模式

**全局单例过滤器：**

```go
package sensitive

import (
    "sync"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

var (
    instance *sensitive.Manager
    once     sync.Once
)

func GetFilter() *sensitive.Manager {
    once.Do(func() {
        instance, _ = sensitive.NewFilter(
            sensitive.StoreOption{Type: sensitive.StoreMemory},
            sensitive.FilterOption{Type: sensitive.FilterAC},
        )
        // 加载词库...
    })
    return instance
}
```

### 2. 中间件模式

**HTTP 中间件：**

```go
func SensitiveCheckMiddleware(filter *sensitive.Manager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 读取请求体
            body, _ := io.ReadAll(r.Body)
            r.Body.Close()
            
            // 检测敏感词
            if filter.IsSensitive(string(body)) {
                http.Error(w, "内容包含敏感词", http.StatusBadRequest)
                return
            }
            
            // 恢复请求体
            r.Body = io.NopCloser(bytes.NewReader(body))
            next.ServeHTTP(w, r)
        })
    }
}

// 使用
mux := http.NewServeMux()
// ... 注册路由
filter, _ := sensitive.NewFilter(...)
http.ListenAndServe(":8080", SensitiveCheckMiddleware(filter)(mux))
```

### 3. 错误处理

**完善的错误处理：**

```go
filter, err := sensitive.NewFilter(...)
if err != nil {
    log.Fatalf("初始化过滤器失败: %v", err)
}

if err := filter.LoadDictEmbed(...); err != nil {
    log.Fatalf("加载词库失败: %v", err)
}

if err := filter.AddWord("新词"); err != nil {
    log.Printf("添加敏感词失败: %v", err)
}
```

## 监控与告警

### 1. 词库监控

```go
// 定期检查词库状态
func monitorWords(filter *sensitive.Manager) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := filter.GetStats()
        
        // 告警条件
        if stats.TotalWords == 0 {
            alert("词库为空!")
        }
        
        if time.Since(stats.LastUpdate) > 24*time.Hour {
            alert("词库超过 24 小时未更新")
        }
    }
}
```

### 2. 性能监控

```go
// 记录检测耗时
func checkWithMetrics(filter *sensitive.Manager, text string) bool {
    start := time.Now()
    result := filter.IsSensitive(text)
    duration := time.Since(start)
    
    // 记录指标
    metrics.Record("sensitive_check_duration", duration)
    
    return result
}
```

### 3. 日志记录

```go
// 记录敏感词检测日志
func checkWithLogging(filter *sensitive.Manager, text string) bool {
    result := filter.IsSensitive(text)
    if result {
        words := filter.FindAll(text)
        log.Printf("检测到敏感词: %v, 文本: %s", words, text)
    }
    return result
}
```

## 安全建议

### 1. 词库备份

**定期备份词库：**

```go
// 每天备份一次
func backupWords(filter *sensitive.Manager) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        content, err := filter.ExportToString()
        if err != nil {
            log.Printf("导出词库失败: %v", err)
            continue
        }
        
        filename := fmt.Sprintf("backup-%s.txt", time.Now().Format("20060102"))
        os.WriteFile(filename, []byte(content), 0644)
    }
}
```

### 2. 词库验证

**启动时验证词库：**

```go
func validateWords(filter *sensitive.Manager) error {
    stats := filter.GetStats()
    if stats.TotalWords == 0 {
        return errors.New("词库为空")
    }
    
    // 测试检测
    if !filter.IsSensitive("测试敏感词") {
        log.Warn("词库可能未正确加载")
    }
    
    return nil
}
```

### 3. 敏感信息处理

**结合工具函数处理敏感信息：**

```go
func sanitizeContent(text string) string {
    // 1. 敏感词替换
    text = filter.Replace(text, '*')
    
    // 2. 邮箱屏蔽
    text = sensitive.MaskEmail(text)
    
    // 3. URL 屏蔽
    text = sensitive.MaskURL(text)
    
    return text
}
```

## 测试建议

### 1. 单元测试

```go
func TestSensitiveFilter(t *testing.T) {
    filter, _ := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    filter.AddWord("测试词")
    time.Sleep(100 * time.Millisecond)
    
    if !filter.IsSensitive("包含测试词的文本") {
        t.Error("应该检测到敏感词")
    }
}
```

### 2. 性能测试

```go
func BenchmarkSensitiveCheck(b *testing.B) {
    filter, _ := sensitive.NewFilter(...)
    filter.LoadDictEmbed(...)
    text := "测试文本"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        filter.IsSensitive(text)
    }
}
```

### 3. 并发测试

```go
func TestConcurrentAccess(t *testing.T) {
    filter, _ := sensitive.NewFilter(...)
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            filter.IsSensitive("测试文本")
        }()
    }
    wg.Wait()
}
```

## 相关文档

- [算法选择指南](./algorithm-guide.md)
- [词库管理详解](./word-management.md)
- [资源管理详解](./lifecycle.md)
- [常见问题](./faq.md)

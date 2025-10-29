# 词库管理详解

`go-sensitive-word` 支持在生产环境中动态管理敏感词库，无需重启服务即可更新词库。

## 动态添加敏感词

### 单个/多个添加

使用 `AddWord` 添加一个或多个敏感词：

```go
// 添加单个词
filter.AddWord("敏感词1")

// 添加多个词
filter.AddWord("敏感词1", "敏感词2", "敏感词3")
```

### 批量添加

使用 `AddWords` 批量添加：

```go
words := []string{"词1", "词2", "词3", "词4"}
err := filter.AddWords(words)
if err != nil {
    log.Fatal(err)
}
```

**性能建议**：批量添加时优先使用 `AddWords`，而非多次调用 `AddWord`。

### 异步处理说明

**重要提示**：DFA/AC 算法通过 channel 异步处理词的添加，需要短暂延迟才能生效：

```go
filter.AddWord("新敏感词")
time.Sleep(100 * time.Millisecond) // 等待异步处理完成
result := filter.IsSensitive("包含新敏感词的文本")
```

**建议延迟时间：**
- DFA：100ms
- AC：100-200ms（窗口合并可能需要更长时间）

## 动态删除敏感词

### 单个/多个删除

使用 `DelWord` 删除一个或多个敏感词：

```go
// 删除单个词
filter.DelWord("旧敏感词1")

// 删除多个词
filter.DelWord("旧敏感词1", "旧敏感词2", "旧敏感词3")
```

### 批量删除

使用 `DelWords` 批量删除：

```go
words := []string{"词1", "词2"}
err := filter.DelWords(words)
```

### 异步处理说明

删除操作同样需要异步处理延迟：

```go
filter.DelWord("旧敏感词")
time.Sleep(100 * time.Millisecond)
result := filter.IsSensitive("包含旧敏感词的文本") // 应该返回 false
```

## 批量替换敏感词

使用 `ReplaceWords` 批量替换（先删旧词，再加新词）：

```go
err := filter.ReplaceWords(
    []string{"旧词1", "旧词2"},  // 要删除的旧词
    []string{"新词1", "新词2"},  // 要添加的新词
)
```

**使用场景：**
- 敏感词更新（旧词不再敏感，替换为新敏感词）
- 词库维护（定期替换过时词汇）

## 清空词库

使用 `Clear` 清空所有敏感词：

```go
err := filter.Clear()
if err != nil {
    log.Fatal(err)
}
```

**注意**：清空后需要重新加载词库才能恢复功能。

## 获取词库统计信息

使用 `GetStats` 获取词库统计信息：

```go
stats := filter.GetStats()
fmt.Printf("总词数: %d\n", stats.TotalWords)
fmt.Printf("最后更新: %s\n", stats.LastUpdate)
```

**返回信息：**
- `TotalWords int64`：当前词库总词数
- `LastUpdate time.Time`：最后更新时间

**使用场景：**
- 监控词库状态
- 性能指标采集
- 运维告警

## 合并词库

使用 `MergeFromManager` 合并另一个 Manager 的词库：

```go
// 创建另一个过滤器
filter2, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
filter2.LoadDictEmbed(/* ... */)

// 合并到 filter1
err := filter1.MergeFromManager(filter2)
```

**使用场景：**
- 多词库合并
- 词库整合

## 从文件刷新词库

使用 `RefreshFromPath` 从文件刷新词库：

```go
// 替换模式（清空后重新加载）
err := filter.RefreshFromPath("/path/to/words.txt", true)

// 追加模式（保留现有词库）
err := filter.RefreshFromPath("/path/to/words.txt", false)
```

**参数说明：**
- `filePath`：文件路径
- `replace`：`true` 表示替换模式，`false` 表示追加模式

**使用场景：**
- 定期从文件更新词库
- 词库版本管理

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    filter, err := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    if err != nil {
        log.Fatal(err)
    }
    defer filter.Close()

    // 1. 添加敏感词
    err = filter.AddWords([]string{"敏感词1", "敏感词2"})
    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(200 * time.Millisecond)

    // 2. 获取统计信息
    stats := filter.GetStats()
    fmt.Printf("当前词数: %d\n", stats.TotalWords)

    // 3. 删除敏感词
    filter.DelWord("敏感词1")
    time.Sleep(200 * time.Millisecond)

    // 4. 批量替换
    filter.ReplaceWords(
        []string{"敏感词2"},
        []string{"新敏感词"},
    )
    time.Sleep(200 * time.Millisecond)

    // 5. 导出词库
    content, _ := filter.ExportToString()
    fmt.Printf("导出词库:\n%s", content)
}
```

## 注意事项

### 1. 异步处理延迟

所有动态操作（添加、删除、替换）都是异步处理的，需要等待处理完成：

```go
filter.AddWord("新词")
time.Sleep(100 * time.Millisecond) // 必须等待
result := filter.IsSensitive("包含新词的文本")
```

### 2. 并发安全

所有词库管理操作都是并发安全的，可以在多 goroutine 中调用：

```go
// 并发添加
go filter.AddWord("词1")
go filter.AddWord("词2")
go filter.AddWord("词3")
```

### 3. 错误处理

建议对关键操作进行错误处理：

```go
if err := filter.AddWord("敏感词"); err != nil {
    log.Printf("添加敏感词失败: %v", err)
}
```

### 4. 生产环境建议

1. **批量操作**：优先使用批量方法（`AddWords`、`DelWords`）
2. **监控统计**：定期调用 `GetStats()` 监控词库状态
3. **优雅关闭**：使用 `Shutdown(ctx)` 确保异步操作完成

## 相关文档

- [API 参考文档](./api-reference.md)
- [词库加载详解](./word-loading.md)
- [最佳实践](./best-practices.md)
- [动态维护示例](../../examples/dynamic/main.go)

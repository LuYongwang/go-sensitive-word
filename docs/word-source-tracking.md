# 词库来源追踪功能详解

`go-sensitive-word` 提供了强大的词库来源追踪功能，可以精确识别每个敏感词所属的词库来源。

## 概述

来源追踪功能允许您：
- ✅ **精确溯源**：追踪每个敏感词来自哪个词库
- ✅ **多来源管理**：一个词可以来自多个词库
- ✅ **统计分析**：基于来源进行数据分析和审计
- ✅ **业务隔离**：不同业务场景使用不同的词库来源标识

## 核心概念

### 来源标识

来源标识（Source）是一个字符串，用于标识词库的来源，例如：
- `"political"` - 政治类型词库
- `"violence"` - 暴恐词库
- `"pornography"` - 色情词库
- `"custom"` - 自定义词库
- `"business"` - 业务词库

### 多来源支持

一个词可以同时属于多个来源，例如：
- `"温云松"` 可能同时出现在 `political` 和 `tencent` 词库中
- 查询时会返回所有来源：`[political tencent]`

## 核心API

### 1. 加载词库并指定来源

#### LoadDictEmbedWithSource

加载内置词库并指定来源标识：

```go
filter.LoadDictEmbedWithSource(sensitive.DictPolitical, "political")
filter.LoadDictEmbedWithSource(sensitive.DictViolence, "violence")
```

**参数：**
- `content`: 词库内容字符串
- `source`: 来源标识

**使用场景：**
- 为不同类别的内置词库指定明确的来源标识
- 便于后续统计和分析

### 2. 添加词并指定来源

#### AddWordsWithSource

批量添加敏感词并指定来源：

```go
customWords := []string{"违禁词A", "违禁词B", "敏感词C"}
err := filter.AddWordsWithSource(customWords, "custom")
if err != nil {
    log.Fatalf("添加词失败: %v", err)
}
```

**参数：**
- `words`: 敏感词列表
- `source`: 来源标识

**特性：**
- 自动归一化：词会被自动归一化处理
- 来源合并：如果词已存在，会追加来源而不会重复

### 3. 查询词的来源

#### GetWordSources

查询单个词的来源列表：

```go
sources := filter.GetWordSources("温云松")
fmt.Printf("来源: %v\n", sources)
// 输出: 来源: [political]
```

**返回值：**
- `[]string`: 词的来源列表，如果词不存在返回 `nil`

#### GetAllWordSources

获取所有词的来源映射：

```go
allSources := filter.GetAllWordSources()
for word, sources := range allSources {
    fmt.Printf("%s: %v\n", word, sources)
}
```

**返回值：**
- `map[string][]string`: 词到来源列表的映射

### 4. 查找敏感词及其来源

#### FindAllWithSource

查找文本中所有敏感词及其来源信息：

```go
text := "这段文本包含违禁词A和违禁词B，以及温云松"
results := filter.FindAllWithSource(text)

for _, result := range results {
    fmt.Printf("词: %s, 来源: %v\n", result.Word, result.Source)
}
// 输出:
// 词: 违禁词A, 来源: [custom business]
// 词: 违禁词B, 来源: [custom]
// 词: 温云松, 来源: [political]
```

**返回值：**
- `[]MatchResult`: 包含词和来源信息的列表

**MatchResult 结构：**
```go
type MatchResult struct {
    Word   string   // 匹配到的敏感词
    Source []string // 该词所属的词库来源列表
}
```

#### FindAllCountWithSource

查找所有敏感词及其出现次数和来源信息：

```go
text := "违禁词A出现了违禁词A两次，还有违禁词B"
countMap := filter.FindAllCount(text)
countWithSource := filter.FindAllCountWithSource(text)

for word, result := range countWithSource {
    count := countMap[word]
    fmt.Printf("词: %s 出现 %d 次, 来源: %v\n", word, count, result.Source)
}
```

**返回值：**
- `map[string]MatchResult`: 词到 MatchResult 的映射（包含来源信息）

## 完整示例

### 示例1：基本使用

```go
package main

import (
    "fmt"
    "log"
    "time"

    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    // 初始化过滤器
    filter, err := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    if err != nil {
        log.Fatalf("初始化失败: %v", err)
    }
    defer filter.Close()

    // 从不同来源加载词库
    filter.LoadDictEmbedWithSource(sensitive.DictPolitical, "political")
    filter.LoadDictEmbedWithSource(sensitive.DictViolence, "violence")

    // 添加自定义词
    customWords := []string{"违禁词A", "违禁词B"}
    filter.AddWordsWithSource(customWords, "custom")

    // 等待异步处理
    time.Sleep(200 * time.Millisecond)

    // 查找敏感词及其来源
    text := "这段文本包含违禁词A和温云松"
    results := filter.FindAllWithSource(text)

    for _, result := range results {
        fmt.Printf("敏感词: %s, 来源: %v\n", result.Word, result.Source)
    }
}
```

### 示例2：多业务场景隔离

```go
// 场景 A：社交平台业务
filterA, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
filterA.LoadDictEmbedWithSource(sensitive.DictPolitical, "social-political")
filterA.LoadDictEmbedWithSource(sensitive.DictViolence, "social-violence")

// 场景 B：电商平台业务
filterB, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
filterB.LoadDictEmbedWithSource(sensitive.DictAdvertisement, "ecommerce-ad")
filterB.LoadDictEmbedWithSource(sensitive.DictCorruption, "ecommerce-corruption")

// 各业务独立检测
resultA := filterA.FindAllWithSource(text) // 使用社交平台词库
resultB := filterB.FindAllWithSource(text) // 使用电商平台词库
```

### 示例3：统计和审计

```go
// 获取所有词的来源分布
allSources := filter.GetAllWordSources()

// 统计各来源的词数
sourceStats := make(map[string]int)
for _, sources := range allSources {
    for _, source := range sources {
        sourceStats[source]++
    }
}

fmt.Println("词库来源分布：")
for source, count := range sourceStats {
    fmt.Printf("  %s: %d 个词\n", source, count)
}
```

## 性能考虑

### 优化策略

1. **合并锁优化**：来源追踪使用与词库相同的锁，避免额外的锁开销
2. **内存优化**：来源列表使用切片，内存占用低
3. **查询优化**：GetWordSources 提供 O(1) 时间复杂度的查询

### 内存占用

对于每个词：
- **基础存储**：词字符串本身
- **来源追踪**：平均每个词 1-3 个来源，每个来源约 10-20 字节
- **总开销**：约增加 20-40% 内存占用

**建议：**
- 10万词汇量：额外占用约 2-4MB 内存
- 对于性能敏感场景，可考虑不使用来源追踪功能

## 最佳实践

### 1. 来源命名规范

使用清晰的来源标识：

```go
// ✅ 推荐：语义清晰的命名
"political"           // 政治类型
"violence"            // 暴恐类型
"custom"              // 自定义
"business-ads"        // 业务广告

// ❌ 不推荐：模糊的命名
"dict1"               // 不清楚来源
"new"                 // 含义不明
"123"                 // 无意义
```

### 2. 按需加载

根据业务需求选择性加载词库：

```go
// 只加载需要的词库
filter.LoadDictEmbedWithSource(sensitive.DictPolitical, "political")
filter.LoadDictEmbedWithSource(sensitive.DictViolence, "violence")
// 不需要的词库不加载，节省资源
```

### 3. 及时清理

对于临时词库，使用完毕后及时清理：

```go
// 临时词库使用
tempFilter.LoadDictEmbedWithSource(content, "temp")

// 使用完毕后清理
defer tempFilter.Close()
```

## 常见问题

### Q: 为什么需要来源追踪？

A: 来源追踪可以帮助您：
- 追踪敏感词的来源，便于审计
- 统计分析不同词库的效果
- 区分业务场景和违规类型
- 支持精细化的审核策略

### Q: 来源追踪会影响性能吗？

A: 略微影响：
- **内存占用**：增加约 20-40%
- **查询性能**：几乎无影响（O(1) 查询）
- **写入性能**：略有下降（需要维护来源映射）

### Q: 如何禁用来源追踪？

A: 使用不带来源的 API：
- 使用 `LoadDictEmbed()` 而不是 `LoadDictEmbedWithSource()`
- 使用 `AddWords()` 而不是 `AddWordsWithSource()`

### Q: 一个词可以来自多个来源吗？

A: 可以。例如：
```go
filter.AddWordsWithSource([]string{"测试词"}, "source1")
filter.AddWordsWithSource([]string{"测试词"}, "source2")
// "测试词" 的来源为 [source1 source2]
```

## 相关文档

- [API 参考文档](./api-reference.md)
- [词库加载详解](./word-loading.md)
- [词库管理详解](./word-management.md)
- [示例代码](../../examples/word-source/main.go)


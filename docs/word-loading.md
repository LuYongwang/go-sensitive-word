# 词库加载详解

`go-sensitive-word` 支持多种方式加载敏感词库，适应不同的使用场景。

## 加载方式概览

| 加载方式 | 方法 | 适用场景 |
|---------|------|---------|
| **内置词库** | `LoadDictEmbed` | 编译时嵌入，无需外部文件 |
| **文件加载** | `LoadDictPath` | 本地文件、配置文件 |
| **回调加载** | `LoadDictCallback` | 数据库、Redis、配置中心 |

## 方式一：加载内置词库（推荐）

### 基本用法

使用 `LoadDictEmbed` 加载编译时嵌入的词库：

```go
err := filter.LoadDictEmbed(
    sensitive.DictPolitical,
    sensitive.DictViolence,
    sensitive.DictPornography,
)
```

### 内置词库列表

项目内置了多个分类词库，可按需选择加载：

| 词库变量 | 说明 | 文件路径 |
|---------|------|---------|
| `DictGFWAdditional` | GFW补充词库 | `wordlists/GFW补充词库.txt` |
| `DictOther` | 其他词库 | `wordlists/其他词库.txt` |
| `DictReactionary` | 反动词库 | `wordlists/反动词库.txt` |
| `DictAdvertisement` | 广告类型 | `wordlists/广告类型.txt` |
| `DictPolitical` | 政治类型 | `wordlists/政治类型.txt` |
| `DictViolence` | 暴恐词库 | `wordlists/暴恐词库.txt` |
| `DictPeopleLife` | 民生词库 | `wordlists/民生词库.txt` |
| `DictGunExplosion` | 涉枪涉爆 | `wordlists/涉枪涉爆.txt` |
| `DictNeteaseFE` | 网易前端过滤敏感词库 | `wordlists/网易前端过滤敏感词库.txt` |
| `DictSexual` | 色情类型 | `wordlists/色情类型.txt` |
| `DictPornography` | 色情词库 | `wordlists/色情词库.txt` |
| `DictAdditional` | 补充词库 | `wordlists/补充词库.txt` |
| `DictCorruption` | 贪腐词库 | `wordlists/贪腐词库.txt` |
| `DictTemporaryTencent` | 零时-Tencent | `wordlists/零时-Tencent.txt` |
| `DictIllegalURL` | 非法网址 | `wordlists/非法网址.txt` |

### 完整示例

```go
package main

import (
    "log"
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

    // 加载多个内置词库
    err = filter.LoadDictEmbed(
        sensitive.DictPolitical,
        sensitive.DictViolence,
        sensitive.DictSexual,
        sensitive.DictPornography,
        sensitive.DictIllegalURL,
    )
    if err != nil {
        log.Fatalf("加载词库失败: %v", err)
    }

    // 使用过滤器...
}
```

### 优点

- ✅ **无需外部文件**：词库编译时嵌入，分发简单
- ✅ **性能优秀**：无需文件 I/O，启动快
- ✅ **版本一致**：词库与代码版本绑定

### 缺点

- ⚠️ **更新需重新编译**：词库更新需要重新编译程序
- ⚠️ **增加二进制体积**：词库嵌入会增加可执行文件大小

## 方式二：从文件加载

### 基本用法

使用 `LoadDictPath` 从文件路径加载词库：

```go
// 加载单个文件
err := filter.LoadDictPath("/path/to/words.txt")

// 加载多个文件
err := filter.LoadDictPath(
    "/path/to/words1.txt",
    "/path/to/words2.txt",
)
```

### 文件格式

词库文件格式：**每行一个敏感词**

```
敏感词1
敏感词2
敏感词3
```

**示例文件内容：**
```
政治敏感词1
政治敏感词2
暴力敏感词1
色情敏感词1
```

### 完整示例

```go
package main

import (
    "log"
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

    // 从文件加载词库
    err = filter.LoadDictPath(
        "wordlists/custom1.txt",
        "wordlists/custom2.txt",
    )
    if err != nil {
        log.Fatalf("加载词库失败: %v", err)
    }

    // 使用过滤器...
}
```

### 刷新词库（支持替换/追加）

使用 `RefreshFromPath` 从文件刷新词库：

```go
// 替换模式：清空现有词库后重新加载
err := filter.RefreshFromPath("/path/to/words.txt", true)

// 追加模式：在现有词库基础上追加
err := filter.RefreshFromPath("/path/to/words.txt", false)
```

**使用场景：**
- 定期从文件更新词库
- 词库版本管理
- 配置热更新

### 优点

- ✅ **更新灵活**：修改文件即可更新词库
- ✅ **易于管理**：词库与代码分离
- ✅ **支持多个文件**：可以分分类别管理

### 缺点

- ⚠️ **需要文件系统**：依赖外部文件
- ⚠️ **启动稍慢**：需要读取文件

查看 [examples/file-load/main.go](../../examples/file-load/main.go) 获取完整示例。

## 方式三：回调函数加载（自定义数据源）

### 基本用法

使用 `LoadDictCallback` 通过回调函数加载词库：

```go
err := filter.LoadDictCallback(func() ([]string, error) {
    // 自定义加载逻辑
    return words, nil
}, "data-source-name")
```

### 从数据库加载

```go
import "database/sql"

filter.LoadDictCallback(func() ([]string, error) {
    rows, err := db.Query("SELECT word FROM sensitive_words WHERE enabled = 1")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var words []string
    for rows.Next() {
        var word string
        if err := rows.Scan(&word); err != nil {
            return nil, err
        }
        words = append(words, word)
    }
    return words, nil
}, "database")
```

### 从 Redis 加载

```go
import "github.com/go-redis/redis/v8"

filter.LoadDictCallback(func() ([]string, error) {
    ctx := context.Background()
    words, err := redisClient.SMembers(ctx, "sensitive:words").Result()
    if err != nil {
        return nil, err
    }
    return words, nil
}, "redis")
```

### 从配置中心加载

```go
filter.LoadDictCallback(func() ([]string, error) {
    // 从配置中心（如 Consul、etcd、Apollo）获取
    config, err := configCenter.Get("sensitive.words")
    if err != nil {
        return nil, err
    }
    
    // 解析配置（假设是 JSON 数组）
    var words []string
    json.Unmarshal([]byte(config), &words)
    return words, nil
}, "config-center")
```

### 完整示例

```go
package main

import (
    "context"
    "log"
    sensitive "github.com/LuYongwang/go-sensitive-word"
    "github.com/go-redis/redis/v8"
)

func main() {
    filter, err := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC},
    )
    if err != nil {
        log.Fatal(err)
    }

    // 从 Redis 加载
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    err = filter.LoadDictCallback(func() ([]string, error) {
        ctx := context.Background()
        words, err := redisClient.SMembers(ctx, "sensitive:words").Result()
        return words, err
    }, "redis")

    if err != nil {
        log.Fatalf("加载词库失败: %v", err)
    }

    // 使用过滤器...
}
```

### 优点

- ✅ **灵活性最高**：支持任意数据源
- ✅ **动态更新**：可以从外部系统实时获取
- ✅ **集中管理**：词库集中存储在数据库/配置中心

### 缺点

- ⚠️ **依赖外部系统**：需要数据库、Redis 等基础设施
- ⚠️ **启动依赖**：启动时需要连接外部系统

查看 [examples/callback/main.go](../../examples/callback/main.go) 获取完整示例。

## 混合使用

可以组合使用多种加载方式：

```go
// 1. 先加载内置词库
filter.LoadDictEmbed(
    sensitive.DictPolitical,
    sensitive.DictViolence,
)

// 2. 再从文件追加
filter.LoadDictPath("wordlists/custom.txt")

// 3. 最后从数据库追加
filter.LoadDictCallback(func() ([]string, error) {
    return db.QueryWords(), nil
}, "database")
```

**注意**：多次调用 `LoadDictEmbed`、`LoadDictPath`、`LoadDictCallback` 都是追加模式，不会清空已有词库。

## 词库格式要求

无论使用哪种加载方式，词库数据格式都需满足：

- **每行一个词**：单个或多个敏感词，每个词一行
- **自动去重**：重复的词会自动去重
- **自动归一化**：词会被自动归一化（大小写、全角半角等）

## 性能考虑

| 加载方式 | 启动速度 | 更新灵活性 | 推荐场景 |
|---------|---------|-----------|---------|
| 内置词库 | ⭐⭐⭐⭐⭐ | ⭐ | 固定词库、简单部署 |
| 文件加载 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 配置文件、版本管理 |
| 回调加载 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 动态更新、集中管理 |

## 最佳实践

1. **生产环境**：优先使用 `LoadDictCallback` 从配置中心或数据库加载
2. **开发环境**：使用 `LoadDictEmbed` 或 `LoadDictPath` 快速启动
3. **定期更新**：使用 `RefreshFromPath` 或定期调用 `LoadDictCallback` 更新词库
4. **监控加载**：记录加载时间和词库大小，便于排查问题

## 相关文档

- [API 参考文档](./api-reference.md)
- [词库管理详解](./word-management.md)
- [最佳实践](./best-practices.md)
- [文件加载示例](../../examples/file-load/main.go)
- [回调加载示例](../../examples/callback/main.go)

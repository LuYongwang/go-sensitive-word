# 快速开始指南

本指南将帮助您快速上手 `go-sensitive-word`，5 分钟内完成第一个敏感词检测程序。

## 安装

```bash
go get -u github.com/LuYongwang/go-sensitive-word@latest
```

## 最简单的例子

```go
package main

import (
    "fmt"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    // 创建过滤器
    filter, err := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterDfa},
    )
    if err != nil {
        panic(err)
    }

    // 添加敏感词
    filter.AddWord("敏感词1", "敏感词2")

    // 检测文本
    text := "这是一段包含敏感词1的文本"
    if filter.IsSensitive(text) {
        fmt.Println("检测到敏感词!")
        words := filter.FindAll(text)
        fmt.Printf("敏感词: %v\n", words)
    }
}
```

## 完整示例

查看 [examples/basic/main.go](../../examples/basic/main.go) 获取完整的基础示例代码。

## 核心概念

### 1. 过滤器初始化

`NewFilter` 需要两个参数：
- `StoreOption`：词库存储方式（目前仅支持 `StoreMemory`）
- `FilterOption`：过滤算法（支持 `FilterDfa` 或 `FilterAC`）

```go
filter, err := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC}, // 推荐生产环境使用 AC
)
```

### 2. 加载词库

#### 方式一：加载内置词库（推荐）

```go
err := filter.LoadDictEmbed(
    sensitive.DictGFWAdditional,
    sensitive.DictPolitical,
    sensitive.DictViolence,
    // ... 更多词库
)
```

内置词库列表：
- `DictGFWAdditional` - GFW补充词库
- `DictOther` - 其他词库
- `DictReactionary` - 反动词库
- `DictAdvertisement` - 广告类型
- `DictPolitical` - 政治类型
- `DictViolence` - 暴恐词库
- `DictPeopleLife` - 民生词库
- `DictGunExplosion` - 涉枪涉爆
- `DictNeteaseFE` - 网易前端过滤敏感词库
- `DictSexual` - 色情类型
- `DictPornography` - 色情词库
- `DictAdditional` - 补充词库
- `DictCorruption` - 贪腐词库
- `DictTemporaryTencent` - 零时-Tencent
- `DictIllegalURL` - 非法网址

#### 方式二：从文件加载

```go
err := filter.LoadDictPath("/path/to/words.txt")
```

详细说明见：[词库加载详解](./word-loading.md)

#### 方式三：从自定义数据源加载

```go
filter.LoadDictCallback(func() ([]string, error) {
    // 从数据库、Redis等加载
    return db.QueryWords(), nil
}, "custom")
```

详细说明见：[词库加载详解](./word-loading.md)

### 3. 文本检测

#### 判断是否包含敏感词

```go
if filter.IsSensitive("包含敏感词的文本") {
    // 处理逻辑
}
```

#### 查找敏感词

```go
// 查找第一个敏感词
word := filter.FindOne("包含敏感词的文本")

// 查找所有敏感词（去重）
words := filter.FindAll("包含敏感词的文本")

// 查找所有敏感词及出现次数
wordCount := filter.FindAllCount("包含敏感词的文本")
```

#### 替换/删除敏感词

```go
// 替换为 * 号
result := filter.Replace("包含敏感词的文本", '*')

// 直接删除敏感词
result := filter.Remove("包含敏感词的文本")
```

### 4. 动态管理词库

```go
// 添加敏感词
filter.AddWord("新敏感词1", "新敏感词2")

// 删除敏感词
filter.DelWord("旧敏感词1", "旧敏感词2")

// 批量替换
filter.ReplaceWords(
    []string{"旧词1", "旧词2"},
    []string{"新词1", "新词2"},
)
```

**重要提示**：DFA/AC 算法通过 channel 异步处理词的添加/删除，需要短暂延迟才能生效：

```go
filter.AddWord("新词")
time.Sleep(100 * time.Millisecond) // 等待异步处理完成
result := filter.IsSensitive("包含新词的文本")
```

详细说明见：[词库管理详解](./word-management.md)

## 下一步

- 📖 查看 [API 参考文档](./api-reference.md) 了解所有 API
- 🔧 学习 [算法选择指南](./algorithm-guide.md) 选择合适的算法
- 🛡️ 了解 [归一化功能详解](./normalization.md) 防御混淆攻击
- 📚 阅读 [最佳实践](./best-practices.md) 准备生产环境部署

## 相关示例

- [基础功能演示](../../examples/basic/main.go)
- [AC 算法示例](../../examples/ac/main.go)
- [综合功能演示](../../examples/comprehensive/main.go)

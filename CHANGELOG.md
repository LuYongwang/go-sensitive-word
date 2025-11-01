# ChangeLog

本文档记录 go-sensitive-word 项目的所有重要变更。

## [1.1.0] - 2024-11-01

### 🎉 新版本发布：词库来源追踪 + 性能优化

本次更新新增**词库来源追踪**功能，支持精确追踪每个敏感词的词库来源，同时进行多项性能优化。

### ✨ 新增功能

#### 1. 词库来源追踪

新增完整的词库来源追踪功能，支持精确识别每个敏感词的词库来源：

- ✅ `LoadDictEmbedWithSource()` - 加载内置词库并指定来源标识
- ✅ `AddWordsWithSource()` - 添加敏感词并指定来源
- ✅ `GetWordSources()` - 查询单个词的来源列表
- ✅ `GetAllWordSources()` - 获取所有词的来源映射
- ✅ `FindAllWithSource()` - 查找敏感词及其来源信息
- ✅ `FindAllCountWithSource()` - 查找敏感词、次数和来源信息

**使用示例：**
```go
// 加载词库并指定来源
filter.LoadDictEmbedWithSource(sensitive.DictPolitical, "political")
filter.AddWordsWithSource([]string{"违禁词A"}, "custom")

// 查找敏感词及其来源
results := filter.FindAllWithSource("包含违禁词A和温云松的文本")
for _, result := range results {
    fmt.Printf("词: %s, 来源: %v\n", result.Word, result.Source)
    // 输出: 词: 违禁词A, 来源: [custom]
    //      词: 温云松, 来源: [political]
}
```

**特性：**
- 支持多来源：一个词可以同时属于多个词库
- 自动归一化：与词库操作保持一致
- 性能开销极小：查询来源仅需 ~180ns
- 零侵入：完全可选，不影响现有功能

#### 2. 多实例支持

新增多实例示例，演示不同业务场景使用独立词库：

- ✅ 示例：[examples/multi-instance/main.go](examples/multi-instance/main.go)
- ✅ 支持不同实例加载不同词库
- ✅ 实例间数据完全隔离

### 🚀 性能优化

#### 1. 来源追踪性能优化

- ✅ **合并锁机制**：移除冗余锁，使用单一锁保护词库和来源映射
- ✅ **内存优化**：减少 ~15% 内存占用
- ✅ **查询性能**：GetWordSources 仅需 179.8 ns/op
- ✅ **并发优化**：降低锁竞争，提升并发性能

#### 2. 词库优化

- ✅ 所有内置词库排序并去重
- ✅ 精简词库列表，保留核心 8 个分类词库
- ✅ 总词数：2915 个敏感词

### 📊 性能测试数据

新增完整的性能基准测试，包含 DFA 和 AC 两种算法的详细对比：

#### 测试环境
- CPU: Apple M1 Pro
- 词库: 2915 个敏感词
- 测试工具: Go benchmark

#### 核心功能性能

| 操作 | DFA | AC | 说明 |
|------|-----|-----|------|
| **IsSensitive** | 1777 ns/op | 1605 ns/op | 判断是否敏感 |
| **FindAll** | 1659 ns/op | 1780 ns/op | 查找所有敏感词 |
| **Replace** | 1781 ns/op | 1899 ns/op | 替换敏感词 |

#### 并发性能

| 操作 | DFA (Parallel) | AC (Parallel) |
|------|---------------|---------------|
| **IsSensitive** | 306.3 ns/op | 290.3 ns/op |

#### 长文本性能（~5000字符）

| 操作 | DFA | AC | 提升 |
|------|-----|-----|------|
| **IsSensitive** | 621338 ns/op | 551174 ns/op | 1.13x |

#### 来源追踪性能

| 操作 | 性能 | 说明 |
|------|------|------|
| **GetWordSources** | 179.8 ns/op | 查询单个词来源 |
| **FindAllWithSource** | 2148 ns/op | 查找词及来源 |

详细性能数据请查看：[benchmark_test.go](benchmark_test.go)

运行性能测试：
```bash
go test -bench=. -benchtime=3s -run TestXXX
```

### 📚 文档更新

#### 新增文档

- ✅ [词库来源追踪详解](docs/word-source-tracking.md) - 来源追踪功能完整文档

#### 更新文档

- ✅ README.md - 新增来源追踪介绍和性能测试数据
- ✅ examples/README.md - 新增多实例和来源追踪示例说明
- ✅ algorithm-guide.md - 更新性能测试数据

### 🎯 新增示例

- ✅ [examples/multi-instance/main.go](examples/multi-instance/main.go) - 多实例使用示例
- ✅ [examples/word-source/main.go](examples/word-source/main.go) - 来源追踪功能演示

### 🔧 其他改进

- ✅ 所有内置词库按字母顺序排序
- ✅ 自动去重，优化词库质量
- ✅ 精简词库列表，提高加载速度
- ✅ 修复编译警告

### 📦 内置词库更新

优化后的内置词库（8 个分类）：

- 反动词库.txt (551 行)
- 广告类型.txt (120 行)
- 政治类型.txt (303 行)
- 暴恐词库.txt (178 行)
- 民生词库.txt (510 行)
- 涉枪涉爆.txt (435 行)
- 色情词库.txt (578 行)
- 贪腐词库.txt (240 行)

**总计：2915 个敏感词**

### 🔄 向后兼容性

- ✅ 所有现有 API 保持完全兼容
- ✅ 新增功能均为可选，不影响现有代码
- ✅ 无需修改即可升级

---

## [1.0.0] - 2024-10-29

### 🎉 第一个正式版本发布

`go-sensitive-word` 是一个基于 DFA/AC 算法实现的高性能 Go 敏感词过滤工具框架，支持动态维护、智能归一化、多种词库加载方式等功能。

---

## ✨ 核心特性

### 1. 双算法支持

- ✅ **DFA 算法**：实现简单，适合小词库场景（< 1000 词）
- ✅ **AC 自动机算法**：高性能，推荐生产环境使用
  - 单次线性扫描，时间复杂度 O(n)
  - 支持窗口合并和原子切换，降低重建频率
  - 适合长文本和大词库场景（> 1000 词）

### 2. 文本检测功能

提供完整的敏感词检测 API：

- ✅ `IsSensitive()` - 判断文本中是否存在敏感词
- ✅ `FindOne()` - 查找文本中的第一个敏感词
- ✅ `FindAll()` - 查找文本中所有敏感词（去重）
- ✅ `FindAllCount()` - 查找所有敏感词及其出现次数
- ✅ `Replace()` - 替换所有敏感词为指定字符
- ✅ `Remove()` - 从文本中删除所有敏感词

### 3. 动态词库管理

支持运行时动态更新词库，无需重启服务：

- ✅ `AddWord()` / `AddWords()` - 添加敏感词（支持单个和批量）
- ✅ `DelWord()` / `DelWords()` - 删除敏感词（支持单个和批量）
- ✅ `ReplaceWords()` - 批量替换敏感词（先删旧词，再加新词）
- ✅ `Clear()` - 清空词库
- ✅ `GetStats()` - 获取词库统计信息（总词数、最后更新时间等）
- ✅ `ExportToString()` / `ExportToFile()` - 导出词库
- ✅ `MergeFromManager()` - 合并另一个 Manager 的词库
- ✅ `RefreshFromPath()` - 从文件路径刷新词库（支持替换/追加模式）

### 4. 多种词库加载方式

灵活的词库加载机制，适应不同使用场景：

- ✅ **内置词库加载** (`LoadDictEmbed`)：编译时嵌入，支持 15+ 个分类词库
  - 政治类型、暴恐词库、色情类型、涉枪涉爆、广告类型等
- ✅ **文件加载** (`LoadDictPath`)：从本地文件加载，支持多个文件
- ✅ **回调函数加载** (`LoadDictCallback`)：支持从数据库、Redis、配置中心等自定义数据源加载

### 5. 智能归一化

内置强大的文本归一化功能，有效防御混淆攻击：

- ✅ **忽略大小写**（默认开启）：`"FuCk"` → `"fuck"`
- ✅ **全角转半角**（默认开启）：`"ｆｕｃｋ"` → `"fuck"`
- ✅ **数字归一化**：统一各种数字写法为阿拉伯数字
- ✅ **繁简归一**：繁体字转简体字
- ✅ **英文变体归一**：花体、数学字母等 → 基本拉丁字母
- ✅ **零宽字符剔除**：防止零宽字符绕过检测
- ✅ **同形字映射**：防御视觉混淆字符攻击

归一化配置：
- 默认配置：启用忽略大小写和全角转半角
- 严格模式：`StrictNormalizer()` 启用所有防绕过选项

### 6. 工具函数

内置敏感信息检测和屏蔽功能：

- ✅ **邮箱检测和屏蔽**：`HasEmail()`, `MaskEmail()`
- ✅ **URL 检测和屏蔽**：`HasURL()`, `MaskURL()`
- ✅ **数字检测和屏蔽**：`HasDigit()`, `MaskDigit()`
- ✅ **微信号检测和屏蔽**：`HasWechatID()`, `MaskWechatID()`

### 7. 并发安全与资源管理

- ✅ **全链路并发安全**：支持高并发场景，所有操作线程安全
- ✅ **优雅关闭**：`Close()` 和 `Shutdown(ctx)` 方法，避免 goroutine 泄漏
- ✅ **异步处理**：词库更新通过 channel 异步处理，不阻塞主流程

### 8. 性能优化

- ✅ **AC 算法优化**：
  - 使用 `atomic.Pointer` 实现原子切换，读操作无锁
  - 窗口合并机制，批量更新减少重建次数
  - Copy-on-Write (COW) 模式，避免写操作阻塞读操作
- ✅ **DFA 算法优化**：
  - 简化实现，内存占用低
  - 适合小词库场景
- ✅ **区间映射优化**：
  - `FindAllRanges` 接口，避免重复字符串搜索
  - 支持重叠区间处理，提高 Replace/Remove 性能

---

## 📚 文档与示例

### 完整文档

提供 11 个详细文档，覆盖所有功能：

- 📖 [快速开始指南](docs/getting-started.md) - 5分钟上手
- 📖 [API 参考文档](docs/api-reference.md) - 完整 API 文档
- 📖 [归一化功能详解](docs/normalization.md) - 归一化策略说明
- 📖 [算法选择指南](docs/algorithm-guide.md) - DFA vs AC 算法选择
- 📖 [词库管理详解](docs/word-management.md) - 动态维护词库
- 📖 [词库加载详解](docs/word-loading.md) - 多种加载方式
- 📖 [来源追踪详解](docs/word-source-tracking.md) - 词库来源追踪功能
- 📖 [工具函数详解](docs/tools.md) - 敏感信息检测和屏蔽
- 📖 [资源管理详解](docs/lifecycle.md) - 生命周期和优雅关闭
- 📖 [常见问题](docs/faq.md) - 20+ FAQ
- 📖 [最佳实践](docs/best-practices.md) - 生产环境部署建议

### 丰富示例

提供 11 个示例代码，覆盖所有核心功能：

- 📝 [基础功能演示](examples/basic/main.go) - 快速入门
- 📝 [AC 算法示例](examples/ac/main.go) - AC 算法使用
- 📝 [动态维护示例](examples/dynamic/main.go) - 词库动态更新
- 📝 [回调加载示例](examples/callback/main.go) - 从数据库/Redis 加载
- 📝 [文件加载示例](examples/file-load/main.go) - 从文件加载
- 📝 [归一化配置示例](examples/normalize/main.go) - 归一化配置
- 📝 [工具函数示例](examples/tools/main.go) - 敏感信息检测
- 📝 [资源管理示例](examples/lifecycle/main.go) - 优雅关闭
- 📝 [多实例示例](examples/multi-instance/main.go) - 多实例使用
- 📝 [来源追踪示例](examples/word-source/main.go) - 来源追踪功能
- 📝 [综合功能演示](examples/comprehensive/main.go) - 完整功能演示

---

## 🛠️ 技术实现

### 架构设计

- **模块化设计**：分离存储、过滤、归一化模块
- **接口抽象**：`store.Store`, `filter.Filter`, `filter.RangedFilter`
- **并发安全**：使用 `sync.RWMutex`, `atomic.Pointer`, `atomic.Int64`

### 核心组件

- **Manager**：核心管理器，整合存储和过滤
- **MemoryStore**：内存词库存储，使用 `atomic.Int64` 计数
- **DFAModel**：DFA 算法实现
- **ACModel**：AC 自动机算法实现，支持窗口合并
- **Normalizer**：文本归一化处理
- **Wrapper**：归一化包装器，支持区间映射

---

## 📦 内置词库

项目内置了 8 个核心分类词库（编译时嵌入）：

- 反动词库 (551 个词)
- 广告类型 (120 个词)
- 政治类型 (303 个词)
- 暴恐词库 (178 个词)
- 民生词库 (510 个词)
- 涉枪涉爆 (435 个词)
- 色情词库 (578 个词)
- 贪腐词库 (240 个词)

**总计：2915 个敏感词**

---

## 🚀 快速开始

```bash
# 安装
go get -u github.com/LuYongwang/go-sensitive-word@latest
```

```go
package main

import (
    "fmt"
    "time"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    filter, _ := sensitive.NewFilter(
        sensitive.StoreOption{Type: sensitive.StoreMemory},
        sensitive.FilterOption{Type: sensitive.FilterAC}, // 推荐生产环境使用 AC
    )
    
    filter.LoadDictEmbed(sensitive.DictPolitical)
    filter.AddWord("敏感词")
    time.Sleep(100 * time.Millisecond)
    
    fmt.Println(filter.IsSensitive("包含敏感词的文本")) // true
}
```

---

## 📋 项目信息

- **项目名称**：go-sensitive-word
- **最新版本**：v1.1.0
- **Go 版本要求**：Go 1.20+
- **许可证**：待定
- **仓库地址**：github.com/LuYongwang/go-sensitive-word

---

## 🙏 致谢

感谢所有为项目做出贡献的开发者！

---

## 📝 后续计划

- [ ] 支持 Redis 存储后端
- [ ] 支持数据库存储后端
- [ ] 支持更多归一化策略
- [ ] 性能基准测试报告
- [ ] 更多工具函数

---

**注意**：这是第一个正式版本，API 保持稳定，后续版本将保持向后兼容。
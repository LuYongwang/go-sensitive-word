# 示例代码目录

本目录包含了 `go-sensitive-word` 项目的所有功能演示示例。

## 📁 目录结构

```
examples/
├── basic/
│   └── main.go                # 基础功能演示（快速入门）
├── ac/
│   └── main.go                # AC 自动机算法使用示例
├── dynamic/
│   └── main.go                # 动态维护词库功能演示
├── callback/
│   └── main.go                # 回调函数加载词库示例
├── file-load/
│   └── main.go                # 从文件加载词库示例
├── normalize/
│   └── main.go                # 归一化配置示例
├── tools/
│   └── main.go                # 工具函数使用示例
├── lifecycle/
│   └── main.go                # 资源管理示例
├── multi-instance/
│   └── main.go                # 多实例使用示例
├── word-source/
│   └── main.go                # 来源追踪功能演示
└── comprehensive/
    └── main.go                # 综合功能演示
```

## 🚀 快速开始

### 1. 基础示例（快速入门）

```bash
go run examples/basic/main.go
```

演示内容：
- 过滤器的初始化和配置
- 词库加载（内置词库）
- 动态添加/删除敏感词
- 所有文本检测功能（IsSensitive, FindOne, FindAll, FindAllCount, Replace, Remove）

### 2. AC 算法示例

```bash
go run examples/ac/main.go
```

演示内容：
- AC 自动机算法的使用
- AC 算法的高性能特性
- 适用于生产环境的配置

### 3. 动态维护示例

```bash
go run examples/dynamic/main.go
```

演示内容：
- 批量添加/删除/替换敏感词
- 获取词库统计信息
- 词库导出（字符串/文件）
- 词库合并

### 4. 回调函数加载示例

```bash
go run examples/callback/main.go
```

演示内容：
- 从数据库加载词库
- 从 Redis 加载词库
- 从配置中心加载词库
- 从多个数据源合并加载
- 使用内联匿名函数

### 5. 文件加载示例

```bash
go run examples/file-load/main.go
```

演示内容：
- 从文件路径加载词库（追加模式）
- 从文件路径刷新词库（替换模式）
- 加载多个文件

### 6. 归一化配置示例

```bash
go run examples/normalize/main.go
```

演示内容：
- 默认归一化配置演示
- 全角/半角字符处理
- 大小写忽略
- 自定义归一化配置

### 7. 工具函数示例

```bash
go run examples/tools/main.go
```

演示内容：
- 邮箱检测和屏蔽
- URL 检测和屏蔽
- 数字检测和屏蔽
- 微信号检测和屏蔽
- 组合使用示例
- 结合敏感词过滤

### 8. 资源管理示例

```bash
go run examples/lifecycle/main.go
```

演示内容：
- 基本关闭（Close）
- 优雅关闭（Shutdown）
- 生产环境最佳实践
- 并发安全演示

### 9. 多实例示例

```bash
go run examples/multi-instance/main.go
```

演示内容：
- 创建多个独立的过滤器实例
- 不同业务场景使用不同词库
- 实例间互不干扰

### 10. 来源追踪示例

```bash
go run examples/word-source/main.go
```

演示内容：
- 为词库指定来源标识
- 查询敏感词的来源
- 查找敏感词及其来源信息
- 获取所有词的来源分布

### 11. 综合功能示例

```bash
go run examples/comprehensive/main.go
```

演示内容：
- **最完整的功能演示**，覆盖所有核心功能
- 适合作为参考实现

## 📚 示例详细说明

### basic/main.go - 基础功能演示

**适合人群：** 新手快速入门

**演示功能：**
- ✅ 过滤器初始化
- ✅ 加载内置词库
- ✅ 动态添加/删除敏感词
- ✅ 所有文本检测和处理功能

### ac/main.go - AC 算法示例

**适合人群：** 需要高性能的生产环境

**演示功能：**
- ✅ AC 自动机算法的使用
- ✅ 与 DFA 算法的区别
- ✅ 生产环境推荐配置

### dynamic/main.go - 动态维护示例

**适合人群：** 需要运行时更新词库的场景

**演示功能：**
- ✅ `AddWords()` - 批量添加
- ✅ `DelWords()` - 批量删除
- ✅ `ReplaceWords()` - 批量替换
- ✅ `GetStats()` - 统计信息
- ✅ `ExportToString()` / `ExportToFile()` - 导出
- ✅ `MergeFromManager()` - 词库合并

### callback/main.go - 回调函数加载

**适合人群：** 词库存储在数据库/Redis/配置中心等

**演示功能：**
- ✅ `LoadDictCallback()` - 回调函数加载
- ✅ 从数据库加载示例
- ✅ 从 Redis 加载示例
- ✅ 从配置中心加载示例
- ✅ 多数据源合并加载

### file-load/main.go - 文件加载示例

**适合人群：** 词库存储在本地文件

**演示功能：**
- ✅ `LoadDictPath()` - 从文件加载（追加模式）
- ✅ `RefreshFromPath()` - 从文件刷新（替换模式）
- ✅ 加载多个文件

### normalize/main.go - 归一化配置示例

**适合人群：** 需要防御混淆攻击的场景

**演示功能：**
- ✅ 默认归一化配置（忽略大小写、全角转半角）
- ✅ 大小写归一化演示
- ✅ 全角/半角字符处理
- ✅ 自定义归一化配置

### tools/main.go - 工具函数示例

**适合人群：** 需要检测和屏蔽邮箱、URL、微信号等敏感信息

**演示功能：**
- ✅ `HasEmail()` / `MaskEmail()` - 邮箱
- ✅ `HasURL()` / `MaskURL()` - URL
- ✅ `HasDigit()` / `MaskDigit()` - 数字
- ✅ `HasWechatID()` / `MaskWechatID()` - 微信号
- ✅ 组合使用示例

### lifecycle/main.go - 资源管理示例

**适合人群：** 生产环境部署

**演示功能：**
- ✅ `Close()` - 基本关闭
- ✅ `Shutdown()` - 优雅关闭
- ✅ 生产环境最佳实践
- ✅ 并发安全测试

### multi-instance/main.go - 多实例示例

**适合人群：** 需要多个独立词库的业务场景

**演示功能：**
- ✅ 创建多个独立的过滤器实例
- ✅ 不同实例加载不同词库
- ✅ 实例间数据隔离

### word-source/main.go - 来源追踪示例

**适合人群：** 需要追踪敏感词来源的场景

**演示功能：**
- ✅ `LoadDictEmbedWithSource()` - 加载词库并指定来源
- ✅ `AddWordsWithSource()` - 添加词并指定来源
- ✅ `GetWordSources()` - 查询单个词的来源
- ✅ `FindAllWithSource()` - 查找词及其来源
- ✅ `GetAllWordSources()` - 获取所有词的来源映射

### comprehensive/main.go - 综合功能演示

**适合人群：** 了解项目全貌

**演示功能：**
- ✅ **所有核心功能的完整演示**
- ✅ 从初始化到资源管理的完整流程
- ✅ 适合作为参考实现

## 🔍 功能覆盖清单

| 功能分类 | 功能点 | 覆盖示例 |
|---------|--------|---------|
| **初始化** | NewFilter | basic/main.go, comprehensive/main.go |
| | StoreMemory | 所有示例 |
| | FilterDfa / FilterAC | basic/main.go, ac/main.go |
| **词库加载** | LoadDictEmbed | basic/main.go, comprehensive/main.go |
| | LoadDictPath | file-load/main.go |
| | LoadDictCallback | callback/main.go |
| **文本检测** | IsSensitive | 所有示例 |
| | FindOne | basic/main.go, comprehensive/main.go |
| | FindAll | 所有示例 |
| | FindAllCount | basic/main.go, comprehensive/main.go |
| **文本处理** | Replace | basic/main.go, comprehensive/main.go |
| | Remove | basic/main.go, comprehensive/main.go |
| **动态维护** | AddWord / AddWords | dynamic/main.go, comprehensive/main.go |
| | DelWord / DelWords | dynamic/main.go, comprehensive/main.go |
| | ReplaceWords | dynamic/main.go, comprehensive/main.go |
| **词库管理** | GetStats | dynamic/main.go, comprehensive/main.go |
| | ExportToString / ExportToFile | dynamic/main.go, comprehensive/main.go |
| | MergeFromManager | dynamic/main.go, comprehensive/main.go |
| | RefreshFromPath | file-load/main.go |
| | Clear | 未单独演示，可参考 comprehensive/main.go |
| **多实例** | 多个独立实例 | multi-instance/main.go |
| **来源追踪** | LoadDictEmbedWithSource | word-source/main.go |
| | AddWordsWithSource | word-source/main.go |
| | GetWordSources | word-source/main.go |
| | FindAllWithSource | word-source/main.go |
| | GetAllWordSources | word-source/main.go |
| **归一化** | 默认归一化 | normalize/main.go, comprehensive/main.go |
| | 自定义归一化 | normalize/main.go |
| **工具函数** | HasEmail / MaskEmail | tools/main.go |
| | HasURL / MaskURL | tools/main.go |
| | HasDigit / MaskDigit | tools/main.go |
| | HasWechatID / MaskWechatID | tools/main.go |
| **资源管理** | Close | lifecycle/main.go, comprehensive/main.go |
| | Shutdown | lifecycle/main.go, comprehensive/main.go |

## 💡 使用建议

1. **新手入门**：从 `basic/main.go` 开始，了解基本用法
2. **生产环境**：参考 `ac/main.go` 和 `lifecycle/main.go`
3. **需要动态更新**：查看 `dynamic/main.go`
4. **词库在外部系统**：参考 `callback/main.go`
5. **多业务场景**：参考 `multi-instance/main.go`
6. **需要来源追踪**：查看 `word-source/main.go`
7. **需要完整参考**：运行 `comprehensive/main.go`

## 🐛 问题反馈

如果示例代码无法运行或发现问题，请提交 Issue。

## 📖 相关文档

- [主项目 README](../README.md)
- [动态维护文档](../docs/dynamic-maintenance.md)
- [回调加载文档](../docs/callback-loader.md)


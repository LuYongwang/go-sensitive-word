# go-sensitive-word

敏感词（敏感词/违禁词/违法词/脏词）检测工具，基于 DFA/AC 算法实现的高性能 Go 敏感词过滤工具框架。

## ✨ 特性

- 🚀 **高性能**: 支持 DFA 和 AC 自动机两种算法，AC 算法适合生产环境
- 🔄 **动态维护**: 支持运行时动态添加/删除/替换敏感词，无需重启
- 🛡️ **智能归一化**: 内置多种归一化策略，有效防御混淆攻击（大小写、全角、相似字符等）
- 📦 **易用性**: 丰富的 API 和示例代码，快速接入
- 🔌 **灵活加载**: 支持内置词库、文件、数据库、Redis 等多种加载方式
- 🛠️ **工具函数**: 内置邮箱、URL、微信号等敏感信息检测和屏蔽功能
- 🔒 **并发安全**: 全链路并发安全，支持高并发场景
- ♻️ **资源管理**: 支持优雅关闭，避免 goroutine 泄漏

## 📋 目录结构

```
go-sensitive-word/
├── examples/              # 示例代码目录
│   ├── basic/             # 基础功能演示
│   ├── ac/                # AC 算法示例
│   ├── dynamic/           # 动态维护示例
│   ├── callback/          # 回调加载示例
│   ├── file-load/         # 文件加载示例
│   ├── normalize/         # 归一化配置示例
│   ├── tools/             # 工具函数示例
│   ├── lifecycle/         # 资源管理示例
│   └── comprehensive/     # 综合功能演示
├── docs/                  # 文档目录
│   ├── getting-started.md        # 快速开始指南
│   ├── api-reference.md           # API 参考文档
│   ├── normalization.md          # 归一化功能详解
│   ├── algorithm-guide.md         # 算法选择指南
│   ├── word-management.md         # 词库管理详解
│   ├── word-loading.md            # 词库加载详解
│   ├── tools.md                   # 工具函数详解
│   ├── lifecycle.md               # 资源管理详解
│   ├── faq.md                     # 常见问题
│   └── best-practices.md          # 最佳实践
├── internal/             # 内部实现
│   ├── filter/           # 过滤算法（DFA/AC）
│   ├── store/           # 词库存储
│   ├── normalize/       # 文本归一化
│   └── jianfan/         # 繁简转换
├── wordlists/           # 内置敏感词库
├── manager.go          # 核心管理器
├── wrapper.go          # 归一化包装器
├── tool.go             # 工具函数
└── README.md           # 本文档
```

## 快速接入

**安装**
```bash
go get -u github.com/LuYongwang/go-sensitive-word@latest
```

**使用示例**
```go
package main

import (
   "fmt"
   "time"
   sensitive "github.com/LuYongwang/go-sensitive-word"
   "log"
)

func main() {
   filter, err := sensitive.NewFilter(
      sensitive.StoreOption{Type: sensitive.StoreMemory}, // 基于内存
      sensitive.FilterOption{Type: sensitive.FilterDfa},  // 基于DFA算法
   )

   if err != nil {
      log.Fatalf("敏感词服务启动失败, err:%v", err)
      return
   }

   // 加载敏感词库
   err = filter.LoadDictEmbed(
      sensitive.DictGFWAdditional,
      sensitive.DictOther,
      // ... 其他词库
   )
   if err != nil {
      log.Fatalf("加载词库发生了错误, err:%v", err)
      return
   }

   // 动态添加自定义敏感词
   err = filter.AddWord("敏感词1", "敏感词2")
   if err != nil {
      log.Fatalf("添加敏感词发生了错误, err:%v", err)
      return
   }

   // 等待异步处理完成
   time.Sleep(100 * time.Millisecond)

   text := "这是一段包含敏感词的文本"

   // 是否有敏感词
   res1 := filter.IsSensitive(text)
   fmt.Printf("res1: %v \n", res1)

   // 找到所有敏感词
   res2 := filter.FindAll(text)
   fmt.Printf("res2: %v \n", res2)

   // 替换敏感词
   res3 := filter.Replace(text, '*')
   fmt.Printf("res3: %v \n", res3)
}
```

## 核心功能

### 文本检测
- `IsSensitive()` - 判断文本中是否存在敏感词
- `FindOne()` - 查找文本中的第一个敏感词
- `FindAll()` - 查找文本中所有敏感词（去重）
- `FindAllCount()` - 查找所有敏感词及其出现次数
- `Replace()` - 替换所有敏感词为指定字符
- `Remove()` - 从文本中删除所有敏感词

### 词库管理
- `AddWord()` / `AddWords()` - 动态添加敏感词
- `DelWord()` / `DelWords()` - 动态删除敏感词
- `ReplaceWords()` - 批量替换敏感词
- `Clear()` - 清空词库
- `GetStats()` - 获取词库统计信息
- `ExportToString()` / `ExportToFile()` - 导出词库

### 词库加载
- `LoadDictEmbed()` - 加载内置嵌入的词库（编译时嵌入）
- `LoadDictPath()` - 从文件路径加载词库
- `LoadDictCallback()` - 通过回调函数加载词库（支持数据库、Redis等）

## 📖 文档导航

- [快速开始指南](./docs/getting-started.md) - 新手入门，5分钟上手
- [API 参考文档](./docs/api-reference.md) - 完整的 API 文档
- [归一化功能详解](./docs/normalization.md) - 智能归一化策略说明
- [算法选择指南](./docs/algorithm-guide.md) - DFA vs AC 算法选择建议
- [词库管理详解](./docs/word-management.md) - 动态维护词库指南
- [词库加载详解](./docs/word-loading.md) - 多种词库加载方式
- [工具函数详解](./docs/tools.md) - 敏感信息检测和屏蔽
- [资源管理详解](./docs/lifecycle.md) - 生命周期和优雅关闭
- [常见问题](./docs/faq.md) - FAQ 和问题解答
- [最佳实践](./docs/best-practices.md) - 生产环境部署建议

## 📝 示例代码

项目提供了丰富的示例代码，覆盖所有核心功能：

| 示例文件 | 说明 |
|---------|------|
| [examples/basic/main.go](./examples/basic/main.go) | 基础功能演示（快速入门） |
| [examples/ac/main.go](./examples/ac/main.go) | AC 自动机算法使用示例 |
| [examples/dynamic/main.go](./examples/dynamic/main.go) | 动态维护词库功能演示 |
| [examples/callback/main.go](./examples/callback/main.go) | 回调函数加载词库示例 |
| [examples/file-load/main.go](./examples/file-load/main.go) | 从文件加载词库示例 |
| [examples/normalize/main.go](./examples/normalize/main.go) | 归一化配置示例 |
| [examples/tools/main.go](./examples/tools/main.go) | 工具函数使用示例 |
| [examples/lifecycle/main.go](./examples/lifecycle/main.go) | 资源管理示例 |
| [examples/comprehensive/main.go](./examples/comprehensive/main.go) | 综合功能演示 |

**运行示例：**
```bash
# 运行基础示例（快速入门）
go run examples/basic/main.go

# 运行 AC 算法示例
go run examples/ac/main.go

# 运行动态维护示例
go run examples/dynamic/main.go
```

详细示例说明请查看：[examples/README.md](./examples/README.md)

## 参考资料
- 基于Java DFA实现的敏感词过滤：https://github.com/houbb/sensitive-word
- unicode字词的神奇组合：https://www.zhihu.com/question/30873035
- unicode违规技巧：https://zhuanlan.zhihu.com/p/545309061
- unicode视觉欺骗：https://zhuanlan.zhihu.com/p/611904676
- unicode字符列表：https://symbl.cc/en/unicode-table
- 汉字结构描述字符：https://zh.wikipedia.org/wiki/%E8%A1%A8%E6%84%8F%E6%96%87%E5%AD%97%E6%8F%8F%E8%BF%B0%E5%AD%97%E7%AC%A6
- 敏感词库：https://github.com/konsheng/Sensitive-lexicon

## 声明

本项目包含了一些敏感词库，其设计目的是为了解决在互联网环境中可能出现的不适当内容，通过技术手段屏蔽这些敏感词，旨在构建一个更健康、更安全的网络空间。

请注意以下几点：

1. **项目目的**：本项目的初衷是为开发者提供一个工具，帮助他们在各类互联网产品中过滤和屏蔽不适当或敏感的内容，从而营造一个良好的网络生态环境。

2. **使用限制**：本项目中的敏感词库仅供技术研究和实现内容过滤功能之用。任何个人或组织不得将本项目中的敏感词库用于传播、分享或其他任何可能导致敏感信息扩散的行为。

3. **责任声明**：使用本项目所产生的任何直接或间接后果，均由使用者自行承担。本项目开发者不对因不当使用造成的任何损失或法律后果负责。

4. **使用规范**：请确保在使用本项目时遵守相关法律法规。禁止将本项目用于任何违法或不正当的用途。

通过下载和使用本项目，即表示您同意并接受上述声明的所有内容。希望本项目能够为您在构建净化网络空间的过程中提供帮助。我们鼓励所有开发者共同努力，营造一个健康、安全的网络环境。
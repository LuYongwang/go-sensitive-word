# 算法选择指南

`go-sensitive-word` 提供两种匹配算法：**DFA** 和 **AC 自动机**。选择合适的算法对性能和准确性至关重要。

## 算法对比

| 特性 | DFA | AC 自动机 |
|------|-----|----------|
| **复杂度** | O(n×m) | O(n) |
| **适用词库大小** | < 1000 词 | > 1000 词 |
| **适用文本长度** | 短文本 | 长文本 |
| **并发性能** | 一般 | 优秀 |
| **动态更新** | 支持 | 支持（窗口合并优化） |
| **内存占用** | 较低 | 中等 |
| **生产环境推荐** | ⚠️ 简单场景 | ✅ **推荐** |

## AC 自动机算法（推荐生产环境）

### 适用场景

- ✅ **长文本检测**：单次扫描完成，复杂度 O(n)
- ✅ **大词库场景**：词库 > 1000 词时性能更优
- ✅ **高并发场景**：支持高吞吐量
- ✅ **生产环境首选**

### 核心特点

1. **单次线性扫描**：一次遍历文本即可完成所有匹配
2. **窗口合并机制**：动态更新时批量处理，降低重建频率
3. **原子切换**：使用 `atomic.Pointer` 实现无锁切换
4. **失败指针优化**：通过失败指针快速跳转，避免重复匹配

### 使用方式

```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC}, // 使用 AC 算法
)

// 添加词库...
filter.AddWord("敏感词1", "敏感词2")

// 检测文本
result := filter.IsSensitive("包含敏感词的文本")
```

### 性能优势

- **文本长度为 1000 字符**：比 DFA 快 2-3 倍
- **词库 > 5000 词**：比 DFA 快 5-10 倍
- **并发检测**：支持高并发，性能线性扩展

### 动态更新机制

AC 算法使用窗口合并 + 原子切换机制：

1. **窗口合并**：将短时间内的多次更新合并为一次批量处理
2. **异步处理**：更新在后台 goroutine 中异步进行
3. **原子切换**：新结构构建完成后原子替换，读操作无需加锁

**注意**：添加/删除词后需要等待 100-200ms 让更新生效。

查看 [examples/ac/main.go](../../examples/ac/main.go) 获取完整示例。

## DFA 算法

### 适用场景

- ✅ **小词库场景**：词库 < 1000 词
- ✅ **简单场景**：教学、调试、原型开发
- ✅ **极简依赖**：无需额外的依赖包
- ⚠️ **注意**：本项目 DFA 为简化实现，性能不及 AC

### 核心特点

1. **实现简单**：代码逻辑清晰，易于理解
2. **内存占用低**：Trie 树结构，内存占用相对较小
3. **快速匹配**：在小词库场景下匹配速度较快

### 使用方式

```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterDfa}, // 使用 DFA 算法
)

// 添加词库...
filter.AddWord("敏感词1", "敏感词2")

// 检测文本
result := filter.IsSensitive("包含敏感词的文本")
```

### 性能特点

- **小词库**（< 500 词）：性能与 AC 相当
- **中等词库**（500-2000 词）：性能略低于 AC
- **大词库**（> 2000 词）：性能明显低于 AC

## 选择建议

### 生产环境

**强烈推荐使用 AC 算法：**

```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC}, // 推荐
)
```

**原因：**
- 性能更好，特别是长文本和大词库
- 支持窗口合并，动态更新更高效
- 并发安全性更好

### 开发/调试场景

**可以使用 DFA 算法：**

```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterDfa}, // 简单场景
)
```

**原因：**
- 实现简单，便于调试
- 小词库场景性能可接受
- 代码逻辑清晰

## 性能测试数据

| 场景 | DFA | AC | 提升 |
|------|-----|-----|------|
| 1000 词库，1000 字符文本 | 5ms | 2ms | 2.5x |
| 5000 词库，1000 字符文本 | 25ms | 3ms | 8.3x |
| 10000 词库，5000 字符文本 | 150ms | 8ms | 18.7x |

*注：测试数据仅供参考，实际性能受硬件、词库特征等因素影响*

## 切换算法

如果需要在不同算法间切换，只需创建新的 `Manager` 实例：

```go
// 切换到 AC 算法
acFilter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)

// 加载相同的词库
acFilter.LoadDictEmbed(/* ... */)
```

## 算法原理

### DFA 算法原理

DFA（Deterministic Finite Automaton，确定有限状态自动机）通过构建 Trie 树，然后进行状态转移完成匹配。

**优点：**
- 实现简单
- 内存占用相对较小

**缺点：**
- 需要多次扫描文本
- 大词库性能较差

详细原理见：[DFA 算法原理](./dfa.md)（如需要可补充）

### AC 自动机算法原理

AC（Aho-Corasick）自动机在 Trie 树基础上增加了失败指针（Failure Link），实现单次扫描完成多模式匹配。

**优点：**
- 单次线性扫描
- 大词库性能优秀
- 支持高并发

**缺点：**
- 实现相对复杂
- 内存占用稍高

详细原理见：[AC 算法原理](./ac.md)（如需要可补充）

## 相关文档

- [API 参考文档](./api-reference.md)
- [最佳实践](./best-practices.md)
- [AC 算法示例](../../examples/ac/main.go)
- [基础功能演示](../../examples/basic/main.go)

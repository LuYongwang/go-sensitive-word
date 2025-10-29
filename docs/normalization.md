# 归一化功能详解

`go-sensitive-word` 内置了强大的文本归一化功能，可以有效防御各种混淆攻击。

## 为什么需要归一化？

攻击者经常使用各种技巧绕过敏感词检测：

- **大小写混淆**: `FuCk` → `fuck`
- **全角/半角混淆**: `ｆｕｃｋ` → `fuck`
- **相似字符替换**: `Ⓕⓤc⒦` → `fuck`
- **零宽字符插入**: `f​uck`（中间有零宽字符）
- **繁简混用**: `五星紅旗` → `五星红旗`

归一化功能将这些变体统一为标准形式，确保检测的准确性。

## 归一化策略

### 1. 忽略大小写（IgnoreCase）

自动将所有字母转换为小写进行匹配。

**示例：**
```
"FuCk" → "fuck"
"TEST" → "test"
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    IgnoreCase: true, // 默认开启
}
```

### 2. 全角转半角（ToHalfWidth）

将全角字母、数字、符号转换为半角。

**示例：**
```
"ｆｕｃｋ" → "fuck"
"１２３" → "123"
"（" → "("
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    ToHalfWidth: true, // 默认开启
}
```

### 3. 数字归一化（IgnoreDigitType）

统一各种数字写法为阿拉伯数字。

**示例：**
```
"9⓿二肆⁹₈③⑸⒋" → "902438354"
"①②③" → "123"
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    IgnoreDigitType: true, // 需要手动开启
}
```

### 4. 繁简归一（IgnoreSimpTrad）

将繁体字转换为简体字。

**示例：**
```
"五星紅旗" → "五星红旗"
"繁體字" → "繁体字"
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    IgnoreSimpTrad: true, // 需要手动开启
}
```

### 5. 英文变体归一（IgnoreEnglishStyle）

将花体、数学字母等变体转换为基本拉丁字母。

**示例：**
```
"Ⓕⓤc⒦" → "fuck"
"𝕋𝔼𝕊𝕋" → "TEST"
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    IgnoreEnglishStyle: true, // 需要手动开启
}
```

### 6. 剔除零宽字符（RemoveZeroWidth）

删除文本中的零宽字符（Zero-Width Characters），防止绕过检测。

**零宽字符类型：**
- Zero Width Space (U+200B)
- Zero Width Non-Joiner (U+200C)
- Zero Width Joiner (U+200D)
- Left-to-Right Mark (U+200E)
- Right-to-Left Mark (U+200F)
- 等

**配置：**
```go
config := sensitive.NormalizerConfig{
    RemoveZeroWidth: true, // 需要手动开启
}
```

### 7. 同形字映射（HomoglyphMap）

将视觉相似的字符映射为标准字符，防止混淆攻击。

**示例：**
```
"а" (西里尔字母) → "a" (拉丁字母)
"о" (西里尔字母) → "o" (拉丁字母)
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    HomoglyphMap: normalize.DefaultHomoglyphMap(), // 使用默认映射表
}
```

### 8. 忽略连续重复字符（IgnoreRepeat）

将连续相同的字符压缩为 1 个。

**示例：**
```
"f***k" → "fk"
"测试试试试" → "测试"
```

**配置：**
```go
config := sensitive.NormalizerConfig{
    IgnoreRepeat: true, // 需要手动开启
}
```

## 默认配置

当前版本默认启用：

- ✅ **忽略大小写**（IgnoreCase）
- ✅ **全角转半角**（ToHalfWidth）

其他策略需要手动配置开启。

## 使用方式

### 基本使用

默认配置已足够应对大多数场景，无需特殊配置：

```go
filter, _ := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)

// 默认已启用：忽略大小写、全角转半角
filter.AddWord("敏感词")
filter.IsSensitive("敏感词") // true
filter.IsSensitive("敏Ｇ词") // true（全角G会被转半角）
```

### 自定义配置

如需更严格的归一化，可以使用 `StrictNormalizer()`：

```go
// 注意：当前版本归一化配置在内部处理，用户无需显式配置
// 如需自定义，请参考 normalize.go 中的实现
```

查看 [examples/normalize/main.go](../../examples/normalize/main.go) 了解更多示例。

## 工作原理

1. **文本归一化**：检测前先将文本转为归一化形式
2. **词库归一化**：添加词库时同样归一化
3. **匹配检测**：在归一化后的文本上进行匹配
4. **结果映射**：将结果映射回原文位置进行处理

**重要**：归一化过程会保留原始文本的位置映射，确保替换/删除操作在原文上正确执行。

## 防御效果

| 攻击方式 | 原始文本 | 归一化后 | 检测结果 |
|---------|---------|---------|---------|
| 大小写混淆 | `FuCk` | `fuck` | ✅ 能检测 |
| 全角/半角 | `ｆｕｃｋ` | `fuck` | ✅ 能检测 |
| 零宽字符 | `f​uck` | `fuck` | ✅ 能检测（开启后） |
| 相似字符 | `Ⓕⓤc⒦` | `fuck` | ✅ 能检测（开启后） |
| 繁简混用 | `五星紅旗` | `五星红旗` | ✅ 能检测（开启后） |

## 性能考虑

- **默认配置**：性能影响很小，建议保持开启
- **严格配置**：会略微增加处理时间，但能显著提升安全性
- **生产环境建议**：根据实际需求选择合适的策略组合

## 相关文档

- [API 参考文档](./api-reference.md)
- [最佳实践](./best-practices.md)
- [归一化配置示例](../../examples/normalize/main.go)

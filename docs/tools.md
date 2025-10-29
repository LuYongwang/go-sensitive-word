# 工具函数详解

除了敏感词检测，`go-sensitive-word` 还提供了一系列工具函数，用于检测和屏蔽邮箱、URL、微信号等敏感信息。

## 功能概览

| 功能 | 检测函数 | 屏蔽函数 |
|------|---------|---------|
| **邮箱** | `HasEmail` | `MaskEmail` |
| **URL** | `HasURL` | `MaskURL` |
| **数字** | `HasDigit` | `MaskDigit` |
| **微信号** | `HasWechatID` | `MaskWechatID` |

## 邮箱检测与屏蔽

### HasEmail

判断字符串中是否存在邮箱地址。

```go
func HasEmail(s string) bool
```

**示例：**
```go
if sensitive.HasEmail("联系我: user@example.com") {
    fmt.Println("包含邮箱")
}
```

### MaskEmail

将字符串中的邮箱地址替换为 `***`。

```go
func MaskEmail(s string) string
```

**示例：**
```go
text := "联系我: user@example.com 或 admin@test.org"
result := sensitive.MaskEmail(text)
// 输出: "联系我: *** 或 ***"
```

**匹配规则：**
- 标准邮箱格式：`user@domain.com`
- 支持子域名：`user@mail.example.com`
- 支持特殊字符：`user+tag@example.com`

## URL 检测与屏蔽

### HasURL

判断字符串中是否存在网址。

```go
func HasURL(s string) bool
```

**示例：**
```go
if sensitive.HasURL("访问 https://example.com") {
    fmt.Println("包含网址")
}
```

### MaskURL

将字符串中的网址替换为 `***`。

```go
func MaskURL(s string) string
```

**示例：**
```go
text := "访问 https://example.com 或 http://test.org"
result := sensitive.MaskURL(text)
// 输出: "访问 *** 或 ***"
```

**匹配规则：**
- 支持 HTTP/HTTPS 协议
- 匹配完整 URL，包括路径和参数

## 数字检测与屏蔽

### HasDigit

判断字符串中是否包含指定个数的数字（大于等于该数字）。

```go
func HasDigit(s string, count int) bool
```

**示例：**
```go
// 判断是否包含至少 6 个数字（如手机号）
if sensitive.HasDigit("我的手机号是13812345678", 6) {
    fmt.Println("包含足够多的数字")
}

// 判断是否包含至少 11 个数字（完整手机号）
if sensitive.HasDigit("13812345678", 11) {
    fmt.Println("包含完整手机号")
}
```

### MaskDigit

将字符串中的所有数字替换为 `*`。

```go
func MaskDigit(s string) string
```

**示例：**
```go
text := "我的手机号是13812345678，QQ是123456789"
result := sensitive.MaskDigit(text)
// 输出: "我的手机号是***********，QQ是*********"
```

## 微信号检测与屏蔽

### HasWechatID

判断字符串中是否存在微信号。

```go
func HasWechatID(s string) bool
```

**示例：**
```go
if sensitive.HasWechatID("加我微信: my_wechat_123") {
    fmt.Println("包含微信号")
}
```

**匹配规则：**
- 以字母开头
- 长度 6-20 个字符
- 可包含字母、数字、下划线、短横线

### MaskWechatID

将字符串中的微信号替换为 `***`。

```go
func MaskWechatID(s string) string
```

**示例：**
```go
text := "加我微信: my_wechat_123 或 my-wechat"
result := sensitive.MaskWechatID(text)
// 输出: "加我微信: *** 或 ***"
```

## 组合使用

可以组合多个工具函数，实现多层防护：

```go
func sanitizeText(text string) string {
    // 1. 屏蔽邮箱
    text = sensitive.MaskEmail(text)
    
    // 2. 屏蔽 URL
    text = sensitive.MaskURL(text)
    
    // 3. 屏蔽微信号
    text = sensitive.MaskWechatID(text)
    
    // 4. 屏蔽长数字（如手机号）
    if sensitive.HasDigit(text, 11) {
        text = sensitive.MaskDigit(text)
    }
    
    return text
}
```

## 完整示例

```go
package main

import (
    "fmt"
    sensitive "github.com/LuYongwang/go-sensitive-word"
)

func main() {
    text := "联系邮箱: user@example.com, 访问 https://example.com, 微信: my_wechat, 手机: 13812345678"

    // 检测
    fmt.Printf("包含邮箱: %v\n", sensitive.HasEmail(text))
    fmt.Printf("包含URL: %v\n", sensitive.HasURL(text))
    fmt.Printf("包含微信号: %v\n", sensitive.HasWechatID(text))
    fmt.Printf("包含手机号: %v\n", sensitive.HasDigit(text, 11))

    // 屏蔽
    result := sensitive.MaskEmail(text)
    result = sensitive.MaskURL(result)
    result = sensitive.MaskWechatID(result)
    if sensitive.HasDigit(result, 11) {
        result = sensitive.MaskDigit(result)
    }

    fmt.Printf("屏蔽后: %s\n", result)
    // 输出: "联系邮箱: ***, 访问 ***, 微信: ***, 手机: ***********"
}
```

## 使用场景

### 1. 内容审核

在内容审核流程中，检测并屏蔽敏感信息：

```go
func reviewContent(content string) bool {
    // 检测敏感信息
    hasEmail := sensitive.HasEmail(content)
    hasURL := sensitive.HasURL(content)
    hasWechat := sensitive.HasWechatID(content)
    
    if hasEmail || hasURL || hasWechat {
        // 记录日志或告警
        log.Printf("内容包含敏感信息: email=%v, url=%v, wechat=%v", 
            hasEmail, hasURL, hasWechat)
        return false
    }
    
    return true
}
```

### 2. 内容展示

在展示用户内容前，自动屏蔽敏感信息：

```go
func displayContent(content string) string {
    // 自动屏蔽敏感信息
    content = sensitive.MaskEmail(content)
    content = sensitive.MaskURL(content)
    content = sensitive.MaskWechatID(content)
    return content
}
```

### 3. 数据脱敏

在导出数据或日志记录时，脱敏处理：

```go
func exportData(userData map[string]string) map[string]string {
    result := make(map[string]string)
    for k, v := range userData {
        // 根据字段类型选择不同的屏蔽策略
        switch k {
        case "email":
            result[k] = sensitive.MaskEmail(v)
        case "phone":
            if sensitive.HasDigit(v, 11) {
                result[k] = sensitive.MaskDigit(v)
            } else {
                result[k] = v
            }
        default:
            result[k] = v
        }
    }
    return result
}
```

## 注意事项

1. **性能考虑**
   - 工具函数使用正则表达式匹配，对长文本可能有性能影响
   - 建议在必要时使用，避免对所有文本都进行屏蔽处理

2. **匹配准确性**
   - 正则表达式可能无法覆盖所有边界情况
   - 如有特殊需求，建议自定义实现

3. **组合使用顺序**
   - 多个屏蔽函数组合使用时，注意处理顺序
   - URL 可能包含邮箱，建议先屏蔽 URL 再屏蔽邮箱

## 相关文档

- [API 参考文档](./api-reference.md)
- [最佳实践](./best-practices.md)
- [工具函数示例](../../examples/tools/main.go)

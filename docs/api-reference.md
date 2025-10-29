# API 参考文档

本文档提供 `go-sensitive-word` 的完整 API 参考。

## 初始化

### NewFilter

创建敏感词过滤器实例。

```go
func NewFilter(storeOpt StoreOption, filterOpt FilterOption) (*Manager, error)
```

**参数：**
- `storeOpt`: 存储配置（目前仅支持 `StoreMemory`）
- `filterOpt`: 过滤器配置（支持 `FilterDfa` 或 `FilterAC`）

**返回值：**
- `*Manager`: 管理器实例
- `error`: 错误信息

**示例：**
```go
filter, err := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
```

**相关示例：**
- [基础功能演示](../../examples/basic/main.go)

## 文本检测功能

### IsSensitive

判断文本中是否存在敏感词。

```go
func (m *Manager) IsSensitive(text string) bool
```

**参数：**
- `text`: 待检测的文本

**返回值：**
- `bool`: 存在敏感词返回 `true`，否则返回 `false`

**示例：**
```go
if filter.IsSensitive("包含敏感词的文本") {
    fmt.Println("检测到敏感词")
}
```

### FindOne

查找文本中的第一个敏感词。

```go
func (m *Manager) FindOne(text string) string
```

**参数：**
- `text`: 待检测的文本

**返回值：**
- `string`: 找到的第一个敏感词，未找到返回空字符串

**示例：**
```go
word := filter.FindOne("包含敏感词的文本")
if word != "" {
    fmt.Printf("找到敏感词: %s\n", word)
}
```

### FindAll

查找文本中所有敏感词（去重）。

```go
func (m *Manager) FindAll(text string) []string
```

**参数：**
- `text`: 待检测的文本

**返回值：**
- `[]string`: 所有敏感词的列表（已去重）

**示例：**
```go
words := filter.FindAll("包含多个敏感词的文本")
fmt.Printf("找到 %d 个敏感词: %v\n", len(words), words)
```

### FindAllCount

查找所有敏感词及其出现次数。

```go
func (m *Manager) FindAllCount(text string) map[string]int
```

**参数：**
- `text`: 待检测的文本

**返回值：**
- `map[string]int`: 敏感词及其出现次数的映射

**示例：**
```go
countMap := filter.FindAllCount("包含重复敏感词的文本")
for word, count := range countMap {
    fmt.Printf("%s: %d次\n", word, count)
}
```

### Replace

替换所有敏感词为指定字符。

```go
func (m *Manager) Replace(text string, replaceChar rune) string
```

**参数：**
- `text`: 待处理的文本
- `replaceChar`: 替换字符（如 `'*'`）

**返回值：**
- `string`: 替换后的文本

**示例：**
```go
result := filter.Replace("包含敏感词的文本", '*')
// 输出: "包含***的文本"
```

### Remove

从文本中删除所有敏感词。

```go
func (m *Manager) Remove(text string) string
```

**参数：**
- `text`: 待处理的文本

**返回值：**
- `string`: 删除敏感词后的文本

**示例：**
```go
result := filter.Remove("包含敏感词的文本")
// 输出: "包含的文本"
```

**相关示例：**
- [基础功能演示](../../examples/basic/main.go)

## 词库管理功能

### AddWord

动态添加敏感词（支持多个）。

```go
func (m *Manager) AddWord(words ...string) error
```

**参数：**
- `words`: 可变参数，要添加的敏感词列表

**返回值：**
- `error`: 错误信息

**注意：** 词的添加是异步处理的，建议添加后等待 100-200ms 再检测。

**示例：**
```go
err := filter.AddWord("敏感词1", "敏感词2", "敏感词3")
if err != nil {
    log.Fatal(err)
}
time.Sleep(100 * time.Millisecond)
```

### AddWords

批量添加敏感词。

```go
func (m *Manager) AddWords(words []string) error
```

**参数：**
- `words`: 要添加的敏感词列表

**返回值：**
- `error`: 错误信息

**示例：**
```go
words := []string{"词1", "词2", "词3"}
err := filter.AddWords(words)
```

### DelWord

动态删除敏感词（支持多个）。

```go
func (m *Manager) DelWord(words ...string) error
```

**参数：**
- `words`: 可变参数，要删除的敏感词列表

**返回值：**
- `error`: 错误信息

**注意：** 词的删除是异步处理的，建议删除后等待 100-200ms 再检测。

**示例：**
```go
err := filter.DelWord("敏感词1", "敏感词2")
time.Sleep(100 * time.Millisecond)
```

### DelWords

批量删除敏感词。

```go
func (m *Manager) DelWords(words []string) error
```

**参数：**
- `words`: 要删除的敏感词列表

**返回值：**
- `error`: 错误信息

**示例：**
```go
words := []string{"词1", "词2"}
err := filter.DelWords(words)
```

### ReplaceWords

批量替换敏感词（先删旧词，再加新词）。

```go
func (m *Manager) ReplaceWords(oldWords []string, newWords []string) error
```

**参数：**
- `oldWords`: 要删除的旧词列表
- `newWords`: 要添加的新词列表

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := filter.ReplaceWords(
    []string{"旧词1", "旧词2"},
    []string{"新词1", "新词2"},
)
```

### Clear

清空词库。

```go
func (m *Manager) Clear() error
```

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := filter.Clear()
```

### GetStats

获取词库统计信息。

```go
func (m *Manager) GetStats() Stats
```

**返回值：**
- `Stats`: 统计信息结构
  - `TotalWords int64`: 总词数
  - `LastUpdate time.Time`: 最后更新时间

**示例：**
```go
stats := filter.GetStats()
fmt.Printf("总词数: %d, 最后更新: %s\n", stats.TotalWords, stats.LastUpdate)
```

### ExportToString

导出词库为字符串。

```go
func (m *Manager) ExportToString() (string, error)
```

**返回值：**
- `string`: 词库字符串（每行一个词）
- `error`: 错误信息

**示例：**
```go
content, err := filter.ExportToString()
if err == nil {
    fmt.Println(content)
}
```

### ExportToFile

导出词库到文件。

```go
func (m *Manager) ExportToFile(filePath string) error
```

**参数：**
- `filePath`: 文件路径

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := filter.ExportToFile("/path/to/export.txt")
```

### MergeFromManager

合并另一个 Manager 的词库。

```go
func (m *Manager) MergeFromManager(other *Manager) error
```

**参数：**
- `other`: 另一个 Manager 实例

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := filter1.MergeFromManager(filter2)
```

### RefreshFromPath

从文件路径刷新词库（支持替换/追加模式）。

```go
func (m *Manager) RefreshFromPath(filePath string, replace bool) error
```

**参数：**
- `filePath`: 文件路径
- `replace`: `true` 表示替换模式（清空后重新加载），`false` 表示追加模式

**返回值：**
- `error`: 错误信息

**示例：**
```go
// 替换模式
err := filter.RefreshFromPath("/path/to/words.txt", true)

// 追加模式
err := filter.RefreshFromPath("/path/to/words.txt", false)
```

**相关示例：**
- [动态维护示例](../../examples/dynamic/main.go)

## 词库加载功能

### LoadDictEmbed

加载内置嵌入的词库（编译时嵌入）。

```go
func (m *Manager) LoadDictEmbed(dicts ...string) error
```

**参数：**
- `dicts`: 可变参数，内置词库变量（如 `DictPolitical`、`DictViolence` 等）

**返回值：**
- `error`: 错误信息

**内置词库变量：**
- `DictGFWAdditional`
- `DictOther`
- `DictReactionary`
- `DictAdvertisement`
- `DictPolitical`
- `DictViolence`
- `DictPeopleLife`
- `DictGunExplosion`
- `DictNeteaseFE`
- `DictSexual`
- `DictPornography`
- `DictAdditional`
- `DictCorruption`
- `DictTemporaryTencent`
- `DictIllegalURL`

**示例：**
```go
err := filter.LoadDictEmbed(
    sensitive.DictPolitical,
    sensitive.DictViolence,
)
```

### LoadDictPath

从文件路径加载词库（支持多个文件）。

```go
func (m *Manager) LoadDictPath(filePaths ...string) error
```

**参数：**
- `filePaths`: 可变参数，文件路径列表

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := filter.LoadDictPath(
    "/path/to/words1.txt",
    "/path/to/words2.txt",
)
```

### LoadDictCallback

通过回调函数加载词库（支持数据库、Redis等）。

```go
func (m *Manager) LoadDictCallback(fn func() ([]string, error), source string) error
```

**参数：**
- `fn`: 回调函数，返回敏感词列表和错误
- `source`: 数据源标识（用于日志）

**返回值：**
- `error`: 错误信息

**示例：**
```go
// 从数据库加载
err := filter.LoadDictCallback(func() ([]string, error) {
    return db.QueryWords(), nil
}, "database")

// 从 Redis 加载
err := filter.LoadDictCallback(func() ([]string, error) {
    return redis.SMembers(ctx, "sensitive:words"), nil
}, "redis")
```

**相关示例：**
- [文件加载示例](../../examples/file-load/main.go)
- [回调加载示例](../../examples/callback/main.go)

## 资源管理

### Close

关闭过滤器，停止所有 goroutine。

```go
func (m *Manager) Close() error
```

**返回值：**
- `error`: 错误信息

**示例：**
```go
defer filter.Close()
```

### Shutdown

优雅关闭（带超时，等待异步处理完成）。

```go
func (m *Manager) Shutdown(ctx context.Context) error
```

**参数：**
- `ctx`: 上下文（用于超时控制）

**返回值：**
- `error`: 错误信息

**示例：**
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
defer filter.Shutdown(ctx)
```

**相关示例：**
- [资源管理示例](../../examples/lifecycle/main.go)

## 工具函数

### HasEmail / MaskEmail

邮箱检测和屏蔽。

```go
func HasEmail(s string) bool
func MaskEmail(s string) string
```

### HasURL / MaskURL

URL 检测和屏蔽。

```go
func HasURL(s string) bool
func MaskURL(s string) string
```

### HasDigit / MaskDigit

数字检测和屏蔽。

```go
func HasDigit(s string, count int) bool
func MaskDigit(s string) string
```

### HasWechatID / MaskWechatID

微信号检测和屏蔽。

```go
func HasWechatID(s string) bool
func MaskWechatID(s string) string
```

**相关示例：**
- [工具函数示例](../../examples/tools/main.go)

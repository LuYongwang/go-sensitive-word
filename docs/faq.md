# 常见问题（FAQ）

本文档收集了 `go-sensitive-word` 使用过程中的常见问题和解答。

## 基础使用

### Q1: 为什么添加词后立即检测可能检测不到？

**A:** DFA/AC 算法通过 channel 异步处理词的添加/删除，需要短暂延迟才能生效。

**解决方案：**
```go
filter.AddWord("敏感词")
time.Sleep(100 * time.Millisecond) // 等待异步处理完成
result := filter.IsSensitive("包含敏感词的文本")
```

**原因说明：**
- DFA/AC 算法为了性能，将词库更新放在后台 goroutine 异步处理
- 更新操作通过 channel 发送到处理队列
- 需要等待 goroutine 处理完队列中的任务

**推荐延迟时间：**
- DFA: 100ms
- AC: 100-200ms（窗口合并可能需要更长时间）

### Q2: DFA 和 AC 算法如何选择？

**A:** 根据场景选择：

| 场景 | 推荐算法 |
|------|---------|
| **生产环境** | AC（推荐） |
| **大词库**（> 1000 词） | AC |
| **长文本检测** | AC |
| **小词库**（< 1000 词） | DFA 或 AC |
| **简单场景/调试** | DFA |

**详细对比：**
- AC：性能更好，适合生产环境
- DFA：实现简单，小词库性能可接受

详见：[算法选择指南](./algorithm-guide.md)

### Q3: 如何自定义归一化配置？

**A:** 当前版本默认归一化已开启"忽略大小写"和"全角转半角"。

如需更严格配置，可参考 `normalize.go` 中的 `StrictNormalizer()`：

```go
// 严格归一化配置（示例）
config := sensitive.StrictNormalizer()
// 注意：当前版本归一化配置在内部处理，用户无需显式配置
```

详见：[归一化功能详解](./normalization.md)

### Q4: 词库更新是否需要重启服务？

**A:** **不需要！** 项目支持动态维护，可以在运行时通过 `AddWord()`、`DelWord()`、`ReplaceWords()` 等方法实时更新词库，无需重启。

**示例：**
```go
// 动态添加
filter.AddWord("新敏感词")

// 动态删除
filter.DelWord("旧敏感词")

// 批量替换
filter.ReplaceWords([]string{"旧词"}, []string{"新词"})
```

详见：[词库管理详解](./word-management.md)

### Q5: 如何从数据库/Redis 加载词库？

**A:** 使用 `LoadDictCallback()` 方法，传入回调函数：

```go
// 从数据库加载
filter.LoadDictCallback(func() ([]string, error) {
    return db.QueryWords(), nil
}, "database")

// 从 Redis 加载
filter.LoadDictCallback(func() ([]string, error) {
    return redis.SMembers(ctx, "sensitive:words"), nil
}, "redis")
```

详见：[词库加载详解](./word-loading.md) | [回调加载示例](../../examples/callback/main.go)

## 性能优化

### Q6: 大词库场景如何优化性能？

**A:** 建议：

1. **使用 AC 算法**（推荐）：
   ```go
   filter, _ := sensitive.NewFilter(
       sensitive.StoreOption{Type: sensitive.StoreMemory},
       sensitive.FilterOption{Type: sensitive.FilterAC},
   )
   ```

2. **批量操作**：
   ```go
   // 推荐：批量添加
   filter.AddWords([]string{"词1", "词2", "词3"})
   
   // 不推荐：多次单次添加
   filter.AddWord("词1")
   filter.AddWord("词2")
   filter.AddWord("词3")
   ```

3. **避免频繁更新**：
   - 将多次小更新合并为一次大更新
   - AC 算法已内置窗口合并机制

详见：[最佳实践](./best-practices.md)

### Q7: 高并发场景如何优化？

**A:** 建议：

1. **使用 AC 算法**：支持高并发，性能线性扩展
2. **预加载词库**：启动时加载完整词库，避免运行时加载
3. **连接池**：如果使用数据库/Redis，确保配置合适的连接池

### Q8: 内存占用如何优化？

**A:** 建议：

1. **选择合适算法**：
   - DFA：内存占用较低
   - AC：内存占用中等，但性能更好

2. **按需加载词库**：
   ```go
   // 只加载必要的词库
   filter.LoadDictEmbed(
       sensitive.DictPolitical,  // 只加载需要的
       sensitive.DictViolence,
   )
   ```

## 功能相关

### Q9: 如何处理重叠的敏感词？

**A:** `Replace()` 和 `Remove()` 方法已自动处理重叠情况：

```go
// 自动处理重叠，不会重复替换
result := filter.Replace("敏感词1敏感词2", '*')
```

内部使用区间映射，确保重叠部分只处理一次。

### Q10: 如何检测特定格式的敏感信息（如手机号、身份证）？

**A:** 使用工具函数：

```go
// 检测手机号（11位数字）
if sensitive.HasDigit(text, 11) {
    masked := sensitive.MaskDigit(text)
}

// 检测邮箱
if sensitive.HasEmail(text) {
    masked := sensitive.MaskEmail(text)
}
```

详见：[工具函数详解](./tools.md)

### Q11: 如何实现自定义的屏蔽规则？

**A:** 可以结合工具函数和敏感词检测：

```go
func customMask(text string) string {
    // 1. 敏感词替换
    text = filter.Replace(text, '*')
    
    // 2. 邮箱屏蔽
    text = sensitive.MaskEmail(text)
    
    // 3. URL 屏蔽
    text = sensitive.MaskURL(text)
    
    // 4. 自定义规则
    // ... 自定义逻辑
    
    return text
}
```

## 错误排查

### Q12: 出现 "store is nil" 错误

**A:** 检查初始化代码：

```go
// 确保正确初始化
filter, err := sensitive.NewFilter(
    sensitive.StoreOption{Type: sensitive.StoreMemory},
    sensitive.FilterOption{Type: sensitive.FilterAC},
)
if err != nil {
    log.Fatal(err)
}
```

### Q13: 检测结果不准确

**A:** 检查以下几点：

1. **归一化一致性**：确保词库词和检测文本使用相同的归一化策略
2. **异步延迟**：添加/删除词后等待 100-200ms
3. **词库加载**：确认词库已正确加载

```go
// 检查词库统计
stats := filter.GetStats()
fmt.Printf("词库总词数: %d\n", stats.TotalWords)
```

### Q14: 性能突然下降

**A:** 可能原因：

1. **词库过大**：使用 AC 算法
2. **频繁更新**：合并多次更新为批量更新
3. **并发竞争**：检查是否有其他性能瓶颈

## 部署相关

### Q15: 如何优雅关闭过滤器？

**A:** 使用 `Shutdown()` 方法：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
defer filter.Shutdown(ctx)
```

详见：[资源管理详解](./lifecycle.md)

### Q16: Docker 容器中如何使用？

**A:** 正常使用即可，注意：

1. **文件路径**：如果使用 `LoadDictPath`，确保文件在容器中
2. **优雅关闭**：容器退出时正确关闭过滤器

```dockerfile
# Dockerfile 示例
FROM golang:1.21
WORKDIR /app
COPY . .
RUN go build -o app .
CMD ["./app"]
```

## 其他问题

### Q17: 支持其他语言吗？

**A:** 当前版本主要支持中文和英文，通过归一化可以处理：

- 大小写混淆（英文）
- 全角/半角（中文、英文）
- 繁简体（中文）

如需支持其他语言，可以扩展归一化配置。

### Q18: 如何贡献代码？

**A:** 欢迎提交 Issue 和 Pull Request。

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 发起 Pull Request

### Q19: 是否有性能测试数据？

**A:** 可以参考：

- [算法选择指南](./algorithm-guide.md) 中的性能对比
- 运行 `go test -bench` 查看基准测试结果

### Q20: 是否线程安全？

**A:** **是的**，所有操作都是并发安全的：

```go
// 可以安全地在多个 goroutine 中使用
go filter.AddWord("词1")
go filter.AddWord("词2")
go filter.IsSensitive("文本")
```

## 获取帮助

如果以上问题未能解决您的问题：

1. 查看 [完整文档](../README.md)
2. 查看 [示例代码](../../examples/)
3. 提交 [Issue](https://github.com/your-repo/issues)

## 相关文档

- [快速开始指南](./getting-started.md)
- [API 参考文档](./api-reference.md)
- [最佳实践](./best-practices.md)

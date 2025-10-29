package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示 AC 自动机算法的使用
// AC 算法适合长文本、大词库场景，性能优于 DFA
func main() {
	fmt.Println("=== AC 自动机算法示例 ===\n")

	// 创建使用 AC 算法的过滤器
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC}, // 使用 AC 算法
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 加载词库
	err = filter.LoadDictEmbed(
		sensitive.DictPolitical,
		sensitive.DictViolence,
	)
	if err != nil {
		log.Fatalf("加载词库失败: %v", err)
	}

	// 添加自定义敏感词
	err = filter.AddWords([]string{"违禁词A", "违禁词B", "违禁词C"})
	if err != nil {
		log.Fatalf("添加词失败: %v", err)
	}

	// 等待异步处理完成（AC 有窗口合并机制，可能需要等待）
	time.Sleep(200 * time.Millisecond)

	// 测试文本
	testText := "这是一个测试文本，包含违禁词A和违禁词B，还有违禁词C"
	fmt.Printf("测试文本: %s\n\n", testText)

	// 基础检测功能
	fmt.Println("=== 基础检测功能 ===")
	fmt.Printf("是否包含敏感词: %v\n", filter.IsSensitive(testText))
	fmt.Printf("找到第一个敏感词: %s\n", filter.FindOne(testText))
	fmt.Printf("找到所有敏感词: %v\n", filter.FindAll(testText))
	fmt.Printf("敏感词及次数: %v\n\n", filter.FindAllCount(testText))

	// 文本处理功能
	fmt.Println("=== 文本处理功能 ===")
	fmt.Printf("替换敏感词: %s\n", filter.Replace(testText, '*'))
	fmt.Printf("移除敏感词: %s\n\n", filter.Remove(testText))

	// 统计信息
	stats := filter.GetStats()
	fmt.Println("=== 统计信息 ===")
	fmt.Printf("总词数: %d\n", stats.TotalWords)
	fmt.Printf("最后更新: %s\n", stats.LastUpdate.Format(time.RFC3339))
	fmt.Printf("来源: %v\n", stats.Source)

	fmt.Println("\n=== AC 算法示例完成 ===")
}

package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示词库来源追踪功能
func main() {
	fmt.Println("=== 词库来源追踪示例 ===")
	fmt.Println()

	// 初始化过滤器
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer filter.Close()

	// 1. 从不同来源加载词库
	fmt.Println("【1】从不同来源加载词库")

	// 加载政治类型词库（使用带来源的方法）
	err = filter.LoadDictEmbedWithSource(sensitive.DictPolitical, "political")
	if err != nil {
		log.Fatalf("加载政治词库失败: %v", err)
	}
	fmt.Println("  ✅ 加载政治类型词库")

	// 加载民生词库
	_ = filter.LoadDictEmbedWithSource(sensitive.DictPeopleLife, "peopleLife")
	_ = filter.LoadDictEmbedWithSource(sensitive.DictAdvertisement, "advertisement")

	// 手动添加一些词并指定来源
	customWords := []string{"违禁词A", "违禁词B", "自定义敏感词"}
	err = filter.AddWordsWithSource(customWords, "custom")
	if err != nil {
		log.Fatalf("添加自定义词失败: %v", err)
	}
	fmt.Println("  ✅ 添加自定义敏感词（来源：custom）")

	// 添加更多词并指定另一个来源
	businessWords := []string{"违禁词A", "违禁词C", "业务敏感词"}
	err = filter.AddWordsWithSource(businessWords, "business")
	if err != nil {
		log.Fatalf("添加业务词失败: %v", err)
	}
	fmt.Println("  ✅ 添加业务敏感词（来源：business）")

	// 等待异步处理完成
	time.Sleep(200 * time.Millisecond)

	fmt.Println()

	// 2. 测试文本匹配
	fmt.Println("【2】文本匹配测试")
	testText := "这段文本包含违禁词A和违禁词B"
	fmt.Printf("  测试文本: %s\n\n", testText)

	// 3. 查找所有敏感词及其来源
	fmt.Println("【3】查找所有敏感词及其来源")
	results := filter.FindAllWithSource(testText)
	fmt.Printf("  找到 %d 个敏感词：\n", len(results))
	for i, result := range results {
		fmt.Printf("    %d. 词: \"%s\"\n", i+1, result.Word)
		fmt.Printf("       来源: %v\n", result.Source)
	}

	fmt.Println()

	// 4. 查找敏感词并统计出现次数和来源
	fmt.Println("【4】查找敏感词并统计出现次数和来源")
	testText2 := "违禁词A出现了违禁词A两次，还有违禁词B"
	fmt.Printf("  测试文本: %s\n", testText2)

	// 注意：FindAllCountWithSource 返回的是每个词及其来源（不包含count）
	// 实际的计数通过 FindAllCount 获取
	countMap := filter.FindAllCount(testText2)
	countWithSource := filter.FindAllCountWithSource(testText2)

	fmt.Printf("  找到 %d 种敏感词：\n", len(countMap))
	for word, result := range countWithSource {
		count := countMap[word]
		fmt.Printf("    - 词: \"%s\" 出现 %d 次\n", word, count)
		fmt.Printf("      来源: %v\n", result.Source)
	}

	fmt.Println()

	// 5. 查询单个词的来源
	fmt.Println("【5】查询单个词的来源")
	testWords := []string{"违禁词A", "违禁词B", "违禁词C", "不存在的词"}
	for _, word := range testWords {
		sources := filter.GetWordSources(word)
		if len(sources) > 0 {
			fmt.Printf("  \"%s\" 的来源: %v\n", word, sources)
		} else {
			fmt.Printf("  \"%s\": 未找到\n", word)
		}
	}

	fmt.Println()

	// 6. 获取所有词的来源映射
	fmt.Println("【6】获取所有词的来源映射")
	allSources := filter.GetAllWordSources()
	fmt.Printf("  词库中共有 %d 个词：\n", len(allSources))

	// 显示前10个词的来源
	count := 0
	for word, sources := range allSources {
		if count >= 10 {
			fmt.Printf("    ... 还有 %d 个词\n", len(allSources)-10)
			break
		}
		fmt.Printf("    \"%s\": %v\n", word, sources)
		count++
	}

	// 7. 获取统计信息
	fmt.Println()
	fmt.Println("【7】词库统计信息")
	stats := filter.GetStats()
	fmt.Printf("  总词数: %d\n", stats.TotalWords)
	fmt.Printf("  来源数: %d\n", len(stats.Source))
	fmt.Printf("  最后更新: %s\n", stats.LastUpdate.Format("2006-01-02 15:04:05"))

	fmt.Println("\n=== 示例完成 ===")
}

package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 综合示例：演示项目的所有核心功能
func main() {
	fmt.Println("=== go-sensitive-word 综合示例 ===\n")

	// ========== 1. 初始化 ==========
	fmt.Println("【1】初始化过滤器")
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC}, // 使用 AC 算法（推荐生产环境）
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer filter.Close() // 确保资源释放

	// ========== 2. 加载词库 ==========
	fmt.Println("\n【2】加载内置词库")
	err = filter.LoadDictEmbed(
		sensitive.DictPolitical,
		sensitive.DictViolence,
	)
	if err != nil {
		log.Fatalf("加载词库失败: %v", err)
	}
	fmt.Println("  ✅ 词库加载完成")

	// ========== 3. 动态添加词 ==========
	fmt.Println("\n【3】动态添加敏感词")
	err = filter.AddWords([]string{"自定义词1", "自定义词2", "自定义词3"})
	if err != nil {
		log.Fatalf("添加词失败: %v", err)
	}
	fmt.Println("  ✅ 已添加 3 个敏感词")

	time.Sleep(200 * time.Millisecond) // 等待异步处理

	// ========== 4. 查看统计信息 ==========
	fmt.Println("\n【4】词库统计信息")
	stats := filter.GetStats()
	fmt.Printf("  总词数: %d\n", stats.TotalWords)
	fmt.Printf("  最后更新: %s\n", stats.LastUpdate.Format(time.RFC3339))
	fmt.Printf("  更新次数: %d\n", stats.UpdateCount)
	fmt.Printf("  来源数: %d\n", len(stats.Source))

	// ========== 5. 文本检测功能 ==========
	fmt.Println("\n【5】文本检测功能")
	testText := "这是一个包含自定义词1和自定义词2的测试文本，还包含自定义词3"
	fmt.Printf("  测试文本: %s\n\n", testText)

	fmt.Println("  a) 判断是否包含敏感词:")
	fmt.Printf("     IsSensitive(): %v\n", filter.IsSensitive(testText))

	fmt.Println("\n  b) 查找第一个敏感词:")
	fmt.Printf("     FindOne(): %s\n", filter.FindOne(testText))

	fmt.Println("\n  c) 查找所有敏感词（去重）:")
	allWords := filter.FindAll(testText)
	fmt.Printf("     FindAll(): %v\n", allWords)

	fmt.Println("\n  d) 查找所有敏感词及出现次数:")
	countMap := filter.FindAllCount(testText)
	fmt.Printf("     FindAllCount(): %v\n", countMap)

	// ========== 6. 文本处理功能 ==========
	fmt.Println("\n【6】文本处理功能")
	fmt.Println("  a) 替换敏感词（用 * 替换）:")
	replaced := filter.Replace(testText, '*')
	fmt.Printf("     Replace(): %s\n", replaced)

	fmt.Println("\n  b) 移除敏感词:")
	removed := filter.Remove(testText)
	fmt.Printf("     Remove(): %s\n", removed)

	// ========== 7. 动态维护 ==========
	fmt.Println("\n【7】动态维护词库")

	fmt.Println("  a) 批量删除敏感词:")
	err = filter.DelWords([]string{"自定义词2"})
	if err != nil {
		log.Printf("删除失败: %v", err)
	} else {
		fmt.Println("    ✅ 已删除 '自定义词2'")
	}

	fmt.Println("\n  b) 批量替换敏感词:")
	err = filter.ReplaceWords(
		[]string{"自定义词1"},
		[]string{"新词A", "新词B"},
	)
	if err != nil {
		log.Printf("替换失败: %v", err)
	} else {
		fmt.Println("    ✅ 已替换 '自定义词1' -> '新词A', '新词B'")
	}

	time.Sleep(200 * time.Millisecond)

	// 验证动态维护效果
	testText2 := "包含新词A和新词B，但不包含自定义词2"
	fmt.Printf("\n  测试文本: %s\n", testText2)
	fmt.Printf("  是否敏感: %v\n", filter.IsSensitive(testText2))
	fmt.Printf("  找到的敏感词: %v\n", filter.FindAll(testText2))

	// ========== 8. 词库导出 ==========
	fmt.Println("\n【8】词库导出")
	exported, err := filter.ExportToString()
	if err != nil {
		log.Printf("导出失败: %v", err)
	} else {
		preview := exported
		if len(preview) > 150 {
			preview = preview[:150] + "..."
		}
		fmt.Printf("  导出为字符串（预览）: %s\n", preview)
		fmt.Printf("  总长度: %d 字符\n", len(exported))
	}

	// ========== 9. 词库合并 ==========
	fmt.Println("\n【9】词库合并")
	filter2, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err == nil {
		filter2.AddWords([]string{"合并词1", "合并词2"})
		time.Sleep(100 * time.Millisecond)

		err = filter.MergeFromManager(filter2)
		if err != nil {
			log.Printf("合并失败: %v", err)
		} else {
			fmt.Println("  ✅ 已合并第二个词库")
			stats = filter.GetStats()
			fmt.Printf("  合并后总词数: %d\n", stats.TotalWords)
		}
		filter2.Close()
	}

	// ========== 10. 归一化功能 ==========
	fmt.Println("\n【10】归一化功能演示")
	filter3, _ := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	filter3.AddWord("http")
	time.Sleep(100 * time.Millisecond)

	normalizeTests := []string{
		"http", // 正常
		"HTTP", // 大写
		"Ｈttp", // 全角 H
		"ＨＴＴＰ", // 全角全部
	}

	fmt.Println("  归一化测试（忽略大小写、全角转半角）:")
	for _, text := range normalizeTests {
		result := filter3.IsSensitive(text)
		fmt.Printf("    '%s' -> %v\n", text, result)
	}

	filter3.Close()

	// ========== 11. 工具函数 ==========
	fmt.Println("\n【11】工具函数")
	mixedText := "联系邮箱: user@example.com, 访问 https://example.com, 微信号: my_wechat"
	fmt.Printf("  原文: %s\n", mixedText)

	result := sensitive.MaskEmail(mixedText)
	result = sensitive.MaskURL(result)
	result = sensitive.MaskWechatID(result)
	fmt.Printf("  综合屏蔽: %s\n", result)

	// ========== 总结 ==========
	fmt.Println("\n=== 功能总结 ===")
	fmt.Println(`
  已演示的功能：
  ✓ 初始化和配置（DFA/AC 算法选择）
  ✓ 词库加载（内置词库、文件、回调函数）
  ✓ 文本检测（IsSensitive, FindOne, FindAll, FindAllCount）
  ✓ 文本处理（Replace, Remove）
  ✓ 动态维护（AddWords, DelWords, ReplaceWords）
  ✓ 词库导出（ExportToString, ExportToFile）
  ✓ 词库合并（MergeFromManager）
  ✓ 归一化（忽略大小写、全角转半角等）
  ✓ 工具函数（邮箱、URL、微信号检测和屏蔽）
  ✓ 资源管理（Close, Shutdown）

  更多示例请查看 examples/ 目录下的其他示例。
  `)

	fmt.Println("=== 综合示例完成 ===")
}

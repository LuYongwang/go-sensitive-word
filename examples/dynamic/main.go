package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示动态维护词库功能
func main() {
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 1. 加载初始词库
	err = filter.LoadDictEmbed(
		sensitive.DictPolitical,
		sensitive.DictViolence,
	)
	if err != nil {
		log.Fatalf("加载词库失败: %v", err)
	}

	// 2. 查看统计信息
	stats := filter.GetStats()
	fmt.Printf("=== 初始词库统计 ===\n")
	fmt.Printf("总词数: %d\n", stats.TotalWords)
	fmt.Printf("最后更新: %s\n", stats.LastUpdate.Format(time.RFC3339))
	fmt.Printf("更新次数: %d\n", stats.UpdateCount)
	fmt.Printf("来源: %v\n\n", stats.Source)

	// 3. 批量添加敏感词
	newWords := []string{"违禁词1", "违禁词2", "违禁词3"}
	err = filter.AddWords(newWords)
	if err != nil {
		log.Fatalf("批量添加失败: %v", err)
	}
	fmt.Printf("✅ 已批量添加 %d 个词\n\n", len(newWords))

	// 4. 批量删除敏感词
	delWords := []string{"违禁词2"}
	err = filter.DelWords(delWords)
	if err != nil {
		log.Fatalf("批量删除失败: %v", err)
	}
	fmt.Printf("✅ 已批量删除 %d 个词\n\n", len(delWords))

	// 5. 批量替换（先删旧词，再加新词）
	oldWords := []string{"违禁词1"}
	newReplacementWords := []string{"新违禁词A", "新违禁词B"}
	err = filter.ReplaceWords(oldWords, newReplacementWords)
	if err != nil {
		log.Fatalf("批量替换失败: %v", err)
	}
	fmt.Printf("✅ 已替换: 删除 %d 个旧词，添加 %d 个新词\n\n", len(oldWords), len(newReplacementWords))

	// 6. 查看更新后的统计信息
	stats = filter.GetStats()
	fmt.Printf("=== 更新后词库统计 ===\n")
	fmt.Printf("总词数: %d\n", stats.TotalWords)
	fmt.Printf("最后更新: %s\n", stats.LastUpdate.Format(time.RFC3339))
	fmt.Printf("更新次数: %d\n", stats.UpdateCount)
	fmt.Println()

	// 7. 导出词库到字符串
	exported, err := filter.ExportToString()
	if err != nil {
		log.Fatalf("导出失败: %v", err)
	}
	previewLen := 100
	if len(exported) < previewLen {
		previewLen = len(exported)
	}
	fmt.Printf("=== 导出词库（前%d字符）===\n%s...\n\n", previewLen, exported[:previewLen])

	// 8. 导出词库到文件
	err = filter.ExportToFile("/tmp/sensitive_words_export.txt")
	if err != nil {
		log.Printf("⚠️  导出到文件失败: %v (可能是权限问题)\n", err)
	} else {
		fmt.Printf("✅ 词库已导出到 /tmp/sensitive_words_export.txt\n\n")
	}

	// 9. 测试过滤功能
	testText := "这是一个包含违禁词3和新违禁词A的测试文本"
	fmt.Printf("=== 测试文本 ===\n%s\n\n", testText)
	fmt.Printf("是否敏感: %v\n", filter.IsSensitive(testText))
	fmt.Printf("找到的敏感词: %v\n", filter.FindAll(testText))

	// 11. 词库合并示例（创建第二个 Manager）
	fmt.Printf("\n=== 词库合并示例 ===\n")
	filter2, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err == nil {
		if err = filter2.AddWords([]string{"合并词1", "合并词2"}); err != nil {
			log.Printf("⚠️  添加词失败: %v\n", err)
		} else if err = filter.MergeFromManager(filter2); err != nil {
			log.Printf("⚠️  合并失败: %v\n", err)
		} else {
			fmt.Printf("✅ 已合并第二个词库\n")
			stats = filter.GetStats()
			fmt.Printf("合并后总词数: %d\n", stats.TotalWords)
		}
	}

	fmt.Printf("\n=== 动态维护功能演示完成 ===\n")
}

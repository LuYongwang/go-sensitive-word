package main

import (
	"fmt"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示归一化配置的使用
// 归一化可以有效防御各种混淆攻击（大小写、全角、相似字符等）
func main() {
	fmt.Println("=== 归一化配置示例 ===\n")

	// 1. 默认归一化（忽略大小写、全角转半角）
	fmt.Println("--- 默认归一化配置 ---")
	filter1, _ := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	filter1.AddWord("test")
	time.Sleep(100 * time.Millisecond)

	testCases := []string{
		"TEST", // 大写
		"ｔｅｓｔ", // 全角
		"TeSt", // 混合大小写
		"ＴｅＳｔ", // 全角大小写混合
	}

	for _, text := range testCases {
		result := filter1.IsSensitive(text)
		fmt.Printf("  '%s' -> %v\n", text, result)
	}

	// 2. 严格归一化（启用所有防绕过选项）
	fmt.Println("\n--- 严格归一化配置（Strict）---")
	fmt.Println("  注意：当前版本默认归一化已足够，StrictNormalizer 可作为参考实现")
	fmt.Println("  在实际使用中，可以根据需要自定义 NormalizerConfig")

	// 演示自定义归一化配置
	fmt.Println("\n--- 归一化示例 ---")

	// 测试全角字符
	filter2, _ := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	filter2.AddWord("http")
	time.Sleep(100 * time.Millisecond)

	fullWidthTests := []string{
		"Ｈttp", // 全角 H
		"http", // 正常
		"ＨＴＴＰ", // 全角全部
	}

	for _, text := range fullWidthTests {
		result := filter2.IsSensitive(text)
		fmt.Printf("  '%s' -> %v\n", text, result)
	}

	// 演示数字归一化（如果需要）
	fmt.Println("\n--- 数字归一化示例（如需启用需修改配置）---")
	fmt.Println("  默认不开启数字归一化，如需启用可参考 normalize.go 中的 StrictNormalizer")
	fmt.Println("  支持的归一化：各种数字写法 → 阿拉伯数字（①→1, 二→2 等）")

	// 演示归一化文本函数
	fmt.Println("\n--- 归一化文本工具函数 ---")
	cfg := sensitive.DefaultNormalizer()
	original := "ＴｅＳｔ　Ｈｔｔｐ"
	normalized := sensitive.NormalizeWord(original, cfg)
	fmt.Printf("  原文: %s\n", original)
	fmt.Printf("  归一化后: %s\n", normalized)

	fmt.Println("\n=== 归一化配置示例完成 ===")
}

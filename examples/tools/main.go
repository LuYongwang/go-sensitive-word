package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示工具函数的使用
// 这些工具函数用于检测和屏蔽邮箱、网址、微信号等敏感信息
func main() {
	fmt.Println("=== 工具函数示例 ===\n")

	// 1. 邮箱检测和屏蔽
	fmt.Println("--- 邮箱检测和屏蔽 ---")
	emailText := "请联系我的邮箱：user@example.com 或 admin@test.org"
	fmt.Printf("  原文: %s\n", emailText)
	fmt.Printf("  包含邮箱: %v\n", sensitive.HasEmail(emailText))
	fmt.Printf("  屏蔽后: %s\n\n", sensitive.MaskEmail(emailText))

	// 2. URL 检测和屏蔽
	fmt.Println("--- URL 检测和屏蔽 ---")
	urlText := "访问 https://example.com 或 http://test.org/path"
	fmt.Printf("  原文: %s\n", urlText)
	fmt.Printf("  包含URL: %v\n", sensitive.HasURL(urlText))
	fmt.Printf("  屏蔽后: %s\n\n", sensitive.MaskURL(urlText))

	// 3. 数字检测和屏蔽
	fmt.Println("--- 数字检测和屏蔽 ---")
	digitText := "我的电话是13800138000，房间号是2024"
	fmt.Printf("  原文: %s\n", digitText)
	fmt.Printf("  包含至少5个数字: %v\n", sensitive.HasDigit(digitText, 5))
	fmt.Printf("  包含至少11个数字: %v\n", sensitive.HasDigit(digitText, 11))
	fmt.Printf("  屏蔽所有数字: %s\n\n", sensitive.MaskDigit(digitText))

	// 4. 微信号检测和屏蔽
	fmt.Println("--- 微信号检测和屏蔽 ---")
	wechatText := "加我微信：test_user123 或 MyWechat_ID"
	fmt.Printf("  原文: %s\n", wechatText)
	fmt.Printf("  包含微信号: %v\n", sensitive.HasWechatID(wechatText))
	fmt.Printf("  屏蔽后: %s\n\n", sensitive.MaskWechatID(wechatText))

	// 5. 组合使用示例
	fmt.Println("--- 组合使用示例 ---")
	combinedText := `联系我：
邮箱: user@example.com
微信: my_wechat_id
电话: 13800138000
网站: https://example.com`

	fmt.Printf("  原文:\n%s\n\n", combinedText)

	// 依次屏蔽各种敏感信息
	result := sensitive.MaskEmail(combinedText)
	result = sensitive.MaskURL(result)
	result = sensitive.MaskWechatID(result)
	result = sensitive.MaskDigit(result)

	fmt.Printf("  综合屏蔽后:\n%s\n", result)

	// 6. 结合敏感词过滤
	fmt.Println("\n--- 结合敏感词过滤 ---")
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	filter.AddWords([]string{"违禁词", "敏感内容"})
	time.Sleep(100 * time.Millisecond)

	compositeText := "这是一个包含违禁词的文本，邮箱是 user@example.com，网站是 https://test.com"
	fmt.Printf("  原文: %s\n", compositeText)

	// 先过滤敏感词
	filtered1 := filter.Replace(compositeText, '*')
	// 再屏蔽邮箱和URL
	filtered2 := sensitive.MaskEmail(filtered1)
	finalResult := sensitive.MaskURL(filtered2)

	fmt.Printf("  最终结果: %s\n", finalResult)

	fmt.Println("\n=== 工具函数示例完成 ===")
}

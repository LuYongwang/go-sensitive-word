package main

import (
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示多实例：不同业务场景使用不同的过滤器实例与词库
func main() {
	// 实例 A：使用 DFA，词库为 [1,2,3,4]
	filterA, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化实例A失败: %v", err)
	}
	defer filterA.Close()

	// 实例 B：使用 AC，词库为 [4,5,6,7,8]
	filterB, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC},
	)
	if err != nil {
		log.Fatalf("初始化实例B失败: %v", err)
	}
	defer filterB.Close()

	// 加载各自词库（互不影响）
	if err := filterA.AddWords([]string{"1", "2", "3", "4"}); err != nil {
		log.Fatalf("实例A加载词库失败: %v", err)
	}
	if err := filterB.AddWords([]string{"4", "5", "6", "7", "8"}); err != nil {
		log.Fatalf("实例B加载词库失败: %v", err)
	}

	// 等待异步处理完成
	time.Sleep(100 * time.Millisecond)

	// 业务文本
	textA := "这段文本包含3"
	textB := "这段文本包含6"

	// 各自场景检测
	fmt.Println("A 场景（应为 true）：", filterA.IsSensitive(textA))
	fmt.Println("B 场景（应为 true）：", filterB.IsSensitive(textB))

	// 交叉验证
	fmt.Println("A 检测B文本（应为 false）：", filterA.IsSensitive(textB))
	fmt.Println("B 检测A文本（应为 true，因为含4为交集示例可选）：", filterB.IsSensitive("包含4的文本"))
}



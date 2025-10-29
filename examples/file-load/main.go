package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示从文件加载词库
func main() {
	fmt.Println("=== 从文件加载词库示例 ===\n")

	// 创建临时测试文件
	tmpDir := os.TempDir()
	testFile := filepath.Join(tmpDir, "test_sensitive_words.txt")

	// 准备测试词库内容
	testWords := `敏感词1
敏感词2
敏感词3
测试词4
`

	// 写入测试文件
	err := os.WriteFile(testFile, []byte(testWords), 0644)
	if err != nil {
		log.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(testFile) // 清理临时文件

	fmt.Printf("创建测试文件: %s\n\n", testFile)

	// 方式 1: 追加模式（追加到现有词库）
	fmt.Println("--- 方式 1: 追加模式加载 ---")
	filter1, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 先添加一些词
	filter1.AddWords([]string{"已有词1", "已有词2"})
	time.Sleep(100 * time.Millisecond)

	// 从文件追加加载
	err = filter1.LoadDictPath(testFile)
	if err != nil {
		log.Fatalf("从文件加载失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	stats1 := filter1.GetStats()
	fmt.Printf("  追加后总词数: %d\n", stats1.TotalWords)
	fmt.Printf("  来源: %v\n\n", stats1.Source)

	// 方式 2: 替换模式（清空后重新加载）
	fmt.Println("--- 方式 2: 替换模式加载 ---")
	filter2, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 先添加一些词
	filter2.AddWords([]string{"旧词1", "旧词2"})
	time.Sleep(100 * time.Millisecond)

	// 使用 RefreshFromPath 替换模式
	err = filter2.RefreshFromPath(testFile, true) // true = 替换模式
	if err != nil {
		log.Fatalf("刷新词库失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	stats2 := filter2.GetStats()
	fmt.Printf("  替换后总词数: %d\n", stats2.TotalWords)
	fmt.Printf("  来源: %v\n\n", stats2.Source)

	// 方式 3: 加载多个文件
	fmt.Println("--- 方式 3: 加载多个文件 ---")
	filter3, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 创建第二个测试文件
	testFile2 := filepath.Join(tmpDir, "test_words2.txt")
	testWords2 := `新词1
新词2
`
	err = os.WriteFile(testFile2, []byte(testWords2), 0644)
	if err != nil {
		log.Fatalf("创建测试文件2失败: %v", err)
	}
	defer os.Remove(testFile2)

	err = filter3.LoadDictPath(testFile, testFile2)
	if err != nil {
		log.Fatalf("加载多个文件失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	stats3 := filter3.GetStats()
	fmt.Printf("  加载后总词数: %d\n", stats3.TotalWords)
	fmt.Printf("  来源: %v\n\n", stats3.Source)

	// 测试加载的词是否生效
	testText := "这是一个包含敏感词1和新词2的测试文本"
	fmt.Println("--- 测试加载的词库 ---")
	fmt.Printf("  测试文本: %s\n", testText)
	fmt.Printf("  是否敏感: %v\n", filter3.IsSensitive(testText))
	fmt.Printf("  找到的敏感词: %v\n", filter3.FindAll(testText))

	fmt.Println("\n=== 从文件加载示例完成 ===")
}

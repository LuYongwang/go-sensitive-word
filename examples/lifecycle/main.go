package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 演示资源管理和生命周期管理
// 在生产环境中，应该正确关闭资源以避免 goroutine 泄漏
func main() {
	fmt.Println("=== 资源管理和生命周期示例 ===\n")

	// 1. 基本使用（自动关闭）
	fmt.Println("--- 基本使用 ---")
	filter1, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	filter1.AddWords([]string{"测试词"})
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("  测试: %v\n", filter1.IsSensitive("这是一个测试词的文本"))

	// 使用完毕后关闭
	err = filter1.Close()
	if err != nil {
		log.Printf("关闭失败: %v", err)
	} else {
		fmt.Println("  ✅ 资源已关闭\n")
	}

	// 2. 优雅关闭（带超时）
	fmt.Println("--- 优雅关闭（Shutdown） ---")
	filter2, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC}, // AC 有异步处理
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	filter2.AddWords([]string{"测试词1", "测试词2"})
	time.Sleep(200 * time.Millisecond)

	// 使用 context 控制关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = filter2.Shutdown(ctx)
	if err != nil {
		log.Printf("优雅关闭失败: %v", err)
	} else {
		fmt.Println("  ✅ 已优雅关闭（等待异步处理完成）\n")
	}

	// 3. 在生产环境中的最佳实践
	fmt.Println("--- 生产环境最佳实践 ---")
	fmt.Println(`
  在生产环境中，应该：

  1. 在应用启动时初始化过滤器
  2. 在整个应用生命周期内复用同一个实例
  3. 在应用关闭时调用 Shutdown() 优雅关闭

  示例代码结构：
  var globalFilter *sensitive.Manager
  
  func init() {
      var err error
      globalFilter, err = sensitive.NewFilter(...)
      if err != nil {
          log.Fatal(err)
      }
      // 加载词库...
  }
  
  func main() {
      // 注册优雅关闭
      go func() {
          // sig := <-signal.Notify(...)
          ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
          defer cancel()
          globalFilter.Shutdown(ctx)
          // os.Exit(0)
      }()
      
      // 启动服务...
  }
  `)

	// 4. 并发安全示例
	fmt.Println("--- 并发安全演示 ---")
	filter3, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterAC},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer filter3.Shutdown(context.Background())

	// 并发添加词
	for i := 0; i < 10; i++ {
		go func(id int) {
			words := []string{fmt.Sprintf("并发词%d", id)}
			if err := filter3.AddWords(words); err != nil {
				log.Printf("添加失败: %v", err)
			}
		}(i)
	}

	// 并发查询
	for i := 0; i < 10; i++ {
		go func(id int) {
			text := fmt.Sprintf("包含并发词%d的测试", id)
			_ = filter3.IsSensitive(text)
		}(i)
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("  ✅ 并发操作完成，无数据竞争")

	stats := filter3.GetStats()
	fmt.Printf("  最终词数: %d\n", stats.TotalWords)

	fmt.Println("\n=== 资源管理示例完成 ===")
}

package main

import (
	"fmt"
	"log"

	sensitive "github.com/LuYongwang/go-sensitive-word"
)

// 示例：从数据库加载词库
func loadFromDatabase() ([]string, error) {
	// 模拟数据库查询
	// 实际使用中，这里可以是真实的数据库查询逻辑
	// 例如：SELECT word FROM sensitive_words WHERE status = 1
	words := []string{
		"违规词1",
		"违规词2",
		"违规词3",
	}
	return words, nil
}

// 示例：从 Redis 加载词库
func loadFromRedis() ([]string, error) {
	// 模拟 Redis 查询
	// 实际使用中，这里可以是真实的 Redis GET/SMEMBERS 等操作
	// 例如：redis.SMembers(ctx, "sensitive:words")
	words := []string{
		"敏感词A",
		"敏感词B",
	}
	return words, nil
}

// 示例：从配置中心加载词库
func loadFromConfigCenter() ([]string, error) {
	// 模拟配置中心查询
	// 实际使用中，这里可以是真实的配置中心 API 调用
	// 例如：configCenter.GetConfig("sensitive-words")
	words := []string{
		"禁用词X",
		"禁用词Y",
	}
	return words, nil
}

// 示例：从多个数据源合并加载
func loadFromMultipleSources() ([]string, error) {
	words := make([]string, 0)

	// 从数据库加载
	dbWords, err := loadFromDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to load from database: %w", err)
	}
	words = append(words, dbWords...)

	// 从 Redis 加载
	redisWords, err := loadFromRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to load from redis: %w", err)
	}
	words = append(words, redisWords...)

	return words, nil
}

func main() {
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 方式 1: 从数据库加载
	fmt.Println("=== 从数据库加载词库 ===")
	err = filter.LoadDictCallback(loadFromDatabase, "database")
	if err != nil {
		log.Printf("从数据库加载失败: %v", err)
	} else {
		fmt.Println("✅ 从数据库加载成功")
	}

	// 方式 2: 从 Redis 加载
	fmt.Println("\n=== 从 Redis 加载词库 ===")
	err = filter.LoadDictCallback(loadFromRedis, "redis")
	if err != nil {
		log.Printf("从 Redis 加载失败: %v", err)
	} else {
		fmt.Println("✅ 从 Redis 加载成功")
	}

	// 方式 3: 从配置中心加载
	fmt.Println("\n=== 从配置中心加载词库 ===")
	err = filter.LoadDictCallback(loadFromConfigCenter, "config-center")
	if err != nil {
		log.Printf("从配置中心加载失败: %v", err)
	} else {
		fmt.Println("✅ 从配置中心加载成功")
	}

	// 方式 4: 从多个数据源合并加载
	fmt.Println("\n=== 从多个数据源合并加载 ===")
	err = filter.LoadDictCallback(loadFromMultipleSources, "multi-source")
	if err != nil {
		log.Printf("从多数据源加载失败: %v", err)
	} else {
		fmt.Println("✅ 从多数据源加载成功")
	}

	// 查看统计信息
	stats := filter.GetStats()
	fmt.Printf("\n=== 词库统计 ===\n")
	fmt.Printf("总词数: %d\n", stats.TotalWords)
	fmt.Printf("来源: %v\n", stats.Source)

	// 使用内联匿名函数
	fmt.Println("\n=== 使用内联匿名函数加载 ===")
	err = filter.LoadDictCallback(func() ([]string, error) {
		// 在这里可以调用任何自定义逻辑
		return []string{"内联词1", "内联词2"}, nil
	}, "inline")
	if err != nil {
		log.Printf("内联加载失败: %v", err)
	} else {
		fmt.Println("✅ 内联加载成功")
	}

	// 测试过滤功能
	testText := "这是一个包含违规词1和敏感词A的测试文本"
	fmt.Printf("\n=== 测试文本 ===\n%s\n", testText)
	fmt.Printf("是否敏感: %v\n", filter.IsSensitive(testText))
	fmt.Printf("找到的敏感词: %v\n", filter.FindAll(testText))
}

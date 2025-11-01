package go_sensitive_word

import (
	"testing"
	"time"
)

// 性能测试：DFA 算法
func BenchmarkDFA_IsSensitive(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterDfa},
	)

	// 加载所有内置词库
	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	// 等待加载完成
	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	// 预热
	_ = filter.IsSensitive(testText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.IsSensitive(testText)
	}
}

func BenchmarkDFA_FindAll(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterDfa},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.FindAll(testText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.FindAll(testText)
	}
}

func BenchmarkDFA_Replace(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterDfa},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.Replace(testText, '*')

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Replace(testText, '*')
	}
}

// 性能测试：AC 算法
func BenchmarkAC_IsSensitive(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.IsSensitive(testText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.IsSensitive(testText)
	}
}

func BenchmarkAC_FindAll(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.FindAll(testText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.FindAll(testText)
	}
}

func BenchmarkAC_Replace(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.Replace(testText, '*')

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Replace(testText, '*')
	}
}

// 性能测试：来源追踪
func BenchmarkSourceTracking_GetWordSources(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbedWithSource(DictPolitical, "political")
	filter.AddWordsWithSource([]string{"违禁词A", "违禁词B"}, "custom")

	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.GetWordSources("违禁词A")
	}
}

func BenchmarkSourceTracking_FindAllWithSource(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbedWithSource(DictPolitical, "political")
	filter.AddWordsWithSource([]string{"违禁词A", "违禁词B"}, "custom")

	time.Sleep(100 * time.Millisecond)

	testText := "这段文本包含违禁词A和违禁词B，以及温云松"
	_ = filter.FindAllWithSource(testText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.FindAllWithSource(testText)
	}
}

// 并发性能测试
func BenchmarkAC_IsSensitive_Parallel(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.IsSensitive(testText)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = filter.IsSensitive(testText)
		}
	})
}

func BenchmarkDFA_IsSensitive_Parallel(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterDfa},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	testText := "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	_ = filter.IsSensitive(testText)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = filter.IsSensitive(testText)
		}
	})
}

// 长文本性能测试
func BenchmarkAC_LongText(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterAC},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	// 生成长文本（约5000字符）
	longText := ""
	for i := 0; i < 500; i++ {
		longText += "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	}

	_ = filter.IsSensitive(longText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.IsSensitive(longText)
	}
}

func BenchmarkDFA_LongText(b *testing.B) {
	filter, _ := NewFilter(
		StoreOption{Type: StoreMemory},
		FilterOption{Type: FilterDfa},
	)

	filter.LoadDictEmbed(
		DictReactionary,
		DictAdvertisement,
		DictPolitical,
		DictViolence,
		DictPeopleLife,
		DictGunExplosion,
		DictPornography,
		DictCorruption,
	)

	time.Sleep(100 * time.Millisecond)

	longText := ""
	for i := 0; i < 500; i++ {
		longText += "这是一个测试文本包含多个敏感词台湾国毒品销售违禁内容"
	}

	_ = filter.IsSensitive(longText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.IsSensitive(longText)
	}
}

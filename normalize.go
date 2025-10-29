package go_sensitive_word

import (
	"github.com/LuYongwang/go-sensitive-word/internal/normalize"
)

// NormalizerConfig 定义文本归一化的策略
type NormalizerConfig struct {
	IgnoreCase         bool          // 忽略大小写（转小写）
	ToHalfWidth        bool          // 全角转半角
	IgnoreRepeat       bool          // 忽略连续重复字符（将连续相同字符压缩为 1 个）
	IgnoreDigitType    bool          // 归一化各种数字写法为阿拉伯数字
	IgnoreSimpTrad     bool          // 繁简归一（繁体字转简体）
	IgnoreEnglishStyle bool          // 归一化英文变体（花体、数学字母等）
	RemoveZeroWidth    bool          // 剔除零宽字符（防止绕过）
	HomoglyphMap       map[rune]rune // 同形字映射表（防止混淆字符绕过）
}

// toInternalConfig 将公开配置转换为内部配置
func (c NormalizerConfig) toInternalConfig() normalize.Config {
	return normalize.Config{
		IgnoreCase:         c.IgnoreCase,
		ToHalfWidth:        c.ToHalfWidth,
		IgnoreRepeat:       c.IgnoreRepeat,
		IgnoreDigitType:    c.IgnoreDigitType,
		IgnoreSimpTrad:     c.IgnoreSimpTrad,
		IgnoreEnglishStyle: c.IgnoreEnglishStyle,
		RemoveZeroWidth:    c.RemoveZeroWidth,
		HomoglyphMap:       c.HomoglyphMap,
	}
}

// DefaultNormalizer 返回默认归一化配置
// 默认启用：忽略大小写、全角转半角；不忽略重复
func DefaultNormalizer() NormalizerConfig {
	return NormalizerConfig{
		IgnoreCase:      true,
		ToHalfWidth:     true,
		IgnoreRepeat:    false,
		RemoveZeroWidth: false, // 默认不剔除零宽字符，可通过配置开启
		HomoglyphMap:    nil,   // 默认不同形字映射，可通过配置开启
	}
}

// StrictNormalizer 返回严格归一化配置
// 启用所有归一化选项，提供最强的防绕过能力
func StrictNormalizer() NormalizerConfig {
	return NormalizerConfig{
		IgnoreCase:         true,
		ToHalfWidth:        true,
		IgnoreRepeat:       true,                            // 开启：忽略连续重复字符
		IgnoreDigitType:    true,                            // 开启：数字归一化
		IgnoreSimpTrad:     true,                            // 开启：繁简归一
		IgnoreEnglishStyle: true,                            // 开启：英文变体归一
		RemoveZeroWidth:    true,                            // 开启：剔除零宽字符
		HomoglyphMap:       normalize.DefaultHomoglyphMap(), // 使用默认同形字映射
	}
}

// NormalizeTextWithMap 对文本做归一化，同时返回从规范化索引到原始索引的映射
// 返回：规范化后的字符串、规范化索引 -> 原始索引 的映射
func NormalizeTextWithMap(s string, cfg NormalizerConfig) (string, []int) {
	return normalize.NormalizeTextWithMap(s, cfg.toInternalConfig())
}

// NormalizeWord 对词库词条做归一化（与文本侧保持同策略）
func NormalizeWord(word string, cfg NormalizerConfig) string {
	return normalize.NormalizeWord(word, cfg.toInternalConfig())
}

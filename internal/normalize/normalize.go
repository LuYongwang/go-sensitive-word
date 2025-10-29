package normalize

import (
	"unicode"

	"github.com/LuYongwang/go-sensitive-word/internal/jianfan"
)

// Config 定义文本归一化的策略（内部实现）
type Config struct {
	IgnoreCase         bool
	ToHalfWidth        bool
	IgnoreRepeat       bool
	IgnoreDigitType    bool
	IgnoreSimpTrad     bool
	IgnoreEnglishStyle bool
	RemoveZeroWidth    bool
	HomoglyphMap       map[rune]rune
}

// 数字映射表：各种数字写法 -> 阿拉伯数字
var digitMap = map[rune]rune{
	'０': '0', '１': '1', '２': '2', '３': '3', '４': '4', '５': '5', '６': '6', '７': '7', '８': '8', '９': '9',
	'⓪': '0', '①': '1', '②': '2', '③': '3', '④': '4', '⑤': '5', '⑥': '6', '⑦': '7', '⑧': '8', '⑨': '9',
	'⓿': '0', '❶': '1', '❷': '2', '❸': '3', '❹': '4', '❺': '5', '❻': '6', '❼': '7', '❽': '8', '❾': '9',
	'⁰': '0', '¹': '1', '²': '2', '³': '3', '⁴': '4', '⁵': '5', '⁶': '6', '⁷': '7', '⁸': '8', '⁹': '9',
	'₀': '0', '₁': '1', '₂': '2', '₃': '3', '₄': '4', '₅': '5', '₆': '6', '₇': '7', '₈': '8', '₉': '9',
	'➀': '1', '➁': '2', '➂': '3', '➃': '4', '➄': '5', '➅': '6', '➆': '7', '➇': '8', '➈': '9',
	'➊': '1', '➋': '2', '➌': '3', '➍': '4', '➎': '5', '➏': '6', '➐': '7', '➑': '8', '➒': '9',
	'㈠': '1', '㈡': '2', '㈢': '3', '㈣': '4', '㈤': '5', '㈥': '6', '㈦': '7', '㈧': '8', '㈨': '9',
	'㊀': '1', '㊁': '2', '㊂': '3', '㊃': '4', '㊄': '5', '㊅': '6', '㊆': '7', '㊇': '8', '㊈': '9',
	// 中文数字
	'零': '0', '一': '1', '二': '2', '三': '3', '四': '4', '五': '5', '六': '6', '七': '7', '八': '8', '九': '9',
	'〇': '0', '壹': '1', '贰': '2', '叁': '3', '肆': '4', '伍': '5', '陆': '6', '柒': '7', '捌': '8', '玖': '9',
}

// 英文变体映射表（花体、数学字母等 -> 基本拉丁字母）
var englishVariantMap = map[rune]rune{
	'Ⓕ': 'F', 'Ⓖ': 'G',
	'Ａ': 'A', 'Ｂ': 'B', 'Ｃ': 'C', 'Ｄ': 'D', 'Ｅ': 'E', 'Ｆ': 'F',
	'Ｇ': 'G', 'Ｈ': 'H', 'Ｉ': 'I', 'Ｊ': 'J', 'Ｋ': 'K', 'Ｌ': 'L',
	'Ｍ': 'M', 'Ｎ': 'N', 'Ｏ': 'O', 'Ｐ': 'P', 'Ｑ': 'Q', 'Ｒ': 'R',
	'Ｓ': 'S', 'Ｔ': 'T', 'Ｕ': 'U', 'Ｖ': 'V', 'Ｗ': 'W', 'Ｘ': 'X',
	'Ｙ': 'Y', 'Ｚ': 'Z', 'ａ': 'a', 'ｂ': 'b', 'ｃ': 'c', 'ｄ': 'd',
	'ｅ': 'e', 'ｆ': 'f', 'ｇ': 'g', 'ｈ': 'h', 'ｉ': 'i', 'ｊ': 'j',
	'ｋ': 'k', 'ｌ': 'l', 'ｍ': 'm', 'ｎ': 'n', 'ｏ': 'o', 'ｐ': 'p',
	'ｑ': 'q', 'ｒ': 'r', 'ｓ': 's', 'ｔ': 't', 'ｕ': 'u', 'ｖ': 'v',
	'ｗ': 'w', 'ｘ': 'x', 'ｙ': 'y', 'ｚ': 'z',
}

// 零宽字符列表
// 这些字符不可见但存在于文本中，常用于绕过敏感词检测
var zeroWidthChars = map[rune]bool{
	'\u200B': true, // Zero Width Space
	'\u200C': true, // Zero Width Non-Joiner
	'\u200D': true, // Zero Width Joiner
	'\uFEFF': true, // Zero Width No-Break Space (BOM)
	'\u200E': true, // Left-to-Right Mark
	'\u200F': true, // Right-to-Left Mark
	'\u202A': true, // Left-to-Right Embedding
	'\u202B': true, // Right-to-Left Embedding
	'\u202C': true, // Pop Directional Formatting
	'\u202D': true, // Left-to-Right Override
	'\u202E': true, // Right-to-Left Override
	'\u2060': true, // Word Joiner
	'\u2061': true, // Function Application
	'\u2062': true, // Invisible Times
	'\u2063': true, // Invisible Separator
	'\u2064': true, // Invisible Plus
}

// 同形字映射表：容易混淆的相似字符 -> 标准字符
var defaultHomoglyphMap = map[rune]rune{
	// 全角字母 -> 半角字母（大写）
	'Ａ': 'A', 'Ｂ': 'B', 'Ｃ': 'C', 'Ｄ': 'D', 'Ｅ': 'E', 'Ｆ': 'F',
	'Ｇ': 'G', 'Ｈ': 'H', 'Ｉ': 'I', 'Ｊ': 'J', 'Ｋ': 'K', 'Ｌ': 'L',
	'Ｍ': 'M', 'Ｎ': 'N', 'Ｏ': 'O', 'Ｐ': 'P', 'Ｑ': 'Q', 'Ｒ': 'R',
	'Ｓ': 'S', 'Ｔ': 'T', 'Ｕ': 'U', 'Ｖ': 'V', 'Ｗ': 'W', 'Ｘ': 'X',
	'Ｙ': 'Y', 'Ｚ': 'Z',
	// 全角字母 -> 半角字母（小写）
	'ａ': 'a', 'ｂ': 'b', 'ｃ': 'c', 'ｄ': 'd', 'ｅ': 'e', 'ｆ': 'f',
	'ｇ': 'g', 'ｈ': 'h', 'ｉ': 'i', 'ｊ': 'j', 'ｋ': 'k', 'ｌ': 'l',
	'ｍ': 'm', 'ｎ': 'n', 'ｏ': 'o', 'ｐ': 'p', 'ｑ': 'q', 'ｒ': 'r',
	'ｓ': 's', 'ｔ': 't', 'ｕ': 'u', 'ｖ': 'v', 'ｗ': 'w', 'ｘ': 'x',
	'ｙ': 'y', 'ｚ': 'z',
	// 易混淆的相似字符
	'０': '0', '１': '1', '２': '2', '３': '3', '４': '4', '５': '5',
	'６': '6', '７': '7', '８': '8', '９': '9',
	'ο': 'o', // Greek small letter omicron -> lowercase o
	'а': 'a', // Cyrillic small letter a -> lowercase a
	'е': 'e', // Cyrillic small letter e -> lowercase e
	'о': 'o', // Cyrillic small letter o -> lowercase o
	'р': 'p', // Cyrillic small letter p -> lowercase p
	'с': 'c', // Cyrillic small letter c -> lowercase c
	'у': 'y', // Cyrillic small letter u -> lowercase y
	'х': 'x', // Cyrillic small letter ha -> lowercase x
	'ι': 'i', // Greek small letter iota -> lowercase i
	'τ': 't', // Greek small letter tau -> lowercase t
	// 更多混淆字符可以继续添加
}

// DefaultHomoglyphMap 返回默认同形字映射表的副本
func DefaultHomoglyphMap() map[rune]rune {
	m := make(map[rune]rune, len(defaultHomoglyphMap))
	for k, v := range defaultHomoglyphMap {
		m[k] = v
	}
	return m
}

// toHalfWidth 将常见全角字符转换为半角（ASCII 范围）
// 规则：
// - 全角空格（U+3000）-> 半角空格（U+0020）
// - 全角 U+FF01~U+FF5E -> 减去 0xFEE0
func toHalfWidth(r rune) rune {
	if r == '\u3000' {
		return ' '
	}
	if r >= '\uff01' && r <= '\uff5e' {
		return r - 0xFEE0
	}
	return r
}

// normalizeRune 根据配置对单个字符做归一化
func normalizeRune(r rune, cfg Config) rune {
	// 1. 零宽字符剔除
	if cfg.RemoveZeroWidth {
		if zeroWidthChars[r] {
			return rune(-1) // 返回特殊标记，后续会跳过
		}
	}

	// 2. 同形字映射
	if cfg.HomoglyphMap != nil {
		if mapped, ok := cfg.HomoglyphMap[r]; ok {
			r = mapped
		}
	}

	// 3. 数字归一化
	if cfg.IgnoreDigitType {
		if digit, ok := digitMap[r]; ok {
			r = digit
		}
	}

	// 4. 繁简归一（繁体 -> 简体）
	if cfg.IgnoreSimpTrad {
		// 单字符维度转换：若为繁体则转为简体
		s := string(r)
		cn := jianfan.T2S(s)
		if len([]rune(cn)) == 1 {
			r = []rune(cn)[0]
		}
	}

	// 5. 英文变体归一
	if cfg.IgnoreEnglishStyle {
		if eng, ok := englishVariantMap[r]; ok {
			r = eng
		}
	}

	// 6. 全角转半角
	if cfg.ToHalfWidth {
		r = toHalfWidth(r)
	}

	// 7. 大小写归一（最后执行）
	if cfg.IgnoreCase {
		r = unicode.ToLower(r)
	}

	return r
}

// NormalizeTextWithMap 对文本做归一化，同时返回从规范化索引到原始索引的映射
// 返回：规范化后的字符串、规范化索引 -> 原始索引 的映射
func NormalizeTextWithMap(s string, cfg Config) (string, []int) {
	runes := []rune(s)
	norm := make([]rune, 0, len(runes))
	idxMap := make([]int, 0, len(runes))

	var last rune
	var hasLast bool
	for i, r := range runes {
		nr := normalizeRune(r, cfg)
		// 跳过零宽字符（标记为 -1）
		if nr == rune(-1) {
			continue
		}
		if cfg.IgnoreRepeat {
			if hasLast && nr == last {
				// 跳过连续重复
				continue
			}
			last = nr
			hasLast = true
		}
		norm = append(norm, nr)
		idxMap = append(idxMap, i)
	}

	return string(norm), idxMap
}

// NormalizeWord 对词库词条做归一化（与文本侧保持同策略）
func NormalizeWord(word string, cfg Config) string {
	n, _ := NormalizeTextWithMap(word, cfg)
	return n
}

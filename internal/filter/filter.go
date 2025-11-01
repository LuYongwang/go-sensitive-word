package filter

// Range 表示匹配区间 [Start, End]，闭区间
type Range struct {
	Start int // 起始位置（包含）
	End   int // 结束位置（包含）
}

// MatchResult 包含匹配词和其来源信息
type MatchResult struct {
	Word   string   // 匹配到的敏感词
	Source []string // 该词所属的词库来源列表
}

// RangedFilter 是可选的扩展接口，返回匹配区间而非字符串
// 如果实现类支持此接口，包装器会优先使用以提高性能
type RangedFilter interface {
	FindAllRanges(text string) []Range
}

type (
	Filter interface {
		FindAll(text string) []string
		FindAllCount(text string) map[string]int
		FindOne(text string) string
		IsSensitive(text string) bool
		Replace(text string, repl rune) string
		Remove(text string) string
	}

	// SourceFilter 支持返回匹配结果及其来源信息的扩展接口
	SourceFilter interface {
		Filter
		FindAllWithSource(text string) []MatchResult
		FindAllCountWithSource(text string) map[string]MatchResult
	}
)

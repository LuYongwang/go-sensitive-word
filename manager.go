package go_sensitive_word

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/LuYongwang/go-sensitive-word/internal/filter"
	"github.com/LuYongwang/go-sensitive-word/internal/filter/ac"
	"github.com/LuYongwang/go-sensitive-word/internal/filter/dfa"
	"github.com/LuYongwang/go-sensitive-word/internal/store"
)

// Manager 是敏感词过滤系统的核心结构，整合了词库存储和过滤算法
type Manager struct {
	store.Store                    // 词库存储接口（支持内存、本地文件、远程等）
	filter.Filter                  // 敏感词匹配算法接口（如 DFA、Aho-Corasick）
	normalizer    NormalizerConfig // 归一化配置，用于确保词库的词和测试文本的归一化一致
}

// NewFilter 初始化过滤器和词库存储
// 参数：storeOption 指定存储方式，filterOption 指定过滤算法
func NewFilter(storeOption StoreOption, filterOption FilterOption) (*Manager, error) {
	var filterStore store.Store
	var myFilter filter.Filter

	switch storeOption.Type {
	case StoreMemory: // 使用内存词库
		filterStore = store.NewMemoryModel()
	default:
		return nil, errors.New("invalid store type")
	}

	switch filterOption.Type {
	case FilterDfa: // 使用 DFA 算法
		dfaModel := dfa.NewDFAModel()
		// 启动监听协程，实时接收新增/删除词的通知
		go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())
		myFilter = dfaModel
	case FilterAC: // 使用 AC 自动机算法
		acModel := ac.NewACModel()
		go acModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())
		myFilter = acModel
	default:
		return nil, errors.New("invalid filter type")
	}

	// 默认开启大小写与全角归一化，使匹配对大小写/全角不敏感
	normalizerCfg := DefaultNormalizer()
	wrapped := newNormalizedFilter(myFilter, normalizerCfg)

	return &Manager{
		Store:      filterStore,
		Filter:     wrapped,
		normalizer: normalizerCfg,
	}, nil
}

// Close 关闭内部资源
func (m *Manager) Close() error {
	if m.Store != nil {
		if c, ok := m.Store.(interface{ Close() error }); ok {
			return c.Close()
		}
	}
	return nil
}

// Shutdown 优雅关闭
func (m *Manager) Shutdown(ctx context.Context) error {
	if m.Store != nil {
		if s, ok := m.Store.(interface{ Shutdown(context.Context) error }); ok {
			return s.Shutdown(ctx)
		}
	}
	return nil
}

// ==================== 动态维护词库增强方法 ====================

// GetStats 获取词库统计信息
func (m *Manager) GetStats() store.Stats {
	if m.Store == nil {
		return store.Stats{}
	}
	return m.Store.GetStats()
}

// AddWord 添加敏感词（支持多个）
// 注意：词会被归一化后再添加到词库，确保与测试文本的归一化策略一致
func (m *Manager) AddWord(words ...string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化，确保词库的词和测试文本的归一化一致
	normalizedWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedWords = append(normalizedWords, normalized)
		}
	}
	if len(normalizedWords) == 0 {
		return nil
	}
	return m.Store.AddWords(normalizedWords)
}

// AddWords 批量添加敏感词
// 注意：词会被归一化后再添加到词库，确保与测试文本的归一化策略一致
func (m *Manager) AddWords(words []string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化，确保词库的词和测试文本的归一化一致
	normalizedWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedWords = append(normalizedWords, normalized)
		}
	}
	if len(normalizedWords) == 0 {
		return nil
	}
	return m.Store.AddWords(normalizedWords)
}

// DelWord 删除敏感词（支持多个）
// 注意：词会被归一化后再删除，确保与词库中的归一化词匹配
func (m *Manager) DelWord(words ...string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化
	normalizedWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedWords = append(normalizedWords, normalized)
		}
	}
	if len(normalizedWords) == 0 {
		return nil
	}
	return m.Store.DelWords(normalizedWords)
}

// DelWords 批量删除敏感词
// 注意：词会被归一化后再删除，确保与词库中的归一化词匹配
func (m *Manager) DelWords(words []string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化
	normalizedWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedWords = append(normalizedWords, normalized)
		}
	}
	if len(normalizedWords) == 0 {
		return nil
	}
	return m.Store.DelWords(normalizedWords)
}

// ReplaceWords 批量替换敏感词（先删除旧词，再添加新词）
// 注意：词会被归一化后再处理，确保与测试文本的归一化策略一致
func (m *Manager) ReplaceWords(oldWords, newWords []string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化
	normalizedOldWords := make([]string, 0, len(oldWords))
	for _, word := range oldWords {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedOldWords = append(normalizedOldWords, normalized)
		}
	}
	normalizedNewWords := make([]string, 0, len(newWords))
	for _, word := range newWords {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedNewWords = append(normalizedNewWords, normalized)
		}
	}
	return m.Store.ReplaceWords(normalizedOldWords, normalizedNewWords)
}

// ExportToFile 导出词库到文件
func (m *Manager) ExportToFile(filepath string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return m.ExportToWriter(f)
}

// ExportToString 导出词库为字符串
func (m *Manager) ExportToString() (string, error) {
	if m.Store == nil {
		return "", errors.New("store is nil")
	}
	return m.Store.ExportToString()
}

// Clear 清空词库
func (m *Manager) Clear() error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	return m.Store.Clear()
}

// MergeFromManager 从另一个 Manager 合并词库
func (m *Manager) MergeFromManager(other *Manager) error {
	if m.Store == nil || other == nil || other.Store == nil {
		return errors.New("invalid store")
	}
	return m.Merge(other.Store)
}

// RefreshFromPath 从文件路径刷新词库（可选：完全替换或追加）
func (m *Manager) RefreshFromPath(path string, replace bool) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	if replace {
		if err := m.Clear(); err != nil {
			return err
		}
	}
	return m.LoadDictPath(path)
}

// LoadDictCallback 通过回调函数加载词库
// 使用场景：从数据库、Redis、配置中心等自定义数据源读取词库
// loader: 回调函数，返回词列表和错误
// source: 词库来源标识（用于统计信息），如 "database", "redis", "config-center" 等
func (m *Manager) LoadDictCallback(loader store.DictLoader, source string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	return m.Store.LoadDictCallback(loader, source)
}

// ==================== 词库来源追踪功能 ====================

// LoadDictEmbedWithSource 加载内置词库并指定来源名称
// content: 词库内容字符串
// source: 来源标识，如 "political", "violence", "custom" 等
func (m *Manager) LoadDictEmbedWithSource(content string, source string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 读取词库内容并按行拆分
	lines := strings.Split(content, "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		words = append(words, line)
	}
	// 使用 AddWordsWithSource 添加，会自动归一化
	return m.AddWordsWithSource(words, source)
}

// AddWordsWithSource 批量添加敏感词并指定来源
// words: 敏感词列表
// source: 来源标识，如 "political", "violence", "custom" 等
func (m *Manager) AddWordsWithSource(words []string, source string) error {
	if m.Store == nil {
		return errors.New("store is nil")
	}
	// 对词进行归一化
	normalizedWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := NormalizeWord(word, m.normalizer)
		if normalized != "" {
			normalizedWords = append(normalizedWords, normalized)
		}
	}
	if len(normalizedWords) == 0 {
		return nil
	}
	return m.Store.AddWordsWithSource(normalizedWords, source)
}

// GetWordSources 获取指定词的来源列表
func (m *Manager) GetWordSources(word string) []string {
	if m.Store == nil {
		return nil
	}
	normalized := NormalizeWord(word, m.normalizer)
	return m.Store.GetWordSources(normalized)
}

// GetAllWordSources 获取所有词的来源映射
func (m *Manager) GetAllWordSources() map[string][]string {
	if m.Store == nil {
		return nil
	}
	return m.Store.GetAllWordSources()
}

// FindAllWithSource 查找文本中所有敏感词及其来源信息
func (m *Manager) FindAllWithSource(text string) []filter.MatchResult {
	words := m.FindAll(text)
	if len(words) == 0 {
		return []filter.MatchResult{}
	}

	result := make([]filter.MatchResult, 0, len(words))
	seen := make(map[string]bool)
	for _, word := range words {
		if seen[word] {
			continue
		}
		seen[word] = true
		sources := m.GetWordSources(word)
		result = append(result, filter.MatchResult{
			Word:   word,
			Source: sources,
		})
	}
	return result
}

// FindAllCountWithSource 查找所有敏感词及其出现次数和来源信息
func (m *Manager) FindAllCountWithSource(text string) map[string]filter.MatchResult {
	countMap := m.FindAllCount(text)
	result := make(map[string]filter.MatchResult, len(countMap))
	for word := range countMap {
		sources := m.GetWordSources(word)
		result[word] = filter.MatchResult{
			Word:   word,
			Source: sources,
		}
	}
	return result
}

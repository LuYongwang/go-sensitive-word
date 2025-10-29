package store

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MemoryModel struct {
	store      map[string]struct{}
	storeMu    sync.RWMutex // 保护词库 map
	totalWords atomic.Int64 // 原子计数，避免 O(n) 的 Count()
	addChan    chan string
	delChan    chan string
	closed     chan struct{}
	mu         sync.RWMutex // 保护统计信息
	stats      Stats
	sources    []string // 记录加载来源
}

func NewMemoryModel() *MemoryModel {
	return &MemoryModel{
		store:   make(map[string]struct{}),
		addChan: make(chan string, 8192),
		delChan: make(chan string, 8192),
		closed:  make(chan struct{}),
		stats: Stats{
			TotalWords:  0,
			LastUpdate:  time.Now(),
			UpdateCount: 0,
			Source:      make([]string, 0),
		},
		sources: make([]string, 0),
	}
}

func (m *MemoryModel) LoadDictPath(paths ...string) error {
	for _, path := range paths {
		err := func(path string) error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()
			err = m.LoadDict(f)
			if err == nil {
				m.mu.Lock()
				m.sources = append(m.sources, "file://"+path)
				m.stats.Source = append(m.stats.Source, "file://"+path)
				m.mu.Unlock()
			}
			return err
		}(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryModel) LoadDictEmbed(contents ...string) error {
	for _, con := range contents {
		reader := strings.NewReader(con)
		if err := m.LoadDict(reader); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryModel) LoadDict(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	count := 0
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		word := strings.TrimSpace(string(line))
		if word == "" {
			continue
		}
		word = strings.ToLower(word)
		m.storeMu.Lock()
		isNew := true
		if _, exists := m.store[word]; !exists {
			m.store[word] = struct{}{}
			m.totalWords.Add(1)
		} else {
			isNew = false
		}
		m.storeMu.Unlock()
		if isNew {
			count++
		}
		select {
		case m.addChan <- word:
		case <-m.closed:
			return errors.New("store closed during load")
		}
	}
	// 更新统计信息
	m.mu.Lock()
	m.stats.TotalWords = int(m.totalWords.Load())
	m.stats.LastUpdate = time.Now()
	m.stats.UpdateCount += count
	m.mu.Unlock()
	return nil
}

func (m *MemoryModel) ReadChan() <-chan string {
	ch := make(chan string)
	go func() {
		m.storeMu.RLock()
		keys := make([]string, 0, len(m.store))
		for key := range m.store {
			keys = append(keys, key)
		}
		m.storeMu.RUnlock()
		for _, key := range keys {
			select {
			case ch <- key:
			case <-m.closed:
				close(ch)
				return
			}
		}
		close(ch)
	}()
	return ch
}

func (m *MemoryModel) ReadString() []string {
	m.storeMu.RLock()
	res := make([]string, 0, len(m.store))
	for key := range m.store {
		res = append(res, key)
	}
	m.storeMu.RUnlock()
	return res
}

func (m *MemoryModel) GetAddChan() <-chan string { return m.addChan }
func (m *MemoryModel) GetDelChan() <-chan string { return m.delChan }

// LoadDictCallback 通过回调函数加载词库
// loader: 回调函数，返回词列表
// source: 词库来源标识（用于统计信息），如 "database", "redis", "config-center" 等
func (m *MemoryModel) LoadDictCallback(loader DictLoader, source string) error {
	if loader == nil {
		return errors.New("loader callback is nil")
	}
	words, err := loader()
	if err != nil {
		return err
	}
	err = m.AddWords(words)
	if err == nil && source != "" {
		m.mu.Lock()
		m.sources = append(m.sources, "callback://"+source)
		m.stats.Source = append(m.stats.Source, "callback://"+source)
		m.mu.Unlock()
	}
	return err
}

func (m *MemoryModel) AddWord(words ...string) error {
	return m.AddWords(words)
}

func (m *MemoryModel) AddWords(words []string) error {
	count := 0
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		// 注意：词应该已经在 Manager 层做了归一化，这里只做 TrimSpace
		// 只有新词才计数
		m.storeMu.Lock()
		if _, exists := m.store[word]; !exists {
			m.store[word] = struct{}{}
			m.totalWords.Add(1)
			count++
		} else {
			m.store[word] = struct{}{} // 确保存在
		}
		m.storeMu.Unlock()
		select {
		case m.addChan <- word:
		case <-m.closed:
			return errors.New("store closed")
		}
	}
	// 更新统计信息
	m.mu.Lock()
	m.stats.TotalWords = int(m.totalWords.Load())
	m.stats.LastUpdate = time.Now()
	m.stats.UpdateCount += count
	m.mu.Unlock()
	return nil
}

func (m *MemoryModel) DelWord(words ...string) error {
	return m.DelWords(words)
}

func (m *MemoryModel) DelWords(words []string) error {
	count := 0
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		// 注意：词应该已经在 Manager 层做了归一化，这里只做 TrimSpace
		m.storeMu.Lock()
		if _, exists := m.store[word]; exists {
			delete(m.store, word)
			m.totalWords.Add(-1)
			count++
		}
		m.storeMu.Unlock()
		select {
		case m.delChan <- word:
		case <-m.closed:
			return errors.New("store closed")
		}
	}
	// 更新统计信息
	m.mu.Lock()
	m.stats.TotalWords = int(m.totalWords.Load())
	m.stats.LastUpdate = time.Now()
	m.stats.UpdateCount += count
	m.mu.Unlock()
	return nil
}

// ReplaceWords 批量替换：先删除旧词，再添加新词
func (m *MemoryModel) ReplaceWords(oldWords, newWords []string) error {
	if err := m.DelWords(oldWords); err != nil {
		return err
	}
	return m.AddWords(newWords)
}

// ExportToWriter 导出词库到 Writer
func (m *MemoryModel) ExportToWriter(w io.Writer) error {
	writer := bufio.NewWriter(w)
	defer func() {
		_ = writer.Flush() // Flush 的错误在这里无法返回，记录但不中断流程
	}()

	m.storeMu.RLock()
	words := make([]string, 0, len(m.store))
	for word := range m.store {
		words = append(words, word)
	}
	m.storeMu.RUnlock()

	for _, word := range words {
		if _, err := writer.WriteString(word + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// ExportToString 导出词库为字符串（每行一个词）
func (m *MemoryModel) ExportToString() (string, error) {
	m.storeMu.RLock()
	words := make([]string, 0, len(m.store))
	for word := range m.store {
		words = append(words, word)
	}
	m.storeMu.RUnlock()

	var builder strings.Builder
	for _, word := range words {
		builder.WriteString(word)
		builder.WriteString("\n")
	}
	return builder.String(), nil
}

// GetStats 获取统计信息
func (m *MemoryModel) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// 返回副本，避免外部修改
	return Stats{
		TotalWords:  int(m.totalWords.Load()),
		LastUpdate:  m.stats.LastUpdate,
		UpdateCount: m.stats.UpdateCount,
		Source:      append([]string{}, m.stats.Source...),
	}
}

// Clear 清空词库
func (m *MemoryModel) Clear() error {
	// 先获取所有词，然后删除
	words := m.ReadString()
	if err := m.DelWords(words); err != nil {
		return err
	}
	// 重置计数器（防止累积误差）
	m.storeMu.Lock()
	m.totalWords.Store(0)
	m.storeMu.Unlock()
	m.mu.Lock()
	m.stats.Source = make([]string, 0)
	m.sources = make([]string, 0)
	m.mu.Unlock()
	return nil
}

// Merge 合并另一个词库
func (m *MemoryModel) Merge(other Store) error {
	words := other.ReadString()
	return m.AddWords(words)
}

func (m *MemoryModel) Close() error {
	select {
	case <-m.closed:
	default:
		close(m.closed)
		close(m.addChan)
		close(m.delChan)
	}
	return nil
}

func (m *MemoryModel) Shutdown(ctx context.Context) error {
	select {
	case <-m.closed:
	default:
		close(m.closed)
	}
	done := make(chan struct{})
	go func() { close(m.addChan); close(m.delChan); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

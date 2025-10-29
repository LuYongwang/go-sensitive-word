package ac

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/LuYongwang/go-sensitive-word/internal/filter"
)

type acNode struct {
	children map[rune]*acNode
	fail     *acNode
	output   []string
}

func newAcNode() *acNode {
	return &acNode{children: make(map[rune]*acNode), output: make([]string, 0)}
}

type ACModel struct {
	rootPtr     atomic.Pointer[acNode] // 原子指针，支持原子切换
	mu          sync.Mutex             // 保护构建过程
	pendingAdds []string               // 待添加的词（窗口合并）
	pendingDels []string               // 待删除的词（窗口合并）
	ticker      *time.Ticker           // 窗口计时器
	done        chan struct{}
}

func NewACModel() *ACModel {
	root := newAcNode()
	model := &ACModel{
		done: make(chan struct{}),
	}
	model.rootPtr.Store(root)
	return model
}

func (m *ACModel) AddWord(word string) {
	if word == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	root := m.rootPtr.Load()
	newRoot := m.cloneNode(root)

	now := newRoot
	for _, r := range word {
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			next = newAcNode()
			now.children[r] = next
			now = next
		}
	}
	now.output = append(now.output, word)

	// 重建失败指针树
	m.buildFailurePointer(newRoot)

	// 原子切换
	m.rootPtr.Store(newRoot)
}

// cloneNode 深拷贝节点及其子树
func (m *ACModel) cloneNode(n *acNode) *acNode {
	newNode := &acNode{
		children: make(map[rune]*acNode),
		output:   make([]string, len(n.output)),
	}
	copy(newNode.output, n.output)
	for r, child := range n.children {
		newNode.children[r] = m.cloneNode(child)
	}
	return newNode
}

func (m *ACModel) AddWords(words ...string) {
	for _, w := range words {
		m.AddWord(w)
	}
}

func (m *ACModel) DelWord(word string) {
	if word == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	root := m.rootPtr.Load()
	newRoot := m.cloneNode(root)

	now := newRoot
	for _, r := range word {
		next, ok := now.children[r]
		if !ok {
			return // 词不存在
		}
		now = next
	}
	for i, w := range now.output {
		if w == word {
			now.output = append(now.output[:i], now.output[i+1:]...)
			break
		}
	}

	// 重建失败指针树
	m.buildFailurePointer(newRoot)

	// 原子切换
	m.rootPtr.Store(newRoot)
}

func (m *ACModel) Delwords(words ...string) {
	for _, w := range words {
		m.DelWord(w)
	}
}

// buildFailurePointer 构建失败指针树（不使用 built 标志，每次重构都重建）
func (m *ACModel) buildFailurePointer(root *acNode) {
	queue := make([]*acNode, 0)
	for _, child := range root.children {
		child.fail = root
		queue = append(queue, child)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for char, child := range current.children {
			queue = append(queue, child)
			temp := current.fail
			for temp != nil && temp.children[char] == nil {
				temp = temp.fail
			}
			if temp == nil {
				child.fail = root
			} else {
				child.fail = temp.children[char]
				child.output = append(child.output, child.fail.output...)
			}
		}
	}
}

// Listen 启动监听协程，支持窗口合并（100ms 或 1000 条）
func (m *ACModel) Listen(addChan, delChan <-chan string) {
	m.mu.Lock()
	m.pendingAdds = make([]string, 0, 1000)
	m.pendingDels = make([]string, 0, 1000)
	m.ticker = time.NewTicker(100 * time.Millisecond) // 100ms 窗口
	m.mu.Unlock()

	go func() {
		defer m.ticker.Stop()
		for {
			select {
			case word := <-addChan:
				if word == "" {
					continue
				}
				m.mu.Lock()
				m.pendingAdds = append(m.pendingAdds, word)
				shouldRebuild := len(m.pendingAdds) >= 1000
				m.mu.Unlock()
				if shouldRebuild {
					m.flushPending()
				}
			case word := <-delChan:
				if word == "" {
					continue
				}
				m.mu.Lock()
				m.pendingDels = append(m.pendingDels, word)
				shouldRebuild := len(m.pendingDels) >= 1000
				m.mu.Unlock()
				if shouldRebuild {
					m.flushPending()
				}
			case <-m.ticker.C:
				m.flushPending()
			case <-m.done:
				m.flushPending() // 最后刷新一次
				return
			}
		}
	}()
}

// flushPending 刷新待处理的词（批量应用）
func (m *ACModel) flushPending() {
	m.mu.Lock()
	if len(m.pendingAdds) == 0 && len(m.pendingDels) == 0 {
		m.mu.Unlock()
		return
	}
	adds := m.pendingAdds
	dels := m.pendingDels
	m.pendingAdds = make([]string, 0, 1000)
	m.pendingDels = make([]string, 0, 1000)
	m.mu.Unlock()

	// 批量应用
	root := m.rootPtr.Load()
	newRoot := m.cloneNode(root)

	// 先删除
	for _, word := range dels {
		now := newRoot
		for _, r := range word {
			next, ok := now.children[r]
			if !ok {
				break
			}
			now = next
		}
		for i, w := range now.output {
			if w == word {
				now.output = append(now.output[:i], now.output[i+1:]...)
				break
			}
		}
	}

	// 再添加
	for _, word := range adds {
		now := newRoot
		for _, r := range word {
			if next, ok := now.children[r]; ok {
				now = next
			} else {
				next = newAcNode()
				now.children[r] = next
				now = next
			}
		}
		now.output = append(now.output, word)
	}

	// 重建失败指针树
	m.buildFailurePointer(newRoot)

	// 原子切换
	m.rootPtr.Store(newRoot)
}

func (m *ACModel) FindAll(text string) []string {
	root := m.rootPtr.Load()
	var matches []string
	seen := make(map[string]struct{})
	now := root
	runes := []rune(text)
	for i, r := range runes {
		for now != root && now.children[r] == nil {
			now = now.fail
		}
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			now = root
		}
		for _, w := range now.output {
			if _, ok := seen[w]; !ok {
				seen[w] = struct{}{}
				matches = append(matches, w)
			}
		}
		_ = i // 防止未使用变量
	}
	return matches
}

func (m *ACModel) FindAllCount(text string) map[string]int {
	root := m.rootPtr.Load()
	counts := make(map[string]int)
	now := root
	for _, r := range text {
		for now != root && now.children[r] == nil {
			now = now.fail
		}
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			now = root
		}
		for _, w := range now.output {
			counts[w]++
		}
	}
	return counts
}

func (m *ACModel) FindOne(text string) string {
	root := m.rootPtr.Load()
	now := root
	for _, r := range text {
		for now != root && now.children[r] == nil {
			now = now.fail
		}
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			now = root
		}
		if len(now.output) > 0 {
			longest := now.output[0]
			for i := 1; i < len(now.output); i++ {
				if len([]rune(now.output[i])) > len([]rune(longest)) {
					longest = now.output[i]
				}
			}
			return longest
		}
	}
	return ""
}

func (m *ACModel) IsSensitive(text string) bool { return m.FindOne(text) != "" }

func (m *ACModel) Replace(text string, repl rune) string {
	ranges := m.FindAllRanges(text)
	runes := []rune(text)
	// 标记区间（处理重叠：按起始位置排序后合并）
	marked := make([]bool, len(runes))
	for _, r := range ranges {
		for i := r.Start; i <= r.End && i < len(marked); i++ {
			marked[i] = true
		}
	}
	for i := range marked {
		if marked[i] {
			runes[i] = repl
		}
	}
	return string(runes)
}

func (m *ACModel) Remove(text string) string {
	ranges := m.FindAllRanges(text)
	runes := []rune(text)
	// 标记要删除的区间
	toDelete := make([]bool, len(runes))
	for _, r := range ranges {
		for i := r.Start; i <= r.End && i < len(toDelete); i++ {
			toDelete[i] = true
		}
	}
	// 构建结果
	result := make([]rune, 0, len(runes))
	for i, r := range runes {
		if !toDelete[i] {
			result = append(result, r)
		}
	}
	return string(result)
}

// FindAllRanges 返回所有匹配的区间，实现 RangedFilter 接口
func (m *ACModel) FindAllRanges(text string) []filter.Range {
	root := m.rootPtr.Load()
	var ranges []filter.Range
	seen := make(map[string]struct{})
	now := root
	runes := []rune(text)

	for i, r := range runes {
		for now != root && now.children[r] == nil {
			now = now.fail
		}
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			now = root
		}
		for _, w := range now.output {
			wordLen := len([]rune(w))
			start := i - wordLen + 1
			if start >= 0 {
				key := string(runes[start : i+1])
				if _, ok := seen[key]; !ok {
					seen[key] = struct{}{}
					ranges = append(ranges, filter.Range{Start: start, End: i})
				}
			}
		}
	}
	return ranges
}

package dfa

import "github.com/LuYongwang/go-sensitive-word/internal/filter"

type dfaNode struct {
	children map[rune]*dfaNode
	isLeaf   bool
}

func newDfaNode() *dfaNode {
	return &dfaNode{children: make(map[rune]*dfaNode), isLeaf: false}
}

type DFAModel struct {
	root *dfaNode
}

func NewDFAModel() *DFAModel {
	return &DFAModel{root: newDfaNode()}
}

func (m *DFAModel) AddWords(words ...string) {
	for _, word := range words {
		m.AddWord(word)
	}
}

func (m *DFAModel) AddWord(word string) {
	if word == "" {
		return
	}
	now := m.root
	for _, r := range word {
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			next = newDfaNode()
			now.children[r] = next
			now = next
		}
	}
	now.isLeaf = true
}

func (m *DFAModel) DelWords(words ...string) {
	for _, word := range words {
		m.DelWord(word)
	}
}

func (m *DFAModel) DelWord(word string) {
	if word == "" {
		return
	}
	runes := []rune(word)
	type pathElem struct {
		node *dfaNode
		ch   rune
	}
	path := make([]pathElem, 0, len(runes)+1)
	now := m.root
	path = append(path, pathElem{node: now})
	for _, r := range runes {
		next, ok := now.children[r]
		if !ok {
			return
		}
		path = append(path, pathElem{node: next, ch: r})
		now = next
	}
	if !now.isLeaf {
		return
	}
	now.isLeaf = false
	for i := len(path) - 1; i >= 1; i-- {
		curr := path[i].node
		parent := path[i-1].node
		ch := path[i].ch
		if len(curr.children) == 0 && !curr.isLeaf {
			delete(parent.children, ch)
		} else {
			break
		}
	}
}

func (m *DFAModel) Listen(addChan, delChan <-chan string) {
	go func() {
		for word := range addChan {
			m.AddWord(word)
		}
	}()
	go func() {
		for word := range delChan {
			m.DelWord(word)
		}
	}()
}

func (m *DFAModel) FindAll(text string) []string {
	var matches []string
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf && start <= pos {
			matches = append(matches, string(runes[start:pos+1]))
		}
		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}
		parent = now
	}
	var res []string
	set := make(map[string]struct{})
	for _, word := range matches {
		if _, ok := set[word]; !ok {
			set[word] = struct{}{}
			res = append(res, word)
		}
	}
	return res
}

func (m *DFAModel) FindAllCount(text string) map[string]int {
	res := make(map[string]int)
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf && start <= pos {
			res[string(runes[start:pos+1])]++
		}
		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}
		parent = now
	}
	return res
}

func (m *DFAModel) FindOne(text string) string {
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf && start <= pos {
			return string(runes[start : pos+1])
		}
		parent = now
	}
	return ""
}

func (m *DFAModel) IsSensitive(text string) bool { return m.FindOne(text) != "" }

// FindAllRanges 返回所有匹配的区间，实现 RangedFilter 接口
func (m *DFAModel) FindAllRanges(text string) []filter.Range {
	var ranges []filter.Range
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	seen := make(map[string]struct{}) // 用字符串作为 key 去重

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf && start <= pos {
			key := string(runes[start : pos+1])
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				ranges = append(ranges, filter.Range{Start: start, End: pos})
			}
		}
		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}
		parent = now
	}
	return ranges
}

func (m *DFAModel) Replace(text string, repl rune) string {
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf && start <= pos {
			for i := start; i <= pos; i++ {
				runes[i] = repl
			}
		}
		parent = now
	}
	return string(runes)
}

func (m *DFAModel) Remove(text string) string {
	var found bool
	var now *dfaNode
	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	filtered := make([]rune, 0, length)
	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]
		if !found || (!now.isLeaf && pos == length-1) {
			filtered = append(filtered, runes[start])
			parent = m.root
			pos = start
			start++
			continue
		}
		if now.isLeaf {
			start = pos + 1
			parent = m.root
		} else {
			parent = now
		}
	}
	filtered = append(filtered, runes[start:]...)
	return string(filtered)
}

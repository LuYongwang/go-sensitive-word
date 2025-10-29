package go_sensitive_word

import (
	"github.com/LuYongwang/go-sensitive-word/internal/filter"
)

// normalizedFilter 对底层 filter.Filter 做归一化包装：
// - 查询时：对文本与字典均做相同归一化
// - 返回时：基于匹配到的规范化片段在原文中定位，返回原文片段
type normalizedFilter struct {
	cfg   NormalizerConfig
	inner filter.Filter
}

func newNormalizedFilter(inner filter.Filter, cfg NormalizerConfig) *normalizedFilter {
	return &normalizedFilter{cfg: cfg, inner: inner}
}

func (nf *normalizedFilter) FindOne(text string) string {
	normText, idxMap := NormalizeTextWithMap(text, nf.cfg)
	hit := nf.inner.FindOne(normText)
	if hit == "" {
		return ""
	}
	// 在规范化文本中再次定位命中的首个区间
	// 简单做法：在 normText 中查找 hit 的首次位置
	// 注意：inner 的实现基于规范化文本，无需再次归一化
	start := -1
	end := -1
	rText := []rune(normText)
	rHit := []rune(hit)
	for i := 0; i+len(rHit) <= len(rText); i++ {
		ok := true
		for j := 0; j < len(rHit); j++ {
			if rText[i+j] != rHit[j] {
				ok = false
				break
			}
		}
		if ok {
			start = i
			end = i + len(rHit) - 1
			break
		}
	}
	if start == -1 {
		return ""
	}
	// 映射回原文区间
	origRunes := []rune(text)
	lo := idxMap[start]
	hi := idxMap[end]
	return string(origRunes[lo : hi+1])
}

func (nf *normalizedFilter) FindAll(text string) []string {
	normText, idxMap := NormalizeTextWithMap(text, nf.cfg)

	// 优先使用 FindAllRanges（如果支持）
	if rf, ok := nf.inner.(filter.RangedFilter); ok {
		ranges := rf.FindAllRanges(normText)
		if len(ranges) == 0 {
			return []string{}
		}
		origRunes := []rune(text)
		res := make([]string, 0, len(ranges))
		for _, r := range ranges {
			lo := idxMap[r.Start]
			hi := idxMap[r.End]
			if lo < len(origRunes) && hi < len(origRunes) {
				res = append(res, string(origRunes[lo:hi+1]))
			}
		}
		return res
	}

	// 降级到 FindAll（需要二次搜索）
	hits := nf.inner.FindAll(normText)
	if len(hits) == 0 {
		return hits
	}
	// 逐个命中在原文中定位
	res := make([]string, 0, len(hits))
	for _, h := range hits {
		rText := []rune(normText)
		rHit := []rune(h)
		start := -1
		end := -1
		for i := 0; i+len(rHit) <= len(rText); i++ {
			ok := true
			for j := 0; j < len(rHit); j++ {
				if rText[i+j] != rHit[j] {
					ok = false
					break
				}
			}
			if ok {
				start = i
				end = i + len(rHit) - 1
				break
			}
		}
		if start == -1 {
			continue
		}
		origRunes := []rune(text)
		lo := idxMap[start]
		hi := idxMap[end]
		res = append(res, string(origRunes[lo:hi+1]))
	}
	return res
}

func (nf *normalizedFilter) FindAllCount(text string) map[string]int {
	normText, idxMap := NormalizeTextWithMap(text, nf.cfg)
	hits := nf.inner.FindAll(normText)
	res := make(map[string]int, len(hits))
	// 简化：基于 FindAll 的唯一集合统计位置，再次在原文映射并计数
	rText := []rune(normText)
	for _, h := range hits {
		rHit := []rune(h)
		count := 0
		for i := 0; i+len(rHit) <= len(rText); i++ {
			ok := true
			for j := 0; j < len(rHit); j++ {
				if rText[i+j] != rHit[j] {
					ok = false
					break
				}
			}
			if ok {
				count++
				i += len(rHit) - 1
			}
		}
		if count == 0 {
			continue
		}
		// 返回原文片段作为 key：取首次出现位置
		// 定位首次
		start := -1
		end := -1
		for i := 0; i+len(rHit) <= len(rText); i++ {
			ok := true
			for j := 0; j < len(rHit); j++ {
				if rText[i+j] != rHit[j] {
					ok = false
					break
				}
			}
			if ok {
				start = i
				end = i + len(rHit) - 1
				break
			}
		}
		if start == -1 {
			continue
		}
		origRunes := []rune(text)
		lo := idxMap[start]
		hi := idxMap[end]
		res[string(origRunes[lo:hi+1])] = count
	}
	return res
}

func (nf *normalizedFilter) IsSensitive(text string) bool {
	return nf.FindOne(text) != ""
}

func (nf *normalizedFilter) Replace(text string, repl rune) string {
	normText, idxMap := NormalizeTextWithMap(text, nf.cfg)
	rOrig := []rune(text)
	marked := make([]bool, len(rOrig))

	// 优先使用 FindAllRanges（如果支持）
	if rf, ok := nf.inner.(filter.RangedFilter); ok {
		ranges := rf.FindAllRanges(normText)
		for _, r := range ranges {
			lo := idxMap[r.Start]
			hi := idxMap[r.End]
			for k := lo; k <= hi && k < len(marked); k++ {
				marked[k] = true
			}
		}
	} else {
		// 降级到 FindAll
		hits := nf.inner.FindAll(normText)
		rNorm := []rune(normText)
		for _, h := range hits {
			rHit := []rune(h)
			for i := 0; i+len(rHit) <= len(rNorm); i++ {
				ok := true
				for j := 0; j < len(rHit); j++ {
					if rNorm[i+j] != rHit[j] {
						ok = false
						break
					}
				}
				if ok {
					lo := idxMap[i]
					hi := idxMap[i+len(rHit)-1]
					for k := lo; k <= hi && k < len(marked); k++ {
						marked[k] = true
					}
					i += len(rHit) - 1
				}
			}
		}
	}

	for i := range marked {
		if marked[i] {
			rOrig[i] = repl
		}
	}
	return string(rOrig)
}

func (nf *normalizedFilter) Remove(text string) string {
	normText, idxMap := NormalizeTextWithMap(text, nf.cfg)
	rOrig := []rune(text)
	del := make([]bool, len(rOrig))

	// 优先使用 FindAllRanges（如果支持）
	if rf, ok := nf.inner.(filter.RangedFilter); ok {
		ranges := rf.FindAllRanges(normText)
		for _, r := range ranges {
			lo := idxMap[r.Start]
			hi := idxMap[r.End]
			for k := lo; k <= hi && k < len(del); k++ {
				del[k] = true
			}
		}
	} else {
		// 降级到 FindAll
		hits := nf.inner.FindAll(normText)
		rNorm := []rune(normText)
		for _, h := range hits {
			rHit := []rune(h)
			for i := 0; i+len(rHit) <= len(rNorm); i++ {
				ok := true
				for j := 0; j < len(rHit); j++ {
					if rNorm[i+j] != rHit[j] {
						ok = false
						break
					}
				}
				if ok {
					lo := idxMap[i]
					hi := idxMap[i+len(rHit)-1]
					for k := lo; k <= hi && k < len(del); k++ {
						del[k] = true
					}
					i += len(rHit) - 1
				}
			}
		}
	}

	out := make([]rune, 0, len(rOrig))
	for i, r := range rOrig {
		if !del[i] {
			out = append(out, r)
		}
	}
	return string(out)
}

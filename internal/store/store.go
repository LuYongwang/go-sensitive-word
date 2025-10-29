package store

import (
	"context"
	"io"
	"time"
)

// Stats 词库统计信息
type Stats struct {
	TotalWords  int       // 总词数
	LastUpdate  time.Time // 最后更新时间
	UpdateCount int       // 更新次数（添加+删除）
	Source      []string  // 词库来源（文件路径、URL等）
}

// DictLoader 词库加载回调函数类型
// 使用方可以实现自己的词库读取逻辑（如从数据库、Redis、配置中心等读取）
// 返回词列表和可能的错误
type DictLoader func() ([]string, error)

type (
	Store interface {
		// 加载词库
		LoadDictPath(path ...string) error
		LoadDictEmbed(contents ...string) error
		LoadDict(reader io.Reader) error
		LoadDictCallback(loader DictLoader, source string) error // 通过回调函数加载词库

		// 读取词库
		ReadChan() <-chan string
		ReadString() []string

		// 动态维护
		GetAddChan() <-chan string
		GetDelChan() <-chan string
		AddWord(words ...string) error
		DelWord(words ...string) error

		// 批量操作
		AddWords(words []string) error                  // 批量添加（新增方法，与 AddWord 功能相同但参数更明确）
		DelWords(words []string) error                  // 批量删除（新增方法）
		ReplaceWords(oldWords, newWords []string) error // 批量替换

		// 导出功能
		ExportToWriter(w io.Writer) error // 导出到 Writer
		ExportToString() (string, error)  // 导出为字符串（每行一个词）

		// 统计信息
		GetStats() Stats // 获取统计信息

		// 词库操作
		Clear() error            // 清空词库
		Merge(other Store) error // 合并另一个词库（去重）

		// 生命周期
		Close() error
		Shutdown(ctx context.Context) error
	}
)

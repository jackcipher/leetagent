//go:build !solution

package ex01_trace

import "time"

// Clock 抽象时间，测试用假时钟。
type Clock interface {
	Now() time.Time
}

// Span 是一次操作（一轮 LLM 调用、一次工具调用……）。
type Span struct {
	Name     string
	Start    time.Time
	EndTime  time.Time // 零值表示还没结束
	Attrs    map[string]string
	Children []*Span
	// TODO: 需要的话可以添加未导出字段（例如指向 Tracer 的指针）
}

// Tracer 记录一棵（或多棵）span 树。
type Tracer struct {
	Clock Clock
	Roots []*Span
}

// Start 开始一个 span；parent 为 nil 时成为根 span（追加到 Roots），否则追加到 parent.Children。
func (t *Tracer) Start(name string, parent *Span) *Span {
	// TODO
	return nil
}

// End 结束 span（记录 EndTime = Clock.Now()）。重复调用只有第一次生效。
func (s *Span) End() {
	// TODO
}

// SetAttr 设置属性（Attrs 为 nil 时创建）。
func (s *Span) SetAttr(key, value string) {
	// TODO
}

// Duration 返回 EndTime-Start；未结束返回 0。
func (s *Span) Duration() time.Duration {
	// TODO
	return 0
}

// Render 把所有根 span 渲染成文本树，每个 span 一行：
//
//	<缩进><name> <耗时>[ k=v ...]
//	缩进：每层 2 个空格；耗时：已结束为 "123ms"（毫秒取整），未结束为 "running"；
//	属性按 key 排序，每个前面一个空格。每行以 "\n" 结尾。
//
// 例：
//
//	chat 1500ms session=abc
//	  llm#1 800ms model=glm
//	  tool:get_post 200ms
func (t *Tracer) Render() string {
	// TODO
	return ""
}

// Find 深度优先（先序）查找第一个名字等于 name 的 span，找不到返回 nil。
func (t *Tracer) Find(name string) *Span {
	// TODO
	return nil
}

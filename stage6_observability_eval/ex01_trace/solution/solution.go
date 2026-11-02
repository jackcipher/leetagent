// Package solution 是 ex01_trace 的参考答案。
package solution

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Span struct {
	Name     string
	Start    time.Time
	EndTime  time.Time
	Attrs    map[string]string
	Children []*Span
	tracer   *Tracer
}

type Tracer struct {
	Clock Clock
	Roots []*Span
}

func (t *Tracer) Start(name string, parent *Span) *Span {
	s := &Span{Name: name, Start: t.Clock.Now(), tracer: t}
	if parent == nil {
		t.Roots = append(t.Roots, s)
	} else {
		parent.Children = append(parent.Children, s)
	}
	return s
}

func (s *Span) End() {
	if s.EndTime.IsZero() {
		s.EndTime = s.tracer.Clock.Now()
	}
}

func (s *Span) SetAttr(key, value string) {
	if s.Attrs == nil {
		s.Attrs = map[string]string{}
	}
	s.Attrs[key] = value
}

func (s *Span) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return 0
	}
	return s.EndTime.Sub(s.Start)
}

func (t *Tracer) Render() string {
	var b strings.Builder
	for _, r := range t.Roots {
		render(&b, r, 0)
	}
	return b.String()
}

func render(b *strings.Builder, s *Span, depth int) {
	b.WriteString(strings.Repeat("  ", depth))
	b.WriteString(s.Name)
	if s.EndTime.IsZero() {
		b.WriteString(" running")
	} else {
		fmt.Fprintf(b, " %dms", s.Duration().Milliseconds())
	}
	keys := make([]string, 0, len(s.Attrs))
	for k := range s.Attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(" " + k + "=" + s.Attrs[k])
	}
	b.WriteString("\n")
	for _, c := range s.Children {
		render(b, c, depth+1)
	}
}

func (t *Tracer) Find(name string) *Span {
	for _, r := range t.Roots {
		if s := find(r, name); s != nil {
			return s
		}
	}
	return nil
}

func find(s *Span, name string) *Span {
	if s.Name == name {
		return s
	}
	for _, c := range s.Children {
		if f := find(c, name); f != nil {
			return f
		}
	}
	return nil
}

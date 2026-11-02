// Package solution 是 m1_forum_tools 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"leetagent/common"
)

type Post struct {
	TID     int    `json:"tid"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Replies int    `json:"replies"`
	Content string `json:"content"`
}

type Summary struct {
	TID     int    `json:"tid"`
	Title   string `json:"title"`
	Replies int    `json:"replies"`
}

type Forum struct {
	Posts []Post
}

func LoadForum(path string) (*Forum, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var posts []Post
	if err := json.Unmarshal(b, &posts); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &Forum{Posts: posts}, nil
}

func (f *Forum) Defs() []common.Tool {
	return []common.Tool{
		common.NewFunctionTool("get_list", "按回复数从多到少列出热门帖子（只有标题，不含正文）",
			`{"type":"object","properties":{"limit":{"type":"integer","description":"最多返回几条，默认 10"}}}`),
		common.NewFunctionTool("get_post", "按帖子编号 tid 读取完整帖子正文",
			`{"type":"object","properties":{"tid":{"type":"integer","description":"帖子编号"}},"required":["tid"]}`),
		common.NewFunctionTool("search_posts", "按关键词搜索帖子标题和正文，返回匹配帖子的 tid 和标题",
			`{"type":"object","properties":{"query":{"type":"string","description":"关键词"}},"required":["query"]}`),
	}
}

type byHot []Summary

func (s byHot) Len() int      { return len(s) }
func (s byHot) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byHot) Less(i, j int) bool {
	if s[i].Replies != s[j].Replies {
		return s[i].Replies > s[j].Replies
	}
	return s[i].TID < s[j].TID
}

type byTID []Summary

func (s byTID) Len() int           { return len(s) }
func (s byTID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byTID) Less(i, j int) bool { return s[i].TID < s[j].TID }

func summarize(p Post) Summary { return Summary{TID: p.TID, Title: p.Title, Replies: p.Replies} }

func toJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

type args struct {
	Limit int    `json:"limit"`
	TID   int    `json:"tid"`
	Query string `json:"query"`
}

func (f *Forum) Call(ctx context.Context, name, argsJSON string) (string, error) {
	var a args
	if strings.TrimSpace(argsJSON) == "" {
		argsJSON = "{}"
	}
	if err := json.Unmarshal([]byte(argsJSON), &a); err != nil {
		return "", fmt.Errorf("bad arguments: %v", err)
	}
	switch name {
	case "get_list":
		if a.Limit <= 0 {
			a.Limit = 10
		}
		list := make([]Summary, 0, len(f.Posts))
		for _, p := range f.Posts {
			list = append(list, summarize(p))
		}
		sort.Sort(byHot(list))
		if len(list) > a.Limit {
			list = list[:a.Limit]
		}
		return toJSON(list)
	case "get_post":
		for _, p := range f.Posts {
			if p.TID == a.TID {
				return toJSON(p)
			}
		}
		return "", fmt.Errorf("post %d not found", a.TID)
	case "search_posts":
		q := strings.ToLower(strings.TrimSpace(a.Query))
		if q == "" {
			return "", fmt.Errorf("query is required")
		}
		list := []Summary{}
		for _, p := range f.Posts {
			if strings.Contains(strings.ToLower(p.Title), q) || strings.Contains(strings.ToLower(p.Content), q) {
				list = append(list, summarize(p))
			}
		}
		sort.Sort(byTID(list))
		return toJSON(list)
	}
	return "", fmt.Errorf("unknown tool: %s", name)
}

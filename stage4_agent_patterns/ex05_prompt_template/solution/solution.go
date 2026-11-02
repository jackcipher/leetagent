// Package solution 是 ex05_prompt_template 的参考答案。
package solution

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"leetagent/common"
)

type Section struct {
	Title string
	Body  string
}

func Build(sections []Section) string {
	var parts []string
	for _, s := range sections {
		body := strings.TrimSpace(s.Body)
		if body == "" {
			continue
		}
		parts = append(parts, "## "+s.Title+"\n"+body)
	}
	return strings.Join(parts, "\n\n")
}

type byName []common.Tool

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].Function.Name < b[j].Function.Name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func RenderToolGuide(tools []common.Tool) string {
	sorted := append(byName(nil), tools...)
	sort.Sort(sorted)
	lines := make([]string, len(sorted))
	for i, t := range sorted {
		lines[i] = "- " + t.Function.Name + ": " + t.Function.Description
	}
	return strings.Join(lines, "\n")
}

var snake = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type paramSchema struct {
	Properties map[string]struct {
		Description string `json:"description"`
	} `json:"properties"`
}

func LintToolDescription(t common.Tool) []string {
	var issues []string
	if !snake.MatchString(t.Function.Name) {
		issues = append(issues, "name: must be snake_case")
	}
	n := utf8.RuneCountInString(t.Function.Description)
	if n < 10 {
		issues = append(issues, "description: too short")
	}
	if n > 200 {
		issues = append(issues, "description: too long")
	}
	var ps paramSchema
	if err := json.Unmarshal(t.Function.Parameters, &ps); err != nil {
		return append(issues, "parameters: invalid json")
	}
	var names []string
	for k, p := range ps.Properties {
		if strings.TrimSpace(p.Description) == "" {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	for _, k := range names {
		issues = append(issues, "parameters."+k+": missing description")
	}
	return issues
}

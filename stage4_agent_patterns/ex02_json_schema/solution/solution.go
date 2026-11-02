// Package solution 是 ex02_json_schema 的参考答案。
package solution

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"leetagent/common"
)

var ErrInvalidAfterRetries = errors.New("structured output still invalid after retries")

type schemaNode struct {
	Type                 string                 `json:"type"`
	Properties           map[string]*schemaNode `json:"properties"`
	Required             []string               `json:"required"`
	AdditionalProperties *bool                  `json:"additionalProperties"`
	Items                *schemaNode            `json:"items"`
	Enum                 []any                  `json:"enum"`
	Minimum              *float64               `json:"minimum"`
	Maximum              *float64               `json:"maximum"`
	MinLength            *int                   `json:"minLength"`
	MaxLength            *int                   `json:"maxLength"`
}

func decode(raw []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	return d.Decode(v)
}

func Validate(schema, data json.RawMessage) []string {
	var s schemaNode
	var v any
	if err := json.Unmarshal(schema, &s); err != nil {
		return []string{"$: invalid json"}
	}
	if err := decode(data, &v); err != nil {
		return []string{"$: invalid json"}
	}
	var errs []string
	errs = check(&s, v, "$", errs)
	sort.Strings(errs)
	return errs
}

func typeOK(t string, v any) bool {
	switch t {
	case "":
		return true
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "null":
		return v == nil
	case "number":
		_, ok := v.(json.Number)
		return ok
	case "integer":
		n, ok := v.(json.Number)
		if !ok {
			return false
		}
		f, err := n.Float64()
		return err == nil && f == math.Trunc(f)
	}
	return false
}

func normalize(v any) any {
	if n, ok := v.(json.Number); ok {
		f, _ := n.Float64()
		return f
	}
	return v
}

func inEnum(enum []any, v any) bool {
	nv := normalize(v)
	for _, e := range enum {
		if f, ok := e.(float64); ok {
			if g, ok := nv.(float64); ok && f == g {
				return true
			}
			continue
		}
		if reflect.DeepEqual(e, nv) {
			return true
		}
	}
	return false
}

func check(s *schemaNode, v any, path string, errs []string) []string {
	if !typeOK(s.Type, v) {
		return append(errs, fmt.Sprintf("%s: expected %s", path, s.Type))
	}
	if len(s.Enum) > 0 && !inEnum(s.Enum, v) {
		errs = append(errs, fmt.Sprintf("%s: not in enum", path))
	}
	switch x := v.(type) {
	case map[string]any:
		errs = checkObject(s, x, path, errs)
	case []any:
		if s.Items != nil {
			for i, item := range x {
				errs = check(s.Items, item, fmt.Sprintf("%s[%d]", path, i), errs)
			}
		}
	case string:
		n := utf8.RuneCountInString(x)
		if s.MinLength != nil && n < *s.MinLength {
			errs = append(errs, fmt.Sprintf("%s: shorter than %d", path, *s.MinLength))
		}
		if s.MaxLength != nil && n > *s.MaxLength {
			errs = append(errs, fmt.Sprintf("%s: longer than %d", path, *s.MaxLength))
		}
	case json.Number:
		f, _ := x.Float64()
		if s.Minimum != nil && f < *s.Minimum {
			errs = append(errs, fmt.Sprintf("%s: less than minimum %v", path, *s.Minimum))
		}
		if s.Maximum != nil && f > *s.Maximum {
			errs = append(errs, fmt.Sprintf("%s: greater than maximum %v", path, *s.Maximum))
		}
	}
	return errs
}

func checkObject(s *schemaNode, obj map[string]any, path string, errs []string) []string {
	for _, r := range s.Required {
		if _, ok := obj[r]; !ok {
			errs = append(errs, fmt.Sprintf("%s.%s: required", path, r))
		}
	}
	for k, val := range obj {
		sub, ok := s.Properties[k]
		if !ok {
			if s.AdditionalProperties != nil && !*s.AdditionalProperties {
				errs = append(errs, fmt.Sprintf("%s.%s: additional property not allowed", path, k))
			}
			continue
		}
		errs = check(sub, val, path+"."+k, errs)
	}
	return errs
}

func ExtractJSON(output string) string {
	s := strings.TrimSpace(output)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		if i := strings.Index(s, "\n"); i >= 0 {
			s = s[i+1:]
		}
		s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "```"))
	}
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return s
	}
	i, j := strings.Index(s, "{"), strings.LastIndex(s, "}")
	if i >= 0 && j > i {
		return s[i : j+1]
	}
	return s
}

func AskStructured(ctx context.Context, c common.Chatter, model string, msgs []common.Message, schema json.RawMessage, retries int) (json.RawMessage, error) {
	history := append([]common.Message(nil), msgs...)
	for attempt := 0; attempt <= retries; attempt++ {
		reply, err := c.Complete(ctx, model, history, nil)
		if err != nil {
			return nil, err
		}
		js := ExtractJSON(reply.Content)
		errs := Validate(schema, json.RawMessage(js))
		if len(errs) == 0 {
			return json.RawMessage(js), nil
		}
		history = append(history,
			common.Message{Role: "assistant", Content: reply.Content},
			common.Message{Role: "user", Content: "你的输出不符合要求：\n" + strings.Join(errs, "\n") + "\n请只输出修正后的 JSON。"},
		)
	}
	return nil, ErrInvalidAfterRetries
}

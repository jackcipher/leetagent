package ex05_slog

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func decodeLines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("not json: %q", line)
		}
		out = append(out, m)
	}
	return out
}

func TestLogRound(t *testing.T) {
	var buf bytes.Buffer
	l := NewJSONLogger(&buf, slog.LevelInfo)
	if l == nil {
		t.Fatal("nil logger")
	}
	LogRound(l, RoundStats{RequestID: "r1", Round: 2, Model: "glm", Duration: 1500 * time.Millisecond, Tools: []string{"get_list", "fetch_one"}})
	LogRound(l, RoundStats{RequestID: "r1", Round: 3, Model: "glm", Duration: time.Millisecond, Err: errors.New("EOF")})
	lines := decodeLines(t, &buf)
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}
	a := lines[0]
	if a["msg"] != "agent round" || a["level"] != "INFO" || a["model"] != "glm" || a["tools"] != "get_list,fetch_one" ||
		a["duration_ms"] != float64(1500) || a["tool_count"] != float64(2) || a["round"] != float64(2) || a["request_id"] != "r1" {
		t.Fatalf("bad first line: %v", a)
	}
	if _, ok := a["err"]; ok {
		t.Fatal("info line should not have err")
	}
	b := lines[1]
	if b["level"] != "ERROR" || b["err"] != "EOF" || b["tool_count"] != float64(0) {
		t.Fatalf("bad second line: %v", b)
	}
}

func TestLevelFilter(t *testing.T) {
	var buf bytes.Buffer
	l := NewJSONLogger(&buf, slog.LevelError)
	LogRound(l, RoundStats{Round: 1})
	if buf.Len() != 0 {
		t.Fatal("info should be filtered at error level")
	}
}

func TestWithRequest(t *testing.T) {
	var buf bytes.Buffer
	l := WithRequest(NewJSONLogger(&buf, slog.LevelInfo), "req-9", "jack")
	l.Info("hello")
	lines := decodeLines(t, &buf)
	if len(lines) != 1 || lines[0]["request_id"] != "req-9" || lines[0]["user"] != "jack" {
		t.Fatalf("got %v", lines)
	}
}

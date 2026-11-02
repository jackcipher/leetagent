//go:build !solution

package ex01_session_store

import (
	"context"
	"database/sql"
	"errors"

	"leetagent/common"

	_ "modernc.org/sqlite" // 纯 Go 的 SQLite 驱动，驱动名 "sqlite"
)

var ErrSessionNotFound = errors.New("session not found")

// Store 用 SQLite 存会话和消息。
type Store struct {
	db *sql.DB
}

// Open 打开（或创建）数据库文件并建表：
//
//	sessions(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, created_at INTEGER)
//	messages(id INTEGER PRIMARY KEY AUTOINCREMENT, session_id INTEGER, role TEXT, content TEXT,
//	         tool_calls TEXT, tool_call_id TEXT)
func Open(path string) (*Store, error) {
	// TODO
	return nil, errors.New("TODO")
}

// Close 关闭数据库。
func (s *Store) Close() error {
	// TODO
	return nil
}

// CreateSession 新建会话，返回 id。
func (s *Store) CreateSession(ctx context.Context, title string) (int64, error) {
	// TODO
	return 0, nil
}

// Append 追加一条消息。ToolCalls 序列化成 JSON 存进 tool_calls 列。会话不存在返回 ErrSessionNotFound。
func (s *Store) Append(ctx context.Context, sessionID int64, m common.Message) error {
	// TODO
	return nil
}

// History 按插入顺序返回会话全部消息。会话不存在返回 ErrSessionNotFound。
func (s *Store) History(ctx context.Context, sessionID int64) ([]common.Message, error) {
	// TODO
	return nil, nil
}

// Package solution 是 ex01_session_store 的参考答案。
package solution

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"leetagent/common"

	_ "modernc.org/sqlite"
)

var ErrSessionNotFound = errors.New("session not found")

type Store struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL DEFAULT '',
  created_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL REFERENCES sessions(id),
  role TEXT NOT NULL,
  content TEXT NOT NULL DEFAULT '',
  tool_calls TEXT NOT NULL DEFAULT '',
  tool_call_id TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id, id);
`

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite 单写者，避免 database is locked
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) CreateSession(ctx context.Context, title string) (int64, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO sessions(title, created_at) VALUES(?, ?)`, title, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) exists(ctx context.Context, id int64) error {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = ?`, id).Scan(&n)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (s *Store) Append(ctx context.Context, sessionID int64, m common.Message) error {
	if err := s.exists(ctx, sessionID); err != nil {
		return err
	}
	calls := ""
	if len(m.ToolCalls) > 0 {
		b, err := json.Marshal(m.ToolCalls)
		if err != nil {
			return err
		}
		calls = string(b)
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO messages(session_id, role, content, tool_calls, tool_call_id) VALUES(?, ?, ?, ?, ?)`,
		sessionID, m.Role, m.Content, calls, m.ToolCallID)
	return err
}

func (s *Store) History(ctx context.Context, sessionID int64) ([]common.Message, error) {
	if err := s.exists(ctx, sessionID); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT role, content, tool_calls, tool_call_id FROM messages WHERE session_id = ? ORDER BY id`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []common.Message
	for rows.Next() {
		var m common.Message
		var calls string
		if err := rows.Scan(&m.Role, &m.Content, &calls, &m.ToolCallID); err != nil {
			return nil, err
		}
		if calls != "" {
			if err := json.Unmarshal([]byte(calls), &m.ToolCalls); err != nil {
				return nil, err
			}
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

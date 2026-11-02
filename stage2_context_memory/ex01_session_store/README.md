# ex01 SQLite 会话持久化（Medium）

## 题目
现在 你的 Agent 项目 每次请求都由前端带上完整 `messages[]`：刷新页面就丢、换设备就丢、前端还能伪造历史。
改成服务端存：前端只传 `session_id` + 新消息，后端从 SQLite 取历史、追加、再存。

## 接口
```go
func Open(path string) (*Store, error)
func (s *Store) Close() error
func (s *Store) CreateSession(ctx context.Context, title string) (int64, error)
func (s *Store) Append(ctx context.Context, sessionID int64, m common.Message) error
func (s *Store) History(ctx context.Context, sessionID int64) ([]common.Message, error)
```

## 要求
- 用 `modernc.org/sqlite`（纯 Go，无需 CGO，已 vendor），驱动名 `"sqlite"`。
- assistant 的 `ToolCalls` 序列化为 JSON 存一列；读出来再反序列化，`reflect.DeepEqual` 要和写入时完全一致（注意 nil vs 空切片）。
- 会话不存在时 Append / History 返回 `ErrSessionNotFound`。
- 历史严格按插入顺序（`ORDER BY id`）。

## 验收
`go test ./stage2_context_memory/ex01_session_store/`

## 坑
- SQLite 同一时间只允许一个写者，`db.SetMaxOpenConns(1)` 能避免 `database is locked`。
- `:memory:` 数据库每个连接是独立的库，连接池会让你"表不见了"，所以测试用临时文件。
- 工具消息（role=tool）也要存，否则下一轮发给模型时 tool_call_id 对不上，上游会 400。

# 提示

1. `sql.Open("sqlite", path)` 后 `db.Exec(schemaSQL)` 建表，用 `CREATE TABLE IF NOT EXISTS`，可以一次执行多条语句。
2. 列都设 `NOT NULL DEFAULT ''`，tool_calls 为空时存空串；读出时 `if calls != "" { json.Unmarshal }`，这样 ToolCalls 保持 nil，DeepEqual 才相等。
3. 判断会话存在：`SELECT COUNT(*) FROM sessions WHERE id = ?`，Scan 到 int，为 0 返回 ErrSessionNotFound。

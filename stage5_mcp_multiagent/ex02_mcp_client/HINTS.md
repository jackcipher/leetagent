# 提示

1. 字段：`mu sync.Mutex; r *bufio.Reader; w io.Writer; nextID int`。写一个通用方法 `rpc(ctx, method string, params any, out any) error`：自增 id、写一行、循环 ReadString('\n') 直到 id 匹配。
2. 请求结构体里 `ID *int \`json:"id,omitempty"\``——通知时传 nil 就不会输出 id。比较 id：`string(resp.ID) == strconv.Itoa(id)`。
3. CallTool 解析 `{content:[{type,text}], isError}`；只拼 type=="text"；isError 时返回 `&ToolError{Text: ...}`。

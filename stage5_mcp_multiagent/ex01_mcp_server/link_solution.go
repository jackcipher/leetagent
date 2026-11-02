//go:build solution

package ex01_mcp_server

import "leetagent/stage5_mcp_multiagent/ex01_mcp_server/solution"

const (
	ProtocolVersion    = solution.ProtocolVersion
	CodeParseError     = solution.CodeParseError
	CodeMethodNotFound = solution.CodeMethodNotFound
	CodeInvalidParams  = solution.CodeInvalidParams
)

type (
	Request  = solution.Request
	RPCError = solution.RPCError
	Response = solution.Response
	Handler  = solution.Handler
	Server   = solution.Server
)

//go:build solution

package ex02_mcp_client

import "leetagent/stage5_mcp_multiagent/ex02_mcp_client/solution"

type (
	MCPTool   = solution.MCPTool
	RPCError  = solution.RPCError
	ToolError = solution.ToolError
	Client    = solution.Client
)

var (
	NewClient     = solution.NewClient
	ToCommonTools = solution.ToCommonTools
)

//go:build solution

package ex03_sse_heartbeat_fallback

import "leetagent/stage7_production/ex03_sse_heartbeat_fallback/solution"

type (
	Heartbeat = solution.Heartbeat
	Locker    = solution.Locker
	TokenSink = solution.TokenSink
	Streamer  = solution.Streamer
)

var (
	WriteSSE           = solution.WriteSSE
	StreamWithFallback = solution.StreamWithFallback
)

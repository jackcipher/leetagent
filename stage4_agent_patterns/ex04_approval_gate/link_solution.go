//go:build solution

package ex04_approval_gate

import "leetagent/stage4_agent_patterns/ex04_approval_gate/solution"

type (
	Approver     = solution.Approver
	Gate         = solution.Gate
	Pending      = solution.Pending
	ChanApprover = solution.ChanApprover
)

var (
	ErrDenied       = solution.ErrDenied
	NewChanApprover = solution.NewChanApprover
)

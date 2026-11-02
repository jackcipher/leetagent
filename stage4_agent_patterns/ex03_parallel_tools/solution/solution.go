// Package solution 是 ex03_parallel_tools 的参考答案。
package solution

import (
	"context"
	"sync"

	"leetagent/common"
)

type Result struct {
	CallID  string
	Content string
	Err     error
}

type job struct {
	ctx    context.Context
	runner common.ToolRunner
	call   common.ToolCall
	out    *Result
	sem    chan struct{}
	wg     *sync.WaitGroup
}

func failed(id string, err error) Result {
	return Result{CallID: id, Content: "error: " + err.Error(), Err: err}
}

func (j job) run() {
	defer j.wg.Done()
	defer releaseSlot(j.sem)
	if err := j.ctx.Err(); err != nil {
		*j.out = failed(j.call.ID, err)
		return
	}
	content, err := j.runner.Call(j.ctx, j.call.Function.Name, j.call.Function.Arguments)
	if err != nil {
		*j.out = failed(j.call.ID, err)
		return
	}
	*j.out = Result{CallID: j.call.ID, Content: content}
}

func releaseSlot(sem chan struct{}) { <-sem }

func ExecuteAll(ctx context.Context, runner common.ToolRunner, calls []common.ToolCall, maxConcurrency int) []Result {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}
	results := make([]Result, len(calls))
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for i, c := range calls {
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			results[i] = failed(c.ID, ctx.Err())
			continue
		}
		if err := ctx.Err(); err != nil {
			<-sem
			results[i] = failed(c.ID, err)
			continue
		}
		wg.Add(1)
		go job{ctx: ctx, runner: runner, call: c, out: &results[i], sem: sem, wg: &wg}.run()
	}
	wg.Wait()
	return results
}

func ToMessages(results []Result) []common.Message {
	msgs := make([]common.Message, len(results))
	for i, r := range results {
		msgs[i] = common.Message{Role: "tool", ToolCallID: r.CallID, Content: r.Content}
	}
	return msgs
}

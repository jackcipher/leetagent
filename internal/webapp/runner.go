package webapp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

const testTimeout = 30 * time.Second

type TestResult struct {
	Passed     bool   `json:"passed"`
	TimedOut   bool   `json:"timedOut"`
	ExitCode   int    `json:"exitCode"`
	DurationMS int64  `json:"durationMs"`
	Output     string `json:"output"`
}

type Runner interface {
	Run(ctx context.Context, problem Problem) TestResult
}

type GoTestRunner struct {
	Root string
}

func (r GoTestRunner) Run(ctx context.Context, problem Problem) TestResult {
	started := time.Now()
	testCtx, cancel := context.WithTimeout(ctx, testTimeout)
	defer cancel()

	rel, err := filepath.Rel(r.Root, problem.Dir)
	if err != nil {
		return TestResult{ExitCode: -1, Output: fmt.Sprintf("resolve exercise path: %v\n", err)}
	}
	cmd := exec.CommandContext(testCtx, "go", "test", "-json", "-timeout", testTimeout.String(), "./"+filepath.ToSlash(rel))
	cmd.Dir = r.Root
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) && testCtx.Err() == nil {
			fmt.Fprintf(&output, "leetagent: run go test: %v\n", err)
		}
	}

	result := TestResult{
		Passed:     err == nil,
		ExitCode:   exitCode(err),
		DurationMS: time.Since(started).Milliseconds(),
		Output:     output.String(),
	}
	if errors.Is(testCtx.Err(), context.DeadlineExceeded) {
		result.TimedOut = true
		result.Passed = false
		result.Output += "\nleetagent: test run exceeded 30 seconds\n"
	}
	return result
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

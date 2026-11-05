# Web Platform Implementation Task

## Goal

Turn the existing Go exercise repository into a LeetCode-like local web experience without creating a second source of truth. Exercise Markdown, `exercise.go`, and `exercise_test.go` remain authoritative; the web application discovers and presents them.

## Product Principles

- Keep the clone-and-IDE workflow fully supported.
- Make the first useful experience available with one command and no frontend toolchain.
- Use `go test` as the judge so browser and terminal results agree.
- Keep progress in the browser; this is a local learning tool, not a hosted service.
- Prefer a focused Go exercise editor over IDE features such as completion or project-wide navigation.

## Phase 1: Local Vertical Slice

Status: first vertical slice complete on 2026-11-04

- [x] Add `leetagent serve` command.
- [x] Discover exercises from the existing stage and capstone directories.
- [x] Expose problem statement, hints, and editable `exercise.go` through a local API.
- [x] Run the selected exercise with `go test -json` and a fixed timeout.
- [x] Build an embedded, dependency-free browser UI with problem navigation, Markdown presentation, source editor, and test console.
- [x] Expose `leetagent/common` as an AST-generated API reference plus browsable read-only source files.
- [x] Add resizable workspace and test panels with locally persisted layout.
- [x] Add Go source formatting without saving invalid or partial code.
- [x] Present `go test -json` as test cases, diagnostics, durations, and optional raw output.
- [x] Persist pass/fail navigation markers in browser storage.
- [x] Add backend tests for catalog, API, source updates, and real test execution.
- [x] Document the local web workflow in the root README.
- [x] Complete browser-level interaction checks for navigation and test execution.
- [x] Verify desktop and mobile layouts in a real browser.
- [x] Apply and browser-verify a light, LeetCode-inspired visual theme across problems, editor, tests, and Common API.
- [x] Add locally bundled, scroll-synchronized Go syntax highlighting without changing editor save or test behavior.
- [x] Capture the real browser workflow with ego-browser and embed an optimized demo GIF in the root README.

### Local API

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/problems` | List discovered exercises grouped by stage metadata. |
| `GET` | `/api/problem?id=<id>` | Load statement, hints, source, and exercise metadata. |
| `GET` | `/api/reference` | Load exported `leetagent/common` declarations and non-test source files. |
| `POST` | `/api/format` | Format submitted Go source without writing it to the repository. |
| `PUT` | `/api/source?id=<id>` | Save only the selected exercise's `exercise.go`. |
| `POST` | `/api/test?id=<id>` | Save the submitted source and run that exercise's tests. |

Problem IDs are selected from the server-side catalog. Client-provided filesystem paths are never accepted.

### Local Execution

1. The browser submits the current editor contents.
2. The server validates the problem ID and source size.
3. The server atomically updates that problem's `exercise.go`.
4. The runner executes `go test -json -timeout 30s ./<problem path>` from the repository root.
5. The API returns exit status, duration, and complete command output.

The local runner intentionally executes with the current user's permissions because it is a development tool for a trusted checkout. This runner must not be exposed directly to the internet.

## Explicitly Out of Scope

The local experience is the product. The following are intentionally not planned:

- A CodeMirror/Monaco-style IDE with completion, project navigation, or other heavyweight editor features.
- Server-side progress, submission history, accounts, or cross-device synchronization.
- A hosted judge, sandbox worker fleet, hidden tests, job queues, or public arbitrary-code execution.
- Profiles, discussions, leaderboards, streaks, or other SaaS/community features.

## Acceptance Criteria For Phase 1

- `go run ./cmd/leetagent serve` serves the application on localhost.
- The UI lists every directory containing `exercise.go` and `exercise_test.go` in the supported course roots.
- Selecting a problem loads its current statement, hints, and source.
- The Common API tab shows the exact shared package available to exercise code.
- Saving changes updates only that problem's `exercise.go`.
- Running tests uses the same package and tests as the terminal workflow.
- Invalid IDs, oversized source, command timeout, and failed tests return useful errors without crashing the server.
- The frontend assets are embedded in the Go binary and require no CDN or Node.js installation.

package webapp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	called string
	result TestResult
}

func (f *fakeRunner) Run(_ context.Context, problem Problem) TestResult {
	f.called = problem.ID
	return f.result
}

func TestServerProblemWorkflow(t *testing.T) {
	root := newTestRepository(t)
	server, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{result: TestResult{Passed: true, Output: "ok\n"}}
	server.runner = runner
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	var listing struct {
		Problems []Problem `json:"problems"`
	}
	getJSON(t, httpServer.URL+"/api/problems", &listing)
	if len(listing.Problems) != 2 {
		t.Fatalf("problem count = %d", len(listing.Problems))
	}

	id := "stage1_foundation/ex01_registry"
	var detail problemDetail
	getJSON(t, httpServer.URL+"/api/problem?id="+id, &detail)
	if !strings.Contains(detail.Statement, "Implement it") || !strings.Contains(detail.Source, "return 42") {
		t.Fatalf("unexpected detail: %+v", detail)
	}

	newSource := "package registry\n\nfunc Answer() int { return 7 }\n"
	result := postSource(t, http.MethodPost, httpServer.URL+"/api/test?id="+id, newSource)
	if result.StatusCode != http.StatusOK {
		defer result.Body.Close()
		body, _ := io.ReadAll(result.Body)
		t.Fatalf("test status = %d: %s", result.StatusCode, body)
	}
	result.Body.Close()
	if runner.called != id {
		t.Fatalf("runner called with %q", runner.called)
	}
	saved, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(id), "exercise.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != newSource {
		t.Fatalf("saved source = %q", saved)
	}
}

func TestServerCommonReference(t *testing.T) {
	root := newTestRepository(t)
	writeTestFile(t, root, "common/types.go", `// Package common contains shared course types.
package common

// Message is a chat message.
type Message struct { Content string }

// NewMessage constructs a Message.
func NewMessage(content string) Message { return Message{Content: content} }
`)
	writeTestFile(t, root, "common/types_test.go", "package common\n")
	server, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	var reference PackageReference
	getJSON(t, httpServer.URL+"/api/reference", &reference)
	if reference.ImportPath != "leetagent/common" || len(reference.Files) != 1 {
		t.Fatalf("reference = %+v", reference)
	}
	if len(reference.Symbols) != 2 || reference.Symbols[0].Name != "Message" || reference.Symbols[1].Name != "NewMessage" {
		t.Fatalf("symbols = %+v", reference.Symbols)
	}
	if !strings.Contains(reference.Symbols[0].Signature, "type Message struct") || reference.Symbols[0].Doc != "Message is a chat message." {
		t.Fatalf("message symbol = %+v", reference.Symbols[0])
	}
}

func TestServerFormatsSource(t *testing.T) {
	server, err := New(newTestRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	response := postSource(t, http.MethodPost, httpServer.URL+"/api/format", "package demo\nfunc Answer( )int{return 42}\n")
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("format status = %d: %s", response.StatusCode, body)
	}
	var result sourceRequest
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Source != "package demo\n\nfunc Answer() int { return 42 }\n" {
		t.Fatalf("formatted source = %q", result.Source)
	}
}

func TestServerRejectsInvalidGoForFormat(t *testing.T) {
	server, err := New(newTestRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	response := postSource(t, http.MethodPost, httpServer.URL+"/api/format", "package demo\nfunc")
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("format status = %d", response.StatusCode)
	}
}

func TestServerRejectsInvalidWrites(t *testing.T) {
	server, err := New(newTestRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	response := postSource(t, http.MethodPut, httpServer.URL+"/api/source?id=../../etc", "package bad\n")
	if response.StatusCode != http.StatusNotFound {
		response.Body.Close()
		t.Fatalf("invalid id status = %d", response.StatusCode)
	}
	response.Body.Close()

	request, err := http.NewRequest(http.MethodPut, httpServer.URL+"/api/source?id=stage1_foundation/ex01_registry", strings.NewReader(`{"source":"package registry"}`))
	if err != nil {
		t.Fatal(err)
	}
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("content type status = %d", response.StatusCode)
	}
}

func TestServerRejectsOversizedSource(t *testing.T) {
	server, err := New(newTestRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	id := "stage1_foundation/ex01_registry"
	response := postSource(t, http.MethodPut, httpServer.URL+"/api/source?id="+id, strings.Repeat("x", maxSourceBytes))
	defer response.Body.Close()
	if response.StatusCode != http.StatusRequestEntityTooLarge {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("oversized source status = %d: %s", response.StatusCode, body)
	}
}

func TestGoTestRunner(t *testing.T) {
	root := newTestRepository(t)
	problem := Problem{ID: "stage1_foundation/ex01_registry", Dir: filepath.Join(root, "stage1_foundation", "ex01_registry")}
	result := (GoTestRunner{Root: root}).Run(context.Background(), problem)
	if !result.Passed || result.ExitCode != 0 || !strings.Contains(result.Output, `"Action":"pass"`) {
		t.Fatalf("result = %+v", result)
	}
}

func getJSON(t *testing.T, url string, target any) {
	t.Helper()
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d", url, response.StatusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatal(err)
	}
}

func postSource(t *testing.T, method, url, source string) *http.Response {
	t.Helper()
	body, err := json.Marshal(sourceRequest{Source: source})
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

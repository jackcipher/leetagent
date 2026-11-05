package webapp

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	ReadHeaderTimeout = 5 * time.Second
	IdleTimeout       = 60 * time.Second
	maxSourceBytes    = 1 << 20
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	root    string
	catalog *Catalog
	runner  Runner
	static  http.Handler
	locks   map[string]*sync.Mutex
}

type problemDetail struct {
	Problem
	Statement string `json:"statement"`
	Hints     string `json:"hints"`
	Source    string `json:"source"`
}

type sourceRequest struct {
	Source string `json:"source"`
}

func New(root string) (*Server, error) {
	catalog, err := Discover(root)
	if err != nil {
		return nil, err
	}
	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("load frontend assets: %w", err)
	}
	locks := make(map[string]*sync.Mutex, len(catalog.problems))
	for _, problem := range catalog.problems {
		locks[problem.ID] = &sync.Mutex{}
	}
	return &Server{
		root:    root,
		catalog: catalog,
		runner:  GoTestRunner{Root: root},
		static:  http.FileServer(http.FS(assets)),
		locks:   locks,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self'; img-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	switch r.URL.Path {
	case "/api/problems":
		s.handleProblems(w, r)
	case "/api/problem":
		s.handleProblem(w, r)
	case "/api/reference":
		s.handleReference(w, r)
	case "/api/format":
		s.handleFormat(w, r)
	case "/api/source":
		s.handleSource(w, r)
	case "/api/test":
		s.handleTest(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/api/") {
			writeError(w, http.StatusNotFound, "API endpoint not found")
			return
		}
		s.static.ServeHTTP(w, r)
	}
}

func (s *Server) handleProblems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"problems": s.catalog.Problems()})
}

func (s *Server) handleProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	problem, ok := s.problemFromRequest(w, r)
	if !ok {
		return
	}
	statement, err := readOptionalFile(filepath.Join(problem.Dir, "README.md"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read problem statement")
		return
	}
	hints, err := readOptionalFile(filepath.Join(problem.Dir, "HINTS.md"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read hints")
		return
	}
	source, err := os.ReadFile(filepath.Join(problem.Dir, "exercise.go"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read exercise source")
		return
	}
	writeJSON(w, http.StatusOK, problemDetail{
		Problem:   problem,
		Statement: statement,
		Hints:     hints,
		Source:    string(source),
	})
}

func (s *Server) handleReference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	reference, err := loadCommonReference(s.root)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read common package reference")
		return
	}
	writeJSON(w, http.StatusOK, reference)
}

func (s *Server) handleFormat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	request, ok := decodeSource(w, r)
	if !ok {
		return
	}
	formatted, err := format.Source([]byte(request.Source))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sourceRequest{Source: string(formatted)})
}

func (s *Server) handleSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		methodNotAllowed(w, http.MethodPut)
		return
	}
	problem, ok := s.problemFromRequest(w, r)
	if !ok {
		return
	}
	request, ok := decodeSource(w, r)
	if !ok {
		return
	}
	lock := s.locks[problem.ID]
	lock.Lock()
	defer lock.Unlock()
	if err := writeSource(problem, request.Source); err != nil {
		writeError(w, http.StatusInternalServerError, "save exercise source")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"saved": true})
}

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	problem, ok := s.problemFromRequest(w, r)
	if !ok {
		return
	}
	request, ok := decodeSource(w, r)
	if !ok {
		return
	}
	lock := s.locks[problem.ID]
	lock.Lock()
	defer lock.Unlock()
	if err := writeSource(problem, request.Source); err != nil {
		writeError(w, http.StatusInternalServerError, "save exercise source")
		return
	}
	writeJSON(w, http.StatusOK, s.runner.Run(r.Context(), problem))
}

func (s *Server) problemFromRequest(w http.ResponseWriter, r *http.Request) (Problem, bool) {
	id := r.URL.Query().Get("id")
	problem, ok := s.catalog.Problem(id)
	if !ok {
		writeError(w, http.StatusNotFound, "unknown problem id")
		return Problem{}, false
	}
	return problem, true
}

func decodeSource(w http.ResponseWriter, r *http.Request) (sourceRequest, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "content type must be application/json")
		return sourceRequest{}, false
	}
	body := http.MaxBytesReader(w, r.Body, maxSourceBytes)
	defer body.Close()
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	var request sourceRequest
	if err := decoder.Decode(&request); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "source exceeds 1 MiB")
		} else {
			writeError(w, http.StatusBadRequest, "invalid JSON request")
		}
		return sourceRequest{}, false
	}
	if !utf8.ValidString(request.Source) {
		writeError(w, http.StatusBadRequest, "source must be valid UTF-8")
		return sourceRequest{}, false
	}
	if _, err := decoder.Token(); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return sourceRequest{}, false
	}
	return request, true
}

func writeSource(problem Problem, source string) error {
	path := filepath.Join(problem.Dir, "exercise.go")
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	temp, err := os.CreateTemp(problem.Dir, ".exercise-*.go")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.WriteString(source); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempName, path)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

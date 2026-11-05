package webapp

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var difficultyPattern = regexp.MustCompile(`[（(](Easy|Medium|Hard)[）)]`)

type Problem struct {
	ID         string `json:"id"`
	Stage      string `json:"stage"`
	StageTitle string `json:"stageTitle"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Order      int    `json:"order"`
	Dir        string `json:"-"`
}

type Catalog struct {
	root     string
	problems []Problem
	byID     map[string]Problem
}

func Discover(root string) (*Catalog, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read repository root: %w", err)
	}

	catalog := &Catalog{root: root, byID: make(map[string]Problem)}
	for _, entry := range entries {
		if !entry.IsDir() || !isCourseRoot(entry.Name()) {
			continue
		}
		stageDir := filepath.Join(root, entry.Name())
		stageTitle := firstHeading(filepath.Join(stageDir, "README.md"))
		if stageTitle == "" {
			stageTitle = entry.Name()
		}
		children, err := os.ReadDir(stageDir)
		if err != nil {
			return nil, fmt.Errorf("read stage %s: %w", entry.Name(), err)
		}
		for _, child := range children {
			if !child.IsDir() {
				continue
			}
			problemDir := filepath.Join(stageDir, child.Name())
			if !regularFile(filepath.Join(problemDir, "exercise.go")) || !regularFile(filepath.Join(problemDir, "exercise_test.go")) {
				continue
			}
			rel, err := filepath.Rel(root, problemDir)
			if err != nil {
				return nil, fmt.Errorf("build problem path: %w", err)
			}
			id := filepath.ToSlash(rel)
			title := firstHeading(filepath.Join(problemDir, "README.md"))
			if title == "" {
				title = child.Name()
			}
			difficulty := "Project"
			if match := difficultyPattern.FindStringSubmatch(title); len(match) == 2 {
				difficulty = match[1]
			}
			problem := Problem{
				ID:         id,
				Stage:      entry.Name(),
				StageTitle: stageTitle,
				Title:      title,
				Difficulty: difficulty,
				Order:      courseOrder(entry.Name()),
				Dir:        problemDir,
			}
			catalog.problems = append(catalog.problems, problem)
			catalog.byID[id] = problem
		}
	}

	sort.Slice(catalog.problems, func(i, j int) bool {
		left, right := catalog.problems[i], catalog.problems[j]
		if left.Order != right.Order {
			return left.Order < right.Order
		}
		return left.ID < right.ID
	})
	if len(catalog.problems) == 0 {
		return nil, errors.New("no exercises found in repository")
	}
	return catalog, nil
}

func (c *Catalog) Problems() []Problem {
	return append([]Problem(nil), c.problems...)
}

func (c *Catalog) Problem(id string) (Problem, bool) {
	p, ok := c.byID[id]
	return p, ok
}

func isCourseRoot(name string) bool {
	if name == "capstone" {
		return true
	}
	if !strings.HasPrefix(name, "stage") {
		return false
	}
	underscore := strings.IndexByte(name, '_')
	return underscore > len("stage")
}

func courseOrder(name string) int {
	if name == "capstone" {
		return 100
	}
	var order int
	if _, err := fmt.Sscanf(name, "stage%d_", &order); err != nil {
		return 99
	}
	return order
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func firstHeading(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func readOptionalFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

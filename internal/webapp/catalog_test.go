package webapp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover(t *testing.T) {
	root := newTestRepository(t)
	catalog, err := Discover(root)
	if err != nil {
		t.Fatal(err)
	}

	problems := catalog.Problems()
	if len(problems) != 2 {
		t.Fatalf("problem count = %d, want 2", len(problems))
	}
	if problems[0].ID != "stage1_foundation/ex01_registry" || problems[0].Difficulty != "Easy" {
		t.Fatalf("first problem = %+v", problems[0])
	}
	if problems[1].ID != "capstone/m1_project" || problems[1].Difficulty != "Project" {
		t.Fatalf("second problem = %+v", problems[1])
	}
	if got, ok := catalog.Problem(problems[0].ID); !ok || got.Dir == "" {
		t.Fatalf("problem lookup = %+v, %v", got, ok)
	}
}

func newTestRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeTestFile(t, root, "go.mod", "module example.test/course\n\ngo 1.26\n")
	writeTestFile(t, root, "stage1_foundation/README.md", "# Stage One\n")
	writeTestFile(t, root, "stage1_foundation/ex01_registry/README.md", "# ex01 Registry (Easy)\n\nImplement it.\n")
	writeTestFile(t, root, "stage1_foundation/ex01_registry/HINTS.md", "# Hints\n\nThink about maps.\n")
	writeTestFile(t, root, "stage1_foundation/ex01_registry/exercise.go", "package registry\n\nfunc Answer() int { return 42 }\n")
	writeTestFile(t, root, "stage1_foundation/ex01_registry/exercise_test.go", "package registry\n\nimport \"testing\"\n\nfunc TestAnswer(t *testing.T) { if Answer() != 42 { t.Fail() } }\n")
	writeTestFile(t, root, "capstone/README.md", "# Capstone\n")
	writeTestFile(t, root, "capstone/m1_project/README.md", "# M1 Project\n")
	writeTestFile(t, root, "capstone/m1_project/exercise.go", "package project\n")
	writeTestFile(t, root, "capstone/m1_project/exercise_test.go", "package project\n")
	writeTestFile(t, root, "stage9_empty/README.md", "# Empty Stage\n")
	return root
}

func writeTestFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

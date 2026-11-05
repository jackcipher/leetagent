package webapp

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PackageReference struct {
	Package    string            `json:"package"`
	ImportPath string            `json:"importPath"`
	Doc        string            `json:"doc"`
	Files      []ReferenceFile   `json:"files"`
	Symbols    []ReferenceSymbol `json:"symbols"`
}

type ReferenceFile struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type ReferenceSymbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Signature string `json:"signature"`
	Doc       string `json:"doc"`
	File      string `json:"file"`
	Line      int    `json:"line"`
}

func loadCommonReference(root string) (PackageReference, error) {
	commonDir := filepath.Join(root, "common")
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		return PackageReference{}, fmt.Errorf("read common package: %w", err)
	}

	ref := PackageReference{Package: "common", ImportPath: "leetagent/common"}
	fset := token.NewFileSet()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(commonDir, name)
		source, err := os.ReadFile(path)
		if err != nil {
			return PackageReference{}, fmt.Errorf("read common/%s: %w", name, err)
		}
		file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
		if err != nil {
			return PackageReference{}, fmt.Errorf("parse common/%s: %w", name, err)
		}
		if ref.Doc == "" && file.Doc != nil {
			ref.Doc = strings.TrimSpace(file.Doc.Text())
		}
		ref.Files = append(ref.Files, ReferenceFile{Name: name, Source: string(source)})
		ref.Symbols = append(ref.Symbols, exportedSymbols(fset, file, name)...)
	}
	if len(ref.Files) == 0 {
		return PackageReference{}, fmt.Errorf("common package has no Go source files")
	}

	sort.Slice(ref.Files, func(i, j int) bool {
		return referenceFileOrder(ref.Files[i].Name) < referenceFileOrder(ref.Files[j].Name)
	})
	sort.SliceStable(ref.Symbols, func(i, j int) bool {
		left, right := ref.Symbols[i], ref.Symbols[j]
		if symbolKindOrder(left.Kind) != symbolKindOrder(right.Kind) {
			return symbolKindOrder(left.Kind) < symbolKindOrder(right.Kind)
		}
		return left.Name < right.Name
	})
	return ref, nil
}

func exportedSymbols(fset *token.FileSet, file *ast.File, filename string) []ReferenceSymbol {
	var symbols []ReferenceSymbol
	for _, declaration := range file.Decls {
		switch declaration := declaration.(type) {
		case *ast.FuncDecl:
			if !declaration.Name.IsExported() {
				continue
			}
			copy := *declaration
			copy.Body = nil
			copy.Doc = nil
			kind := "func"
			if declaration.Recv != nil {
				kind = "method"
			}
			symbols = append(symbols, ReferenceSymbol{
				Name:      declaration.Name.Name,
				Kind:      kind,
				Signature: formatDeclaration(fset, &copy),
				Doc:       commentText(declaration.Doc),
				File:      filename,
				Line:      fset.Position(declaration.Pos()).Line,
			})
		case *ast.GenDecl:
			for _, spec := range declaration.Specs {
				symbol, ok := exportedSpec(fset, declaration, spec, filename)
				if ok {
					symbols = append(symbols, symbol)
				}
			}
		}
	}
	return symbols
}

func exportedSpec(fset *token.FileSet, declaration *ast.GenDecl, spec ast.Spec, filename string) (ReferenceSymbol, bool) {
	copy := &ast.GenDecl{Tok: declaration.Tok, Specs: []ast.Spec{spec}}
	switch spec := spec.(type) {
	case *ast.TypeSpec:
		if !spec.Name.IsExported() {
			return ReferenceSymbol{}, false
		}
		kind := "type"
		if _, ok := spec.Type.(*ast.InterfaceType); ok {
			kind = "interface"
		}
		return ReferenceSymbol{
			Name:      spec.Name.Name,
			Kind:      kind,
			Signature: formatDeclaration(fset, copy),
			Doc:       firstCommentText(spec.Doc, declaration.Doc),
			File:      filename,
			Line:      fset.Position(spec.Pos()).Line,
		}, true
	case *ast.ValueSpec:
		var exported []string
		for _, name := range spec.Names {
			if name.IsExported() {
				exported = append(exported, name.Name)
			}
		}
		if len(exported) == 0 {
			return ReferenceSymbol{}, false
		}
		return ReferenceSymbol{
			Name:      strings.Join(exported, ", "),
			Kind:      strings.ToLower(declaration.Tok.String()),
			Signature: formatDeclaration(fset, copy),
			Doc:       firstCommentText(spec.Doc, declaration.Doc),
			File:      filename,
			Line:      fset.Position(spec.Pos()).Line,
		}, true
	default:
		return ReferenceSymbol{}, false
	}
}

func formatDeclaration(fset *token.FileSet, node any) string {
	var output bytes.Buffer
	if err := format.Node(&output, fset, node); err != nil {
		return ""
	}
	return output.String()
}

func firstCommentText(groups ...*ast.CommentGroup) string {
	for _, group := range groups {
		if text := commentText(group); text != "" {
			return text
		}
	}
	return ""
}

func commentText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return strings.TrimSpace(group.Text())
}

func referenceFileOrder(name string) string {
	order := map[string]string{
		"types.go":      "0",
		"client.go":     "1",
		"scripted.go":   "2",
		"mockserver.go": "3",
	}
	if prefix, ok := order[name]; ok {
		return prefix + name
	}
	return "9" + name
}

func symbolKindOrder(kind string) int {
	switch kind {
	case "interface":
		return 0
	case "type":
		return 1
	case "var", "const":
		return 2
	case "func":
		return 3
	case "method":
		return 4
	default:
		return 5
	}
}

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type finding struct {
	File   string
	Line   int
	Kind   string
	Name   string
	Pkg    string
	Reason string
}

func main() {
	var root string
	flag.StringVar(&root, "root", ".", "repo root")
	flag.Parse()

	fset := token.NewFileSet()
	var findings []finding

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "node_modules", ".idea", ".vscode":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if strings.HasPrefix(rel, "examples/") {
			return nil
		}

		file, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			findings = append(findings, finding{File: rel, Line: 1, Kind: "file", Name: rel, Reason: parseErr.Error()})
			return nil
		}

		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !d.Name.IsExported() {
					continue
				}
				kind := "func"
				if d.Recv != nil {
					kind = "method"
				}
				if d.Doc == nil || len(strings.TrimSpace(d.Doc.Text())) == 0 {
					pos := fset.Position(d.Pos())
					findings = append(findings, finding{File: rel, Line: pos.Line, Kind: kind, Name: d.Name.Name, Pkg: file.Name.Name, Reason: "missing doc comment"})
				}

			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if !s.Name.IsExported() {
							continue
						}
						doc := s.Doc
						if doc == nil {
							doc = d.Doc
						}
						if doc == nil || len(strings.TrimSpace(doc.Text())) == 0 {
							pos := fset.Position(s.Pos())
							findings = append(findings, finding{File: rel, Line: pos.Line, Kind: "type", Name: s.Name.Name, Pkg: file.Name.Name, Reason: "missing doc comment"})
						}

					case *ast.ValueSpec:
						for _, n := range s.Names {
							if !n.IsExported() {
								continue
							}
							doc := s.Doc
							if doc == nil {
								doc = d.Doc
							}
							if doc == nil || len(strings.TrimSpace(doc.Text())) == 0 {
								pos := fset.Position(n.Pos())
								kind := strings.ToLower(d.Tok.String())
								findings = append(findings, finding{File: rel, Line: pos.Line, Kind: kind, Name: n.Name, Pkg: file.Name.Name, Reason: "missing doc comment"})
							}
						}
					}
				}
			}
		}

		return nil
	})

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Pkg != findings[j].Pkg {
			return findings[i].Pkg < findings[j].Pkg
		}
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Name < findings[j].Name
	})

	counts := map[string]int{}
	for _, f := range findings {
		counts[f.Pkg]++
	}

	fmt.Printf("Doc audit findings: %d\n", len(findings))
	pkgs := make([]string, 0, len(counts))
	for k := range counts {
		pkgs = append(pkgs, k)
	}
	sort.Strings(pkgs)
	for _, p := range pkgs {
		fmt.Printf("- %s: %d\n", p, counts[p])
	}
	fmt.Println("\nFindings:")
	for _, f := range findings {
		fmt.Printf("%s:%d: [%s] %s (%s)\n", f.File, f.Line, f.Kind, f.Name, f.Reason)
	}
}

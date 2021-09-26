package engine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strings"
)

type Engine interface {
	IdentifyModelTargets() ([]ModelTarget, error)
}

type ModelTarget interface {
	GetDocumentText() (string, error)
	GetModel() (interface{}, error)
	GetLocation() (string, error)
}

type modelTarget struct {
	astNode            ast.Node
	docText            string
	structTypeNode     *ast.StructType
	definitionPosition token.Position
}

func (mt *modelTarget) GetDocumentText() (string, error) {
	return mt.docText, nil
}
func (mt *modelTarget) GetModel() (interface{}, error) {
	return mt.structTypeNode, nil
}
func (mt *modelTarget) GetLocation() (string, error) {
	return mt.definitionPosition.String(), nil
}

type engine struct {
	root         fs.FS
	modelEntries []ModelTarget
}

var (
	debugLoggingOn = strings.ToLower(os.Getenv("AUTUMN_DEBUG_LOGGING_ON")) == "true"
)

func NewEngine(root fs.FS) (Engine, error) {
	var (
		fileEntries  = make(map[string]fs.DirEntry)
		modelEntries []ModelTarget
	)

	// Walk the directory to identify all go files.
	if err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".go") {
			fileEntries[path] = d
			if debugLoggingOn {
				fmt.Printf("identified file %q as a go file\n", path)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Now process files
	for fileName, _ := range fileEntries {
		fileContents, err := fs.ReadFile(root, fileName)
		if err != nil {
			return nil, err
		}

		targets, err := modelTargetFromText(fileName, string(fileContents))
		if err != nil {
			return nil, err
		}

		modelEntries = append(modelEntries, targets...)
	}

	return &engine{
		root:         root,
		modelEntries: modelEntries,
	}, nil
}

var autumnModelIdentifier = "@Autumn:Model"

func modelTargetFromText(name, text string) ([]ModelTarget, error) {
	var targets []ModelTarget

	// Parse the file
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, name, text, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Process the comments in the file
	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for node, comments := range cmap {
		for _, commentList := range comments {
			for _, comment := range commentList.List {
				if strings.Contains(comment.Text, autumnModelIdentifier) {
					decl, ok := node.(*ast.GenDecl)
					if !ok {
						if debugLoggingOn {
							fmt.Println("not a valid GenDecl")
						}
						continue
					} else if decl.Tok != token.TYPE {
						if debugLoggingOn {
							fmt.Println("not a type declaration")
						}
						continue
					}

					for _, spec := range decl.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							if debugLoggingOn {
								fmt.Println("not a struct type")
							}
							continue
						}

						structType, ok := typeSpec.Type.(*ast.StructType)
						if !ok {
							if debugLoggingOn {
								fmt.Println("type spec not for a struct")
							}
							continue
						}

						targets = append(targets, &modelTarget{
							astNode:            node,
							docText:            typeSpec.Doc.Text(),
							structTypeNode:     structType,
							definitionPosition: fset.Position(typeSpec.Pos()),
						})
					}
				}
			}
		}
	}

	return targets, nil
}

func (e *engine) IdentifyModelTargets() ([]ModelTarget, error) {
	return e.modelEntries, nil
}

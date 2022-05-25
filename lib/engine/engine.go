package engine

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strings"
)

// Engine is our entrypoint into how we load model targets for any and all
// future processing.
type Engine interface {
	// IdentifyModelTargets runs the engine identifying all identifiable
	// ModelTargets.
	IdentifyModelTargets() ([]ModelTarget, error)
}

// ModelTarget is a model struct that has been identified that we want to
// generate code for.
type ModelTarget interface {
	// Name returns the name of the model target.
	Name() (string, error)
	// PkgName returns the name of the package that the model is a part of.
	PkgName() string
	// GetDocumentText returns any comments attributed to the model target.
	GetDocumentText() (string, error)
	// GetModel returns the struct type node.
	GetModel() (interface{}, error)
	// GetLocation returns the location of the struct in its source file.
	GetLocation() (string, error)
	// ToTemplateVariables returns the model target as variables that can be
	// used by engine consumers for generating code.
	ToTemplateVariables() map[string]interface{}
}

// modelTarget is our local implementation of ModelTarget.
type modelTarget struct {
	astNode            ast.Node
	docText            string
	typeNode           *ast.TypeSpec
	structTypeNode     *ast.StructType
	definitionPosition token.Position
	pkgName            string
}

func (mt *modelTarget) Name() (string, error) {
	if mt.typeNode == nil {
		return "", errors.New("model target is missing a type node")
	}
	return mt.typeNode.Name.String(), nil
}

func (mt *modelTarget) PkgName() string {
	return mt.pkgName
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
func (mt *modelTarget) ToTemplateVariables() map[string]interface{} {
	return map[string]interface{}{
		"Name":             mt.typeNode.Name.String(),
		"ModelPackageName": mt.pkgName,
	}
}

type engine struct {
	root         fs.FS
	modelEntries []ModelTarget
}

var (
	debugLoggingOn = strings.ToLower(os.Getenv("AUTUMN_DEBUG_LOGGING_ON")) == "true"
)

// NewEngine returns a new engine rooted at the goven fs.FS root.
// It returns an error if it fails to walk the directory at the root of the
// fs.FS.
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

// modelTargetFromText creates a model target from source code when it finds
// a struct with the given annotation.
//
// NOTE: in a future iteration, it would be useful to also support providing
// model struct names via a config file in addition to using annotations.
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
							typeNode:           typeSpec,
							structTypeNode:     structType,
							definitionPosition: fset.Position(typeSpec.Pos()),
							pkgName:            f.Name.String(),
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

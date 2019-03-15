package documenting

import (
	"bytes"
	"errors"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os/exec"
)

type SearchingVisitor struct {
	Target string
	Found  bool

	CommentGroup *ast.CommentGroup
	Ident        *ast.Ident
}

func (self *SearchingVisitor) Visit(node ast.Node) (w ast.Visitor) {

	switch node.(type) {

	case *ast.CommentGroup:
		self.CommentGroup = node.(*ast.CommentGroup)
	case *ast.Ident:
		self.Ident = node.(*ast.Ident)
		if self.Ident.Name == self.Target {
			self.Found = true
		}
	}
	if self.Found {
		return nil
	}
	return self
}

func GetTextDocForFilesIdent(
	file string, name string,
	indent, preIndent string, width int,
) (string, error) {

	v := new(SearchingVisitor)
	v.Target = name

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	ast.Walk(v, f)

	if !v.Found {
		return "", errors.New("not found")
	}

	b := &bytes.Buffer{}

	doc.ToText(b, v.CommentGroup.Text(), indent, preIndent, width)

	return b.String(), nil
}

func UseGoDocInDir(dir string, sym string) (string, error) {

	args := []string{"doc", "."}
	if sym != "" {
		args = append(args, sym)
	}

	c := exec.Command("go", args...)
	c.Dir = dir

	b := &bytes.Buffer{}

	c.Stdout = b

	err := c.Run()
	if err != nil {
		return "", nil
	}

	return b.String(), nil
}

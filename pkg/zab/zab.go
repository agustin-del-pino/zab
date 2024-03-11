package zab

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

type structVisitor []*ast.TypeSpec

func (v *structVisitor) Visit(n ast.Node) ast.Visitor {
	switch s := n.(type) {
	case *ast.TypeSpec:
		if _, ok := s.Type.(*ast.StructType); ok {
			*v = append(*v, s)
		}
	}

	return v
}

func CreateBuilder(n *ast.TypeSpec) []byte {
	s := n.Type.(*ast.StructType)
	nm := fmt.Sprintf("%sBuilder", n.Name.Name)

	b := bytes.NewBufferString("// ")
	b.WriteString(nm)
	b.WriteString(" is the builder of ")
	b.WriteString(n.Name.Name)
	b.WriteString(".\ntype ")
	b.WriteString(nm)
	b.WriteString(" ")
	b.WriteString(n.Name.Name)
	b.WriteString("\n\n")

	for _, fld := range s.Fields.List {
		fn := fld.Names[0].Name
		ty, ok := fld.Type.(*ast.Ident)

		if !ok {
			continue
		}
		b.WriteString("// Set")
		b.WriteString(fn)
		b.WriteString(" returns the builder instance by applying 'v' to ")
		b.WriteString(n.Name.Name)
		b.WriteString(".")
		b.WriteString(fn)
		b.WriteString(".\nfunc (b *")
		b.WriteString(nm)
		b.WriteString(") Set")
		b.WriteString(fn)

		b.WriteString("(v ")
		b.WriteString(ty.Name)
		b.WriteString(") *")
		b.WriteString(nm)
		b.WriteString(" {\nb.")
		b.WriteString(fn)
		b.WriteString(" = v\nreturn b\n}\n\n")

	}
	return b.Bytes()
}

func GetStructs(n *ast.File) []*ast.TypeSpec {
	var v structVisitor
	ast.Walk(&v, n)
	return v
}

func ReadFile(f string) (*ast.File, *token.FileSet, error) {
	s, sErr := os.ReadFile(f)

	if sErr != nil {
		return nil, nil, sErr
	}
	fs := token.NewFileSet()
	n, nErr := parser.ParseFile(fs, f, s, parser.AllErrors)

	if nErr != nil {
		log.Fatal(nErr)
	}

	return n, fs, nErr
}

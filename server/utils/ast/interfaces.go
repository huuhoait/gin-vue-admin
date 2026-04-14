package ast

import (
	"go/ast"
	"io"
)

type Ast interface {
	// Parse parse file/GenerationCode
	Parse(filename string, writer io.Writer) (file *ast.File, err error)
	// Rollback rollback
	Rollback(file *ast.File) error
	// Injection inject
	Injection(file *ast.File) error
	// Format FormatInputOut
	Format(filename string, writer io.Writer, file *ast.File) error
}

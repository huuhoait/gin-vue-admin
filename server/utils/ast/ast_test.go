package ast

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestAst(t *testing.T) {
	filename := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva", "plugin.go")
	if _, err := os.Stat(filename); err != nil {
		t.Skipf("skip: %s not present: %v", filename, err)
	}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		t.Error(err)
		return
	}
	err = ast.Print(fileSet, file)
	if err != nil {
		t.Error(err)
		return
	}
	err = printer.Fprint(os.Stdout, token.NewFileSet(), file)
	if err != nil {
		panic(err)
	}

}

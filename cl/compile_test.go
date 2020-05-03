package cl

import (
	"fmt"
	"testing"

	"github.com/qiniu/qlang/ast/asttest"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/parser"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"

	_ "github.com/qiniu/qlang/lib/builtin"
)

func init() {
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(log.Ldebug)
}

// -----------------------------------------------------------------------------

var fsTestBasic = asttest.NewSingleFileFS("/foo", "bar.ql", `
	println("Hello", "xsw", "- nice to meet you!")
	println("Hello, world!")
`)

func TestBasic(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestBasic, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(14) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

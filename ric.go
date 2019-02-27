package goScript

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/golang/glog"

	"github.com/complyue/hbigo/pkg/errors"
)

func RunInContext(scriptName, code string, context interface{}) (val interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf(`Error running code in context:
-*-%s-*-
%s
=*-%s-*=`, scriptName, code, scriptName)
			err = errors.Errorf("Error run in context: %+v", r)
		}
	}()

	ctxt := createContext(context)

	pf, err := parser.ParseFile(&token.FileSet{}, scriptName, fmt.Sprintf(`
package hbi
func init() {
%s
}
`, code), 0)
	if err != nil {
		return nil, err
	}
	for _, decl := range pf.Decls {
		if initf, ok := decl.(*ast.FuncDecl); ok && "init" == initf.Name.Name {
			for _, stmt := range initf.Body.List {
				val, err = eval(stmt, ctxt)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return val, err
}

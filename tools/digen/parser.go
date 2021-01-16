package digen

import (
	"fmt"
	"github.com/loilo-inc/logos/set"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"sort"
	"strings"
)

type Manifest struct {
	Imports []string
	Fields  [][]string
}

func (m *Manifest) IsEmpty() bool {
	return len(m.Imports) == 0 && len(m.Fields) == 0
}

func ParseManifest(src interface{}) (*Manifest, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	imports := set.NewStringSet()
	m := file.Scope.Lookup("Manifest")
	if m == nil {
		return &Manifest{}, nil
	}
	for _, im := range file.Imports {
		value := fmt.Sprintf("%s", im.Path.Value)
		if im.Name != nil {
			value = fmt.Sprintf("%s %s", im.Name.Name, value)
		}
		imports.Add(value)
	}
	var fields [][]string
	if spec, ok := m.Decl.(*ast.TypeSpec); ok {
		if st, ok := spec.Type.(*ast.StructType); ok {
			for _, field := range st.Fields.List {
				name := field.Names[0].String()
				t := printType(field.Type)
				fields = append(fields, []string{
					name, strings.Join(t, ""),
				})
			}
		}
	}
	sortedImports := imports.Values()
	sort.Strings(sortedImports)
	return &Manifest{
		Imports: sortedImports,
		Fields:  fields,
	}, nil
}

func printType(expr ast.Expr) []string {
	switch expr := expr.(type) {
	case *ast.Ident:
		if expr.Obj != nil {
			switch decl := expr.Obj.Decl.(type) {
			case *ast.TypeSpec:
				return []string{decl.Name.String()}
			}
			a := expr.Obj.Decl.(int)
			log.Printf("%d", a)
			return []string{expr.Obj.Name}
		}
		return []string{expr.String()}
	case *ast.InterfaceType:
		return []string{"interface{}"}
	case *ast.StarExpr:
		ret := []string{"*"}
		ret = append(ret, printType(expr.X)...)
		return ret
	case *ast.StructType:
		return []string{"struct{}"}
	case *ast.ArrayType:
		ret := []string{
			"[",
		}
		ret = append(ret, printType(expr.Len)...)
		ret = append(ret, "]")
		ret = append(ret, printType(expr.Elt)...)
		return ret
	case *ast.BasicLit:
		return []string{expr.Value}
	case nil:
		return []string{}
	case *ast.SelectorExpr:
		var ret []string
		ret = append(ret, printType(expr.X)...)
		ret = append(ret, ".")
		ret = append(ret, printType(expr.Sel)...)
		return ret
	default:
		// Trigger panic
		v := reflect.ValueOf(expr)
		log.Printf("unknwon expr: %s", v.Type().String())
		return []string{"interface{}"}
	}
}

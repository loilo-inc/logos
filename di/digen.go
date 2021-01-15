package di

import (
	"fmt"
	"github.com/loilo-inc/logos/set"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

func Generate(outDir string, pkg string, manifest interface{}) {
	value := reflect.ValueOf(manifest)
	valueType := value.Type()

	fileTmpl, err := template.New("file").Parse(file)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	attrTmpl, err := template.New("attr").Parse(attr)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	packages := set.NewStringSet()
	wrap := func(s string) string {
		return fmt.Sprintf("\"%s\"", s)
	}
	packages.Add(fmt.Sprintf("logosdi %s", wrap("github.com/loilo-inc/logos/di")))
	var values []map[string]string
	out, _ := os.Create(filepath.Join(outDir, "di_ext.go"))

	for i := 0; i < value.NumField(); i++ {
		f := value.Field(i)
		t := f.Type()
		v := valueType.Field(i)
		data := map[string]string{
			"Field": v.Name,
			"Type":  t.String(),
		}
		fmt.Println(data["Field"])
		packages.Add(wrap(t.PkgPath()))
		values = append(values, data)
	}
	imports := strings.Join(packages.Values(), "\n")
	if err := fileTmpl.Execute(out, map[string]string{
		"Package": pkg,
		"Imports": imports,
	}); err != nil {
		log.Fatalf(err.Error())
	}
	for _, v := range values {
		if err := attrTmpl.Execute(out, v); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

var file = `package {{.Package}}

import (
{{.Imports}}
)

`
var attr = `
type {{.Field}}Factory func() {{.Type}}
func (d *D) Get{{.Field}}() {{.Type}} {
  v := d.Get("{{.Field}}").({{.Type}})
	return v
}
func (b *B) Set{{.Field}}(f {{.Field}}Factory) {
  b.Set("{{.Field}}", f)
}

`

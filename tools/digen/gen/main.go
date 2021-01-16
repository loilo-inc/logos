package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	m, err := ioutil.ReadFile("./di/di.go")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out, err := os.Create("./cmd/digen/main.go")
	if err != nil {
		log.Fatalf(err.Error())
	}
	tmpl, err := template.ParseFiles("./tools/digen/gen/main.tmpl")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err := tmpl.Execute(out, map[string]string{
		"Body": string(m),
	}); err != nil {
		log.Fatalf(err.Error())
	}
}

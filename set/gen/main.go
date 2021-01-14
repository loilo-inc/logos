package main

import (
	"fmt"
	"log"
	"os"
)
import "text/template"

type parameter struct {
	Iface  string
	Struct string
	Type   string
	Values string
}

func main() {
	tmpl, err := template.ParseFiles("./set/gen/set.tmpl", "./set/gen/test.tmpl")
	if err != nil {
		log.Fatalf(err.Error())
	}
	params := []parameter{
		{Iface: "StringSet", Struct: "stringSet", Type: "string", Values: "[]string{\"1\",\"2\"}"},
		{Iface: "IntSet", Struct: "intSet", Type: "int", Values: "[]int{1,2}"},
		{Iface: "Int64Set", Struct: "int64Set", Type: "int64", Values: "[]int64{1,2}"},
		{Iface: "BoolSet", Struct: "boolSet", Type: "bool", Values: "[]bool{false,true}"},
	}
	for _, p := range params {
		dest := fmt.Sprintf("./set/%s_set.go", p.Type)
		destFile, _ := os.Create(dest)
		testDest := fmt.Sprintf("./set/%s_set_test.go", p.Type)
		testDestFile, _ := os.Create(testDest)
		if err := tmpl.ExecuteTemplate(destFile, "set.tmpl", p); err != nil {
			log.Fatalf(err.Error())
		}
		if err := tmpl.ExecuteTemplate(testDestFile, "test.tmpl", p); err != nil {
			log.Fatalf(err.Error())
		}

	}
}

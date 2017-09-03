package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {

	iobuf := bytes.NewBuffer(nil)

	s1, _ := template.ParseFiles("header.tmpl", "content.tmpl", "footer.tmpl")
	s1.ExecuteTemplate(iobuf, "header", nil)
	s1.ExecuteTemplate(iobuf, "content", nil)
	s1.ExecuteTemplate(iobuf, "footer", nil)
	s1.Execute(iobuf, nil)

	fmt.Println(iobuf.String())
}

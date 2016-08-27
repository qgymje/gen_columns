package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"reflect"
)

type User struct {
	Name string `json:"name" column:"name"`
	Age  int    `json:"age" column:"age"`
}

var (
	tagName = flag.String("tag", "column", "请指定tag")
)

// 0. 指定tagname
// 1. 扫描文件所有的struct, 并且读取所有的struct name
// 2. 读取所有的FieldName
// 3. 将数据填写到template中
// 4. 生成一个lower case structName_column.go文件

var temp2 = `
//package PackageName

type {{.StructName}}Column struct {
{{range $key, $value := .Columns}}
  {{ $key }}: string
{{end}}
}

var {{.StructName}}Columns  {{.StructName}}Column

func init() {
{{range $key, $value := .Columns}}
{{ $.StructName}}.{{$key}} = "{{$value}}"
{{end}}
}

`

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.Parse()

	var u User

	name := GetStructName(u)
	columns := GetFiledWithTag(u, *tagName)

	data := map[string]interface{}{
		"StructName": name,
		"Columns":    columns,
	}

	template.Must(template.New("temp2").Parse(temp2)).Execute(os.Stdout, data)
}

func GetStructName(i interface{}) string {
	t := reflect.TypeOf(i)
	return t.Name()
}

func GetFiledWithTag(i interface{}, tagName string) map[string]string {
	columns := make(map[string]string)
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Name
		tag := t.Field(i).Tag.Get(tagName)
		columns[field] = tag
	}
	return columns
}

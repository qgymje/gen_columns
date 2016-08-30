package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Barrage struct {
	UserID    bson.ObjectId `bson:"user_id"`
	RoomID    bson.ObjectId `bson:"room_id"`
	Message   string        `bson:"message"`
	CreatedAt time.Time     `bson:"created_at"`
}

var (
	tagName = flag.String("tag", "column", "请指定tag")
)

// 0. 指定tagname
// 1. 扫描文件所有的struct, 并且读取所有的struct name
// 1.1 读取指定文件到
// 1.2 使用ast token.NewFileSet
// 2. 读取所有的FieldName
// 3. 将数据填写到template中
// 4. 生成一个lower case structName_column.go文件

var temp2 = `
//package PackageName

type {{.StructName}}Column struct {
{{range $key, $value := .Columns}}
  {{ $key }} string
{{end}}
}

var {{.StructName}}Columns  {{.StructName}}Column

func init() {
{{range $key, $value := .Columns}}
{{ $.StructName}}Columns.{{$key}} = "{{$value}}"
{{end}}
}

`

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.Parse()

	var u Barrage

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

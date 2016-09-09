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

var (
	tagName = flag.String("tag", "column", "请指定tag")
)

type BroadcastRoom struct {
	ID          bson.ObjectId   `bson:"_id"`
	UserID      bson.ObjectId   `json:"userID" bson:"userId"`           //用户ID
	RoomID      int64           `bson:"roomID"`                         // RoomID号
	Name        string          `json:"name" bson:"name"`               //标题
	Cover       string          `json:"cover" bson:"cover"`             //封面图片地址
	Domain      string          `json:"-" bson:"domain"`                //个性域名
	Channel     []string        `json:"channel" bson:"channel"`         //领域
	Score       int32           `json:"score" bson:"score"`             //分数
	Tags        []string        `json:"tags" bson:"tags"`               //标签
	IsPlaying   bool            `json:"isPlaying" bson:"isPlaying"`     //是否正在直播
	AdminUsers  []bson.ObjectId `json:"adminUsers" bson:"adminUsers"`   //管理员ID
	ValidStatus int8            `json:"validStatus" bson:"validStatus"` //0 申请中未审核,1 审核通过, -1 审核失败
	Orientation int8            `json:"orientation" bson:"orientation"` //横竖屏 0 未设置 1 横屏 2 竖屏
	CreatedTime time.Time       `json:"createdTime" bson:"createdTime"` //创建时间
	UpdatedTime time.Time       `json:"updatedTime" bson:"updatedTime"` //更新时间
}

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
{{range $key, $value := .Columns}} {{ $key }} string 
{{end}}
}

var {{.StructName}}Columns  {{.StructName}}Column

func init() {
{{range $key, $value := .Columns}} {{ $.StructName}}Columns.{{$key}} = "{{$value}}" 
{{end}}
}
`

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.Parse()

	var s BroadcastRoom

	name := GetStructName(s)
	columns := GetFiledWithTag(s, *tagName)

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

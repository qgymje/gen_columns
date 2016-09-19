package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	tagName = flag.String("tag", "bson", "tag name")
	path    = flag.String("path", "./", "file path")
	debug   = flag.Bool("debug", false, "show generated data")
)

var temp = `
package {{.PackageName}}

type _{{.StructName}}Column struct {
{{range $key, $value := .Columns}} {{ $key }} string 
{{end}}
}

// {{.StructName}}Columns {{lower .StructName}} columns name
var {{.StructName}}Columns  _{{.StructName}}Column

func init() {
{{range $key, $value := .Columns}} {{ $.StructName}}Columns.{{$key}} = "{{$value}}" 
{{end}}
}
`

// TempData 表示生成template所需要的数据结构
type TempData struct {
	FileName    string
	PackageName string
	StructName  string
	Columns     map[string]string
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.Parse()

	filepath.Walk(*path, func(filename string, f os.FileInfo, _ error) error {
		if filepath.Ext(filename) == ".go" {
			if strings.Contains(filename, "_column.go") {
				return nil
			}
			return handleFile(filename)
		}
		return nil
	})
}

func handleFile(filename string) error {
	var tempData TempData
	tempData.Columns = make(map[string]string)

	fset := token.NewFileSet()
	var src interface{}
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		panic(err)
	}

	//ast.Print(fset, f)
	tempData.PackageName = f.Name.Name
	tempData.FileName = filename

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, s := range x.Specs {
					vSpec := s.(*ast.TypeSpec)
					if _, ok := vSpec.Type.(*ast.StructType); ok {
						tempData.StructName = vSpec.Name.Name
					}
				}
			}

		case *ast.StructType:
			for _, f := range x.Fields.List {
				if f.Tag != nil {
					tag := handleTags(f.Tag.Value)
					tempData.Columns[f.Names[0].Name] = tag
				}
			}
		}
		return true
	})

	if *debug {
		//spew.Dump(tempData)
		tempData.writeTo(os.Stdout)
	}
	return tempData.WriteToFile()
}

func handleTags(tags string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s:"(.*?)"`, *tagName))
	matchs := re.FindStringSubmatch(tags)
	if len(matchs) >= 1 {
		return matchs[1]
	}
	return ""
}

func (d *TempData) handleFilename() {
	absPath, _ := filepath.Abs(d.FileName)
	basePath := filepath.Dir(absPath)
	d.FileName = basePath + "/" + strings.ToLower(d.StructName) + "_column.go"
}

func (d *TempData) writeTo(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("temp").Funcs(funcMap).Parse(temp)).Execute(w, d)
}

// WriteToFile 将生成好的模块文件写到本地
func (d *TempData) WriteToFile() error {
	d.handleFilename()
	file, err := os.Create(d.FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = d.writeTo(&buf)
	formatted, _ := format.Source(buf.Bytes())
	file.Write(formatted)
	return err
}

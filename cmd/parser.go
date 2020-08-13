package cmd

import (
	"bytes"
	"fmt"
	"github.com/vetcher/go-astra"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	//jsonName  = "json"
	tableName = "TableName"
	temp      = `
		package table

		type _{{.StructName}} struct {
			TableName string
		{{range $key, $value := .Columns}} {{ $key }} TableField 
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"
		{{range $key, $value := .Columns}} 
		{{ $.StructName}}.{{$key}} = TableField{
			Name: "{{index $value 0}}",
			Json: "{{index $value 1}}",
		} 
		{{end}}
		}
		`
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	FileName    string
	PackageName string
	StructName  string
	TableName   string
	Columns     map[string][]string
}

//
func handleFile(filename string) error {
	tempData := new(TempData)

	tempData.Columns = make(map[string][]string)

	fset := token.NewFileSet()
	var src interface{}
	_, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		showError(err)
		panic(err)
	}

	////////////
	file, err := astra.ParseFile(filename,
		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
			astra.IgnoreInterfaces|astra.IgnoreTypes)
	if err != nil {
		showError(err)
	}
	tempData.PackageName = file.Name
	//
	functions := make(map[string]bool)
	for _, fun := range file.Methods {
		if fun.Name == "TableName" {
			functions[fun.Receiver.Type.String()] = true
		}
	}

	for _, stru := range file.Structures {
		tempData.FileName = filename
		tempData.StructName = stru.Name
		tempData.TableName = parseDoc(strings.Join(stru.Docs, " "))
		if tempData.TableName == "" {
			tempData.TableName = tempData.StructName
		}

		for _, field := range stru.Fields {
			var _namejson = make([]string, 3)
			for k, v := range field.Tags {
				if k == "json" {
					_namejson[1] = v[0]
				} else if k == XORM_TAG {
					_namejson[0] = parseTagsForXORM(v)
				} else if k == GORM_TAG {
					_namejson[0] = parseTagsForGORM(v)
				}
			}
			_namejson[2] = field.Type.String()

			if _namejson[1] == "" {
				if _namejson[0] == "" {
					_namejson[1] = getFieldName(field.Name)
				} else {
					_namejson[1] = _namejson[0]
				}
			}
			if _namejson[0] == "" {
				if _namejson[1] == "" {
					_namejson[0] = getFieldName(field.Name)
				} else {
					_namejson[0] = _namejson[1]
				}
			}
			tempData.Columns[field.Name] = _namejson
		}
		//if functions["*"+stru.Name] != true {
		err = tempData.appendToModel(filename, tempData.StructName)
		if err != nil {
			showError(err)
			return err
		}
		//}

		if len(tempData.StructName) == 0 ||
			tempData.StructName[:1] == strings.ToLower(tempData.StructName[:1]) ||
			len(tempData.Columns) == 0 {
			return nil
		}

		if debug {
			err = tempData.writeTo(os.Stdout)
		}

		err := tempData.writeToFile()
		if err != nil {
			showError(err.Error())
			return err
		}
	}

	return err
}

func parseTagsForXORM(matchs []string) string {
	if len(matchs) >= 1 {
		_matchs := regexp.MustCompile(`'(.*?)'`).FindStringSubmatch(matchs[0])
		if len(_matchs) >= 1 {
			return _matchs[1]
		}
	}
	return ""
}

func parseTagsForGORM(matchs []string) string {
	if len(matchs) >= 1 {
		_matchs := regexp.MustCompile(`(?i:column):(.*?)(?:;|$)`).FindStringSubmatch(matchs[0])
		if len(_matchs) >= 1 {
			return _matchs[1]
		}
	}
	return ""
}

func parseDoc(doc string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i:%s)[: ]+(.*)`, tableName))
	matchs := re.FindStringSubmatch(doc)

	if len(matchs) >= 1 {
		return strings.TrimSpace(matchs[1])
	}
	return ""
}

func getFilepath(filename string) string {
	absPath, _ := filepath.Abs(filename)
	return filepath.Join(filepath.Dir(absPath), "table")
}

func (d *TempData) handleFilename() {
	d.FileName = filepath.Join(getFilepath(d.FileName), strings.ToLower(d.StructName)+"_table.go")
}

func (d *TempData) writeTo(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("temp").Funcs(funcMap).Parse(temp)).Execute(w, d)
}

// writeToFile 将生成好的模块文件写到本地
func (d *TempData) writeToFile() error {
	d.handleFilename()
	file, err := os.Create(d.FileName)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = d.writeTo(&buf)
	formatted, _ := format.Source(buf.Bytes())
	_, err = file.Write(formatted)
	return err
}

var model_str = `
package {{.PackageName}}

import (
	"sync"
)

var (
	{{lower .StructName}}Pool = sync.Pool{
		New: func() interface{} {
			return &{{.StructName}}{}
		},
	}
)

func New{{.StructName}}() *{{.StructName}} {
	return {{lower .StructName}}Pool.Get().(*{{.StructName}})
}

func (p *{{.StructName}}) Free() {
	{{range $key, $value := .Columns}}p.{{$key}} = {{getTypeValue $value}}				
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}

func (*{{.StructName}}) TableName() string {
	return table.{{.StructName}}.TableName
}
	`

func (d *TempData) appendToModel(fileName, tableName string) error {
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"getTypeValue": func(t []string) interface{} {
			if len(t) < 3 {
				return `""`
			}
			var ret interface{}
			switch t[2] {
			case "string":
				ret = `""`
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64":
				ret = 0
			case "time.Time":
				ret = `time.Time{}`
			default:
				ret = `""`
			}
			return ret
		},
	}

	err := template.Must(template.New("temp").Funcs(funcMap).Parse(model_str)).Execute(&buf, d)
	if err != nil {
		showError(err)
		return err
	}
	absPath, _ := filepath.Abs(fileName)
	fileName = filepath.Join(filepath.Dir(absPath), strings.ToLower(d.StructName)+"_sorm.go")
	file, err := os.Create(fileName)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()

	//file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
	//if err != nil {
	//	showError(err.Error())
	//	return err
	//}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		showError(err)
		return err
	}

	return nil
}

//func (d *TempData) writeBaseFile() error {
//	baseFilename := filepath.Join(getFilepath(d.FileName), "base.go")
//
//	file, err := os.Create(baseFilename)
//	if err != nil {
//		showError(err.Error())
//		return err
//	}
//	defer file.Close()
//	var buf bytes.Buffer
//	_ = template.Must(template.New("temp").Parse(base)).Execute(&buf, d)
//	formatted, _ := format.Source(buf.Bytes())
//	_, err = file.Write(formatted)
//	if err != nil {
//		showError(err.Error())
//	}
//	return err
//}

func getFieldName(name string) string {
	bs := bytes.NewBuffer([]byte{})
	for i, s := range name {
		if s >= 65 && s <= 90 {
			s += 32
			if i == 0 {
				bs.WriteByte(byte(s))
			} else {
				bs.WriteByte(byte(95))
				bs.WriteByte(byte(s))
			}
			continue
		}
		bs.WriteByte(byte(s))
	}
	return bs.String()
}

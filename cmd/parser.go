package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/vetcher/go-astra"
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	FileName    string
	PackageName string
	StructName  string
	TableName   string
	Columns     map[string][]string
}

//handleFile 处理model文件
func handleFile(filename string) error {
	tempData := new(TempData)

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
	//functions := make(map[string]bool)
	//for _, fun := range file.Methods {
	//	if fun.Name == "TableName" {
	//		functions[fun.Receiver.Type.String()] = true
	//	}
	//}

	for _, stru := range file.Structures {
		tempData.Columns = make(map[string][]string)
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
		//如果struct名称为空,或者是一个私有struct,或者field为空,返回
		if len(tempData.StructName) == 0 ||
			tempData.StructName[:1] == strings.ToLower(tempData.StructName[:1]) ||
			len(tempData.Columns) == 0 {
			return nil
		}

		if debug {
			return tempData.writeTo(os.Stdout)
		}
		//写model文件
		//if !functions["*"+stru.Name] && !functions[stru.Name] {
		err = tempData.writeModel(filename)
		if err != nil {
			showError(err)
			return err
		}
		//}

		//写table文件
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
func getBaseFilename(filename string) string {
	f := filepath.Base(filename)
	pos := strings.LastIndex(f, ".")
	if pos == -1 {
		return f
	}
	return f[:pos]
}

func (d *TempData) handleFilename() string {
	return filepath.Join(getFilepath(d.FileName), getBaseFilename(d.FileName)+"_"+d.StructName+"_table.go")
}

func (d *TempData) writeTo(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("tableTpl").Funcs(funcMap).Parse(tableTpl)).Execute(w, d)
}

// writeToFile 将生成好的模块文件写到本地
func (d *TempData) writeToFile() error {
	file, err := os.Create(d.handleFilename())
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

func getFieldName(name string) string {
	bs := bytes.NewBuffer([]byte{})

	pre_lower := true //前一个字母是小写
	for i, s := range name {
		//如果是大写字母
		if s >= 65 && s <= 90 {
			s += 32 //转成小写
			if i == 0 {
				bs.WriteByte(byte(s))
			} else {
				if pre_lower {
					bs.WriteByte(byte(95)) //写下划线
				}
				bs.WriteByte(byte(s))
			}
			pre_lower = false
			continue
		}
		pre_lower = true
		bs.WriteByte(byte(s))
	}
	return bs.String()
}

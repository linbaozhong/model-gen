package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
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
	jsonName  = "json"
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
	base = `
		package table

		const (
			Quote_Char = "` + "`" + `"
		)

		type TableField struct {
			Name string
			Json string
		}
		//Eq 等于
		func (f TableField) Eq() string {
			return f.generate("=")
		}
		//Gt 大于
		func (f TableField) Gt() string {
			return f.generate(">")
		}
		//Gte 大于等于
		func (f TableField) Gte() string {
			return f.generate(">=")
		}
		//Lt 小于
		func (f TableField) Lt() string {
			return f.generate("<")
		}
		//Lte 小于等于
		func (f TableField) Lte() string {
			return f.generate("<=")
		}
		//Ue 不等于
		func (f TableField)Ue() string {
			return f.generate("<>")
		}
		//Bt BETWEEN
		func (f TableField)Bt() string {
			return f.QuoteName() + " BETWEEN ? AND ?"
		}
		//In IN
		func (f TableField)In() string {
			return f.QuoteName() + " IN (?)"
		}
		
		func (f TableField) QuoteName() string {
			return Quote_Char + f.Name + Quote_Char
		}
		func (f TableField) generate(op string) string {
			return f.QuoteName() + op + "?"
		}

		`
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	FileName    string
	PackageName string
	StructName  string
	TableName   string
	Columns     map[string][3]string
}

func handleFile(filename string) error {
	var tempData TempData
	tempData.Columns = make(map[string][3]string)

	fset := token.NewFileSet()
	var src interface{}
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	//ast.Print(fset, f)
	tempData.PackageName = f.Name.Name
	tempData.FileName = filename
	var hasTableNameFunc bool

	ast.Inspect(f, func(n ast.Node) bool {

		switch x := n.(type) {

		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, s := range x.Specs {
					vSpec := s.(*ast.TypeSpec)
					if _, ok := vSpec.Type.(*ast.StructType); ok {
						tempData.StructName = vSpec.Name.Name

						if x.Doc != nil {
							tempData.TableName = parseDoc(x.Doc.List[len(x.Doc.List)-1].Text)
						}
						if tempData.TableName == "" {
							tempData.TableName = tempData.StructName
						}
					}
				}
			}

		case *ast.StructType:
			for _, f := range x.Fields.List {
				var _namejson [3]string
				if f.Tag != nil {
					_namejson = parseTags(f.Tag.Value)
				}
				if _namejson[0] == "" {
					_namejson[0] = strings.ToLower(f.Names[0].Name)
				}
				if _namejson[1] == "" {
					_namejson[1] = strings.ToLower(f.Names[0].Name)
				}
				_namejson[2] = fmt.Sprintf("%v", f.Type)
				tempData.Columns[f.Names[0].Name] = _namejson
			}

		case *ast.FuncDecl:
			if x.Name.Name == "TableName" &&
				x.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name == tempData.StructName {
				hasTableNameFunc = true
			}
		}
		return true
	})

	if !hasTableNameFunc {
		tempData.appendToModel(filename, tempData.StructName)
	}

	if len(tempData.StructName) == 0 ||
		tempData.StructName[:1] == strings.ToLower(tempData.StructName[:1]) ||
		len(tempData.Columns) == 0 {
		return nil
	}
	if debug {
		tempData.writeTo(os.Stdout)
	}
	if err := tempData.writeBaseFile(); err != nil {
		return err
	}
	return tempData.writeToFile()
}

func parseTags(tags string) [3]string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i:%s):"(.*?)"`, tagName))
	matchs := re.FindStringSubmatch(tags)

	var col_name string
	if tagName == XORM_TAG {
		col_name = parseTagsForXORM(matchs)
	} else {
		col_name = parseTagsForGORM(matchs)
	}
	json_name := getJsonName(tags)
	if json_name == "" && col_name != "" {
		json_name = col_name
	}
	return [3]string{col_name, json_name, ""}
}

func parseTagsForXORM(matchs []string) string {
	if len(matchs) >= 1 {
		_matchs := regexp.MustCompile(`'(.*?)'`).FindStringSubmatch(matchs[1])
		if len(_matchs) >= 1 {
			return _matchs[1]
		}
	}
	return ""
}

func parseTagsForGORM(matchs []string) string {
	if len(matchs) >= 1 {
		_matchs := regexp.MustCompile(`(?i:column):(.*?)(?:;|$)`).FindStringSubmatch(matchs[1])
		if len(_matchs) >= 1 {
			return _matchs[1]
		}
	}
	return ""
}

func getJsonName(tags string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i:%s):"(.*?)"`, jsonName))
	matchs := re.FindStringSubmatch(tags)
	if len(matchs) >= 1 {
		return matchs[1]
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
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = d.writeTo(&buf)
	formatted, _ := format.Source(buf.Bytes())
	file.Write(formatted)
	return err
}
func (d *TempData) appendToModel(fileName, tableName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	str := `
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
			//todo:初始化每个字段
			{{range $key, $value := .Columns}}p.{{$key}} = ""
			{{end}}
			{{lower .StructName}}Pool.Put(p)
		}

		func (*{{.StructName}}) TableName() string {
			return table.{{.StructName}}.TableName
		}
	`
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	template.Must(template.New("temp").Funcs(funcMap).Parse(str)).Execute(&buf, d)
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}
func (d *TempData) writeBaseFile() error {
	baseFilename := filepath.Join(getFilepath(d.FileName), "base.go")

	file, err := os.Create(baseFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = template.Must(template.New("temp").Parse(base)).Execute(&buf, d)
	formatted, _ := format.Source(buf.Bytes())
	file.Write(formatted)
	return err
}

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var model_str = `
package {{.PackageName}}

import (
	"sync"
	"libs/utils"
	{{if .HasTime}}"time"{{end}}
	"libs/types"
	"{{.ModulePath}}/table"
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

//Free 
func (p *{{.StructName}}) Free() {
	if p == nil {
		return
	}
	{{range $key, $value := .Columns}}p.{{$key}} = {{getTypeValue $value}}				
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}

//TableName
func (*{{.StructName}}) TableName() string {
	return table.{{.StructName}}.TableName
}

//ToMap struct转map
func (p *{{.StructName}}) ToMap(cols...table.TableField) map[string]interface{} {
	l := len(cols)
	if l == 0 {
		return map[string]interface{}{
			{{range $key, $value := .Columns}}table.{{$.StructName}}.{{$key}}.Name:p.{{$key}},
			{{end}}
		}
	}

	m := make(map[string]interface{},l)
	for i := 0; i < l; i++ {
		col := cols[i]
		switch col.Name {
		{{range $key, $value := .Columns}}case table.{{$.StructName}}.{{$key}}.Name:
			m[col.Name] = p.{{$key}}
		{{end}}
		}
	}
	return m
}

//ToJSON struct转json
func (p *{{.StructName}}) ToJSON(cols...table.TableField) types.Smap {
	m := p.ToMap()
	clen := len(cols)
	if clen == 0 {
		l := len(table.{{.StructName}}.ColumnNames)
		sm := make(types.Smap, l)
		var cn string
		for i := 0; i < l; i++ {
			cn = table.{{.StructName}}.ColumnNames[i]
			sm.Set(table.{{.StructName}}.ColumnName2Json[cn], m[cn])
		}
		return sm
	}

	sm := make(types.Smap, clen)
	for i := 0; i < clen; i++ {
		col := cols[i]
		sm.Set(table.{{.StructName}}.ColumnName2Json[col.Name], m[col.Name])
	}
	return sm
}

//ToCnJSON struct转json，key被替换为字段描述
func (p *{{.StructName}}) ToCnJSON(cols...table.TableField) types.Smap {
	m := p.ToMap()
	clen := len(cols)
	if clen == 0 {
		l := len(table.{{.StructName}}.ColumnNames)
		sm := make(types.Smap, l)
		var cn string
		for i := 0; i < l; i++ {
			cn = table.{{.StructName}}.ColumnNames[i]
			sm.Set(table.{{.StructName}}.ColumnName2Comment[cn], m[cn])
		}
		return sm
	}

	sm := make(types.Smap, clen)
	for i := 0; i < clen; i++ {
		col := cols[i]
		sm.Set(table.{{.StructName}}.ColumnName2Comment[col.Name], m[col.Name])
	}
	return sm
}

//TranslateJSON 将json格式对象的key从列名转为列描述
func (p *{{.StructName}}) TranslateJSON(bean interface{}) (types.Smap, error) {
	var m types.Smap
	if s, ok := bean.(map[string]interface{}); ok {
		m = s
	} else if s, ok := bean.(types.Smap); ok {
		m = s
	} else if s, ok := bean.(string); ok {
		e := utils.JSON.Unmarshal([]byte(s), &m)
		if e != nil {
			return nil, e
		}
	} else if s, ok := bean.(*{{.StructName}}); ok {
		m = s.ToMap()
	} else {
		return nil, Err_Type
	}
	sm := types.Smap{}
	for k, v := range m {
		sm.Set(table.{{.StructName}}.ColumnName2Comment[k], v)
	}
	return sm, Err_Type
}

//SliceToJSON slice转json
func (p *{{.StructName}}) SliceToJSON(sls []*{{.StructName}},cols...table.TableField) []types.Smap {
	slen := len(sls)
	clen := len(cols)
	ms := make([]types.Smap, 0, slen)
	if clen == 0 {
		var (
			cn string
			sm types.Smap
			m map[string]interface{}
		)
		l := len(table.{{.StructName}}.ColumnNames)
		for i := 0; i < slen; i++ {
			m = sls[i].ToMap()
			sm = make(types.Smap, l)
			for i := 0; i < l; i++ {
				cn = table.{{.StructName}}.ColumnNames[i]
				sm.Set(table.{{.StructName}}.ColumnName2Json[cn], m[cn])
			}
			
			ms = append(ms, sm)
		}
		return ms
	}
	for i := 0; i < slen; i++ {
		m := sls[i].ToMap()
		sm := make(types.Smap, clen)
		for i := 0; i < clen; i++ {
			col := cols[i]
			sm.Set(table.{{.StructName}}.ColumnName2Json[col.Name], m[col.Name])
		}

		ms = append(ms, sm)
	}
	return ms
}

	`

func (d *TempData) writeToModel(fileName string) error {
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
			case "bool":
				ret = `false`
			default:
				ret = 0
			}
			return ret
		},
		//"marshal": JSONValue,
	}

	e := template.Must(template.New("tableTpl").Funcs(funcMap).Parse(model_str)).Execute(&buf, d)
	if e != nil {
		showError(e)
		return e
	}

	absPath, _ := filepath.Abs(fileName)
	fileName = filepath.Join(filepath.Dir(absPath), "zzz_"+d.StructName+".go")

	f, e := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer f.Close()

	//e = f.Truncate(0)
	//if e != nil {
	//	showError(e.Error())
	//	return e
	//}

	_, e = f.Write(buf.Bytes())
	if e != nil {
		showError(e)
		return e
	}

	return nil
}

//
////JSONValue redis存储序列化数据转换方法
//func JSONValue(t []string, v string) string {
//	//return `utils.JSONValue(p.` + v + `)`
//	switch t[2] {
//	case "string":
//		return `"\"" + strings.ReplaceAll(p.` + v + `, "\"", "\\\"") + "\""`
//	case "int", "int8", "int16", "int32", "int64":
//		return `strconv.FormatInt(int64(p.` + v + `), 10)`
//	case "uint", "uint8", "uint16", "uint32", "uint64":
//		return `strconv.FormatUint(uint64(p.` + v + `), 10)`
//	case "float32", "float64":
//		return `strconv.FormatFloat(float64(p.` + v + `), 'f', -1, 64)`
//	case "time.Time":
//		return `"\"" + p.` + v + `.Format(time.RFC3339Nano) + "\""`
//	case "bool":
//		return `strconv.FormatBool(p.` + v + `)`
//	case "[]byte":
//		return `"\"" + base64.StdEncoding.EncodeToString(p.` + v + `) + "\""`
//	case "types.BigUint":
//		return `p.` + v + `.String()`
//	case "types.Money":
//		return `strconv.FormatFloat((float64(p.` + v + `) / 100), 'f', -1, 64)`
//	default:
//		return `utils.JSONValue(p.` + v + `)`
//	}
//}

////MarshalJSON
//func (p *{{.StructName}}) MarshalJSON() ([]byte, error) {
//	var buf bytes.Buffer
//	buf.WriteByte('{')
//	{{range $key, $value := .Columns}}buf.WriteString("\"{{index $value 1}}\":" + {{marshal $value $key}} + ",")
//	{{end}}
//	return append(buf.Bytes()[:buf.Len()-1], '}'), nil
//}

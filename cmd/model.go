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
	{{if .HasTime}}"time"{{end}}
	"sync"
	"{{.Module}}/table"
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
	{{range $key, $value := .Columns}}p.{{$key}} = {{getTypeValue $value}}				
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}

//TableName
func (*{{.StructName}}) TableName() string {
	return table.{{.StructName}}.TableName
}

//Insert
func (p *{{.StructName}}) Insert() (int64,error) {
	return 0,nil
}

//Update
func (p *{{.StructName}}) Update() (int64,error) {
	return 0,nil
}

//Delete
func (p *{{.StructName}}) Delete() (int64,error) {
	return 0,nil
}

//func (p *{{.StructName}}) ToMap() map[string]interface{} {
//	m := make(map[string]interface{}, {{len .Columns}})
//	{{range $key, $value := .Columns}}m[table.{{$.StructName}}.{{$key}}.Name] = p.{{$key}}
//	{{end}}
//	return m
//}
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
			default:
				ret = 0
			}
			return ret
		},
	}

	err := template.Must(template.New("tableTpl").Funcs(funcMap).Parse(model_str)).Execute(&buf, d)
	if err != nil {
		showError(err)
		return err
	}

	absPath, _ := filepath.Abs(fileName)
	//fileName = filepath.Join(filepath.Dir(absPath), getBaseFilename(d.FileName)+"_"+d.StructName+"_sorm.go")
	fileName = filepath.Join(filepath.Dir(absPath), "zzz_"+d.StructName+".go")

	var (
		file *os.File
	)

	file, err = os.Create(fileName)

	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		showError(err)
		return err
	}

	return nil
}

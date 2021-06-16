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
	"context"
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

//这里是cache的示例，建议在其他的文件中实现
//var (
//	{{lower .StructName}}_cache     = redis.C(conf.App.Mode).Expiration(time.Minute)
//	{{lower .StructName}}_ids_cache = redis.C(conf.App.Mode, "ids").Expiration(time.Minute)
//)
//
//func init() {
//	{{lower .StructName}}_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
//		id := utils.Interface2UInt64(k, 0)
//		if id < 1 {
//			return nil, errors.New("key is invalid")
//		}
//
//		m := New{{.StructName}}()
//		db := lib.DB().Table(table.{{.StructName}}.TableName)
//		has, e := db.ID(id).Get(m)
//		if has {
//			return m, nil
//		}
//		if e != nil {
//			log.Logs.DBError(db, e)
//		}
//		return nil, e
//	})
//	//
//	sharemp_ids_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
//		//todo:实现从数据库拉取数据的逻辑
//		return nil,nil
//	})
//}

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
func (p *{{.StructName}}) Insert(db lib.Session, cols ...string) (int64,error) {
	if len(cols) == 0 {
		return db.InsertOne(p)
	}
	return db.Cols(cols...).InsertOne(p)
}

//Update
func (p *{{.StructName}}) Update(db lib.Session, id uint64, bean ...interface{}) (int64,error) {
	if len(bean) == 0 {
		return db.ID(id).Update(p)
	}
	return db.ID(id).Update(bean[0])
}

//Delete
func (p *{{.StructName}}) Delete(db lib.Session, id uint64) (int64,error) {
	return  db.ID(id).Delete(p)
}

//Get
func (p *{{.StructName}}) Get(db Session,id uint64) (bool, error) {
	cm, e := {{lower .StructName}}_cache.Get(context.TODO(), id)
	if e != nil {
		return false, e
	}
	if val, ok := cm.(*{{.StructName}}); ok {
		*p = *val
		return ok, nil
	}
	return false, errors.New("类型错误")
}

//FindIDs
//args: size,index
func (p *{{.StructName}}) FindIDs(db Session,query string, vals []interface{}, args ...int) ([]uint64, error) {
	ids := make([]uint64, 0)
	db.Where(query, vals...)

	if len(args) > 0 {
		if len(args) > 1 {
			db.Limit(args[0], args[1]*args[0])
		} else {
			db.Limit(args[0])
		}
	}
	e := db.Find(&ids)
	return ids, e
}

//Find
//args: size,index
func (p *{{.StructName}}) Find(db Session, query string, vals []interface{}, args ...int) ([]*{{.StructName}}, error) {
	ids, e := p.FindIDs(db, query, vals, args...)
	if len(ids) == 0 {
		return nil, e
	}
	list := make([]*{{.StructName}}, 0, len(ids))
	for _, id := range ids {
		m := New{{.StructName}}()
		b, _ := m.Get(db, id)
		if b {
			list = append(list, m)
		}
	}
	return list, nil
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

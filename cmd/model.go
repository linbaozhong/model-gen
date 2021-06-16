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
	"dao/lib"
	"errors"
	"internal/log"
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
//	}).DeserializeModel(func() interface{} {
//		return New{{.StructName}}()
//	})
//	//
//	{{lower .StructName}}_ids_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
//		key, ok := k.(*lib.CacheKey)
//		if !ok {
//			return nil, errors.New("key is invalid")
//		}
//
//		db := lib.DB().Where(key.Query, key.Vals...).Limit(Max_Size_Limit)
//		ids := make([]uint64, 0)
//
//		e := db.Find(&ids)
//		return ids, e
//	}).DeserializeFunc(func(bean interface{}) (interface{}, error) {
//		list := make([]uint64, 0)
//		e := utils.JSON.UnmarshalFromString(bean.(string), list)
//		return list, e
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
func (p *{{.StructName}}) Insert(db Session, cols ...string) (int64,error) {
	if len(cols) > 0 {
		db.Cols(cols...)
	}
	i64,e := db.InsertOne(p)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return i64, e
}

//Update
func (p *{{.StructName}}) Update(db Session, id uint64, bean ...interface{}) (int64,error) {
	var (
		i64 int64
		e error
	)
	if len(bean) == 0 {
		i64,e =  db.ID(id).Update(p)
	} else {
		i64,e = db.ID(id).Update(bean[0])
	}
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return i64, e
}

//Delete
func (p *{{.StructName}}) Delete(db Session, id uint64) (int64,error) {
	i64,e := db.ID(id).Delete(p)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return i64, e
}

//Get
func (p *{{.StructName}}) Get(db Session,id uint64) (bool, error) {
	cm, e := {{lower .StructName}}_cache.Get(context.TODO(), id)
	if e != nil {
		log.Logs.Error(e)
		return false, e
	}
	if val, ok := cm.(*{{.StructName}}); ok {
		*p = *val
		return ok, nil
	}

	e = errors.New("类型错误")
	log.Logs.Error(e)
	return false, e
}

////FindIDs
//func (p *{{.StructName}}) FindIDs(db Session,query string, vals []interface{}, args ...int) ([]uint64, error) {
//	ids := make([]uint64, 0)
//	db.Where(query, vals...)
//
//	if len(args) > 0 {
//		if len(args) > 1 {
//			db.Limit(args[0], args[1]*args[0])
//		} else {
//			db.Limit(args[0])
//		}
//	}
//	e := db.Find(&ids)
//	return ids, e
//}

//Find
//args: size,index
func (p *{{.StructName}}) Find(db Session, query string, vals []interface{}, size, index int) ([]*{{.StructName}}, error) {
	k := lib.NewCacheKey(query,vals)

	ids, e := {{lower .StructName}}_ids_cache.LGet(context.TODO(), k, int64(size*index), int64(size*(index+1)))
	if len(ids) == 0 {
		log.Logs.Error(e)
		return nil, e
	}

	ms, e := {{lower .StructName}}_cache.Gets(context.TODO(), ids)
	if e != nil {
		log.Logs.Error(e)
		return nil, e
	}
	list := make([]*{{.StructName}}, 0, len(ms))
	for _, m := range ms {
		if mm, ok := m.(*{{.StructName}}); ok {
			list = append(list, mm)
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

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
	{{if or .HasTime .HasCache}}"time"{{end}}
{{if .HasPrimaryKey}}
	{{if .HasCache}}"context"
	"internal/cache/redis"
	"internal/conf"
	"libs/utils"{{end}}
	"internal/log"
	"libs/types"
	"{{.ModulePath}}/table"
{{end}}
)

var (
	{{lower .StructName}}Pool = sync.Pool{
		New: func() interface{} {
			return &{{.StructName}}{}
		},
	}
)

{{if and .HasCache .HasPrimaryKey}}
var (
	{{lower .StructName}}_cache     = redis.NewClient(conf.App.Mode,"{{lower .StructName}}").Expiration({{.CacheData}})
	{{lower .StructName}}_ids_cache = redis.NewClient(conf.App.Mode, "{{lower .StructName}}_ids").Expiration({{.CacheList}})
)

func init() {
	{{lower .StructName}}_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
		id := utils.Interface2Uint64(k, 0)
		if id < 1 {
			return nil, InvalidKey
		}

		m := New{{.StructName}}()
		db := Db().Table(table.{{.StructName}}.TableName)
		has, e := db.Where(table.{{.StructName}}.PrimaryKey.Eq(),id).
			Get(m)
		if has {
			return m, nil
		}
		if e != nil {
			log.Logs.DBError(db, e)
		}
		return nil, e
	}).DeserializeModel(func() interface{} {
		return New{{.StructName}}()
	})
	//
	{{lower .StructName}}_ids_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
		cond, ok := k.(table.ISqlBuilder)
		if !ok {
			return nil, InvalidKey
		}
		
		query, args := cond.GetCondition()
		
		db := Db().Table(table.{{.StructName}}.TableName).
			Cols(table.{{.StructName}}.PrimaryKey.Name).
			Limit({{.CacheLimit}})
		if query != ""{
			db.Where(query, args...)
		}
		ids := make([]uint64, 0)
		e := db.Find(&ids)
		if e != nil {
			log.Logs.DBError(db, e)
		}
		return ids, e
	}).DeserializeFunc(func(bean interface{}) (interface{}, error) {
		ids := make([]uint64, 0)
		e := utils.JSON.UnmarshalFromString(utils.Interface2String(bean), &ids)
		if e != nil {
			log.Logs.Error(e)
		}
		return ids, e
	})
}
{{end}}

func New{{.StructName}}() *{{.StructName}} {
	return {{lower .StructName}}Pool.Get().(*{{.StructName}})
}

//Free 
func (p *{{.StructName}}) Free() {
	{{range $key, $value := .Columns}}p.{{$key}} = {{getTypeValue $value}}				
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}
{{if .HasPrimaryKey}}
//TableName
func (*{{.StructName}}) TableName() string {
	return table.{{.StructName}}.TableName
}

//Insert 
func (p *{{.StructName}}) Insert(db types.Session, cols ...string) (int64,error) {
	if len(cols) > 0 {
		db.Cols(cols...)
	}

	i64,e := db.Table(table.{{.StructName}}.TableName).InsertOne(p)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		{{lower .StructName}}_ids_cache.Empty(context.TODO())
	}
{{end}}
	return i64, e
}

//InsertBatch
func (p *{{.StructName}}) InsertBatch(db types.Session, beans []interface{}, cols ...string) (int64, error) {
	if len(cols) > 0 {
		db.Cols(cols...)
	}

	i64, e := db.Table(table.{{.StructName}}.TableName).Insert(beans...)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		{{lower .StructName}}_ids_cache.Empty(context.TODO())
	}
{{end}}
	return i64, e
}

//Update
func (p *{{.StructName}}) Update(db types.Session, id uint64, bean ...interface{}) (int64,error) {
	var (
		i64 int64
		e error
	)

	db.Table(table.{{.StructName}}.TableName).
		Where(table.{{.StructName}}.PrimaryKey.Eq(),id)

	if len(bean) == 0 {
		i64,e = db.Update(p)
	} else {
		i64,e = db.Update(bean[0])
	}

	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnChange(id)
	}
{{end}}
	return i64, e
}

//UpdateBatch
func (p *{{.StructName}}) UpdateBatch(db types.Session, cond table.ISqlBuilder, bean ...interface{}) (int64, error) {
	var (
		i64 int64
		e   error
	)

	db.Table(table.{{.StructName}}.TableName)
	if cond != nil {
		query, args := cond.GetCondition()
		if query != "" {
			db.Where(query, args...)
		}
	}

	if len(bean) == 0 {
		i64, e = db.Update(p)
	} else {
		i64, e = db.Update(bean[0])
	}

	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnBatchChange(cond)
	}
{{end}}
	return i64, e
}

//Delete
func (p *{{.StructName}}) Delete(db types.Session, id uint64) (int64,error) {
	i64,e := db.Table(table.{{.StructName}}.TableName).
		Where(table.{{.StructName}}.PrimaryKey.Eq(),id).
		Delete(p)

	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnChange(id)
	}
{{end}}
	return i64, e
}

//DeleteBatch
func (p *{{.StructName}}) DeleteBatch(db types.Session, cond table.ISqlBuilder) (int64, error) {
	db.Table(table.{{.StructName}}.TableName)
	if cond != nil {
		query, args := cond.GetCondition()
		if query != "" {
			db.Where(query, args...)
		}
	}
	i64, e := db.Delete(p)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnBatchChange(cond)
	}
{{end}}
	return i64, e
}

//Get
func (p *{{.StructName}}) Get(db types.Session,id uint64) (bool, error) {
{{if .HasCache}}
	cm, e := {{lower .StructName}}_cache.Get(context.TODO(), id)
	if e != nil {
		log.Logs.Error(e)
		return false, e
	}
	if val, ok := cm.(*{{.StructName}}); ok {
		*p = *val
		return ok, nil
	}

	log.Logs.Error(Err_Type)
	return false, e
{{else}}
	has, e := db.Table(table.{{.StructName}}.TableName).Where(table.{{.StructName}}.PrimaryKey.Eq(),id).
		Get(p)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return has,e
{{end}}
}

//Find
func (p *{{.StructName}}) Find(db types.Session, cond table.ISqlBuilder, size, index int) ([]*{{.StructName}}, error) {
{{if .HasCache}}
	ids, e := {{lower .StructName}}_ids_cache.LGet(context.TODO(), cond, int64(size*(index-1)), int64(size*index))
	if len(ids) == 0 {
		log.Logs.Error(e)
		return nil, e
	}

	ms, e := {{lower .StructName}}_cache.Gets(context.TODO(), ids...)
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
{{else}}
	list := make([]*{{.StructName}}, 0)

	db.Table(table.{{.StructName}}.TableName)
	if cond != nil {
		query, args := cond.GetCondition()
		if query != "" {
			db.Where(query, args...)
		}
	}
	e := db.Limit(size,size*(index-1)).
		Find(&list)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{end}}
	return list, nil
}

//ToMap
func (p *{{.StructName}}) ToMap(cols...table.TableField) map[string]interface{} {
	if len(cols) == 0{
		return map[string]interface{}{
			{{range $key, $value := .Columns}}table.{{$.StructName}}.{{$key}}.Name:p.{{$key}},
			{{end}}
		}
	}

	m := make(map[string]interface{},len(cols))
	for _, col := range cols {
		switch col.Name {
		{{range $key, $value := .Columns}}case table.{{$.StructName}}.{{$key}}.Name:
			m[col.Name] = p.{{$key}}
		{{end}}
		}
	}
	return m
}

//ToJSON
func (p *{{.StructName}}) ToJSON(cols...table.TableField) types.Smap {
	if len(cols) == 0{
		return types.Smap{
			{{range $key, $value := .Columns}}table.{{$.StructName}}.{{$key}}.Json:p.{{$key}},
			{{end}}
		}
	}

	m := make(types.Smap,len(cols))
	for _, col := range cols {
		switch col.Json {
		{{range $key, $value := .Columns}}case table.{{$.StructName}}.{{$key}}.Json:
			m[col.Json] = p.{{$key}}
		{{end}}
		}
	}
	return m
}

{{if .HasCache}}
//OnChange
func (p *{{.StructName}}) OnChange(id uint64) error {
	return {{lower .StructName}}_cache.Remove(context.TODO(), id)
}

//OnBatchChange
func (p *{{.StructName}}) OnBatchChange(cond table.ISqlBuilder) {
	ids := make([]interface{}, 0)

	query, args := cond.GetCondition()
	db := Db().Table(table.{{.StructName}}.TableName).Where(query, args...)
	e := db.Find(&ids)

	if e != nil {
		log.Logs.DBError(db, e)
	}

	if len(ids) > 0 {
		{{lower .StructName}}_cache.Remove(context.TODO(), ids...)
	}
}
func {{.StructName}}Cache() *redis.RedisBroker {
	return {{lower .StructName}}_cache
}

func {{.StructName}}IDsCache() *redis.RedisBroker {
	return {{lower .StructName}}_ids_cache
}
{{end}}

//func (p *{{.StructName}}) getInsert(cols ...string) (sql string, params []interface{}, e error) {
//	sb := table.NewSqlBuilder()
//	defer sb.Free()
//
//	sb.Table(p)
//
//	m := p.ToMap(cols...)
//	for k, v := range m {
//		sb.Set(k, v)
//	}
//
//	sql, params, e = sb.Insert()
//	return
//}


{{end}}
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
	}

	err := template.Must(template.New("tableTpl").Funcs(funcMap).Parse(model_str)).Execute(&buf, d)
	if err != nil {
		showError(err)
		return err
	}

	absPath, _ := filepath.Abs(fileName)
	fileName = filepath.Join(filepath.Dir(absPath), "zzz_"+d.StructName+".go")

	////文件已存在
	//_, e := os.Stat(fileName)
	//if e == nil {
	//	return nil
	//}
	var (
		f *os.File
	)

	f, err = os.Create(fileName)

	if err != nil {
		showError(err.Error())
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		showError(err)
		return err
	}

	return nil
}

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var dao_str = `
package dao

import (
{{if .HasPrimaryKey}}
	{{if .HasCache}}"time"
	"context"
	"libs/utils"
	"internal/cache/redis"
	"internal/conf"{{end}}
{{end}}
	"internal/log"
	"libs/types"
	"{{.ModulePath}}"
	"{{.ModulePath}}/table"
)

type {{lower .StructName}} struct {
}

var (
	{{.StructName}} {{lower .StructName}}
)

{{if and .HasCache .HasPrimaryKey}}
var (
	{{lower .StructName}}_cache     = redis.NewClient(conf.App.Mode,"{{lower .StructName}}").Expiration({{.CacheData}})
	{{lower .StructName}}_ids_cache = redis.NewClient(conf.App.Mode, "{{lower .StructName}}_ids").Expiration({{.CacheList}})
	{{lower .StructName}}_count_cache = redis.NewClient(conf.App.Mode, "{{lower .StructName}}_count").Expiration({{.CacheList}})
)

func init() {
	{{lower .StructName}}_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
		if k == nil {
			return nil, InvalidKey
		}

		m := models.New{{.StructName}}()
		db := models.Db().Table(table.{{.StructName}}.TableName)
		has, e := db.Where(table.{{.StructName}}.PrimaryKey.Eq(),k).
			Get(m)
		if has {
			return m, nil
		}
		m.Free()
		if e != nil {
			log.Logs.DBError(db, e)
		}
		return nil, e
	}).DeserializeModel(func() interface{} {
		return models.New{{.StructName}}()
	})
	//
	{{lower .StructName}}_ids_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
		cond, ok := k.(table.ISqlBuilder)
		if !ok {
			return nil, InvalidKey
		}
		
		db := models.Db().Table(table.{{.StructName}}.TableName)
		db.Cols(table.{{.StructName}}.PrimaryKey.Quote())

		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}
		if s := cond.GetOrderBy(); s != "" {
			db.OrderBy(s)
		}
		if size, start := cond.GetLimit(); size > 0 {
			db.Limit(size, start)
		} else {
			db.Limit({{.CacheLimit}})
		}

		ids := make([]interface{}, 0)
		e := db.Find(&ids)
		if e != nil {
			log.Logs.DBError(db, e)
		}
		return ids, e
	})
	//
	{{lower .StructName}}_count_cache.LoaderFunc(func(k interface{}) (interface{}, error) {
		cond, ok := k.(table.ISqlBuilder)
		if !ok {
			return nil, InvalidKey
		}
		
		db := models.Db().Table(table.{{.StructName}}.TableName)

		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}

		i64, e := db.Count()
		if e != nil {
			log.Logs.DBError(db, e)
		}
		return i64, e
	})
}
{{end}}

{{if .HasPrimaryKey}}
//Insert 新增一条数据
func (p {{lower .StructName}}) Insert(x interface{}, bean *models.{{.StructName}}, cols ...string) (int64,error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	if len(cols) > 0 {
		db.Cols(cols...)
	}

	i64, e := db.InsertOne(bean)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnListChange()
	}
{{end}}
	return i64, e
}

//InsertBatch 批量新增数据
func (p {{lower .StructName}}) InsertBatch(x interface{}, beans []*models.{{.StructName}}, cols ...string) (int64, error) {
	l := len(beans)
	if l == 0 {
		return 0, Err_Type
	}
	db := getDB(x, table.{{.StructName}}.TableName)

	if len(cols) > 0 {
		db.Cols(cols...)
	}
	ibeans := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ibeans[i] = beans[i]
	}
	i64, e := db.Insert(ibeans...)
	if e != nil {
		log.Logs.DBError(db, e)
	}
{{if .HasCache}}
	if i64 > 0 {
		p.OnListChange()
	}
{{end}}
	return i64, e
}

//Update 根据主键修改一条数据
func (p {{lower .StructName}}) Update(x interface{}, id types.BigUint, bean interface{}) (int64,error) {
	var (
		i64   int64
		e     error
	)

	db := getDB(x, table.{{.StructName}}.TableName)
	db.Where(table.{{.StructName}}.PrimaryKey.Eq(), id).
		Limit(1)
	i64, e = db.Update(bean)
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

//UpdateBatch 根据cond条件批量修改数据
func (p {{lower .StructName}}) UpdateBatch(x interface{}, cond table.ISqlBuilder, bean interface{}) (int64, error) {
	var (
		i64 int64
		e   error
	)
	
	db := getDB(x, table.{{.StructName}}.TableName)
	if cond != nil {
		if cols := cond.GetCols(); len(cols) > 0 {
			db.Cols(cols...)
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if size, start := cond.GetLimit(); size > 0 {
			db.Limit(size, start)
		}
	}
	i64, e = db.Update(bean)
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

//Delete 根据主键删除一条数据
func (p {{lower .StructName}}) Delete(x interface{}, id types.BigUint) (int64,error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	i64,e := db.Where(table.{{.StructName}}.PrimaryKey.Eq(),id).
		Limit(1).
		Delete()

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

//DeleteBatch 根据cond条件批量删除数据
func (p {{lower .StructName}}) DeleteBatch(x interface{}, cond table.ISqlBuilder) (int64, error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	if cond != nil {
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if size, start := cond.GetLimit(); size > 0 {
			db.Limit(size, start)
		}
	}
	i64, e := db.Delete()
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

//Get 根据主键从Cache中获取一条数据
func (p {{lower .StructName}}) Get(x interface{},id types.BigUint) (*models.{{.StructName}}, error) {
{{if .HasCache}}
	cm, e := {{lower .StructName}}_cache.Get(context.TODO(), id)
	if e != nil {
		log.Logs.Error(e)
		return nil, e
	}
	if val, ok := cm.(*models.{{.StructName}}); ok {
		return val, nil
	}

	log.Logs.Error(Err_Type)
	return nil, Err_Type
{{else}}
	return p.GetNoCache(x,id)
{{end}}
}

//GetNoCache 根据主键从数据库中获取一条数据
func (p {{lower .StructName}}) GetNoCache(x interface{},id types.BigUint, cols ...table.TableField) (*models.{{.StructName}},error) {
	var bean = models.New{{.StructName}}()
	db := getDB(x, table.{{.StructName}}.TableName)
	//
	if len(cols) > 0 {
		_cols := make([]string, 0, len(cols))
		for _, col := range cols {
			_cols = append(_cols, col.Name)
		}
		db.Cols(_cols...)
	}

	has, e := db.Where(table.{{.StructName}}.PrimaryKey.Eq(),id).Limit(1).
		Get(bean)
	if has {
		return bean, nil
	}
	bean.Free()
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return nil, e
}

//IDs 根据cond条件从cache中获取主键slice
func (p {{lower .StructName}}) IDs(x interface{}, cond table.ISqlBuilder, size, index int) ([]interface{}, error) {
{{if .HasCache}}
	if size == 0 {
		size = {{.CacheLimit}}
	}

	if index == 0 {
		index = 1
	}
	return {{lower .StructName}}_ids_cache.LGet(context.TODO(), cond, int64(size*(index-1)), int64(size*index))
{{else}}
	return p.IDsNoCache(x,cond,size,index)
{{end}}
}

//IDsNoCache 根据cond条件从数据库中获取主键slice
func (p {{lower .StructName}}) IDsNoCache(x interface{}, cond table.ISqlBuilder, size, index int) ([]interface{}, error) {
	return getColumn(x,table.{{.StructName}}.TableName, table.{{.StructName}}.PrimaryKey, cond, size, index)
}

//GetColumn 根据cond条件从数据库中单列slice
func (p {{lower .StructName}}) GetColumn(x interface{}, col table.TableField, cond table.ISqlBuilder, size, index int) ([]interface{}, error) {
	return getColumn(x,table.{{.StructName}}.TableName, col, cond, size, index)
}

//Sum 对某个字段进行求和
func (p {{lower .StructName}}) Sum(x interface{}, cond table.ISqlBuilder, col table.TableField) (float64, error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	if cond != nil {
		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}
	}

	sum, e := db.Sum(p, col.Name)
	if e != nil {
		log.Logs.Error(e)
		return 0, e
	}
	return sum, nil
}

//Sums 对某几个字段进行求和
func (p {{lower .StructName}}) Sums(x interface{}, cond table.ISqlBuilder, args ...table.TableField) ([]float64, error) {
	if len(args) == 0 {
		return nil, Param_Missing
	}
	
	cols := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		cols[i] = args[i].Name
	}

	db := getDB(x, table.{{.StructName}}.TableName)

	if cond != nil {
		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}
	}

	sums, e := db.Sums(p, cols...)
	if e != nil {
		log.Logs.Error(e)
		return nil, e
	}
	return sums, nil
}

//Count 根据cond条件从cache中获取数据总数
func (p {{lower .StructName}}) Count(x interface{}, cond table.ISqlBuilder) (int64, error) {
{{if .HasCache}}
	i, e := {{lower .StructName}}_count_cache.Get(context.TODO(), cond)
	if e != nil {
		log.Logs.Error(e)
		return 0, e
	}
	return utils.Interface2Int64(i), nil
{{else}}
	return p.CountNoCache(x,cond)
{{end}}
}


//CoundNoCache 根据cond条件从数据库中获取数据列表
func (p {{lower .StructName}}) CountNoCache(x interface{}, cond table.ISqlBuilder) (int64, error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	if cond != nil {
		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}
	}
	i64, e := db.Count()
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return i64, nil
}

// Gets
func (p {{lower .StructName}}) Gets(x interface{}, ids []interface{}) ([]*models.{{.StructName}}, error) {
{{if .HasCache}}
	if len(ids) == 0 {
		return nil, nil
	}
	ms, e := {{lower .StructName}}_cache.Gets(context.TODO(), ids...)
	if e != nil {
		log.Logs.Error(e)
		return nil, e
	}
	list := make([]*models.{{.StructName}}, 0, len(ms))
	for _, m := range ms {
		if mm, ok := m.(*models.{{.StructName}}); ok {
			list = append(list, mm)
		}
	}
	return list, nil
{{else}}
	return p.GetsNoCache(x, ids)
{{end}}
}

// GetsNoCache
func (p {{lower .StructName}}) GetsNoCache(x interface{}, ids []interface{}) ([]*models.{{.StructName}}, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	db := getDB(x, table.{{.StructName}}.TableName)

	list := make([]*models.{{.StructName}}, 0)
	e := db.In(table.{{.StructName}}.PrimaryKey.Name,ids...).Find(&list)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return list, nil
}

//Find 根据cond条件从cache中获取数据列表
func (p {{lower .StructName}}) Find(x interface{}, cond table.ISqlBuilder, size, index int) ([]*models.{{.StructName}}, error) {
{{if .HasCache}}
	ids, e := p.IDs(x,cond,size,index)
	if len(ids) == 0 {
		return nil, e
	}

	return p.Gets(x, ids)
{{else}}
	return p.FindNoCache(x,cond,size,index)
{{end}}
}

//FindNoCache 根据cond条件从数据库中获取数据列表
func (p {{lower .StructName}}) FindNoCache(x interface{}, cond table.ISqlBuilder, size, index int) ([]*models.{{.StructName}}, error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	list := make([]*models.{{.StructName}}, 0)

	if cond == nil {
		if size > 0 {
			if index == 0 {
				index = 1
			}
			db.Limit(size, size*(index-1))
		}
	} else {
		if joins := cond.GetJoin(); len(joins) > 0 {
			for _, join := range joins {
				db.Join(join[0], join[1], join[2])
			}
		}
		if cols := cond.GetCols(); len(cols) > 0 {
			db.Cols(cols...)
		}
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
		if s := cond.GetGroupBy(); s != "" {
			db.GroupBy(s)
		}
		if s := cond.GetHaving(); s != "" {
			db.Having(s)
		}
		if s := cond.GetOrderBy(); s != "" {
			db.OrderBy(s)
		}
		if size > 0 {
			if index == 0 {
				index = 1
			}
			db.Limit(size, size*(index-1))
		} else if i, start := cond.GetLimit(); i > 0 {
			db.Limit(i, start)
		}
	}

	e := db.Find(&list)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return list, nil
}

//FindOne 根据cond条件从cache中获取一条数据
func (p {{lower .StructName}}) FindOne(x interface{}, cond table.ISqlBuilder) (*models.{{.StructName}}, error) {
	if cond != nil {
		cond.Limit(1)
	}
	f, e := p.Find(x, cond, 1, 1)
	if e != nil {
		return nil,e
	}
	if len(f) > 0 {
		return f[0],nil
	}
	return nil, nil
}

//FindOneNoCache 根据cond条件从数据库中获取一条数据
func (p {{lower .StructName}}) FindOneNoCache(x interface{}, cond table.ISqlBuilder) (*models.{{.StructName}},error) {
	if cond != nil {
		cond.Limit(1)
	}
	f, e := p.FindNoCache(x, cond, 1, 1)
	if e != nil {
		return nil,e
	}
	if len(f) > 0 {
		return f[0],nil
	}
	return nil, nil
}

//Exists 是否存在符合条件cond的记录
func (p {{lower .StructName}}) Exists(x interface{}, cond table.ISqlBuilder) (bool, error) {
	db := getDB(x, table.{{.StructName}}.TableName)

	db.Cols(table.{{.StructName}}.PrimaryKey.Name)
	if cond != nil {
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
	}
	var bean = models.New{{.StructName}}()
	defer bean.Free()
	has, e := db.Limit(1).Get(bean)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return has, e
}

{{if .HasCache}}
//OnChange
func (p {{lower .StructName}}) OnChange(id types.BigUint) {
	{{lower .StructName}}_cache.Remove(context.TODO(), id)
	//p.OnListChange()
}

//OnBatchChange
func (p {{lower .StructName}}) OnBatchChange(cond table.ISqlBuilder) {
	db := models.Db().Table(table.{{.StructName}}.TableName).
			Cols(table.{{.StructName}}.PrimaryKey.Quote())
	if cond != nil {
		if s, args := cond.GetWhere(); s != "" {
			db.Where(s, args...)
		}
	}
	ids := make([]interface{}, 0)
	e := db.Find(&ids)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	if len(ids) > 0 {
		{{lower .StructName}}_cache.Remove(context.TODO(), ids...)
		//p.OnListChange()
	}
}
//OnListChange
func (p {{lower .StructName}}) OnListChange() {
	{{lower .StructName}}_ids_cache.Empty(context.TODO())
	{{lower .StructName}}_count_cache.Empty(context.TODO())
}

func (p {{lower .StructName}})Cache() *redis.RedisBroker {
	return {{lower .StructName}}_cache
}

func (p {{lower .StructName}})IDsCache() *redis.RedisBroker {
	return {{lower .StructName}}_ids_cache
}
{{end}}
{{end}}

`

func (d *TempData) writeToDao(fileName string) error {
	if !d.HasPrimaryKey {
		return nil
	}

	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	err := template.Must(template.New("daoTpl").Funcs(funcMap).Parse(dao_str)).Execute(&buf, d)
	if err != nil {
		showError(err)
		return err
	}

	absPath, _ := filepath.Abs(fileName)
	fileName = filepath.Join(filepath.Dir(absPath), "dao", d.StructName+"_dao.go")

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
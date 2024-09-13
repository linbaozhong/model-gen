package cmd

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

func writeDaoBaseFile(filename string, d *TempData) error {
	baseFilename, _ := filepath.Abs(filename)

	f, e := os.OpenFile(baseFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer f.Close()

	// e = f.Truncate(0)
	// if e != nil {
	//	showError(e.Error())
	//	return e
	// }

	var buf bytes.Buffer
	_ = template.Must(template.New("daoBaseTpl").Parse(daoBaseTpl)).Execute(&buf, d)
	formatted, _ := format.Source(buf.Bytes())
	_, e = f.Write(formatted)
	if e != nil {
		showError(e.Error())
	}
	return e
}

var daoBaseTpl = `
package dao

import (
	"context"
	"{{.ModulePath}}"
	"database/sql"
	"errors"
	"internal/log"
	"golang.org/x/sync/singleflight"
	jsoniter "github.com/json-iterator/go"
	tbl "github.com/linbaozhong/model-gen/pkg/table"
)

var (
	InvalidKey    = errors.New("The key is invalid")
	Err_Type      = errors.New("The type is wrong")
	Err_NoRows    = sql.ErrNoRows
	Param_Missing = errors.New("Parameters are missing")

	json = jsoniter.ConfigCompatibleWithStandardLibrary
	sg   singleflight.Group
)

//Transaction 事务处理
func Transaction(f func(*{{.PackageName}}.Session) (interface{}, error)) (result interface{}, err error) {
	session := {{.PackageName}}.Db().NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			session.Rollback()
		}
	}()

	result, err = f(session)
	if err != nil {
		return result, err
	}

	if err = session.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

//
func Sync(k string, f func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return sg.Do(k, f)
}

func queryInterfaces(x *{{.PackageName}}.Session, cond tbl.ISqlBuilder) ([]map[string]interface{}, error) {
	sql, e := cond.Select()
	if e != nil {
		return nil, e
	}
	return x.QueryInterface(sql...)
}

func getSession(x interface{}, tablename string) *{{.PackageName}}.Session {
	db, ok := x.(*{{.PackageName}}.Engine)
	if ok && db != nil {
		return db.Table(tablename)
	}
	sess, ok := x.(*{{.PackageName}}.Session)
	if ok && sess != nil {
		return sess.Table(tablename)
	}
	return nil
}
func getDB(x interface{}, tablename string) *{{.PackageName}}.Session {
	sess := getSession(x, tablename)
	if sess != nil {
		return sess
	}
	if ctx, ok := x.(context.Context); ok {
		if db := ctx.Value("db"); db != nil {
			if sess = getSession(db, tablename); sess != nil {
				return sess
			}
		}
	}
	return {{.PackageName}}.Db().Table(tablename)
}

func getColumn(x interface{}, tablename string, col string, cond tbl.ISqlBuilder, size, index int) ([]interface{}, error) {
	db := getDB(x, tablename)

	cls := make([]interface{}, 0)

	db.Cols(col)
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
		if s, args, e := cond.GetWhere(); e != nil {
			return nil, e
		} else if s != "" {
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

	e := db.Find(&cls)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return cls, e
}


func getContext(x interface{}) context.Context {
	if ctx, ok := x.(context.Context); ok {
		return ctx
	}
	return context.Background()
}

		`

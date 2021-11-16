package cmd

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

func writeDaoBaseFile(filename, modulePath string) error {
	baseFilename, _ := filepath.Abs(filename)

	file, err := os.Create(baseFilename)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = template.Must(template.New("daoBaseTpl").Parse(daoBaseTpl)).Execute(&buf, modulePath)
	formatted, _ := format.Source(buf.Bytes())
	_, err = file.Write(formatted)
	if err != nil {
		showError(err.Error())
	}
	return err
}

//
var daoBaseTpl = `
package dao

import (
	"{{.}}"
	"{{.}}/table"
	"errors"
	"internal/log"
)

var (
	InvalidKey    = errors.New("The key is invalid")
	Err_Type      = errors.New("The type is wrong")
	Param_Missing = errors.New("Parameters are missing")
)

//Transaction 事务处理
func Transaction(f func(*models.Session) (interface{}, error)) (interface{}, error) {
	session := models.Db().NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return nil, err
	}

	result, err := f(session)
	if err != nil {
		return result, err
	}

	if err := session.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

func getDB(x interface{}, tablename string) *models.Session {
	db, ok := x.(*models.Session)
	if ok && db != nil {
		return db.Table(tablename)
	}
	return models.Db().Table(tablename)
}

func getColumn(x interface{}, tablename string, col table.TableField, cond table.ISqlBuilder, size, index int) ([]interface{}, error) {
	db := getDB(x, tablename)

	cls := make([]interface{}, 0)

	db.Cols(col.Quote())
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

	e := db.Find(&cls)
	if e != nil {
		log.Logs.DBError(db, e)
	}
	return cls, e
}

		`

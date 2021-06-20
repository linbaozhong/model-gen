package cmd

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

func writeBuildFile(filename string) error {
	buildFilename, _ := filepath.Abs(filename)

	file, err := os.Create(buildFilename)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = template.Must(template.New("buildTpl").Parse(buildTpl)).Execute(&buf, nil)
	formatted, _ := format.Source(buf.Bytes())
	_, err = file.Write(formatted)
	if err != nil {
		showError(err.Error())
	}
	return err
}

var buildTpl = `
package table

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type JoinType string

const (
	Inner_Join JoinType = " INNER JOIN "
	Left_Join  JoinType = " LEFT JOIN "
	Right_Join JoinType = " RIGHT JOIN "

	operator_and = " AND "
	operator_or  = " OR "
	placeholder  = "?"
)

var (
	// ErrTableEmpty table not set
	ErrTableEmpty = errors.New("table empty")
	// ErrInsertEmpty insert content not set
	ErrInsertEmpty = errors.New("insert content empty")
	// ErrUpdateEmpty update content not set
	ErrUpdateEmpty = errors.New("update content empty")
)

type ITableName interface {
	TableName() string
}
type ISqlBuilder interface {
	Table(m interface{}) ISqlBuilder
	GetQuery() (string, []interface{})
	Select() (string, []interface{}, error)
	//Insert() (string, error)

	Cols(args ...TableField) ISqlBuilder
	Eq(f TableField, v interface{}) ISqlBuilder
	Gt(f TableField, v interface{}) ISqlBuilder
	Gte(f TableField, v interface{}) ISqlBuilder
	Lt(f TableField, v interface{}) ISqlBuilder
	Lte(f TableField, v interface{}) ISqlBuilder
	Ue(f TableField, v interface{}) ISqlBuilder
	Bt(f TableField, v1, v2 interface{}) ISqlBuilder
	Like(f TableField, v interface{}) ISqlBuilder
	Llike(f TableField, v interface{}) ISqlBuilder
	Rlike(f TableField, v interface{}) ISqlBuilder
	Null(f TableField) ISqlBuilder
	UnNull(f TableField) ISqlBuilder
	Join(t JoinType, l, r TableField) ISqlBuilder
	Limit(size int, offset ...int) ISqlBuilder

	And() ISqlBuilder
	AndWhere(sb ISqlBuilder) ISqlBuilder

	Or() ISqlBuilder
	OrWhere(sb ISqlBuilder) ISqlBuilder

	Where() string
	Params() []interface{}

	GroupBy(cols ...TableField) ISqlBuilder
	Having(sb ISqlBuilder) ISqlBuilder
	OrderBy(cols ...TableField) ISqlBuilder
	Asc(cols ...TableField) ISqlBuilder
	Desc(cols ...TableField) ISqlBuilder

	Free()
}
type sqlBuilder struct {
	table        string
	cols         []TableField
	where        strings.Builder
	params       []interface{}
	groupBy      strings.Builder
	having       strings.Builder
	havingParams []interface{}
	orderBy      strings.Builder
	limit        string
	join         string

	andOr bool
}

var (
	sqlBuilderPool = sync.Pool{New: func() interface{} {
		return &sqlBuilder{
			andOr: true,
		}
	}}
)

func NewSqlBuilder() *sqlBuilder {
	return sqlBuilderPool.Get().(*sqlBuilder)
}

//Free
func (p *sqlBuilder) Free() {
	p.table = ""
	p.cols = []TableField{}
	p.where.Reset()
	p.params = []interface{}{}
	p.groupBy.Reset()
	p.having.Reset()
	p.havingParams = []interface{}{}
	p.orderBy.Reset()
	p.limit = ""

	p.andOr = true

	sqlBuilderPool.Put(p)
}

//Table
func (p *sqlBuilder) Table(m interface{}) ISqlBuilder {
	if name, ok := m.(string); ok {
		p.table = name
	} else if iface, ok := m.(ITableName); ok {
		p.table = iface.TableName()
	}
	return p
}

//Select
func (p *sqlBuilder) Select() (string, []interface{}, error) {
	if p.table == "" {
		return "", nil, ErrTableEmpty
	}
	var buf strings.Builder
	//SELECT
	buf.WriteString("SELECT ")
	if len(p.cols) == 0 {
		buf.WriteString("*")
	} else {
		for i, col := range p.cols {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(col.Quote())
		}
	}
	//FROM TABLE
	buf.WriteString(" FROM " + Quote_Char + p.table + Quote_Char)
	//JOIN
	if p.join != "" {
		buf.WriteString(p.join)
	}
	//WHERE
	sql, params := p.GetQuery()
	if sql != "" {
		buf.WriteString(" WHERE " + sql)
	}

	return buf.String(), params, nil
}

//GetQuery
func (p *sqlBuilder) GetQuery() (string, []interface{}) {
	defer p.Free()

	var buf strings.Builder
	//WHERE
	if p.where.Len() > 0 {
		buf.WriteString(p.Where())
	}
	//GROUP BY
	if p.groupBy.Len() > 0 {
		buf.WriteString(" GROUP BY " + p.groupBy.String())
	}
	//HAVING
	if p.having.Len() > 0 {
		buf.WriteString(" HAVING " + p.having.String())
	}
	//ORDER BY
	if p.orderBy.Len() > 0 {
		buf.WriteString(" ORDER BY " + p.orderBy.String())
	}
	//LIMIT
	if p.limit != "" {
		buf.WriteString(p.limit)
	}

	return buf.String(), p.Params()
}

//JOIN
func (p *sqlBuilder) Join(t JoinType, l, r TableField) ISqlBuilder {
	p.join = string(t) + Quote_Char + r.Table + Quote_Char + " ON " + r.Quote() + " = " + l.Quote()
	return p
}

//LIMIT
func (p *sqlBuilder) Limit(size int, offset ...int) ISqlBuilder {
	if len(offset) == 0 {
		p.limit = " LIMIT " + strconv.Itoa(size)
	} else {
		p.limit = " LIMIT " + strconv.Itoa(size) + " OFFSET " + strconv.Itoa(offset[0])
	}
	return p
}

//GROUPBY
func (p *sqlBuilder) GroupBy(cols ...TableField) ISqlBuilder {
	if len(cols) == 0 {
		return p
	}
	for i, col := range cols {
		if i > 0 {
			p.groupBy.WriteByte(',')
		}
		p.groupBy.WriteString(col.Quote())
	}
	return p
}

//HAVING
func (p *sqlBuilder) Having(sb ISqlBuilder) ISqlBuilder {
	defer sb.Free()
	p.having.WriteString(sb.Where())
	p.havingParams = sb.Params()
	return p
}

//OrderBy
func (p *sqlBuilder) OrderBy(cols ...TableField) ISqlBuilder {
	return p.Asc(cols...)
}

//Asc
func (p *sqlBuilder) Asc(cols ...TableField) ISqlBuilder {
	if len(cols) == 0 {
		return p
	}
	for i, col := range cols {
		if i > 0 {
			p.orderBy.WriteByte(',')
		}
		p.orderBy.WriteString(col.Quote())
	}
	return p
}

//Desc
func (p *sqlBuilder) Desc(cols ...TableField) ISqlBuilder {
	if len(cols) == 0 {
		return p
	}
	for i, col := range cols {
		if i > 0 {
			p.orderBy.WriteByte(',')
		}
		p.orderBy.WriteString(col.Quote() + " DESC")
	}
	return p
}

//UnNull
func (p *sqlBuilder) UnNull(f TableField) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " IS NOT NULL")

	p.andOr = false
	return p
}

//Null
func (p *sqlBuilder) Null(f TableField) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " IS NULL")

	p.andOr = false
	return p
}

//Rlike
func (p *sqlBuilder) Rlike(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " LIKE CONCAT(" + placeholder + ",'%')")
	p.params = append(p.params, v)

	p.andOr = false
	return p
}

//Llike
func (p *sqlBuilder) Llike(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " LIKE CONCAT('%'," + placeholder + ")")
	p.params = append(p.params, v)

	p.andOr = false
	return p
}

//Like
func (p *sqlBuilder) Like(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " LIKE CONCAT('%'," + placeholder + ",'%')")
	p.params = append(p.params, v)

	p.andOr = false
	return p
}

//Bt
func (p *sqlBuilder) Bt(f TableField, v1, v2 interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " BETWEEN " + placeholder + " AND " + placeholder)
	p.params = append(p.params, v1, v2)

	p.andOr = false
	return p
}

//Ue
func (p *sqlBuilder) Ue(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " <> ")
}

//Lte
func (p *sqlBuilder) Lte(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " <= ")
}

//Lt
func (p *sqlBuilder) Lt(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " < ")
}

//Gte
func (p *sqlBuilder) Gte(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " >= ")
}

//Gt
func (p *sqlBuilder) Gt(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " > ")
}

//Eq
func (p *sqlBuilder) Eq(f TableField, v interface{}) ISqlBuilder {
	return p.toWhere(f, v, " = ")
}

//Cols
func (p *sqlBuilder) Cols(args ...TableField) ISqlBuilder {
	p.cols = args
	return p
}

//And 算术方法之间默认为 AND 逻辑
func (p *sqlBuilder) And() ISqlBuilder {
	if !p.andOr {
		p.where.WriteString(operator_and)
		p.andOr = true
	}
	return p
}

//AndWhere
func (p *sqlBuilder) AndWhere(sb ISqlBuilder) ISqlBuilder {
	defer sb.Free()

	if sb.Where() == "" {
		return p
	}

	p.And()
	return p.subCond(sb)
}

//Or
func (p *sqlBuilder) Or() ISqlBuilder {
	if !p.andOr {
		p.where.WriteString(operator_or)
		p.andOr = true
	}
	return p
}

//OrWhere
func (p *sqlBuilder) OrWhere(sb ISqlBuilder) ISqlBuilder {
	defer sb.Free()
	if sb.Where() == "" {
		return p
	}

	p.Or()
	return p.subCond(sb)
}

//Where
func (p *sqlBuilder) Where() string {
	return p.where.String()
}

//Params
func (p *sqlBuilder) Params() []interface{} {
	params := []interface{}{}
	params = append(params, p.params...)
	params = append(params, p.havingParams...)

	return params
}

////
//subCond 子条件
func (p *sqlBuilder) subCond(sb ISqlBuilder) ISqlBuilder {
	p.where.WriteString(" ( ")
	p.where.WriteString(sb.Where())
	p.where.WriteString(" ) ")

	if len(sb.Params()) > 0 {
		p.params = append(p.params, sb.Params()...)
	}

	p.andOr = false
	return p
}

func (p *sqlBuilder) toWhere(f TableField, v interface{}, op string) *sqlBuilder {
	p.prepare()
	p.where.WriteString(f.generate(op))
	p.params = append(p.params, v)

	p.andOr = false
	return p
}
func (p *sqlBuilder) prepare() *sqlBuilder {
	if !p.andOr {
		p.where.WriteString(operator_and)
		p.andOr = true
	}
	return p
}
`

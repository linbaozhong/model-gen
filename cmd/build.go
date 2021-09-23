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
	"libs/types"
	"libs/utils"
	"reflect"
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

type IInsert interface {
	Insert(x interface{}, cols ...string) (int64, error)
	InsertBatch(x interface{}, beans []interface{}, cols ...string) (int64, error)
}

type IUpdate interface {
	Update(x interface{}, id uint64, bean ...interface{}) (int64, error)
	UpdateBatch(x interface{}, cond ISqlBuilder, bean ...interface{}) (int64, error)
}

type IDelete interface {
	Delete(x interface{}, id uint64) (int64, error)
	DeleteBatch(x interface{}, cond ISqlBuilder) (int64, error)
}

type IModel interface {
	Get(x interface{}, id uint64) (bool, error)
	//Find(x interface{}, query string, vals []interface{}, size, index int) ([]interface{}, error)
	ToMap(cols ...TableField) map[string]interface{}
	ToJSON(cols ...TableField) types.Smap
}

type ISqlBuilder interface {
	Table(m interface{}) ISqlBuilder
	GetCondition() (string, []interface{})
	Select() ([]interface{}, error)
	Insert() ([]interface{}, error)
	Update() ([]interface{}, error)
	Delete() ([]interface{}, error)

	Distinct() ISqlBuilder
	Cols(args ...interface{}) ISqlBuilder
	GetCols() []string
	Where(sql string, v interface{}) *sqlBuilder
	//等于
	Eq(f TableField, v interface{}) ISqlBuilder
	//大于
	Gt(f TableField, v interface{}) ISqlBuilder
	//大于等于
	Gte(f TableField, v interface{}) ISqlBuilder
	//小于
	Lt(f TableField, v interface{}) ISqlBuilder
	//小于等于
	Lte(f TableField, v interface{}) ISqlBuilder
	//IN
	In(f TableField, v ...interface{}) ISqlBuilder
	//NOT IN
	UnIn(f TableField, v ...interface{}) ISqlBuilder
	//不等于
	Ue(f TableField, v interface{}) ISqlBuilder
	//BETWEEN
	Bt(f TableField, v1, v2 interface{}) ISqlBuilder
	//LIKE
	Like(f TableField, v interface{}) ISqlBuilder
	//
	Llike(f TableField, v interface{}) ISqlBuilder
	//
	Rlike(f TableField, v interface{}) ISqlBuilder
	//IS NULL
	Null(f TableField) ISqlBuilder
	//IS NOT NULL
	UnNull(f TableField) ISqlBuilder
	//JOIN
	Join(t JoinType, l, r TableField) ISqlBuilder
	GetJoin() string
	//LIMIT
	Limit(size int, start ...int) ISqlBuilder
	GetLimit() (size int, start int)

	And() ISqlBuilder
	AndWhere(sb ISqlBuilder) ISqlBuilder

	Or() ISqlBuilder
	OrWhere(sb ISqlBuilder) ISqlBuilder

	GetWhere() (string, []interface{})
	GetWhereString() string
	GetParams() []interface{}

	GroupBy(cols ...TableField) ISqlBuilder
	GetGroupBy() string

	Having(s string) ISqlBuilder
	GetHaving() string

	OrderBy(cols ...TableField) ISqlBuilder
	Asc(cols ...TableField) ISqlBuilder
	Desc(cols ...TableField) ISqlBuilder
	Rand() ISqlBuilder
	GetOrderBy() string

	//
	Set(f TableField, v interface{}) ISqlBuilder
	Incr(f TableField, v ...interface{}) ISqlBuilder
	Decr(f TableField, v ...interface{}) ISqlBuilder
	SetExpr(f TableField, expr string) ISqlBuilder

	Free()
}
type sqlBuilder struct {
	table       string
	distinct    bool
	cols        []interface{}
	where       strings.Builder
	whereParams []interface{}
	groupBy     strings.Builder
	having      strings.Builder
	//havingParams []interface{}
	orderBy strings.Builder
	limit   string
	limitSize  int
	limitStart int
	join    string

	andOr bool

	updateCols     []string
	updateExprCols []string
	updateExprVals []string
	updateParams   []interface{}
}

var (
	sqlBuilderPool = sync.Pool{New: func() interface{} {
		return &sqlBuilder{
			andOr: true,
		}
	}}
)

//X NewSqlBuilder的简称
func X() *sqlBuilder {
	return sqlBuilderPool.Get().(*sqlBuilder)
}

//NewSqlBuilder 实例化一个 *sqlBuilder
func NewSqlBuilder() *sqlBuilder {
	return sqlBuilderPool.Get().(*sqlBuilder)
}

//Free
func (p *sqlBuilder) Free() {
	p.table = ""
	p.distinct = false
	p.cols = []interface{}{}
	p.where.Reset()
	p.whereParams = []interface{}{}
	p.groupBy.Reset()
	p.having.Reset()
	//p.havingParams = []interface{}{}
	p.orderBy.Reset()
	p.limit = ""
	p.limitStart = 0
	p.limitSize = 0
	p.join = ""

	p.andOr = true

	p.updateCols = []string{}
	p.updateExprCols = []string{}
	p.updateExprVals = []string{}
	p.updateParams = []interface{}{}

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

//Insert
func (p *sqlBuilder) Insert() ([]interface{}, error) {
	defer p.Free()

	if p.table == "" {
		return nil, ErrTableEmpty
	}
	if len(p.updateCols) == 0 && len(p.updateExprCols) == 0 {
		return nil, ErrUpdateEmpty
	}

	var buf strings.Builder
	//INSERT
	buf.WriteString("INSERT INTO " + Quote_Char + p.table + Quote_Char)
	//VALUES
	var cols = make([]string, len(p.updateCols)+len(p.updateExprCols))
	copy(cols, p.updateCols)
	copy(cols[len(p.updateCols):], p.updateExprCols)

	buf.WriteString(" ( " + strings.Join(cols, ", ") + " ) VALUES ( ")
	for i := 0; i < len(p.updateCols); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("?")
	}
	var cLen = len(p.updateExprVals)
	if cLen > 0 {
		if len(p.updateCols) > 0 {
			buf.WriteString(", ")
		}
		for i := 0; i < cLen; i++ {
			buf.WriteString(p.updateExprVals[i])
			if i < cLen-1 {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteString(" ) ")

	r := make([]interface{}, 0, len(p.updateParams)+1)
	r = append(r, buf.String())
	r = append(r, p.updateParams...)

	return r, nil
}

//Update
func (p *sqlBuilder) Update() ([]interface{}, error) {
	defer p.Free()

	if p.table == "" {
		return nil, ErrTableEmpty
	}
	if len(p.updateCols) == 0 && len(p.updateExprCols) == 0 {
		return nil, ErrUpdateEmpty
	}

	var buf strings.Builder
	//UPDATE
	buf.WriteString("UPDATE ")
	//TABLE
	buf.WriteString(Quote_Char + p.table + Quote_Char + " SET ")
	//SET
	if len(p.updateCols) > 0 {
		for i, col := range p.updateCols {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(col + " = ?")
		}
	}
	var vLen = len(p.updateExprVals)
	if vLen > 0 {
		if len(p.updateCols) > 0 {
			buf.WriteString(", ")
		}
		for i := 0; i < vLen; i++ {
			buf.WriteString(p.updateExprCols[i] + " = " + p.updateExprVals[i])
			if i < vLen-1 {
				buf.WriteString(", ")
			}
		}
	}
	//WHERE
	sql, params := p.condition()
	var _params = make([]interface{}, len(p.updateParams)+len(params))
	copy(_params, p.updateParams)

	if sql != "" {
		buf.WriteString(" WHERE " + sql)
		copy(_params[len(p.updateParams):], params)
	}

	r := make([]interface{}, 0, len(_params)+1)
	r = append(r, buf.String())
	r = append(r, _params...)

	return r, nil
}

//Delete
func (p *sqlBuilder) Delete() ([]interface{}, error) {
	defer p.Free()

	if p.table == "" {
		return nil, ErrTableEmpty
	}

	var buf strings.Builder
	//UPDATE
	buf.WriteString("DELETE FROM ")
	//TABLE
	buf.WriteString(Quote_Char + p.table + Quote_Char + " ")
	//WHERE
	sql, params := p.condition()
	var _params = make([]interface{}, len(p.updateParams)+len(params))
	copy(_params, p.updateParams)

	if sql != "" {
		buf.WriteString(" WHERE " + sql)
		copy(_params[len(p.updateParams):], params)
	}

	r := make([]interface{}, 0, len(_params)+1)
	r = append(r, buf.String())
	r = append(r, _params...)

	return r, nil
}

//Select
func (p *sqlBuilder) Select() ([]interface{}, error) {
	defer p.Free()

	if p.table == "" {
		return nil, ErrTableEmpty
	}
	var buf strings.Builder
	//SELECT
	buf.WriteString("SELECT ")
	if len(p.cols) == 0 {
		buf.WriteString("*")
	} else {
		if p.distinct {
			buf.WriteString("DISTINCT ")
		}
		buf.WriteString(p.getColString())
	}
	//FROM TABLE
	buf.WriteString(" FROM " + Quote_Char + p.table + Quote_Char)
	//JOIN
	if p.join != "" {
		buf.WriteString(p.join)
	}
	//if len(p.join) == 3 {
	//	buf.WriteString(p.join[0] + p.join[1] + " ON " + p.join[2])
	//}
	//WHERE
	sql, params := p.condition()
	if sql != "" {
		buf.WriteString(" WHERE " + sql)
	}
	//LIMIT
	if p.limit != "" {
		buf.WriteString(p.limit)
	}

	r := make([]interface{}, 0, len(params)+1)
	r = append(r, buf.String())
	r = append(r, params...)

	return r, nil
}

//GetCondition
func (p *sqlBuilder) GetCondition() (string, []interface{}) {
	defer p.Free()
	return p.condition()
}

//
func (p *sqlBuilder) Distinct() ISqlBuilder {
	p.distinct = true
	return p
}

//JOIN
func (p *sqlBuilder) Join(t JoinType, l, r TableField) ISqlBuilder {
	p.join += string(t) + Quote_Char + r.Table + Quote_Char + " ON " + l.Quote() + " = " + r.Quote() + " "
	//p.join = append(p.join, string(t), Quote_Char+r.Table+Quote_Char, r.Quote()+" = "+l.Quote())
	return p
}

func (p *sqlBuilder) GetJoin() string {
	return p.join
}

//LIMIT
func (p *sqlBuilder) Limit(size int, start ...int) ISqlBuilder {
	p.limitSize = size
	if len(start) > 0 {
		p.limitStart = start[0]
	}
	p.limit = " LIMIT " + strconv.Itoa(p.limitSize) + " OFFSET " + strconv.Itoa(p.limitStart)

	return p
}

//GetLimit
func (p *sqlBuilder) GetLimit() (size int, start int) {
	size = p.limitSize
	start = p.limitStart
	return
}

//GROUPBY
func (p *sqlBuilder) GroupBy(cols ...TableField) ISqlBuilder {
	if len(cols) == 0 {
		return p
	}
	for _, col := range cols {
		if p.groupBy.Len() > 0 {
			p.groupBy.WriteByte(',')
		}
		p.groupBy.WriteString(col.Quote())
	}
	return p
}

func (p *sqlBuilder) GetGroupBy() string {
	return p.groupBy.String()
}

//HAVING
func (p *sqlBuilder) Having(s string) ISqlBuilder {
	if s == "" {
		return p
	}
	p.having.WriteString(s)
	return p
}

func (p *sqlBuilder) GetHaving() string {
	return p.having.String()
}

func (p *sqlBuilder) GetOrderBy() string {
	return p.orderBy.String()
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
	for _, col := range cols {
		if p.orderBy.Len() > 0 {
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
	for _, col := range cols {
		if p.orderBy.Len() > 0 {
			p.orderBy.WriteByte(',')
		}
		p.orderBy.WriteString(col.Quote() + " DESC")
	}
	return p
}

//orderby rand()
func (p *sqlBuilder) Rand() ISqlBuilder {
	p.orderBy.WriteString(" rand()")
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
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Llike
func (p *sqlBuilder) Llike(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " LIKE CONCAT('%'," + placeholder + ")")
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Like
func (p *sqlBuilder) Like(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " LIKE CONCAT('%'," + placeholder + ",'%')")
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Bt
func (p *sqlBuilder) Bt(f TableField, v1, v2 interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString(f.Quote() + " BETWEEN " + placeholder + " AND " + placeholder)
	p.whereParams = append(p.whereParams, v1, v2)

	p.andOr = false
	return p
}


//In
func (p *sqlBuilder) In(f TableField, v ...interface{}) ISqlBuilder {
	if len(v) == 0 {
		return p
	}

	vv := reflect.ValueOf(v[0])
	if vv.Kind() == reflect.Slice {
		l := vv.Len()
		if l == 0 {
			return p
		}
		p.prepare()
		p.where.WriteString(f.Quote() + " IN (" + strings.Repeat(placeholder+",", l)[:2*l-1] + ") ")
		for i := 0; i < l; i++ {
			p.whereParams = append(p.whereParams, vv.Index(i).Interface())
		}
	} else {
		p.prepare()
		p.where.WriteString(f.Quote() + " IN (" + strings.Repeat(placeholder+",", len(v))[:2*len(v)-1] + ") ")
		p.whereParams = append(p.whereParams, v...)
	}

	p.andOr = false
	return p
}

//UnIn
func (p *sqlBuilder) UnIn(f TableField, v ...interface{}) ISqlBuilder {
	if len(v) == 0 {
		return p
	}
	
	vv := reflect.ValueOf(v[0])
	if vv.Kind() == reflect.Slice {
		l := vv.Len()
		if l == 0 {
			return p
		}
		p.prepare()
		p.where.WriteString(f.Quote() + " NOT IN (" + strings.Repeat(placeholder+",", l)[:2*l-1] + ") ")
		for i := 0; i < l; i++ {
			p.whereParams = append(p.whereParams, vv.Index(i).Interface())
		}
	} else {
		p.prepare()
		p.where.WriteString(f.Quote() + " NOT IN (" + strings.Repeat(placeholder+",", len(v))[:2*len(v)-1] + ") ")
		p.whereParams = append(p.whereParams, v...)
	}
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

//Where
func (p *sqlBuilder) Where(sql string, v interface{}) *sqlBuilder {
	p.prepare()
	p.where.WriteString(sql)
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Cols
func (p *sqlBuilder) Cols(args ...interface{}) ISqlBuilder {
	p.cols = args
	return p
}

func (p *sqlBuilder) GetCols() []string {
	if len(p.cols) == 0 {
		return []string{}
	}
	s := make([]string, 0, len(p.cols))
	for _, col := range p.cols {
		if _f, ok := col.(TableField); ok {
			s = append(s, _f.Quote())
		} else if _f, ok := col.(string); ok {
			s = append(s, _f)
		}
	}
	return s
}

func (p *sqlBuilder) getColString() string {
	if len(p.cols) == 0 {
		return ""
	}
	buf := strings.Builder{}
	for i, col := range p.cols {
		if i > 0 {
			buf.WriteByte(',')
		}
		if _f, ok := col.(TableField); ok {
			buf.WriteString(_f.Quote())
		} else if _f, ok := col.(string); ok {
			buf.WriteString(_f)
		}
	}
	return buf.String()
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

	if sb.GetWhereString() == "" {
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
	if sb.GetWhereString() == "" {
		return p
	}

	p.Or()
	return p.subCond(sb)
}

//GetWhere
func (p *sqlBuilder) GetWhere() (string, []interface{}) {
	return p.where.String(), p.whereParams
}

//GetWhereString
func (p *sqlBuilder) GetWhereString() string {
	return p.where.String()
}

//GetParams
func (p *sqlBuilder) GetParams() []interface{} {
	//params := []interface{}{}
	//params = append(params, p.whereParams...)
	//params = append(params, p.havingParams...)
	//
	//return params
	return p.whereParams
}

//condition
func (p *sqlBuilder) condition() (string, []interface{}) {
	var buf strings.Builder
	//WHERE
	if p.where.Len() > 0 {
		buf.WriteString(p.GetWhereString())
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

	return buf.String(), p.GetParams()
}

////
//Set
func (p *sqlBuilder) Set(f TableField, v interface{}) ISqlBuilder {
	p.updateCols = append(p.updateCols, f.QuoteName())
	p.updateParams = append(p.updateParams, v)
	return p
}

//Incr
func (p *sqlBuilder) Incr(f TableField, v ...interface{}) ISqlBuilder {
	var _v = "1"
	if len(v) > 0 {
		_v = utils.Interface2String(v[0])
	}
	//p.updateCols = append(p.updateCols, f.QuoteName()+" = "+f.QuoteName()+"+"+_v)
	return p.SetExpr(f, f.QuoteName()+"+"+_v)
}

//Decr
func (p *sqlBuilder) Decr(f TableField, v ...interface{}) ISqlBuilder {
	var _v = "1"
	if len(v) > 0 {
		_v = utils.Interface2String(v[0])
	}
	//p.updateCols = append(p.updateCols, f.QuoteName()+" = "+f.QuoteName()+"-"+_v)
	return p.SetExpr(f, f.QuoteName()+"-"+_v)
}

//SetExpr
func (p *sqlBuilder) SetExpr(f TableField, expr string) ISqlBuilder {
	p.updateExprCols = append(p.updateExprCols, f.QuoteName())
	p.updateExprVals = append(p.updateExprVals, expr)
	return p
}

//String
func (p *sqlBuilder) String() string {
	var buf strings.Builder
	sql, pm := p.condition()
	buf.WriteString(sql)
	buf.WriteString("@params:")
	for _, i := range pm {
		buf.WriteString(utils.Interface2String(i) + "|")
	}
	return buf.String()
}

////
//subCond 子条件
func (p *sqlBuilder) subCond(sb ISqlBuilder) ISqlBuilder {
	s := sb.GetWhereString()
	if s == "" {
		return p
	}

	p.where.WriteString(" ( ")
	p.where.WriteString(s)
	p.where.WriteString(" ) ")

	if len(sb.GetParams()) > 0 {
		p.whereParams = append(p.whereParams, sb.GetParams()...)
	}

	p.andOr = false
	return p
}

func (p *sqlBuilder) toWhere(f TableField, v interface{}, op string) *sqlBuilder {
	p.prepare()
	p.where.WriteString(f.generate(op))
	p.whereParams = append(p.whereParams, v)

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

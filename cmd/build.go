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

	f, e := os.OpenFile(buildFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
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

	var buf bytes.Buffer
	_ = template.Must(template.New("buildTpl").Parse(buildTpl)).Execute(&buf, nil)
	formatted, _ := format.Source(buf.Bytes())
	_, e = f.Write(formatted)
	if e != nil {
		showError(e.Error())
	}
	return e
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
type command string

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	operator_and = " AND "
	operator_or  = " OR "
	placeholder  = "?"

	command_insert command = "INSERT INTO "
	command_select command = "SELECT "
	command_update command = "UPDATE "
	command_delete command = "DELETE FROM "
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

type IDataAccess interface {
	Get(x interface{}, id uint64) (bool, error)
	//Find(x interface{}, query string, vals []interface{}, size, index int) ([]interface{}, error)
	ToMap(cols ...TableField) map[string]interface{}
	ToJSON(cols ...TableField) types.Smap
}

type ISqlBuilder interface {
	ToSQL() (string, []interface{}, error)
	Table(m interface{}) ISqlBuilder
	GetCondition() (string, []interface{})
	Select() ([]interface{}, error)
	Insert() ([]interface{}, error)
	Update() ([]interface{}, error)
	Delete() ([]interface{}, error)

	Distinct() ISqlBuilder
	Cols(args ...interface{}) ISqlBuilder
	Omit(args ...interface{}) ISqlBuilder
	GetColsX(args []string) []string
	GetCols() []string
	GetOmit() []string
	Where(sql string, v ...interface{}) *sqlBuilder
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
	GetJoin() [][3]string
	//LIMIT
	Limit(size int, start ...int) ISqlBuilder
	GetLimit() (size int, start int)

	And(sb ISqlBuilder) ISqlBuilder

	Or(sb ISqlBuilder) ISqlBuilder

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
	GetUpdate() ([]string, []interface{})
	Incr(f TableField, v ...interface{}) ISqlBuilder
	GetIncr() []Expr
	Decr(f TableField, v ...interface{}) ISqlBuilder
	GetDecr() []Expr
	SetExpr(f TableField, expr string) ISqlBuilder
	Replace(f TableField, o, n string) ISqlBuilder
	GetExpr() []Expr
	//Sum(cols ...TableField) ISqlBuilder
	//GetSum() []string

	Free()
}

//Insert
func I(t interface{}) ISqlBuilder {
	return get(t, command_insert)
}

//Update
func U(t interface{}) ISqlBuilder {
	return get(t, command_update)
}

//Delete
func D(t interface{}) ISqlBuilder {
	return get(t, command_delete)
}

//Select
func S(t interface{}) ISqlBuilder {
	return get(t, command_select)
}

func get(t interface{}, cmd command) ISqlBuilder {
	p := sqlBuilderPool.Get().(*sqlBuilder)
	p.setTable(t)
	p.setCmd(cmd)
	return p
}

//Table
func Table(t interface{}) ISqlBuilder {
	return NewSqlBuilder().setTable(t)
}

////Columns
//func Columns(cols ...interface{}) []string {
//	_cols := make([]string, 0, len(cols))
//	for i := 0; i < len(cols); i++ {
//		if c, ok := cols[i].(TableField); ok {
//			_cols = append(_cols, c.Name)
//		} else {
//			_cols = append(_cols, utils.Interface2String(cols[i]))
//		}
//	}
//	return _cols
//}

//JOIN
func Join(t JoinType, l, r TableField) ISqlBuilder {
	return NewSqlBuilder().Join(t, l, r)
}

//UnNull Is Not Null
func UnNull(f TableField) ISqlBuilder {
	return NewSqlBuilder().UnNull(f)
}

//Null Is Null
func Null(f TableField) ISqlBuilder {
	return NewSqlBuilder().Null(f)
}

//Rlike
func Rlike(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().Rlike(f, v)
}

//Llike
func Llike(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().Llike(f, v)
}

//Like
func Like(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().Like(f, v)
}

//Bt Between
func Bt(f TableField, v1, v2 interface{}) ISqlBuilder {
	return NewSqlBuilder().Bt(f, v1, v2)
}

//In
func In(f TableField, v ...interface{}) ISqlBuilder {
	return NewSqlBuilder().In(f, v...)
}

//UnIn Not In
func UnIn(f TableField, v ...interface{}) ISqlBuilder {
	return NewSqlBuilder().UnIn(f, v...)
}

//Ue !=
func Ue(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " <> ")
}

//Lte <=
func Lte(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " <= ")
}

//Lt <
func Lt(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " < ")
}

//Gte >=
func Gte(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " >= ")
}

//Gt >
func Gt(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " > ")
}

//Eq =
func Eq(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().toWhere(f, v, " = ")
}

//Where
func Where(sql string, v ...interface{}) *sqlBuilder {
	return NewSqlBuilder().Where(sql, v...)
}

//Cols
func Cols(args ...interface{}) ISqlBuilder {
	return NewSqlBuilder().Cols(args...)
}

//Omit
func Omit(args ...interface{}) ISqlBuilder {
	return NewSqlBuilder().Omit(args...)
}

//OrderBy
func OrderBy(cols ...TableField) ISqlBuilder {
	return NewSqlBuilder().Asc(cols...)
}

//Asc
func Asc(cols ...TableField) ISqlBuilder {
	return NewSqlBuilder().Asc(cols...)
}

//Desc
func Desc(cols ...TableField) ISqlBuilder {
	return NewSqlBuilder().Desc(cols...)
}

//Set
func Set(f TableField, v interface{}) ISqlBuilder {
	return NewSqlBuilder().Set(f, v)
}

//Incr
func Incr(f TableField, v ...interface{}) ISqlBuilder {
	return NewSqlBuilder().Incr(f, v...)
}

//Decr
func Decr(f TableField, v ...interface{}) ISqlBuilder {
	return NewSqlBuilder().Decr(f, v...)
}

//SetExpr
func SetExpr(f TableField, expr string) ISqlBuilder {
	return NewSqlBuilder().SetExpr(f, expr)
}

//Replace
func Replace(f TableField, o, n string) ISqlBuilder {
	return NewSqlBuilder().Replace(f, o, n)
}

// Expr represents an SQL express
type Expr struct {
	ColName string
	Arg     interface{}
}

type sqlBuilder struct {
	cmd         command
	table       string
	distinct    bool
	cols        []interface{}
	omit        []interface{}
	where       strings.Builder
	whereParams []interface{}
	groupBy     strings.Builder
	having      strings.Builder
	//havingParams []interface{}
	orderBy    strings.Builder
	limit      string
	limitSize  int
	limitStart int
	join       [][3]string

	andOr bool

	updateCols   []string
	updateParams []interface{}
	incrCols     []Expr
	decrCols     []Expr
	exprCols     []Expr
	//sumCols      []string
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
	p.cmd = ""
	p.table = ""
	p.distinct = false
	p.cols = p.cols[0:0]	//[]interface{}{}
	p.omit = p.omit[0:0]
	p.where.Reset()
	p.whereParams = p.whereParams[0:0]	// []interface{}{}
	p.groupBy.Reset()
	p.having.Reset()
	//p.havingParams = p.havingParams[0:0]	// []interface{}{}
	p.orderBy.Reset()
	p.limit = ""
	p.limitStart = 0
	p.limitSize = 0
	p.join = [][3]string{}

	p.andOr = true

	p.updateCols = p.updateCols[0:0]	// []string{}
	p.updateParams = p.updateParams[0:0]	// []interface{}{}
	p.incrCols = p.incrCols[0:0]	// []Expr{}
	p.decrCols = p.decrCols[0:0]	// []Expr{}
	p.exprCols = p.exprCols[0:0]	// []Expr{}
	//p.sumCols = p.sumCols[0:0]	// []string{}

	sqlBuilderPool.Put(p)
}

//Table
func (p *sqlBuilder) Table(t interface{}) ISqlBuilder {
	return p.setTable(t)
}

func (p *sqlBuilder) ToSQL() (string, []interface{}, error) {
	defer p.Free()
	switch p.cmd {
	case command_insert:
		return p.getInsert()
	case command_update:
		return p.getUpdate()
	case command_delete:
		return p.getDelete()
	default:
		return p.getSelect()
	}
}

//Insert
func (p *sqlBuilder) Insert() ([]interface{}, error) {
	defer p.Free()

	sql, args, e := p.getInsert()
	if e != nil {
		return nil, e
	}
	r := make([]interface{}, 0, len(p.updateParams)+1)
	r = append(r, sql)
	r = append(r, args...)

	return r, nil
}

//Update
func (p *sqlBuilder) Update() ([]interface{}, error) {
	defer p.Free()
	sql, args, e := p.getUpdate()
	if e != nil {
		return nil, e
	}
	r := make([]interface{}, 0, len(p.updateParams)+1)
	r = append(r, sql)
	r = append(r, args...)

	return r, nil
}

//Delete
func (p *sqlBuilder) Delete() ([]interface{}, error) {
	defer p.Free()

	sql, args, e := p.getDelete()
	if e != nil {
		return nil, e
	}
	r := make([]interface{}, 0, len(p.updateParams)+1)
	r = append(r, sql)
	r = append(r, args...)

	return r, nil
}

//Select
func (p *sqlBuilder) Select() ([]interface{}, error) {
	defer p.Free()

	sql, args, e := p.getSelect()
	if e != nil {
		return nil, e
	}
	r := make([]interface{}, 0, len(p.updateParams)+1)
	r = append(r, sql)
	r = append(r, args...)

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
	p.join = append(p.join, [3]string{
		string(t),
		Quote_Char + r.Table + Quote_Char,
		l.Quote() + " = " + r.Quote(),
	})
	return p
}

func (p *sqlBuilder) GetJoin() [][3]string {
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
	p.where.WriteString(" NOT ISNULL(" + placeholder + ")")
	p.whereParams = append(p.whereParams,f.Quote())

	p.andOr = false
	return p
}

//Null
func (p *sqlBuilder) Null(f TableField) ISqlBuilder {
	p.prepare()
	p.where.WriteString(" ISNULL(" + placeholder + ")")
	p.whereParams = append(p.whereParams,f.Quote())

	p.andOr = false
	return p
}

//Rlike
func (p *sqlBuilder) Rlike(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString("(" + f.Quote() + " LIKE CONCAT(" + placeholder + ",'%'))")
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Llike
func (p *sqlBuilder) Llike(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString("(" + f.Quote() + " LIKE CONCAT('%'," + placeholder + "))")
	p.whereParams = append(p.whereParams, v)

	p.andOr = false
	return p
}

//Like
func (p *sqlBuilder) Like(f TableField, v interface{}) ISqlBuilder {
	p.prepare()
	p.where.WriteString("(" + f.Quote() + " LIKE CONCAT('%'," + placeholder + ",'%'))")
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
func (p *sqlBuilder) Where(sql string, v ...interface{}) *sqlBuilder {
	p.prepare()
	p.where.WriteString(sql)
	p.whereParams = append(p.whereParams, v...)

	p.andOr = false
	return p
}

//Cols
func (p *sqlBuilder) Cols(args ...interface{}) ISqlBuilder {
	p.cols = args
	return p
}

//Omit
func (p *sqlBuilder) Omit(args ...interface{}) ISqlBuilder {
	p.omit = args
	return p
}

//GetColsX
func (p *sqlBuilder) GetColsX(args []string) []string {
	var s = p.GetCols()
	if len(s) == 0 {
		s = args
	}
	if len(s) > 0 {
		o := p.GetOmit()
		if len(o) == 0 {
			return s
		}
		m := make(map[string]struct{}, len(o))
		for i := 0; i < len(o); i++ {
			m[o[i]] = struct{}{}
		}
		_s := make([]string, 0, len(s))
		for i := 0; i < len(s); i++ {
			if _, ok := m[s[i]]; !ok {
				_s = append(_s, s[i])
			}
		}
		return _s
	}
	return s
}

//GetCols
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
		} else if _fs, ok := col.([]string); ok {
			s = append(s, _fs...)
		}
	}
	return s
}

//GetOmit
func (p *sqlBuilder) GetOmit() []string {
	if len(p.omit) == 0 {
		return []string{}
	}
	s := make([]string, 0, len(p.omit))
	for _, col := range p.omit {
		if _f, ok := col.(TableField); ok {
			s = append(s, _f.Quote())
		} else if _f, ok := col.(string); ok {
			s = append(s, _f)
		} else if _fs, ok := col.([]string); ok {
			s = append(s, _fs...)
		}
	}
	return s
}

//And 算术方法之间默认为 AND 逻辑
func (p *sqlBuilder) And(sb ISqlBuilder) ISqlBuilder {
	defer sb.Free()
	if sb.GetWhereString() == "" {
		return p
	}

	if !p.andOr {
		p.where.WriteString(operator_and)
		p.andOr = true
	}
	p.where.WriteString("(")
	p.subCond(sb)
	p.where.WriteString(")")
	return p
}

//Or
func (p *sqlBuilder) Or(sb ISqlBuilder) ISqlBuilder {
	defer sb.Free()
	if sb.GetWhereString() == "" {
		return p
	}

	if !p.andOr {
		p.where.WriteString(operator_or)
		p.andOr = true
	}
	p.where.WriteString("(")
	p.subCond(sb)
	p.where.WriteString(")")
	return p
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
	return p.whereParams
}

////
//Set
func (p *sqlBuilder) Set(f TableField, v interface{}) ISqlBuilder {
	p.updateCols = append(p.updateCols, f.PureQuote())
	p.updateParams = append(p.updateParams, v)
	return p
}

//Incr
func (p *sqlBuilder) Incr(f TableField, v ...interface{}) ISqlBuilder {
	var _v = "1"
	if len(v) > 0 {
		_v = utils.Interface2String(v[0])
	}
	p.incrCols = append(p.incrCols, Expr{
		ColName: f.PureQuote(),
		Arg:     _v,
	})
	return p
}

//Decr
func (p *sqlBuilder) Decr(f TableField, v ...interface{}) ISqlBuilder {
	var _v = "1"
	if len(v) > 0 {
		_v = utils.Interface2String(v[0])
	}
	p.decrCols = append(p.decrCols, Expr{
		ColName: f.Quote(),
		Arg:     _v,
	})
	return p
}

//SetExpr
func (p *sqlBuilder) SetExpr(f TableField, expr string) ISqlBuilder {
	p.exprCols = append(p.exprCols, Expr{
		ColName: f.PureQuote(),
		Arg:     expr,
	})
	return p
}

//// Sum
//func (p *sqlBuilder) Sum(cols ...TableField) ISqlBuilder {
//	for _, f := range cols {
//		p.sumCols = append(p.sumCols, f.Quote())
//	}
//	return p
//}

//Replace
func (p *sqlBuilder) Replace(f TableField, o, n string) ISqlBuilder {
	p.exprCols = append(p.exprCols, Expr{
		ColName: f.PureQuote(),
		Arg:     "REPLACE(" + f.PureQuote() + ",'" + o + "','" + n + "')",
	})
	return p
}

func (p *sqlBuilder) GetIncr() []Expr {
	return p.incrCols
}

func (p *sqlBuilder) GetDecr() []Expr {
	return p.decrCols
}

func (p *sqlBuilder) GetExpr() []Expr {
	return p.exprCols
}

//// GetSum
//func (p *sqlBuilder) GetSum() []string {
//	return p.sumCols
//}

func (p *sqlBuilder) GetUpdate() ([]string, []interface{}) {
	return p.updateCols, p.updateParams
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

	p.where.WriteString(s)

	if len(sb.GetParams()) > 0 {
		p.whereParams = append(p.whereParams, sb.GetParams()...)
	}

	p.andOr = false
	return p
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

func (p *sqlBuilder) toWhere(f TableField, v interface{}, op string) *sqlBuilder {
	p.prepare()
	p.where.WriteString("(" + f.generate(op) + ")")
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

//setTable
func (p *sqlBuilder) setTable(t interface{}) ISqlBuilder {
	if name, ok := t.(string); ok {
		p.table = name
	} else if iface, ok := t.(ITableName); ok {
		p.table = iface.TableName()
	}
	return p
}
func (p *sqlBuilder) setCmd(c command) ISqlBuilder {
	p.cmd = c
	return p
}

func (p *sqlBuilder) getInsert() (string, []interface{}, error) {
	if p.table == "" {
		return "", nil, ErrTableEmpty
	}
	if len(p.updateCols) == 0 {
		return "", nil, ErrUpdateEmpty
	}
	var buf strings.Builder
	//INSERT
	buf.WriteString(string(command_insert) + Quote_Char + p.table + Quote_Char)
	//VALUES
	var cols = make([]string, len(p.updateCols))
	copy(cols, p.updateCols)

	buf.WriteString(" ( " + strings.Join(cols, ", ") + " ) VALUES ( ")
	for i := 0; i < len(cols); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("?")
	}

	buf.WriteString(" ) ")

	return buf.String(), p.updateParams, nil
}

func (p *sqlBuilder) getUpdate() (string, []interface{}, error) {
	if p.table == "" {
		return "", nil, ErrTableEmpty
	}

	if len(p.updateCols) == 0 &&
		len(p.incrCols) == 0 &&
		len(p.decrCols) == 0 &&
		len(p.exprCols) == 0 {
		return "", nil, ErrUpdateEmpty
	}

	var buf strings.Builder
	//UPDATE
	buf.WriteString(string(command_update) + Quote_Char + p.table + Quote_Char + " SET ")
	//SET
	cols := make([]string, 0, 5)
	for _, col := range p.updateCols {
		cols = append(cols, col+" = ?")
	}
	for _, col := range p.incrCols {
		cols = append(cols, col.ColName+" = "+col.ColName+" + ?")
		p.updateParams = append(p.updateParams, col.Arg)
	}
	for _, col := range p.decrCols {
		cols = append(cols, col.ColName+" = "+col.ColName+" - ?")
		p.updateParams = append(p.updateParams, col.Arg)
	}
	for _, col := range p.exprCols {
		cols = append(cols, col.ColName+" = "+utils.Interface2String(col.Arg))
	}

	buf.WriteString(strings.Join(cols, ", "))
	//WHERE
	sql, params := p.condition()
	var _params = make([]interface{}, len(p.updateParams)+len(params))
	copy(_params, p.updateParams)

	if sql != "" {
		buf.WriteString(" WHERE " + sql)
		copy(_params[len(p.updateParams):], params)
	}

	return buf.String(), _params, nil
}

func (p *sqlBuilder) getDelete() (string, []interface{}, error) {
	if p.table == "" {
		return "", nil, ErrTableEmpty
	}
	var buf strings.Builder
	//DELETE
	buf.WriteString(string(command_delete) + Quote_Char + p.table + Quote_Char + " ")
	//WHERE
	sql, params := p.condition()

	if sql != "" {
		buf.WriteString(" WHERE " + sql)
	}

	return buf.String(), params, nil
}

func (p *sqlBuilder) getSelect() (string, []interface{}, error) {
	if p.table == "" {
		return "", nil, ErrTableEmpty
	}
	var buf strings.Builder
	//SELECT
	buf.WriteString(string(command_select))
	if len(p.cols) == 0 {
		buf.WriteString("*")
	} else {
		if p.distinct {
			buf.WriteString("DISTINCT ")
		}
		buf.WriteString(strings.Join(p.GetCols(), ","))
	}
	//FROM TABLE
	buf.WriteString(" FROM " + Quote_Char + p.table + Quote_Char)
	for _, j := range p.join {
		buf.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}
	//WHERE
	sql, params := p.condition()
	if sql != "" {
		buf.WriteString(" WHERE " + sql)
	}
	//LIMIT
	if p.limit != "" {
		buf.WriteString(p.limit)
	}
	return buf.String(), params, nil
}

`

package cmd

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

func writeBaseFile(filename string) error {
	baseFilename, _ := filepath.Abs(filename)

	file, err := os.Create(baseFilename)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	_ = template.Must(template.New("baseTpl").Parse(baseTpl)).Execute(&buf, nil)
	formatted, _ := format.Source(buf.Bytes())
	_, err = file.Write(formatted)
	if err != nil {
		showError(err.Error())
	}
	return err
}

var baseTpl = `
		package table

		import "strings"

		const (
			Quote_Char = "` + "`" + `"
		)

		//Select 生成select字段字符串
		func Select(fields ...interface{}) string {
			l := len(fields)
			if l == 0 {
				return ""
			}
			buf := strings.Builder{}
			for i, f := range fields {
				if i > 0 {
					buf.WriteByte(',')
				}
				if _f, ok := f.(TableField); ok {
					buf.WriteString(_f.Quote())
				} else if _f, ok := f.(string); ok {
					buf.WriteString(_f)
				} else {
					
				}
			}
			return buf.String()
		}

		type TableField struct {
			Name string
			Json string
			Table string
		}
		//Eq 等于
		func (f *TableField) Eq() string {
			return f.generate("=")
		}
		//Gt 大于
		func (f *TableField) Gt() string {
			return f.generate(">")
		}
		//Gte 大于等于
		func (f *TableField) Gte() string {
			return f.generate(">=")
		}
		//Lt 小于
		func (f *TableField) Lt() string {
			return f.generate("<")
		}
		//Lte 小于等于
		func (f *TableField) Lte() string {
			return f.generate("<=")
		}
		//Ue 不等于
		func (f *TableField)Ue() string {
			return f.generate("<>")
		}
		//Bt BETWEEN
		func (f *TableField)Bt() string {
			return f.Quote() + " BETWEEN ? AND ?"
		}
		//Like LIKE
		func (f *TableField) Like() string {
			return f.Quote() + " LIKE CONCAT('%',?,'%')"
		}

		//Like 左like
		func (f *TableField) Llike() string {
			return f.Quote() + " LIKE CONCAT('%',?)"
		}
		
		//Like 右like
		func (f *TableField) Rlike() string {
			return f.Quote() + " LIKE CONCAT(?,'%')"
		}

		//JOIN
		func (f *TableField) Join(joinOp string, col TableField) (string, string, string) {
			return strings.ToUpper(joinOp), col.Table, col.Quote() + "=" + f.Quote()
		}

		//AsName
		func (f *TableField) AsName(s string) string {
			return f.Quote() + " AS " + s
		}

		func (f *TableField) Quote() string {
			if f.Table == "" {
				return Quote_Char + f.Name + Quote_Char
			}
			return f.Table + "." + Quote_Char + f.Name + Quote_Char
		}
		func (f *TableField) generate(op string) string {
			return f.Quote() + op + "?"
		}

		///////////////
		type conditions struct {
			Sql  strings.Builder
			Args []interface{}
		}
		
		func (p *conditions) Where() {
			fmt.Println(p.Sql.String(), p.Args)
		}
		
		func Insert(a ...interface{}) {
			var (
				l  = len(a)
				fs = make([]string, 0, l)
				vs = make([]interface{}, 0, l)
				f  TableField
				ok bool
			)
			for i := 0; i < l; i++ {
				if f, ok = a[i].(TableField); ok {
					fs = append(fs, f.Quote())
					i++
					vs = append(vs, a[i])
				} else {
					continue
				}
			}
			var sql strings.Builder
			sql.WriteString("insert into " + f.Table + " (" + strings.Join(fs, ",") + ") values (")
			sql.WriteString(strings.Repeat("?,", len(vs))[0 : len(vs)*2-1])
			sql.WriteString(")")
			fmt.Println(sql.String(), vs)
		}
		
		func Update(a ...interface{}) *conditions {
			var (
				l  = len(a)
				fs = make([]string, 0, l)
				vs = make([]interface{}, 0, l)
				f  TableField
				ok bool
			)
			for i := 0; i < l; i++ {
				if f, ok = a[i].(TableField); ok {
					fs = append(fs, f.Quote()+"=?")
					i++
					vs = append(vs, a[i])
				} else {
					continue
				}
			}
			var sql strings.Builder
			sql.WriteString("update " + f.Table + " set " + strings.Join(fs, ","))
			return &conditions{Sql: sql, Args: vs}
		}

		`

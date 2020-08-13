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
	_ = template.Must(template.New("temp").Parse(base)).Execute(&buf, nil)
	formatted, _ := format.Source(buf.Bytes())
	_, err = file.Write(formatted)
	if err != nil {
		showError(err.Error())
	}
	return err
}

var base = `
		package table

		const (
			Quote_Char = "` + "`" + `"
		)

		type TableField struct {
			Name string
			Json string
		}
		//Eq 等于
		func (f TableField) Eq() string {
			return f.generate("=")
		}
		//Gt 大于
		func (f TableField) Gt() string {
			return f.generate(">")
		}
		//Gte 大于等于
		func (f TableField) Gte() string {
			return f.generate(">=")
		}
		//Lt 小于
		func (f TableField) Lt() string {
			return f.generate("<")
		}
		//Lte 小于等于
		func (f TableField) Lte() string {
			return f.generate("<=")
		}
		//Ue 不等于
		func (f TableField)Ue() string {
			return f.generate("<>")
		}
		//Bt BETWEEN
		func (f TableField)Bt() string {
			return f.QuoteName() + " BETWEEN ? AND ?"
		}
		//Like LIKE
		func (f TableField) Like() string {
			return f.QuoteName() + " LIKE CONCAT('%',?,'%')"
		}

		//Like 左like
		func (f TableField) Llike() string {
			return f.QuoteName() + " LIKE CONCAT('%',?)"
		}
		
		//Like 右like
		func (f TableField) Rlike() string {
			return f.QuoteName() + " LIKE CONCAT(?,'%')"
		}
		
		func (f TableField) QuoteName() string {
			return Quote_Char + f.Name + Quote_Char
		}
		func (f TableField) generate(op string) string {
			return f.QuoteName() + op + "?"
		}

		`

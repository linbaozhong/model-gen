package cmd

var (
	tableName = "TableName"
	tableTpl  = `
		package table
		import (
			"strings"
		)
		type _{{.StructName}} struct {
			TableName string
		{{range $key, $value := .Columns}} {{ $key }} TableField 
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"
		{{range $key, $value := .Columns}} 
		{{ $.StructName}}.{{$key}} = TableField{
			Name: "{{index $value 0}}",
			Json: "{{index $value 1}}",
			Table: {{$.StructName}}.TableName,
		} 
		{{end}}
		}

		func (*_{{.StructName}}) Select(fields ...TableField) string {
			if len(fields) == 0 {
				return ""
			}
			buf := strings.Builder{}
			for i, f := range fields {
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(f.Quote())
			}
			return buf.String()
		}
		`
)

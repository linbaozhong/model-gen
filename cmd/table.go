package cmd

var (
	tableName = "TableName"
	tableTpl  = `
		package table

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
			Table: {{.StructName}}.TableName,
		} 
		{{end}}
		}
		`
)

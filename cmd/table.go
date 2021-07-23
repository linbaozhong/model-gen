package cmd

var (
	//
	tableName = "TableName"
	//
	tableTpl = `
		package table

		type _{{.StructName}} struct {
			TableName string
			PrimaryKey TableField
		{{range $key, $value := .Columns}} {{ $key }} TableField 
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"

			{{.StructName}}.PrimaryKey = TableField{
				Name: "{{index .PrimaryKey 0}}",
				Json: "{{index .PrimaryKey 1}}",
				Table: {{$.StructName}}.TableName,
			}
		{{range $key, $value := .Columns}} 
		{{ $.StructName}}.{{$key}} = TableField{
			Name: "{{index $value 0}}",
			Json: "{{index $value 1}}",
			Table: {{$.StructName}}.TableName,
		} 
		{{end}}
		}
		`
)

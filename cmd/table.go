package cmd

var (
	//
	tableName = "TableName"
	//
	tableTpl = `
		package table

		type _{{.StructName}} struct {
			TableName string
			{{if .HasPrimaryKey}}PrimaryKey TableField{{end}}
		{{range $key, $value := .Columns}} {{ $key }} TableField 
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"
		{{if .HasPrimaryKey}}
			{{.StructName}}.PrimaryKey = TableField{
				Name: "{{index .PrimaryKey 0}}",
				Json: "{{index .PrimaryKey 1}}",
				Table: {{$.StructName}}.TableName,
			}
		{{end}}
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

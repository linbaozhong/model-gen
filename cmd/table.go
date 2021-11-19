package cmd

var (
	//
	tableName = "TableName"
	//
	tableTpl = `
		package table

		type _{{.StructName}} struct {
			TableName string
			ColumnNames        []string          //列名
			ColumnName2Comment map[string]string //列名和列描述映射
			ColumnName2Json    map[string]string //列名和JSON Key映射
			{{if .HasPrimaryKey}}PrimaryKey TableField{{end}}
		{{range $key, $value := .Columns}} {{ $key }} TableField 
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"
			{{ $.StructName}}.ColumnNames = make([]string,0,{{len .Columns}})
			{{.StructName}}.ColumnName2Json = make(map[string]string,{{len .Columns}})
			{{.StructName}}.ColumnName2Comment = make(map[string]string,{{len .Columns}})

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
			//Comment: "{{index $value 3}}",
			Table: {{$.StructName}}.TableName,
		} 
		{{ $.StructName}}.ColumnNames = append({{ $.StructName}}.ColumnNames,"{{index $value 0}}")
		{{ $.StructName}}.ColumnName2Json["{{index $value 0}}"] = "{{index $value 1}}"
		{{ $.StructName}}.ColumnName2Comment["{{index $value 0}}"] = "{{index $value 3}}"
		{{end}}
		}
		`
)

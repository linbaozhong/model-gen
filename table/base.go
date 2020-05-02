package table

const (
	Quote_Char = "`"
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
func (f TableField) Ue() string {
	return f.generate("<>")
}

//Bt BETWEEN
func (f TableField) Bt() string {
	return f.QuoteName() + " BETWEEN ? AND ?"
}

//In IN
func (f TableField) In() string {
	return f.QuoteName() + " IN (?)"
}

func (f TableField) QuoteName() string {
	return Quote_Char + f.Name + Quote_Char
}
func (f TableField) generate(op string) string {
	return f.QuoteName() + op + "?"
}

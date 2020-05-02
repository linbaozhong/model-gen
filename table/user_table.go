package table

type _User struct {
	TableName string
	Age       TableField
	ID        TableField
	Name      TableField
	NickName  TableField
}

var (
	User _User
)

func init() {
	User.TableName = "user"

	User.Age = TableField{
		Name: "age",
		Json: "age",
	}

	User.ID = TableField{
		Name: "id",
		Json: "id",
	}

	User.Name = TableField{
		Name: "name",
		Json: "name",
	}

	User.NickName = TableField{
		Name: "nick_name",
		Json: "nick_name",
	}

}

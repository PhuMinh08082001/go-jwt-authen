package model

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID       int64
	UserName string
	Password string
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

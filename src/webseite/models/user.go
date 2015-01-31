package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id       int32 `orm:"auto"`
	Email    string
	Salt     string
	PassHash string
	Avatar   string
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User))
}

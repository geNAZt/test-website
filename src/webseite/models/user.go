package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id   int
	Name string
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User))
}

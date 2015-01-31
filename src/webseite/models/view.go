package models

import "github.com/astaxie/beego/orm"

type View struct {
	Id      int32 `orm:"auto"`
	Name    string
	Servers []*Server `orm:"reverse(many)"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(View))
}

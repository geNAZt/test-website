package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Ping struct {
	Id     int32   `orm:"auto"`
	Server *Server `orm:"rel(fk)"`
	Online int32
	Time   time.Time `orm:"auto_now"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Ping))
}

package models

import "github.com/astaxie/beego/orm"

type Server struct {
	Id                      int32 `orm:"auto"`
	Name                    string
	Website                 string
	Ip                      string
	Record                  int32
	DownloadAnimatedFavicon bool
	Pings                   []*Ping `orm:"reverse(many)"`
	Views                   []*View `orm:"rel(m2m)"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Server))
}

package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	_ "net/http/pprof"
	_ "webseite/controllers/websocket"
	_ "webseite/models"
	_ "webseite/routers"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:80", nil))
	}()

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "go:go@/go?charset=utf8", 50, 300)
}

func main() {
	o := orm.NewOrm()
	o.Using("default") // Using default, you can use other database

	// Database alias.
	name := "default"

	// Drop table and re-create.
	force := false

	// Print log.
	verbose := true

	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	beego.Run()
}

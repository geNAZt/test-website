package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	_ "net/http/pprof"
	_ "webseite/controllers/websocket"
	_ "webseite/models"
	_ "webseite/routers"
)

func init() {
	// Start the profiler if needed
	if v, err := beego.AppConfig.Bool("ProfilerOn"); err == nil && v == true {
		address := beego.AppConfig.String("ProfilerIP") + ":" + beego.AppConfig.String("ProfilerPort")

		beego.BeeLogger.Info("Booting up profiler http server at " + address)

		go func() {
			beego.BeeLogger.Error("Error in Profiler HTTP Server: %v", http.ListenAndServe(address, nil))
		}()

		beego.BeeLogger.Info("Profiler HTTP Server up at " + address)
	}

	// Get the DB Configuration
	dbDriver := beego.AppConfig.String("DBDriver")
	dbUser := beego.AppConfig.String("DBUser")
	dbPass := beego.AppConfig.String("DBPass")
	dbHost := beego.AppConfig.String("DBHost")
	dbDatabase := beego.AppConfig.String("DBDatabase")

	dbMinPool, err := beego.AppConfig.Int("DBMinPool")
	if err != nil {
		beego.BeeLogger.Warn("No DBMinPool configured. Using DBMinPool of 5 as default")
		dbMinPool = 5
	}

	dbMaxConnections, err := beego.AppConfig.Int("DBMaxConnections")
	if err != nil {
		beego.BeeLogger.Warn("No DBMaxConnections configured. Using DBMaxConnections of 10 as default")
		dbMaxConnections = 10
	}

	// Check which Database we use
	if dbDriver == "mysql" {
		beego.BeeLogger.Info("Using %s as Database Driver", dbDriver)
		orm.RegisterDriver("mysql", orm.DR_MySQL)

		beego.BeeLogger.Info("Connecting to MySQL Server: %s@%s/%s using a MinPool of %d, MaxConnections of %d", dbUser, dbHost, dbDatabase, dbMinPool, dbMaxConnections)
		orm.RegisterDataBase("default", "mysql", dbUser+":"+dbPass+"@"+dbHost+"/"+dbDatabase+"?charset=utf8", dbMinPool, dbMaxConnections)
	}

	// TODO add support for more database drivers
}

func main() {
	// Check if we are in dev env
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	// Error.
	err := orm.RunSyncdb("default", false, false)
	if err != nil {
		fmt.Println(err)
	}

	// Run the Webserver
	beego.Run()
}

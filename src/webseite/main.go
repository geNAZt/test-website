package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	_ "net/http/pprof"
	"os"
	_ "webseite/controllers/websocket"
	_ "webseite/models"
	_ "webseite/routers"
	"webseite/storage"
	_ "webseite/template"
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
	beego.BeeLogger.Info("Using %s as Database Driver", dbDriver)
	if dbDriver == "mysql" {
		orm.RegisterDriver(dbDriver, orm.DR_MySQL)
		beego.BeeLogger.Info("Connecting to MySQL Server: %s@%s/%s using a MinPool of %d, MaxConnections of %d", dbUser, dbHost, dbDatabase, dbMinPool, dbMaxConnections)
		orm.RegisterDataBase("default", dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbDatabase+"?charset=utf8", dbMinPool, dbMaxConnections)
	} else {
		panic("Currently there is only MySQL Support. Cya..")
	}

	file, err := os.OpenFile("geNAZt.jpg", 0, 0666)

	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, stat.Size())
	file.Read(buffer)
	file.Close()

	storage := storage.GetStorage()
	storage.Store(buffer, "avatar/geNAZt.jpg")

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

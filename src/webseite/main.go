package main

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	_ "webseite/controllers/websocket"
	_ "webseite/models"
	_ "webseite/routers"
	"webseite/storage"
	"webseite/tasks"
	wtemp "webseite/template"
)

func init() {
	// GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

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

	// Build up the text template engine so we can parse css/js
	funcMap := template.FuncMap{
		"asset": wtemp.AssetResolver,
	}

	// Be sure that any static content is inside the storage
	storage := storage.GetStorage()
	errWalk := filepath.Walk("static/", func(path string, info os.FileInfo, err error) error {
		strippedPath := strings.Replace(path, "static"+string(filepath.Separator), "", -1)

		if !info.IsDir() && !storage.Exists(strippedPath) {
			beego.BeeLogger.Info("Storing file " + strippedPath + " into given storage")

			file, err := os.OpenFile(path, 0, 0666)
			if err != nil {
				return err
			}

			buffer := make([]byte, info.Size())
			read, errRead := file.Read(buffer)
			if errRead != nil {
				return errRead
			}

			errClose := file.Close()
			if errClose != nil {
				return errClose
			}

			if read == 0 {
				beego.BeeLogger.Info("Skipping 0 byte file " + strippedPath)
				return nil
			}

			ext := filepath.Ext(strippedPath)
			if ext == ".css" || ext == ".js" {
				contentString := string(buffer)
				templateEngine := template.New(strippedPath)
				templateEngine.Delims("![", "]!")
				templateEngine.Funcs(funcMap)
				template, errTemplate := templateEngine.Parse(contentString)
				if errTemplate != nil {
					return errTemplate
				}

				newbytes := bytes.NewBufferString("")
				errTemplateExecute := template.Execute(newbytes, nil)
				if errTemplateExecute != nil {
					return errTemplateExecute
				}

				buffer = newbytes.Bytes()
			}

			success, errStore := storage.Store(buffer, strippedPath)
			if errStore != nil {
				return errStore
			}

			if !success {
				beego.BeeLogger.Info("Storage could not store file " + strippedPath)
			}
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}
}

func main() {
	// Error.
	err := orm.RunSyncdb("default", false, false)
	if err != nil {
		fmt.Println(err)
	}

	// Build up and init all tasks we need
	tasks.InitTasks()

	// Run the Webserver
	beego.Run()
}

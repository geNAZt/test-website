package storage

import "github.com/astaxie/beego"

type Storage interface {
	Store([]byte, string) (bool, error)
}

var useOpenStack bool

func init() {
	useOpenStack = false

	if v, err := beego.AppConfig.Bool("OpenStackOn"); err == nil && v == true {
		useOpenStack = true
	}
}

func GetStorage() Storage {
	if useOpenStack {
		return new(openStackStorage)
	}

	return new(fileSystemStorage)
}

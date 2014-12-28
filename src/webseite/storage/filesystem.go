package storage

import (
	"github.com/astaxie/beego"
	"os"
)

type fileSystemStorage struct {
}

func init() {
	exist := exists(beego.AppConfig.String("FilestorageDir"))
	if !exist {
		os.MkdirAll(beego.AppConfig.String("FilestorageDir"), 0666)
	}
}

func (s *fileSystemStorage) Store(bytes []byte, filename string) (bool, error) {
	file, err := os.Create(beego.AppConfig.String("FilestorageDir") + "/" + filename)
	if err != nil {
		return false, err
	}

	written, err := file.Write(bytes)
	if err != nil || written != len(bytes) {
		return false, err
	}

	errClose := file.Close()
	if errClose != nil {
		return false, errClose
	}

	return true, nil
}

func exists(path string) bool {
	stat, err := os.Stat(path)

	if err == nil {
		return stat.IsDir()
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

package storage

import (
	"errors"
	"github.com/astaxie/beego"
	"os"
)

type fileSystemStorage struct {
}

var (
	storageDir string
	staticUrl  string
)

func init() {
	storageDir = beego.AppConfig.String("FilestorageDir")
	staticUrl = beego.AppConfig.String("FilestorageUrl")

	exist := exists(storageDir)
	if !exist {
		os.MkdirAll(storageDir, 0666)
	}

	beego.SetStaticPath(staticUrl, storageDir)
}

func (s *fileSystemStorage) Store(bytes []byte, filename string) (bool, error) {
	file, err := os.Create(storageDir + "/" + filename)
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

func (s *fileSystemStorage) Exists(filename string) bool {
	stat, err := os.Stat(storageDir + "/" + filename)
	return err != nil && stat.Size() > 0
}

func (s *fileSystemStorage) GetUrl(filename string) (string, error) {
	if s.Exists(filename) {
		return staticUrl + "/" + filename, nil
	}

	return "", errors.New("Not found")
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

package storage

import (
	"errors"
	"github.com/astaxie/beego"
	"os"
	"path"
	"webseite/cache"
)

type fileSystemStorage struct {
}

var (
	storageDir    string
	staticUrl     string
	fsExistsCache *cache.TimeoutCache
)

func init() {
	storageDir = beego.AppConfig.String("FilestorageDir")
	staticUrl = beego.AppConfig.String("FilestorageUrl")

	exist := exists(storageDir)
	if !exist {
		os.MkdirAll(storageDir, 0666)
	}

	// Build up cache
	tempCache, err := cache.NewTimeoutCache(60)
	if err != nil {
		panic(err)
	}

	fsExistsCache = tempCache

	beego.SetStaticPath(staticUrl, storageDir)
}

func (s *fileSystemStorage) Store(bytes []byte, filename string) (bool, error) {
	fullPath := storageDir + "/" + filename
	fullDir := path.Dir(fullPath)

	if !exists(fullDir) {
		os.MkdirAll(fullDir, 0666)
	}

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
	value, ok := fsExistsCache.Get(filename)
	if !ok {
		stat, err := os.Stat(storageDir + "/" + filename)
		exist := err == nil && stat != nil && stat.Size() > 0
		fsExistsCache.Add(filename, exist)
		return exist
	}

	return value.(bool)
}

func (s *fileSystemStorage) Delete(filename string) (bool, error) {
	err := os.Remove(filename)
	if err == nil {
		return true, nil
	}

	return false, err
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

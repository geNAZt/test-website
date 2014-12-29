package template

import (
	"github.com/astaxie/beego"
	"webseite/storage"
)

func init() {
	beego.AddFuncMap("asset", AssetResolver)
}

func AssetResolver(filename string) string {
	storage := storage.GetStorage()
	url, err := storage.GetUrl(filename)
	if err != nil {
		beego.BeeLogger.Info("Could not resolve asset: %s", filename)
		return ""
	}

	return url
}

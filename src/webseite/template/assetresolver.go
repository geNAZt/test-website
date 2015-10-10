package template

import (
	"github.com/astaxie/beego"
	"reflect"
	"webseite/storage"
)

func init() {
	beego.AddFuncMap("asset", AssetResolver)
	beego.AddFuncMap("isset", IsSet)
    beego.AddFuncMap("getAvatar", GetAvatar)
}

func GetAvatar(url string) string {
    if url == "default" {
        return AssetResolver("/avatar/default.png")
    }

    return AssetResolver("/avatar/default.png")
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

func IsSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > int64(kv.Int()) {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

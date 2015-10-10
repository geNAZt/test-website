package template

import (
	"github.com/astaxie/beego"
	"reflect"
	"webseite/storage"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func init() {
	beego.AddFuncMap("asset", AssetResolver)
	beego.AddFuncMap("isset", IsSet)
    beego.AddFuncMap("getAvatar", GetAvatar)
}

func GetAvatar(url string) string {
    if url == "default" {
        return AssetResolver("avatar/default.png")
    } else if strings.HasPrefix(url, "gravatar:") {
		// Generate md5 hash out of the email appended in the avatar
		split := strings.Split(url, ":")

		hasher := md5.New()
		hasher.Write([]byte(split[1]))
		return "http://www.gravatar.com/avatar/" + hex.EncodeToString(hasher.Sum(nil)) + "?s=40"
	}

    return AssetResolver("avatar/default.png")
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

package webapi

import "fmt"

type QrcodeTokenData struct {
	Uid    string `json:"uid"`
	Time   int64  `json:"time"`
	Sign   string `json:"sign"`
	Qrcode string `json:"qrcode"`
}

type QrcodeStatusData struct {
	Status  int    `json:"status,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Version string `json:"version,omitempty"`
}

const (
	_ApiFormatQrcodeToken = "https://qrcodeapi.115.com/api/1.0/%s/1.0/token"
	_ApiFormatQrcodeLogin = "https://passportapi.115.com/app/1.0/%s/1.0/login/qrcode"
	_UrlFormatQrcodeImage = "https://qrcodeapi.115.com/api/1.0/%s/1.0/qrcode?qrfrom=1&client=%d&uid=%s"
)

var (
	_AppIdMap = map[string]int{
		"web": 0,
		// Client ID for app is always 7
		"mac":     7,
		"linux":   7,
		"windows": 7,
	}
)

func QrcodeTokenApi(appType string) string {
	return fmt.Sprintf(_ApiFormatQrcodeToken, appType)
}

func QrcodeLoginApi(appType string) string {
	return fmt.Sprintf(_ApiFormatQrcodeLogin, appType)
}

func QrcodeImageUrl(appType, userId string) string {
	appId := _AppIdMap[appType]
	return fmt.Sprintf(_UrlFormatQrcodeImage, appType, appId, userId)
}
package api

type AppVersionInfo struct {
	CreatedTime int64  `json:"created_time"`
	VersionCode string `json:"version_code"`
	VersionUrl  string `json:"version_url"`
}

type AppVersionResult struct {
	Android    AppVersionInfo `json:"android"`
	LinuxApp   AppVersionInfo `json:"linux_115"`
	MacBrowser AppVersionInfo `json:"mac"`
	MacApp     AppVersionInfo `json:"mac_115"`
	WinBrowser AppVersionInfo `json:"win"`
	WinApp     AppVersionInfo `json:"window_115"`
}

type AppVersionSpec struct {
	_JsonpApiSpec[AppVersionResult, StandardResp]
}

func (s *AppVersionSpec) Init() *AppVersionSpec {
	s._JsonpApiSpec.Init(
		"https://appversion.115.com/1/web/1.0/api/chrome", "get_version",
	)
	return s
}

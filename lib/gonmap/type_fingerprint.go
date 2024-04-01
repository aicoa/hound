package gonmap

type FingerPrint struct {
	ProbeName        string
	MatchRegexString string
	Service          string
	ProductName      string
	Version          string
	Info             string
	Hostname         string
	Os               string
	DeviceType       string
}

func parseVersionInfo(s string) *FingerPrint {
	m := &match{} // 创建一个临时 match 对象以调用 getVersionInfo 方法
	return &FingerPrint{
		ProductName: m.getVersionInfo(s, "PRODUCTNAME"),
		Version:     m.getVersionInfo(s, "VERSION"),
		Info:        m.getVersionInfo(s, "INFO"),
		Hostname:    m.getVersionInfo(s, "HOSTNAME"),
		Os:          m.getVersionInfo(s, "Os"),
		DeviceType:  m.getVersionInfo(s, "DEVICE"),
	}
}

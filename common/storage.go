package common

var (
	ScanDone, TaskBrake bool
	// 端口扫描组别
	Database   string = "1433,1521,3306,5432,6379,9200,11211,27017"
	Enterprise string = "21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880"
	All        string = "1-65535"
	HighRisk   string = "21,22,23,53,80,443,8080,8000,139,445,3389,1521,3306,6379,7001,2375,27017,11211"
	// 资产收集数组
	HoldAsset   = [][]string{{"公司名称", "股权比例", "投资数额", "域名", ""}}
	WechatAsset = [][]string{{"公众号名称", "微信号", "公众号LOGO", "二维码", "简介", ""}}
	HunterAsset = [][]string{{"公司名称 | 域名", "资产数量"}}
	// WEB扫描
	VulHeader       = []string{"#", "漏洞名称", "风险等级", "漏洞地址", "详情"}
	ScanResult      = [][]string{{"#", "网站地址", "状态码", "长度", "标题", "指纹", ""}} // WebScan
	PortBurstResult = [][]string{{"协议", "主机", "用户名", "密码", ""}}
	SubdomainResult = [][]string{{"#", "子域名", "IP", ""}}
	// 扫描到指纹后需要扫描到的URL以及对应的指纹
	UrlFingerMap = map[string][]string{}
	// IP 域名目标
	WaitPortScan []string
	Userdict     = map[string][]string{
		"ftp":        {"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
		"telnet":     {"root", "admin"},
		"mysql":      {"root", "mysql"},
		"mssql":      {"sa", "sql"},
		"smb":        {"administrator", "admin", "guest"},
		"rdp":        {"administrator", "admin", "guest"},
		"postgresql": {"postgres", "admin"},
		"vnc":        {"admin", "administrator", "root"},
		"redis":      {},
		"ssh":        {"root", "admin", "ssh"},
		"mongodb":    {"root", "admin"},
		"oracle":     {"sys", "system", "admin", "test", "web", "orcl"},
		"weblogic":   {"system", "weblogic", "admin", "wlcsystem", "wlsystem", "wlcsystem", "wladmin", "wluser", "wlservice", "wljmx", "wljms", "wlservlet", "wldeploy", "joe", "wlinternal", "wlportal", "wlapp", "wlw", "wlwlang", "portaladmin", "portal", "mary"},
		"致远OA":       {"system", "group-admin", "admin1", "audit-admin"},
	}
	Passwords = []string{"sipisystem", "wlcsystem", "security", "portaladmin", "guest", "123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "1234567", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789", "weblogic", "sys123"}
)

var (
	Vulnerability = []VulnerabilityInfo{}
)

type VulnerabilityInfo struct {
	Id        string
	Name      string
	RiskLevel string
	Url       string
	TransInfo
}

type TransInfo struct {
	Request  string
	Response string
	ExtInfo  string // 拓展信息
}

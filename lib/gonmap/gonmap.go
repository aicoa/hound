package gonmap

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// 指纹匹配所需要的数据均来自github
var nmap *Nmap

var ProbeCount = 0   //探测器数量
var MatchCount = 0   //匹配器数量
var UsedProbenum = 0 //已使用探测器数量
var UsedMatchnum = 0 //已使用匹配器数量

var logger = Logger(log.New(os.Stderr, "[Gonmap]", log.Ldate|log.Ltime|log.Lshortfile))

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func init() {
	initWithFilter(9)
}

func initWithFilter(filter int) {
	// 初始化NMAP探针库
	repairNMAPString()
	initializeNmap(filter)

	// 加载探针数据
	loadProbes()

	// 初始化探针映射
	initializePortProbeMap()

	// 修复fallback
	nmap.fixFallback()

	// 新增自定义指纹信息
	customNMAPMatch()

	// 优化检测逻辑及端口对应的默认探针
	optimizeNMAPProbes()

	// 排序
	sortProbes()

	// 输出统计数据状态
	statistical()
}

func statistical() {
	ProbeCount = len(nmap.probeSort)
	MatchCount = 0 // 重置全局变量，避免累加
	UsedProbenum = 0
	UsedMatchnum = 0

	for _, p := range nmap.ProbeNameMap {
		MatchCount += len(p.matchGroup)
	}

	// 确保portProbeMap[0]存在
	if probes, ok := nmap.portProbeMap[0]; ok {
		UsedProbenum = len(probes)
		for _, probeName := range probes {
			if probe, ok := nmap.ProbeNameMap[probeName]; ok {
				UsedMatchnum += len(probe.matchGroup)
			}
		}
	}
}

func initializeNmap(filter int) {
	nmap = &Nmap{
		exclude:            emptyPortList,
		ProbeNameMap:       make(map[string]*probe),
		probeSort:          []string{},
		portProbeMap:       make(map[int]ProbeList),
		filter:             filter,
		timeout:            time.Second,
		probeUsed:          emptyProbeList,
		bypassAllProbePort: []int{161, 137, 139, 135, 389, 443, 548, 1433, 6379, 1883, 5432, 1521, 3389, 3388, 3389, 33890, 33900},
		sslSecondProbeMap:  []string{"TCP_TerminalServerCookie", "TCP_TerminalServer"},
		allProbeMap:        []string{"TCP_GetRequest", "TCP_NULL", "TCP_JDWP"},
		sslProbeMap:        []string{"TCP_TLSSessionReq", "TCP_SSLSessionReq", "TCP_SSLv23SessionReq"},
	}
}

func loadProbes() {
	nmap.loads(nmapServiceProbes + nmapCustomizeProbes)
}

func initializePortProbeMap() {
	for i := 0; i <= 65535; i++ {
		nmap.portProbeMap[i] = []string{}
	}
}

func sortProbes() {
	nmap.sslSecondProbeMap = nmap.sortByRarity(nmap.sslSecondProbeMap)
	nmap.allProbeMap = nmap.sortByRarity(nmap.allProbeMap)
	nmap.sslProbeMap = nmap.sortByRarity(nmap.sslProbeMap)
	for index, value := range nmap.portProbeMap {
		nmap.portProbeMap[index] = nmap.sortByRarity(value)
	}
}

// 其他辅助函数的定义...
func optimizeNMAPProbes() {
	nmap.ProbeNameMap["TCP_GenericLines"].sslports = nmap.ProbeNameMap["TCP_GenericLines"].sslports.append(993, 994, 456, 995)
	//优化检测逻辑，及端口对应的默认探针
	nmap.portProbeMap[993] = append([]string{"TCP_GenericLines"}, nmap.portProbeMap[993]...)
	nmap.portProbeMap[994] = append([]string{"TCP_GenericLines"}, nmap.portProbeMap[994]...)
	nmap.portProbeMap[995] = append([]string{"TCP_GenericLines"}, nmap.portProbeMap[995]...)
	nmap.portProbeMap[465] = append([]string{"TCP_GenericLines"}, nmap.portProbeMap[465]...)
	nmap.portProbeMap[3390] = append(nmap.portProbeMap[3390], "TCP_TerminalServer")
	nmap.portProbeMap[3390] = append(nmap.portProbeMap[3390], "TCP_TerminalServerCookie")
	nmap.portProbeMap[33890] = append(nmap.portProbeMap[33890], "TCP_TerminalServer")
	nmap.portProbeMap[33890] = append(nmap.portProbeMap[33890], "TCP_TerminalServerCookie")
	nmap.portProbeMap[33900] = append(nmap.portProbeMap[33900], "TCP_TerminalServer")
	nmap.portProbeMap[33900] = append(nmap.portProbeMap[33900], "TCP_TerminalServerCookie")
	nmap.portProbeMap[7890] = append(nmap.portProbeMap[7890], "TCP_Socks5")
	nmap.portProbeMap[7891] = append(nmap.portProbeMap[7891], "TCP_Socks5")
	nmap.portProbeMap[4000] = append(nmap.portProbeMap[4000], "TCP_Socks5")
	nmap.portProbeMap[2022] = append(nmap.portProbeMap[2022], "TCP_Socks5")
	nmap.portProbeMap[6000] = append(nmap.portProbeMap[6000], "TCP_Socks5")
	nmap.portProbeMap[7000] = append(nmap.portProbeMap[7000], "TCP_Socks5")
	//将TCP_GetRequest的fallback参数设置为NULL探针，避免漏资产
	nmap.ProbeNameMap["TCP_GenericLines"].fallback = "TCP_NULL"
	nmap.ProbeNameMap["TCP_GetRequest"].fallback = "TCP_NULL"
	nmap.ProbeNameMap["TCP_TerminalServerCookie"].fallback = "TCP_GetRequest"
	nmap.ProbeNameMap["TCP_TerminalServer"].fallback = "TCP_GetRequest"
}
func repairNMAPString() {
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, "${backquote}", "`")
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `q|GET / HTTP/1.0\r\n\r\n|`,
		`q|GET / HTTP/1.0\r\nHost: {Host}\r\nUser-Agent: Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)\r\nAccept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2\r\nAccept: */*\r\n\r\n|`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `\1`, `$1`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?=\\)`, `(?:\\)`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?=[\w._-]{5,15}\r?\n$)`, `(?:[\w._-]{5,15}\r?\n$)`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?:[^\r\n]*r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?:[^\r\n]*\r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?:[^\r\n]+\r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!2526)`, ``)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!400)`, ``)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!\0\0)`, ``)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!/head>)`, ``)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!HTTP|RTSP|SIP)`, ``)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!.*[sS][sS][hH]).*`, `.*`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!\xff)`, `.`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?!x)`, `[^x]`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?<=.)`, `(?:.)`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `(?<=\?)`, `(?:\?)`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `\x20\x02\x00.`, `\x20\x02..`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `match rtmp`, `# match rtmp`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `nmap`, `pamn`)
	nmapServiceProbes = strings.ReplaceAll(nmapServiceProbes, `Nmap`, `pamn`)
}

func customNMAPMatch() {
	//新增自定义指纹信息
	nmap.AddMatch("TCP_GetRequest", `echo m|^GET / HTTP/1.0\r\n\r\n$|s`)
	nmap.AddMatch("TCP_GetRequest", `mongodb m|.*It looks like you are trying to access MongoDB.*|s p/MongoDB/`)
	nmap.AddMatch("TCP_GetRequest", `http m|^HTTP/1\.[01] \d\d\d (?:[^\r\n]+\r\n)*?Server: ([^\r\n]+)| p/$1/`)
	nmap.AddMatch("TCP_GetRequest", `http m|^HTTP/1\.[01] \d\d\d|`)
	nmap.AddMatch("TCP_NULL", `mysql m|.\x00\x00..j\x04Host '.*' is not allowed to connect to this MariaDB server| p/MariaDB/`)
	nmap.AddMatch("TCP_NULL", `mysql m|.\x00\x00..j\x04Host '.*' is not allowed to connect to this MySQL server| p/MySQL/`)
	nmap.AddMatch("TCP_NULL", `mysql m|.\x00\x00\x00\x0a(\d+\.\d+\.\d+)\x00.*caching_sha2_password\x00| p/MariaDB/ v/$1/`)
	nmap.AddMatch("TCP_NULL", `mysql m|.\x00\x00\x00\x0a(\d+\.\d+\.\d+)\x00.*caching_sha2_password\x00| p/MariaDB/ v/$1/`)
	nmap.AddMatch("TCP_NULL", `mysql m|.\x00\x00\x00\x0a([\d.-]+)-MariaDB\x00.*mysql_native_password\x00| p/MariaDB/ v/$1/`)
	nmap.AddMatch("TCP_NULL", `redis m|-DENIED Redis is running in.*| p/Redis/ i/Protected mode/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Welcome to visit (.*) series router!.*|s p/$1 Router/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^Username: ??|`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Telnet service is disabled or Your telnet session has expired due to inactivity.*|s i/Disabled/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Telnet connection from (.*) refused.*|s i/Refused/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Command line is locked now, please retry later.*\x0d\x0a\x0d\x0a|s i/Locked/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Warning: Telnet is not a secure protocol, and it is recommended to use Stelnet.*|s`)
	nmap.AddMatch("TCP_NULL", `telnet m|^telnetd:|s`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Quopin CLI for (.*)\x0d\x0a\x0d\x0a|s p/$1/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^\x0d\x0aHello, this is FRRouting \(version ([\d.]+)\).*|s p/FRRouting/ v/$1/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*User Access Verification.*Username:|s`)
	nmap.AddMatch("TCP_NULL", `telnet m|^Connection failed.  Windows CE Telnet Service cannot accept anymore concurrent users.|s o/Windows/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^\x0d\x0a\x0d\x0aWelcome to the host.\x0d\x0a.*|s o/Windows/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^.*Welcome Visiting Huawei Home Gateway\x0d\x0aCopyright by Huawei Technologies Co., Ltd.*Login:|s p/Huawei/`)
	nmap.AddMatch("TCP_NULL", `telnet m|^..\x01..\x03..\x18..\x1f|s p/Huawei/`)
	nmap.AddMatch("TCP_NULL", `smtp m|^220 ([a-z0-1.-]+).*| h/$1/`)
	nmap.AddMatch("TCP_NULL", `ftp m|^220 H3C Small-FTP Server Version ([\d.]+).* | p/H3C Small-FTP/ v/$1/`)
	nmap.AddMatch("TCP_NULL", `ftp m|^421[- ]Service not available..*|`)
	nmap.AddMatch("TCP_NULL", `ftp m|^220[- ].*filezilla.*|i p/FileZilla/`)
	nmap.AddMatch("TCP_TerminalServerCookie", `ms-wbt-server m|^\x03\0\0\x13\x0e\xd0\0\0\x124\0\x02.*\0\x02\0\0\0| p/Microsoft Terminal Services/ o/Windows/ cpe:/o:microsoft:windows/a`)
	nmap.AddMatch("TCP_redis-server", `redis m|^.*redis_version:([.\d]+)\n|s p/Redis key-value store/ v/$1/ cpe:/a:redislabs:redis:$1/`)
	nmap.AddMatch("TCP_redis-server", `redis m|^-NOAUTH Authentication required.|s p/Redis key-value store/`)
}

var regexpFirstNum = regexp.MustCompile(`^\d`)

func FixProtocol(oldProtocol string) string {
	// 使用switch语句简化条件判断
	switch oldProtocol {
	case "ssl/http":
		return "https"
	case "http-proxy":
		return "http"
	case "ms-wbt-server":
		return "rdp"
	case "microsoft-ds":
		return "smb"
	case "netbios-ssn":
		return "netbios"
	case "oracle-tns":
		return "oracle"
	case "msrpc":
		return "rpc"
	case "ms-sql-s":
		return "mssql"
	case "domain":
		return "dns"
	case "svnserve":
		return "svn"
	case "ibm-db2":
		return "db2"
	case "socks-proxy":
		return "socks5"
	default:
		if strings.HasPrefix(oldProtocol, "ssl/") {
			return oldProtocol[4:] + "-ssl"
		}
		if regexpFirstNum.MatchString(oldProtocol) {
			return "S" + oldProtocol
		}
		return strings.ReplaceAll(oldProtocol, "_", "-")
	}
}

func GuessProtocol(port int) string {
	// 首先调用 nmapServices 函数获取切片
	services := nmapServices()

	// 然后检查端口号是否在 services 切片的范围内
	if port < 0 || port >= len(services) {
		return "http" // 或者返回其他默认协议或错误
	}

	protocol := services[port]
	if protocol == "unknown" {
		return "http"
	}
	return protocol
}

func SetFilter(filter int) {
	initWithFilter(filter)
}

func SetLogger(v Logger) {
	logger = v
}

// 功能类
func New() *Nmap {
	n := *nmap
	return &n
}

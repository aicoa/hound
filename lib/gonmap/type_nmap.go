package gonmap

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Nmap struct {
	exclude      PortList          // 排除的端口
	portProbeMap map[int]ProbeList // 端口到探针列表的映射
	portNameMap  map[string]*probe // 探针名到探针的映射
	ProbeNameMap map[string]*probe // 探针名到探针的映射
	probeSort    ProbeList         // 探针排序列表
	probeUsed    ProbeList         // 已使用的探针列表
	filter       int               // 过滤器
	timeout      time.Duration     // 超时时间
	//需要绕过所有探针的端口列表
	bypassAllProbePort PortList
	sslSecondProbeMap  ProbeList // ssl第二次探测的探针列表
	allProbeMap        ProbeList // 所有探针列表
	sslProbeMap        ProbeList // ssl探针列表
}

func (this *Nmap) getResponse(ip string, port int, tls bool, timeout time.Duration, p *probe) (status Status, response *Response) {
	if port == 53 {
		if DnsScan(ip, port) {
			return Matched, &dnsResPonse
		} else {
			return Closed, nil
		}
	}
	raw, tls, err := p.scan(ip, port, tls, timeout, 10240)
	if err != nil {
		if strings.Contains(err.Error(), "STEP1") {
			return Closed, nil
		} else if strings.Contains(err.Error(), "STEP2") {
			return Closed, nil
		}
		if p.protocol == "UDP" && strings.Contains(err.Error(), "refused") {
			return Closed, nil
		}
		return Open, nil
	}
	response = &Response{
		Raw:         raw,
		Tls:         tls,
		Fingerprint: &FingerPrint{},
	}
	fingerprint := this.getFinger(raw, tls, p.name)
	response.Fingerprint = fingerprint
	if fingerprint.Service == "" {
		return NoMatched, response
	} else {
		return Matched, response
	}
}

func (this *Nmap) getFinger(raw string, tls bool, requestName string) *FingerPrint {
	data := this.ConvertToUTF8(raw)
	probe := this.ProbeNameMap[requestName]
	finger := probe.match(data)
	if tls {
		if finger.Service == "http" {
			finger.Service = "https"
		}
	}
	if finger.Service != "" || this.ProbeNameMap[requestName].fallback == "" {
		//标记探针名称
		finger.ProbeName = requestName
		return finger
	}
	//尝试使用fallback指定的探针
	fallback := this.ProbeNameMap[requestName].fallback
	fallbackProbe := this.ProbeNameMap[fallback]
	for fallback != "" {
		finger = fallbackProbe.match(data)
		fallback = this.ProbeNameMap[fallback].fallback
		if finger.Service != "" {
			break
		}
	}
	//标记探针名称
	finger.ProbeName = requestName
	return finger
}
func (n *Nmap) ScanTimeOut(ip string, port int, timeout time.Duration) (status Status, response *Response) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	var resChan = make(chan bool)
	defer func() {
		close(resChan)
		cancel()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if fmt.Sprint(r) != "send on closed channel" {
					panic(r)
				}
			}
		}()
		status, response = n.Scan(ip, port)
		resChan <- true
	}()
	select {
	case <-ctx.Done():
		return Closed, nil
	case <-resChan:
		return status, response
	}
}

func (n *Nmap) Scan(ip string, port int) (status Status, response *Response) {
	var probeNames ProbeList
	if n.bypassAllProbePort.exist(port) {
		probeNames = append(n.portProbeMap[port], n.allProbeMap...)
	} else {
		probeNames = append(n.allProbeMap, n.portProbeMap[port]...)
	}
	probeNames = append(probeNames, n.sslProbeMap...)
	//去除重复tanzhen
	if len(probeNames) == 0 { // 检查 probeNames 是否为空
		return Closed, nil // 或者返回一个适当的错误
	}
	probeNames = probeNames.removeDuplicate()
	firtProbe := probeNames[0]
	status, response = n.getRealResponse(ip, port, n.timeout, firtProbe)

	if status == Closed || status == Matched {
		return status, response
	}
	otherProbes := probeNames[1:]
	return n.getRealResponse(ip, port, 2*time.Second, otherProbes...)
}
func (n *Nmap) getResponseByProbes(host string, port int, timeout time.Duration, probes ...string) (status Status, response *Response) {
	var responseNotMatch *Response
	for _, requestName := range probes {
		if n.probeUsed.exist(requestName) {
			continue
		}
		n.probeUsed = append(n.probeUsed, requestName)
		p := n.ProbeNameMap[requestName]

		status, response = n.getResponse(host, port, p.sslports.exist(port), timeout, p)

		if status == Closed || status == Matched {
			responseNotMatch = nil
			break
		}
		if status == NoMatched {
			responseNotMatch = response
		}
	}
	if responseNotMatch != nil {
		response = responseNotMatch
	}
	return status, response
}

func (n *Nmap) getRealResponse(host string, port int, timeout time.Duration, probes ...string) (status Status, response *Response) {
	status, response = n.getResponseByProbes(host, port, timeout, probes...)
	if status != Matched {
		return status, response
	}
	if response.Fingerprint.Service == "ssl" {
		status, response := n.getResponseBySSL2thProbe(host, port, timeout)
		if status == Matched {
			return Matched, response
		}
	}
	return status, response
}
func (this *Nmap) getResponseByHTTPS(ip string, port int, timeout time.Duration) (status Status, response *Response) {
	return this.getResponse(ip, port, true, timeout, this.ProbeNameMap["TCP_GetRequest"])
}

func (n *Nmap) getResponseBySSL2thProbe(ip string, port int, timeout time.Duration) (status Status, response *Response) {
	status, response = n.getResponseByProbes(ip, port, timeout, n.sslSecondProbeMap...)
	if status == Closed || response.Fingerprint.Service == "ssl" {
		status, response = n.getResponseByHTTPS(ip, port, timeout)
	}
	if status == Matched && response.Fingerprint.Service != "ssl" {
		if response.Fingerprint.Service == "http" {
			response.Fingerprint.Service = "https"
		}
		return Matched, response
	}
	return NoMatched, nil
}

// 考虑到存在别的编码的可能，增加适配
func (n *Nmap) ConvertToUTF8(s string) string {
	reader := strings.NewReader(s)
	transformer := transform.NewReader(reader, charmap.ISO8859_1.NewDecoder())
	bytes, err := io.ReadAll(transformer)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// 配置类
func (this *Nmap) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}
func (this *Nmap) OpenDeepIdentify() {
	//与Nmap同理，-sV参数进行深度解析
	this.allProbeMap = this.probeSort
}

func (this *Nmap) AddMatch(probeName string, expr string) {
	probe := this.ProbeNameMap[probeName]
	probe.loadMatch(expr, false)
}

// initialize class
func (this *Nmap) loads(str string) {
	lines := strings.Split(str, "\n")
	var probeGroups [][]string
	var probeLines []string
	for _, line := range lines {
		if !this.isCommand(line) {
			continue
		}
		commandName := line[:strings.Index(line, " ")]
		if commandName == "Exclude" {
			this.loadExclude(line)
			continue
		}
		if commandName == "Probe" { //probe命令
			if len(probeLines) > 0 {
				probeGroups = append(probeGroups, probeLines)
				probeLines = []string{}
			}
		}
		probeLines = append(probeLines, line)
	}
	probeGroups = append(probeGroups, probeLines)
	for _, probeLines := range probeGroups {
		probe := parseProbe(probeLines)
		this.addProbe(*probe)
	}
}

func (this *Nmap) loadExclude(s string) {
	if !portGroupRegx.MatchString(s) {
		panic("端口组参数不正确")
	}
	this.exclude = parsePortList(s)
}
func (this *Nmap) addProbe(p probe) {
	this.probeSort = append(this.probeSort, p.name)
	this.ProbeNameMap[p.name] = &p
	//建立端口扫描对应表，将根据端口号决定使用何种请求包
	//如果端口列表为空，则为全端口
	if p.rarity > this.filter {
		return
	}
	//用0来记录所有已使用的探针
	this.portProbeMap[0] = append(this.portProbeMap[0], p.name)
	for _, port := range p.ports {
		this.portProbeMap[port] = append(this.portProbeMap[port], p.name)
	}
	for _, port := range p.sslports {
		this.portProbeMap[port] = append(this.portProbeMap[port], p.name)
	}
}
func (n *Nmap) fixFallback() {
	for probeName, probeType := range n.ProbeNameMap {
		fallback := probeType.fallback
		if fallback == "" {
			continue
		}
		if _, ok := n.ProbeNameMap["TCP_"+fallback]; ok {
			n.ProbeNameMap[probeName].fallback = "TCP_" + fallback
		} else {
			n.ProbeNameMap[probeName].fallback = "UDP_" + fallback
		}
	}
}
func (this *Nmap) isCommand(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[:1] == "#" {
		return false
	}
	//去除异常命令
	commandName := s[:strings.Index(s, " ")]
	commandArr := []string{
		"Exclude", "Probe", "match", "softmatch", "ports", "sslports", "totalwaitms", "tcpwrappedms", "rarity", "fallback",
	}
	for _, item := range commandArr {
		if item == commandName {
			return true
		}
	}
	return false
}

func (n *Nmap) sortByRarity(list ProbeList) ProbeList {
	if len(list) == 0 {
		return list
	}

	// 创建一个包含探针名称和对应稀有度的结构的切片
	sortedList := make([]struct {
		name   string
		rarity int
	}, len(list))

	for i, probeName := range list {
		sortedList[i].name = probeName
		sortedList[i].rarity = n.ProbeNameMap[probeName].rarity
	}

	// 使用sort.Slice稳定排序，基于rarity字段
	sort.SliceStable(sortedList, func(i, j int) bool {
		return sortedList[i].rarity < sortedList[j].rarity
	})

	// 从排序后的结构中提取探针名称
	for i, v := range sortedList {
		list[i] = v.name
	}

	return list
}
func DnsScan(host string, port int) bool {
	domainServer := fmt.Sprintf("%s:%d", host, port)
	c := dns.Client{
		Timeout: 2 * time.Second,
	}
	m := dns.Msg{}
	// 最终都会指向一个ip 也就是typeA, 这样就可以返回所有层的cname.
	m.SetQuestion("www.baidu.com.", dns.TypeA)
	_, _, err := c.Exchange(&m, domainServer)
	if err != nil {
		return false
	}
	return true
}

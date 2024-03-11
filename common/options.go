package common

import (
	"hound/lib/poc"
	"hound/lib/util"
	"strings"
)

var (
	ReverseCeyeApiKey = "ba446c3277a60555ad9e74a6f0cb4290"
	ReverseCeyeDomain = "xrn0nb.ceye.io" //github上随意找

	ReverseCeyeLive bool
)

type Options struct {
	//拉到数组里面批量搞，下面都同理，看变量名应该就明白了
	PocsDirectory []string
	Targets       util.SafeSlice
	Target        []string

	// PoC file or directory to scan
	PocFile string

	// search PoC by keyword , eg: -s tomcat
	Search         string
	SearchKeywords []string

	// pocs to run based on severity. Possible values: info, low, medium, high, critical
	Severity string

	SeverityKeywords []string

	//扫描数量
	count int
	//当前扫描数量
	Currentnum uint32

	// maximum number of requests to send per second (default 150)
	RateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	Concurrency int

	// maximum number of requests to send per second (default 150)
	ReverseRateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	ReverseConcurrency int

	// number of times to retry a failed request (default 2)
	Retries int

	MaxHostError int

	// time to wait in seconds before timeout (default 9)
	Timeout int

	// http/socks5 proxy to use
	Proxy string

	MaxRespBodySize int
}

func NewOptions(target []string, keyword, severity, proxy string) *Options {
	options := &Options{}
	options.Target = target
	options.Search = keyword
	options.Severity = severity
	options.RateLimit = Profile.WebScan.Thread
	options.Concurrency = 25
	options.ReverseRateLimit = 50
	options.ReverseConcurrency = 20
	options.MaxRespBodySize = 2
	options.Retries = 2
	options.Timeout = DefaultWebTimeout
	options.MaxHostError = 4
	options.Proxy = proxy
	poc.SelectFolderReadPocPath(options.PocFile)
	return options
}

// 关键词数量判断
func (O *Options) SetSearchKeyword() bool {
	if len(O.Search) == 0 {
		return false
	} else {
		arr := strings.Split(O.Search, ",")
		if len(arr) >= 1 {
			for _, v := range arr {
				O.SearchKeywords = append(O.SearchKeywords, strings.TrimSpace(v))
			}
			return true
		}
	}
	return false
}

func (this *Options) CheckPocsKeyWord(id string) bool {
	if len(this.SearchKeywords) > 0 {
		for _, v := range this.SearchKeywords {
			v = strings.ToLower(v)
			if strings.Contains(strings.ToLower(id), v) {
				return true
			}
		}
	}
	return false
}

// 风险等级筛选功能部分
func (this *Options) SetServerityKeyword(id string) bool {
	if len(this.Severity) > 0 {
		arr := strings.Split(this.Severity, ",")
		if len(arr) >= 1 {
			this.SeverityKeywords = append(this.SeverityKeywords, arr...)
			return true
		}
	}
	return false
}

// 关键字筛选功能
func (this *Options) CheckPocsServerityKeyWord(serverity string) bool {
	if len(this.SearchKeywords) > 0 {
		for _, v := range this.SeverityKeywords {
			if strings.EqualFold(serverity, v) {
				return true
			}
		}
	}
	return false
}

// 对poc中的关键字进行筛选，搜索关键字筛选
func (this *Options) FilterPocServeritySearch(pocId, serverity string) bool {
	var isShowed bool
	if len(this.Search) > 0 && this.SetSearchKeyword() && len(o.serverity) > 0 && this.SetServerityKeyword() {
		isShowed = true
	} else if len(this.Severity) > 0 && this.SetServerityKeyword() {
		if this.CheckPocsServerityKeyWord(serverity) {
			isShowed = true
		}
	} else if len(this.Search) > 0 && this.SetSearchKeyword() {
		isShowed = true
	} else {
		isShowed = false
	}
	return isShowed
}

// 区分dnslog poc
func (this *Options) ReversePocs(allpocs []poc.Poc) ([]poc.Poc, []poc.Poc) {
	reverse := []poc.Poc{}
	other := []poc.Poc{}
	for _, v := range allpocs {
		flag := false
		for _, item := range poc.Set {
			key := item.Key.(string)
			if strings.EqualFold(key, "reverse") {
				flag = true
				break
			}
		}
		if flag {
			reverse = append(reverse, v)
		} else {
			other = append(other, v)
		}
	}
	return reverse, other
}

// 自建poc功能部分
func (this *Options) CreatePocsList(active, finger2Poc bool) []poc.Poc {
	var pocSlice []poc.Poc
	if len(poc.LocalTestList) > 0 {
		for _, YamlOfPoc := range poc.LocalTestList {
			if p, err := poc.LocalReadPocByPath(YamlOfPoc); err == nil {
				pocSlice = append(pocSlice, p)
			}
		}
	}
	// 如果开启仅指纹扫描，或者指纹poc扫描，则屏蔽其他poc，且优先级大于扫描指纹的poc
	if !active || finger2Poc {
		for _, pocEmbedYaml := range poc.EmbedFileList {
			if p, err := poc.LocalReadPocByPath(pocEmbedYaml); err == nil {
				pocSlice = append(pocSlice, p)
			}
		}
	}

	newPocSlice := []poc.Poc{}
	for _, poc := range pocSlice {
		//筛选
		if this.FilterPocServeritySearch(poc.Id, poc.Info.Severity) {
			newPocSlice = append(newPocSlice, poc)
		}
	}
	return newPocSlice
}

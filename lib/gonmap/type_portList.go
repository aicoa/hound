package gonmap

import (
	"regexp"
	"strconv"
	"strings"
)

// 正则表达式：匹配单个端口号或端口范围。
var portRangeRegx = regexp.MustCompile(`^(\d+)(?:-(\d+))?$`)

// 正则表达式：匹配一组端口号或端口范围，由逗号分隔。
var portGroupRegx = regexp.MustCompile(`^(\d+(?:-\d+)?)(?:,\d+(?:-\d+)?)*$`)

// PortList 定义了端口列表的类型。
type PortList []int

// emptyPortList 提供了一个空的PortList实例。
var emptyPortList = PortList([]int{})

// parsePortList 解析一个表示端口范围的字符串，并返回一个PortList。
// 如果字符串不符合端口范围的格式，则会触发panic。
func parsePortList(express string) PortList {
	if !portGroupRegx.MatchString(express) {
		panic("port expression string invalid")
	}

	var list PortList
	for _, expr := range strings.Split(express, ",") {
		rArr := portRangeRegx.FindStringSubmatch(expr)
		startPort, _ := strconv.Atoi(rArr[1]) // 忽略错误，因为正则表达式已经验证了格式
		endPort := startPort
		if rArr[2] != "" {
			endPort, _ = strconv.Atoi(rArr[2]) // 同上
		}
		for num := startPort; num <= endPort; num++ {
			list = append(list, num)
		}
	}
	return list.removeDuplicate()
}

// removeDuplicate 移除PortList中的重复端口。
func (p PortList) removeDuplicate() PortList {
	result := make([]int, 0, len(p))
	temp := make(map[int]struct{})
	for _, item := range p {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// exist 检查PortList中是否存在指定的端口。
func (p PortList) exist(port int) bool {
	for _, num := range p {
		if num == port {
			return true
		}
	}
	return false
}

// append 向PortList中添加新的端口，并移除重复项。
func (p PortList) append(ports ...int) PortList {
	p = append(p, ports...)
	p = (p).removeDuplicate()
	return p
}

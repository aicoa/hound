package common

import (
	"fmt"
	"hound/lib/util"
	"net"
	"strings"
)

func ParseIPs(ipString string) (ips []string) {
	var excludedIPs []string
	if ipString != "" {
		for _, line := range ParseTargets(ipString, Mode_Other) {
			if strings.HasPrefix(line, "!") {
				excludedIPs = append(excludedIPs, ParseIP(line[1:])...)
			} else {
				ips = append(ips, ParseIP(line)...)
			}
		}
	}
	for _, excludedIP := range excludedIPs {
		ips = util.RemoveElement(ips, excludedIP)
	}
	return ips
}

func ParseIP(ipString string) []string {
	var result []string
	if strings.Contains(ipString, "-") {
		result = append(result, parseIPRange(ipString)...)
	} else if strings.Contains(ipString, ",") {
		ipArray := strings.Split(ipString, ",")
		for _, ip := range ipArray {
			if strings.Contains(ip, "/") {
				parsedIPs, err := parseCIDR(ip)
				if err != nil {
					fmt.Println("Error parsing CIDR:", err)
					continue
				}
				result = append(result, parsedIPs...)
			} else {
				result = append(result, ip)
			}
		}
	} else {
		result = append(result, ipString)
	}
	return result
}

func parseCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func parseIPRange(ipRange string) []string {
	ipParts := strings.Split(ipRange, "-")
	startIP := net.ParseIP(ipParts[0])
	endIP := net.ParseIP(ipParts[1])

	if endIP == nil || startIP == nil {
		fmt.Println("Invalid IP range format.")
		return nil
	}

	var ips []string
	for ip := startIP; !ip.Equal(endIP); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	ips = append(ips, endIP.String()) // Include the last IP

	return ips
}

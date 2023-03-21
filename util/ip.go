package util

import "net"

// IsPublicIP 判断Ip地址是否是公网Ip。是返回 true， 不是返回 false
func IsPublicIP(ipAddress string) bool {
	parseIp := net.ParseIP(ipAddress)
	if parseIp == nil || parseIp.IsPrivate() || parseIp.IsLoopback() {
		return false
	}
	return true
}

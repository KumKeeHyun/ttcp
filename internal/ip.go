package internal

import (
	"fmt"
	"net"
)

// intToIP converts IPv4 number to net.IP
func IntToIP(ipNum uint32) net.IP {
	ip := make(net.IP, 4)
	NativeEndian.PutUint32(ip, ipNum)
	return ip
}

// IPToInt converts net.IP to IPv4 number
func IPToInt(ip net.IP) uint32 {
	return NativeEndian.Uint32(ip)
}

func StringToIPv4(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("%s is not ip format", ipStr)
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("%s is not ipv4 format", ip)
	}
	return ipv4, nil
}
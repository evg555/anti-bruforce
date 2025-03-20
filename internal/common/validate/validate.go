package validate

import (
	"net"
	"strings"
)

func IsValidSubnet(subnet string) bool {
	if IsEmpty(subnet) {
		return false
	}
	_, _, err := net.ParseCIDR(subnet)
	return err == nil
}

func IsValidIPAddress(ip string) bool {
	if IsEmpty(ip) {
		return false
	}

	valid := net.ParseIP(ip)
	return valid != nil
}

func IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

package common

import "net"

func IsValidSubnet(subnet string) bool {
	_, _, err := net.ParseCIDR(subnet)
	return err == nil
}

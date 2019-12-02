package generic

import (
	"fmt"
	"net"
)

// CheckIPAddress determines if the hostname string is a valid IP address (either IPv4 or IPv6). If the validation fails, CheckIPAddress returns an error with a message that explains the validation failure. If validation succeeds, CheckIPAddress returns nil.
func CheckIPAddress(address string) error {
	ip := net.ParseIP(address)
	if ip == nil {
		return fmt.Errorf("not a valid IP address: '%s'", address)
	}

	return nil
}

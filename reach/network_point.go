package reach

import (
	"net"
	"strings"
)

// A NetworkPoint is a point of termination for an analyzed network vector (on either the source or destination side), such that there is no further subdivision of a source or destination possible beyond the network point. For example, the CIDR block "10.0.1.0/24" contains numerous individual IP addresses, and the analysis result might vary depending on which of these individual IP addresses is used in real network traffic. To break this problem down, such that an analysis result is as definitive as possible, each individual IP address must be analyzed, one at a time. Each IP address could be considered a network point, whereas the CIDR block could not be considered a network point.
type NetworkPoint struct {
	IPAddress net.IP
	Lineage   []ResourceReference
	Factors   []Factor
}

// String returns the text representation of the NetworkPoint
func (point NetworkPoint) String() string {
	var generations []string

	for i := len(point.Lineage) - 1; i >= 0; i-- {
		generations = append(generations, point.Lineage[i].ID)
	}

	generations = append(generations, point.IPAddress.String())

	return strings.Join(generations, " -> ")
}

// Domain returns the domain of the network point. It determines this by returning the domain of the first resource reference in the point's lineage. If the point has no lineage, an empty string is returned.
func (point NetworkPoint) Domain() string {
	for _, resourceRef := range point.Lineage {
		return resourceRef.Domain
	}

	return ""
}

// IPAddressIsInternetAccessible determines if the IP address for this network point can be accessed from the Internet.
func (point NetworkPoint) IPAddressIsInternetAccessible() bool {
	ip := point.IPAddress

	if ip.IsLoopback() || ip.IsUnspecified() { // TODO: Figure out other disqualifying criteria for an Internet-accessible IP address.
		return false
	}

	// For background on the following networks, see https://en.wikipedia.org/wiki/Private_network

	if cidrBlockContainsIP("10.0.0.0/8", ip) {
		return false
	}

	if cidrBlockContainsIP("172.16.0.0/12", ip) {
		return false
	}

	if cidrBlockContainsIP("192.168.0.0/16", ip) {
		return false
	}

	if cidrBlockContainsIP("100.64.0.0/10", ip) {
		return false
	}

	if cidrBlockContainsIP("fd00::/8", ip) {
		return false
	}

	if cidrBlockContainsIP("169.254.0.0/16", ip) {
		return false
	}

	if cidrBlockContainsIP("fc00::/7", ip) {
		return false
	}

	if cidrBlockContainsIP("fe80::/10", ip) {
		return false
	}

	return true
}

// IPv4 determines if the network point's IP address is an IPv4 address. If not, one can assume it's an IPv6 address.
func (point NetworkPoint) IPv4() bool {
	return point.IPAddress.To4() != nil
}

func (point NetworkPoint) trafficContents() []TrafficContent {
	var components []TrafficContent

	for _, factor := range point.Factors {
		components = append(components, factor.Traffic)
	}

	return components
}

func cidrBlockContainsIP(cidr string, ip net.IP) bool {
	_, network, _ := net.ParseCIDR(cidr)

	if network == nil {
		return false
	}

	return network.Contains(ip)
}

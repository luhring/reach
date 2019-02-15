package aws

import (
	"net"
	"strings"
)

// Instance is an internal representation of an EC2 instance, containing only properties relevant to this program
type Instance struct {
	ID                   string
	NameTag              string
	PrivateIPv4Addresses []net.IP
	PublicIPv4Addresses  []net.IP
	IPv6Addresses        []net.IP
	SecurityGroupIDs     []*string
	State                string
	SubnetID             string
	VpcID                string
}

func (instance *Instance) isRunning() bool {
	const runningStateText = "running"
	return strings.EqualFold(instance.State, runningStateText)
}

// GetFriendlyName returns the name tag if it exists, and otherwise returns the instance ID
func (instance *Instance) GetFriendlyName() string {
	nameTag := instance.NameTag

	if nameTag != "" {
		return nameTag
	}

	id := instance.ID

	if id != "" {
		return id
	}

	return "[unnamed instance]"
}

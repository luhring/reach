package set

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	minimumPort = 0
	maximumPort = 65535
)

type PortSet struct {
	set set
}

func NewEmptyPortSet() PortSet {
	return PortSet{
		set: newEmptySet(),
	}
}

func NewPortSetFromPortValue(port uint16) (*PortSet, error) {
	if err := validatePort(port); err != nil {
		return nil, fmt.Errorf("unable to use port: %v", err)
	}

	return &PortSet{
		set: NewSetFromSingleValue(port),
	}, nil
}

func NewPortSetFromRange(lowPort, highPort uint16) (*PortSet, error) {
	if err := validatePort(lowPort); err != nil {
		return nil, fmt.Errorf("unable to use lowPort: %v", err)
	}

	if err := validatePort(highPort); err != nil {
		return nil, fmt.Errorf("unable to use highPort: %v", err)
	}

	return &PortSet{
		set: newSetFromRange(lowPort, highPort),
	}, nil
}

func NewPortSetFromAWSPortRange(portRange *ec2.PortRange) (*PortSet, error) {
	if portRange == nil {
		return nil, fmt.Errorf("input portRange was nil")
	}

	from := aws.Int64Value(portRange.From)
	to := aws.Int64Value(portRange.To)

	return NewPortSetFromRange(uint16(from), uint16(to))
}

func NewPortSetFromAWSIPPermission(permission *ec2.IpPermission) (*PortSet, error) {
	if permission == nil {
		return nil, fmt.Errorf("input IpPermission was nil")
	}

	from := aws.Int64Value(permission.FromPort)
	to := aws.Int64Value(permission.ToPort)

	return NewPortSetFromRange(uint16(from), uint16(to))
}

func (s PortSet) Intersect(other PortSet) PortSet {
	return PortSet{
		set: s.set.intersect(other.set),
	}
}

func (s PortSet) Merge(other PortSet) PortSet {
	return PortSet{
		set: s.set.merge(other.set),
	}
}

// Subtract OTHER set from set (= set - other set)
func (s PortSet) Subtract(other PortSet) PortSet {
	return PortSet{
		set: s.set.subtract(other.set),
	}
}

func (s PortSet) String() string {
	return s.set.String()
}

func validatePort(port uint16) error {
	if port < minimumPort || port > maximumPort {
		return fmt.Errorf(
			"port number %v is not valid, must be between %v and %v",
			port,
			minimumPort,
			maximumPort,
		)
	}

	return nil
}

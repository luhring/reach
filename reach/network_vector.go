package reach

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

// A NetworkVector represents the path between two network points that's able to be analyzed in terms of what kind of network traffic is allowed to flow from point to point.
type NetworkVector struct {
	ID          string
	Source      NetworkPoint
	Destination NetworkPoint
	Traffic     *TrafficContent
}

// NewNetworkVector creates a new network vector given a source and a destination network point.
func NewNetworkVector(source, destination NetworkPoint) (NetworkVector, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return NetworkVector{}, err
	}

	return NetworkVector{
		ID:          u.String(),
		Source:      source,
		Destination: destination,
	}, nil
}

// String returns the text representation of a NetworkVector.
func (v NetworkVector) String() string {
	output := ""
	output += fmt.Sprintf("* network vector ID: %s\n", v.ID)
	output += fmt.Sprintf("* source network point: %s\n* destination network point: %s\n", v.Source.String(), v.Destination.String())

	if v.Traffic != nil {
		output += "\n"
		output += v.Traffic.String()
		output += "\n"
	}

	return output
}

// SourcePerspective returns an analyzable Perspective based on the NetworkVector's source network point.
func (v NetworkVector) SourcePerspective() Perspective {
	return Perspective{
		Self:      v.Source,
		Other:     v.Destination,
		SelfRole:  SubjectRoleSource,
		OtherRole: SubjectRoleDestination,
	}
}

// DestinationPerspective returns an analyzable Perspective based on the NetworkVector's destination network point.
func (v NetworkVector) DestinationPerspective() Perspective {
	return Perspective{
		Self:      v.Destination,
		Other:     v.Source,
		SelfRole:  SubjectRoleDestination,
		OtherRole: SubjectRoleSource,
	}
}

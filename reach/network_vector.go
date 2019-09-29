package reach

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

type NetworkVector struct {
	ID          string
	Source      NetworkPoint
	Destination NetworkPoint
	Traffic     *TrafficContent
}

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

func (v NetworkVector) TrafficComponents() []TrafficContent {
	var components []TrafficContent

	components = append(components, v.Source.TrafficComponents()...)
	components = append(components, v.Destination.TrafficComponents()...)

	return components
}

func (v NetworkVector) NetTraffic() (TrafficContent, error) {
	vectorTrafficComponents := v.TrafficComponents()

	resultingTraffic, err := NewTrafficContentFromIntersectingMultiple(vectorTrafficComponents)
	if err != nil {
		return TrafficContent{}, err
	}

	return resultingTraffic, nil
}

func (v NetworkVector) String() string {
	output := ""
	output += fmt.Sprintf("* network vector ID: %s\n", v.ID)
	output += fmt.Sprintf("* source: %s\n* destination: %s\n", v.Source.String(), v.Destination.String())

	if v.Traffic != nil {
		output += "\n"
		output += v.Traffic.String()
	}

	return output
}

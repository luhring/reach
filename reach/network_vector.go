package reach

import "github.com/nu7hatch/gouuid"

type NetworkVector struct {
	ID          string       `json:"id"`
	Source      NetworkPoint `json:"source"`
	Destination NetworkPoint `json:"destination"`
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

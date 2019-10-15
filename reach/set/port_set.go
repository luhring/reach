package set

import (
	"encoding/json"
	"fmt"
)

const (
	minimumPort = 0
	maximumPort = 65535
)

type PortSet struct {
	set Set
}

func NewEmptyPortSet() PortSet {
	return PortSet{
		set: newEmptySet(),
	}
}

func NewFullPortSet() PortSet {
	return PortSet{
		set: newCompleteSet(),
	}
}

func NewPortSetFromRange(lowPort, highPort uint16) (PortSet, error) {
	if err := validatePort(lowPort); err != nil {
		return PortSet{}, fmt.Errorf("unable to use lowPort: %v", err)
	}

	if err := validatePort(highPort); err != nil {
		return PortSet{}, fmt.Errorf("unable to use highPort: %v", err)
	}

	return PortSet{
		set: newSetFromRange(lowPort, highPort),
	}, nil
}

func (s PortSet) Complete() bool {
	return s.set.Complete()
}

func (s PortSet) Empty() bool {
	return s.set.Empty()
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

func (s PortSet) RangeStrings() []string {
	return s.set.rangeStrings()
}

func (s PortSet) String() string {
	return s.set.String()
}

func (s PortSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
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

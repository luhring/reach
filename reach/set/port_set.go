package set

import (
	"encoding/json"
)

// A PortSet represents a set of network traffic in terms of ports, suitable for describing TCP or UDP traffic. The PortSet type itself does not specify that the described content is for a particular IP protocol (like TCP or UDP).
type PortSet struct {
	set Set
}

// NewEmptyPortSet returns a new, empty PortSet.
func NewEmptyPortSet() PortSet {
	return PortSet{
		set: newEmptySet(),
	}
}

// NewFullPortSet returns a new, full PortSet.
func NewFullPortSet() PortSet {
	return PortSet{
		set: newCompleteSet(),
	}
}

// NewPortSetFromRange returns a new PortSet containing all ports contained in the specified range, inclusively.
func NewPortSetFromRange(lowPort, highPort uint16) PortSet {
	return PortSet{
		set: newSetFromRange(lowPort, highPort),
	}
}

// Complete returns a boolean indicating whether or not the PortSet is complete.
func (s PortSet) Complete() bool {
	return s.set.Complete()
}

// Empty returns a boolean indicating whether or not the PortSet is empty.
func (s PortSet) Empty() bool {
	return s.set.Empty()
}

// Intersect takes the set intersection of two sets of ports and returns the result. Because the PortSet type does not specify whether the content is TCP or UDP, that check must be performed by the consumer.
func (s PortSet) Intersect(other PortSet) PortSet {
	return PortSet{
		set: s.set.intersect(other.set),
	}
}

// Merge takes the set merging of two sets of ports and returns the result. Because the PortSet type does not specify whether the content is TCP or UDP, that check must be performed by the consumer.
func (s PortSet) Merge(other PortSet) PortSet {
	return PortSet{
		set: s.set.merge(other.set),
	}
}

// Subtract takes the "other" set and subtracts it from the calling set and returns the result. Because the PortSet type does not specify whether the content is TCP or UDP, that check must be performed by the consumer.
func (s PortSet) Subtract(other PortSet) PortSet {
	return PortSet{
		set: s.set.subtract(other.set),
	}
}

// RangeStrings returns a slice of strings, where each string describes a continuous range of ports within the PortSet.
func (s PortSet) RangeStrings() []string {
	return s.set.rangeStrings()
}

// String returns the string representation of the PortSet.
func (s PortSet) String() string {
	return s.set.String()
}

// MarshalJSON returns the JSON representation of the PortSet.
func (s PortSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

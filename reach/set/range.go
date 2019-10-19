package set

import (
	"fmt"
	"strconv"
)

// Range defines a continuous, inclusive range of values.
type Range struct {
	first uint16
	last  uint16
}

// String returns the string representation of the Range.
func (r Range) String() string {
	if r.first == r.last {
		return strconv.Itoa(int(r.first))
	}

	return fmt.Sprintf("%d-%d", r.first, r.last)
}

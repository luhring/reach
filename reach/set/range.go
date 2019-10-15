package set

import (
	"fmt"
	"strconv"
)

type Range struct {
	first uint16
	last  uint16
}

func (r Range) String() string {
	if r.first == r.last {
		return strconv.Itoa(int(r.first))
	}

	return fmt.Sprintf("%d-%d", r.first, r.last)
}

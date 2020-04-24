package reach

import "fmt"

type UniversalReference struct {
	Implicit bool // Is infrastructure implied by referenced resource instead of being the resource itself
	R        ResourceReference
}

func (i UniversalReference) Equal(other UniversalReference) bool {
	if i.Implicit != other.Implicit {
		return false
	}

	return i.R.Equal(other.R)
}

func (i UniversalReference) String() string {
	var implicitSuffix string
	if i.Implicit {
		implicitSuffix = " (implicit infrastructure)"
	}

	return fmt.Sprintf("%s%s", i.R, implicitSuffix)
}

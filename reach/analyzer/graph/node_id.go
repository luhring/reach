package graph

import "fmt"

type NodeID uint64

func (id NodeID) String() string {
	return fmt.Sprint(uint64(id))
}

package graph

import (
	"github.com/luhring/reach/reach"
)

type Node struct {
	id   NodeID
	p    *reach.Point
	prev *NodeID // nil if this is the first Node
}

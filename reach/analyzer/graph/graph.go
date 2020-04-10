package graph

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/luhring/reach/reach"
)

type Graph struct {
	m     sync.RWMutex // See notes for sync.Map re: potentially using that instead of this
	r     *rand.Rand
	nodes map[NodeID]*Node
}

func New() Graph {
	nodes := make(map[NodeID]*Node)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return Graph{
		m:     sync.RWMutex{},
		r:     r,
		nodes: nodes,
	}
}

func (g *Graph) Add(p *reach.Point, prev *NodeID) (NodeID, error) {
	id := g.GenerateNodeID()

	if g.validPrev(prev) == false {
		return 0, fmt.Errorf("cannot add node when specified prev node is invalid: %s", *prev)
	}

	g.m.Lock() // begin critical section
	defer g.m.Unlock()

	g.nodes[id] = &Node{
		id:   id,
		p:    p,
		prev: prev,
	}
	return id, nil
}

func (g *Graph) Get(id NodeID) (Node, bool) {
	g.m.RLock()
	defer g.m.RUnlock()

	n := g.nodes[id]
	if n == nil {
		return Node{}, false
	}

	return *n, true
}

func (g *Graph) GenerateNodeID() NodeID {
	i := g.r.Uint64()
	return NodeID(i)
}

func (g *Graph) exists(id NodeID) bool {
	g.m.RLock()
	defer g.m.RUnlock()

	n := g.nodes[id]
	return n != nil
}

func (g *Graph) validPrev(prev *NodeID) bool {
	if prev == nil {
		return true
	}

	return g.exists(*prev)
}

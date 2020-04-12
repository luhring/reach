package reach

import (
	"errors"
	"fmt"
)

// RCProvider is a possibly temporary wrapper for RC, until we figure out the best way to facilitate resource querying
type RCProvider struct {
	rc *ResourceCollection
}

func NewRCProvider(rc *ResourceCollection) *RCProvider {
	return &RCProvider{rc: rc}
}

func (p *RCProvider) Get(ref InfrastructureReference) (Resource, error) {
	if ref.Implicit == false {
		r := p.rc.Get(ref.R)
		if r == nil {
			return Resource{}, fmt.Errorf("not found in resource collection: %+v", ref)
		}

		return *r, nil
	}

	return Resource{}, errors.New("implicit infrastructure not available from RCProvider yet")
}

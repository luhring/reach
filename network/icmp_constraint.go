package network

import (
	"errors"
	"fmt"
)

type ICMPConstraint struct {
	Type int64 // -1 for all types
	Code int64 // -1 for all codes
	V6   bool  // false for ICMP (v4), true for ICMPv6
}

func NewICMPConstraintForAllV4Traffic() *ICMPConstraint {
	return &ICMPConstraint{
		Type: all,
		Code: all,
		V6:   false,
	}
}

func NewICMPConstraintForAllV6Traffic() *ICMPConstraint {
	return &ICMPConstraint{
		Type: all,
		Code: all,
		V6:   true,
	}
}

func (c *ICMPConstraint) GetType() int64 {
	return c.Type
}

func (c *ICMPConstraint) GetCode() int64 {
	return c.Code
}

func (c *ICMPConstraint) isV4() bool {
	return c.V6 == false
}

func (c *ICMPConstraint) isV6() bool {
	return c.V6
}

func (c *ICMPConstraint) AllV4Types() bool {
	return c.isV4() && c.GetType() == all
}

func (c *ICMPConstraint) AllV4Traffic() bool {
	return c.AllV4Types()
}

func (c *ICMPConstraint) AllV4Codes() bool {
	return c.isV4() && (c.AllV4Types() || c.GetCode() == all)
}

func (c *ICMPConstraint) AllV6Types() bool {
	return c.isV6() && c.GetType() == all
}

func (c *ICMPConstraint) AllV6Traffic() bool {
	return c.AllV6Types()
}

func (c *ICMPConstraint) AllV6Codes() bool {
	return c.isV6() && (c.AllV6Codes() || c.GetCode() == all)
}

func (c *ICMPConstraint) isSameVersionAs(other *ICMPConstraint) bool {
	return c.isV4() == other.isV4()
}

func (c *ICMPConstraint) IntersectionWith(other *ICMPConstraint) *ICMPConstraint {
	if other == nil {
		return nil
	}

	if c.isSameVersionAs(other) == false {
		return nil
	}

	// ICMP versions are the same

	if c.isV4() {
		if c.AllV4Traffic() && other.AllV4Traffic() {
			return NewICMPConstraintForAllV4Traffic()
		}

		if c.AllV4Traffic() {
			return other
		}

		if other.AllV4Traffic() {
			return c
		}
	}

	if c.isV6() {
		if c.AllV6Traffic() && other.AllV6Traffic() {
			return NewICMPConstraintForAllV6Traffic()
		}

		if c.AllV6Traffic() {
			return other
		}

		if other.AllV6Traffic() {
			return c
		}
	}

	if c.Type != other.GetType() {
		return nil
	}

	// types are the same

	if c.Code != other.GetCode() {
		return nil
	}

	// codes are the same

	return c
}

func (c *ICMPConstraint) mergeWith(other *ICMPConstraint) (*ICMPConstraint, error) {
	if intersection := c.IntersectionWith(other); intersection == nil {
		return nil, errors.New("unable to merge ICMP constraints due to lack of intersection")
	}

	if c.AllV4Traffic() || other.AllV4Traffic() {
		return NewICMPConstraintForAllV4Traffic(), nil
	}

	if c.AllV6Traffic() || other.AllV6Traffic() {
		return NewICMPConstraintForAllV6Traffic(), nil
	}

	return c, nil
}

func (c *ICMPConstraint) Describe() string {
	if c.isV4() {
		if c.AllV4Traffic() {
			return "all ICMPv4 types and codes"
		}

		typeName := fmt.Sprintf("type %v (%v)", c.Type, c.getV4TypeName())

		var codeName string

		if c.AllV4Codes() {
			codeName = "all codes"
		} else {
			codeName = fmt.Sprintf("code %v", c.Code)
		}

		return fmt.Sprintf("%v, %v", typeName, codeName)
	}

	if c.AllV6Traffic() {
		return "all ICMPv6 types and codes"
	}

	typeName := fmt.Sprintf("type %v (%v)", c.Type, c.getV6TypeName())

	var codeName string

	if c.AllV6Codes() {
		codeName = "all codes"
	} else {
		codeName = fmt.Sprintf("code %v", c.Code)
	}

	return fmt.Sprintf("%v, %v", typeName, codeName)
}

func (c *ICMPConstraint) getV4TypeName() string {
	if typeName := icmpv4TypeNames[c.Type]; typeName != "" {
		return typeName
	}

	return "unknown type"
}

func (c *ICMPConstraint) getV6TypeName() string {
	if typeName := icmpv6TypeNames[c.Type]; typeName != "" {
		return typeName
	}

	return "unknown type"
}

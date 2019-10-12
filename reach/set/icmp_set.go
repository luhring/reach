package set

import (
	"fmt"
	"strings"
)

const (
	AllICMPTypes = -1
	AllICMPCodes = -1

	minimumICMPType = 0
	maximumICMPType = 255
	minimumICMPCode = 0
	maximumICMPCode = 255
)

type ICMPSet struct {
	set Set
}

func NewEmptyICMPSet() ICMPSet {
	return ICMPSet{
		set: newEmptySet(),
	}
}

func NewFullICMPSet() ICMPSet {
	return ICMPSet{
		set: newCompleteSet(),
	}
}

func NewICMPSetFromICMPType(icmpType uint8) (ICMPSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return ICMPSet{}, fmt.Errorf("unable to use icmpType: %v", err)
	}

	startingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), minimumICMPCode)
	endingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), maximumICMPCode)

	set := newSetFromRange(startingICMPTypeCodeIndex, endingICMPTypeCodeIndex)

	return ICMPSet{
		set: set,
	}, nil
}

func NewICMPSetFromICMPTypeCode(icmpType, icmpCode uint8) (ICMPSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return ICMPSet{}, fmt.Errorf("unable to use icmpType: %v", err)
	}

	if err := validateICMPCode(icmpCode); err != nil {
		return ICMPSet{}, fmt.Errorf("unable to use icmpCode: %v", err)
	}

	typeCodeIndex := encodeICMPTypeCode(uint(icmpType), uint(icmpCode))
	set := NewSetForSingleValue(typeCodeIndex)

	return ICMPSet{
		set: set,
	}, nil
}

func (s ICMPSet) Complete() bool {
	return s.set.Complete()
}

func (s ICMPSet) Empty() bool {
	return s.set.Empty()
}

func allTypes(first, last ICMPTypeCode) (bool, string) {
	if first.icmpType == minimumICMPType && first.icmpCode == minimumICMPCode && last.icmpType == maximumICMPType && last.icmpCode == maximumICMPCode {
		return true, "all traffic"
	}

	return false, ""
}

func allCodesForOneTypeV4(first, last ICMPTypeCode) (bool, string) {
	if first.icmpType != last.icmpType {
		return false, ""
	}

	// same type!

	if first.icmpCode == minimumICMPCode && last.icmpCode == maximumICMPCode {
		return true, fmt.Sprintf("ICMPv4 type \"%s\" (all traffic)", GetICMPv4TypeName(first.icmpType))
	}

	return false, ""
}

func allCodesForOneTypeV6(first, last ICMPTypeCode) (bool, string) {
	if first.icmpType != last.icmpType {
		return false, ""
	}

	// same type!

	if first.icmpCode == minimumICMPCode && last.icmpCode == maximumICMPCode {
		return true, fmt.Sprintf("ICMPv6 type \"%s\" (all traffic)", GetICMPv6TypeName(first.icmpType))
	}

	return false, ""
}

func (s ICMPSet) RangeStringsV4() []string {
	var result []string

	for _, rangeItem := range s.set.ranges() {
		firstICMPTypeCode := decodeICMPTypeCode(rangeItem.first)
		lastICMPTypeCode := decodeICMPTypeCode(rangeItem.last)

		if isAllTypes, name := allTypes(firstICMPTypeCode, lastICMPTypeCode); isAllTypes {
			return []string{name}
		}

		if isAllCodesForType, name := allCodesForOneTypeV4(firstICMPTypeCode, lastICMPTypeCode); isAllCodesForType {
			result = append(result, name)
			continue
		}

		rangeString := fmt.Sprintf("%s - %s", firstICMPTypeCode.StringV4(), lastICMPTypeCode.StringV4())
		result = append(result, rangeString)
	}

	return result
}

func (s ICMPSet) RangeStringsV6() []string {
	var result []string

	for _, rangeItem := range s.set.ranges() {
		firstICMPTypeCode := decodeICMPTypeCode(rangeItem.first)
		lastICMPTypeCode := decodeICMPTypeCode(rangeItem.last)

		if isAllTypes, name := allTypes(firstICMPTypeCode, lastICMPTypeCode); isAllTypes {
			return []string{name}
		}

		if isAllCodesForType, name := allCodesForOneTypeV6(firstICMPTypeCode, lastICMPTypeCode); isAllCodesForType {
			result = append(result, name)
			continue
		}

		rangeString := fmt.Sprintf("%s - %s", firstICMPTypeCode.StringV6(), lastICMPTypeCode.StringV6())
		result = append(result, rangeString)
	}

	return result
}

func (s ICMPSet) StringV4() string {
	if s.Empty() {
		return "[empty]"
	}
	return strings.Join(s.RangeStringsV4(), ", ")
}

func (s ICMPSet) StringV6() string {
	if s.Empty() {
		return "[empty]"
	}
	return strings.Join(s.RangeStringsV6(), ", ")
}

func (s ICMPSet) Intersect(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.intersect(other.set),
	}
}

func (s ICMPSet) Merge(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.merge(other.set),
	}
}

func (s ICMPSet) Subtract(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.subtract(other.set),
	}
}

func validateICMPType(icmpType uint8) error {
	if icmpType < minimumICMPType || icmpType > maximumICMPType {
		return fmt.Errorf(
			"icmpType value %v is not valid, must be between %v and %v (inclusive)",
			icmpType,
			minimumICMPType,
			maximumICMPType,
		)
	}

	return nil
}

func validateICMPCode(icmpCode uint8) error {
	if icmpCode < minimumICMPCode || icmpCode > maximumICMPCode {
		return fmt.Errorf(
			"icmpCode value %v is not valid, must be between %v and %v (inclusive)",
			icmpCode,
			minimumICMPCode,
			maximumICMPCode,
		)
	}

	return nil
}

func encodeICMPTypeCode(icmpType, icmpCode uint) uint16 {
	const bitSize = 8

	return uint16((icmpType << bitSize) | icmpCode)
}

func decodeICMPTypeCode(value uint16) ICMPTypeCode {
	const bitSize = 8

	var icmpType uint8 = uint8((value & 0b1111111100000000) >> bitSize)

	var icmpCode uint8 = uint8((value & 0b0000000011111111))

	return ICMPTypeCode{icmpType, icmpCode}
}

package set

import (
	"fmt"
	"strings"
)

// Constants to handle ingestion of data that refers to all ICMP types or codes.
const (
	AllICMPTypes = -1
	AllICMPCodes = -1
)

const (
	minimumICMPType = 0
	maximumICMPType = 255
	minimumICMPCode = 0
	maximumICMPCode = 255
)

// ICMPSet is a set of ICMP traffic, expressed as ICMP types and codes. ICMPSet can be used for either ICMPv4 or ICMPv6. For more information on how ICMPv4 and ICMPv6 use types and codes to describe IP traffic, see RFC 792 and RFC 4443, respectively.
type ICMPSet struct {
	set Set
}

// NewEmptyICMPSet returns a new, empty ICMPSet.
func NewEmptyICMPSet() ICMPSet {
	return ICMPSet{
		set: newEmptySet(),
	}
}

// NewFullICMPSet returns a new, full ICMPSet.
func NewFullICMPSet() ICMPSet {
	return ICMPSet{
		set: newCompleteSet(),
	}
}

// NewICMPSetFromICMPType returns a new ICMPSet containing all codes gor a given type.
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

// NewICMPSetFromICMPTypeCode returns a new ICMPSet containing the given type and code.
func NewICMPSetFromICMPTypeCode(icmpType, icmpCode uint8) (ICMPSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return ICMPSet{}, fmt.Errorf("unable to use icmpType: %v", err)
	}

	if err := validateICMPCode(icmpCode); err != nil {
		return ICMPSet{}, fmt.Errorf("unable to use icmpCode: %v", err)
	}

	typeCodeIndex := encodeICMPTypeCode(uint(icmpType), uint(icmpCode))
	set := newSetForSingleValue(typeCodeIndex)

	return ICMPSet{
		set: set,
	}, nil
}

// Complete returns a boolean indicating whether or not the ICMPSet is complete.
func (s ICMPSet) Complete() bool {
	return s.set.Complete()
}

// Empty returns a boolean indicating whether or not the ICMPSet is empty.
func (s ICMPSet) Empty() bool {
	return s.set.Empty()
}

func allTypesV4(first, last ICMPTypeCode) (bool, string) {
	if first.icmpType == minimumICMPType && first.icmpCode == minimumICMPCode && last.icmpType == maximumICMPType && last.icmpCode == maximumICMPCode {
		return true, "ICMPv4 (all traffic)"
	}

	return false, ""
}

func allTypesV6(first, last ICMPTypeCode) (bool, string) {
	if first.icmpType == minimumICMPType && first.icmpCode == minimumICMPCode && last.icmpType == maximumICMPType && last.icmpCode == maximumICMPCode {
		return true, "ICMPv6 (all traffic)"
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

// RangeStringsV4 returns a slice of strings, where each string describes an individual ICMPv4 type component of the ICMPSet.
func (s ICMPSet) RangeStringsV4() []string {
	var result []string

	for _, rangeItem := range s.set.ranges() {
		firstICMPTypeCode := decodeICMPTypeCode(rangeItem.first)
		lastICMPTypeCode := decodeICMPTypeCode(rangeItem.last)

		if isAllTypes, name := allTypesV4(firstICMPTypeCode, lastICMPTypeCode); isAllTypes {
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

// RangeStringsV6 returns a slice of strings, where each string describes an individual ICMPv6 type component of the ICMPSet.
func (s ICMPSet) RangeStringsV6() []string {
	var result []string

	for _, rangeItem := range s.set.ranges() {
		firstICMPTypeCode := decodeICMPTypeCode(rangeItem.first)
		lastICMPTypeCode := decodeICMPTypeCode(rangeItem.last)

		if isAllTypes, name := allTypesV6(firstICMPTypeCode, lastICMPTypeCode); isAllTypes {
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

// StringV4 returns the string representation of the ICMPSet, assuming that the set describes ICMPv4 content.
func (s ICMPSet) StringV4() string {
	if s.Empty() {
		return "[empty]"
	}
	return strings.Join(s.RangeStringsV4(), ", ")
}

// StringV6 returns the string representation of the ICMPSet, assuming that the set describes ICMPv6 content.
func (s ICMPSet) StringV6() string {
	if s.Empty() {
		return "[empty]"
	}
	return strings.Join(s.RangeStringsV6(), ", ")
}

// Intersect takes the set intersection of two sets of ICMP traffic and returns the result. Because the ICMPSet type does not specify whether the content is ICMPv4 or ICMPv6, that check must be performed by the consumer.
func (s ICMPSet) Intersect(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.intersect(other.set),
	}
}

// Merge takes the set merging of two sets of ICMP traffic and returns the result. Because the ICMPSet type does not specify whether the content is ICMPv4 or ICMPv6, that check must be performed by the consumer.
func (s ICMPSet) Merge(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.merge(other.set),
	}
}

// Subtract takes the input set and subtracts it from the calling set of ICMP traffic and returns the result. Because the ICMPSet type does not specify whether the content is ICMPv4 or ICMPv6, that check must be performed by the consumer.
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

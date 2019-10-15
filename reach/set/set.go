package set

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const chunkSize = 64
const fullChunk = math.MaxUint64
const numberOfChunksInSet = 1024

type Set struct {
	complete bool // if 'complete' is true, 'chunks' should not be accessed
	empty    bool // if 'empty' is true, 'chunks' should not be accessed
	chunks   *[numberOfChunksInSet]uint64
}

func newEmptySet() Set {
	return Set{
		complete: false,
		empty:    true,
		chunks:   nil,
	}
}

func newCompleteSet() Set {
	return Set{
		complete: true,
		empty:    false,
		chunks:   nil,
	}
}

func newSetWithChunks(chunks [1024]uint64) Set {
	return Set{
		complete: false,
		empty:    false,
		chunks:   &chunks,
	}
}

func newSetFromRange(low, high uint16) Set {
	// TODO: if range makes set complete, skip chunk approach in favor of 'complete' bool

	// 0 - 63  --   i=0
	// 64 - 127 --  i=1
	// 128 - 191 -- i=2

	// 11100000 -> items (e.g. ports) 0-2, inclusive

	const chunkSize = 64
	var resultChunks [numberOfChunksInSet]uint64

	// TODO: Fast-forward i to first applicable chunk (instead of skipping to it chunk by chunk)

	resultIsStillEmpty := true
	resultIsStillComplete := true

	for i := 0; i < numberOfChunksInSet; i++ {
		// what are bounds of current chunk?
		chunkStartValue := uint16(chunkSize * i)
		chunkEndValue := uint16(chunkSize*i + (chunkSize - 1))

		if low > chunkEndValue {
			resultIsStillComplete = false
			continue // we don't need to start writing ones yet
		}

		if high < chunkStartValue {
			resultIsStillComplete = false
			break // we won't need to write any ones any more
		}

		if low <= chunkStartValue && high >= chunkEndValue {
			resultChunks[i] = fullChunk
			resultIsStillEmpty = false
			continue
		}

		var startBitPosition, endBitPosition uint8

		if low <= chunkStartValue {
			startBitPosition = 0
		} else {
			startBitPosition = calculateBitPositionWithinChunk(chunkSize, i, int(low))
		}

		if high >= chunkEndValue {
			endBitPosition = chunkSize - 1
		} else {
			endBitPosition = calculateBitPositionWithinChunk(chunkSize, i, int(high))
		}

		var currentChunk uint64

		currentChunk = createUint64ForBitPositionRange(startBitPosition, endBitPosition, chunkSize)

		resultChunks[i] = currentChunk

		if resultIsStillEmpty && !chunkIsEmpty(resultChunks[i]) {
			resultIsStillEmpty = false
		}

		if resultIsStillComplete && !chunkIsFull(resultChunks[i]) {
			resultIsStillComplete = false
		}
	}

	if resultIsStillEmpty {
		return newEmptySet()
	}

	if resultIsStillComplete {
		return newCompleteSet()
	}

	return newSetWithChunks(resultChunks)
}

func NewSetForSingleValue(val uint16) Set {
	return newSetFromRange(val, val)
}

func (s Set) Complete() bool {
	return s.complete
}

func (s Set) Empty() bool {
	return s.empty
}

func (s Set) rangeStrings() []string {
	var result []string

	for _, rangeItem := range s.ranges() {
		result = append(result, rangeItem.String())
	}

	return result
}

func (s Set) ranges() []Range {
	var result []Range
	var curRangeStart uint16
	var chunk uint64
	var err error
	midRange := false

	if s.complete {
		result = append(result, Range{0, chunkSize*numberOfChunksInSet - 1})
		return result
	} else if s.empty {
		return nil
	}

	for chunkIdx := 0; chunkIdx < numberOfChunksInSet; chunkIdx++ {
		chunk = s.chunks[chunkIdx]

		// faster processing by checking the entire chunk
		if chunkIsFull(chunk) {
			if !midRange {
				// edge case: new range on start of chunk; start tracking a range...
				midRange = true
				curRangeStart, err = calculateValueAtPosition(chunkSize, chunkIdx, 0)
				if err != nil {
					panic(err)
				}
			}
			continue
		}

		if chunkIsEmpty(chunk) {
			if midRange {
				// edge case: end of range on start of a new chunk; terminate...
				curRangeEnd, err := calculateValueAtPosition(chunkSize, chunkIdx, 0)
				if err != nil {
					panic(err)
				}
				result = append(result, Range{curRangeStart, curRangeEnd - 1})
				midRange = false
			}
			continue
		}

		// this is a mixed chunk (not full or empty), determine the edges...
		for chunkSubIdx := 0; chunkSubIdx < chunkSize; chunkSubIdx++ {
			isBitSet := (1<<uint64(chunkSize-chunkSubIdx-1))&chunk != 0

			if midRange && !isBitSet {
				// terminate the current tracked range...
				curRangeEnd, err := calculateValueAtPosition(chunkSize, chunkIdx, chunkSubIdx)
				if err != nil {
					panic(err)
				}
				result = append(result, Range{curRangeStart, curRangeEnd - 1})
				midRange = false
			}

			if !midRange && isBitSet {
				// start tracking a range...
				midRange = true
				curRangeStart, err = calculateValueAtPosition(chunkSize, chunkIdx, chunkSubIdx)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// if we were tracking a range up till the last port value, terminate
	if midRange {
		curRangeEnd, err := calculateValueAtPosition(chunkSize, numberOfChunksInSet-1, chunkSize-1)
		if err != nil {
			panic(err)
		}
		result = append(result, Range{curRangeStart, curRangeEnd})
	}

	return result
}

func (s Set) String() string {
	if s.Empty() {
		return "[empty]"
	}
	return strings.Join(s.rangeStrings(), ", ")
}

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.rangeStrings())
}

// equals tests if two sets are equivalent
func (s Set) equals(other Set) bool {
	if s.Complete() && other.Complete() {
		return true
	}

	if s.Complete() != other.Complete() {
		return false
	}

	if s.Empty() && other.Empty() {
		return true
	}

	if s.Empty() != other.Empty() {
		return false
	}

	if s.chunks != nil && other.chunks != nil {
		for i := 0; i < numberOfChunksInSet; i++ {
			if s.chunks[i] != other.chunks[i] {
				return false
			}
		}

		return true
	}

	return false // unexpected situation
}

// intersect set with other set
func (s Set) intersect(other Set) Set {
	if s.Empty() || other.Empty() {
		return newEmptySet()
	}

	if s.Complete() && other.Complete() {
		return newCompleteSet()
	}

	if s.Complete() {
		return other
	}

	if other.Complete() {
		return s
	}

	var resultChunks [numberOfChunksInSet]uint64

	resultIsStillEmpty := true
	resultIsStillComplete := true

	for i := 0; i < numberOfChunksInSet; i++ {
		resultChunks[i] = s.chunks[i] & other.chunks[i]

		if resultIsStillEmpty && !chunkIsEmpty(resultChunks[i]) {
			resultIsStillEmpty = false
		}

		if resultIsStillComplete && !chunkIsFull(resultChunks[i]) {
			resultIsStillComplete = false
		}
	}

	if resultIsStillEmpty {
		return newEmptySet()
	}

	if resultIsStillComplete {
		return newCompleteSet()
	}

	return newSetWithChunks(resultChunks)
}

// merge set with other set
func (s Set) merge(other Set) Set {
	if s.Complete() || other.Complete() {
		return newCompleteSet()
	}

	if s.Empty() && other.Empty() {
		return newEmptySet()
	}

	if s.Empty() {
		return other
	}

	if other.Empty() {
		return s
	}

	var resultChunks [numberOfChunksInSet]uint64

	resultIsStillEmpty := true
	resultIsStillComplete := true

	for i := 0; i < numberOfChunksInSet; i++ {
		resultChunks[i] = s.chunks[i] | other.chunks[i]

		if resultIsStillEmpty && !chunkIsEmpty(resultChunks[i]) {
			resultIsStillEmpty = false
		}

		if resultIsStillComplete && !chunkIsFull(resultChunks[i]) {
			resultIsStillComplete = false
		}
	}

	if resultIsStillEmpty {
		return newEmptySet()
	}

	if resultIsStillComplete {
		return newCompleteSet()
	}

	return newSetWithChunks(resultChunks)
}

// subtract 'other' set from set (= set - other set)
func (s Set) subtract(other Set) Set {
	if s.Empty() || other.Complete() {
		return newEmptySet()
	}

	if other.Empty() {
		return s
	}

	if s.Complete() {
		return other.invert()
	}

	var resultChunks [numberOfChunksInSet]uint64

	resultIsStillEmpty := true
	resultIsStillComplete := true

	for i := 0; i < numberOfChunksInSet; i++ {
		resultChunks[i] = s.chunks[i] &^ other.chunks[i]

		if resultIsStillEmpty && !chunkIsEmpty(resultChunks[i]) {
			resultIsStillEmpty = false
		}

		if resultIsStillComplete && !chunkIsFull(resultChunks[i]) {
			resultIsStillComplete = false
		}
	}

	if resultIsStillEmpty {
		return newEmptySet()
	}

	if resultIsStillComplete {
		return newCompleteSet()
	}

	return newSetWithChunks(resultChunks)
}

// invert the set
func (s Set) invert() Set {
	if s.Empty() {
		return newCompleteSet()
	}

	if s.Complete() {
		return newEmptySet()
	}

	var resultChunks [numberOfChunksInSet]uint64

	for i := 0; i < numberOfChunksInSet; i++ {
		resultChunks[i] = fullChunk ^ s.chunks[i]
	}

	return newSetWithChunks(resultChunks)
}

func createUint64ForBitPositionRange(start, end, chunkSize uint8) uint64 { // range is zero-based
	startShift := chunkSize - start
	endShift := chunkSize - (end + 1)

	// prevent shift overflow

	var block uint64

	if startShift == 64 {
		block = math.MaxUint64
	} else {
		block = 1<<startShift - 1
	}

	var hole uint64

	if endShift == 64 {
		hole = math.MaxUint64
	} else {
		hole = 1<<endShift - 1
	}

	return block ^ hole
}

func calculateValueAtPosition(chunkSize, chunkIndex, chunkSubIndex int) (uint16, error) {
	if chunkSubIndex > chunkSize {
		return 0, fmt.Errorf("chunk index (%d) is greater than chunk size (%d)", chunkSubIndex, chunkSize)
	}
	return uint16((chunkSize * chunkIndex) + chunkSubIndex), nil
}

func calculateBitPositionWithinChunk(chunkSize, chunkIndex, value int) uint8 {
	return uint8(value - (chunkSize * chunkIndex))
}

func chunkIsEmpty(chunk uint64) bool {
	return chunk == 0
}

func chunkIsFull(chunk uint64) bool {
	return chunk == math.MaxUint64
}

package set

import (
	"fmt"
	"math"
	"testing"
)

func TestNewEmptySet(t *testing.T) {
	s := newEmptySet()

	if s.complete != false {
		t.Error("complete should be false")
	}

	if s.empty != true {
		t.Error("empty should be true")
	}

	if s.chunks != nil {
		t.Error("chunks should be nil")
	}
}

func TestNewCompleteSet(t *testing.T) {
	s := newCompleteSet()

	if s.complete != true {
		t.Error("complete should be true")
	}

	if s.empty != false {
		t.Error("empty should be false")
	}

	if s.chunks != nil {
		t.Error("chunks should be nil")
	}
}

func TestNewSetWithChunks(t *testing.T) {
	var c [1024]uint64
	s := newSetWithChunks(c)

	if s.complete != false {
		t.Error("complete should be false")
	}

	if s.empty != false {
		t.Error("empty should be false")
	}

	if s.chunks == nil {
		t.Error("chunks shouldn't be nil")
	}
}

func TestNewSetFromRange(t *testing.T) {
	cases := []struct {
		name string
		low  uint16
		high uint16
	}{
		{
			"single value",
			88,
			88,
		},
		{
			"range within single chunk",
			22,
			30,
		},
		{
			"range spanning two chunks",
			22,
			80,
		},
		{
			"range spanning many chunks",
			22,
			20000,
		},
		{
			"upper bound test",
			1,
			65535,
		},
		{
			"lower bound test",
			0,
			65534,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newSetFromRange(tc.low, tc.high)

			if s.complete != false {
				t.Error("complete should be false")
			}

			if s.empty != false {
				t.Error("empty should be false")
			}

			if s.chunks == nil {
				t.Error("chunks shouldn't be nil")
			}
		})
	}
}

func TestNewSetForSingleValue(t *testing.T) {
	cases := []struct {
		name  string
		value uint16
	}{
		{
			"value in first chunk",
			2,
		},
		{
			"value in middle chunk",
			100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSetForSingleValue(tc.value)

			if s.complete != false {
				t.Error("complete should be false")
			}

			if s.empty != false {
				t.Error("empty should be false")
			}

			if s.chunks == nil {
				t.Error("chunks shouldn't be nil")
			}
		})
	}
}

func TestIsComplete(t *testing.T) {
	cases := []struct {
		setDescription string
		set            Set
		expected       bool
	}{
		{
			setDescription: "new empty set",
			set:            newEmptySet(),
			expected:       false,
		},
		{
			setDescription: "new complete set",
			set:            newCompleteSet(),
			expected:       true,
		},
		{
			setDescription: "new set for single value",
			set:            NewSetForSingleValue(42),
			expected:       false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.setDescription, func(t *testing.T) {
			actual := tc.set.Complete()
			if actual != tc.expected {
				t.Errorf("expectedValue %t but got %t", tc.expected, actual)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		setDescription string
		set            Set
		expected       bool
	}{
		{
			setDescription: "new empty set",
			set:            newEmptySet(),
			expected:       true,
		},
		{
			setDescription: "new complete set",
			set:            newCompleteSet(),
			expected:       false,
		},
		{
			setDescription: "new set for single value",
			set:            NewSetForSingleValue(42),
			expected:       false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.setDescription, func(t *testing.T) {
			actual := tc.set.Empty()
			if actual != tc.expected {
				t.Errorf("expectedValue %t but got %t", tc.expected, actual)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		name     string
		a        Set
		b        Set
		expected bool
	}{
		{
			"two empty sets",
			newEmptySet(),
			newEmptySet(),
			true,
		},
		{
			"two complete sets",
			newCompleteSet(),
			newCompleteSet(),
			true,
		},
		{
			"complete set and partial set",
			newCompleteSet(),
			newSetFromRange(1, 2),
			false,
		},
		{
			"empty set and partial set",
			newEmptySet(),
			newSetFromRange(1, 2),
			false,
		},
		{
			"empty set and complete set",
			newEmptySet(),
			newCompleteSet(),
			false,
		},
		{
			"equivalent partial sets",
			newSetFromRange(100, 200),
			newSetFromRange(100, 200),
			true,
		},
		{
			"non-equivalent partial sets",
			newSetFromRange(100, 200),
			newSetFromRange(50, 200),
			false,
		},
		{
			"equivalent single-value sets",
			NewSetForSingleValue(5555),
			NewSetForSingleValue(5555),
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.a.equals(tc.b)
			if actual != tc.expected {
				t.Errorf("expectedValue %t but got %t", tc.expected, actual)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	cases := []struct {
		name     string
		a        Set
		b        Set
		expected Set
	}{
		{
			"mutually exclusive sets",
			newSetFromRange(2, 1024),
			newSetFromRange(10000, 12000),
			newEmptySet(),
		},
		{
			"mutually exclusive but neighboring sets",
			newSetFromRange(2, 1024),
			newSetFromRange(1025, 2000),
			newEmptySet(),
		},
		{
			"complete set with partial set",
			newCompleteSet(),
			newSetFromRange(1025, 2000),
			newSetFromRange(1025, 2000),
		},
		{
			"partial set with complete set",
			newSetFromRange(1025, 2000),
			newCompleteSet(),
			newSetFromRange(1025, 2000),
		},
		{
			"complete set with empty set",
			newCompleteSet(),
			newEmptySet(),
			newEmptySet(),
		},
		{
			"two overlapping partial sets (large)",
			newSetFromRange(100, 1000),
			newSetFromRange(500, 50000),
			newSetFromRange(500, 1000),
		},
		{
			"two overlapping partial sets (small)",
			newSetFromRange(80, 90),
			newSetFromRange(85, 95),
			newSetFromRange(85, 90),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.a.intersect(tc.b)
			if actual.equals(tc.expected) == false {
				t.Fail()
			}
		})
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name     string
		a        Set
		b        Set
		expected Set
	}{
		{
			"mutually exclusive but neighboring sets",
			newSetFromRange(2, 1024),
			newSetFromRange(1025, 2000),
			newSetFromRange(2, 2000),
		},
		{
			"complete set with partial set",
			newCompleteSet(),
			newSetFromRange(1025, 2000),
			newCompleteSet(),
		},
		{
			"partial set with complete set",
			newSetFromRange(1025, 2000),
			newCompleteSet(),
			newCompleteSet(),
		},
		{
			"complete set with empty set",
			newCompleteSet(),
			newEmptySet(),
			newCompleteSet(),
		},
		{
			"two overlapping partial sets (large)",
			newSetFromRange(100, 1000),
			newSetFromRange(500, 50000),
			newSetFromRange(100, 50000),
		},
		{
			"two overlapping partial sets (small)",
			newSetFromRange(80, 90),
			newSetFromRange(85, 95),
			newSetFromRange(80, 95),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.a.merge(tc.b)
			if actual.equals(tc.expected) == false {
				t.Fail()
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		name     string
		a        Set
		b        Set
		expected Set
	}{
		{
			"mutually exclusive but neighboring sets",
			newSetFromRange(2, 1024),
			newSetFromRange(1025, 2000),
			newSetFromRange(2, 1024),
		},
		{
			"complete set with partial set",
			newCompleteSet(),
			newSetFromRange(0, 2000),
			newSetFromRange(2001, 65535),
		},
		{
			"partial set with complete set",
			newSetFromRange(1025, 2000),
			newCompleteSet(),
			newEmptySet(),
		},
		{
			"complete set with empty set",
			newCompleteSet(),
			newEmptySet(),
			newCompleteSet(),
		},
		{
			"two overlapping partial sets (large)",
			newSetFromRange(100, 1000),
			newSetFromRange(500, 50000),
			newSetFromRange(100, 499),
		},
		{
			"two overlapping partial sets (small)",
			newSetFromRange(80, 95),
			newSetFromRange(85, 95),
			newSetFromRange(80, 84),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.a.subtract(tc.b)
			if actual.equals(tc.expected) == false {
				t.Fail()
			}
		})
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		name     string
		s        Set
		expected string
	}{
		{
			"lower bound",
			newSetFromRange(0, 100),
			"0-100",
		},
		{
			"upper bound",
			newSetFromRange(65000, 65535),
			"65000-65535",
		},
		{
			"cross-chunk bound",
			newSetFromRange(50, 1000),
			"50-1000",
		},
		{
			"terminate range across chunk",
			newSetFromRange(50, 63),
			"50-63",
		},
		{
			"start range across chunk",
			newSetFromRange(64, 70),
			"64-70",
		},
		{
			"full set",
			newCompleteSet(),
			"0-65535",
		},
		{
			"empty set",
			newEmptySet(),
			"<empty>",
		},
		{
			"mixed ranges",
			newSetFromRange(50, 1000).merge(newSetFromRange(1100, 1200)),
			"50-1000, 1100-1200",
		},
		{
			"single value (middle)",
			NewSetForSingleValue(12),
			"12",
		},
		{
			"upper bound (single value)",
			NewSetForSingleValue(65535),
			"65535",
		},
		{
			"lower bound (single value)",
			NewSetForSingleValue(0),
			"0",
		},
		{
			"upper + lower bound (single values)",
			NewSetForSingleValue(0).merge(NewSetForSingleValue(65535)),
			"0, 65535",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.s.String()
			if actual != tc.expected {
				t.Errorf("expected '%s' but got '%s'", tc.expected, actual)
			}
		})

	}
}

func TestInvert(t *testing.T) {
	cases := []struct {
		name     string
		set      Set
		expected Set
	}{
		{
			"complete set",
			newCompleteSet(),
			newEmptySet(),
		},
		{
			"empty set",
			newEmptySet(),
			newCompleteSet(),
		},
		{
			"low set",
			newSetFromRange(0, 1000),
			newSetFromRange(1001, 65535),
		},
		{
			"high set",
			newSetFromRange(1000, 65535),
			newSetFromRange(0, 999),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.set.invert()

			if actual.equals(tc.expected) == false {
				t.Fail()
			}
		})
	}
}

func TestCreateUint64ForBitPositionRange(t *testing.T) {
	cases := []struct {
		start     uint8
		end       uint8
		chunkSize uint8
		expected  uint64
	}{
		{
			1,
			2,
			4,
			6,
		},
		{
			0,
			0,
			4,
			8,
		},
		{
			2,
			3,
			4,
			3,
		},
		{
			1,
			3,
			4,
			7,
		},
		{
			0,
			3,
			4,
			15,
		},
		{
			1,
			1,
			4,
			4,
		},
		{
			3,
			3,
			4,
			1,
		},
		{
			2,
			6,
			8,
			62,
		},
		{
			0,
			63,
			64,
			math.MaxUint64,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("range %d - %d (of %d)", tc.start, tc.end, tc.chunkSize), func(t *testing.T) {
			actual := createUint64ForBitPositionRange(tc.start, tc.end, tc.chunkSize)

			if actual != tc.expected {
				t.Errorf("expectedValue %b (%d) but got %b (%d)", tc.expected, tc.expected, actual, actual)
			}
		})
	}
}

func TestCalculateBitPositionWithinChunk(t *testing.T) {
	cases := []struct {
		chunkSize  int
		chunkIndex int
		value      int
		expected   uint8
	}{
		{
			4,
			0,
			1,
			1,
		},
		{
			8,
			1,
			12,
			4,
		},
		{
			2,
			2,
			4,
			0,
		},
		{
			64,
			1,
			66,
			2,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case #%v", i), func(t *testing.T) {
			actual := calculateBitPositionWithinChunk(tc.chunkSize, tc.chunkIndex, tc.value)
			if actual != tc.expected {
				t.Errorf("expectedValue %v but got %v", tc.expected, actual)
			}
		})
	}
}

func TestChunkIsEmpty(t *testing.T) {
	cases := []struct {
		chunk    uint64
		expected bool
	}{
		{
			0,
			true,
		},
		{
			1,
			false,
		},
		{
			1000,
			false,
		},
		{
			math.MaxUint64,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("chunk is %v", tc.chunk), func(t *testing.T) {
			actual := chunkIsEmpty(tc.chunk)
			if actual != tc.expected {
				t.Errorf("expectedValue: %t", tc.expected)
			}
		})
	}
}

func TestChunkIsFull(t *testing.T) {
	cases := []struct {
		chunk    uint64
		expected bool
	}{
		{
			0,
			false,
		},
		{
			1,
			false,
		},
		{
			1000,
			false,
		},
		{
			math.MaxUint64,
			true,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("chunk is %v", tc.chunk), func(t *testing.T) {
			actual := chunkIsFull(tc.chunk)
			if actual != tc.expected {
				t.Errorf("expectedValue: %t", tc.expected)
			}
		})
	}
}

func TestCalculateValueAtPosition(t *testing.T) {
	cases := []struct {
		chunkSize     int
		chunkIndex    int
		chunkSubIndex int
		expectedValue int
		expectedError bool
	}{
		{
			4,
			0,
			1,
			1,
			false,
		},
		{
			8,
			1,
			7,
			15,
			false,
		},
		{
			2,
			2,
			4,
			-1,
			true,
		},
		{
			64,
			1,
			66,
			-1,
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case #%v", i), func(t *testing.T) {
			actual, err := calculateValueAtPosition(tc.chunkSize, tc.chunkIndex, tc.chunkSubIndex)
			if !tc.expectedError && err != nil {
				t.Fatalf("expectedValue no error but got %v", err)
			}
			if tc.expectedError && err == nil {
				t.Fatalf("expectedValue no error but got none")
			}
			if actual != tc.expectedValue {
				t.Fatalf("expectedValue %v but got %v", tc.expectedValue, actual)
			}
		})
	}
}

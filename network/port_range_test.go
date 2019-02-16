package network

import (
	"fmt"
	"testing"
)

func TestDoPortRangesIntersect(t *testing.T) {
	testCases := []struct {
		// starting with this, make all the unit tests up to par, and adjusted for expanded Port range logic (with Protocols, etc.)
		firstPortRange  *PortRange
		secondPortRange *PortRange
		expectedOutput  bool
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 200,
				HighPort:                500,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 300,
				HighPort:                400,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 100,
				HighPort:                200,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 300,
				HighPort:                400,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 100,
				HighPort:                300,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 300,
				HighPort:                400,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 500,
				HighPort:                600,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 100,
				HighPort:                200,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 500,
				HighPort:                600,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 100,
				HighPort:                500,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{ ///
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 15,
				HighPort:                25,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 15,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 100,
				HighPort:                500,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		result := doPortRangesIntersect(testCase.firstPortRange, testCase.secondPortRange)

		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected intersect test result for %+v and %+v to be %t, but it was %t.",
				testCase.firstPortRange,
				testCase.secondPortRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestGetHigherOfTwoNumbers(t *testing.T) {
	const lowNumber = 2
	const highNumber = 17

	if getHigherOfTwoNumbers(lowNumber, highNumber) != highNumber {
		t.Fail()
	}
}

func TestGetLowerOfTwoNumbers(t *testing.T) {
	const lowNumber = 2
	const highNumber = 17

	if getLowerOfTwoNumbers(lowNumber, highNumber) != lowNumber {
		t.Fail()
	}
}

func TestPortRangeDoesDescribeOnlyASinglePort(t *testing.T) {
	testCases := []struct {
		portRange      *PortRange
		expectedOutput bool
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1000,
				HighPort:                2000,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 500,
				HighPort:                500,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := testCase.portRange.doesDescribeOnlyASinglePort()

		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected result of doesPortRangeDescribeOnlyASinglePort for %+v to be %t, but it was %t.",
				testCase.portRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestIsValidPortNumber(t *testing.T) {
	invalidPortNumbers := []int64{
		0,
		-200,
		7654321,
	}

	for _, invalidPortNumber := range invalidPortNumbers {
		if isValidPortNumber(invalidPortNumber) {
			t.Errorf(
				"Expected %d not to be evaluated as a valid port number.",
				invalidPortNumber,
			)
		}
	}

	validPortNumbers := []int64{
		1,
		443,
		60000,
	}

	for _, validPortNumber := range validPortNumbers {
		if false == isValidPortNumber(validPortNumber) {
			t.Errorf(
				"Expected %d to be evaluated as a valid port number.",
				validPortNumber,
			)
		}
	}
}

func TestDefragmentPortRanges(t *testing.T) {
	testCases := []struct {
		input          []*PortRange
		expectedOutput []*PortRange
	}{
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 30,
					HighPort:                40,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                40,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 15,
					HighPort:                25,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 30,
					HighPort:                40,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                25,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 30,
					HighPort:                40,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
			},
		},
	}

	for _, testCase := range testCases {
		outputPortRanges := DefragmentPortRanges(testCase.input)
		if false == arePortRangesSlicesEqual(outputPortRanges, testCase.expectedOutput) {
			t.Errorf(
				"Expected defragmentation to result in %+v, but it resulted in %+v.",
				testCase.expectedOutput,
				outputPortRanges,
			)
		}
	}
}

func TestArePortRangesSlicesEqual(t *testing.T) {
	testCases := []struct {
		firstInput     []*PortRange
		secondInput    []*PortRange
		expectedOutput bool
	}{
		{
			firstInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 40,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			secondInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 40,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			expectedOutput: true,
		},
		{
			firstInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 40,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			secondInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 3,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			expectedOutput: false,
		},
		{
			firstInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 40,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			secondInput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 40,
					HighPort:                50,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		output := arePortRangesSlicesEqual(testCase.firstInput, testCase.secondInput)
		if output != testCase.expectedOutput {
			t.Errorf(
				"Expected equal test would result in %t, but it resulted in %t.",
				testCase.expectedOutput,
				output,
			)
		}
	}
}

func TestPortRangeIsValid(t *testing.T) {
	testCases := []struct {
		portRange      *PortRange
		expectedOutput bool
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 -2,
				HighPort:                2,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 0,
				HighPort:                5,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                4,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1,
				HighPort:                80000,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			false,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			true,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		result := testCase.portRange.isValid()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected port range validity test result for %+v to be %t, but it was %t.",
				testCase.portRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestPortRangeGetIntersection(t *testing.T) {
	testCases := []struct {
		firstPortRange  *PortRange
		secondPortRange *PortRange
		expectedOutput  *PortRange
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 15,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 15,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 16,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			nil,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                16,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			nil,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                15,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			nil,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		result := testCase.firstPortRange.getIntersection(testCase.secondPortRange)

		if result == nil && testCase.expectedOutput == nil {
			continue
		}

		if result == nil && testCase.expectedOutput != nil {
			t.Errorf(
				"Expected non-nil intersection of %+v and %+v, but result was nil.",
				testCase.firstPortRange,
				testCase.secondPortRange,
			)
			continue
		}

		if result != nil && testCase.expectedOutput == nil {
			t.Errorf(
				"Expected nil intersection of %+v and %+v, but result was %+v.",
				testCase.firstPortRange,
				testCase.secondPortRange,
				result,
			)
			continue
		}

		if *result != *testCase.expectedOutput {
			t.Errorf(
				"Expected intersection result of %+v and %+v to be %+v, but it was %+v.",
				testCase.firstPortRange,
				testCase.secondPortRange,
				testCase.expectedOutput,
				result,
			)
			continue
		}
	}
}

func TestMergePortRanges(t *testing.T) {
	testCases := []struct {
		firstPortRange      *PortRange
		secondPortRange     *PortRange
		expectedMergeResult *PortRange
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                30,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 20,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                30,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 31,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 31,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 31,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			nil,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1,
				HighPort:                20,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 15,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			nil,
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1,
				HighPort:                100,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 40,
				HighPort:                40,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 1,
				HighPort:                100,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
		},
	}

	for _, testCase := range testCases {
		result := mergePortRanges(testCase.firstPortRange, testCase.secondPortRange)
		if result == nil && testCase.expectedMergeResult == nil {
			continue
		}

		if result == nil && testCase.expectedMergeResult != nil {
			t.Errorf(
				"Expected non-nil merge result of %v and %v, but result was nil.",
				testCase.firstPortRange,
				testCase.secondPortRange,
			)
			continue
		}

		if result != nil && testCase.expectedMergeResult == nil {
			t.Errorf(
				"Expected nil merge result of %v and %v, but result was %v.",
				testCase.firstPortRange,
				testCase.secondPortRange,
				result,
			)
			continue
		}

		if *result != *testCase.expectedMergeResult {
			t.Errorf(
				"Expected merge result of %v and %v to be %v, but it was %v.",
				testCase.firstPortRange,
				testCase.secondPortRange,
				testCase.expectedMergeResult,
				result,
			)
			continue
		}
	}
}

func TestArePortRangesJuxtaposed(t *testing.T) {
	testCases := []struct {
		portRanges     [2]*PortRange
		expectedOutput bool
	}{
		{
			[2]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 21,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			true,
		},
		{
			[2]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 22,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			false,
		},
		{
			[2]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 21,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			false,
		},
		{
			[2]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 10,
					HighPort:                20,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := arePortRangesJuxtaposed(testCase.portRanges)
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected juxtaposition test result to be %t, but it was %t.",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestSortPortRanges(t *testing.T) {
	testCases := []struct {
		portRanges     []*PortRange
		expectedOutput []*PortRange
	}{
		{
			portRanges: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 19,
					HighPort:                34,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			expectedOutput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 19,
					HighPort:                34,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			portRanges: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 19,
					HighPort:                34,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			expectedOutput: []*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 19,
					HighPort:                34,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 20,
					HighPort:                30,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
	}

	for _, testCase := range testCases {
		sortPortRanges(testCase.portRanges)
		if false == arePortRangesSlicesEqual(testCase.portRanges, testCase.expectedOutput) {
			t.Error("Port ranges were not sorted correctly.")
		}
	}
}

func TestGetIntersectionBetweenTwoListsOfPortRanges(t *testing.T) {
	testCases := []struct {
		firstListOfPortRanges  []*PortRange
		secondListOfPortRanges []*PortRange
		expectedOutput         []*PortRange
	}{
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 443,
					HighPort:                443,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 443,
					HighPort:                443,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 443,
					HighPort:                443,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 443,
					HighPort:                443,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                80,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			nil,
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                80,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 53,
					HighPort:                53,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 53,
					HighPort:                53,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                80,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                100,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 90,
					HighPort:                110,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 53,
					HighPort:                53,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                110,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                100,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 90,
					HighPort:                110,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 50,
					HighPort:                60,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
			},
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 50,
					HighPort:                60,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 80,
					HighPort:                110,
					DoesSpecifyAllProtocols: false,
					Protocol:                "tcp",
				},
			},
		},
	}

	for _, testCase := range testCases {
		result := IntersectPortRangeSlices(
			testCase.firstListOfPortRanges,
			testCase.secondListOfPortRanges,
		)

		if false == arePortRangesSlicesEqual(result, testCase.expectedOutput) {
			expandedExpectedOutput := "\n"
			for _, portRange := range testCase.expectedOutput {
				expandedExpectedOutput += fmt.Sprintf(
					"%+v\n",
					portRange,
				)
			}
			expandedResult := "\n"
			for _, portRange := range result {
				expandedResult += fmt.Sprintf(
					"%+v\n",
					portRange,
				)
			}

			t.Errorf(
				"Expected intersection to result in %v, but it resulted in %v.",
				expandedExpectedOutput,
				expandedResult,
			)
		}
	}
}

func TestDescribeListOfPortRanges(t *testing.T) {
	testCases := []struct {
		listOfPortRanges []*PortRange
		expectedOutput   string
	}{
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 443,
					HighPort:                460,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
			},
			"(ALL protocols) 443 - 460\n",
		},
		{
			[]*PortRange{
				{
					DoesSpecifyAllPorts:     true,
					LowPort:                 0,
					HighPort:                0,
					DoesSpecifyAllProtocols: true,
					Protocol:                "",
				},
				{
					DoesSpecifyAllPorts:     false,
					LowPort:                 53,
					HighPort:                53,
					DoesSpecifyAllProtocols: false,
					Protocol:                "udp",
				},
			},
			"(ALL protocols) ALL ports\nUDP 53\n",
		},
	}

	for _, testCase := range testCases {
		result := DescribeListOfPortRanges(testCase.listOfPortRanges)
		if result != testCase.expectedOutput {
			t.Error("Description of list of port ranges resulted in a different string than expected.")
		}
	}
}

func TestPortRangeDescribe(t *testing.T) {
	testCases := []struct {
		portRange      *PortRange
		expectedOutput string
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			"TCP 443",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			"TCP ALL ports",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 80,
				HighPort:                85,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			"UDP 80 - 85",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			"(ALL protocols) ALL ports",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 40,
				HighPort:                41,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			"(ALL protocols) 40 - 41",
		},
	}

	for _, testCase := range testCases {
		result := testCase.portRange.describe()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected portRange.describe() for %+v to result in \"%s\", but it resulted in \"%s\".",
				testCase.portRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestPortRangeDescribeProtocol(t *testing.T) {
	testCases := []struct {
		portRange      *PortRange
		expectedOutput string
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			"(ALL protocols)",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                4443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			"TCP",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                4443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			"UDP",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                4443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "UDP",
			},
			"UDP",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "icmp",
			},
			"ICMP",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "icmpv6",
			},
			"ICMPv6",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "58",
			},
			"ICMPv6",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 80,
				HighPort:                80,
				DoesSpecifyAllProtocols: false,
				Protocol:                "59",
			},
			"(IP protocol 59)",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "1000",
			},
			"(IP protocol 1000)",
		},
	}

	for _, testCase := range testCases {
		result := testCase.portRange.describeProtocol()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected portRange.describeProtocol() for %+v to result in \"%s\", but it resulted in \"%s\".",
				testCase.portRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestPortRangeDescribePorts(t *testing.T) {
	testCases := []struct {
		portRange      *PortRange
		expectedOutput string
	}{
		{
			&PortRange{
				DoesSpecifyAllPorts:     true,
				LowPort:                 0,
				HighPort:                0,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			"ALL ports",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 10,
				HighPort:                2000,
				DoesSpecifyAllProtocols: false,
				Protocol:                "udp",
			},
			"10 - 2000",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 80,
				HighPort:                81,
				DoesSpecifyAllProtocols: true,
				Protocol:                "",
			},
			"80 - 81",
		},
		{
			&PortRange{
				DoesSpecifyAllPorts:     false,
				LowPort:                 443,
				HighPort:                443,
				DoesSpecifyAllProtocols: false,
				Protocol:                "tcp",
			},
			"443",
		},
	}

	for _, testCase := range testCases {
		result := testCase.portRange.describePorts()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected portRange.describePorts() for %+v to result in \"%s\", but it resulted in \"%s\".",
				testCase.portRange,
				testCase.expectedOutput,
				result,
			)
		}
	}
}

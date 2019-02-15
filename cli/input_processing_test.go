package cli

import "testing"

func TestDoesFirstItemMatchBeginningSubstringOfSecondItem(t *testing.T) {
	testCases := []struct {
		firstItem      string
		secondItem     string
		expectedOutput bool
	}{
		{
			"abc",
			"abcdef",
			true,
		},
		{
			"abcdef",
			"abc",
			false,
		},
		{
			"abc",
			"zabcdef",
			false,
		},
	}

	for _, testCase := range testCases {
		result := doesFirstItemMatchBeginningSubstringOfSecondItem(testCase.firstItem, testCase.secondItem)
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected %t but got %t for items: \"%s\" and \"%s\"",
				testCase.expectedOutput,
				result,
				testCase.firstItem,
				testCase.secondItem,
			)
		}
	}
}

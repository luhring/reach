package aws

import (
	"testing"
)

func TestInstanceIsRunning(t *testing.T) {
	testCases := []struct {
		instance       *Instance
		expectedOutput bool
	}{
		{
			&Instance{
				ID:      "i-12345",
				NameTag: "MyInstance",
				State:   "running",
			},
			true,
		},
		{
			&Instance{
				ID:      "i-12345",
				NameTag: "MyInstance",
				State:   "stopped",
			},
			false,
		},
		{
			&Instance{
				ID:      "i-12345",
				NameTag: "MyInstance",
				State:   "shutting-down",
			},
			false,
		},
		{
			&Instance{
				ID:      "i-12345",
				NameTag: "MyInstance",
				State:   "pending",
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := testCase.instance.isRunning()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected output to be %t, but it was %t. For instance:\n%+v",
				testCase.expectedOutput,
				result,
				testCase.instance,
			)
		}
	}
}

func TestInstanceGetFriendlyName(t *testing.T) {
	testCases := []struct {
		instance       *Instance
		expectedOutput string
	}{
		{
			&Instance{
				ID:      "i-12345",
				NameTag: "MyInstance",
			},
			"MyInstance",
		},
		{
			&Instance{
				ID: "i-12345",
			},
			"i-12345",
		},
		{
			&Instance{},
			"[unnamed instance]",
		},
	}

	for _, testCase := range testCases {
		result := testCase.instance.GetFriendlyName()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected GetFriendlyName() to return \"%s\", but it returned \"%s\". For instance:\n%+v",
				testCase.expectedOutput,
				result,
				testCase.instance,
			)
		}
	}
}

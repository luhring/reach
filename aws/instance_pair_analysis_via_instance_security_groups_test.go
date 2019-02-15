package aws

import "testing"

func TestDoesProtocolImplyAllPortsAreAccessible(t *testing.T) {
	testCases := []struct {
		protocol       string
		expectedOutput bool
	}{
		{
			"tcp",
			false,
		},
		{
			"udp",
			false,
		},
		{
			"icmp",
			false,
		},
		{
			"icmpv6",
			false,
		},
		{
			"58",
			false,
		},
		{
			"59",
			true,
		},
		{
			"msg",
			true,
		},
	}

	for _, testCase := range testCases {
		result := doesProtocolImplyAllPortsAreAccessible(testCase.protocol)
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected %t but got %t for protocol \"%s\".",
				testCase.expectedOutput,
				result,
				testCase.protocol,
			)
		}
	}
}

package traffic

import (
	"testing"

	"github.com/luhring/reach/reach/set"
)

func TestIntersect(t *testing.T) {
	cases := []struct {
		name     string
		a        ProtocolContent
		b        ProtocolContent
		expected ProtocolContent
	}{
		{
			name:     "all and all",
			a:        newProtocolContentWithPortsFull(ProtocolUDP),
			b:        newProtocolContentWithPortsFull(ProtocolUDP),
			expected: newProtocolContentWithPortsFull(ProtocolUDP),
		},
		{
			name:     "none and none",
			a:        newProtocolContentWithPortsEmpty(ProtocolUDP),
			b:        newProtocolContentWithPortsEmpty(ProtocolUDP),
			expected: newProtocolContentWithPortsEmpty(ProtocolUDP),
		},
		{
			name:     "different protocols",
			a:        newProtocolContentWithPortsFull(ProtocolUDP),
			b:        newProtocolContentWithPortsFull(ProtocolTCP),
			expected: newProtocolContentWithPortsEmpty(ProtocolUDP),
		},
		{
			name:     "all and none",
			a:        newProtocolContentWithPortsFull(ProtocolUDP),
			b:        newProtocolContentWithPortsEmpty(ProtocolUDP),
			expected: newProtocolContentWithPortsEmpty(ProtocolUDP),
		},
		{
			name:     "single port and all",
			a:        newProtocolContentWithPorts(ProtocolUDP, set.NewPortSetFromRange(53, 53)),
			b:        newProtocolContentWithPortsFull(ProtocolUDP),
			expected: newProtocolContentWithPorts(ProtocolUDP, set.NewPortSetFromRange(53, 53)),
		},
		{
			name:     "all and single port",
			a:        newProtocolContentWithPortsFull(ProtocolUDP),
			b:        newProtocolContentWithPorts(ProtocolUDP, set.NewPortSetFromRange(53, 53)),
			expected: newProtocolContentWithPorts(ProtocolUDP, set.NewPortSetFromRange(53, 53)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outcome := tc.a.intersect(tc.b)
			if outcome.String() != tc.expected.String() {
				t.Errorf("outcome was not expected (outcome: %s, expected: %s", outcome, tc.expected)
			}
		})
	}
}

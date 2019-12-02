package generic

import (
	"fmt"
	"testing"
)

func TestCheckIPAddress(t *testing.T) {
	cases := []struct {
		address string
		valid   bool
	}{
		{
			"1.2.3.4",
			true,
		},
		{
			"2001:0db8:0000:0000:0000:8a2e:0370:7334",
			true,
		},
		{
			"2001:db8::8a2e:370:7334",
			true,
		},
		{
			"example.com",
			false,
		},
		{
			"one.two.three.four",
			false,
		},

		{
			"",
			false,
		},
		{
			".",
			false,
		},
		{
			"....",
			false,
		},
		{
			"100.200.300.400",
			false,
		},
		{
			"10.20.40",
			false,
		},
		{
			"1.2.3.4/32",
			false,
		},
		{
			"i-abcdef123456",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s, valid = %t", tc.address, tc.valid), func(t *testing.T) {
			err := CheckIPAddress(tc.address)
			resultValid := err == nil

			if tc.valid != resultValid {
				t.Fail()
			}
		})
	}
}

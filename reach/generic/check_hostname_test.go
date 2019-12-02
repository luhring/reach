package generic

import (
	"fmt"
	"testing"
)

func TestCheckHostname(t *testing.T) {
	cases := []struct {
		hostname string
		valid    bool
	}{
		{
			"example.com",
			true,
		},
		{
			"sub.example.org",
			true,
		},
		{
			"another.sub.example.co.uk",
			true,
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
			"...",
			false,
		},
		{
			".com",
			false,
		},
		{
			"example..com",
			false,
		},
		{
			"example",
			false,
		},
		{
			"i-abcdef123456",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s, valid = %t", tc.hostname, tc.valid), func(t *testing.T) {
			err := CheckHostname(tc.hostname)
			resultValid := err == nil

			if tc.valid != resultValid {
				t.Fail()
			}
		})
	}
}

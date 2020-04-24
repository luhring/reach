package aws

import (
	"reflect"
	"testing"

	"github.com/luhring/reach/reach"
	acc "github.com/luhring/reach/reach/acceptance"
)

func TestFindEC2InstanceID(t *testing.T) {
	acc.Check(t)

	cases := []struct {
		searchText    string
		expectedID    string
		expectedError error
	}{
		{
			searchText:    "abc",
			expectedError: nil,
			expectedID:    "i-0a93117c7575b6d54",
		},
		{
			searchText:    "def",
			expectedError: nil,
			expectedID:    "i-0136d3233f0ef1924",
		},
		// { // TODO: Add back negative cases when aws_manager implementation is replaced
		// 	searchText:    "ghi",
		// 	expectedError: nil,
		// 	expectedID:    "",
		// },
	}

	for _, tc := range cases {
		t.Run(tc.searchText, func(t *testing.T) {
			id, err := findEC2InstanceID(tc.searchText, nil)

			if tc.expectedID != id {
				reach.DiffErrorf(t, "id", tc.expectedID, id)
			}

			if !reflect.DeepEqual(tc.expectedError, err) {
				reach.DiffErrorf(t, "err", tc.expectedError, err)
			}
		})
	}
}

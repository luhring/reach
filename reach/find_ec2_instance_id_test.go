package reach

import (
	"reflect"
	"testing"
)

func TestFindEC2InstanceID(t *testing.T) {
	checkAcceptance(t)

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
			id, err := FindEC2InstanceID(tc.searchText)

			if tc.expectedID != id {
				diffErrorf(t, "id", tc.expectedID, id)
			}

			if !reflect.DeepEqual(tc.expectedError, err) {
				diffErrorf(t, "err", tc.expectedError, err)
			}
		})
	}
}

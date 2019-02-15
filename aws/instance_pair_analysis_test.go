package aws

import "testing"

func TestInstancePairGenerateMessageForWhenAccessExists(t *testing.T) {
	testCases := []struct {
		instancePair   *InstancePair
		expectedOutput string
	}{
		{
			&InstancePair{
				&Instance{
					ID:      "i-12345",
					NameTag: "MyInstanceOne",
				},
				&Instance{
					ID:      "i-23456",
					NameTag: "MyInstanceTwo",
				},
			},
			"Instance 'MyInstanceOne' is able to access instance 'MyInstanceTwo'.\n",
		},
		{
			&InstancePair{
				&Instance{
					ID:      "i-12345",
					NameTag: "MyInstanceOne",
				},
				&Instance{
					ID:      "i-23456",
					NameTag: "",
				},
			},
			"Instance 'MyInstanceOne' is able to access instance 'i-23456'.\n",
		},
	}

	for _, testCase := range testCases {
		result := testCase.instancePair.generateMessageForWhenAccessExists()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected \"%s\" but got \"%s\".",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestGenerateMessageForWhenAccessDoesNotExist(t *testing.T) {
	testCases := []struct {
		instancePair   *InstancePair
		expectedOutput string
	}{
		{
			&InstancePair{
				&Instance{
					ID:      "i-12345",
					NameTag: "MyInstanceOne",
				},
				&Instance{
					ID:      "i-23456",
					NameTag: "MyInstanceTwo",
				},
			},
			"Instance 'MyInstanceOne' is unable to access instance 'MyInstanceTwo'.\n",
		},
		{
			&InstancePair{
				&Instance{
					ID:      "i-12345",
					NameTag: "MyInstanceOne",
				},
				&Instance{
					ID:      "i-23456",
					NameTag: "",
				},
			},
			"Instance 'MyInstanceOne' is unable to access instance 'i-23456'.\n",
		},
	}

	for _, testCase := range testCases {
		result := testCase.instancePair.generateMessageForWhenAccessDoesNotExist()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected \"%s\" but got \"%s\".",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestInstancePairAreInstancesInSameSubnet(t *testing.T) {
	testCases := []struct {
		instancePair   *InstancePair
		expectedOutput bool
	}{
		{
			&InstancePair{
				&Instance{
					SubnetID: "12345",
				},
				&Instance{
					SubnetID: "12345",
				},
			},
			true,
		},
		{
			&InstancePair{
				&Instance{
					SubnetID: "12345",
				},
				&Instance{
					SubnetID: "abcde",
				},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := testCase.instancePair.areInstancesInSameSubnet()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected result of areInstancesInSameSubnet to be %t, but it was %t.",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

func TestInstancePairAreInstancesInSameVpc(t *testing.T) {
	testCases := []struct {
		instancePair   *InstancePair
		expectedOutput bool
	}{
		{
			&InstancePair{
				&Instance{
					VpcID: "12345",
				},
				&Instance{
					VpcID: "12345",
				},
			},
			true,
		},
		{
			&InstancePair{
				&Instance{
					VpcID: "12345",
				},
				&Instance{
					VpcID: "abcde",
				},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := testCase.instancePair.areInstancesInSameVpc()
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected result of areInstancesInSameVpc to be %t, but it was %t.",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

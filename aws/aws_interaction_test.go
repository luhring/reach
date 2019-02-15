package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestGetNameTagValueFromTags(t *testing.T) {
	testCases := []struct {
		tags           []*ec2.Tag
		expectedOutput string
	}{
		{
			[]*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("MyInstance"),
				},
			},
			"MyInstance",
		},
		{
			[]*ec2.Tag{
				{
					Key:   aws.String("AnotherTag"),
					Value: aws.String("MyInstance"),
				},
			},
			"",
		},
		{
			[]*ec2.Tag{
				{
					Key:   aws.String("AnotherTag"),
					Value: aws.String("Name"),
				},
			},
			"",
		},
		{
			[]*ec2.Tag{},
			"",
		},
	}

	for _, testCase := range testCases {
		result := getNameTagValueFromTags(testCase.tags)
		if result != testCase.expectedOutput {
			t.Errorf(
				"Expected \"%s\" but got \"%s\".",
				testCase.expectedOutput,
				result,
			)
		}
	}
}

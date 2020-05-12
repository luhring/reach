package apiclient

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const errAWSAuthFailure = "AuthFailure"

func awsErrMessage(err awserr.Error) (message string) {
	switch err.Code() {
	case credentials.ErrNoValidProvidersFoundInChain.Code():
		message = "The AWS SDK was unable to find your API credentials.\n" +
			"\n" +
			"To ensure you've set up your credentials correctly, check out AWS's documentation here:\n" +
			"https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
	case errAWSAuthFailure:
		message = "AWS API authentication failure.\n" +
			"\n" +
			"Please make sure that:\n" +
			"  - your credentials are correct\n" +
			"  - you have the necessary permissions in AWS"
	default:
		message = err.Message()
	}
	return
}

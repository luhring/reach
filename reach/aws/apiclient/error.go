package apiclient

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const errAWSAuthFailure = "AuthFailure"

func awsErrMessage(err awserr.Error) (message string) {
	switch err.Code() {
	case credentials.ErrNoValidProvidersFoundInChain.Code():
		message = "The AWS SDK was unable to find your API credentials. To ensure you've set up your credentials correctly, check out AWS's documentation here: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
	case errAWSAuthFailure:
		message = "AWS API authentication failed. Please make sure your credentials are correct."
	default:
		message = err.Message()
	}
	return
}

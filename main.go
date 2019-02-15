package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	instanceIdentificationInputs := cli.GetInstanceIdentificationInputsFromCli()
	session, err := aws.CreateAwsSession()

	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	ec2Client := ec2.New(session)
	allInstances, err := aws.GetAllInstancesUsingEc2Client(ec2Client)

	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	instancePair := cli.GetInstancePairFromIdentificationInputs(&instanceIdentificationInputs, allInstances)
	resultMessage := instancePair.AnalyzeNetworkAccessAndGenerateResultMessage(ec2Client)

	fmt.Print(resultMessage)
}

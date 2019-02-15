package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/luhring/cnct/network"
)

// InstancePair represents two instances on which analyses will be based
type InstancePair [2]*Instance

type networkAccessDescription struct {
	doesAccessExist bool
	explanation     string
	accessiblePorts []*network.PortRange
}

// AnalyzeNetworkAccessAndGenerateResultMessage ...
func (instancePair *InstancePair) AnalyzeNetworkAccessAndGenerateResultMessage(ec2Client *ec2.EC2) string {
	if false == instancePair.areInstancesInSameVpc() {
		return "This version of the application is unable to determine network accessibility for instances that do not reside in the same VPC.  Please check back in a later version.\n"
	}

	if false == instancePair.areInstancesInSameSubnet() {
		return "This version of the application is unable to determine network accessibility for instances that do not reside in the same subnet.  Please check back in a later version.\n"
	}

	networkAccessDescription :=
		instancePair.describeNetworkAccessToSecondInstanceFromFirstInstance(ec2Client)

	if false == networkAccessDescription.doesAccessExist {
		return instancePair.generateMessageForWhenAccessDoesNotExist() + networkAccessDescription.explanation
	}

	accessiblePortsDescription := "Accessible ports:\n" + network.DescribeListOfPortRanges(networkAccessDescription.accessiblePorts)

	return instancePair.generateMessageForWhenAccessExists() + "\n" + accessiblePortsDescription + "\n"
}

func (instancePair *InstancePair) describeNetworkAccessToSecondInstanceFromFirstInstance(ec2Client *ec2.EC2) *networkAccessDescription {
	analysisOfInstanceStates := instancePair.analyzeNetworkAccessViaInstanceStates()
	analysisOfInstanceSecurityGroups := instancePair.analyzeNetworkAccessViaInstanceSecurityGroups(ec2Client)
	// insert more analyses here...

	explanation := ""
	doesAccessExist := true

	if false == analysisOfInstanceStates.doesAnalysisSuggestThatAccessExists {
		doesAccessExist = false
		explanation = analysisOfInstanceStates.generateExplanationForLackOfAccessDueToInstanceStates()
	} else if false == analysisOfInstanceSecurityGroups.doesAnalysisSuggestThatAccessExists {
		doesAccessExist = false
		explanation = *analysisOfInstanceSecurityGroups.generateExplanationOfLackOfAccess()
	}

	return &networkAccessDescription{
		doesAccessExist,
		explanation,
		analysisOfInstanceSecurityGroups.accessiblePorts,
	}
}

func (instancePair *InstancePair) areInstancesInSameVpc() bool {
	return instancePair[0].VpcID == instancePair[1].VpcID
}

func (instancePair *InstancePair) areInstancesInSameSubnet() bool {
	return instancePair[0].SubnetID == instancePair[1].SubnetID
}

func (instancePair *InstancePair) generateMessageForWhenAccessExists() string {
	message := fmt.Sprintf(
		"Instance '%s' is able to access instance '%s'.\n",
		instancePair[0].GetFriendlyName(),
		instancePair[1].GetFriendlyName(),
	)

	return message
}

func (instancePair *InstancePair) generateMessageForWhenAccessDoesNotExist() string {
	message := fmt.Sprintf(
		"Instance '%s' is unable to access instance '%s'.\n",
		instancePair[0].GetFriendlyName(),
		instancePair[1].GetFriendlyName(),
	)

	return message
}

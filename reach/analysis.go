package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"net"
	"strings"
)

type Analyzer struct {
	AWSSession *session.Session
}

type EC2Instance struct {
	ID                string
	NameTag           string
	State             string
	NetworkInterfaces []NetworkInterface
}

type SecurityGroup struct {
}

type NetworkInterface struct {
	ID                 string
	PrivateIPAddresses []net.IP
	PublicIPAddress    net.IP
	SecurityGroupIDs   []string
	SubnetID           string
	VPCID              string
}

func NewAnalyzer() *Analyzer {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Analyzer{
		AWSSession: awsSession,
	}
}

type InstanceToInstanceVector struct {
	From *EC2Instance
	To   *EC2Instance
}

func (a *Analyzer) AnalyzeVector(vector *InstanceToInstanceVector) {
	// ec2Client := ec2.New(a.AWSSession)

	fmt.Printf("from:\n%v\n", vector.From)
	fmt.Printf("to:\n%v\n", vector.To)
}

func (a *Analyzer) FindEC2Instance(identifier string) (*EC2Instance, error) {
	// identifier could be ID or name tag
	ec2Client := ec2.New(a.AWSSession)

	allInstances, err := GetAllInstancesUsingEC2Client(ec2Client)
	if err != nil {
		return nil, err
	}

	instance, err := getInstanceThatMatchesInput(identifier, allInstances)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func GetAllInstancesUsingEC2Client(ec2Client *ec2.EC2) ([]*EC2Instance, error) {
	describeInstancesOutput, err := ec2Client.DescribeInstances(nil)

	if err != nil {
		return nil, err
	}

	allReservations := describeInstancesOutput.Reservations
	return getAllInstancesFromReservations(allReservations), nil
}

func getAllInstancesFromReservations(reservations []*ec2.Reservation) []*EC2Instance {
	var ec2Instances []*EC2Instance

	for _, r := range reservations {
		for _, i := range r.Instances {
			newItem := ingestEC2Instance(i)
			ec2Instances = append(ec2Instances, newItem)
		}
	}

	return ec2Instances
}

func ingestEC2Instance(instance *ec2.Instance) *EC2Instance {
	networkInterfaces := make([]NetworkInterface, len(instance.NetworkInterfaces))

	for i, networkInterface := range instance.NetworkInterfaces {
		networkInterfaces[i] = ingestNetworkInterface(networkInterface)
	}

	return &EC2Instance{
		ID:                aws.StringValue(instance.InstanceId),
		NameTag:           getNameTagValueFromTags(instance.Tags),
		State:             aws.StringValue(instance.State.Name),
		NetworkInterfaces: networkInterfaces,
	}
}

func ingestNetworkInterface(networkInterface *ec2.InstanceNetworkInterface) NetworkInterface {
	privateIPAddresses := make([]net.IP, len(networkInterface.PrivateIpAddresses))

	for i, address := range networkInterface.PrivateIpAddresses {
		privateIPAddresses[i] = ingestPrivateIPAddress(address)
	}

	securityGroupIDs := make([]string, len(networkInterface.Groups))

	for i, group := range networkInterface.Groups {
		securityGroupIDs[i] = ingestSecurityGroupID(group)
	}

	var publicIPAddress net.IP

	if assoc := networkInterface.Association; assoc != nil {
		if pubIP := assoc.PublicIp; pubIP != nil {
			publicIPAddress = net.ParseIP(aws.StringValue(pubIP))
		}
	}

	return NetworkInterface{
		ID:                 aws.StringValue(networkInterface.NetworkInterfaceId),
		PrivateIPAddresses: privateIPAddresses,
		PublicIPAddress:    publicIPAddress,
		SecurityGroupIDs:   securityGroupIDs,
		SubnetID:           aws.StringValue(networkInterface.SubnetId),
		VPCID:              aws.StringValue(networkInterface.VpcId),
	}
}

func ingestSecurityGroupID(identifier *ec2.GroupIdentifier) string {
	return aws.StringValue(identifier.GroupId)
}

func ingestPrivateIPAddress(address *ec2.InstancePrivateIpAddress) net.IP {
	if address.PrivateIpAddress == nil {
		return nil
	}

	ip := net.ParseIP(aws.StringValue(address.PrivateIpAddress))

	if ip == nil {
		log.Printf("unable to parse IP address '%s'\n", aws.StringValue(address.PrivateIpAddress))
	}

	return ip
}

func getNameTagValueFromTags(tags []*ec2.Tag) string {
	const keyForNameTag = "Name"

	for _, tag := range tags {
		if *tag.Key == keyForNameTag {
			return aws.StringValue(tag.Value)
		}
	}

	return ""
}

func getInstanceThatMatchesInput(
	input string,
	allInstances []*EC2Instance,
) (*EC2Instance, error) {
	var indicesOfInstanceIDSubstringMatches []int
	var indicesOfInstanceNameTagExactMatches []int
	var indicesOfInstanceNameTagSubstringMatches []int
	const noMatchingIndices = -1
	const instanceIDPrefix = "i-"

	matchingIndex := noMatchingIndices

	for index, instance := range allInstances {
		if len(input) >= 3 && strings.EqualFold(input[:2], instanceIDPrefix) {
			if strings.EqualFold(input, instance.ID) {
				// exact match to instance ID (using "i-a1b2c3..." format)
				return instance, nil
			}

			if doesFirstItemMatchBeginningSubstringOfSecondItem(input, instance.ID) {
				indicesOfInstanceIDSubstringMatches =
					append(indicesOfInstanceIDSubstringMatches, index)
			}
		} else if strings.EqualFold(input, instance.NameTag) {
			indicesOfInstanceNameTagExactMatches =
				append(indicesOfInstanceNameTagExactMatches, index)
		} else if doesFirstItemMatchBeginningSubstringOfSecondItem(input, instance.NameTag) {
			indicesOfInstanceNameTagSubstringMatches = append(indicesOfInstanceNameTagSubstringMatches, index)
		}
	}

	countOfInstancesWithIDSubstringMatches := len(indicesOfInstanceIDSubstringMatches)
	countOfInstancesWithNameTagExactMatches := len(indicesOfInstanceNameTagExactMatches)
	countOfInstancesWithNameTagSubstringMatches := len(indicesOfInstanceNameTagSubstringMatches)

	if countOfInstancesWithIDSubstringMatches == 1 {
		matchingIndex = indicesOfInstanceIDSubstringMatches[0]
	} else if countOfInstancesWithNameTagExactMatches == 1 {
		matchingIndex = indicesOfInstanceNameTagExactMatches[0]
	} else if countOfInstancesWithNameTagSubstringMatches == 1 {
		matchingIndex = indicesOfInstanceNameTagSubstringMatches[0]
	}

	if matchingIndex != noMatchingIndices {
		return allInstances[matchingIndex], nil
	}

	if countOfInstancesWithIDSubstringMatches > 1 {
		return nil, fmt.Errorf("the input, '%s', was found the IDs of more than one instance and thus could not be used to uniquely identify an instance", input)
	}

	if countOfInstancesWithNameTagExactMatches > 1 {
		return nil, fmt.Errorf("the input, '%s', exactly matched the 'Name' tags of more than one instance, and thus could not be used to uniquely identify an instance", input)
	}

	if countOfInstancesWithNameTagSubstringMatches > 1 {
		return nil, fmt.Errorf("the input, '%s', was found in the 'Name' tags of more than one instance and thus could not be used to uniquely identify an instance", input)
	}

	return nil, fmt.Errorf("no instances found in EC2 that match the input, '%s'", input)
}

func doesFirstItemMatchBeginningSubstringOfSecondItem(firstItem string, secondItem string) bool {
	lengthOfFirstItem := len(firstItem)

	if lengthOfFirstItem > len(secondItem) {
		return false
	}

	truncatedSecondItem := secondItem[:lengthOfFirstItem]
	return strings.EqualFold(firstItem, truncatedSecondItem)
}

package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/luhring/reach/network"
	"strings"
)

type Analyzer struct {
	AWSSession     *session.Session
	EC2Client      *ec2.EC2
	SecurityGroups map[string]*SecurityGroup
}

func NewAnalyzer() *Analyzer {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Client := ec2.New(awsSession)

	var securityGroups = make(map[string]*SecurityGroup)

	return &Analyzer{
		AWSSession:     awsSession,
		EC2Client:      ec2Client,
		SecurityGroups: securityGroups,
	}
}

func (a *Analyzer) findSecurityGroup(id string) (*SecurityGroup, error) {
	// lookup from cache, return if found
	if group := a.SecurityGroups[id]; group != nil {
		return group, nil
	}

	// lookup from AWS API. If found, store in cache and return
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(id),
		},
	}
	output, err := a.EC2Client.DescribeSecurityGroups(input)
	if err != nil {
		return nil, fmt.Errorf("error: unable to find security group: %v", err)
	}

	groups := output.SecurityGroups
	if len(groups) > 1 {
		return nil, fmt.Errorf("error: found more than one security group matching id: %v", id)
	}

	group := groups[0]

	securityGroup, err := NewSecurityGroup(group)
	if err != nil {
		return nil, fmt.Errorf("error: unable to find security group: %v", err)
	}

	// Save ingested group to local cache
	a.SecurityGroups[id] = securityGroup

	return securityGroup, nil
}

func (a *Analyzer) AnalyzeVector(instanceVector *InstanceVector) {
	fmt.Printf("source:\n%v\n", instanceVector.Source)
	fmt.Printf("destination:\n%v\n", instanceVector.Destination)
	fmt.Println("")

	interfaceVectors := a.createInterfaceVectors(instanceVector)

	if len(interfaceVectors) < 1 {
		fmt.Println("No network interface vectors to analyze.")
		return
	}

	fmt.Printf("Analyzing %v network interface vector(s)...\n\n", len(interfaceVectors))

	var allowedTraffic []*network.TrafficAllowance

	for i, v := range interfaceVectors {
		vectorNumber := i + 1
		fmt.Printf(
			"Vector %v) Analyzing security group rules to determine traffic allowed from source interface (%v) to destination interface (%v):\n\n",
			vectorNumber,
			v.Source.Name,
			v.Destination.Name,
		)

		reachablePortsViaSecurityGroups := v.getAllowedTrafficViaSecurityGroups()

		if len(reachablePortsViaSecurityGroups) >= 1 {
			allowedTraffic = append(allowedTraffic, reachablePortsViaSecurityGroups...)
		}
	}

	allowedTraffic = network.ConsolidateTrafficAllowances(allowedTraffic)

	description := network.DescribeListOfTrafficAllowances(allowedTraffic)
	fmt.Printf("\nAllowed traffic for vector:\n\n%v\n", description)
}

func (a *Analyzer) createInterfaceVectors(instanceVector *InstanceVector) []InterfaceVector {
	if instanceVector.Source.doesStateAllowAccess() == false {
		// Source instance is not running -- ideally, we'd aggregate reasons for the explanation
		// instead of returning early
		return nil
	}

	if instanceVector.Destination.doesStateAllowAccess() == false {
		// Destination instance is not running
		return nil
	}

	var interfaceVectors []InterfaceVector

	for _, fromInterface := range instanceVector.Source.NetworkInterfaces {
		for _, toInterface := range instanceVector.Destination.NetworkInterfaces {
			newVector := InterfaceVector{
				Source:      fromInterface,
				Destination: toInterface,
				PortRange:   instanceVector.PortRange,
			}
			interfaceVectors = append(interfaceVectors, newVector)
		}
	}

	return interfaceVectors
}

func (a *Analyzer) CreateInstanceVector(fromIdentifier, toIdentifier string) (*InstanceVector, error) {
	var vector InstanceVector

	from, err := a.findEC2Instance(fromIdentifier)
	if err != nil {
		return nil, fmt.Errorf("unable to create instance vector: %v", err)
	}

	vector.Source = from

	to, err := a.findEC2Instance(toIdentifier)
	if err != nil {
		return nil, fmt.Errorf("unable to create instance vector: %v", err)
	}

	vector.Destination = to

	return &vector, nil
}

func (a *Analyzer) findEC2Instance(identifier string) (*EC2Instance, error) {
	// identifier could be ID or name tag

	allInstances, err := a.getAllEC2Instances()
	if err != nil {
		return nil, err
	}

	instance, err := getInstanceThatMatchesInput(identifier, allInstances)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a *Analyzer) getAllEC2Instances() ([]*EC2Instance, error) {
	describeInstancesOutput, err := a.EC2Client.DescribeInstances(nil)

	if err != nil {
		return nil, err
	}

	allReservations := describeInstancesOutput.Reservations
	allEC2Instances, err := a.getAllInstancesFromReservations(allReservations)
	if err != nil {
		return nil, fmt.Errorf("error getting all EC2Instances: %v", err)
	}

	return allEC2Instances, nil
}

func (a *Analyzer) getAllInstancesFromReservations(reservations []*ec2.Reservation) ([]*EC2Instance, error) {
	var ec2Instances []*EC2Instance

	for _, r := range reservations {
		for _, i := range r.Instances {
			newItem, err := NewEC2Instance(i, a.findSecurityGroup)
			if err != nil {
				return nil, fmt.Errorf("error getting instance from reservation: %v", err)
			}

			ec2Instances = append(ec2Instances, newItem)
		}
	}

	return ec2Instances, nil
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

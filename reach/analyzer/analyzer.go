package analyzer

import (
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
)

type Analyzer struct {
}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	// Eventually, this dependency wiring should depend on a passed in config.
	provider := api.NewResourceProvider()

	resources := make(map[string]map[string]map[string]reach.Resource)

	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone {
			switch subject.Kind {
			case aws.SubjectKindEC2Instance:
				ec2InstanceSubject := subject.Properties.(aws.EC2InstanceSubject)
				id := ec2InstanceSubject.ID

				ec2Instance, err := provider.GetEC2Instance(id)
				if err != nil {
					log.Fatalf("couldn't get resource: %v", err)
				}
				resources = reach.EnsureResourcePathExists(resources, aws.ResourceDomainAWS, aws.ResourceKindEC2Instance)
				resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][ec2Instance.ID] = ec2Instance.ToResource()

				dependencies, err := ec2Instance.GetDependencies(provider)
				if err != nil {
					return nil, err
				}
				resources = reach.MergeResources(resources, dependencies)
			default:
				log.Fatal("unsupported subject kind")
			}
		}
	}

	return &reach.Analysis{
		Subjects:  subjects,
		Resources: resources,
	}, nil
}

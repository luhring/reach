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
	provider := api.NewResourceGetter()

	var resources []reach.Resource // TODO: Make this a generic store

	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone {
			switch subject.Kind {
			case aws.SubjectKindEC2Instance:
				ec2Props := subject.Properties.(aws.EC2InstanceSubject)
				id := ec2Props.ID
				ec2Instance, err := provider.GetEC2Instance(id)
				if err != nil {
					log.Fatalf("couldn't get resource: %v", err)
				}
				resource := reach.Resource{
					Kind:       aws.ResourceKindEC2Instance,
					Properties: ec2Instance,
				}
				resources = append(resources, resource)
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

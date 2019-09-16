package reach

type Helper interface {
	GetResourcesForSubject(subject Subject) ([]Resource, error)
}

func Analyze(subjects ...Subject) (*Analysis, error) {
	// for each subject, determine needed types —- use correct resource getter to add resource tree to resource store (which is maybe type agnostic?)

	// calculate factors for reachability between subject(s) and destination(s)

	// compute results

	// return full Analysis struct

	// if subjects == nil || len(subjects) < 2 {
	// 	return nil, errors.New("not enough subjects to analyze") // TODO: test for not enough of each role, etc.
	// }
	//
	// for _, subject := range subjects {
	// 	if subject.Kind == SubjectKindEC2Instance {
	// 		ec2Properties := subject.Properties.(aws.EC2InstanceSubjectProperties)
	// 		_ = ec2Properties.ID
	// 		// TODO: load graph!
	// 	} else {
	// 		return nil, fmt.Errorf("unrecognized subject kind '%s'", subject.Kind)
	// 	}
	// }

	analysis := newAnalysis(subjects, nil)
	return analysis, nil
}

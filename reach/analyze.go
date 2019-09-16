package reach

type Helper interface {
	GetResourcesForSubject(subject Subject) ([]Resource, error)
}

// func Analyze(provider reach.ResourceGetter, subjects ...*Subject) (*Analysis, error) {
// 	// for each subject, determine needed types —- use correct resource getter to add resource tree to resource store (which is maybe type agnostic?)
//
// 	var resources []Resource
//
// 	for _, subject := range subjects {
// 		if subject.Role != SubjectRoleNone {
// 			resources = append(resources, subject.GetResources()...)
// 		}
// 	}
//
// 	// calculate factors for reachability between subject(s) and destination(s)
//
// 	// compute results
//
// 	// return full Analysis struct
//
// 	analysis := newAnalysis(subjects, resources)
// 	return analysis, nil
// }

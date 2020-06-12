package reach

import (
	"encoding/json"
	"log"
)

// Analysis is the central structure of a Reach analysis. It describes what subjects were analyzed, what resources were retrieved, and a collection of network vectors between all source-to-destination pairings of subjects.
type Analysis struct {
	Subjects []Subject
	Paths    []AnalyzedPath
}

// NewAnalysis simply creates a new Analysis struct.
func NewAnalysis(subjects []Subject, paths []AnalyzedPath) *Analysis {
	return &Analysis{
		Subjects: subjects,
		Paths:    paths,
	}
}

// ToJSON outputs the Analysis as a JSON string.
func (a *Analysis) ToJSON() string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// MergedTraffic gets the TrafficContent results of each of the analysis's network vectors and returns them as a merged TrafficContent.
func (a *Analysis) MergedTraffic() (TrafficContent, error) {
	result := newTrafficContent()

	for _, path := range a.Paths {
		t := path.TrafficForward()
		mergedTrafficContent, err := result.Merge(t)
		if err != nil {
			return TrafficContent{}, err
		}

		result = mergedTrafficContent
	}

	return result, nil
}

// MergedReturnTraffic gets the return TrafficContent results of each of the analysis's network vectors and returns them as a merged TrafficContent.
func (a *Analysis) MergedReturnTraffic() (TrafficContent, error) {
	// result := newTrafficContent()
	//
	// for _, path := range a.Paths {
	// 	if t := path.ReturnTraffic; t != nil { // TODO: Need to figure out return traffic!
	// 		mergedTrafficContent, err := result.Merge(*t)
	// 		if err != nil {
	// 			return TrafficContent{}, err
	// 		}
	//
	// 		result = mergedTrafficContent
	// 	}
	// }

	panic("not implemented currently!") // TODO: Don't panic! (Come back and fix implementation)
	// return result, nil
}

// PassesAssertReachable determines if the analysis implies the source can reach the destination over at least one protocol whose return path is unobstructed.
func (a Analysis) PassesAssertReachable() bool {
	// forwardTrafficCanReach := false
	//
	// // For each vector, see if there is an obstructed path
	// for _, vector := range a.NetworkVectors {
	// 	if !vector.Traffic.None() {
	// 		forwardTrafficCanReach = true
	//
	// 		for _, p := range vector.Traffic.Protocols() {
	// 			// is return path obstructed (at all) for this protocol?
	// 			if protocolReturnTraffic := vector.ReturnTraffic.protocol(p); !protocolReturnTraffic.complete() {
	// 				return false
	// 			}
	// 		}
	// 	}
	// }
	//
	// if !forwardTrafficCanReach {
	// 	return false
	// }

	return true
}

// PassesAssertNotReachable determines if the analysis implies the source has no way to send network traffic to the destination.
func (a Analysis) PassesAssertNotReachable() bool {
	// Here, we want to be more careful / conservative. If any traffic can get out to destination, fail, regardless of return traffic.

	forwardTraffic, err := a.MergedTraffic()
	if err != nil {
		return false
	}

	if !forwardTraffic.None() {
		return false
	}

	return true
}

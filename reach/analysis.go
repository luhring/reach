package reach

import (
	"encoding/json"
	"log"
)

// Analysis is the central structure of a Reach analysis. It describes what subjects were analyzed, what resources were retrieved, and a collection of network vectors between all source-to-destination pairings of subjects.
type Analysis struct {
	Subjects       []*Subject
	Resources      *ResourceCollection
	NetworkVectors []NetworkVector
}

// NewAnalysis simply creates a new Analysis struct.
func NewAnalysis(subjects []*Subject, resources *ResourceCollection, networkVectors []NetworkVector) *Analysis {
	return &Analysis{
		Subjects:       subjects,
		Resources:      resources,
		NetworkVectors: networkVectors,
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

	for _, v := range a.NetworkVectors {
		if t := v.Traffic; t != nil {
			mergedTrafficContent, err := result.Merge(*v.Traffic)
			if err != nil {
				return TrafficContent{}, err
			}

			result = mergedTrafficContent
		}
	}

	return result, nil
}

package reach

import (
	"encoding/json"
	"log"
)

type Analysis struct {
	Subjects       []*Subject
	Resources      *ResourceCollection
	NetworkVectors []NetworkVector
}

func NewAnalysis(subjects []*Subject, resources *ResourceCollection, networkVectors []NetworkVector) *Analysis {
	return &Analysis{
		Subjects:       subjects,
		Resources:      resources,
		NetworkVectors: networkVectors,
	}
}

func (a *Analysis) ToJSON() string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (a *Analysis) MergedTraffic() (TrafficContent, error) {
	result := NewTrafficContent()

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

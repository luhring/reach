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

package reach

import (
	"encoding/json"
	"log"
)

type Analysis struct {
	Subjects       []*Subject                                `json:"subjects"`
	Resources      map[string]map[string]map[string]Resource `json:"resources"`
	NetworkVectors []NetworkVector                           `json:"networkVectors"`
}

func newAnalysis(subjects []*Subject, resources map[string]map[string]map[string]Resource, networkVectors []NetworkVector) *Analysis {
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

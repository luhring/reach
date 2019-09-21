package reach

import (
	"encoding/json"
	"log"
)

type Analysis struct {
	Subjects  []*Subject                                `json:"subjects"`
	Resources map[string]map[string]map[string]Resource `json:"resources"`
}

func newAnalysis(subjects []*Subject, resources map[string]map[string]map[string]Resource) *Analysis {
	return &Analysis{
		Subjects:  subjects,
		Resources: resources,
	}
}

func (a *Analysis) ToJSON() string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

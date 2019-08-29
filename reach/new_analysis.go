package reach

import (
	"encoding/json"
	"log"
)

type NewAnalysis struct {
	Subjects []subject `json:"subjects"`
}

func newAnalysis(subjects []subject) *NewAnalysis {
	return &NewAnalysis{
		Subjects: subjects,
	}
}

func (a *NewAnalysis) ToJSON() string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

package reach

import "fmt"

const (
	errAnalyzePrefix     = "analysis failed"
	errAnalyzeNilSubject = "one or more subject inputs was nil"
)

func Analyze(source, dest *subject) (*NewAnalysis, error) {
	if source == nil || dest == nil {
		return nil, fmt.Errorf("%s: %s", errAnalyzePrefix, errAnalyzeNilSubject)
	}

	analysis := newAnalysis([]subject{
		*source,
		*dest,
	})

	return analysis, nil
}

package reach

import "github.com/luhring/reach/reach/traffic"

// AnalyzedPath is a path that includes fields for the analysis results of this path.
type AnalyzedPath struct {
	Path
	Traffic               traffic.Content
	ConnectionPredictions ConnectionPredictionSet
}

// TrafficContentsFromAnalyzedPaths returns the set of forward-bound traffic that can traverse the entirety of any of the input paths.
func TrafficContentsFromAnalyzedPaths(paths []AnalyzedPath) []traffic.Content {
	var result []traffic.Content

	for _, p := range paths {
		ft := p.TrafficForward()
		result = append(result, ft)
	}

	return result
}

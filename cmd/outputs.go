package cmd

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach"
)

func handleShowJSON(analysis reach.Analysis) {
	fmt.Println(analysis.ToJSON())
}

func handleShowPaths(analysis reach.Analysis) {
	var pathDescriptions []string

	for _, p := range analysis.Paths {
		pathDescriptions = append(pathDescriptions, fmt.Sprint(p))
	}

	fmt.Print(strings.Join(pathDescriptions, "\n"))
}

func handleDefaultOutput(analysis reach.Analysis) {
	paths := analysis.Paths
	tcs := reach.TrafficContentsFromAnalyzedPaths(paths)
	mergedTraffic := reach.MergeTraffic(tcs...)

	fmt.Print("network traffic allowed from source to destination:" + "\n")
	fmt.Print(mergedTraffic.ColorStringWithSymbols())

	switch len(paths) {
	case 0:
		// ignore
	case 1:
		p := paths[0]
		printConnectionPredictionWarnings(p)
	default:
		// handling this case with care; this view isn't optimized for multi-path output!
		printMergedResultsWarning()
	}
}

package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/luhring/reach/reach"
)

func printMergedResultsWarning() {
	const mergedResultsWarning = "WARNING: Reach detected more than one network path between the source and destination. The analysis result shown above is the merging of all network paths' analysis results, and important analysis information might be obscured or omitted. The impact that infrastructure configuration has on actual network reachability can vary based on the way hosts are configured to use their network interfaces, and Reach is unable to access any configuration internal to a host. To see the complete network reachability for all individual network paths, run the command again with '--" + pathsFlag + "'.\n"
	_, _ = fmt.Fprint(os.Stderr, "\n"+mergedResultsWarning)
}

func printConnectionPredictionWarnings(path reach.AnalyzedPath) {
	connectionPredictions := path.ConnectionPredictions

	var protocolsToCheck []reach.Protocol
	for p := range connectionPredictions {
		protocolsToCheck = append(protocolsToCheck, p)
	}

	sort.Slice(protocolsToCheck, func(i, j int) bool {
		return protocolsToCheck[i].DisplayOrder() < protocolsToCheck[j].DisplayOrder()
	})

	var warnings []string
	for _, protocol := range protocolsToCheck {
		if prediction, ok := connectionPredictions[protocol]; ok && prediction.ShouldWarn() {
			warnings = append(warnings, connectionPredictionWarning(protocol, prediction))
		}
	}

	if len(warnings) >= 1 {
		fmt.Println("\n" + strings.Join(warnings, "\n"))
	}
}

func connectionPredictionWarning(protocol reach.Protocol, prediction reach.ConnectionPrediction) string {
	var warning string

	switch prediction {
	case reach.ConnectionPredictionFailure:
		warning = ansi.Color(warningMessageForProtocolPrognosis(protocol, "will fail"), "red+b")
	case reach.ConnectionPredictionPossibleFailure:
		warning = ansi.Color(warningMessageForProtocolPrognosis(protocol, "might fail"), "yellow+b")
	}

	return warning
}

func warningMessageForProtocolPrognosis(protocol reach.Protocol, prognosis string) string {
	return fmt.Sprintf("Warning: %s communication attempts %s because %s traffic from the destination has an impeded path back to the source.", protocol, prognosis, protocol)
}

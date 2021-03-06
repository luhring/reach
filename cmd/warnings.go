package cmd

import (
	"fmt"
	"os"

	"github.com/luhring/reach/reach"
)

func printMergedResultsWarning() {
	const mergedResultsWarning = "WARNING: Reach detected more than one network path between the source and destination. Reach calls these paths \"network vectors\". The analysis result shown above is the merging of all network vectors' analysis results. The impact that infrastructure configuration has on actual network reachability might vary based on the way hosts are configured to use their network interfaces, and Reach is unable to access any configuration internal to a host. To see the network reachability across individual network vectors, run the command again with '--" + vectorsFlag + "'.\n"
	_, _ = fmt.Fprint(os.Stderr, "\n"+mergedResultsWarning)
}

func warnIfAnyVectorHasRestrictedReturnTraffic(vectors []reach.NetworkVector) {
	for _, v := range vectors {
		if !v.ReturnTraffic.All() {
			const restrictedVectorReturnTraffic = "WARNING: One or more of the analyzed network vectors has restrictions on network traffic allowed to return from the destination to the source. For details, run the command again with '--" + vectorsFlag + "'.\n"
			_, _ = fmt.Fprintf(os.Stderr, "\n"+restrictedVectorReturnTraffic)

			return
		}
	}
}

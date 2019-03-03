package reach

import (
	"github.com/luhring/reach/network"
)

type Analysis struct {
	trafficAllowances []*network.TrafficAllowance
	explanation       Explanation
}

func newAnalysisWithNoTrafficAllowances(explanation Explanation) Analysis {
	return Analysis{
		[]*network.TrafficAllowance{},
		explanation,
	}
}

func (a *Analysis) Results() string {
	return network.DescribeListOfTrafficAllowances(a.trafficAllowances)
}

func (a *Analysis) Explanation() string {
	return a.explanation.Render()
}

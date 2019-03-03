package reach

type Analysis struct {
	trafficAllowances []*TrafficAllowance
	explanation       Explanation
}

func newAnalysisWithNoTrafficAllowances(explanation Explanation) Analysis {
	return Analysis{
		[]*TrafficAllowance{},
		explanation,
	}
}

func (a *Analysis) Results() string {
	return describeListOfTrafficAllowances(a.trafficAllowances)
}

func (a *Analysis) Explanation() string {
	return a.explanation.render()
}

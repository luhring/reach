package reach

type SimpleTrafficContent struct {
	indicator trafficContentIndicator
	protocols map[Protocol]SimpleProtocolContent
}

func SimplifyTrafficContent(tc TrafficContent) SimpleTrafficContent {
	// First, honor the overriding indicator if it's set
	if tc.indicator != trafficContentIndicatorUnset {
		return SimpleTrafficContent{
			indicator: tc.indicator,
		}
	}

	protocols := make(map[Protocol]SimpleProtocolContent)

	for protocol, content := range tc.protocols {
		if content == nil {
			panic("Unable to simplify traffic content due to nil content, please submit an issue with this message!")
		}

		protocols[protocol] = SimplifyProtocolContent(*content)
	}

	return SimpleTrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// TODO: need TC-style methods for getting result for protocol, where method implementation considers the indicator

func (s SimpleTrafficContent) All() bool {
	return s.indicator == trafficContentIndicatorAll
}

func (s SimpleTrafficContent) None() bool {
	return s.indicator == trafficContentIndicatorNone
}

func (s SimpleTrafficContent) Protocol(p Protocol) SimpleProtocolContent {
	if s.All() {
		return SimpleProtocolContentAll
	}

	if s.None() {
		return SimpleProtocolContentNone
	}

	return s.protocols[p]
}

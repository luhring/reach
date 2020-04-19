package reach

type SimpleProtocolContent int

const (
	SimpleProtocolContentNone SimpleProtocolContent = iota
	SimpleProtocolContentSome
	SimpleProtocolContentAll
)

func SimplifyProtocolContent(content ProtocolContent) SimpleProtocolContent {
	if content.complete() {
		return SimpleProtocolContentAll
	}

	if content.empty() {
		return SimpleProtocolContentNone
	}

	return SimpleProtocolContentSome
}

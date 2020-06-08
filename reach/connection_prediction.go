package reach

// ConnectionPrediction describes a prediction of success for a network connection using a particular protocol
type ConnectionPrediction int

// The possible values for ConnectionPrediction
const (
	ConnectionPredictionUnknown ConnectionPrediction = iota
	ConnectionPredictionSuccess
	ConnectionPredictionPossibleFailure
	ConnectionPredictionFailure
)

func (cp ConnectionPrediction) String() string {
	switch cp {
	case ConnectionPredictionSuccess:
		return "success"
	case ConnectionPredictionPossibleFailure:
		return "possible failure"
	case ConnectionPredictionFailure:
		return "failure"
	default:
		return "unknown"
	}
}

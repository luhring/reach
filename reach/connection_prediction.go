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

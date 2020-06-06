package reach

type ConnectionPrediction int

const (
	ConnectionPredictionUnknown ConnectionPrediction = iota
	ConnectionPredictionGuaranteedSuccess
	ConnectionPredictionPotentialFailure
	ConnectionPredictionGuaranteedFailure
)

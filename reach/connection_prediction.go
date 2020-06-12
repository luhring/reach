package reach

import "encoding/json"

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
		return "possible-failure"
	case ConnectionPredictionFailure:
		return "failure"
	default:
		return "unknown"
	}
}

// MarshalJSON returns the JSON representation fo the ConnectionPrediction.
func (cp ConnectionPrediction) MarshalJSON() ([]byte, error) {
	return json.Marshal(cp.String())
}

// ShouldWarn returns a bool to indicate whether this prediction warrants a warning to the user.
func (cp ConnectionPrediction) ShouldWarn() bool {
	return cp == ConnectionPredictionPossibleFailure || cp == ConnectionPredictionFailure
}

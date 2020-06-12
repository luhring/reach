package reach

import (
	"encoding/json"

	"github.com/luhring/reach/reach/traffic"
)

// ConnectionPredictionSet wraps the common map association between protocols and their predictions.
type ConnectionPredictionSet map[traffic.Protocol]ConnectionPrediction

// MarshalJSON returns the JSON representation of the ConnectionPredictionSet.
func (cps ConnectionPredictionSet) MarshalJSON() ([]byte, error) {
	result := make(map[string]string)
	for protocol, prediction := range cps {
		result[protocol.String()] = prediction.String()
	}
	return json.Marshal(result)
}

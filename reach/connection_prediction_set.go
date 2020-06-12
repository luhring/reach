package reach

import "encoding/json"

// ConnectionPredictionSet wraps the common map association between protocols and their predictions.
type ConnectionPredictionSet map[Protocol]ConnectionPrediction

// MarshalJSON returns the JSON representation of the ConnectionPredictionSet.
func (cps ConnectionPredictionSet) MarshalJSON() ([]byte, error) {
	result := make(map[string]string)
	for protocol, prediction := range cps {
		result[protocol.String()] = prediction.String()
	}
	return json.Marshal(result)
}

package reach

type TrafficConnectionPrediction struct {
	protocols map[Protocol]ProtocolConnectionPrediction
}

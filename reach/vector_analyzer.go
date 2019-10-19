package reach

// A VectorAnalyzer can calculate analysis factors for a given network vector.
type VectorAnalyzer interface {
	Factors(v NetworkVector) ([]Factor, NetworkVector, error)
}

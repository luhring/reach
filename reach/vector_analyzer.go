package reach

type VectorAnalyzer interface {
	Factors(v NetworkVector) ([]Factor, NetworkVector, error)
	// Explain(v NetworkVector) string
}

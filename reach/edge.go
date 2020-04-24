package reach

type Edge struct {
	Tuple             IPTuple            // The IP tuple state for IP packets traveling along this edge
	EndRef            UniversalReference // The Ref for the infrastructure at the latter end of the edge
	ConnectsInterface bool               // Does edge connect a network interface to its attached entity
}

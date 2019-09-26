package reach

type Resource struct {
	Kind       string      `json:"kind"`
	Properties interface{} `json:"properties"`
}

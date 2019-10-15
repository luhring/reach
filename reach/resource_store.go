package reach

type ResourceStore interface {
	ExportAll() []Resource
	Save(kind string, resource interface{})
}

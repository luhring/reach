package reach

func FindEC2InstanceID(searchText string) (string, error) {
	mgr := NewAWSManager() // TODO: Replace implementation for AWS client operations

	instance, err := mgr.findEC2Instance(searchText)
	if err != nil {
		return "", err
	}

	return instance.ID, nil
}

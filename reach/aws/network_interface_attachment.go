package aws

type NetworkInterfaceAttachment struct {
	ID                        string
	ElasticNetworkInterfaceID string
	DeviceIndex               int64 // e.g. 0 for "eth0"
}

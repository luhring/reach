package aws

type NetworkInterfaceAttachment struct {
	ID                        string `json:"id"`
	ElasticNetworkInterfaceID string `json:"elasticNetworkInterfaceID"`
	DeviceIndex               int64  `json:"deviceIndex"` // e.g. 0
}

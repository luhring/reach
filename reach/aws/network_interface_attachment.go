package aws

type NetworkInterfaceAttachment struct {
	ID                        string `json:"id"`
	ElasticNetworkInterfaceID string `json:"elasticNetworkInterfaceID"`
	DeviceName                string `json:"deviceName"` // e.g. "eth0"
}

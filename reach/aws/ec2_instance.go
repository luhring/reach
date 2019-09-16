package aws

const ResourceKindEC2Instance = "EC2Instance"

type EC2Instance struct {
	ID                          string                       `json:"id"`
	NameTag                     string                       `json:"nameTag"`
	State                       string                       `json:"state"`
	NetworkInterfaceAttachments []NetworkInterfaceAttachment `json:"networkInterfaceAttachments"`
}

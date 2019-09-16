package aws

type SecurityGroupReference struct {
	ID        string `json:"id"`
	AccountID string `json:"accountID"`
	NameTag   string `json:"nameTag"`
	GroupName string `json:"groupName"`
}

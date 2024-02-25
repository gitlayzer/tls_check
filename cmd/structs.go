package cmd

type CertInfo struct {
	Domain    string `json:"domain"`
	Subject   string `json:"subject"`
	ExpiresOn string `json:"expires_on"`
	DaysLeft  int    `json:"days_left"`
}

package signalmessengerutils

type Envelope struct {
	Source                   string `json:"source"`
	SourceNumber             string `json:"sourceNumber"`
	SourceUuid               string `json:"sourceUuid"`
	SourceName               string `json:"sourceName"`
	SourceDevice             int    `json:"sourceDevice"`
	Timestamp                int64  `json:"timestamp"`
	ServerReceivedTimestamp  int64  `json:"serverReceivedTimestamp"`
	ServerDeliveredTimestamp int64  `json:"serverDeliveredTimestamp"`
	DataMessage              struct {
		Timestamp          int64  `json:"timestamp"`
		Message            string `json:"message"`
		ExpiresInSeconds   int    `json:"expiresInSeconds"`
		IsExpirationUpdate bool   `json:"isExpirationUpdate"`
		ViewOnce           bool   `json:"viewOnce"`
	} `json:"dataMessage"`
}

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
	} `json:"dataMessage,omitempty"`
	SyncMessage  *SyncMessage  `json:"syncMessage,omitempty"`
	ReceiptMessage *ReceiptMessage `json:"receiptMessage,omitempty"`
}

type SyncMessage struct {
	SentMessage struct {
		Destination      string `json:"destination"`
		DestinationNumber string `json:"destinationNumber"`
		DestinationUuid   string `json:"destinationUuid"`
		Timestamp        int64  `json:"timestamp"`
		Message          string `json:"message"`
		ExpiresInSeconds int    `json:"expiresInSeconds"`
		IsExpirationUpdate bool   `json:"isExpirationUpdate"`
		ViewOnce         bool   `json:"viewOnce"`
	} `json:"sentMessage,omitempty"`
}

type ReceiptMessage struct {
	When       int64  `json:"when"`
	IsDelivery bool   `json:"isDelivery"`
	IsRead     bool   `json:"isRead"`
	IsViewed   bool   `json:"isViewed"`
	Timestamps []int64 `json:"timestamps"`
}

package signalmessengerutils

// Message represents a signal message.
type Message struct {
	Envelope *Envelope `json:"envelope"`
	Account  string    `json:"account"`
}

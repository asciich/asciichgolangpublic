package messengerinterfaces

type Message interface{
	GetContentAsString() (string, error)
	GetSenderAccountAsString() (string, error)
	GetRecipientsAsStringSlice() ([]string, error)
	GetTimestampMilliseconds() (int64, error)

	IsDataMessage() (bool, error)
	IsSenderAccount(string) (bool, error)
}
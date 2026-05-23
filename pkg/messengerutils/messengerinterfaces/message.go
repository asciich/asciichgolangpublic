package messengerinterfaces

type Message interface{
	GetContentAsString() (string, error)
	GetSenderAccountAsString() (string, error)
	GetTimestampMilliseconds() (int64, error)
}
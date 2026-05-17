package messengerinterfaces

type Message interface{
	GetContentAsString() (string, error)
}
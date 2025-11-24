package documentinterfaces

type Element interface {
	GetPlainText() (plainText string, err error)
}

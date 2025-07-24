package documentbase

type Element interface{
	GetPlainText() (plainText string, err error)
}

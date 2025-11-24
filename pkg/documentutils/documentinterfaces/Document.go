package documentinterfaces

type Document interface {
	AddTitleByString(title string) (err error)
	GetElements() (elements []Element)
}

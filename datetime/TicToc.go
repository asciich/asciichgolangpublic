package datetime

import (
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TicToc struct {
	title  string
	tStart *time.Time
}

func MustTic(title string, verbose bool) (t *TicToc) {
	t, err := Tic(title, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return t
}

func NewTicToc() (t *TicToc) {
	return new(TicToc)
}

func Tic(title string, verbose bool) (t *TicToc, err error) {
	if title == "" {
		return nil, tracederrors.TracedError("title is empty string")
	}

	t = NewTicToc()

	err = t.SetTitle(title)
	if err != nil {
		return nil, err
	}

	t.Start(verbose)
	return t, nil
}

func TicWithoutTitle(verbose bool) (t *TicToc) {
	t = NewTicToc()
	t.Start(verbose)
	return t
}

func (t *TicToc) GetTStart() (tStart *time.Time, err error) {
	if t.tStart == nil {
		return nil, tracederrors.TracedErrorf("tStart not set")
	}

	return t.tStart, nil
}

func (t *TicToc) GetTitle() (title string, err error) {
	if t.title == "" {
		return "", tracederrors.TracedErrorf("title not set")
	}

	return t.title, nil
}

func (t *TicToc) GetTitleOrDefaultIfUnset() (title string) {
	if t.title != "" {
		return t.title
	}

	return "TicToc"
}

func (t *TicToc) MustGetTStart() (tStart *time.Time) {
	tStart, err := t.GetTStart()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tStart
}

func (t *TicToc) MustGetTitle() (title string) {
	title, err := t.GetTitle()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return title
}

func (t *TicToc) MustSetTStart(tStart *time.Time) {
	err := t.SetTStart(tStart)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TicToc) MustSetTitle(title string) {
	err := t.SetTitle(title)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TicToc) MustToc(verbose bool) (elapsedTime *time.Duration) {
	elapsedTime, err := t.Toc(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return elapsedTime
}

func (t *TicToc) SetTStart(tStart *time.Time) (err error) {
	if tStart == nil {
		return tracederrors.TracedErrorf("tStart is nil")
	}

	t.tStart = tStart

	return nil
}

func (t *TicToc) SetTitle(title string) (err error) {
	if title == "" {
		return tracederrors.TracedErrorf("title is empty string")
	}

	t.title = title

	return nil
}

func (t *TicToc) Start(verbose bool) {
	timeNow := time.Now()
	t.tStart = &timeNow

	title := t.GetTitleOrDefaultIfUnset()

	if verbose {
		logging.LogInfof("TicToc timer '%s': started", title)
	}
}

func (t *TicToc) Toc(verbose bool) (elapsedTime *time.Duration, err error) {
	tStart, err := t.GetTStart()
	if err != nil {
		return nil, err
	}

	elapsedDurationValue := time.Since(*tStart)
	elapsedDuration := &elapsedDurationValue

	title := t.GetTitleOrDefaultIfUnset()

	if verbose {
		elapsedDurationString, err := FormatDurationAsString(elapsedDuration)
		if err != nil {
			return nil, err
		}

		logging.LogInfof("TicToc timer '%s': elapsed duration: %s", title, elapsedDurationString)
	}

	return elapsedDuration, nil
}

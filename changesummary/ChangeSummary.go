package changesummary

import (
	"log"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A ChangeSummary is used to return details about if/what/how much was actually changed.
// Since a lot of functions are written idempotent the ChangeSummary should be used to inform the caller about what actually was done to reach a desired state.
type ChangeSummary struct {
	numberOfChanges int
	childSummaries  []*ChangeSummary
}

func NewChangeSummary() (c *ChangeSummary) {
	return new(ChangeSummary)
}

func (c *ChangeSummary) AddChildSummary(childSummary *ChangeSummary) (err error) {
	if childSummary == nil {
		return tracederrors.TracedErrorNil("childSummary")
	}

	c.childSummaries = append(c.childSummaries, childSummary)

	return nil
}

func (c *ChangeSummary) GetChildSummaries() (childSummaries []*ChangeSummary, err error) {
	if c.childSummaries == nil {
		return nil, tracederrors.TracedErrorf("childSummaries not set")
	}

	if len(c.childSummaries) <= 0 {
		return nil, tracederrors.TracedErrorf("childSummaries has no elements")
	}

	return c.childSummaries, nil
}

func (c *ChangeSummary) GetIsChanged() (isChanged bool) {
	return c.numberOfChanges != 0
}

func (c *ChangeSummary) GetNumberOfChanges() (numberOfChanges int) {
	return c.numberOfChanges
}

func (c *ChangeSummary) IncrementNumberOfChanges() {
	c.numberOfChanges += 1
}

func (c *ChangeSummary) IsChanged() (isChanged bool) {
	if c.GetIsChanged() {
		return true
	}

	for _, child := range c.childSummaries {
		if child.IsChanged() {
			return true
		}
	}

	return false
}

func (c *ChangeSummary) MustAddChildSummary(childSummary *ChangeSummary) {
	err := c.AddChildSummary(childSummary)
	if err != nil {
		log.Panic(err)
	}
}

func (c *ChangeSummary) MustGetChildSummaries() (childSummaries []*ChangeSummary) {
	childSummaries, err := c.GetChildSummaries()
	if err != nil {
		log.Panic(err)
	}

	return childSummaries
}

func (c *ChangeSummary) MustSetChildSummaries(childSummaries []*ChangeSummary) {
	err := c.SetChildSummaries(childSummaries)
	if err != nil {
		log.Panic(err)
	}
}

func (c *ChangeSummary) MustSetNumberOfChanges(numberOfChanges int) {
	err := c.SetNumberOfChanges(numberOfChanges)
	if err != nil {
		log.Panic(err)
	}
}

func (c *ChangeSummary) SetChanged(isChanged bool) {
	c.SetIsChanged(isChanged)
}

func (c *ChangeSummary) SetChildSummaries(childSummaries []*ChangeSummary) (err error) {
	if childSummaries == nil {
		return tracederrors.TracedErrorf("childSummaries is nil")
	}

	if len(childSummaries) <= 0 {
		return tracederrors.TracedErrorf("childSummaries has no elements")
	}

	c.childSummaries = childSummaries

	return nil
}

func (c *ChangeSummary) SetIsChanged(isChanged bool) {
	if c.IsChanged() == isChanged {
		return
	}

	if isChanged {
		c.numberOfChanges = 1
	} else {
		c.numberOfChanges = 0
	}
}

func (c *ChangeSummary) SetNumberOfChanges(numberOfChanges int) (err error) {
	if numberOfChanges < 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for numberOfChanges", numberOfChanges)
	}

	c.numberOfChanges = numberOfChanges

	return nil
}

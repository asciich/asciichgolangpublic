package logging

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/changesummary"
)

var overrideFunctionLogByChangeSummary func(changeSummary *changesummary.ChangeSummary, message string)

func OverrideLogByChangeSummary(overrideFunction func(changeSummary *changesummary.ChangeSummary, message string)) {
	overrideFunctionLogByChangeSummary = overrideFunction
}

func LogByChangeSummary(changeSummary *changesummary.ChangeSummary, message string) {
	if overrideFunctionLogByChangeSummary != nil {
		overrideFunctionLogByChangeSummary(changeSummary, message)
		return
	}

	isChanged := false

	if changeSummary != nil {
		isChanged = changeSummary.IsChanged()
	}

	if isChanged {
		LogChanged(message)
	} else {
		LogInfo(message)
	}
}

var overrideFunctionLogByChangeSummaryf func(changeSummary *changesummary.ChangeSummary, message string, args ...interface{})

func OverrideLogByChangeSummaryf(overrideFunction func(changeSummary *changesummary.ChangeSummary, message string, args ...interface{})) {
	overrideFunctionLogByChangeSummaryf = overrideFunction
}

func LogByChangeSummaryf(changeSummary *changesummary.ChangeSummary, message string, args ...interface{}) {
	if overrideFunctionLogByChangeSummaryf != nil {
		overrideFunctionLogByChangeSummaryf(changeSummary, message, args)
		return
	}

	formattedMessage := fmt.Sprintf(message, args...)

	LogByChangeSummary(changeSummary, formattedMessage)
}

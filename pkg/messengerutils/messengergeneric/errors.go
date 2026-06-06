package messengergeneric

import "errors"

var ErrEmptyMessageSlice = errors.New("message slice is empty")
var ErrNoDataMessageFound = errors.New("no data message found")
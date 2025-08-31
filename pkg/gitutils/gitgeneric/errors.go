package gitgeneric

import "errors"

var ErrGitRepositoryDoesNotExist = errors.New("gitRepository does not exist")
var ErrGitRepositoryHeadNotFound = errors.New("gitRepository head not found")

package apperor

import "errors"

var (
	ErrNoEffect     = errors.New("no effect") // Может пригодиться для Delete
	ErrRepoNotFound = errors.New("not found")
)

package core

import "errors"

// ErrStorageUnavailable indicates storage backend is not accessible
var ErrStorageUnavailable = errors.New("storage unavailable")

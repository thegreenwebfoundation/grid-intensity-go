package unfccc

import (
	"errors"
)

var (
	ErrNoMatchingRegion error = errors.New("no region matched that region code")
)

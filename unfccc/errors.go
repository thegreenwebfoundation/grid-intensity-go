package unfccc

import (
	"errors"
)

var (
	ErrNoMatchingRegion error = errors.New("No region matched that region code")
)

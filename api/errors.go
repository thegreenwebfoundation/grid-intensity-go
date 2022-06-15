package api

import "errors"

var (
	ErrNoRegionProvided = errors.New("no region provided")
	ErrTimeout          = errors.New("timed out")
)

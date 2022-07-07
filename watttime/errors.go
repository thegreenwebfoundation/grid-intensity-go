package watttime

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrNoRegionProvided           error = errors.New("no region was provided")
	ErrNoMarginalIntensityPresent error = errors.New("no marginal intensity present")
	ErrNoRelativeIntensityPresent error = errors.New("no relative intensity present")
	ErrReceived403Forbidden       error = errors.New("received 403 forbidden")
	ErrReceivedNon200Status       error = errors.New("received non-200 status")
)

func errBadStatus(resp *http.Response) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("could not read error response: %w", err)
	} else {
		err = errors.New(string(data))
	}

	return fmt.Errorf("%s - %s: %w", resp.Status, err, ErrReceivedNon200Status)
}

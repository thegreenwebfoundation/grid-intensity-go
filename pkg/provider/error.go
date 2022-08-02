package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrInvalidRegion              error = errors.New("region is not supported by this provider")
	ErrNoMarginalIntensityPresent error = errors.New("no marginal intensity present")
	ErrNoRelativeIntensityPresent error = errors.New("no relative intensity present")
	ErrNoResponse                 error = errors.New("no data was received in response, try again later")
	ErrUnknownResponse            error = errors.New("unknown index received")
	ErrReceivedNon200Status       error = errors.New("received non-200 status")
	ErrReceived403Forbidden       error = errors.New("received 403 forbidden")
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

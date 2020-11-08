package carbonintensity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrOnlyUK               error = errors.New("only UK is supported by this provider")
	ErrNoResponse           error = errors.New("no data was received in response, try again later")
	ErrUnknownResponse      error = errors.New("unknown index received")
	ErrReceivedNon200Status error = errors.New("received non-200 status")
)

func errUnknownIndex(index string) error {
	return fmt.Errorf("'%s': %w", index, ErrUnknownResponse)
}

func errBadStatus(resp *http.Response) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("could not read error response: %w", err)
	} else {
		err = errors.New(string(data))
	}

	return fmt.Errorf("%s - %s: %w", resp.Status, err, ErrReceivedNon200Status)
}

package utils

import (
	"errors"
	"fmt"
	Url "net/url"
	"net"
)

var ErrIncorrectUrlSchema = errors.New("Url schema is incorrect")
var ErrNotSupported = errors.New("The following is not supported")

func IsCorrectUrl(url string) (bool, error) {
	uri, err := Url.ParseRequestURI(url)
	if err != nil {
		return false, err
	}
	
	switch uri.Scheme {
	case "http":
	case "https":
	default:
		return false, fmt.Errorf("%v: %w", url, ErrIncorrectUrlSchema)
	}

	_, err = net.LookupHost(uri.Host)
	if err != nil {
		return false, err
	}

	return true, nil
}

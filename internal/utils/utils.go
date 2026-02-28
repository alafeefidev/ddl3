package utils

import (
	"errors"
	"fmt"
	"strings"
)

var ErrIncorrectUrlSchema = errors.New("Url schema is incorrect")
var ErrNotSupported       = errors.New("The following is not supported") 

func IsCorrectHttpUrl(url string) error {
	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("%v: %w", url, ErrIncorrectUrlSchema)
	}
	return nil
}

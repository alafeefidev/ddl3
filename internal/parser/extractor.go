package parser

import (
	"errors"
	"strings"

	"github.com/alafeefidev/ddl3/internal/parser/hls"
)

type ParsingConfig struct {
	url string
}

func LoadFromText(content *string) error {
	*content = strings.Trim(*content, " ")

	if strings.Contains(*content, hls.ExtM3u) {
		return errors.New("Manage hls downloads")
	} else if strings.Contains(*content, "</MPD>") && strings.Contains(*content, "<MPD") {
		//TODO
		return errors.New("Manage dash/mpd downloads")
	} else {
		return errors.New("Extension not supported")
	}
}

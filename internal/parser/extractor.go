package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alafeefidev/ddl3/internal/httputils"
	"github.com/alafeefidev/ddl3/internal/parser/hls"
	"github.com/alafeefidev/ddl3/internal/utils"
)

type ParsingConfig struct {
	Url string
	OriginalUrl string
	BaseUrl string
}

func LoadFromUri(uri string) error {
	if err := utils.IsCorrectHttpUrl(uri); err == nil {
		client := &httputils.HttpClient{BaseUrl: uri}
		resp, err := client.ReqContentGet("")
		if err != nil {
			return err
		}
		//TODO return LoadFromText(*resp)
		resp = resp // remove
		return nil  // remove

	} else if err != nil {
		return err
	} else {
		return fmt.Errorf("Unsupported schema %w", utils.ErrNotSupported)
	}
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

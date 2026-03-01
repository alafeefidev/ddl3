package parser

import (
	"errors"
	"fmt"
	"strings"
	"encoding/xml"

	"github.com/alafeefidev/ddl3/internal/httputils"
	"github.com/alafeefidev/ddl3/internal/parser/hls"
	MPD "github.com/alafeefidev/ddl3/internal/parser/mpd"
	"github.com/alafeefidev/ddl3/internal/utils"
)

type ParsingConfig struct {
	Url         string
	OriginalUrl string
	BaseUrl     string
}

//TODO change return to be able to return mpd and hls
func LoadFromUri(uri string) (*MPD.Mpd, error) {
	if ok, err := utils.IsCorrectUrl(uri); ok {
		client := httputils.NewHttpClient(uri)
		resp, err := client.ReqContentGet("")
		if err != nil {
			return nil, err
		}
		
		mpd, err := LoadFromText(resp)
		if err != nil {
			return nil, err
		}
		return mpd, nil

	} else {
		return nil, fmt.Errorf("Unsupported url %w: %w", utils.ErrNotSupported, err)
	}
}

func LoadFromText(content []byte) (*MPD.Mpd, error) {
	s := strings.TrimSpace(string(content))

	if strings.Contains(s, hls.ExtM3u) {
		return nil, errors.New("Manage hls downloads")
	} else if strings.Contains(s, "</MPD>") && strings.Contains(s, "<MPD") {
		var mpd MPD.Mpd
		if err := xml.Unmarshal([]byte(s), &mpd); err != nil {
			return nil, err
		}
		return &mpd, nil

	} else {
		return nil, errors.New("Extension not supported")
	}
}

package utils

import (
	"errors"
	"fmt"
	"net"
	Url "net/url"
	"strings"
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

func ResolveUrl(urls ...string) string {
	final := ""
	for _, url := range urls {
		if url == "" {
			continue
		}
		if strings.Contains(url, "://") {
			final = url
		} else {
			final = strings.TrimRight(final, "/") + "/" + url
		}
	}
	return final
}

func StripUrlFilename(url string) string {
	if i := strings.LastIndex(url, "/"); i>= 0 {
		return url[:i+1]
	}
	return url
}

func CodecRepr(codec string) string {
	// Return codec with fallback to provided codec string if not found
	c := strings.ToLower(codec)
	switch {
	case strings.HasPrefix(c, "avc1"):
		return "H.264"
	case strings.HasPrefix(c, "hvc1"):
		return "H.265"
	case strings.HasPrefix(c, "hev1"):
		return "HEVC"
	case strings.HasPrefix(c, "av01"):
		return "AV1"
	case strings.HasPrefix(c, "vp09"):
		return "VP9"
	case strings.EqualFold(c, "mp4a.40.2"):
		return "AAC-LC"
	case strings.EqualFold(c, "mp4a.40.5"):
		return "HE-AAC (v1)"
	case strings.EqualFold(c, "mp4a.40.29"):
		return "HE-AAC v2"
	case strings.EqualFold(c, "ac-3"):
		return "Dolby AC-3"
	case strings.EqualFold(c, "ec-3"):
		return "Dolby E-AC-3 (Atmos)"
	default:
		return c
	}
}

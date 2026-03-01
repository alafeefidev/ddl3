package enums

import "fmt"

type MediaType int

const (
	VideoType MediaType = iota
	AudioType
	SubtitleType
)

var mediaName = map[MediaType]string{
	VideoType:    "video",
	AudioType:    "audio",
	SubtitleType: "text",
}

var mediaTypeByName = map[string]MediaType{
	"video": VideoType,
	"audio": AudioType,
	"text":  SubtitleType,
}

func ParseMediaType(s string) (MediaType, error) {
	ft, ok := mediaTypeByName[s]
	if !ok {
		return -1, fmt.Errorf("unknown file type: %q", s)
	}
	return ft, nil

}

func (ft MediaType) String() string {
	return mediaName[ft]
}

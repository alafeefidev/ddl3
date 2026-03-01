package mpd

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/alafeefidev/ddl3/internal/enums"
)

// https://standards.iso.org/ittf/PubliclyAvailableStandards/MPEG-DASH_schema_files/DASH-MPD.xsd
// https://github.com/MPEGGroup/DASHSchema/blob/6th-Ed/DASH-MPD.xsd
type Mpd struct {
	// Skipped: profiles, availabilityEndTime, minimumUpdatePeriod, minBufferTime
	XMLName                   xml.Name  `xml:"MPD"`
	XMLNS                     *string   `xml:"xmlns,attr"`
	Type                      *string   `xml:"type,attr"`                  // static & dynamic
	AvailabilityStartTime     *string   `xml:"availabilityStartTime,attr"` // Where to start
	PublishTime               *string   `xml:"publishTime,attr"`
	MediaPresentationDuration *string   `xml:"mediaPresentationDuration,attr"` // Total length of file
	TimeShiftBufferDepth      *string   `xml:"timeShiftBufferDepth,attr"`
	MaxSegmentDuration        *string   `xml:"maxSegmentDuration,attr"` // Max length of each segmentation, def PT1M
	BaseUrl                   *string   `xml:"BaseURL,omitempty"`
	Periods                   []*Period `xml:"Period,omitempty"`
}

// Part of a content, usually 1 with a duration
type Period struct {
	Duration       *string          `xml:"duration,attr"`
	Id             *string          `xml:"id,attr"`
	BaseUrl        *string          `xml:"BaseURL,omitempty"` // combine with Mpd BaseUrl
	AdaptationSets []*AdaptationSet `xml:"AdaptationSet,omitempty"`
}

// Media streams, videos, audios, subs
type AdaptationSet struct {
	//TODO maybe add fallback for MimeType to contentType?
	MimeType                  *string              `xml:"mimeType,attr"`    // ex. audio/mp4
	ContentType               *string              `xml:"contentType,attr"` // ex. video, audio, text
	FrameRate                 *string              `xml:"frameRate,attr"`
	BaseUrl                   *string              `xml:"BaseURL,omitempty"` // combine with Period, Mpd BaseUrl
	Codecs                    *string              `xml:"codecs,attr"`
	Width                     *uint64              `xml:"width,attr"`
	Height                    *uint64              `xml:"height,attr"`
	Lang                      *string              `xml:"lang,attr"`
	Representations           []*Representation    `xml:"Representation,omitempty"`
	ContentProtections        []*ContentProtection `xml:"ContentProtection,omitempty"`
	AudioChannelConfiguration *struct {
		Value *string `xml:"value,attr"`
	} `xml:"AudioChannelConfiguration,omitempty"`
}

// Different types of an adaptation, quality, and more
type Representation struct {
	Id          *string `xml:"id,attr"`
	FrameRate   *string `xml:"frameRate,attr"`
	MimeType    *string `xml:"mimeType,attr"`
	ContentType *string `xml:"contentType,attr"` // ex. video, audio, text
	Bandwith    *uint64 `xml:"bandwidth,attr"`
	Codecs      *string `xml:"codecs,attr"`
	Lang        *string `xml:"lang,attr"`
	Width       *uint64 `xml:"width,attr"`
	Height      *uint64 `xml:"height,attr"`
	BaseUrl     *string `xml:"BaseURL,omitempty"` // combine with AdaptationSet, Period, Mpd BaseUrl
}

// https://dashif.org/identifiers/content_protection/
type ContentProtection struct {
	SchemeIdUri *string `xml:"schemeIdUri,attr"` // For identifying encryption type
	Value       *string `xml:"value,attr"`
	DefaultKid  *string `xml:"urn:mpeg:cenc:2013 default_KID,attr"` // default encryption
	PSSH        *string `xml:"urn:mpeg:cenc:2013 pssh"`             // Mostly widevine
	Pro         *string `xml:"urn:microsoft:playready pro"`         // playready
}

type MediaEncryption struct {
	Type  string  // cenc, widevine, playready
	Value string  // DefaultKid for cenc, pssh for widevine and playready
	Pro   *string // only for playready
}

func (a *AdaptationSet) GetAllEncryption() ([]MediaEncryption, error) {
	var cps []MediaEncryption
	//TODO manage if error to just skip, hmmmmmmm.... thinking ahhhh
	for _, cp := range a.ContentProtections {
		en, err := cp.GetEncryption()
		if err != nil {
			return nil, err
		}
		cps = append(cps, *en)
	}
	return cps, nil
}

func (cp *ContentProtection) GetEncryption() (*MediaEncryption, error) {
	// Extra verification maybe needed: strings.EqualFold(*cp.Value, "cenc")
	if cp.DefaultKid != nil {
		return &MediaEncryption{
			Type:  "cenc",
			Value: *cp.DefaultKid,
		}, nil
	}

	if cp.SchemeIdUri != nil {
		if strings.EqualFold(*cp.SchemeIdUri, "urn:uuid:EDEF8BA9-79D6-4ACE-A3C8-27DCD51D21ED") {
			// Widevine
			if cp.PSSH != nil {
				return &MediaEncryption{
					Type:  "widevine",
					Value: *cp.PSSH,
				}, nil
			}
			return nil, fmt.Errorf("no pssh key found in widevine element")
		}

		if strings.EqualFold(*cp.SchemeIdUri, "urn:uuid:9A04F079-9840-4286-AB92-E65BE0885F95") {
			// PlayReady
			if cp.PSSH != nil && cp.Pro != nil {
				return &MediaEncryption{
					Type:  "playready",
					Value: *cp.PSSH,
					Pro:   cp.Pro, // hmmm, same reference
				}, nil
			}
			return nil, fmt.Errorf("no pssh or pro key found in playready element")
		}
	}
	return nil, fmt.Errorf("not a proper content encryption")
}

func (a *AdaptationSet) GetAudioChannels() (int, error) {
	if chConf := a.AudioChannelConfiguration; chConf != nil {
		if Ch := chConf.Value; Ch != nil {
			ch, err := strconv.Atoi(*Ch)
			if err != nil {
				return -1, fmt.Errorf("invalid audio channels value %w", err)
			}
			return ch, nil
		}
	}

	return -1, fmt.Errorf("no audio channels found")
}

func (r *Representation) GetBandwith() (uint64, error) {
	// Only to video and audio, original bandwith in B/s
	if r.Bandwith != nil {
		return *r.Bandwith / 1024, nil //KiB/s - Kbps
	}

	return 0, fmt.Errorf("no bandwith in the representation")
}

func (a *AdaptationSet) GetLang() (string, error) {
	// Mostly for audio and subtitles only, need to check
	if a.Lang != nil {
		return *a.Lang, nil
	}

	for _, rep := range a.Representations {
		if rep.Lang != nil {
			return *rep.Lang, nil
		}
	}

	return "", fmt.Errorf("no lang found in AdaptationSet")
}

func (a *AdaptationSet) GetCodecAll() (string, error) {
	// Getting codec with fallback to codec in Representation
	// Maybe need to switch, dunno currently
	if cod, err := a.GetCodec(); err == nil {
		return cod, nil
	}

	for _, repr := range a.Representations {
		if cod, err := repr.GetCodec(); err == nil {
			return cod, nil
		}
	}

	return "", fmt.Errorf("no codec found in mpd")
}

func (a *AdaptationSet) GetCodec() (string, error) {
	// Mostly for audio and subtitles AdaptationSet
	if a.Codecs != nil {
		return codecRepr(*a.Codecs), nil
	}
	return "", fmt.Errorf("no codec found for the AdaptationSet")
}

func (r *Representation) GetCodec() (string, error) {
	// Mostly to video Representations have to check more examples
	if r.Codecs != nil {
		return codecRepr(*r.Codecs), nil
	}
	return "", fmt.Errorf("no codec found for the representation")
}

func codecRepr(codec string) string {
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

func (a *AdaptationSet) GetResolutionRepr() (string, error) {
	// MAYBE, if movie, have to check
	if a.Width != nil && a.Height != nil {
		return fmt.Sprintf("%dx%d", *a.Width, *a.Height), nil
	}
	return "", fmt.Errorf("no available resolution info")
}

func (r *Representation) GetResolutionRepr() (string, error) {
	// Only to video Representation, main
	if r.Width != nil && r.Height != nil {
		return fmt.Sprintf("%dx%d", *r.Width, *r.Height), nil
	}
	return "", fmt.Errorf("no available resolution info")
}

func (a *AdaptationSet) GetType() (enums.MediaType, error) {
	if a.MimeType != nil {
		s := strings.Split(*a.MimeType, "/")[0]
		return enums.ParseMediaType(s)
	}

	if a.ContentType != nil {
		return enums.ParseMediaType(*a.ContentType)
	}

	for _, rep := range a.Representations {
		if rep.MimeType != nil {
			s := strings.Split(*rep.MimeType, "/")[0]
			return enums.ParseMediaType(s)
		}
		if rep.ContentType != nil {
			return enums.ParseMediaType(*rep.ContentType)
		}
	}

	return -1, fmt.Errorf("could not determine media type")
}

func (m *Mpd) IsLive() (bool, error) {
	if strings.EqualFold(*m.Type, "static") {
		return false, nil
	}
	if strings.EqualFold(*m.Type, "dynamic") {
		return true, nil
	}

	return false, fmt.Errorf("%v is not static or dynamic", *m.Type)
}

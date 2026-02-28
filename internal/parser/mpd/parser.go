package mpd

import (
	"encoding/xml"
	"log"
	"strings"
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

type Period struct {
	Duration       *string          `xml:"duration,attr"`
	Id             *string          `xml:"id,attr"`
	BaseUrl        *string          `xml:"BaseURL,omitempty"` // combine with Mpd BaseUrl
	AdaptationSets []*AdaptationSet `xml:"AdaptationSet,omitempty"`
}

type AdaptationSet struct {
	// Skipped (not exclusive): contentType, frameRate
	MimeType  *string `xml:"mimeType,attr"` // ex. audio/mp4
	FrameRate *string `xml:"frameRate,attr"`
	BaseUrl   *string `xml:"BaseURL,omitempty"` // combine with Period, Mpd BaseUrl
	// Representations    []*Representation    `xml:"Representation,omitempty"`
	// ContentProtections []*ContentProtection `xml:"ContentProtection,omitempty"`
}

func (m *Mpd) IsLive() bool {
	if strings.EqualFold(*m.Type, "static") {
		return false
	} else if strings.EqualFold(*m.Type, "dynamic") {
		return true
	} else {
		log.Printf("%v is not static or dynamic\n", *m.Type)
		return false
	}

}

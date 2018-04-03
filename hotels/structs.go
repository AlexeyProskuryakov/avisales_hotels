package hotels

import (
	"fmt"
	"crypto/md5"
)

type Location struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Address     string  `json:"address"`
}

type Hotel struct {
	HotelId     string
	SourceId    string
	SourceType  string
	Name        string
	Photos      []string
	Description string
	Location    Location
	Stars       uint64
}

func NewHotel(sourceType string) *Hotel {
	return &Hotel{SourceType: sourceType}
}

func (h *Hotel) SetId(SourceId string) {
	h.SourceId = SourceId
	h.HotelId = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v_%v", h.SourceType, h.SourceId))))
}

package parsers

import (
	"encoding/xml"
	"io"
	"fmt"
	"strings"
	"errors"
	"strconv"
	"regexp"
	"hotels/hotels"
	"hotels/cfg"
)

var descrRegex = regexp.MustCompile(`\[\/?[a-z-]+\]`)

func prepareDescription(description string) string {
	return descrRegex.ReplaceAllString(description, "")
}

func getValueByPath(path string, data map[string]interface{}) (*interface{}) {
	if val, ok := data[path]; ok == true {
		result := val
		return &result
	} else {
		return nil
	}
}

type XMLParser struct {
	FileParseEngine
	Config cfg.ParserConfig
}

func (p *XMLParser) generateHotel(data map[string]interface{}) (*hotels.Hotel, error) {
	// in data {path:data},  in fieldPaths {field:path}
	config := p.Config.FieldsMapping
	result := hotels.NewHotel(cfg.XML)

	name := getValueByPath(config.Name, data)
	if name != nil {
		result.Name = (*name).(string)
	}

	id := getValueByPath(config.Id, data)
	if id != nil {
		result.SetId((*id).(string))
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise id [%s] path", config.Id))
	}

	description := getValueByPath(config.Description, data)
	if description != nil {
		result.Description = prepareDescription((*description).(string))
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise description [%s] path", config.Description))
	}

	stars := getValueByPath(config.Stars.Key, data)
	if stars != nil {
		starsString := (*stars).(string)
		starsUint, err := strconv.ParseUint(starsString, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Can not parse stars config [%s] at [%s]", config.Stars.Key, starsString))
		}
		result.Stars = starsUint
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise stars [%s] path", config.Stars.Key))
	}

	lat, lon := getValueByPath(config.Location.Lat, data), getValueByPath(config.Location.Lon, data)
	if lat != nil && lon != nil {
		latStr, lonStr := (*lat).(string), (*lon).(string)
		latFl, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Can not parse location lat [%s] at [%s]", config.Location.Lat, latStr))
		}
		lonFl, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Can not parse location lon [%s] at [%s]", config.Location.Lat, latStr))
		}
		result.Location.Lat = latFl
		result.Location.Lon = lonFl
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise location lat/lon [%s]/[%s] path", config.Location.Lat, config.Location.Lon))
	}

	country := getValueByPath(config.Location.Country, data)
	if country != nil {
		result.Location.Country = (*country).(string)
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise location country [%s] path", config.Stars.Key))
	}

	countryCode := getValueByPath(config.Location.CountryCode, data)
	if country != nil {
		result.Location.CountryCode = (*countryCode).(string)
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise location country code [%s] path", config.Stars.Key))
	}

	city := getValueByPath(config.Location.City, data)
	if city != nil {
		result.Location.City = (*city).(string)
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise location city [%s] path", config.Stars.Key))
	}

	address := getValueByPath(config.Location.Address, data)
	if city != nil {
		result.Location.Address = (*address).(string)
	} else {
		return nil, errors.New(fmt.Sprintf("Can not recognise location address [%s] path", config.Stars.Key))
	}

	return result, nil

}

func (p XMLParser) ParseFileData(r io.Reader, hotelChan chan *hotels.Hotel) error {
	d := xml.NewDecoder(r)

	var parseHotel, parsePhotos bool
	var path, currentTag string

	elementData := map[string]interface{}{}
	photosData := []string{}

	for {
		t, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr != io.EOF {
				fmt.Printf("Error at reading xml file %s", tokenErr)
				break
			}
			return tokenErr
		}

		switch t := t.(type) {
		case xml.StartElement:
			currentTag = t.Name.Local
			if currentTag == p.Config.In {
				parseHotel = true
				continue
			}

			if parseHotel == true && currentTag != p.Config.In {
				path += fmt.Sprintf("/%s", currentTag)
			}

			if currentTag == p.Config.FieldsMapping.Photos.In {
				parsePhotos = true
			}

		case xml.EndElement:
			currentTag = t.Name.Local
			if currentTag == p.Config.In {
				parseHotel = false
				hotel, err := p.generateHotel(elementData)
				if err != nil {
					fmt.Printf("Error at parsing xml element: %s %s", t, err)
					continue
				}
				hotel.Photos = photosData

				fmt.Printf("XML Processing hotel %v %v \n", hotel.SourceId, hotel.Name)
				hotelChan <- hotel

				elementData = map[string]interface{}{}
				photosData = []string{}
				path = ""
				currentTag = ""
			}

			if parseHotel == true && currentTag != p.Config.In {
				path = strings.TrimSuffix(path, fmt.Sprintf("/%s", currentTag))
			}

			if currentTag == p.Config.FieldsMapping.Photos.In {
				parsePhotos = false
			}

		case xml.CharData:
			tagData := strings.TrimSpace(string([]byte(t)))
			if parseHotel == true && currentTag != p.Config.In {
				elementData[path] = tagData
			}
			if parsePhotos == true && currentTag == p.Config.FieldsMapping.Photos.Key {
				if tagData != "" {
					photosData = append(photosData, tagData)
				}
			}
		}
	}
	return nil
}

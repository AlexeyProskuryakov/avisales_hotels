package parsers

import (
	"hotels/cfg"
	"hotels/hotels"
	"io"
	"bufio"
	"fmt"
	"encoding/json"
	"github.com/oliveagle/jsonpath"
)

type JsonParser struct {
	FileParseEngine
	Config cfg.ParserConfig
}

func retrievePhotos(jsonData interface{}, in, key string) ([]string, error) {
	result := []string{}
	images, err := jsonpath.JsonPathLookup(jsonData, in)
	if err != nil {
		return nil, err
	}
	for _, image := range images.([]interface{}) {
		imageInfo := image.(map[string]interface{})
		if data, ok := imageInfo[key]; ok == true {
			result = append(result, data.(string))
		}
	}
	return result, nil
}

func (p *JsonParser) parseJsonElement(jsonData interface{}) (*hotels.Hotel, error) {
	hotel := hotels.NewHotel(cfg.JSON)

	id, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Id)
	if err != nil {
		fmt.Printf("Can not read id property from json, because: %v\n", err)
		return nil, err
	}
	hotel.SetId((id.(string)))

	name, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Name)
	if err != nil {
		fmt.Printf("Can not read name property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Name = (name.(string))

	descr, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Description)
	if err != nil {
		fmt.Printf("Can not read name description from json, because: %v\n", err)
		return nil, err
	}
	hotel.Description = (descr.(string))

	starsValue, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Stars.Key)
	if err != nil {
		fmt.Printf("Can not read stars property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Stars = uint64(p.Config.FieldsMapping.Stars.FormValue(starsValue.(float64)))

	photosCfg := p.Config.FieldsMapping.Photos
	photos, err := retrievePhotos(jsonData, photosCfg.In, photosCfg.Key)
	if err != nil {
		fmt.Printf("Can not read photos property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Photos = photos

	address, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.Address)
	if err != nil {
		fmt.Printf("Can not read address property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.Address = (address.(string))

	city, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.City)
	if err != nil {
		fmt.Printf("Can not read city property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.City = (city.(string))

	country, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.Country)
	if err != nil {
		fmt.Printf("Can not read name country from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.Country = (country.(string))

	countryCode, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.CountryCode)
	if err != nil {
		fmt.Printf("Can not read name countryCode from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.CountryCode = (countryCode.(string))

	lat, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.Lat)
	if err != nil {
		fmt.Printf("Can not read lat property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.Lat = (lat.(float64))

	lon, err := jsonpath.JsonPathLookup(jsonData, p.Config.FieldsMapping.Location.Lon)
	if err != nil {
		fmt.Printf("Can not read lon property from json, because: %v\n", err)
		return nil, err
	}
	hotel.Location.Lon = (lon.(float64))

	return hotel, nil
}

func (p JsonParser) ParseFileData(file io.Reader, hotelChan chan *hotels.Hotel) error {
	reader := bufio.NewReader(file)
	var line string
	var err error
	counter := 0
	for {
		counter += 1
		line, err = reader.ReadString('\n')
		if line != "" {
			var jsonData interface{}
			unmarshalErr := json.Unmarshal([]byte(line), &jsonData)
			if unmarshalErr != nil {
				fmt.Printf("Error at unmarshalling: [%v] \njson data: [%v]\n", err, line)
				continue
			}
			hotel, err := p.parseJsonElement(jsonData)
			if err != nil {
				fmt.Printf("Error at parsing json element [%v]\n", err)
			}

			fmt.Printf("JSON Processing hotel %v %v %v \n", counter, hotel.SourceId, hotel.Name)
			hotelChan <- hotel
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error at reading json file data [%v]\n", err)
				return err
			}
		}

	}
	return nil
}

package parsers

import (
	"hotels/cfg"
	"io"
	"hotels/hotels"
	"bufio"
	"encoding/csv"
	"fmt"
	"errors"
	"strconv"
	"regexp"
)

func loadNamesToPositionMapping(input []string) map[string]int {
	result := make(map[string]int, len(input))
	for i, el := range input {
		result[el] = i
	}
	return result
}

type CSVParser struct {
	FileParseEngine
	Config cfg.ParserConfig
}

func getDataByName(name string, record []string, mapping map[string]int) (string, error) {
	if position, ok := mapping[name]; ok {
		return record[position], nil
	} else {
		return "", errors.New(fmt.Sprintf("Can not found [%v] field in csv file", name))
	}
}

func (p *CSVParser) parseCSV(record []string, mapping map[string]int) (*hotels.Hotel, error) {
	hotel := hotels.NewHotel(cfg.CSV)

	fieldVal, err := getDataByName(p.Config.FieldsMapping.Id, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.SetId(fieldVal)

	nameVal, err := getDataByName(p.Config.FieldsMapping.Name, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Name = nameVal

	description, err := getDataByName(p.Config.FieldsMapping.Description, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Description = description

	stars, err := getDataByName(p.Config.FieldsMapping.Stars.Key, record, mapping)
	if err != nil {
		return nil, err
	}
	starsInt, err := strconv.ParseUint(stars, 10, 64)
	if err != nil {
		return nil, errors.New("Can't recognise stars value of csv file")
	}
	hotel.Stars = starsInt

	address, err := getDataByName(p.Config.FieldsMapping.Location.Address, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Location.Address = address

	city, err := getDataByName(p.Config.FieldsMapping.Location.City, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Location.City = city

	country, err := getDataByName(p.Config.FieldsMapping.Location.Country, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Location.Country = country

	countryCode, err := getDataByName(p.Config.FieldsMapping.Location.CountryCode, record, mapping)
	if err != nil {
		return nil, err
	}
	hotel.Location.CountryCode = countryCode

	lat, err := getDataByName(p.Config.FieldsMapping.Location.Lat, record, mapping)
	if err != nil {
		return nil, err
	}
	latFl, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, errors.New("Can't recognise latitude value of csv file")
	}
	hotel.Location.Lat = latFl

	lon, err := getDataByName(p.Config.FieldsMapping.Location.Lon, record, mapping)
	if err != nil {
		return nil, err
	}
	lonFl, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return nil, errors.New("Can't recognise latitude value of csv file")
	}
	hotel.Location.Lon = lonFl

	photos := []string{}
	for name, position := range mapping {
		if matched, err := regexp.MatchString(p.Config.FieldsMapping.Photos.Key, name); matched && err == nil {
			photos = append(photos, record[position])
		} else if err != nil {
			return nil, errors.New(fmt.Sprintf("Bad regexp in photos.key for csv file"))
		}
	}

	hotel.Photos = photos
	return hotel, nil
}

func (p CSVParser) ParseFileData(file io.Reader, hotelChan chan *hotels.Hotel) error {
	reader := csv.NewReader(bufio.NewReader(file))
	headerLoaded := false
	var namesToPositionMapping map[string]int
	counter := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error at reading csv file data [%v]\n", err)
				return err
			}
		}

		if headerLoaded == false {
			namesToPositionMapping = loadNamesToPositionMapping(record)
			headerLoaded = true
			continue
		}

		hotel, err := p.parseCSV(record, namesToPositionMapping)
		if err != nil {
			return err
		}
		fmt.Printf("CSV Processing hotel %s %s \n", counter, hotel.SourceId, hotel.Name)
		hotelChan <- hotel
	}
	return nil
}

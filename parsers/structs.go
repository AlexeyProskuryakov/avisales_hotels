package parsers

import (
	"hotels/hotels"
	"os"
	"io"
	"fmt"
)

type FileParseEngine interface {
	ParseFileData(file io.Reader, hotelChan chan *hotels.Hotel) error
}

type FileParser struct {
	FileParseEngine
}

func (p *FileParser) Parse(fileName string) (<-chan *hotels.Hotel, error) {
	file, err := os.Open(fileName)
	if (err != nil) {
		return nil, err
	}
	hotelChan := make(chan *hotels.Hotel)

	if (file != nil) {
		go func() {
			err := p.ParseFileData(file, hotelChan)
			if err != nil {
				fmt.Printf("Error at parsing file data %v\n", err)
			}

			file.Close()
			close(hotelChan)
		}()
	}
	return hotelChan, nil
}

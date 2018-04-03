package main

import (
	"fmt"
	"hotels/cfg"
	"hotels/db"
	"hotels/parsers"
	"os"
	"strings"
)

func main() {
	config, err := cfg.LoadConfig()
	if (err != nil) {
		fmt.Printf("Error at loading config file")
		return
	}
	hotelDb := db.NewHotelsDb(config.Db)

	fileName := os.Args[1]
	pPosition := strings.LastIndex(fileName, ".")
	var fileExt string
	if pPosition != -1 {
		fileExt = fileName[pPosition+1:]
	} else {
		fmt.Printf("I recognise type of file using extension. Please provide json|xml|csv extension for your file\n")
		return
	}

	var engine parsers.FileParseEngine
	switch fileExt {
	case cfg.CSV:
		engine = parsers.CSVParser{Config: config.Parsers.CSV}
	case cfg.XML:
		engine = parsers.XMLParser{Config: config.Parsers.XML}
	case cfg.JSON:
		engine = parsers.JsonParser{Config: config.Parsers.JSON}
	default:
		fmt.Printf("I haven't parser for [%v] extension", fileExt)
		return
	}

	parser := parsers.FileParser{FileParseEngine: engine}

	hotelsChan, err := parser.Parse(fileName)
	if err != nil {
		fmt.Printf("Error at parsing xml file %v", err)
	}

	hotelDb.StoreHotels(hotelsChan)
}

package cfg

import (
	"reflect"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const (
	XML  = "xml"
	JSON = "json"
	CSV  = "csv"

	_div = "/"
	_pow = "*"
)

type CountProcessField struct {
	Key      string  `json:"key"`
	CalcType string  `json:"calc_type"`
	Factor   float64 `json:"factor"`
}

func (cpf *CountProcessField) FormValue(input float64) float64 {
	switch cpf.CalcType {
	case _div:
		return input / cpf.Factor
	case _pow:
		return input * cpf.Factor
	default:
		return input
	}
}

type LocationConfig struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Address     string `json:"address"`
}

type PhotosField struct {
	In  string `json:"in"`
	Key string `json:"key"`
}

type FieldsMappingConfig struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Photos      PhotosField    `json:"photos"`
	Description string         `json:"description"`
	Location    LocationConfig `json:"location"`
	Stars       CountProcessField
}

type ParserConfig struct {
	FieldsMapping FieldsMappingConfig `json:"fields_mapping"`
	In            string              `json:"in"`
}
type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Db     DbConfig `json:"db"`
	Locale string   `json:"locale"`
	Parsers struct {
		JSON ParserConfig `json:"json"`
		CSV  ParserConfig `json:"csv"`
		XML  ParserConfig `json:"xml"`
	} `json:"parsers"`
}

func changeLocaleMarker(obj FieldsMappingConfig, localeString string) FieldsMappingConfig {
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	changeLocaleRec(copy, original, localeString)
	return copy.Interface().(FieldsMappingConfig)
}

func changeLocaleRec(copy reflect.Value, original reflect.Value, change string) {
	switch original.Kind() {
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			changeLocaleRec(copy.Field(i), original.Field(i), change)
		}
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			changeLocaleRec(copy.Index(i), original.Index(i), change)
		}

	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			changeLocaleRec(copyValue, originalValue, change)
			copy.SetMapIndex(key, copyValue)
		}

	case reflect.String:
		translatedString := original.Interface().(string)
		translatedString = strings.Replace(translatedString, "%loc", change, -1)
		copy.SetString(translatedString)
	default:
		copy.Set(original)
	}
}

func LoadConfig() (*Config, error) {
	var result Config
	data, err := ioutil.ReadFile("config.json")
	if (err != nil) {
		fmt.Printf("Error at reading config file %v", err)
		return nil, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Printf("Error at read parse config %v", err)
		return nil, err
	}
	result.Parsers.XML.FieldsMapping = changeLocaleMarker(result.Parsers.XML.FieldsMapping, result.Locale)
	result.Parsers.JSON.FieldsMapping = changeLocaleMarker(result.Parsers.JSON.FieldsMapping, result.Locale)
	result.Parsers.CSV.FieldsMapping = changeLocaleMarker(result.Parsers.CSV.FieldsMapping, result.Locale)
	return &result, nil

}

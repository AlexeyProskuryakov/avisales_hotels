package db

import (
	"hotels/hotels"
	"fmt"
	_ "github.com/lib/pq"
	"hotels/cfg"
	"database/sql"
	"errors"
)

type HotelsDB struct {
	info string
	db   *sql.DB
}

const (
	hotel_insert = `INSERT INTO hotels (hotel_id, source, source_id, name, description, stars, country_code, country, city, address, lat, lon)  VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	photo_insert = `INSERT INTO photos (hotel_id, url) VALUES ($1, $2)`
)

func (p *HotelsDB) storeHotel(hotel *hotels.Hotel) error {
	tx, err := p.db.Begin()
	if err != nil {
		fmt.Printf("Error at begin transaction %v\n", err)
		return err
	}

	{
		stmt, err := tx.Prepare(hotel_insert)
		if err != nil {
			fmt.Printf("Error at prepare insert hotel statement %v\n", err)
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(hotel.HotelId, hotel.SourceType, hotel.SourceId, hotel.Name, hotel.Description, hotel.Stars, hotel.Location.CountryCode, hotel.Location.Country, hotel.Location.City, hotel.Location.Address, hotel.Location.Lat, hotel.Location.Lon); err != nil {
			fmt.Printf("Error at execute insert hotel statement %v\n", err)
			tx.Rollback()
			return err
		}

	}

	{
		for _, photo := range hotel.Photos {
			stmt, err := tx.Prepare(photo_insert)
			if err != nil {
				fmt.Printf("Error at prepare insert photo statement %v\n", err)
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(hotel.HotelId, photo); err != nil {
				fmt.Printf("Error at execute insert photo statement %v\n", err)
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (p *HotelsDB) processHotelsChan(hotelChan <-chan *hotels.Hotel) <-chan bool {
	done := make(chan bool)

	go func() {
		for {
			hotel, more := <-hotelChan
			if !more {
				fmt.Printf("No more hotels. Stopping persisting.\n")
				done <- true
				break
			}
			errResult := p.storeHotel(hotel)
			if errResult != nil {
				fmt.Printf("Can not store hotel %v %v, because %v.\n", hotel.SourceType, hotel.SourceId, errResult)
			}

		}
	}()
	return done
}

func NewHotelsDb(config cfg.DbConfig) *HotelsDB {
	result := HotelsDB{}

	pgInfo := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`, config.Host, config.Port, config.User, config.Password, config.DbName)
	result.info = pgInfo

	return &result
}

func (p *HotelsDB) openConnection() error {
	db, err := sql.Open("postgres", p.info)
	if err != nil {
		fmt.Sprintf("Error at open connection %v\n", p.info)
		return err
	}
	p.db = db
	return nil
}

func (p *HotelsDB) closeConnection() error {
	if (p.db == nil) {
		return errors.New("Connection not exists")
	}
	return p.db.Close()
}

func (p *HotelsDB) StoreHotels(hotelsChan <-chan *hotels.Hotel) {
	p.openConnection()
	defer p.closeConnection()

	result := p.processHotelsChan(hotelsChan)
	<-result
}

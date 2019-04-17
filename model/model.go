package model

import (
	"database/sql"
	"fmt"
	"pik-arenda/base"
	"sync"

	"github.com/lib/pq"
)

// Database head struct pg
type Database struct {
	db *sql.DB
}

var (
	_once     sync.Once
	_database *Database
)

// GetDatabase get singleton
func GetDatabase() (*Database, error) {
	var err error

	_once.Do(func() {
		_database = new(Database)
		_database, err = newDatabase()
	})
	return _database, err
}

func newDatabase() (*Database, error) {
	db := new(Database)
	if err := db.connectPSQL(); err != nil {
		return db, fmt.Errorf("[db CONNECT] %v", err)
	}
	return db, nil
}

func (d *Database) connectPSQL() (err error) {
	d.db, err = sql.Open("postgres", base.GetConfig().ConnStr)
	if err != nil {
		return err
	}
	return d.db.Ping()
}

// InputFlat add flat into database
func (d *Database) InputFlat(houseID, roomNumber, floor, roomCount int, cost, space float64) (id int, err error) {
	err = d.db.QueryRow(`insert into flat (house_id, number, cost, space, floor, room_count) values($1,$2,$3,$4,$5,$6) RETURNING id`,
		houseID,
		roomNumber,
		cost,
		space,
		floor,
		roomCount).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, err
}

// SelectInfoByName return id dictionary by name
func (d *Database) SelectInfoByName(table, value string) (id int, err error) {

	tableQuoted := pq.QuoteIdentifier(table)
	err = d.db.QueryRow(fmt.Sprintf("SELECT id from %s where name = $1", tableQuoted), value).Scan(&id)

	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	return id, nil
}

// SelectHouseID return house id
func (d *Database) SelectHouseID(cityID, districtID, streetID, number int, literal string) (id int, err error) {
	err = d.db.QueryRow(`select id 
from house 
where city_id = $1 
  and district_id = $2 
  and street_id = $3 
  and number = $4 
  and literal = $5 `, cityID, districtID, streetID, number, literal).Scan(&id)

	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	return id, nil
}

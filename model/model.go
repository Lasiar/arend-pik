package model

import (
	"database/sql"
	"fmt"
	"pik-arenda/base"
	"sync"

	_ "github.com/lib/pq"
)

// Flat flat struct
type Flat struct {
	FlatDataPublic
	flatDataPrivate
}

type FlatDataPublic struct {
	ID          int     `json:"id,omitempty"`
	City        string  `json:"city,omitempty"`
	District    string  `json:"district,omitempty"`
	Street      string  `json:"street,omitempty"`
	HouseNumber int     `json:"house,omitempty"`
	Literal     string  `json:"literal,omitempty"`
	FloorCount  int     `json:"floor_count,omitempty"`
	Floor       int     `json:"floor,omitempty"`
	RoomCount   int     `json:"room_count,omitempty"`
	RoomNumber  int     `json:"room_number,omitempty"`
	Space       float64 `json:"space,omitempty"`
	Cost        float64 `json:"cost,omitempty"`
}

type flatDataPrivate struct {
	literal sql.NullString
}

func (fd *Flat) validate() {
	if fd.flatDataPrivate.literal.Valid {
		fd.FlatDataPublic.Literal = fd.flatDataPrivate.literal.String
	}
}

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

package model

import (
	"database/sql"
	"fmt"
)

// FlatData struc for save info
type FlatData struct {
	flatDataPublic
	flatDataPrivate
}

type flatDataPublic struct {
	ID         int     `json:"id"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	Street     string  `json:"street"`
	House      int     `json:"house"`
	FloorCount int     `json:"floor_count"`
	Floor      int     `json:"floor"`
	RoomCount  int     `json:"room_count"`
	RoomNumber int     `json:"room_number"`
	Literal    string  `json:"literal"`
	Space      float64 `json:"space"`
	Cost       float64 `json:"cost"`
}

type flatDataPrivate struct {
	literal sql.NullString
}

func (fd *FlatData) validate() {
	if fd.flatDataPrivate.literal.Valid {
		fd.flatDataPublic.Literal = fd.flatDataPrivate.literal.String
	}
}

// EventDataSelectBuilder структура для удобной работы с запросами по выборке событий
type EventDataSelectBuilder struct {
	db         *sql.DB
	MaxCount   int
	Offset     int
	Reverse    bool
	baseWhere  string
	baseFrom   string
	baseSelect string
	rawSelect  string
	City       string
	Street     string
	District   string
	Space      float64
	Cost       float64
}

// NewSelectBuilder зоздает новый экземпляр SelectBuilder
func (d *Database) NewSelectBuilder() *EventDataSelectBuilder {
	sb := &EventDataSelectBuilder{db: d.db,
		baseWhere: ` WHERE flat.house_id = house.id
	and house.city_id = city.id
	and house.district_id = district.id
	and house.street_id = street.id`,
		baseFrom: ` FROM
	flat,
	house,
	district,
	street,
	city`,
		baseSelect: ``,
		rawSelect: `SELECT
	flat.id as "id",
	flat."cost",
	flat."space",
	flat.room_count,
	flat.floor,
	flat."number",
	street."name",
	city."name",
	district."name"
		`,
	}
	sb.MaxCount = 100
	return sb
}

func (sb *EventDataSelectBuilder) makeOrder(args *[]interface{}) string {
	query := ""
	if sb.MaxCount > 0 {
		query += " LIMIT " + sb.addArg(args, sb.MaxCount)
	}
	if sb.Offset > 0 {
		query += " OFFSET " + sb.addArg(args, sb.Offset)
	}
	return query
}

func (sb *EventDataSelectBuilder) makeWhere(args *[]interface{}) string {
	query := sb.baseWhere
	if sb.City != "" {
		query += " AND city.name = " + sb.addArg(args, sb.City)
	}

	if sb.Cost != 0 {
		query += " AND flat.cost = " + sb.addArg(args, sb.Cost)
	}

	if sb.Street != "" {
		query += " AND street.id = " + sb.addArg(args, sb.Street)
	}

	if sb.District != "" {
		query += " AND district.name = " + sb.addArg(args, sb.District)
	}

	if sb.Space != 0 {
		query += " AND flat.space = " + sb.addArg(args, sb.Space)
	}

	return query
}

func (sb *EventDataSelectBuilder) addArg(args *[]interface{}, arg interface{}) string {
	*args = append(*args, arg)
	return fmt.Sprintf("$%d", len(*args))
}

// Select select
func (sb *EventDataSelectBuilder) Select() ([]*FlatData, error) {
	args := make([]interface{}, 0)
	res := make([]*FlatData, 0)

	query := sb.rawSelect + sb.baseFrom + sb.makeWhere(&args) + sb.makeOrder(&args)

	fmt.Println(query)
	rows, err := sb.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		fd := &FlatData{}
		err = rows.Scan(&fd.ID,
			&fd.Cost,
			&fd.Space,
			&fd.RoomCount,
			&fd.Floor,
			&fd.RoomNumber,
			&fd.Street,
			&fd.City,
			&fd.District)

		if err != nil {
			return nil, err
		}
		fd.validate()
		res = append(res, fd)
	}

	return res, nil
}

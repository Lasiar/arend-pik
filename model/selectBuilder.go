package model

import (
	"database/sql"
	"fmt"
)

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
	FlatDataPublic
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

// TODO: cyclomatic complexity 13  is high
func (sb *EventDataSelectBuilder) makeWhere(args *[]interface{}) string {
	query := sb.baseWhere
	if sb.ID != 0 {
		query += ` AND flat.id = ` + sb.addArg(args, sb.ID)
	}

	if sb.City != "" {
		query += ` AND city."name" = ` + sb.addArg(args, sb.City)
	}

	if sb.District != "" {
		query += ` AND district."name" =` + sb.addArg(args, sb.District)
	}

	if sb.Street != "" {
		query += ` AND street."name" = ` + sb.addArg(args, sb.Street)
	}

	if sb.HouseNumber != 0 {
		query += ` AND  house."number" = ` + sb.addArg(args, sb.HouseNumber)
	}

	if sb.Literal != "" {
		query += ` AND house.literal = ` + sb.addArg(args, sb.Literal)

	}

	if sb.FloorCount != 0 {
		query += ` AND house.floor_count = ` + sb.addArg(args, sb.FloorCount)
	}

	if sb.Floor != 0 {
		query += ` AND flat.floor = ` + sb.addArg(args, sb.Floor)
	}

	if sb.RoomCount != 0 {
		query += `  AND flat.room_count = ` + sb.addArg(args, sb.RoomCount)
	}

	if sb.RoomNumber != 0 {
		query += ` AND flat."number" = ` + sb.addArg(args, sb.RoomNumber)
	}

	if sb.Space != 0 {
		query += ` AND flat.space = ` + sb.addArg(args, sb.Space)
	}

	if sb.Cost != 0 {
		query += ` AND flat.cost = ` + sb.addArg(args, sb.Cost)
	}

	return query
}

func (sb *EventDataSelectBuilder) addArg(args *[]interface{}, arg interface{}) string {
	*args = append(*args, arg)
	return fmt.Sprintf("$%d", len(*args))
}

// Select select
func (sb *EventDataSelectBuilder) Select() ([]*Flat, error) {
	args := make([]interface{}, 0)
	res := make([]*Flat, 0)

	query := sb.rawSelect + sb.baseFrom + sb.makeWhere(&args) + sb.makeOrder(&args)

	fmt.Println(query)
	rows, err := sb.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		fd := &Flat{}
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

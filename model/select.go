package model

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

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

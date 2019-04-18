package model

// InsertFlat add flat into database
func (d *Database) InsertFlat(houseID, roomNumber, floor, roomCount int, cost, space float64) (id int, err error) {
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

package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pik-arenda/model"
)

func add() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		flat := new(model.Flat)

		if err := json.NewDecoder(r.Body).Decode(&flat); err != nil {
			SetError(r, err)
			return
		}

		db, err := model.GetDatabase()
		if err != nil {
			SetError(r, err)
			return
		}

		cityID, err := checker(db, "city", flat.City)
		if err != nil {
			SetError(r, err)
			return
		}

		districtID, err := checker(db, "district", flat.District)
		if err != nil {
			SetError(r, err)
			return
		}

		streetID, err := checker(db, "street", flat.Street)
		if err != nil {
			SetError(r, err)
			return
		}

		houseID, err := db.SelectHouseID(cityID, districtID, streetID, flat.HouseNumber, flat.Literal)
		if err != nil {
			SetError(r, fmt.Errorf("select house id: %v", err))
			return
		}
		if houseID < 1 {
			SetError(r, fmt.Errorf("hose not found"))
			return
		}

		flatID, err := db.InsertFlat(houseID, flat.RoomNumber, flat.Floor, flat.RoomCount, flat.Cost, flat.Space)
		if err != nil {
			SetError(r, fmt.Errorf("input flat: %v", err))
			return
		}

		resp := struct {
			ID int `json:"id"`
		}{}

		resp.ID = flatID

		SetResponse(r, resp)
	})
}
func checker(db *model.Database, table, value string) (int, error) {
	id, err := db.SelectInfoByName(table, value)
	if err != nil {
		return id, err
	}

	if id < 1 {
		return id, fmt.Errorf("%s not found", table)
	}
	return id, nil
}

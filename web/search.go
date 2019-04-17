package web

import (
	"encoding/json"
	"log"
	"net/http"
	"pik-arenda/model"
)

func search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flat := new(struct {
			Flat
			PageNumber int `json:"page_number"`
			RowsOnPage int `json:"rows_on_page"`
		})

		if err := json.NewDecoder(r.Body).Decode(&flat); err != nil {
			SetError(r, err)
			return
		}

		db, err := model.GetDatabase()
		if err != nil {
			log.Fatal(err)
		}

		sb := db.NewSelectBuilder()

		sb.Space = flat.Space
		sb.District = flat.District
		sb.City = flat.City

		res, err := sb.Select()
		if err != nil {
			log.Fatal(err)
		}

		SetResponse(r, res)
	})
}

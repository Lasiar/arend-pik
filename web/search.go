package web

import (
	"encoding/json"
	"log"
	"net/http"
	"pik-arenda/model"
)

func search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flatSearch := new(struct {
			model.FlatDataPublic
			MaxCount int `json:"max_count"`
			Offset   int `json:"offset"`
		})

		if err := json.NewDecoder(r.Body).Decode(&flatSearch); err != nil {
			SetError(r, err)
			return
		}

		db, err := model.GetDatabase()
		if err != nil {
			log.Fatal(err)
		}

		sb := db.NewSelectBuilder()
		sb.FlatDataPublic = flatSearch.FlatDataPublic
		sb.MaxCount = flatSearch.MaxCount
		sb.Offset = flatSearch.Offset

		res, err := sb.Select()
		if err != nil {
			log.Fatal(err)
		}

		SetResponse(r, res)
	})
}

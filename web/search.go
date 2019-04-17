package web

import (
	"encoding/json"
	"net/http"
)

// SearchFlat search
type SearchFlat struct {
	Flat
	PageNumber int `json:"page_number"`
	RowsOnPage int `json:"rows_on_page"`
}

func search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flat := new(SearchFlat)

		if err := json.NewDecoder(r.Body).Decode(&flat); err != nil {
			SetError(r, err)
			return
		}

	})
}

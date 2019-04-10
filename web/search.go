package web

import (
	"fmt"
	"net/http"
)

func search(w http.ResponseWriter, r *http.Request) *errorHanding {
	fmt.Fprint(w, "search")
	return nil
}

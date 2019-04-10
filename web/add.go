package web

import (
	"fmt"
	"net/http"
)

func add(w http.ResponseWriter, r *http.Request) *errorHanding {
	fmt.Fprint(w, "added")
	return nil
}


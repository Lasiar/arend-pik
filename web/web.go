package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pik-arenda/base"
	"time"
)

// Run start web server
func Run() {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/search", search)
	serveMux.HandleFunc("/add", add)

	apiMux := JSONWriteHandler(serveMux)

	server := &http.Server{
		Addr:           base.GetConfig().Port,
		Handler:        apiMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("start listen on: %s \n", base.GetConfig().Port)
	if err := server.ListenAndServe(); err != nil {
		log.Panicf("error run web: %v", err)
	}
}

// JSONWriteHandler хандлер для ответа в виде json
func JSONWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "POST, OPTIONS")
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}

		if err := r.Context().Err(); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			// switch err {
			// case ErrForbidden:
			//	w.WriteHeader(http.StatusForbidden)
			// case ErrNotFound:
			//	w.WriteHeader(http.StatusNotFound)

			w.WriteHeader(http.StatusInternalServerError)
			//}

			log.Println(err)

			if _, err := io.WriteString(w, err.Error()); err != nil {
				log.Println(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)

		data := r.Context().Value(ResponseDataKey)
		if data == nil {
			return
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	})
}

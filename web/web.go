package web

import (
	"encoding/json"
	"log"
	"net/http"
	"pik-arenda/base"
	"time"
)

type errorHanding struct {
	Error   error
	Message string
}

type webHandler func(http.ResponseWriter, *http.Request) *errorHanding

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Run start web server
func Run() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/search", webHandler(search))
	serveMux.Handle("/add", webHandler(add))

	apiMux := middlewareCORS(serveMux)

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

func (wh webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := wh(w, r)
	if e == nil {
		return
	}

	request := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}
	encoder := json.NewEncoder(w)

	request.Message = e.Message

	log.Printf("[WEB] %v %v [METНOD] %v [URL] %v [USER AGENT] %v", e.Message, e.Error, r.Method, r.URL, r.UserAgent())

	w.WriteHeader(http.StatusInternalServerError)

	if err := encoder.Encode(request); err != nil {
		log.Printf("[WEB] %v", err)
	}
}

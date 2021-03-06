package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/tdewolff/auth"
)

type API struct {
	db     *sqlx.DB
	devURL string
}

func New(db *sqlx.DB) *API {
	return &API{
		db,
		"",
	}
}

func (api *API) SetDevURL(devURL string) {
	api.devURL = devURL
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api.devURL != "" {
		w.Header().Set("Access-Control-Allow-Origin", api.devURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == "OPTIONS" {
		return
	}

	router := httprouter.New()
	router.GET("/api/user", api.GetUser)
	router.ServeHTTP(w, r)
}

func (api *API) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, clients := auth.FromContext(r.Context())

	client := clients["google"]
	_ = client

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

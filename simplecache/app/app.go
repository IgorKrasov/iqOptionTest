package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type setObject struct {
	Key     string      `json:"key"`
	Expired int         `json:"expired"`
	Value   interface{} `json:"value"`
}

type setHObject struct {
	Key     string                 `json:"key"`
	Expired int                    `json:"expired"`
	Value   map[string]interface{} `json:"value"`
}

type App struct {
	cache  *cache
	Router *mux.Router
}

func NewApp() *App {
	return &App{
		cache:  NewCache(time.Duration(1 * time.Second)),
		Router: mux.NewRouter(),
	}
}

func (a *App) Initialize() {
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/get/{key}", a.get).Methods("GET")
	a.Router.HandleFunc("/set", a.set).Methods("POST")
	a.Router.HandleFunc("/keys", a.keys).Methods("GET")
	a.Router.HandleFunc("/unset/{key}", a.unset).Methods("DELETE")
	a.Router.HandleFunc("/rpush", a.rpush).Methods("POST")
	a.Router.HandleFunc("/pop/{key}", a.pop).Methods("GET")
	a.Router.HandleFunc("/lgetall/{key}", a.lgetall).Methods("GET")
	a.Router.HandleFunc("/lget/{key}/{id:[0-9]+}", a.lget).Methods("GET")
	a.Router.HandleFunc("/hset", a.hset).Methods("POST")
	a.Router.HandleFunc("/hgetall/{key}", a.hgetall).Methods("GET")
	a.Router.HandleFunc("/hget/{key}/{dictKey}", a.hget).Methods("GET")
}

func (a *App) Run(addr string) {
	fmt.Println("run server")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) set(w http.ResponseWriter, r *http.Request) {
	var so setObject
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&so); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	a.cache.set(so.Key, so.Value, so.Expired)

	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

func (a *App) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	object, err := a.cache.get(key)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}
	respondWithJSON(w, http.StatusOK, object)
}

func (a *App) keys(w http.ResponseWriter, r *http.Request) {
	keys := a.cache.keys()
	respondWithJSON(w, http.StatusOK, keys)
}

func (a *App) unset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	a.cache.deleteItem(key)

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) rpush(w http.ResponseWriter, r *http.Request) {
	var so setObject
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&so); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	_, err := a.cache.rpush(so.Key, so.Value, so.Expired)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

func (a *App) lgetall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	object, err := a.cache.lgetall(key)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	respondWithJSON(w, http.StatusOK, object)
}

func (a *App) lget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	object, err := a.cache.lget(key, id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	respondWithJSON(w, http.StatusOK, object)
}

func (a *App) pop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	object, err := a.cache.pop(key)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	respondWithJSON(w, http.StatusOK, object)
}

func (a *App) hset(w http.ResponseWriter, r *http.Request) {
	var sho setHObject
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sho); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Println(sho.Value)
	if err := a.cache.hset(sho.Key, sho.Value, sho.Expired); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

func (a *App) hgetall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	object, err := a.cache.hgetall(key)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	respondWithJSON(w, http.StatusOK, object)
}

func (a *App) hget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	dictKey := vars["dictKey"]

	object, err := a.cache.hget(key, dictKey)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	respondWithJSON(w, http.StatusOK, object)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

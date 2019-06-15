package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type User struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

var globalS *mgo.Session

const (
	MGODB      = "test"
	COLLECTION = "user"
)

func init() {
	var url = os.Getenv("MONGO_DB_URL")
	s, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("Create Session: %s\n", err)
	}
	globalS = s
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", findUser).Methods("GET")
	fmt.Println("starting user service on port 5000")
	http.ListenAndServe(":5000", nil)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user.Id = bson.NewObjectId()
	if err := insertUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	responseWithJson(w, http.StatusCreated, user)
}

func findUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var users []User
	users, err := findAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	responseWithJson(w, http.StatusOK, users)
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func connect() (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(MGODB).C(COLLECTION)
	return s, c
}

func insertUser(users ...interface{}) error {
	ms, c := connect()
	defer ms.Close()
	return c.Insert(users...)
}

func findAll() ([]User, error) {
	var result []User
	ms, c := connect()
	defer ms.Close()
	err := c.Find(nil).Select(nil).All(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

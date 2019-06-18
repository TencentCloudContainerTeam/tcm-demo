package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Id    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"user" json:"name"`
	Image string        `json:"image"`
}

const (
	MGODB      = "test"
	COLLECTION = "users"
	EGRESSURL  = "http://httpbin.org/anything"
	IMAGEURL   = "https://cdn1.iconfinder.com/data/icons/DarkGlass_Reworked/128x128/apps/user-3.png"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", findUserByName).Methods("GET")
	fmt.Println("starting user service on port 7000")
	http.ListenAndServe(":7000", router)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	headers := getForwardHeaders(r)

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	result, err := findOneByName(user.Name)
	if err != nil && err != mgo.ErrNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if err == mgo.ErrNotFound {
		user.Id = bson.NewObjectId()
		if err := insertUser(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		result.Image = getImageUrlFromHttpBin(headers)
		responseWithJson(w, http.StatusCreated, result)
		return
	}

	user.Image = getImageUrlFromHttpBin(headers)
	responseWithJson(w, http.StatusCreated, user)
}

func findUserByName(w http.ResponseWriter, r *http.Request) {
	headers := getForwardHeaders(r)
	query := r.URL.Query()
	name := query.Get("name")
	user, err := findOneByName(name)
	if err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	user.Image = getImageUrlFromHttpBin(headers)
	responseWithJson(w, http.StatusOK, user)
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func connect() (*mgo.Session, *mgo.Collection, error) {
	var url = os.Getenv("MONGO_DB_URL")
	s, err := mgo.DialWithTimeout(url, 1*time.Second)
	if err != nil {
		log.Printf("Create Session: %s\n", err)
		return nil, nil, err
	}
	c := s.DB(MGODB).C(COLLECTION)
	return s, c, nil
}

func insertUser(users ...interface{}) error {
	ms, c, err := connect()
	if err != nil {
		return err
	}
	defer ms.Close()
	return c.Insert(users...)
}

func findOneByName(name string) (User, error) {
	var result User
	ms, c, err := connect()
	if err != nil {
		return result, err
	}
	defer ms.Close()
	err = c.Find(bson.M{"user": name}).Select(nil).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func getImageUrlFromHttpBin(headers map[string]string) string {
	image := make(map[string]interface{})
	image["url"] = IMAGEURL
	bytesData, err := json.Marshal(image)
	if err != nil {
		return ""
	}
	httpRequest, err := http.NewRequest("POST", EGRESSURL, bytes.NewReader(bytesData))
	httpRequest.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpRequest.Header.Add(k, v)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return ""
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return ""
	}

	if httpResponse.StatusCode != http.StatusOK {
		return ""
	}

	type dataJson struct {
		ImageUrl string `json:"url"`
	}
	type httpBinRsp struct {
		DataJson dataJson `json:"json"`
	}
	resp := &httpBinRsp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return ""
	}

	return resp.DataJson.ImageUrl
}

func getForwardHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	forwardHeaders := []string{
		"user",
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}

	for _, h := range forwardHeaders {
		if v := r.Header.Get(h); v != "" {
			headers[h] = v
		}
	}

	return headers
}

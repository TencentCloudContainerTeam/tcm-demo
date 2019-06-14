package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Product struct {
	ID    int64
	Title string
	Image string
	Price float64
	Sales int64
	Stock int64
}

func main() {
	http.HandleFunc("/discount", discountController)
	fmt.Println("staring discount service on port 4000")
	http.ListenAndServe(":4000", nil)
}

func discountController(w http.ResponseWriter, r *http.Request) {
	mockIDs := []int64{1, 2, 3, 4, 5, 6}

	products, err := getProducts(mockIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	js, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getProducts(ids []int64) ([]*Product, error) {
	query := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	url := fmt.Sprintf("http://127.0.0.1:5000/products?ids=%s", query)

	p := []*Product{}
	err := getJson(url, &p)
	return p, err
}

var client = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

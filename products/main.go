package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	http.HandleFunc("/products", productsController)
	fmt.Println("staring products service on port 5000")
	http.ListenAndServe(":5000", nil)
}

func productsController(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	v := r.URL.Query()
	query := strings.Split(v.Get("ids"), ",")

	for _, q := range query {
		id, err := strconv.ParseInt(q, 10, 64)
		if err != nil {
			//panic(err)
			continue // ignore
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parameter ids is required"))
		return
	}

	fmt.Printf("querying products of %v\n", ids)

	var products []*Product
	for _, id := range ids {
		product := getProduct(id)
		if product != nil {
			products = append(products, product)
		}
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

var mockDB = map[int64]Product{
	1: Product{
		ID:    1,
		Title: "test1",
		Image: "testddd1",
	},
	2: {
		ID:    2,
		Title: "test2",
		Image: "testddd2",
	},
}

func getProduct(id int64) *Product {
	if p, ok := mockDB[id]; ok {
		return &p
	}
	return nil
}

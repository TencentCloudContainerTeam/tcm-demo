package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Discount struct {
	Banner   string     `json:"banner,omitempty"`
	Products []*Product `json:"products"`
}

type Product struct {
	ID    int64   `json:"id"`
	Title string  `json:"title"`
	Image string  `json:"image"`
	Price float64 `json:"price"`
	Sales int64   `json:"sales"`
	Stock int64   `json:"stock"`
}

func main() {
	http.HandleFunc("/discount", discountController)
	fmt.Println("staring discount service on port 7000")
	http.ListenAndServe(":7000", nil)
}

func discountController(w http.ResponseWriter, r *http.Request) {
	mockIDs := []int64{1, 2, 3, 4, 5, 6}

	headers := getForwardHeaders(r)

	products, err := getProducts(mockIDs, headers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	discount := Discount{
		// Banner:   "v2 Banner todo", // v2
		Products: products,
	}

	js, err := json.Marshal(discount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getProducts(ids []int64, headers map[string]string) ([]*Product, error) {
	query := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	url := fmt.Sprintf("http://products.base.svc.cluster.local:7000/products?ids=%s", query)

	p := []*Product{}
	err := getJson(url, &p, headers)
	return p, err
}

var client = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}, headers map[string]string) error {
	reqest, err := http.NewRequest("GET", url, nil)

	for k, v := range headers {
		reqest.Header.Add(k, v)
	}

	if err != nil {
		panic(err)
	}
	response, err := client.Do(reqest)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
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

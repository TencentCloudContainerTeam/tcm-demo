package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Product struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	Sales     int64   `json:"sales"`
	Stock     int64   `json:"stock"`
	Favorites int64   `json:"favorites,omitempty"`
}

func main() {
	http.HandleFunc("/products", productsController)
	fmt.Println("staring products service on port 7000")
	http.ListenAndServe(":7000", nil)
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
	headers := getForwardHeaders(r)

	products := getProduct(ids, headers)

	js, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getProduct(ids []int64, headers map[string]string) []*Product {
	var products []*Product

	var idsStr []string
	for _, id := range ids {
		if p, ok := mockDB[id]; ok {
			products = append(products, &p)
			idsStr = append(idsStr, strconv.FormatInt(id, 10))
		}
	}

	if len(products) == 0 {
		return products
	}

	var waitgroup sync.WaitGroup
	query := strings.Join(idsStr, ",")

	waitgroup.Add(1)
	go func() {
		result := make(map[int64]int64)
		getJson(fmt.Sprintf("http://sales.base.svc.cluster.local:7000/sales?ids=%s", query), &result, headers)
		for _, p := range products {
			if v, ok := result[p.ID]; ok {
				p.Sales = v
			}
		}
		waitgroup.Done()
	}()

	waitgroup.Add(1)
	go func() {
		result := make(map[int64]int64)
		getJson(fmt.Sprintf("http://stock.base.svc.cluster.local:7000/stock?ids=%s", query), &result, headers)
		for _, p := range products {
			if v, ok := result[p.ID]; ok {
				p.Stock = v
			}
		}
		waitgroup.Done()
	}()

	waitgroup.Add(1)
	go func() {
		result := make(map[int64]int64)
		getJson(fmt.Sprintf("http://favorites.base.svc.cluster.local:7000/favorites?ids=%s", query), &result, headers)
		for _, p := range products {
			fmt.Println(result[p.ID])
			if v, ok := result[p.ID]; ok {
				p.Favorites = v
			}
		}
		waitgroup.Done()
	}()

	waitgroup.Wait()
	return products
}

func getJson(url string, target interface{}, headers map[string]string) error {
	var client = &http.Client{Timeout: 10 * time.Second}
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

var mockDB = map[int64]Product{
	1: Product{
		ID:    1,
		Title: "Stainless Steel Sports Bottle",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/SM744_silver_mu_large.jpg?v=1529788095",
		Price: 8,
	},
	2: {
		ID:    2,
		Title: "Jaeger Socks",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/DSC_2308_large.jpg?v=1550501412",
		Price: 8,
	},
	3: {
		ID:    3,
		Title: "The Ladies envoy Short Sleeve Concert Tee",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/DM130L_greyfrost_model_mu_large.jpg?v=1505771048",
		Price: 15,
	},
	4: {
		ID:    4,
		Title: "Prometheus Full-Zip Hoodie",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/Prometheus_3739_Front_large.jpg?v=1469471991",
		Price: 35,
	},
	5: {
		ID:    5,
		Title: "Fluentd Decal",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/fluentd_decal_Artboard_1_large.png?v=1528318244",
		Price: 0.5,
	},
	6: {
		ID:    6,
		Title: "The Men's OpenTracing Short Sleeve Concert Tee",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/OpenTracing_Mens_Tee_mu_large.jpg?v=1476819035",
		Price: 15,
	},
	7: {
		ID:    7,
		Title: "linkerd Full-Zip Hoodie",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/Hoodie_mu_38b1d323-3e51-421a-b904-4639bec2062a_large.jpg?v=1549202334",
		Price: 35,
	},
	8: {
		ID:    8,
		Title: "The Ladies gRPC Short Sleeve Concert Tee",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/gRPC_ladies_tee_large.jpg?v=1554062191",
		Price: 15,
	},
	9: {
		ID:    9,
		Title: "CNI Socks",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/DSC_2302_large.jpg?v=1550500895",
		Price: 8,
	},
	10: {
		ID:    10,
		Title: "The Men's containerd Short Sleeve Concert Tee",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/3001C_mu_2eea78bd-4601-413a-ab40-6fbd1a36becf_large.jpg?v=1492272070",
		Price: 15,
	},
	11: {
		ID:    11,
		Title: "The Ladies CoreDNS Short Sleeve Concert Tee",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/CoreDNS_tee_mu_large.jpg?v=1489461006",
		Price: 15,
	},
	12: {
		ID:    12,
		Title: "rkt Decal",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/rkt_decal_Artboard_1_large.png?v=1528318547",
		Price: 0.5,
	},
	13: {
		ID:    13,
		Title: "etcd Full-Zip Hoodie",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/etcd_hoodie_large.jpg?v=1544994217",
		Price: 35,
	},
	14: {
		ID:    14,
		Title: "Helm Socks",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/DSC_2306_large.jpg?v=1550501348",
		Price: 8,
	},
	15: {
		ID:    15,
		Title: "Gloss White H2go Force Copper Vacuum Insulated Thermal Bottle",
		Image: "https://cdn.shopify.com/s/files/1/1300/8977/products/Prometheus_Force_Water_Bottle_large.jpg?v=1469027980",
		Price: 19,
	},
}

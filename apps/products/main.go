package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int64   `json:"id"`
	Title string  `json:"title"`
	Image string  `json:"image"`
	Price float64 `json:"price"`
	Sales int64   `json:"sales"`
	Stock int64   `json:"stock"`
}

func main() {
	http.HandleFunc("/products", productsController)
	fmt.Println("staring products service on port 5000")
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

func getProduct(id int64) *Product {
	if p, ok := mockDB[id]; ok {
		return &p
	}
	return nil
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
		Title: "",
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

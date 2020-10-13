package product

import (
	"encoding/json"
	"fmt"
	"github.com/lucaslra/inventory-mgmt/cors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const productsRouteBasePath = "products"

func SetupRoutes(apiBasePath string) {
	productListHandler := http.HandlerFunc(listHandler)
	productItemHandler := http.HandlerFunc(handler)

	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsRouteBasePath), cors.Middleware(productListHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsRouteBasePath), cors.Middleware(productItemHandler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", productsRouteBasePath))
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var product *Product
	product, err = getProduct(productId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(productJson)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductID != productId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		deleteProduct(productId)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lucaslra/inventory-mgmt/database"
	"github.com/lucaslra/inventory-mgmt/product"
	"github.com/lucaslra/inventory-mgmt/receipt"
	"net/http"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	_ = http.ListenAndServe(":5000", nil)
}

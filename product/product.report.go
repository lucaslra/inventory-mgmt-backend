package product

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type ProductReportFilter struct {
	NameFilter         string `json:"productName"`
	ManufacturerFilter string `json:"manufacturer"`
	SKUFilter          string `json:"sku"`
}

func handleProductReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var productFilter ProductReportFilter
		err := json.NewDecoder(r.Body).Decode(&productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		products, err := searchForProductData(productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		modFuncMap := template.FuncMap{"mod": func(x, i int) bool {
			return x%i == 0
		}}
		titleFuncMap := template.FuncMap{"title": func(v string) string {
			return strings.Title(v)
		}}

		t := template.New("report.gotmpl").Funcs(modFuncMap).Funcs(titleFuncMap)
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var tmpl bytes.Buffer
		if len(products) > 0 {
			err = t.Execute(&tmpl, products)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rdr := bytes.NewReader(tmpl.Bytes())
		w.Header().Set("Content-Disposition", "Attachment")
		http.ServeContent(w, r, "report.html", time.Now(), rdr)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

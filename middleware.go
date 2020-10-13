package main

import (
	"fmt"
	"net/http"
	"time"
)

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Println("Started")
		handler.ServeHTTP(w, r)
		fmt.Println("Took", time.Since(start))
	})
}

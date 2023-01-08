package main

import (
	"appstore/backend"
	"appstore/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
    fmt.Println("started-service")

	// ES backend
	backend.InitElasticsearchBackend()

	// GCS backend
	backend.InitGCSBackend()

	// HTTP router
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}

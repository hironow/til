package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"./internal"
)

func main() {
	ctx := context.Background()
	projectID := "hironow-datastore-dev"

	info := &internal.ConnectionInfo{
		ProjectID: projectID,
	}

	server, err := internal.InitializeServer(ctx, info, true)
	if err != nil {
		log.Fatalf("Failed to init server: %+v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}

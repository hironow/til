package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	projectID := "hironow-datastore-dev"

	info := &ConnectionInfo{
		ProjectID: projectID,
	}

	server, err := InitializeServer(ctx, info, true)
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

func NewServer(indexService *IndexService, userService *UserService, bookService *BookService) *http.ServeMux {
	mux := http.NewServeMux()

	for pattern, handlerFunc := range indexService.handlerFuncMap {
		mux.Handle(pattern, handlerFunc)
	}
	for pattern, handlerFunc := range userService.handlerFuncMap {
		mux.Handle(pattern, handlerFunc)
	}
	for pattern, handlerFunc := range bookService.handlerFuncMap {
		mux.Handle(pattern, handlerFunc)
	}

	return mux
}

// Service

type IndexService struct {
	handlerFuncMap map[string]http.HandlerFunc
}

func NewIndexService() *IndexService {
	m := make(map[string]http.HandlerFunc)

	m["/"] = func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "Hello, World!")
	}

	return &IndexService{handlerFuncMap: m}
}

package internal

import (
	"fmt"
	"net/http"
)

// Service

type IndexService struct {
	handlerFuncMap map[string]http.HandlerFunc
}

func NewIndexService() *IndexService {
	m := make(map[string]http.HandlerFunc)

	m["/"] = indexHandler

	return &IndexService{handlerFuncMap: m}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

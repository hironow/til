package internal

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Service

type IndexService struct {
}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (i *IndexService) SetRouter(r *mux.Router) {
	r.HandleFunc("/", i.indexHandler).Methods("GET")
}

func (i *IndexService) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, World!")
}

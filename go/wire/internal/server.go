package internal

import (
	"github.com/gorilla/mux"
)

func NewServer(indexService *IndexService, userService *UserService, bookService *BookService) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	indexService.SetRouter(router)
	userService.SetRouter(router)
	bookService.SetRouter(router)

	return router
}

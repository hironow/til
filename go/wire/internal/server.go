package internal

import "net/http"

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

//+build wireinject

package main

import (
	"context"
	"net/http"

	"github.com/google/wire"
)

func InitializeEvent(phrase string) (Event, error) {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}, nil
}

func InitializeUserStore(ctx context.Context, info *ConnectionInfo, debug bool) (*UserStore, error) {
	wire.Build(NewUserStore, NewDefaultConfig, NewClient)
	return &UserStore{}, nil
}

func InitializeBookStore(ctx context.Context, info *ConnectionInfo, debug bool) (*BookStore, error) {
	wire.Build(NewBookStore, NewDefaultConfig, NewClient)
	return &BookStore{}, nil
}

func InitializeServer(ctx context.Context, info *ConnectionInfo, debug bool) (*http.ServeMux, error) {
	wire.Build(
		NewServer,
		NewIndexService, NewUserService, NewBookService,
		NewUserStore, NewBookStore,
		NewClient, NewDefaultConfig)
	return &http.ServeMux{}, nil
}

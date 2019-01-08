//+build wireinject

package internal

import (
	"context"

	"github.com/google/wire"
	"github.com/gorilla/mux"
)

func InitializeUserStore(ctx context.Context, info *ConnectionInfo, debug bool) (*UserStore, error) {
	wire.Build(NewUserStore, NewDefaultConfig, NewClient)
	return &UserStore{}, nil
}

func InitializeBookStore(ctx context.Context, info *ConnectionInfo, debug bool) (*BookStore, error) {
	wire.Build(NewBookStore, NewDefaultConfig, NewClient)
	return &BookStore{}, nil
}

func InitializeServer(ctx context.Context, info *ConnectionInfo, debug bool) (*mux.Router, error) {
	wire.Build(
		NewServer,
		NewIndexService, NewUserService, NewBookService,
		NewUserStore, NewBookStore,
		NewClient, NewDefaultConfig)
	return &mux.Router{}, nil
}

package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/satori/go.uuid"
	"google.golang.org/api/iterator"
)

// Entity

type UserID string

func NewUserID() (UserID, error) {
	uID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return UserID(uID.String()), nil
}

func GenerateUserKey(id UserID) *datastore.Key {
	return datastore.NameKey("User", string(id), nil)
}

type User struct {
	ID   UserID
	Name string
}

func NewUser(name string) (*User, error) {
	id, err := NewUserID()
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Name: name}, nil
}

// Store

type UserStore struct {
	client *datastore.Client
	debug  bool
}

func NewUserStore(cfg *Config, client *datastore.Client) *UserStore {
	return &UserStore{client: client, debug: cfg.Debug}
}

func (u *UserStore) Get(ctx context.Context, id UserID) (*User, error) {
	var user User
	err := u.client.Get(ctx, GenerateUserKey(id), &user)

	if u.debug {
		log.Printf("User: %+v", user)
	}

	return &user, err
}

func (u *UserStore) Put(ctx context.Context, user *User) error {
	_, err := u.client.Put(ctx, GenerateUserKey(user.ID), user)
	return err
}

func (u *UserStore) Delete(ctx context.Context, id UserID) error {
	return u.client.Delete(ctx, GenerateUserKey(id))
}

func (u *UserStore) List(ctx context.Context, cursorString string) ([]*User, error) {
	q := datastore.NewQuery("User")
	if cursorString != "" {
		cursor, err := datastore.DecodeCursor(cursorString)
		if err != nil {
			return nil, err
		}
		q.Start(cursor)
	}

	var users []*User
	t := u.client.Run(ctx, q)
	for {
		var user User
		_, err := t.Next(&user)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if u.debug {
		for i, user := range users {
			log.Printf("#%d User: %+v", i, user)
		}
	}

	return users, nil
}

// Service

type UserService struct {
	handlerFuncMap map[string]http.HandlerFunc
}

func NewUserService(userStore *UserStore) *UserService {
	m := make(map[string]http.HandlerFunc)

	m["/users/"] = func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			users, err := userStore.List(ctx, "")
			if err != nil {
				fmt.Fprintf(w, "get user list error: %+v", err)
				return
			}
			for i, user := range users {
				fmt.Fprintf(w, "#%d User: %+v", i, user)
			}
			return
		}

		user, err := userStore.Get(ctx, UserID(id))
		if err != nil {
			fmt.Fprintf(w, "get user error: %+v", err)
			return
		}
		fmt.Fprintf(w, "User: %+v", user)
	}

	return &UserService{handlerFuncMap: m}
}

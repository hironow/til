package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
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
		log.Printf("User: %#v", user)
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
			log.Printf("#%d User: %#v", i, user)
		}
	}

	return users, nil
}

// Service

type UserService struct {
	userStore *UserStore
}

func NewUserService(userStore *UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (u *UserService) SetRouter(r *mux.Router) {
	r.HandleFunc("/users", u.usersHandler).Methods("GET", "POST")
	r.HandleFunc("/users/{userID}", u.userHandler).Methods("GET", "DELETE")
}

func (u *UserService) usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case "GET":
		users, err := u.userStore.List(ctx, "")
		if err != nil {
			fmt.Fprintf(w, "get user list error: %#v", err)
			return
		}
		if len(users) == 0 {
			fmt.Fprint(w, "no users!")
		}

		for i, user := range users {
			fmt.Fprintf(w, "#%d User: %#v\n", i, user)
		}

		w.WriteHeader(http.StatusOK)
		return
	case "POST":
		userName := r.FormValue("name")

		user, err := NewUser(userName)
		if err != nil {
			fmt.Fprintf(w, "create user error: %#v", err)
			return
		}

		err = u.userStore.Put(ctx, user)
		if err != nil {
			fmt.Fprintf(w, "create user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		w.WriteHeader(http.StatusOK)
		return
	}
}

func (u *UserService) userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	switch r.Method {
	case "GET":
		user, err := u.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		w.WriteHeader(http.StatusOK)
		return
	case "DELETE":
		user, err := u.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		err = u.userStore.Delete(ctx, user.ID)
		if err != nil {
			fmt.Fprintf(w, "delete user error: %#v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

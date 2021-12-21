package handler

import (
	"net/http"

	"github.com/Fe4p3b/go-backend-2/lesson1/internal/app/repos/catalog"
)

type handler struct {
	environmentStore catalog.EnvironmentStore
	userStore        catalog.UserStore
}

func NewHandler(e catalog.EnvironmentStore, u catalog.UserStore) *handler {
	return &handler{
		environmentStore: e,
		userStore:        u,
	}
}

func (h *handler) ServeHTTP() error {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/environments", environmentsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}

func usersHandler(w http.ResponseWriter, r *http.Request)
func environmentsHandler(w http.ResponseWriter, r *http.Request)

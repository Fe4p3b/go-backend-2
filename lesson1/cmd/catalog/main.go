package main

import (
	"log"

	"github.com/Fe4p3b/go-backend-2/lesson1/internal/api/handler"
	"github.com/Fe4p3b/go-backend-2/lesson1/internal/db/pg"
)

func main() {
	db, err := pg.NewDB("127.0.0.1:5432")
	if err != nil {
		log.Fatal(err)
	}

	e := pg.NewEnvironments(db)
	u := pg.NewUsers(db)

	h := handler.NewHandler(e, u)

	if err := h.ServeHTTP(); err != nil {
		log.Fatal(err)
	}
}

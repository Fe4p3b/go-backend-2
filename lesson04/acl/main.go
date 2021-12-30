package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Fe4p3b/go-backend-2/lesson04/acl/red"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db         *sql.DB
	measurable = red.MeasurableHandler

	router = mux.NewRouter()
	web    = http.Server{
		Handler: router,
	}
)

func init() {
	router.
		HandleFunc("/identity", measurable(GetIdentityHandler)).
		Methods(http.MethodGet)

	var err error
	db, err = sql.Open("mysql", "root:test@tcp(mysql:3306)/test")
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("starting server")
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
	log.Println("server stopped")
}

const sqlSelectToken = `
SELECT IF(COUNT(*),'true','false') from tokens WHERE token = ?
`

func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	rr, err := db.Query(sqlSelectToken, r.FormValue("token"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rr.Close()

	var isTokenCorrect bool
	for rr.Next() {
		err = rr.Scan(&isTokenCorrect)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	log.Println(isTokenCorrect)

	if isTokenCorrect {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

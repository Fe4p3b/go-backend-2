package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/echo/v4"
)

var es, _ = elasticsearch.NewDefaultClient()

// var esclient, _ = elastic.NewClient(elastic.SetURL("http://localhost:9200"),
// 	elastic.SetSniff(false),
// 	elastic.SetHealthcheck(false))

func main() {
	e := echo.New()
	e.POST("/items", saveItem)
	e.GET("/items", getItems)
	e.GET("/items/:id", getItems)

	// log.Println(elasticsearch.Version)
	// res, err := es.Info()
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }
	// defer res.Body.Close()
	// log.Println(res)

	log.Fatal(e.Start(":8080"))
}

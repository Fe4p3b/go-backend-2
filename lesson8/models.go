package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var elasticIndex string = "item"

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
}

func (i *Item) Save() error {
	item, err := json.Marshal(i)

	if err != nil {
		return err
	}

	request := esapi.IndexRequest{
		Index:      elasticIndex,
		DocumentID: i.ID,
		Body:       strings.NewReader(string(item)),
	}
	request.Do(context.Background(), es)

	// _, err = esclient.Index().
	// 	Index(elasticIndex).
	// 	BodyJson(string(item)).
	// 	Do(context.Background())

	// if err != nil {
	// 	return err
	// }

	return nil
}

/* func (i *Item) Search(matches map[string]string) ([]Item, error) {
	var items []Item

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", "cookie"))

	searchService := esclient.Search().Index(elasticIndex).SearchSource(searchSource)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var item Item
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		items = append(items, item)
	}

	return items, nil
} */

func (i *Item) getQuery(id string, name string, description string, price string, quantity string) map[string]interface{} {
	should := make([]interface{}, 0, 5)

	if id != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		})
	}

	if name != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		})
	}

	if description != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"description": description,
			},
		})
	}

	if price != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"price": price,
			},
		})
	}

	if quantity != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"quantity": quantity,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	return query
}

func (i *Item) Search(id string, name string, description string, price string, quantity string) ([]Item, error) {
	var buffer bytes.Buffer

	query := i.getQuery(id, name, description, price, quantity)

	_ = json.NewEncoder(&buffer).Encode(query)
	resp, err := es.Search(es.Search.WithIndex(elasticIndex), es.Search.WithBody(&buffer))
	if err != nil {
		return nil, err
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source Item `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	_ = json.NewDecoder(resp.Body).Decode(&hits) // XXX: error omitted

	items := make([]Item, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		items[i].ID = hit.Source.ID
		items[i].Name = hit.Source.Name
		items[i].Description = hit.Source.Description
		items[i].Price = hit.Source.Price
		items[i].Quantity = hit.Source.Quantity
	}

	return items, nil
}

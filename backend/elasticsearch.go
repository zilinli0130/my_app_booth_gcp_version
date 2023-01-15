package backend

import (
	"context"
	"fmt"

	"appstore/constants"

	"github.com/olivere/elastic/v7"
)

// Global variable
var (
    ESBackend *ElasticsearchBackend
)

// ElasticsearchBackend class
type ElasticsearchBackend struct {
    client *elastic.Client
}

// Initialize the singleton backend object
func InitElasticsearchBackend() {

    // Create the ES client
    newClient, err := elastic.NewClient(
        elastic.SetURL(constants.ES_URL),
        elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
    if err != nil {
        panic(err)
    }

    // Check if the APP_INDEX exists or not here
    exists, err := newClient.IndexExists(constants.APP_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    // Create the APP_INDEX if not exists
    if !exists {

        // App schema
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "title":      { "type": "text"},
                    "description":  { "type": "text" },
                    "price":      { "type": "keyword", "index": false },
                    "url":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := newClient.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    // Check if the APP_INDEX exists or not
    exists, err = newClient.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    // Create the USER_INDEX if not exists
    if !exists {

        // User schema
        mapping := `{
                     "mappings": {
                         "properties": {
                            "username": {"type": "keyword"},
                            "password": {"type": "keyword"},
                            "age": {"type": "long", "index": false},
                            "gender": {"type": "keyword", "index": false}
                         }
                    }
                }`
        _, err = newClient.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    fmt.Println("Indexes are created.")
    ESBackend = &ElasticsearchBackend{client: newClient}
}

// Search in ES
// func (receiver) name (input) (output)
// receiver.name()
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

// Save to ES
func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}

// Delete from ES
func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}




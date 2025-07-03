package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var Client *elasticsearch.Client

func InitElasticSearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // ganti jika beda
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	Client = es
}

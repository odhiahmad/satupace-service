// helper/elasticsearch.go
package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	es "github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/entity"
)

func CreateElasticProductIndex() error {
	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"tokenizer": map[string]interface{}{
					"edge_ngram_tokenizer": map[string]interface{}{
						"type":     "edge_ngram",
						"min_gram": 1,
						"max_gram": 20,
						"token_chars": []string{
							"letter", "digit",
						},
					},
				},
				"analyzer": map[string]interface{}{
					"edge_ngram_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "edge_ngram_tokenizer",
						"filter":    []string{"lowercase"},
					},
					"edge_ngram_search_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter":    []string{"lowercase"},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				// Fields yang menggunakan autocomplete
				"name": map[string]interface{}{
					"type":            "text",
					"analyzer":        "edge_ngram_analyzer",
					"search_analyzer": "edge_ngram_search_analyzer",
				},
				"description": map[string]interface{}{
					"type":            "text",
					"analyzer":        "edge_ngram_analyzer",
					"search_analyzer": "edge_ngram_search_analyzer",
				},

				// Field lain yang sering dicari atau ditampilkan
				"brand_name": map[string]interface{}{
					"type": "text",
				},
				"category_name": map[string]interface{}{
					"type": "text",
				},
				"variant_names": map[string]interface{}{
					"type": "text",
				},

				// Metadata produk
				"id": map[string]interface{}{
					"type": "integer",
				},
				"business_id": map[string]interface{}{
					"type": "integer",
				},
				"base_price": map[string]interface{}{
					"type": "float",
				},
				"sku": map[string]interface{}{
					"type": "keyword",
				},
				"stock": map[string]interface{}{
					"type": "integer",
				},
				"is_active": map[string]interface{}{
					"type": "boolean",
				},
				"is_available": map[string]interface{}{
					"type": "boolean",
				},
				"track_stock": map[string]interface{}{
					"type": "boolean",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_time",
				},
				"updated_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_time",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("failed to encode index mapping: %w", err)
	}

	// Hapus index lama jika perlu (opsional)
	// es.Client.Indices.Delete([]string{"products"})

	res, err := es.Client.Indices.Create("products", es.Client.Indices.Create.WithBody(&buf))
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index creation error: %s", res.String())
	}

	return nil
}

func IndexProductToElastic(product *entity.Product) error {
	// Bangun struktur data untuk ES
	data := map[string]interface{}{
		"id":           product.Id,
		"business_id":  product.BusinessId,
		"name":         product.Name,
		"description":  product.Description,
		"has_variant":  product.HasVariant,
		"is_active":    product.IsActive,
		"is_available": product.IsAvailable,
		"track_stock":  product.TrackStock,
		"created_at":   product.CreatedAt.Format(time.RFC3339),
		"updated_at":   product.UpdatedAt.Format(time.RFC3339),
	}

	if product.BasePrice != nil {
		data["base_price"] = *product.BasePrice
	}
	if product.Image != nil {
		data["image"] = *product.Image
	}
	if product.SKU != nil {
		data["sku"] = *product.SKU
	}
	if product.Brand != nil {
		data["brand_name"] = product.Brand.Name
	}
	if product.Stock != nil {
		data["stock"] = *product.Stock
	}
	if product.Category != nil {
		data["category_name"] = product.Category.Name
	}
	if product.Tax != nil {
		data["tax_name"] = product.Tax.Name
		data["tax_amount"] = product.Tax.Amount
	}
	if product.Unit != nil {
		data["unit_name"] = product.Unit.Name
	}
	if product.Discount != nil {
		data["discount_name"] = product.Discount.Name
		data["discount_is_percentage"] = product.Discount.IsPercentage
		data["discount_amount"] = product.Discount.Amount
	}
	if len(product.Variants) > 0 {
		var variantNames []string
		for _, v := range product.Variants {
			variantNames = append(variantNames, v.Name)
		}
		data["variant_names"] = variantNames
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal product for Elasticsearch: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "products",
		DocumentID: strconv.Itoa(product.Id),
		Body:       bytes.NewReader(payload),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("error indexing product to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error for product ID %d: %s", product.Id, res.String())
	}

	return nil
}

func SearchProductElastic(keyword string, businessId int) ([]int, error) {
	var result []int

	// Build query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  keyword,
							"fields": []string{"name", "description", "brand_name", "unit_name", "category_name"},
						},
					},
					{
						"match": map[string]interface{}{
							"business_id": businessId, // fix: field name sesuai index
						},
					},
				},
			},
		},
	}

	// Debug log untuk query
	queryJSON, _ := json.MarshalIndent(query, "", "  ")
	fmt.Println("🔍 ES Query:", string(queryJSON))

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Println("❌ Error encode query:", err)
		return nil, err
	}

	// Kirim request ke Elasticsearch
	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex("products"),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)
	if err != nil {
		fmt.Println("❌ Error request ke Elasticsearch:", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Println("❌ Response error dari Elasticsearch:", res.String())
		return nil, fmt.Errorf("error from Elasticsearch: %s", res.String())
	}

	// Decode response dari ES
	var esResp struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		fmt.Println("❌ Error decode response:", err)
		return nil, err
	}

	// Log hasil ID dari Elasticsearch
	fmt.Printf("✅ Ditemukan %d produk di Elasticsearch\n", len(esResp.Hits.Hits))

	for _, hit := range esResp.Hits.Hits {
		id, err := strconv.Atoi(hit.ID)
		if err != nil {
			fmt.Println("⚠️ Gagal konversi _id:", hit.ID)
			continue
		}
		result = append(result, id)
	}

	return result, nil
}

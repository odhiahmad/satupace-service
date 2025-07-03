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
	if product.Stock != nil {
		data["stock"] = *product.Stock
	}
	if product.ProductCategory != nil {
		data["category_name"] = product.ProductCategory.Name
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

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  keyword,
							"fields": []string{"name", "description", "brand"},
						},
					},
					{
						"match": map[string]interface{}{
							"businessId": businessId,
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex("products"),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error from Elasticsearch: %s", res.String())
	}

	var esResp struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, err
	}

	for _, hit := range esResp.Hits.Hits {
		id, _ := strconv.Atoi(hit.ID)
		result = append(result, id)
	}

	return result, nil
}

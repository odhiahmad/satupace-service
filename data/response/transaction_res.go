package response

import "time"

type TransactionItemAttributeResponse struct {
	Id                 int     `json:"id"`
	ProductAttributeId int     `json:"product_attribute_id"`
	AdditionalPrice    float64 `json:"additional_price"`
}

type TransactionItemResponse struct {
	Id                 int                                `json:"id"`
	ProductId          *int                               `json:"product_id,omitempty"`
	BundleId           *int                               `json:"bundle_id,omitempty"`
	ProductAttributeId *int                               `json:"product_attribute_id,omitempty"`
	ProductVariantId   *int                               `json:"product_variant_id,omitempty"`
	Quantity           int                                `json:"quantity"`
	UnitPrice          float64                            `json:"unit_price"`
	Price              float64                            `json:"price"`
	Discount           *float64                           `json:"discount"`
	Promo              *float64                           `json:"promo"`
	Rating             *float64                           `json:"rating"`
	Attributes         []TransactionItemAttributeResponse `json:"attributes"`
}

type TransactionResponse struct {
	Id              int                       `json:"transaction_id"`
	BusinessId      int                       `json:"business_id"`
	CustomerId      *int                      `json:"customer_id,omitempty"`
	PaymentMethodId *int                      `json:"payment_method_id,omitempty"`
	BillNumber      string                    `json:"bill_number"`
	Items           []TransactionItemResponse `json:"items"`
	Total           float64                   `json:"total"`
	Discount        *float64                  `json:"discount"`
	Promo           *float64                  `json:"promo"`
	Status          *string                   `json:"status"`
	Rating          *float64                  `json:"rating"`
	Notes           *string                   `json:"notes"`
	AmountReceived  *float64                  `json:"amount_received"`
	Change          *float64                  `json:"change"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

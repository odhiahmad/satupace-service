package response

type TransactionItemAttributeResponse struct {
	Id                 int     `json:"id"`
	ProductAttributeId int     `json:"product_attribute_id"`
	AdditionalPrice    float64 `json:"additional_price"`
}

type TransactionItemResponse struct {
	Id                 int                                `json:"id"`
	ProductId          *int                               `json:"product_id,omitempty"`
	Product            []ProductResponse                  `json:"product"`
	BundleId           *int                               `json:"bundle_id,omitempty"`
	Bundle             []BundleResponse                   `json:"bundle"`
	ProductAttributeId *int                               `json:"product_attribute_id,omitempty"`
	ProductVariantId   *int                               `json:"product_variant_id,omitempty"`
	Quantity           int                                `json:"quantity"`
	Attributes         []TransactionItemAttributeResponse `json:"attributes"`
	Discount           float64                            `json:"discount"`
	Promo              float64                            `json:"promo"`
	Tax                float64                            `json:"tax"`
	UnitPrice          float64                            `json:"unitPrice"`
	Total              float64                            `json:"total"`
}

type TransactionResponse struct {
	Id              int                       `json:"transaction_id"`
	BusinessId      int                       `json:"business_id"`
	CustomerId      *int                      `json:"customer_id,omitempty"`
	PaymentMethodId *int                      `json:"payment_method_id,omitempty"`
	BillNumber      string                    `json:"bill_number"`
	Items           []TransactionItemResponse `json:"items"`
	FinalPrice      float64                   `json:"finalPrice"`
	BasePrice       float64                   `json:"basePrice"`
	Discount        float64                   `json:"discount"`
	Promo           float64                   `json:"promo"`
	Tax             float64                   `json:"tax"`
	Status          string                    `json:"status"`
	Rating          *float64                  `json:"rating"`
	Notes           *string                   `json:"notes"`
	AmountReceived  *float64                  `json:"amount_received"`
	Change          *float64                  `json:"change"`
	Paid            float64                   `json:"paid_at"`
}

package request

type TransactionItemAttributeCreate struct {
	ProductAttributeId int     `json:"product_attribute_id" validate:"required"`
	AdditionalPrice    float64 `json:"additional_price" validate:"required"`
}

type TransactionItemAttributeUpdate struct {
	Id                 int     `json:"id"`
	ProductAttributeId int     `json:"product_attribute_id"`
	AdditionalPrice    float64 `json:"additional_price"`
}

type TransactionItemCreate struct {
	ProductId          *int                             `json:"product_id,omitempty"`
	BundleId           *int                             `json:"bundle_id,omitempty"`
	ProductAttributeId *int                             `json:"product_attribute_id,omitempty"`
	ProductVariantId   *int                             `json:"product_variant_id,omitempty"`
	Quantity           int                              `json:"quantity"`
	Discount           float64                          `json:"discount"`
	Promo              float64                          `json:"promo"`
	Tax                float64                          `json:"tax"`
	UnitPrice          float64                          `json:"unitPrice"`
	Total              float64                          `json:"total"`
	Rating             *float64                         `json:"rating"`
	Attributes         []TransactionItemAttributeCreate `json:"attributes"`
}

type TransactionItemUpdate struct {
	Id                 int                              `json:"id"`
	ProductId          *int                             `json:"product_id,omitempty"`
	BundleId           *int                             `json:"bundle_id,omitempty"`
	ProductAttributeId *int                             `json:"product_attribute_id,omitempty"`
	ProductVariantId   *int                             `json:"product_variant_id,omitempty"`
	Quantity           int                              `json:"quantity"`
	UnitPrice          float64                          `json:"unitPrice"`
	Total              float64                          `json:"total"`
	Discount           float64                          `json:"discount"`
	Promo              float64                          `json:"promo"`
	Tax                float64                          `json:"tax"`
	Rating             *float64                         `json:"rating"`
	Attributes         []TransactionItemAttributeUpdate `json:"attributes"`
}

type TransactionCreateRequest struct {
	BusinessId int                     `json:"business_id" validate:"required"`
	CustomerId *int                    `json:"customer_id,omitempty"`
	Items      []TransactionItemCreate `json:"items" validate:"required,dive"`
}

type TransactionPaymentRequest struct {
	Id              int      `json:"id"`
	CustomerId      *int     `json:"customer_id,omitempty"`
	PaymentMethodId *int     `json:"payment_method_id,omitempty"`
	Rating          *float64 `json:"rating"`
	Notes           *string  `json:"notes"`
	AmountReceived  *float64 `json:"amount_received"`
}

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
	Quantity           int                              `json:"quantity" validate:"required,gt=0"`
	Price              float64                          `json:"price"`    // Hapus required karena dihitung di backend
	Discount           float64                          `json:"discount"` // Dihitung juga
	Promo              float64                          `json:"promo"`    // Dihitung juga
	Rating             *float64                         `json:"rating"`   // Optional
	Attributes         []TransactionItemAttributeCreate `json:"attributes"`
}

type TransactionItemUpdate struct {
	Id                 int                              `json:"id"`
	ProductId          *int                             `json:"product_id,omitempty"`
	BundleId           *int                             `json:"bundle_id,omitempty"`
	ProductAttributeId *int                             `json:"product_attribute_id,omitempty"`
	ProductVariantId   *int                             `json:"product_variant_id,omitempty"`
	Quantity           int                              `json:"quantity"`
	UnitPrice          float64                          `json:"unit_price"`
	Price              float64                          `json:"price"`
	Discount           float64                          `json:"discount"`
	Promo              float64                          `json:"promo"`
	Rating             *float64                         `json:"rating"`
	Attributes         []TransactionItemAttributeUpdate `json:"attributes"`
}

type TransactionCreateRequest struct {
	BusinessId int                     `json:"business_id" validate:"required"`
	CustomerId *int                    `json:"customer_id,omitempty"`
	Items      []TransactionItemCreate `json:"items" validate:"required,dive"`
}

type TransactionUpdateRequest struct {
	Id              int                     `json:"id"`
	CustomerId      *int                    `json:"customer_id,omitempty"`
	PaymentMethodId *int                    `json:"payment_method_id,omitempty"`
	BillNumber      string                  `json:"bill_number"`
	Items           []TransactionItemUpdate `json:"items"`
	Total           float64                 `json:"total"`
	Discount        float64                 `json:"discount"`
	Promo           float64                 `json:"promo"`
	Status          string                  `json:"status"`
	Rating          *float64                `json:"rating"`
	Notes           *string                 `json:"notes"`
	AmountReceived  *float64                `json:"amount_received"`
	Change          *float64                `json:"change"`
}

package response

import "time"

type TransactionItemAttributeResponse struct {
	Id                 int     `json:"id"`
	ProductAttributeId int     `json:"product_attribute_id"`
	AdditionalPrice    float64 `json:"additional_price"`
}

type TransactionItemResponse struct {
	Id                 int                                `json:"id"`
	ProductId          *int                               `json:"product_id"`
	Product            []ProductResponse                  `json:"product"`
	BundleId           *int                               `json:"bundle_id"`
	Bundle             []BundleResponse                   `json:"bundle"`
	ProductAttributeId *int                               `json:"product_attribute_id"`
	ProductVariantId   *int                               `json:"product_variant_id"`
	Quantity           int                                `json:"quantity"`
	Attributes         []TransactionItemAttributeResponse `json:"attributes"`
	Discount           float64                            `json:"discount"`
	Promo              float64                            `json:"promo"`
	Tax                float64                            `json:"tax"`
	UnitPrice          float64                            `json:"unitPrice"`
	Total              float64                            `json:"total"`
	BasePrice          float64                            `json:"basePrice"`
	SellPrice          float64                            `json:"sellPrice"`
}

type TransactionResponse struct {
	Id              int                       `json:"transaction_id"`
	BusinessId      int                       `json:"business_id"`
	CustomerId      *int                      `json:"customer_id"`
	Cashier         UserBusinessResponse      `json:"cashier"`
	PaymentMethodId *int                      `json:"payment_method_id"`
	BillNumber      string                    `json:"bill_number"`
	Items           []TransactionItemResponse `json:"items"`
	FinalPrice      float64                   `json:"final_price"`
	BasePrice       float64                   `json:"base_price"`
	SellPrice       float64                   `json:"sell_price"`
	Discount        float64                   `json:"discount"`
	Promo           float64                   `json:"promo"`
	Tax             float64                   `json:"tax"`
	Status          string                    `json:"status"`
	Rating          *float64                  `json:"rating"`
	Notes           *string                   `json:"notes"`
	AmountReceived  *float64                  `json:"amount_received"`
	Change          *float64                  `json:"change"`
	PaidAt          *time.Time                `json:"paid_at"`
	RefundedAt      *time.Time                `json:"refunded_at"`
	RefundedBy      *int                      `json:"refunded_by"`
	RefundReason    *string                   `json:"refund_reason"`
	IsRefunded      *bool                     `json:"is_refunded"`
	CanceledAt      *time.Time                `json:"canceled_at"`
	CanceledBy      *int                      `json:"canceled_by"`
	CanceledReason  *string                   `json:"canceled_reason"`
	IsCanceled      *bool                     `json:"is_canceled"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

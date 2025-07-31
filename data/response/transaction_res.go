package response

import (
	"time"

	"github.com/google/uuid"
)

type TransactionItemAttributeResponse struct {
	Id                 uuid.UUID `json:"id"`
	ProductAttributeId uuid.UUID `json:"product_attribute_id"`
	AdditionalPrice    float64   `json:"additional_price"`
}

type TransactionItemResponse struct {
	Id                 uuid.UUID                          `json:"id"`
	ProductId          *uuid.UUID                         `json:"product_id"`
	Product            []ProductResponse                  `json:"product"`
	BundleId           *uuid.UUID                         `json:"bundle_id"`
	Bundle             []BundleResponse                   `json:"bundle"`
	ProductAttributeId *uuid.UUID                         `json:"product_attribute_id"`
	ProductVariantId   *uuid.UUID                         `json:"product_variant_id"`
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
	Id              uuid.UUID                 `json:"transaction_id"`
	BusinessId      uuid.UUID                 `json:"business_id"`
	CustomerId      *uuid.UUID                `json:"customer_id"`
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
	RefundedBy      *uuid.UUID                `json:"refunded_by"`
	RefundReason    *string                   `json:"refund_reason"`
	IsRefunded      bool                      `json:"is_refunded"`
	CanceledAt      *time.Time                `json:"canceled_at"`
	CanceledBy      *uuid.UUID                `json:"canceled_by"`
	CanceledReason  *string                   `json:"canceled_reason"`
	IsCanceled      bool                      `json:"is_canceled"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

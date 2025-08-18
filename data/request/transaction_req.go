package request

import "github.com/google/uuid"

type TransactionItemAttributeCreate struct {
	ProductAttributeId uuid.UUID `json:"product_attribute_id" validate:"required"`
	AdditionalPrice    float64   `json:"additional_price" validate:"required"`
}

type TransactionItemAttributeUpdate struct {
	Id                 uuid.UUID `json:"id"`
	ProductAttributeId uuid.UUID `json:"product_attribute_id"`
	AdditionalPrice    float64   `json:"additional_price"`
}

type TransactionItemCreate struct {
	ProductId          *uuid.UUID                       `json:"product_id"`
	BundleId           *uuid.UUID                       `json:"bundle_id"`
	ProductAttributeId *uuid.UUID                       `json:"product_attribute_id"`
	ProductVariantId   *uuid.UUID                       `json:"product_variant_id"`
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
	Id                 uuid.UUID                        `json:"id"`
	ProductId          *uuid.UUID                       `json:"product_id"`
	BundleId           *uuid.UUID                       `json:"bundle_id"`
	ProductAttributeId *uuid.UUID                       `json:"product_attribute_id"`
	ProductVariantId   *uuid.UUID                       `json:"product_variant_id"`
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
	BusinessId   uuid.UUID               `json:"business_id" validate:"required"`
	CustomerName *string                 `json:"customer_name"`
	OrderTypeId  uuid.UUID               `json:"service_type" validate:"required"`
	TableId      *uuid.UUID              `json:"table"`
	Items        []TransactionItemCreate `json:"items" validate:"required,dive"`
	Notes        *string                 `json:"notes"`
}

type TransactionPaymentRequest struct {
	Id              uuid.UUID `json:"id"`
	CashierId       uuid.UUID `json:"cashier_id"`
	PaymentMethodId *int      `json:"payment_method_id"`
	Rating          *float64  `json:"rating"`
	AmountReceived  *float64  `json:"amount_received"`
}

type TransactionCancelRequest struct {
	Id         uuid.UUID  `json:"id"`
	CustomerId *uuid.UUID `json:"customer_id"`
	BusinessId uuid.UUID  `json:"business_id" validate:"required"`
	UserId     uuid.UUID  `json:"user_id"`
	Reason     *string    `json:"reason"`
}

type TransactionRefundRequest struct {
	Id         uuid.UUID  `json:"id"`
	CustomerId *uuid.UUID `json:"customer_id"`
	BusinessId uuid.UUID  `json:"business_id" validate:"required"`
	UserId     uuid.UUID  `json:"user_id"`
	Reason     *string    `json:"reason"`
}

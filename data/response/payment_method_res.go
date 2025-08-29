package response

type PaymentMethodResponse struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Nama string `json:"nama"`
}

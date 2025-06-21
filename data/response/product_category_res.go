package response

type ProductCategoryResponse struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	BusinessId int    `json:"business_id"`
	ParentId   *int   `json:"parent_id,omitempty"`
}

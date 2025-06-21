package request

type ProductCategoryCreate struct {
	Name       string `json:"name" validate:"required"`
	BusinessId int    `json:"business_id" validate:"required"`
	ParentId   *int   `json:"parent_id,omitempty"`
}

type ProductCategoryUpdate struct {
	Id       int    `validate:"required"`
	Name     string `json:"name" validate:"required"`
	ParentId *int   `json:"parent_id,omitempty"`
}

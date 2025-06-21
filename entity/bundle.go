package entity

type Bundle struct {
	Id          int `gorm:"primaryKey"`
	BusinessId  int
	Business    Business `gorm:"foreignKey:BusinessId"`
	Name        string   `gorm:"type:varchar(255)"`
	Description string   `gorm:"type:varchar(255)"`
	Image       string   `gorm:"type:varchar(255)"`
	BasePrice   float64
	Discount    float64
	Promo       float64
	Stock       int
	FinalPrice  float64
	Items       []BundleItem `gorm:"foreignKey:BundleId"`
	IsAvailable bool         `gorm:"not null;column:is_available"`
	IsActive    bool         `gorm:"not null;column:is_active"`
}

func (u *Bundle) Prepare() error {
	u.IsActive = true
	return nil
}

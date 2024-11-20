package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type BusinessBranchRepository interface {
	InsertBusinessBranch(businessBranch entity.BusinessBranch)
	FindById(businessBranchId int) (businessBranch entity.BusinessBranch, err error)
	FindAll() []entity.BusinessBranch
	Delete(businessBranchId int)
}

type BusinessBranchConnection struct {
	Db *gorm.DB
}

func NewBusinessBranchRepository(Db *gorm.DB) BusinessBranchRepository {
	return &BusinessBranchConnection{Db: Db}
}

func (t *BusinessBranchConnection) InsertBusinessBranch(businessBranch entity.BusinessBranch) {
	result := t.Db.Create(&businessBranch)

	helper.ErrorPanic(result.Error)
}

func (t *BusinessBranchConnection) FindById(businessBranchId int) (businessBranchs entity.BusinessBranch, err error) {
	var businessBranch entity.BusinessBranch
	result := t.Db.Find(&businessBranch, businessBranchId)
	if result != nil {
		return businessBranch, nil
	} else {
		return businessBranch, errors.New("tag is not found")
	}
}

func (t *BusinessBranchConnection) FindAll() []entity.BusinessBranch {
	var businessBranch []entity.BusinessBranch
	result := t.Db.Find(&businessBranch)
	helper.ErrorPanic(result.Error)
	return businessBranch
}

func (t *BusinessBranchConnection) Delete(businessBranchId int) {
	var businessBranchs entity.BusinessBranch
	result := t.Db.Where("id = ?", businessBranchId).Delete(&businessBranchs)
	helper.ErrorPanic(result.Error)
}

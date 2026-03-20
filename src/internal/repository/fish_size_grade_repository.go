package repository

import (
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FishSizeGradeRepository --output=./mocks --outpkg=mocks --filename=fish_size_grade_repository.go --structname=MockFishSizeGradeRepository --with-expecter=false
type FishSizeGradeRepository interface {
	WithTx(tx *gorm.DB) FishSizeGradeRepository
	List() ([]*model.FishSizeGrade, error)
	GetByID(id int) (*model.FishSizeGrade, error)
	GetByIDs(ids []int) ([]*model.FishSizeGrade, error)
}

type fishSizeGradeRepository struct {
	db *gorm.DB
}

func NewFishSizeGradeRepository(db *gorm.DB) FishSizeGradeRepository {
	return &fishSizeGradeRepository{db: db}
}

func (r *fishSizeGradeRepository) WithTx(tx *gorm.DB) FishSizeGradeRepository {
	return &fishSizeGradeRepository{db: tx}
}

func (r *fishSizeGradeRepository) List() ([]*model.FishSizeGrade, error) {
	var grades []*model.FishSizeGrade
	err := r.db.Where("deleted_at IS NULL").Order("sort_index ASC").Find(&grades).Error
	return grades, err
}

func (r *fishSizeGradeRepository) GetByID(id int) (*model.FishSizeGrade, error) {
	var grade model.FishSizeGrade
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&grade).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *fishSizeGradeRepository) GetByIDs(ids []int) ([]*model.FishSizeGrade, error) {
	var grades []*model.FishSizeGrade
	err := r.db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&grades).Error
	return grades, err
}

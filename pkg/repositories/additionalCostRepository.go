package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IAdditionalCostRepository interface {
	Create(pond *models.AdditionalCost) (*models.AdditionalCost, error)
	TakeById(id int) (*models.AdditionalCost, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.AdditionalCost, error)
	Update(pond *models.AdditionalCost) error
	WithTrx(trxHandle *gorm.DB) IAdditionalCostRepository
}

type additionalCostRepositoryImp struct {
	dbContext *gorm.DB
}

func NewAdditionalCostRepository(db *gorm.DB) IAdditionalCostRepository {
	return &additionalCostRepositoryImp{
		dbContext: db,
	}
}

func (rp additionalCostRepositoryImp) Create(request *models.AdditionalCost) (*models.AdditionalCost, error) {
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp additionalCostRepositoryImp) TakeById(id int) (*models.AdditionalCost, error) {
	var result *models.AdditionalCost
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found AdditionalCost TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp additionalCostRepositoryImp) WithTrx(trxHandle *gorm.DB) IAdditionalCostRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}

	return &additionalCostRepositoryImp{
		dbContext: trxHandle,
	}
}

// func (rp additionalCostRepositoryImp) TakePage(clientId, page, pageSize int, orderBy, keyword string, billType *string, farmGroupId *int) (*[]models.BillWithFarmGroupName, int64, error) {
// 	var result *[]models.BillWithFarmGroupName
// 	var total int64

// 	joinFarmGroup := fmt.Sprintf("LEFT JOIN %s ON %s.\"FarmGroupId\" = %s.\"Id\"", dbconst.TFarmGroup, dbconst.TBill, dbconst.TFarmGroup)

// 	firstWhereClause := fmt.Sprintf("%s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ?", dbconst.TFarmGroup, dbconst.TBill)
// 	whereClient := fmt.Sprintf("%s.\"ClientId\" = ?", dbconst.TFarmGroup)

// 	query := rp.dbContext.Table(dbconst.TBill).Select(dbconst.TBill+".*").Select(dbconst.TFarmGroup+".\"Name\"").Joins(joinFarmGroup).Order(orderBy).Where(firstWhereClause, false, false).Where(whereClient, clientId)
// 	// .Select(dbconst.TFarmGroup+".\"Name\"")
// 	if keyword != "" {
// 		whereKeyword := fmt.Sprintf("(%s.\"Other\" LIKE ? OR %s.\"Type\" LIKE ? OR %s.\"Name\" LIKE ? OR %s.\"Code\" LIKE ?)", dbconst.TBill, dbconst.TBill, dbconst.TFarmGroup, dbconst.TFarmGroup)
// 		query = query.Where(whereKeyword, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
// 	}

// 	if billType != nil {
// 		whereMode := fmt.Sprintf("%s.\"Type\" = ?", dbconst.TBill)
// 		query = query.Where(whereMode, *billType)
// 	}

// 	if farmGroupId != nil {
// 		whereFarm := fmt.Sprintf("%s.\"Id\" = ?", dbconst.TFarmGroup)
// 		query = query.Where(whereFarm, *farmGroupId)
// 	}

// 	if err := query.Limit(1).Count(&total).Limit(pageSize).Offset(page * pageSize).Find(&result).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	return result, total, nil
// }

func (rp additionalCostRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.AdditionalCost, error) {
	var result *models.AdditionalCost
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found AdditionalCost FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp additionalCostRepositoryImp) Update(request *models.AdditionalCost) error {
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

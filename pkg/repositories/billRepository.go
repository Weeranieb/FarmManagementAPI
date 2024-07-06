package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IBillRepository interface {
	Create(pond *models.Bill) (*models.Bill, error)
	TakeById(id int) (*models.Bill, error)
	TakePage(clientId, page, pageSize int, orderBy, keyword string, billType *string, farmGroupId *int) (*[]models.BillWithFarmGroupName, int64, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Bill, error)
	Update(pond *models.Bill) error
}

type billRepositoryImp struct {
	dbContext *gorm.DB
}

func NewBillRepository(db *gorm.DB) IBillRepository {
	return &billRepositoryImp{
		dbContext: db,
	}
}

func (rp billRepositoryImp) Create(request *models.Bill) (*models.Bill, error) {
	if err := rp.dbContext.Table(dbconst.TBill).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp billRepositoryImp) TakeById(id int) (*models.Bill, error) {
	var result *models.Bill
	if err := rp.dbContext.Table(dbconst.TBill).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Bill TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp billRepositoryImp) TakePage(clientId, page, pageSize int, orderBy, keyword string, billType *string, farmGroupId *int) (*[]models.BillWithFarmGroupName, int64, error) {
	var result *[]models.BillWithFarmGroupName
	var total int64

	joinFarmGroup := fmt.Sprintf("LEFT JOIN %s ON %s.\"FarmGroupId\" = %s.\"Id\"", dbconst.TFarmGroup, dbconst.TBill, dbconst.TFarmGroup)

	firstWhereClause := fmt.Sprintf("%s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ?", dbconst.TFarmGroup, dbconst.TBill)
	whereClient := fmt.Sprintf("%s.\"ClientId\" = ?", dbconst.TFarmGroup)

	query := rp.dbContext.Table(dbconst.TBill).Select(dbconst.TBill+".*").Select(dbconst.TFarmGroup+".\"Name\"").Joins(joinFarmGroup).Order(orderBy).Where(firstWhereClause, false, false).Where(whereClient, clientId)
	// .Select(dbconst.TFarmGroup+".\"Name\"")
	if keyword != "" {
		whereKeyword := fmt.Sprintf("(%s.\"Other\" LIKE ? OR %s.\"Type\" LIKE ? OR %s.\"Name\" LIKE ? OR %s.\"Code\" LIKE ?)", dbconst.TBill, dbconst.TBill, dbconst.TFarmGroup, dbconst.TFarmGroup)
		query = query.Where(whereKeyword, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if billType != nil {
		whereMode := fmt.Sprintf("%s.\"Type\" = ?", dbconst.TBill)
		query = query.Where(whereMode, *billType)
	}

	if farmGroupId != nil {
		whereFarm := fmt.Sprintf("%s.\"Id\" = ?", dbconst.TFarmGroup)
		query = query.Where(whereFarm, *farmGroupId)
	}

	if err := query.Limit(1).Count(&total).Limit(pageSize).Offset(page * pageSize).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (rp billRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Bill, error) {
	var result *models.Bill
	if err := rp.dbContext.Table(dbconst.TBill).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Bill FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp billRepositoryImp) Update(request *models.Bill) error {
	if err := rp.dbContext.Table(dbconst.TBill).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

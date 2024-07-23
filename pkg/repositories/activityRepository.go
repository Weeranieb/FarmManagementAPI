package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IActivityRepository interface {
	Create(request *models.Activity) (*models.Activity, error)
	TakeById(id int) (*models.Activity, error)
	TakePage(clientId, page, pageSize int, orderBy, keyword string, mode *string, farmId *int) (*[]models.ActivityPage, int64, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Activity, error)
	Update(request *models.Activity) error
	WithTrx(trxHandle *gorm.DB) IActivityRepository
}

type activityRepositoryImp struct {
	dbContext *gorm.DB
}

func NewActivityRepository(db *gorm.DB) IActivityRepository {
	return &activityRepositoryImp{
		dbContext: db,
	}
}

func (rp activityRepositoryImp) Create(request *models.Activity) (*models.Activity, error) {
	if err := rp.dbContext.Table(dbconst.TActivitiy).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp activityRepositoryImp) TakeById(id int) (*models.Activity, error) {
	var result *models.Activity
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp activityRepositoryImp) TakePage(clientId, page, pageSize int, orderBy, keyword string, mode *string, farmId *int) (*[]models.ActivityPage, int64, error) {
	var result *[]models.ActivityPage
	var total int64

	joinActivePond := fmt.Sprintf("LEFT JOIN %s ON %s.\"ActivePondId\" = %s.\"Id\"", dbconst.TActivePond, dbconst.TActivitiy, dbconst.TActivePond)
	joinPond := fmt.Sprintf("LEFT JOIN %s ON %s.\"PondId\" = %s.\"Id\"", dbconst.TPond, dbconst.TActivePond, dbconst.TPond)
	joinFarm := fmt.Sprintf("LEFT JOIN %s ON %s.\"FarmId\" = %s.\"Id\"", dbconst.TFarm, dbconst.TPond, dbconst.TFarm)

	firstWhereClause := fmt.Sprintf("%s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ?", dbconst.TActivitiy, dbconst.TActivePond, dbconst.TPond, dbconst.TFarm)
	whereClient := fmt.Sprintf("%s.\"ClientId\" = ?", dbconst.TFarm)

	query := rp.dbContext.Table(dbconst.TActivitiy).Select(fmt.Sprintf("%s.*, %s.\"Name\" as \"FarmName\", %s.\"Name\" as \"PondName\"", dbconst.TActivitiy, dbconst.TFarm, dbconst.TPond)).Joins(joinActivePond).Joins(joinPond).Joins(joinFarm).Order(orderBy).Where(firstWhereClause, false, false, false, false).Where(whereClient, clientId)

	if keyword != "" {
		whereKeyword := fmt.Sprintf("(%s.\"Code\" LIKE ? OR %s.\"Name\" LIKE ?)", dbconst.TPond, dbconst.TPond)
		query = query.Where(whereKeyword, "%"+keyword+"%", "%"+keyword+"%")
	}

	if mode != nil {
		whereMode := fmt.Sprintf("%s.\"Mode\" = ?", dbconst.TActivitiy)
		query = query.Where(whereMode, *mode)
	}

	if farmId != nil {
		whereFarm := fmt.Sprintf("%s.\"Id\" = ?", dbconst.TFarm)
		query = query.Where(whereFarm, *farmId)
	}

	if err := query.Limit(1).Count(&total).Limit(pageSize).Offset(page * pageSize).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (rp activityRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Activity, error) {
	var result *models.Activity
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp activityRepositoryImp) WithTrx(trxHandle *gorm.DB) IActivityRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp activityRepositoryImp) Update(request *models.Activity) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}

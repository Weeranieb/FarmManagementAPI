package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IClientRepository interface {
	Create(request *models.Client) (*models.Client, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error)
	Update(request *models.Client) error
	TakeById(id int) (*models.Client, error)
	WithTrx(trxHandle *gorm.DB) IClientRepository
	GetClientWithFarms(clientId *int) ([]*models.ClientWithFarms, error)
	TakeAll() ([]*models.Client, error)
	TakePage(keyword string) ([]*models.Client, error)
}

type clientRepositoryImp struct {
	dbContext *gorm.DB
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &clientRepositoryImp{
		dbContext: db,
	}
}

func (rp clientRepositoryImp) Create(request *models.Client) (*models.Client, error) {
	if err := rp.dbContext.Table(dbconst.TClient).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp clientRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error) {
	var result *models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Clien FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp clientRepositoryImp) TakeById(id int) (*models.Client, error) {
	var result *models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Client TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp clientRepositoryImp) WithTrx(trxHandle *gorm.DB) IClientRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp clientRepositoryImp) Update(request *models.Client) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}

func (rp clientRepositoryImp) GetClientWithFarms(clientId *int) ([]*models.ClientWithFarms, error) {
	var clientsWithFarms []*models.ClientWithFarms

	// Build the join clause
	joinClause := fmt.Sprintf("JOIN %s ON %s.\"ClientId\" = %s.\"Id\"", dbconst.TClient, dbconst.TFarm, dbconst.TClient)

	// Prepare the query
	query := rp.dbContext.Table(dbconst.TFarm).
		Select(fmt.Sprintf("%s.*, %s.\"Id\" as \"ClientId\", %s.\"Name\" as \"ClientName\"", dbconst.TFarm, dbconst.TClient, dbconst.TClient)).
		Joins(joinClause).
		Where(fmt.Sprintf("%s.\"DelFlag\" = ? AND %s.\"DelFlag\" = ?", dbconst.TFarm, dbconst.TClient), false, false)

	// Add filter by clientId if provided
	if clientId != nil {
		query = query.Where(fmt.Sprintf("%s.\"ClientId\" = ?", dbconst.TFarm), *clientId)
	}

	// Execute the query
	err := query.Scan(&clientsWithFarms).Error

	if err != nil {
		return nil, err
	}

	return clientsWithFarms, nil
}

func (rp clientRepositoryImp) TakeAll() ([]*models.Client, error) {
	var result []*models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"DelFlag\" = ?", false).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp clientRepositoryImp) TakePage(keyword string) ([]*models.Client, error) {
	var result []*models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"Name\" LIKE ? AND \"DelFlag\" = ?", "%"+keyword+"%", false).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

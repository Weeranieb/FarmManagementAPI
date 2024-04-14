package repositories

import (
	"boonmafarm/api/pkg/models"

	"gorm.io/gorm"
)

type IFarmGroupRepository interface {
	Create(request *models.FarmGroup) (*models.FarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error)
	GetFarmGroupWithFarms(id, clientId int) (*models.GetFarmGroupResponse, error)
	Update(request *models.FarmGroup) error
	TakeById(id int) (*models.FarmGroup, error)
}

type FarmGroupRepository struct {
	dbContext *gorm.DB
}

func NewFarmGroupRepository(db *gorm.DB) IFarmGroupRepository {
	return &FarmGroupRepository{
		dbContext: db,
	}
}

func (rp FarmGroupRepository) Create(request *models.FarmGroup) (*models.FarmGroup, error) {
	if err := rp.dbContext.Table("FarmGroups").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table("FarmGroups").Where(query, args...).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp FarmGroupRepository) TakeById(id int) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table("FarmGroups").Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// FIXME Test this function
func (rp FarmGroupRepository) GetFarmGroupWithFarms(id, clientId int) (*models.GetFarmGroupResponse, error) {
	var result models.GetFarmGroupResponse

	// Perform the left join query to fetch data from FarmGroups and Farms tables
	if err := rp.dbContext.Table("FarmGroups").
		Select("FarmGroups.Id, FarmGroups.Code, FarmGroups.Name, FarmGroups.Description, FarmGroups.CreatedBy, FarmGroups.UpdatedBy, FarmGroups.CreatedAt, FarmGroups.UpdatedAt, FarmGroups.ClientId, FarmGroups.Active, FarmGroups.DeletedAt, Farms.Id AS FarmId, Farms.Code AS FarmCode, Farms.Name AS FarmName, Farms.Description AS FarmDescription, Farms.CreatedBy AS FarmCreatedBy, Farms.UpdatedBy AS FarmUpdatedBy, Farms.CreatedAt AS FarmCreatedAt, Farms.UpdatedAt AS FarmUpdatedAt, Farms.ClientId AS FarmClientId, Farms.Active AS FarmActive, Farms.DeletedAt AS FarmDeletedAt").
		Joins("LEFT JOIN FarmOnFarmGroups ON FarmGroups.Id = FarmOnFarmGroups.FarmGroupId").
		Joins("LEFT JOIN Farms ON FarmOnFarmGroups.FarmId = Farms.Id").
		Where("FarmGroups.Id = ? AND FarmGroups.ClientId = ?", id, clientId).
		Scan(&result).
		Error; err != nil {
		return nil, err
	}
	// Return the populated result
	return &result, nil
}

func (rp FarmGroupRepository) Update(request *models.FarmGroup) error {
	if err := rp.dbContext.Table("FarmGroups").Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IActivePondRepository interface {
	Create(activePond *models.ActivePond) (*models.ActivePond, error)
	TakeById(id int) (*models.ActivePond, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.ActivePond, error)
	Update(activePond *models.ActivePond) error
	GetListWithActive(farmId int) ([]*models.PondWithActive, error)
}

type activePondRepositoryImp struct {
	dbContext *gorm.DB
}

func NewActivePondRepository(db *gorm.DB) IActivePondRepository {
	return &activePondRepositoryImp{
		dbContext: db,
	}
}

func (rp activePondRepositoryImp) Create(request *models.ActivePond) (*models.ActivePond, error) {
	if err := rp.dbContext.Table(dbconst.TActivePond).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp activePondRepositoryImp) TakeById(id int) (*models.ActivePond, error) {
	var result *models.ActivePond
	if err := rp.dbContext.Table(dbconst.TActivePond).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Active Pond TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp activePondRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.ActivePond, error) {
	var result *models.ActivePond
	if err := rp.dbContext.Table(dbconst.TActivePond).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found active Pond FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp activePondRepositoryImp) Update(request *models.ActivePond) error {
	if err := rp.dbContext.Table(dbconst.TActivePond).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

func (rp activePondRepositoryImp) GetListWithActive(farmId int) ([]*models.PondWithActive, error) {
	var result []*models.PondWithActive

	rawSql := `SELECT
    p."Id" AS "Id",
    p."Code" AS "Code",
    p."Name" AS "Name",
    ap_latest."Id" AS "ActivePondId",
    COALESCE(ap_latest."StartDate" IS NOT NULL, false) AS "HasHistory"
FROM
    "Ponds" p
LEFT JOIN
    (
        SELECT
            ap_inner."Id",
            ap_inner."PondId",
            ap_inner."StartDate"
        FROM
            "ActivePonds" ap_inner
        WHERE
            ap_inner."DelFlag" IS NULL OR ap_inner."DelFlag" = false
        AND
            ap_inner."StartDate" = (
                SELECT
                    MAX(ap_inner2."StartDate")
                FROM
                    "ActivePonds" ap_inner2
                WHERE
                    ap_inner2."PondId" = ap_inner."PondId"
                AND
                    (ap_inner2."DelFlag" IS NULL OR ap_inner2."DelFlag" = false)
            )
    ) ap_latest ON p."Id" = ap_latest."PondId"
WHERE
    (p."DelFlag" IS NULL OR p."DelFlag" = false AND p."FarmId" = ?)
ORDER BY
    p."Id";`

	if err := rp.dbContext.Raw(rawSql, farmId).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

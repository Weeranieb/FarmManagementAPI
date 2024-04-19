package models

import (
	"boonmafarm/api/pkg/models/constants"
	"errors"
	"time"
)

// Activity represents an activity.
type Activity struct {
	Id             int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivePondId   int       `json:"activePondId" gorm:"column:ActivePondId"`
	ToActivePondId *int      `json:"toActivePondId" gorm:"column:ToActivePondId"`
	Mode           string    `json:"mode" gorm:"column:Mode"`
	MerchantId     *int      `json:"merchantId" gorm:"column:MerchantId"`
	Amount         *int      `json:"amount" gorm:"column:Amount"`
	FishType       *string   `json:"fishType" gorm:"column:FishType"`
	FishWeight     *float64  `json:"fishWeight" gorm:"column:FishWeight"`
	FishUnit       *string   `json:"fishUnit" gorm:"column:FishUnit"`
	PricePerUnit   *float64  `json:"pricePerUnit" gorm:"column:PricePerUnit"`
	ActivityDate   time.Time `json:"activityDate" gorm:"column:ActivityDate"`
	Base
}

type CreateActivityRequest struct {
	ActivePondId   int             `json:"activePondId" gorm:"column:ActivePondId"`
	ToActivePondId *int            `json:"toActivePondId,omitempty" gorm:"column:ToActivePondId"`
	Mode           string          `json:"mode" gorm:"column:Mode"`
	MerchantId     *int            `json:"merchantId,omitempty" gorm:"column:MerchantId"`
	Amount         *int            `json:"amount,omitempty" gorm:"column:Amount"`
	FishType       *string         `json:"fishType,omitempty" gorm:"column:FishType"`
	FishWeight     *float64        `json:"fishWeight,omitempty" gorm:"column:FishWeight"`
	PricePerUnit   *float64        `json:"pricePerUnit,omitempty" gorm:"column:PricePerUnit"`
	FishUnit       *string         `json:"fishUnit" gorm:"column:FishUnit"`
	ActivityDate   time.Time       `json:"activityDate" gorm:"column:ActivityDate"`
	SellDetail     []AddSellDetail `json:"sellDetails,omitempty"`
}

type CreateActivityResponse struct {
	Activity
	SellDetail []SellDetail `json:"sellDetails,omitempty"`
}

// Validation Add
func (a CreateActivityRequest) Validation() error {
	if a.ActivePondId == 0 {
		return errors.New(ErrActivePondIdEmpty)
	}
	if a.Mode == "" {
		return errors.New(ErrModeEmpty)
	}
	if a.ActivityDate.IsZero() {
		return errors.New(ErrActivityDateEmpty)
	}

	switch constants.ActivityType(a.Mode) {
	case constants.FillType:
		if a.Amount == nil {
			return errors.New(ErrAmountEmpty)
		}
		if a.FishType == nil {
			return errors.New(ErrFishTypeEmpty)
		}
		if a.FishWeight == nil {
			return errors.New(ErrFishWeightEmpty)
		}
		if a.FishUnit == nil && (*a.FishType == constants.Kilogram || *a.FishType == constants.Keed) {
			return errors.New(ErrFishUnitEmpty)
		}
		if a.PricePerUnit == nil {
			return errors.New(ErrPricePerUnitEmpty)
		}
	case constants.MoveType:
		if a.ToActivePondId == nil {
			return errors.New(ErrToActivePondIdEmpty)
		}
		if a.Amount == nil {
			return errors.New(ErrAmountEmpty)
		}
		if a.FishUnit == nil && (*a.FishType == constants.Kilogram || *a.FishType == constants.Keed) {
			return errors.New(ErrFishUnitEmpty)
		}
		if a.FishType == nil {
			return errors.New(ErrFishTypeEmpty)
		}
		if a.FishWeight == nil {
			return errors.New(ErrFishWeightEmpty)
		}
		if a.PricePerUnit == nil {
			return errors.New(ErrPricePerUnitEmpty)
		}
	case constants.SellType:
		if a.MerchantId == nil {
			return errors.New(ErrMerchantIdEmpty)
		}
		if len(a.SellDetail) == 0 {
			return errors.New(ErrSellDetailEmpty)
		}
		for _, sellDetail := range a.SellDetail {
			if err := sellDetail.Validation(); err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid mode")
	}

	return nil
}

// Transfer Add
func (a CreateActivityRequest) Transfer(activity *Activity, sellDetail *[]SellDetail) {
	activity.ActivePondId = a.ActivePondId
	activity.ToActivePondId = a.ToActivePondId
	activity.Mode = a.Mode
	activity.MerchantId = a.MerchantId
	activity.Amount = a.Amount
	activity.FishType = a.FishType
	activity.FishUnit = a.FishUnit
	activity.FishWeight = a.FishWeight
	activity.PricePerUnit = a.PricePerUnit
	activity.ActivityDate = a.ActivityDate

	tempSellDetail := make([]SellDetail, len(a.SellDetail))
	for i, sellDetail := range a.SellDetail {
		var temp SellDetail
		temp.Amount = sellDetail.Amount
		temp.FishType = sellDetail.FishType
		temp.FishUnit = sellDetail.FishUnit
		temp.PricePerUnit = sellDetail.PricePerUnit
		temp.Size = sellDetail.Size
		tempSellDetail[i] = temp
	}

	if len(tempSellDetail) > 0 {
		*sellDetail = tempSellDetail
	}
}

const (
	ErrActivePondIdEmpty   = "active pond id is empty"
	ErrModeEmpty           = "mode is empty"
	ErrActivityDateEmpty   = "activity date is empty"
	ErrToActivePondIdEmpty = "to active pond id is empty"
	ErrMerchantIdEmpty     = "merchant id is empty"
	ErrSellDetailEmpty     = "sell detail is empty"
)

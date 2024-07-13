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

type ActivityPage struct {
	Activity
	TotalWeight float64 `json:"totalWeight"`
	Unit        string  `json:"unit"`
	FarmName    string  `json:"farmName" gorm:"column:FarmName"`
	PondName    string  `json:"pondName" gorm:"column:PondName"`
}

type CreateFillActivityRequest struct {
	PondId       int       `json:"pondId" gorm:"column:PondId"`
	Amount       int       `json:"amount,omitempty" gorm:"column:Amount"`
	FishType     string    `json:"fishType,omitempty" gorm:"column:FishType"`
	FishWeight   float64   `json:"fishWeight,omitempty" gorm:"column:FishWeight"`
	PricePerUnit float64   `json:"pricePerUnit,omitempty" gorm:"column:PricePerUnit"`
	FishUnit     string    `json:"fishUnit" gorm:"column:FishUnit"`
	ActivityDate time.Time `json:"activityDate" gorm:"column:ActivityDate"`
	IsNewPond    bool      `json:"isNewPond,omitempty" gorm:"column:IsNewPond"`
}

type CreateMoveActivityRequest struct {
	PondId       int       `json:"pondId" gorm:"column:PondId"`
	ToPondId     int       `json:"toPondId" gorm:"column:ToPondId"`
	Amount       int       `json:"amount,omitempty" gorm:"column:Amount"`
	FishType     string    `json:"fishType,omitempty" gorm:"column:FishType"`
	FishWeight   float64   `json:"fishWeight,omitempty" gorm:"column:FishWeight"`
	PricePerUnit float64   `json:"pricePerUnit,omitempty" gorm:"column:PricePerUnit"`
	FishUnit     string    `json:"fishUnit" gorm:"column:FishUnit"`
	ActivityDate time.Time `json:"activityDate" gorm:"column:ActivityDate"`
	IsNewPond    bool      `json:"isNewPond,omitempty" gorm:"column:IsNewPond"`
	IsClose      bool      `json:"isClose,omitempty" gorm:"column:IsClose"`
}

type CreateSellActivityRequest struct {
	PondId       int          `json:"pondId" gorm:"column:PondId"`
	MerchantId   int          `json:"merchantId" gorm:"column:MerchantId"`
	ActivityDate time.Time    `json:"activityDate" gorm:"column:ActivityDate"`
	SellDetail   []SellDetail `json:"sellDetails,omitempty"`
	IsClose      bool         `json:"isClose,omitempty" gorm:"column:IsClose"`
}

type ActivityWithSellDetail struct {
	Activity
	SellDetail []SellDetail `json:"sellDetails,omitempty"`
}

// Validation Add
func (a CreateFillActivityRequest) Validation() error {
	if a.PondId == 0 {
		return errors.New(ErrActivePondIdEmpty)
	}
	if a.Amount == 0 {
		return errors.New(ErrAmountEmpty)
	}
	if a.FishUnit == "" {
		return errors.New(ErrFishUnitEmpty)
	}
	if a.FishType == "" {
		return errors.New(ErrFishTypeEmpty)
	}
	if a.FishWeight == 0 {
		return errors.New(ErrFishWeightEmpty)
	}
	if a.PricePerUnit == 0 {
		return errors.New(ErrPricePerUnitEmpty)
	}
	if a.ActivityDate.IsZero() {
		return errors.New(ErrActivityDateEmpty)
	}
	return nil
}

func (a CreateMoveActivityRequest) Validation() error {
	if a.PondId == 0 {
		return errors.New(ErrActivePondIdEmpty)
	}
	if a.ToPondId == 0 {
		return errors.New(ErrToPondIdEmpty)
	}
	if a.Amount == 0 {
		return errors.New(ErrAmountEmpty)
	}
	if a.FishUnit == "" {
		return errors.New(ErrFishUnitEmpty)
	}
	if a.FishType == "" {
		return errors.New(ErrFishTypeEmpty)
	}
	if a.FishWeight == 0 {
		return errors.New(ErrFishWeightEmpty)
	}
	if a.PricePerUnit == 0 {
		return errors.New(ErrPricePerUnitEmpty)
	}
	if a.ActivityDate.IsZero() {
		return errors.New(ErrActivityDateEmpty)
	}
	return nil
}

func (a CreateSellActivityRequest) Validation() error {
	if a.PondId == 0 {
		return errors.New(ErrActivePondIdEmpty)
	}
	if a.MerchantId == 0 {
		return errors.New(ErrMerchantIdEmpty)
	}
	if len(a.SellDetail) == 0 {
		return errors.New(ErrSellDetailEmpty)
	}
	if a.ActivityDate.IsZero() {
		return errors.New(ErrActivityDateEmpty)
	}
	return nil
}

// Transfer Add Fill
func (a CreateFillActivityRequest) Transfer(activity *Activity, activePondId int) {
	activity.Mode = string(constants.FillType)
	activity.ActivePondId = activePondId
	activity.Amount = &a.Amount
	activity.FishType = &a.FishType
	activity.FishWeight = &a.FishWeight
	activity.FishUnit = &a.FishUnit
	activity.PricePerUnit = &a.PricePerUnit
	activity.ActivityDate = a.ActivityDate
}

// Transfer Add Move
func (a CreateMoveActivityRequest) Transfer(activity *Activity, fromActivePondId int, toActivePondId int) {
	activity.Mode = string(constants.MoveType)
	activity.ActivePondId = fromActivePondId
	activity.ToActivePondId = &toActivePondId
	activity.Amount = &a.Amount
	activity.FishType = &a.FishType
	activity.FishWeight = &a.FishWeight
	activity.FishUnit = &a.FishUnit
	activity.PricePerUnit = &a.PricePerUnit
	activity.ActivityDate = a.ActivityDate
}

// Transfer Add Sell
func (a CreateSellActivityRequest) Transfer(activity *Activity, sellDetail *[]SellDetail, activePondId int) {
	activity.ActivePondId = activePondId
	activity.MerchantId = &a.MerchantId
	activity.ActivityDate = a.ActivityDate
	activity.Mode = string(constants.SellType)

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

// Transfer Add
// func (a CreateActivityRequest) Transfer(activity *Activity, sellDetail *[]SellDetail) {
// 	activity.ActivePondId = a.ActivePondId
// 	activity.ToActivePondId = a.ToActivePondId
// 	activity.Mode = a.Mode
// 	activity.MerchantId = a.MerchantId
// 	activity.Amount = a.Amount
// 	activity.FishType = a.FishType
// 	activity.FishUnit = a.FishUnit
// 	activity.FishWeight = a.FishWeight
// 	activity.PricePerUnit = a.PricePerUnit
// 	activity.ActivityDate = a.ActivityDate

// 	tempSellDetail := make([]SellDetail, len(a.SellDetail))
// 	for i, sellDetail := range a.SellDetail {
// 		var temp SellDetail
// 		temp.Amount = sellDetail.Amount
// 		temp.FishType = sellDetail.FishType
// 		temp.FishUnit = sellDetail.FishUnit
// 		temp.PricePerUnit = sellDetail.PricePerUnit
// 		temp.Size = sellDetail.Size
// 		tempSellDetail[i] = temp
// 	}

// 	if len(tempSellDetail) > 0 {
// 		*sellDetail = tempSellDetail
// 	}
// }

const (
	ErrActivePondIdEmpty   = "active pond id is empty"
	ErrModeEmpty           = "mode is empty"
	ErrActivityDateEmpty   = "activity date is empty"
	ErrToActivePondIdEmpty = "to active pond id is empty"
	ErrMerchantIdEmpty     = "merchant id is empty"
	ErrSellDetailEmpty     = "sell detail is empty"
	ErrToPondIdEmpty       = "to pond id is empty"
)

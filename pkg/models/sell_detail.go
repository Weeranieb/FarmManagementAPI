package models

import (
	"boonmafarm/api/pkg/models/constants"
	"errors"
)

// SellDetail represents a sell detail in the system.
type SellDetail struct {
	Id           int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	SellId       int     `json:"sellId" gorm:"column:SellId"`
	Size         string  `json:"size" gorm:"column:Size"`
	FishType     string  `json:"fishType" gorm:"column:FishType"`
	Amount       float64 `json:"amount" gorm:"column:Amount"`
	FishUnit     string  `json:"fishUnit" gorm:"column:FishUnit"`
	PricePerUnit float64 `json:"pricePerUnit" gorm:"column:PricePerUnit"`
	Base
}

type AddSellDetail struct {
	Size         string  `json:"size" gorm:"column:Size"`
	FishType     string  `json:"fishType" gorm:"column:FishType"`
	Amount       float64 `json:"amount" gorm:"column:Amount"`
	FishUnit     string  `json:"fishUnit" gorm:"column:FishUnit"`
	PricePerUnit float64 `json:"pricePerUnit" gorm:"column:PricePerUnit"`
}

// Validation Add
func (a AddSellDetail) Validation() error {
	if a.Size == "" {
		return errors.New(ErrSizeEmpty)
	}
	if a.FishType == "" {
		return errors.New(ErrFishTypeEmpty)
	}
	if a.Amount == 0 {
		return errors.New(ErrAmountEmpty)
	}
	if a.FishUnit == "" && (a.FishType == constants.Kilogram || a.FishType == constants.Keed) {
		return errors.New(ErrFishUnitEmpty)
	}
	if a.PricePerUnit == 0 {
		return errors.New(ErrPricePerUnitEmpty)
	}
	return nil
}

// transfer
func (a AddSellDetail) Transfer(sellDetail *SellDetail, sellId int) {
	sellDetail.SellId = sellId
	sellDetail.Amount = a.Amount
	sellDetail.Size = a.Size
	sellDetail.FishType = a.FishType
	sellDetail.FishUnit = a.FishUnit
	sellDetail.PricePerUnit = a.PricePerUnit
}

const (
	ErrSizeEmpty         = "size is empty"
	ErrFishTypeEmpty     = "fish type is empty"
	ErrAmountEmpty       = "amount is empty"
	ErrFishUnitEmpty     = "fish unit is empty"
	ErrPricePerUnitEmpty = "price per unit is empty"
	ErrFishWeightEmpty   = "fish weight is empty"
)

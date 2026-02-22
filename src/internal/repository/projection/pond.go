package projection

import (
	"time"

	"gorm.io/gorm"
)

// PondFillQueryRow is used to scan the single-query result (pond + farm.client_id + optional active_pond).
// Custom repo model for GetByIDWithFarmAndActivePond.
type PondFillQueryRow struct {
	PondId        int            `gorm:"column:pond_id"`
	PondFarmId    int            `gorm:"column:pond_farm_id"`
	PondName      string         `gorm:"column:pond_name"`
	PondStatus    string         `gorm:"column:pond_status"`
	PondDeletedAt gorm.DeletedAt `gorm:"column:pond_deleted_at"`
	PondCreatedAt time.Time      `gorm:"column:pond_created_at"`
	PondCreatedBy string         `gorm:"column:pond_created_by"`
	PondUpdatedAt time.Time      `gorm:"column:pond_updated_at"`
	PondUpdatedBy string         `gorm:"column:pond_updated_by"`
	ClientId      int            `gorm:"column:client_id"`
	ApId          *int           `gorm:"column:ap_id"`
	ApPondId      *int           `gorm:"column:ap_pond_id"`
	ApStartDate   *time.Time     `gorm:"column:ap_start_date"`
	ApEndDate     *time.Time     `gorm:"column:ap_end_date"`
	ApIsActive    *bool          `gorm:"column:ap_is_active"`
	ApTotalCost   *string        `gorm:"column:ap_total_cost"`
	ApTotalProfit *string        `gorm:"column:ap_total_profit"`
	ApNetResult   *string        `gorm:"column:ap_net_result"`
	ApTotalFish   *int           `gorm:"column:ap_total_fish"`
	ApCreatedAt   *time.Time     `gorm:"column:ap_created_at"`
	ApCreatedBy   *string        `gorm:"column:ap_created_by"`
	ApUpdatedAt   *time.Time     `gorm:"column:ap_updated_at"`
	ApUpdatedBy   *string        `gorm:"column:ap_updated_by"`
}

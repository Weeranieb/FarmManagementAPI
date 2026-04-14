package model

import "encoding/json"

type ScanLog struct {
	Id               int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivePondId     int             `json:"activePondId" gorm:"column:active_pond_id;not null"`
	Month            string          `json:"month" gorm:"column:month;not null"`
	ImagePaths       json.RawMessage `json:"imagePaths" gorm:"column:image_paths;type:jsonb;not null"`
	RawResponse      string          `json:"rawResponse" gorm:"column:raw_response;type:text"`
	ExtractedData    json.RawMessage `json:"extractedData" gorm:"column:extracted_data;type:jsonb"`
	ConfidenceScores json.RawMessage `json:"confidenceScores" gorm:"column:confidence_scores;type:jsonb"`
	Status           string          `json:"status" gorm:"column:status;not null;default:'pending_review'"`
	ReviewedBy       string          `json:"reviewedBy" gorm:"column:reviewed_by"`
	BaseModel
}

func (ScanLog) TableName() string {
	return "scan_logs"
}

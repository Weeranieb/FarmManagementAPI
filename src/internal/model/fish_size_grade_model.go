package model

type FishSizeGrade struct {
	Id        int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"column:name;not null;uniqueIndex"`
	SortIndex int    `json:"sortIndex" gorm:"column:sort_index;not null;index"`
	BaseModel
}

func (FishSizeGrade) TableName() string {
	return "fish_size_grades"
}

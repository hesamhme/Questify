package entities

type City struct {
	ID   int32  `gorm:"type:int;primaryKey"`
	Name string `gorm:"type:varchar(200);not null"`
}

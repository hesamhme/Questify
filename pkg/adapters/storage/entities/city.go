package entities

type City struct {
	ID   uint   `gorm:"type:int;primaryKey"`
	Name string `gorm:"type:varchar(200);not null"`
}

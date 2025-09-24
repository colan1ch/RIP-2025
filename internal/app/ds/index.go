package ds

type Index struct {
	ID          int `gorm:"primaryKey"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false"`
	Image       string `gorm:"type: varchar(128)"`
	Name        string `gorm:"type:varchar(25);not null"`
	Description string `gorm:"type: varchar(512)"`
	Cardinality int
	RowsCount int
	RecievedRows int
}
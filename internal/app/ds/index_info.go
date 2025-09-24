package ds

type IndexInfo struct {
	ID          int `gorm:"primaryKey"`
	Image       string
	Name        string `gorm:"type:varchar(25);not null"`

	Cardinality  int

	RowsCount int
	RecievedRows int
}
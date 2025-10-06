package ds

type IndexesQuery struct {
	ID        uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	QueryID uint `gorm:"not null;uniqueIndex:idx_order_chat"`
	IndexID    uint `gorm:"not null;uniqueIndex:idx_order_chat"`

	RowsCount int
	RecievedRows int

	
	Query Query `gorm:"foreignKey:QueryID"`
	Index    Index    `gorm:"foreignKey:IndexID"`
}

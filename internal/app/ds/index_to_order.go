package ds

type IndexesOrder struct {
	ID        uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	OrderID uint `gorm:"not null;uniqueIndex:idx_order_chat"`
	IndexID    uint `gorm:"not null;uniqueIndex:idx_order_chat"`

	RowsCount int
	RecievedRows int

	
	Order Order `gorm:"foreignKey:OrderID"`
	Index    Index    `gorm:"foreignKey:IndexID"`
}

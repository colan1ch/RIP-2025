package ds

type IndexesQuery struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	QueryID uint `gorm:"not null;uniqueIndex:idx_order_chat"`
	IndexID uint `gorm:"not null;uniqueIndex:idx_order_chat"`

	RowsCount    int
	RecievedRows int
	Cardinality  int

	Query Query `gorm:"foreignKey:QueryID"`
	Index Index `gorm:"foreignKey:IndexID"`
}

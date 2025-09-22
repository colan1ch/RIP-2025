package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	orders map[int]*Order
}

func NewRepository() (*Repository, error) {
    repo := &Repository{
        orders: map[int]*Order{
            1: {
                OrderDate: "16.09.25",
                ExecutionTime: 1092,
                OrderID: 1,
                IndexesParametrs: []IndexesToOrder{
                    {IndexId: 1, Cardinality: 130, RowsCount: 564, RecievedRows: 144},
                    {IndexId: 3, Cardinality: 22, RowsCount: 426, RecievedRows: 52},
                    {IndexId: 5, Cardinality: 54, RowsCount: 823, RecievedRows: 14},
                },
            },
        },
    }
    return repo, nil
}

type Index struct {
	ID          int
	Image       string
	Name        string
	Description string
}

type IndexesToOrder struct{
	IndexId int
	Cardinality int
	RowsCount int
	RecievedRows int
}

type Order struct{
	OrderID int
	OrderDate string
	IndexesParametrs []IndexesToOrder
	ExecutionTime float64
}

func (r *Repository) GetIndexes() ([]Index, error) {

	indexes := []Index{
		{
			ID:          1,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card1.jpeg",
			Name:        "Player(club_id)",
			Description: "Индекс для поиска игроков по\nклубу (частые запросы по\nигрокам определенного клуба)",
			// Cardinality: 130,
			// RowsCount: 564,
			// RecievedRows: 52,
		},
		{
			ID:          2,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card2.jpeg",
			Name:        "Player(position)",
			Description: "Индекс для поиска игроков по\nпозиции (фильтрация по\n амплуа)",
			// Cardinality: 86,
			// RowsCount: 1027,
			// RecievedRows: 33,
		},
		{
			ID:          3,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card3.jpeg",
			Name:        "Match(match_date)",
			Description: "Индекс для поиска матчей по\nдате (анализ матчей за\nпериод)",
			// Cardinality: 52,
			// RowsCount: 879,
			// RecievedRows: 645,
		},
		{
			ID:          4,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card4.jpeg",
			Name:        "Match\n(tournament_id)",
			Description: "Индекс для поиска матчей по\nтурниру (анализ матчей\nконкретного турнира)",
			// Cardinality: 36,
			// RowsCount: 377,
			// RecievedRows: 52,
		},
		{
			ID:          5,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card5.jpeg",
			Name:        "Player_stats(player_id)",
			Description: "Индекс для статистики\nигрока в матче",
			// Cardinality: 87,
			// RowsCount: 950,
			// RecievedRows: 13,
		},
		{
			ID:          6,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card6.jpeg",
			Name:        "Player_Stats(goals)",
			Description: "Индекс для поиска лучших\nбомбардиров (анализ по голам)",
			// Cardinality: 83,
			// RowsCount: 398,
			// RecievedRows: 5,
		},
		{
			ID:          7,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card7.jpeg",
			Name:        "Coach(club_id)",
			Description: "Индекс для поиска тренеров по клубу (поиск тренерского штаба)",
			// Cardinality: 52,
			// RowsCount: 444,
			// RecievedRows: 13,
		},
	}
	if len(indexes) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return indexes, nil
}

func (r *Repository) GetIndex(id int) (Index, error) {
	indexes, err := r.GetIndexes()
	if err != nil {
		return Index{}, err
	}

	for _, index := range indexes {
		if int(index.ID) == id {
			return index, nil
		}
	}
	return Index{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetIndexesByName(name string) ([]Index, error) {
	var indexes []Index
    var err error
	indexes, err = r.GetIndexes()
       if err != nil {
	       return nil, err
       }

       var result []Index
       for _, index := range indexes {
	       if strings.Contains(strings.ToLower(index.Name), strings.ToLower(name)) {
		       result = append(result, index)
	       }
       }

       return result, nil
}


func (r *Repository) GetOrder(id int) *Order {
	// orderIndexes := map[int]Order{
	// 	1: {
	// 		OrderDate: "16.09.25",
	// 		ExecutionTime: 1092,
	// 		OrderID: 1,
	// 		IndexesParametrs: []IndexesToOrder{ // Поля м-м
	// 			{
	// 				IndexId: 1,
	// 				Cardinality: 130,
	// 				RowsCount: 564,
	// 				RecievedRows: 144,
	// 			},
	// 			{
	// 				IndexId: 3,
	// 				Cardinality: 22,
	// 				RowsCount: 426,
	// 				RecievedRows: 52,
	// 			},
	// 			{
	// 				IndexId:      5,
	// 				Cardinality: 54,
	// 				RowsCount: 823,
	// 				RecievedRows: 14,
	// 			},
	// 		},
	// 	},
	// }
	// // orderIndexes[1].IndexesParametrs[0], orderIndexes[1].IndexesParametrs[1] = orderIndexes[1].IndexesParametrs[1], orderIndexes[1].IndexesParametrs[0]
	// return orderIndexes[id]
	return r.orders[id]
}


func (r *Repository) GetOrderIndexes(id int) []Index{
	order := r.GetOrder(id)

	var indexesInGroup []Index
	for _, _index := range order.IndexesParametrs {
		index, err := r.GetIndex(_index.IndexId)
		if err == nil {
			indexesInGroup = append(indexesInGroup, index)
		}
	}
	return indexesInGroup
}


func (r *Repository) GetOrderCount(id int) int {
	orderIndexes := r.GetOrder(id)
	return len(orderIndexes.IndexesParametrs)
}

func (r *Repository) GetOrderId() int {
  return 1
}

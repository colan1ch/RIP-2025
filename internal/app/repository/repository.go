package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	queries map[int]*Query
}

func NewRepository() (*Repository, error) {
    repo := &Repository{
        queries: map[int]*Query{
            1: {
                QueryDate: "16.09.25",
                ExecutionTime: 1092,
                QueryID: 1,
                IndexesParametrs: []IndexesToQuery{
                    {IndexId: 1, Cardinality: 130, RowsCount: 564, RecievedRows: 144, PositionInQuery: 1},
                    {IndexId: 3, Cardinality: 22, RowsCount: 426, RecievedRows: 52, PositionInQuery: 2},
                    {IndexId: 5, Cardinality: 54, RowsCount: 823, RecievedRows: 14, PositionInQuery: 3},
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
	Cardinality int
	RowsCount int
	RecievedRows int
	TableField string
	PositionInQuery int
}

type IndexesToQuery struct{
	IndexId int
	Cardinality int
	RowsCount int
	RecievedRows int
	TableField string
	PositionInQuery int
}

type Query struct{
	QueryID int
	QueryDate string
	IndexesParametrs []IndexesToQuery
	ExecutionTime float64
}

func (r *Repository) GetIndexes() ([]Index, error) {

	indexes := []Index{
		{
			ID:          1,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card1.jpeg",
			Name:        "Player(club_id)",
			Description: "Индекс для поиска игроков по\nклубу (частые запросы по\nигрокам определенного клуба)",
			Cardinality: 130,
			RowsCount: 564,
			RecievedRows: 52,
			PositionInQuery: 1,
		},
		{
			ID:          2,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card2.jpeg",
			Name:        "Player(position)",
			Description: "Индекс для поиска игроков по\nпозиции (фильтрация по\n амплуа)",
			Cardinality: 86,
			RowsCount: 1027,
			RecievedRows: 33,
		},
		{
			ID:          3,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card3.jpeg",
			Name:        "Match(match_date)",
			Description: "Индекс для поиска матчей по\nдате (анализ матчей за\nпериод)",
			Cardinality: 52,
			RowsCount: 879,
			RecievedRows: 645,
			PositionInQuery: 2,
		},
		{
			ID:          4,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card4.jpeg",
			Name:        "Match\n(tournament_id)",
			Description: "Индекс для поиска матчей по\nтурниру (анализ матчей\nконкретного турнира)",
			Cardinality: 36,
			RowsCount: 377,
			RecievedRows: 52,
		},
		{
			ID:          5,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card5.jpeg",
			Name:        "Player_stats(player_id)",
			Description: "Индекс для статистики\nигрока в матче",
			Cardinality: 87,
			RowsCount: 950,
			RecievedRows: 13,
			PositionInQuery: 3,
		},
		{
			ID:          6,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card6.jpeg",
			Name:        "Player_Stats(goals)",
			Description: "Индекс для поиска лучших\nбомбардиров (анализ по голам)",
			Cardinality: 83,
			RowsCount: 398,
			RecievedRows: 5,
		},
		{
			ID:          7,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/card7.jpeg",
			Name:        "Coach(club_id)",
			Description: "Индекс для поиска тренеров по клубу (поиск тренерского штаба)",
			Cardinality: 52,
			RowsCount: 444,
			RecievedRows: 13,
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


func (r *Repository) GetQuery(id int) *Query {
	// queryIndexes := map[int]Query{
	// 	1: {
	// 		QueryDate: "16.09.25",
	// 		ExecutionTime: 1092,
	// 		QueryID: 1,
	// 		IndexesParametrs: []IndexesToQuery{ // Поля м-м
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
	// // queryIndexes[1].IndexesParametrs[0], queryIndexes[1].IndexesParametrs[1] = queryIndexes[1].IndexesParametrs[1], queryIndexes[1].IndexesParametrs[0]
	// return queryIndexes[id]
	return r.queries[id]
}


func (r *Repository) GetQueryIndexes(id int) []Index{
	query := r.GetQuery(id)

	var indexesInGroup []Index
	for _, _index := range query.IndexesParametrs {
		index, err := r.GetIndex(_index.IndexId)
		if err == nil {
			indexesInGroup = append(indexesInGroup, index)
		}
	}
	return indexesInGroup
}


func (r *Repository) GetQueryCount(id int) int {
	queryIndexes := r.GetQuery(id)
	return len(queryIndexes.IndexesParametrs)
}

func (r *Repository) GetQueryId() int {
  return 1
}

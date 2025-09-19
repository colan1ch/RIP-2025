package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	// researchIndexes []Index
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Index struct {
	ID          int
	Image       string
	Name        string
	Description string
	Cardinality int
	RowsCount   int
	RecievedRows int
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


func (r *Repository) GetResearchIndexes(id int) []Index {
  requestIndexes := map[int][]int{1: {1, 3, 5}}
  var indexesInGroup []Index
  for _, indexID := range requestIndexes[id] {
    index, err := r.GetIndex(indexID)
    if err == nil {
      indexesInGroup = append(indexesInGroup, index)
    }
  }
  return indexesInGroup
}

func (r *Repository) GetResearchCount(id int) int {
	return len(r.GetResearchIndexes(id))
}

func (r *Repository) GetResearchId() int {
  return 1
}

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
	ExecutionTime float32
}

func (r *Repository) GetIndexes() ([]Index, error) {

	indexes := []Index{
		{
			ID:          1,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/b-tree.png",
			Name:        "B-Tree index",
			Description: "Это самый распространенный и стандартный тип индекса.\nОн отлично подходит для поиска по диапазону и точечного поиска.",
			Cardinality: 130,
			RowsCount: 564,
			ExecutionTime: 0,
		},
		{
			ID:          2,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/hash_index.png",
			Name:        "Hash index",
			Description: "Значительно быстрее B-tree для простых сравнений на равенство,\nно не поддерживает сортировку и диапазонные запросы.",
			Cardinality: 86,
			RowsCount: 1027,
			ExecutionTime: 0,
		},
		{
			ID:          3,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/gist_index.png",
			Name:        "GiST index",
			Description: "Геометрические данные (поиск пересечений, соседей),\nполнотекстовый поиск (вектор tsvector), диапазоны.",
			Cardinality: 52,
			RowsCount: 879,
			ExecutionTime: 0,
		},
		{
			ID:          4,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/sp-gist.png",
			Name:        "SP-GiST index",
			Description: "Эффективный поиск по данным, которые можно рекурсивно разделять\nна непересекающиеся области: IP-адреса, точки на плоскости, строки.",
			Cardinality: 36,
			RowsCount: 377,
			ExecutionTime: 0,
		},
		{
			ID:          5,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/gin_index.png",
			Name:        "GIN index",
			Description: "Идеален для массивов (поиск массива, содержащего элемент),\nJSONB (поиск по ключу или значению) и полнотекстового поиска.",
			Cardinality: 87,
			RowsCount: 950,
			ExecutionTime: 0,
		},
		{
			ID:          6,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/brin_index.png",
			Name:        "BRIN index",
			Description: "Хранит не каждое значение, а только минимальное и максимальное\nзначение для каждого блока (или диапазона) данных на диске.",
			Cardinality: 83,
			RowsCount: 398,
			ExecutionTime: 0,
		},
		{
			ID:          7,
			Image:       "http://127.0.0.1:9000/sqlanalyzer/unique_index.png",
			Name:        "Unique index",
			Description: "Особый вид индекса, который обеспечивает уникальность данных\nв проиндексированном столбце или группе столбцов.",
			Cardinality: 52,
			RowsCount: 444,
			ExecutionTime: 0,
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

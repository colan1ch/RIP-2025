package apitypes

import "LAB1/internal/app/ds"

type IndexesQueryJSON struct {
	ID           uint `json:"id"`
	QueryID      uint `json:"query_id"`
	IndexID      uint `json:"index_id"`
	RowsCount    int  `json:"rows_count"`
	RecievedRows int  `json:"recieved_rows"`
	Cardinality  int  `json:"cardinality"`
}

func IndexesQueryToJSON(indexesQuery ds.IndexesQuery) IndexesQueryJSON {
	return IndexesQueryJSON{
		ID:           indexesQuery.ID,
		QueryID:      indexesQuery.QueryID,
		IndexID:      indexesQuery.IndexID,
		RowsCount:    indexesQuery.RowsCount,
		RecievedRows: indexesQuery.RecievedRows,
		Cardinality:  indexesQuery.Cardinality,
	}
}

func IndexesQueryFromJSON(indexesQueryJSON IndexesQueryJSON) ds.IndexesQuery {
	return ds.IndexesQuery{
		// ID:        indexesQueryJSON.ID,
		// QueryID: indexesQueryJSON.QueryID,
		// IndexID:	indexesQueryJSON.IndexID,
		RowsCount: indexesQueryJSON.RowsCount,
		// RecievedRows: indexesQueryJSON.RecievedRows,
		Cardinality: indexesQueryJSON.Cardinality,
	}
}

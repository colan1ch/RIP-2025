package apitypes

import "LAB1/internal/app/ds"

type IndexJSON struct {
	ID          int `json:"id"`
	IsDelete    bool   `json:"is_delete"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func IndexToJSON(index ds.Index) IndexJSON {
	return IndexJSON{
		ID:          index.ID,
		IsDelete:    index.IsDelete,
		Image:       index.Image,
		Name:        index.Name,
		Description: index.Description,
	}
}

func IndexFromJSON(index IndexJSON) ds.Index {
	return ds.Index{
		ID:          index.ID,
		IsDelete:    index.IsDelete,
		Image:       index.Image,
		Name:        index.Name,
		Description: index.Description,
	}
}

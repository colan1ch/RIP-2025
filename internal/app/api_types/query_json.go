package apitypes

import (
	"LAB1/internal/app/ds"
	"time"
)


type QueryJSON struct {
	ID                int `json:"id"`
	DateQuery       string `json:"date_query"`
	Status          string `json:"status"`
	DateCreate   time.Time `json:"date_create"`
	DateForm    *time.Time `json:"date_form"`
	DateFinish  *time.Time `json:"date_finish"`
	CreatorLogin    string `json:"creator_login"`
	ModeratorLogin *string `json:"moderator_login"`
	ExecutionTime      int `json:"execution_time"`
}

func QueryToJSON(query ds.Query, creatorLogin string, moderatorLogin string) QueryJSON {
	var dateForm, dateFinish *time.Time
	if query.DateForm.Valid {
		dateForm = &query.DateForm.Time
	}

	if query.DateFinish.Valid {
		dateFinish = &query.DateFinish.Time
	}

	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}

	return QueryJSON{
		ID:                query.ID,
		DateQuery:    query.DateQuery,
		Status:          query.Status,
		DateCreate:   query.DateCreate,
		DateForm:    dateForm,
		DateFinish:  dateFinish,
		CreatorLogin:    creatorLogin,
		ModeratorLogin: mLogin,
		ExecutionTime:      query.ExecutionTime,
	}
}

func QueryFromJSON(query QueryJSON) ds.Query {
	if query.DateQuery == "" {
		return ds.Query{}
	}
	return ds.Query{
		DateQuery: query.DateQuery,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}

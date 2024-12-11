package storage

import (
	"net/http"
	"strconv"
)

type FilterQuery struct {
	Limit  int `json:"limit" validate:"gte=0,lte=20"`
	Offset int `json:"offset" validate:"gte=0"`
}

func (filterQuery *FilterQuery) Parse(r *http.Request) error {
	var limit int
	var offset int
	var err error
	query := r.URL.Query()

	if l := query.Get("limit"); l != "" {
		if limit, err = strconv.Atoi(l); err != nil {
			return err
		}

		filterQuery.Limit = limit
	}

	if o := query.Get("offset"); o != "" {
		if offset, err = strconv.Atoi(o); err != nil {
			return err
		}

		filterQuery.Offset = offset
	}

	return nil
}

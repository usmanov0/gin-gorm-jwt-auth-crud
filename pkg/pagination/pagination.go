package pagination

import "gorm.io/gorm"

type PaginateRes struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"currentPage"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	LastPage    int         `json:"lastPage"`
	PerPage     int         `json:"perPage"`
	Total       int64       `json:"total"`
}

type PaginationInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func Paginate(db *gorm.DB, page, limit int, rawFunc func(*gorm.DB) *gorm.DB, output interface{}) (PaginateRes, error) {
	offset := (page - 1) * limit

	query := db
	if rawFunc != nil {
		query = rawFunc(query)
	}
	var total int64
	query.Model(output).Count(&total)

	err := query.Offset(offset).Limit(limit).Find(output).Error
	if err != nil {
		return PaginateRes{}, nil
	}

	to := offset + limit
	if to > int(total) {
		to = int(total)
	}

	return PaginateRes{
		Data:        output,
		CurrentPage: page,
		From:        offset + 1,
		To:          to,
		LastPage:    (int(total) + limit - 1) / limit,
		PerPage:     limit,
		Total:       total,
	}, nil
}

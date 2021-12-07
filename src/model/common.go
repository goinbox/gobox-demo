package model

import (
	"github.com/goinbox/mysql"
)

func MakeInt64SliceSqlColQueryItem(name string, value []int64) *mysql.SqlColQueryItem {
	if len(value) <= 0 {
		return nil
	}

	if len(value) == 1 {
		return &mysql.SqlColQueryItem{
			Name:      name,
			Condition: mysql.SqlCondEqual,
			Value:     value[0],
		}
	}

	return &mysql.SqlColQueryItem{
		Name:      name,
		Condition: mysql.SqlCondIn,
		Value:     value,
	}
}

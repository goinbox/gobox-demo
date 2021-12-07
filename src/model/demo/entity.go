package demo

import "gdemo/model"

const (
	StatusOnline  = 1
	StatusOffline = 0

	ColumnID     = "id"
	ColumnName   = "name"
	ColumnStatus = "status"
)

type Entity struct {
	model.BaseEntity

	Name   string
	Status int
}

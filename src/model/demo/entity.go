package demo

import "gdemo/model"

type Entity struct {
	model.BaseEntity

	Name   string `mysql:"name"`
	Status int    `mysql:"status"`
}

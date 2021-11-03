package model

type BaseEntity struct {
	ID       int64  `mysql:"id"`
	AddTime  string `mysql:"add_time"`
	EditTime string `mysql:"edit_time"`
}

package entity

type SqlBaseEntity struct {
	Id       int64  `mysql:"id" json:"id" redis:"id"`
	AddTime  string `mysql:"add_time" json:"add_time" redis:"add_time"`
	EditTime string `mysql:"edit_time" json:"edit_time" redis:"edit_time"`
}

package model

type BaseEntity struct {
	ID       *int64 `column:"id"`
	AddTime  *string
	EditTime *string
}
